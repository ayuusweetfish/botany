#include "judge.h"
#include "child.h"

#include <stdio.h>
#include <string.h>
#include <sys/wait.h>

int match(const char *mid, int num_parties, const char *parties[], char **msg)
{
    pid_t ch = fork();
    if (ch == 0) {
        child_enter_box();

        const char *argv[num_parties + 3];
        argv[0] = "./match.sh";
        argv[1] = mid;
        for (int i = 0; i < num_parties; i++)
            argv[i + 2] = parties[i];
        argv[num_parties + 2] = NULL;

        execv("./match.sh", (char *const *)argv);
        exit(42);   // Unreachable
    } else {
        int wstatus;
        waitpid(ch, &wstatus, 0);
        if (WIFEXITED(wstatus)) {
            return WEXITSTATUS(wstatus);
        } else if (WIFSIGNALED(wstatus)) {
            snprintf(*msg, 1024, "Match terminated by signal %s",
                strsignal(WTERMSIG(wstatus)));
            return -1;
        } else {
            strcpy(*msg, "Unknown internal anomalies");
            return -1;
        }
    }
}
