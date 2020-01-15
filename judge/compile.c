#ifdef BOT_POSIX_COMPLIANT
#define pipe2(__pipe, __flags) pipe(__pipe)
#else
#define _GNU_SOURCE
#endif

#include "judge.h"
#include "child.h"

#include <assert.h>
#include <errno.h>
#include <fcntl.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/stat.h>
#include <sys/wait.h>
#include <unistd.h>

int compile(const char *sid, const char *lang, const char *contents, char **msg)
{
    int fd_pipe[2];
    if (pipe2(fd_pipe, O_NONBLOCK | O_CLOEXEC) != 0) {
        *msg = (char *)malloc(64);
        snprintf(*msg, 64, "pipe() failed: %s\n", strerror(errno));
        return 1;
    }

    pid_t ch = fork();
    if (ch == 0) {
        puts("child");
        dup2(fd_pipe[1], STDOUT_FILENO);
        close(fd_pipe[0]);

        child_enter_box();

        // Create submission folder
        char path[64];
        snprintf(path, sizeof path, "submissions/%s", sid);
        if (mkdir(path, 0755) != 0 && errno != EEXIST) {
            printf("mkdir(%s) failed: %s\n", path, strerror(errno));
            exit(1);
        }

        // Write code, literally write code
        snprintf(path, sizeof path, "submissions/%s/code.%s", sid, lang);
        write_file(path, contents);

        impose_rlimits();
        execl("./compile.sh", "./compile.sh", sid, lang, NULL);
        exit(42);   // Unreachable
    } else {
        puts("parent waiting");
        int wstatus;
        waitpid(ch, &wstatus, 0);
        puts("parent wait done");
        *msg = (char *)malloc(1024);
        ssize_t len = read(fd_pipe[0], *msg, 1023);
        puts("parent read done");
        (*msg)[len > 0 ? len : 0] = '\0';
        close(fd_pipe[0]);
        close(fd_pipe[1]);
        if (WIFEXITED(wstatus)) {
            return WEXITSTATUS(wstatus);
        } else if (WIFSIGNALED(wstatus)) {
            snprintf(*msg, 1024, "Compiler terminated by signal %s",
                strsignal(WTERMSIG(wstatus)));
            return -1;
        } else {
            strcpy(*msg, "Unknown internal anomalies");
            return -1;
        }
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
