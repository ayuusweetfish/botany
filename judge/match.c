#include "judge.h"
#include "child.h"

#include <string.h>

void match(const char *mid, int num_parties, const char *parties[])
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
        // TODO: Check for failure
    }
}
