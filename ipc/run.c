#include "ipc.h"

#include <errno.h>
#include <signal.h>
#include <stdio.h>

typedef struct _childproc {
    pid_t pid;
    int fd_send, fd_recv;
} childproc;

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
        kill(pid, SIGSTOP);
    }

    return ret;
}

int main()
{
    childproc cp = create_child("/bin/cat");
    kill(cp.pid, SIGCONT);

    write(cp.fd_send, "rua", 3);
    write(cp.fd_send, "ruwww", 5);

    char buf[10] = { 0 };
    read(cp.fd_recv, buf, sizeof buf);
    puts(buf);

    return 0;
}
