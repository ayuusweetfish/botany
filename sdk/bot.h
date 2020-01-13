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

/* Judge side interfaces */

typedef struct _childproc {
    pid_t pid;
    /* fd_send is the child's stdin, fd_recv is stdout
       Parent writes to fd_send and reads from fd_recv */
    int fd_send, fd_recv;
    int fd_log;
} childproc;

/*
  Creates the child.
  Child processes are normally paused, but during `child_recv()`
  the process is resumed, and paused again after its response arrives.
 */
childproc child_create(const char *cmd, const char *log);
/* Terminates the child and flushes all output */
void child_finish(childproc proc);
/* Sends to and receives from a child */
void child_send(childproc proc, const char *str);
/*
  Returns a string on success, and NULL on failure.
  The returned string, if non-null, should be free()'d.
  `o_len` holds the length on success, or an error code on failure.
  See constant definitions at top of the file, or use `bot_strerr()`
  to get the error description.
 */
char *child_recv(childproc proc, size_t *o_len, int timeout);

/* Player side interfaces */

/* Sends to stdout and receives from stdin */
void bot_send(const char *s);
char *bot_recv();   /* Returned string should be free()'d */

#ifdef __cplusplus
}
#endif

#endif
