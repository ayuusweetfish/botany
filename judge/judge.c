#include "judge.h"

#include <hiredis/hiredis.h>
#include <blake2.h>

#include <assert.h>
#include <ctype.h>
#include <inttypes.h>
#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <time.h>
#include <unistd.h>

#define GROUP_NAME      "judge_group"
#define CLI_NAME_FMT    "judge_%d"

#define COMPILE_RESULT_LIST "compile_result"
#define MATCH_RESULT_LIST   "match_result"

#define COMPILE_STREAM  "compile"
#define MATCH_STREAM    "match"
#define SUBMISSION_HASH "submission"

#define WLOG(__fmt)         fprintf(stderr, "[%s] " __fmt "\n", wid)
#define WLOGF(__fmt, ...)   fprintf(stderr, "[%s] " __fmt "\n", wid, __VA_ARGS__)

const char *judge_chroot;

static redisContext *rctx;
static char wid[32];

static const char *sig_key = NULL;
static size_t sig_key_len;

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
    while ((c = getopt(argc, argv, "ha:p:i:d:k:")) != -1) {
        switch (c) {
        case 'h':
            printf("Usage: %s [-h] [-a redis_addr] [-p redis_port] "
                "[-i worker_id] [-d chroot_path] [-k sig_key]\n", argv[0]);
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
        case 'd':
            judge_chroot = optarg;
            break;
        case 'k':
            sig_key = optarg;
            break;
        }
    }
    if (redis_addr == NULL) redis_addr = "127.0.0.1";
    if (redis_port == 0) redis_port = 6379;
    if (worker_id == 0) worker_id = (int)getpid();
    if (judge_chroot == NULL) judge_chroot = "/";
    if (sig_key == NULL) sig_key = "aha";
    sig_key_len = strlen(sig_key);

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
    bool pending = true;

    while (1) {
        const char *start_id = (pending ? "0" : ">");
        reply = redisCommand(rctx, "XREADGROUP GROUP " GROUP_NAME
            " %s COUNT 1 BLOCK 1000 STREAMS "
            COMPILE_STREAM " " MATCH_STREAM " %s %s",
            wid, start_id, start_id);

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

        bool processed = false; // For pending check
        for (size_t i = 0; i < streamCount; i++) {
            reply = streams[i];

            assert(reply->type == REDIS_REPLY_ARRAY);
            assert(reply->elements == 2);
            assert(reply->element[0]->type == REDIS_REPLY_STRING);
            stream_name = reply->element[0]->str;
            reply = reply->element[1];

            assert(reply->type == REDIS_REPLY_ARRAY);
            if (pending && reply->elements == 0) continue;
            assert(reply->elements == 1);
            reply = reply->element[0];
            processed = true;

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

        if (pending && !processed) {
            WLOG("Pending tasks cleared");
            pending = false;
        }
    }

    return 0;
}

void retrieve_submission(const char *sid, char **lang, char **contents)
{
    // Generate signature
    const char *hex = "0123456789abcdef";

    char ts[32], sig[65];
    uint8_t digest[32];
    char cmd[256];

retry:
    snprintf(ts, sizeof ts, "%" PRId64, (int64_t)time(NULL));
    blake2b(digest, ts, sig_key, 32, strlen(ts), sig_key_len);
    for (int i = 0; i < 32; i++) {
        sig[i * 2 + 0] = hex[digest[i] >> 4];
        sig[i * 2 + 1] = hex[digest[i] & 15];
    }
    sig[64] = '\0';

    // Run HTTP request through the cURL binary
    snprintf(cmd, sizeof cmd,
        "curl \"http://localhost:3434/api/judge/%s?ts=%s&sig=%s\" 2>/dev/null", sid, ts, sig);

    FILE *fp = popen(cmd, "r");
    if (fp == NULL) {
        WLOG("Cannot retrieve submission, retrying");
        usleep(1000000);
        goto retry;
    }

    // Reuse `cmd`
    if (fgets(cmd, sizeof cmd, fp) == NULL) {
        WLOG("Incorrect response format, retrying");
        usleep(1000000);
        pclose(fp);
        goto retry;
    }
    // Extract content (code) length and language
    size_t i, j;
    size_t code_len = 0;
    for (i = 0; cmd[i] >= '0' && cmd[i] <= '9'; i++)
        code_len = code_len * 10 + cmd[i] - '0';
    while (isspace(cmd[i])) i++;    // Skip all whitespace
    j = strlen(cmd) - 1;
    while (j > i && isspace(cmd[j])) j--;
    cmd[j + 1] = '\0';
    *lang = strdup(cmd + i);

    char *buf = (char *)malloc(code_len + 1);
    i = 0;
    while (i < code_len) {
        if ((j = fread(buf + i, 1, code_len - i, fp)) == 0) {
            WLOG("Incorrect response format, retrying");
            usleep(1000000);
            pclose(fp);
            goto retry;
        }
        i += j;
    }
    buf[code_len] = '\0';
    *contents = buf;

    pclose(fp);
}

