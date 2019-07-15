#ifndef _IPC_H_
#define _IPC_H_

#include <unistd.h>

void ipc_send(int pipe, size_t len, char *payload);
char *ipc_recv(int pipe, size_t *o_len, int timeout);

#endif
