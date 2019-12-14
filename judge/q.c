#include <hiredis/hiredis.h>
#include <assert.h>
#include <stdio.h>
#include <string.h>
#include <unistd.h>

#define GROUP_NAME      "compile_group"
#define STREAM_NAME     "compile"
#define CLI_NAME_FMT    "compile_worker_%d"
#define CB_LIST_NAME    "compile_result"

#define WLOG(__fmt)         fprintf(stderr, "[%s] " __fmt "\n", wid)
#define WLOGF(__fmt, ...)   fprintf(stderr, "[%s] " __fmt "\n", wid, __VA_ARGS__)

int main()
{
    int pid = (int)getpid();
    char wid[32];
    snprintf(wid, sizeof wid, CLI_NAME_FMT, pid);

    redisContext *rctx = redisConnect("127.0.0.1", 6379);
    if (rctx == NULL) {
        WLOG("Cannot create redis context");
        return 1;
    } else if (rctx->err) {
        WLOGF("Redis connection error: %s", rctx->errstr);
        return 1;
    }

    WLOG("Starting");

    redisReply *reply;

    while (1) {
        reply = redisCommand(rctx, "XREADGROUP GROUP " GROUP_NAME
            " %s COUNT 1 BLOCK 1000 STREAMS " STREAM_NAME " >", wid);

        if (reply->type == REDIS_REPLY_NIL) continue;
        if (reply->type == REDIS_REPLY_ERROR) {
            WLOGF("Redis returned error: %s", reply->str);
            usleep(1000000);    // Prevent flooding logs
            continue;
        }

        // Parse reply
        assert(reply->type == REDIS_REPLY_ARRAY);
        assert(reply->elements == 1);
        reply = reply->element[0];

        assert(reply->type == REDIS_REPLY_ARRAY);
        assert(reply->elements == 2);
        reply = reply->element[1];

        assert(reply->type == REDIS_REPLY_ARRAY);
        assert(reply->elements == 1);
        reply = reply->element[0];

        assert(reply->type == REDIS_REPLY_ARRAY);
        assert(reply->elements == 2);

        assert(reply->element[0]->type == REDIS_REPLY_STRING);
        const char *redis_id = reply->element[0]->str;

        redisReply *kv = reply->element[1];
        assert(kv->type == REDIS_REPLY_ARRAY);

        const char *contents = NULL, *sid = NULL;
        for (int i = 0; i + 1 < kv->elements; i += 2) {
            assert(kv->element[i]->type == REDIS_REPLY_STRING);
            assert(kv->element[i + 1]->type == REDIS_REPLY_STRING);
            if (strcmp(kv->element[i]->str, "contents") == 0) {
                contents = kv->element[i + 1]->str;
            } else if (strcmp(kv->element[i]->str, "sid") == 0) {
                sid = kv->element[i + 1]->str;
            }
        }

        assert(contents != NULL && sid != NULL);

        // Update status
        WLOGF("Compiling: %s (%s)", sid, redis_id);
        reply = redisCommand(rctx, "RPUSH " CB_LIST_NAME " %s 1 Compiling", sid);

        // Compilation work
        usleep(1000000);

        // Done!
        WLOGF("Done:      %s (%s)", sid, redis_id);
        reply = redisCommand(rctx, "XACK " STREAM_NAME " " GROUP_NAME " %s", redis_id);
        static int cnt = 0;
        if (++cnt % 3 != 0) {
            // Success
            reply = redisCommand(rctx, "RPUSH " CB_LIST_NAME " %s 9 Done!", sid);
        } else {
            // Failure
            reply = redisCommand(rctx, "RPUSH " CB_LIST_NAME " %s -1 Compilation error", sid);
        }
    }

    return 0;
}
