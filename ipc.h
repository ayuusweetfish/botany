#ifndef _IPC_H_
#define _IPC_H_

#include <unistd.h>

#define IPC_ERR_NONE    0
#define IPC_ERR_FMT     1
#define IPC_ERR_POLL    2
#define IPC_ERR_SYSCALL 3

void ipc_send(int pipe, size_t len, char *payload);

/*
  Receives data from a file descriptor with a given timeout.
  Returns a pointer to the data and stores the length in *o_len.
  In case of errors, returns NULL and stores the error code (see above) in *o_len.
 */
char *ipc_recv(int pipe, size_t *o_len, int timeout);

#endif
