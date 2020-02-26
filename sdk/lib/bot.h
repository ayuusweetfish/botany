#ifndef _BOT_H_
#define _BOT_H_

#ifdef __cplusplus
extern "C" {
#endif

#include <unistd.h>

#define BOT_ERR_NONE    0   /* No error */
#define BOT_ERR_FMT     1   /* Incorrect message format, does not happen if player uses this library */
#define BOT_ERR_SYSCALL 2   /* Failure during system calls */
#define BOT_ERR_TOOLONG 3   /* Message too long */
#define BOT_ERR_CLOSED  4   /* Pipe closed, usually caused by program exiting */
#define BOT_ERR_TIMEOUT 5   /* Time out */

const char *bot_strerr(size_t code);

/* Judge side interfaces */

typedef struct _bot_player {
    pid_t pid;
    /* fd_send is the child's stdin, fd_recv is stdout
       Parent writes to fd_send and reads from fd_recv */
    int fd_send, fd_recv;
    int fd_log;
} bot_player;

/* Creates all children from command line arguments */
bot_player *bot_player_all(int argc, char *const argv[], int *num);
/* Terminates the child and flushes all output */
void bot_player_finish(bot_player *procs, int num);
/* Sends to and receives from a child */
void bot_player_send(bot_player proc, const char *str);
/*
  Returns a string on success, and NULL on failure.
  The returned string, if non-null, should be free()'d.
  `o_len` holds the length on success, or an error code on failure.
  See constant definitions at top of the file, or use `bot_strerr()`
  to get the error description.
 */
char *bot_player_recv(bot_player proc, size_t *o_len, int timeout);

/* Player side interfaces */

/* Sends to stdout and receives from stdin */
void bot_send(const char *s);
char *bot_recv();   /* Returned string should be free()'d */

#ifdef __cplusplus
}
#endif

#endif
