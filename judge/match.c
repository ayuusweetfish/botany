#include "judge.h"
#include "child.h"

#include <stdio.h>
#include <string.h>
#include <sys/wait.h>

int match(const char *mid, const char *judge, int num_parties, const char *parties[], char **msg, char ***logs)
{
    int fd_pipe[2];
    if (pipe(fd_pipe) != 0) {
        *msg = (char *)malloc(64);
        snprintf(*msg, 64, "pipe() failed: %s\n", strerror(errno));
        return 1;
    }

    pid_t ch = fork();
    if (ch == 0) {
        dup2(fd_pipe[1], STDOUT_FILENO);
        close(fd_pipe[0]);

        child_enter_box();

        const char *argv[num_parties + 4];
        argv[0] = "./match.sh";
        argv[1] = mid;
        argv[2] = judge;
        for (int i = 0; i < num_parties; i++)
            argv[i + 3] = parties[i];
        argv[num_parties + 3] = NULL;

        execv("./match.sh", (char *const *)argv);
        exit(42);   // Unreachable
    } else {
        int wstatus;
        waitpid(ch, &wstatus, 0);
        *msg = (char *)malloc(1024);
        ssize_t len = read(fd_pipe[0], *msg, 1023);
        (*msg)[len > 0 ? len : 0] = '\0';
        close(fd_pipe[0]);
        close(fd_pipe[1]);
        if (WIFEXITED(wstatus)) {
            // Read participant logs
            *logs = (char **)malloc(sizeof(char *) * num_parties);
            for (int i = 0; i < num_parties; i++) {
                (*logs)[i] = (char *)malloc(65536);
                snprintf((*logs)[i], 65536, "Log from party %d\n", i + 1);
            }

            return WEXITSTATUS(wstatus);
        } else if (WIFSIGNALED(wstatus)) {
            snprintf(*msg, 1024, "Match terminated by signal %s",
                strsignal(WTERMSIG(wstatus)));
            *logs = NULL;
            return -1;
        } else {
            strcpy(*msg, "Unknown internal anomalies");
            *logs = NULL;
            return -1;
        }
    }
}