void process_compile(redisReply *kv)
{
    redisReply *reply;

    const char *sid = NULL;
    for (int i = 0; i + 1 < kv->elements; i += 2) {
        assert(kv->element[i]->type == REDIS_REPLY_STRING);
        assert(kv->element[i + 1]->type == REDIS_REPLY_STRING);
        if (strcmp(kv->element[i]->str, "sid") == 0) {
            sid = kv->element[i + 1]->str;
        }
    }
    assert(sid != NULL);

    char *lang, *contents;
    retrieve_submission(sid, &lang, &contents);

    // Update status
    WLOGF("Compiling: %s", sid);
    reply = redisCommand(rctx, "RPUSH " COMPILE_RESULT_LIST " %s 1 Compiling", sid);

    // Compilation work
    char *msg;
    int retcode = compile(sid, lang, contents, &msg);
    free(lang);
    free(contents);

    // Done!
    WLOGF("Done:      %s (exit code %d)", sid, retcode);
    // Convert to server's status code
    // 9: success, -1: compilation error
    int code = (retcode == 0 ? 9 : -1);
    reply = redisCommand(rctx, "RPUSH " COMPILE_RESULT_LIST " %s %d %s", sid, code, msg);
    free(msg);
}

void process_match(redisReply *kv)
{
    redisReply *reply;

    const char *mid = NULL;
    const char *judge = NULL;
    int num_parties = 0;
    for (int i = 0; i + 1 < kv->elements; i += 2) {
        assert(kv->element[i]->type == REDIS_REPLY_STRING);
        assert(kv->element[i + 1]->type == REDIS_REPLY_STRING);
        if (strcmp(kv->element[i]->str, "mid") == 0) {
            mid = kv->element[i + 1]->str;
        } else if (strcmp(kv->element[i]->str, "judge") == 0) {
            judge = kv->element[i + 1]->str;
        } else if (strcmp(kv->element[i]->str, "num_parties") == 0) {
            num_parties = (int)strtol(kv->element[i + 1]->str, NULL, 10);
        }
    }

    assert(mid != NULL && judge != NULL && num_parties != 0);

    char *parties[num_parties];
    for (int i = 0; i + 1 < kv->elements; i += 2) {
        if (memcmp(kv->element[i]->str, "party_", 6) == 0) {
            int index = (int)strtol(kv->element[i]->str + 6, NULL, 10);
            if (index >= 0 && index < num_parties)
                parties[index] = kv->element[i + 1]->str;
        }
    }

    // Update status
    WLOGF("Running:   %s", mid);
    for (int i = 0; i < num_parties; i++) {
        WLOGF("  Party #%d: %s", i, parties[i]);
    }

    reply = redisCommand(rctx, "RPUSH " MATCH_RESULT_LIST " %s 1 Compiling", mid);
    for (int i = 0; i < num_parties; i++) {
        if (!is_compiled(parties[i])) {
            const char *sid = parties[i];
            char *lang, *contents;
            retrieve_submission(sid, &lang, &contents);

            // TODO: Assert that compilation succeeds
            char *msg;
            compile(sid, lang, contents, &msg);
            free(lang);
            free(contents);
            free(msg);
        }
    }

    reply = redisCommand(rctx, "RPUSH " MATCH_RESULT_LIST " %s 2 Running", mid);

    // Match work
    char *msg;
    int retcode = match(mid, judge, num_parties, (const char **)parties, &msg);

    WLOGF("Done:      %s (exit code %d)", mid, retcode);
    reply = redisCommand(rctx, "RPUSH " MATCH_RESULT_LIST " %s 9 Done", mid);
}
