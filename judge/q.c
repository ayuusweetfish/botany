#include <hiredis/hiredis.h>
#include <assert.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>

#define GROUP_NAME      "judge_group"
#define CLI_NAME_FMT    "judge_%d"

#define COMPILE_RESULT_LIST "compile_result"
#define MATCH_RESULT_LIST   "match_result"

#define COMPILE_STREAM  "compile"
#define MATCH_STREAM    "match"

#define WLOG(__fmt)         fprintf(stderr, "[%s] " __fmt "\n", wid)
#define WLOGF(__fmt, ...)   fprintf(stderr, "[%s] " __fmt "\n", wid, __VA_ARGS__)

static redisContext *rctx;
static char wid[32];

void process_compile(redisReply *kv);
void process_match(redisReply *kv);

int main(int argc, char *argv[])
{
    // Arguments
    const char *redis_addr = NULL;
    int redis_port = 0;
    int worker_id = 0;

    // Parse command line
    int c;
    while ((c = getopt(argc, argv, "ha:p:i:")) != -1) {
        switch (c) {
        case 'h':
            printf("Usage: %s [-h] [-a redis_addr] [-p redis_port] [-i worker_id]\n", argv[0]);
            exit(0);
        case 'a':
            redis_addr = optarg;
            break;
        case 'p':
            redis_port = (int)strtol(optarg, NULL, 0);
            break;
        case 'i':
            worker_id = (int)strtol(optarg, NULL, 0);
            break;
        }
    }
    if (redis_addr == NULL) redis_addr = "127.0.0.1";
    if (redis_port == 0) redis_port = 6379;
    if (worker_id == 0) worker_id = (int)getpid();

    // Get to work
    snprintf(wid, sizeof wid, CLI_NAME_FMT, worker_id);

retry:
    WLOGF("Connecting to Redis at %s:%d", redis_addr, redis_port);

    rctx = redisConnect(redis_addr, redis_port);
    if (rctx == NULL) {
        WLOG("Cannot create Redis context");
        usleep(1000000);
        goto retry;
    } else if (rctx->err) {
        WLOGF("Redis connection error: %s", rctx->errstr);
        usleep(1000000);
        goto retry;
    }

    WLOG("Starting");

    redisReply *reply;

    while (1) {
        reply = redisCommand(rctx, "XREADGROUP GROUP " GROUP_NAME
            " %s COUNT 1 BLOCK 1000 STREAMS " COMPILE_STREAM " " MATCH_STREAM " > >", wid);

        if (reply == NULL) {
            WLOG("Connection broken, retrying");
            usleep(1000000);
            goto retry;
        }
        if (reply->type == REDIS_REPLY_NIL) continue;
        if (reply->type == REDIS_REPLY_ERROR) {
            WLOGF("Redis returned error: %s", reply->str);
            usleep(1000000);    // Prevent flooding logs
            continue;
        }

        // Parse reply
        char *stream_name;

        assert(reply->type == REDIS_REPLY_ARRAY);
        size_t streamCount = reply->elements;
        redisReply **streams = reply->element;

        for (size_t i = 0; i < streamCount; i++) {
            reply = streams[i];

            assert(reply->type == REDIS_REPLY_ARRAY);
            assert(reply->elements == 2);
            assert(reply->element[0]->type == REDIS_REPLY_STRING);
            stream_name = reply->element[0]->str;
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

            if (strcmp(stream_name, COMPILE_STREAM) == 0) {
                process_compile(kv);
            } else if (strcmp(stream_name, MATCH_STREAM) == 0) {
                process_match(kv);
            }

            reply = redisCommand(rctx,
                "XACK %s " GROUP_NAME " %s", stream_name, redis_id);
        }
    }

    return 0;
}

void process_compile(redisReply *kv)
{
    redisReply *reply;

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
    WLOGF("Compiling: %s", sid);
    reply = redisCommand(rctx, "RPUSH " COMPILE_RESULT_LIST " %s 1 Compiling", sid);

    // Compilation work
    usleep(1000000);

    // Done!
    WLOGF("Done:      %s", sid);
    static int cnt = 0;
    if (++cnt % 3 != 0) {
        // Success
        reply = redisCommand(rctx, "RPUSH " COMPILE_RESULT_LIST " %s 9 Done!", sid);
    } else {
        // Failure
        reply = redisCommand(rctx, "RPUSH " COMPILE_RESULT_LIST " %s -1 Compilation error", sid);
    }
}

void process_match(redisReply *kv)
{
    WLOG("Running:   -");
    usleep(1000000);
    WLOG("Done:      -");
}