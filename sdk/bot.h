#ifndef _BOT_H_
#define _BOT_H_

#ifdef __cplusplus
extern "C" {
#endif

#include <unistd.h>

#define BOT_ERR_NONE    0
#define BOT_ERR_FMT     1
#define BOT_ERR_SYSCALL 2
#define BOT_ERR_TOOLONG 3
#define BOT_ERR_CLOSED  4
#define BOT_ERR_TIMEOUT 5

int bot_send_blob(int pipe, size_t len, const char *payload);

/*
  Receives data from a file descriptor with a given timeout.
  Returns a pointer to the data and stores the length in *o_len.
  In case of errors, returns NULL, stores the error code (see above) in *o_len,
  and prints related messages to stderr.
 */
char *bot_recv_blob(int pipe, size_t *o_len, int timeout);

const char *bot_strerr(int code);

void bot_send(const char *s);
char *bot_recv();

#ifdef __cplusplus
}
#endif

#endif
