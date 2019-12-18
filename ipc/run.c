#include "ipc.h"

#include <errno.h>
#include <signal.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

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
        if (execl(path, path, NULL) != 0) {
            fprintf(stderr, "exec(%s) failed with errno %d\n", path, errno);
            exit(1);
        }
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

/* Modifies contents in `s` */
static inline char *str_head(char *s, size_t n)
{
    size_t m = strlen(s);
    if (m > n) {
        if (n >= 3) s[n - 3] = '.';
        if (n >= 2) s[n - 2] = '.';
        if (n >= 1) s[n - 1] = '.';
        s[n] = '\0';
    }
    return s;
}

/*
    Protocol:
    - Startup: I "<side>" side - 0: first to move, 1: next to move
    - Move: O "<row> <col>" row, col - 0..2
    - Response: I "<row> <col>" row, col - 0..2
*/

int main()
{
    childproc par[2];
    par[0] = create_child("./a.out");
    par[1] = create_child("./a.out");

    char buf[8];

    kill(par[0].pid, SIGCONT);
    ipc_send(par[0].fd_send, 0, "0");
    kill(par[0].pid, SIGSTOP);

    kill(par[1].pid, SIGCONT);
    ipc_send(par[1].fd_send, 0, "1");
    kill(par[1].pid, SIGSTOP);

    char *resp;
    size_t len;

    int move = 0;
    int win = -1;
    int row = -1, col = -1;
    int board[3][3];
    memset(board, -1, sizeof board);

    for (; win == -1; free(resp), move ^= 1) {
        snprintf(buf, sizeof buf, "%d %d", row, col);
        ipc_send(par[move].fd_send, 0, buf);

        kill(par[move].pid, SIGCONT);
        resp = ipc_recv(par[move].fd_recv, &len, 1000);
        kill(par[move].pid, SIGSTOP);

        if (resp == NULL) {
            fprintf(stderr, "Side #%d errors with %d, considered resignation\n",
                move, (int)len);
            win = move ^ 1;
            continue;
        }

        if (sscanf(resp, "%d%d", &row, &col) != 2 ||
            (row < 0 || row >= 3) ||
            (col < 0 || col >= 3))
        {
            fprintf(stderr, "Side #%d format incorrect (%s), considered resignation\n",
                move, str_head(resp, 10));
            win = move ^ 1;
            continue;
        }

        if (board[row][col] != -1) {
            fprintf(stderr, "Side #%d invalid move at (%d, %d), considered resignation\n",
                move, row, col);
            win = move ^ 1;
            continue;
        }

        fprintf(stderr, "Side #%d moves at (%d, %d)\n", move, row, col);
        board[row][col] = move;
    }

    return 0;
}
