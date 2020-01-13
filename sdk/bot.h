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

/* General interfaces; usually not needed */

int bot_send_blob(int pipe, size_t len, const char *payload);

/*
  Receives data from a file descriptor with a given timeout.
  Returns a pointer to the data and stores the length in *o_len.
  In case of errors, returns NULL, stores the error code (see above) in *o_len,
  and prints related messages to stderr.
 */
char *bot_recv_blob(int pipe, size_t *o_len, int timeout);

const char *bot_strerr(int code);

/* Note: all strings returned by *_recv() need to be free()'d */
/* Judge side interfaces */

typedef struct _childproc {
    pid_t pid;
    /* fd_send is the child's stdin, fd_recv is stdout
       Parent writes to fd_send and reads from fd_recv */
    int fd_send, fd_recv;
    int fd_log;
} childproc;

/* Creates the child and pauses it */
childproc child_create(const char *cmd, const char *log);
/* Terminates the child and flushes all output */
void child_finish(childproc proc);
/* Pauses the child */
#define child_pause(__cp)   kill((__cp).pid, SIGSTOP)
/* Resumes the child */
#define child_resume(__cp)  kill((__cp).pid, SIGCONT)

/* Sends to and receives from a child */
void child_send(childproc proc, const char *str);
char *child_recv(childproc proc, size_t *o_len, int timeout);

/* Player side interfaces */

/* Sends to stdout and receives from stdin */
void bot_send(const char *s);
char *bot_recv();

#ifdef __cplusplus
}
#endif

#endif
