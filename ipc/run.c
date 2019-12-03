#include "ipc.h"

#include <errno.h>
#include <signal.h>
#include <stdio.h>
#include <stdlib.h>

typedef struct _childproc {
    pid_t pid;
    /* fd_send is the child's stdin, fd_recv is stdout
       Parent writes to fd_send and reads from fd_recv */
    int fd_send, fd_recv;
} childproc;

#define pause_child(__cp)   kill((__cp).pid, SIGSTOP)
#define resume_child(__cp)  kill((__cp).pid, SIGCONT)

childproc create_child(const char *path)
{
    childproc ret;
    ret.pid = -1;

    int fd_send[2], fd_recv[2];
    if (pipe(fd_send) != 0 || pipe(fd_recv) != 0) {
        fprintf(stderr, "pipe() failed with errno %d\n", errno);
        return ret;
    }

    pid_t pid = fork();
    if (pid == -1) {
        fprintf(stderr, "fork() failed with errno %d\n", errno);
        return ret;
    }

    if (pid == 0) {
        /* Child process */
        dup2(fd_send[0], STDIN_FILENO);
        dup2(fd_recv[1], STDOUT_FILENO);
        close(fd_send[1]);
        close(fd_recv[0]);
        execl(path, path, NULL);
    } else {
        /* Parent process */
        close(fd_send[0]);
        close(fd_recv[1]);
        ret.pid = pid;
        ret.fd_send = fd_send[1];
        ret.fd_recv = fd_recv[0];
        pause_child(ret);
    }

    return ret;
}

void run_as_child()
{
    /* Receives a string, xor each character by 1 and return */
    size_t len, i;
    char *s = ipc_recv(STDIN_FILENO, &len, 1000);
    for (i = 0; i < len; i++) s[i] ^= 1;
    ipc_send(STDOUT_FILENO, len, s);

    free(s);
}

void run_as_parent(const char *path)
{
    /* Sends a string to the child and reads its response */
    childproc cp = create_child(path);
    resume_child(cp);

    ipc_send(cp.fd_send, 5, "ruwww");

    size_t len, i;
    char *s = ipc_recv(cp.fd_recv, &len, 1000);
    for (i = 0; i < len; i++) putchar(s[i]);

    free(s);
}

/* Example: ./a.out ./a.out */
int main(int argc, char *argv[])
{
    if (argc == 1)
        run_as_child();
    else
        run_as_parent(argv[1]);

    return 0;
}
