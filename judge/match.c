#include "judge.h"
#include "child.h"

#include <stdio.h>
#include <string.h>
#include <sys/wait.h>

int match(const char *mid, const char *judge, int num_parties, const char *parties[], char **msg, char ***logs)
{
    int pipe_stdout[2];
    if (pipe(pipe_stdout) != 0) {
        *msg = (char *)malloc(64);
        snprintf(*msg, 64, "pipe() failed: %s\n", strerror(errno));
        *logs = NULL;
        return 1;
    }

    pid_t ch = fork();
    if (ch == 0) {
        dup2(pipe_stdout[1], STDOUT_FILENO);
        close(pipe_stdout[0]);

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
        *msg = (char *)malloc(MATCH_REPORT_LEN);
        ssize_t len = read(pipe_stdout[0], *msg, MATCH_REPORT_LEN - 1);
        (*msg)[len > 0 ? len : 0] = '\0';
        close(pipe_stdout[0]);
        close(pipe_stdout[1]);
        if (WIFEXITED(wstatus)) {
            // Read participant logs
            *logs = (char **)malloc(sizeof(char *) * num_parties);
            char log_path[1024];
            for (int i = 0; i < num_parties; i++) {
                (*logs)[i] = (char *)malloc(MATCH_LOG_LEN);
                size_t len = 0;

                snprintf(log_path, sizeof log_path, "%s/var/botany/matches/%s/%d.log",
                    judge_chroot, mid, i);
                FILE *f = fopen(log_path, "r");
                if (f != NULL) {
                    len = fread((*logs)[i], 1, MATCH_LOG_LEN - 1, f);
                    fclose(f);
                }
                (*logs)[i][len] = '\0';
            }

            return WEXITSTATUS(wstatus);
        } else if (WIFSIGNALED(wstatus)) {
            snprintf(*msg, MATCH_REPORT_LEN, "Match terminated by signal %s",
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
