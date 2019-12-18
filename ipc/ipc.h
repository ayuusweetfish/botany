#ifndef _IPC_H_
#define _IPC_H_

#include <unistd.h>

#define IPC_ERR_NONE    0
#define IPC_ERR_FMT     1
#define IPC_ERR_SYSCALL 2
#define IPC_ERR_TOOLONG 3
#define IPC_ERR_CLOSED  4
#define IPC_ERR_TIMEOUT 5

int ipc_send(int pipe, size_t len, const char *payload);

/*
  Receives data from a file descriptor with a given timeout.
  Returns a pointer to the data and stores the length in *o_len.
  In case of errors, returns NULL, stores the error code (see above) in *o_len,
  and prints related messages to stderr.
 */
char *ipc_recv(int pipe, size_t *o_len, int timeout);

#endif
