#include "judge.h"
#include "child.h"

#include <assert.h>
#include <errno.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/stat.h>
#include <sys/wait.h>
#include <unistd.h>

void compile(const char *sid, const char *contents)
{
    pid_t ch = fork();
    if (ch == 0) {
        child_enter_box();

        // Create submission folder
        char path[64];
        snprintf(path, sizeof path, "submissions/%s", sid);
        if (mkdir(path, 0755) != 0 && errno != EEXIST) {
            printf("mkdir(%s) failed: %s\n", path, strerror(errno));
            exit(1);
        }

        // Write code, literally write code
        snprintf(path, sizeof path, "submissions/%s/lang", sid);
        write_file(path, "cpp");
        snprintf(path, sizeof path, "submissions/%s/code", sid);
        write_file(path, contents);

        execl("./compile.sh", "./compile.sh", sid, NULL);
        exit(42);   // Unreachable
    } else {
        int wstatus;
        waitpid(ch, &wstatus, 0);
        // TODO: Check for failure
    }
}

bool is_compiled(const char *sid)
{
    pid_t ch = fork();
    if (ch == 0) {
        child_enter_box();

        char path[64];
        snprintf(path, sizeof path, "submissions/%s/bin", sid);
        exit(access(path, X_OK) == 0);
    } else {
        int wstatus;
        waitpid(ch, &wstatus, 0);
        assert(WIFEXITED(wstatus));
        int rc = WEXITSTATUS(wstatus);
        return (rc == 1);
    }
}
