#include "ipc.h"

#include <errno.h>
#include <fcntl.h>
#include <signal.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

typedef struct _childproc {
    pid_t pid;
    /* fd_send is the child's stdin, fd_recv is stdout
       Parent writes to fd_send and reads from fd_recv */
    int fd_send, fd_recv;
    int fd_log;
} childproc;

#define child_pause(__cp)   kill((__cp).pid, SIGSTOP)
#define child_resume(__cp)  kill((__cp).pid, SIGCONT)
#define child_kill(__cp)    kill((__cp).pid, SIGKILL)

childproc child_create(const char *path, const char *log)
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
        int fd_log = open(log, O_WRONLY | O_CREAT, 0644);
        if (fd_log == -1) {
            fprintf(stderr, "open(%s) failed with errno %d\n", log, errno);
            exit(1);
        }
        dup2(fd_send[0], STDIN_FILENO);
        dup2(fd_recv[1], STDOUT_FILENO);
        dup2(fd_log, STDERR_FILENO);
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
        child_pause(ret);
    }

    return ret;
}

void child_finish(childproc proc)
{
    fsync(proc.fd_log);
    child_kill(proc);
}

void child_send(childproc proc, const char *str)
{
    ipc_send(proc.fd_send, 0, str);
}

/* Returned value should be free()'d */
char *child_recv(childproc proc, size_t *o_len, int timeout)
{
    child_resume(proc);
    char *resp = ipc_recv(proc.fd_recv, o_len, timeout);
    child_pause(proc);
    return resp;
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
    Usage:
    - ./run <prog-1> <prog-2> <log-1> <log-2>

    Protocol:
    - Startup: I "<side>" side - 0: first to move, 1: next to move
    - Move: O "<row> <col>" row, col - 0..2
    - Response: I "<row> <col>" row, col - 0..2
*/

int main(int argc, char *argv[])
{
    if (argc < 5) return 1;

    childproc par[2];
    par[0] = child_create(argv[1], argv[3]);
    par[1] = child_create(argv[2], argv[4]);

    child_send(par[0], "0");
    child_send(par[1], "1");

    char buf[8];
    char *resp;
    size_t len;

    int move = 0;
    int win = -1, count = 0;
    int row = -1, col = -1;
    int board[3][3];
    memset(board, -1, sizeof board);

    for (; win == -1 && count < 9; free(resp), move ^= 1) {
        snprintf(buf, sizeof buf, "%d %d", row, col);
        child_send(par[move], buf);
        resp = child_recv(par[move], &len, 1000);

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

        // Check winning condition
        for (int i = 0; i < 3; i++)
            if ((board[i][0] == move && board[i][1] == move && board[i][2] == move) ||
                (board[0][i] == move && board[1][i] == move && board[2][i] == move))
            {
                win = move;
                break;
            }
        if ((board[0][0] == move && board[1][1] == move && board[2][2] == move) ||
            (board[0][2] == move && board[1][1] == move && board[2][0] == move))
        {
            win = move;
        }
        if (win != -1) {
            fprintf(stderr, "Side #%d wins!\n", move);
        }
    }

    printf("{\n  \"winner\": %d,\n  \"board\": \"", win);
    for (int i = 0; i < 3; i++) {
        for (int j = 0; j < 3; j++)
            putchar(".ox"[board[i][j] + 1]);
        printf("\\n");
    }
    printf("\"\n}\n");

    child_finish(par[0]);
    child_finish(par[1]);

    return 0;
}
