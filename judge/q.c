#include <hiredis/hiredis.h>
#include <assert.h>
#include <stdio.h>
#include <string.h>

#define GROUP_NAME      "compile_group"
#define STREAM_NAME     "compile"
#define CLI_NAME_PFX    "compile_worker_"

int main()
{
    redisContext *rctx = redisConnect("127.0.0.1", 6379);
    if (rctx == NULL) {
        puts("Cannot create redis context");
        return 1;
    } else if (rctx->err) {
        printf("Error: %s\n", rctx->errstr);
        return 1;
    }

    redisReply *reply;

    while (1) {
        reply = redisCommand(rctx, "XREADGROUP GROUP " GROUP_NAME
            " " CLI_NAME_PFX "_3 COUNT 1 BLOCK 1000 STREAMS " STREAM_NAME " >");

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
        printf("Redis ID: %s\n", reply->element[0]->str);

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

        printf("SID: %s\nContents:\n%s\n====\n", sid, contents);

        break;
    }

    return 0;
}
