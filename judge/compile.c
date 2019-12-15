#include "judge.h"

#include <errno.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/wait.h>
#include <unistd.h>

void compile(const char *sid, const char *contents)
{
    pid_t ch = fork();
    if (ch == 0) {
        usleep(1000000);
        if (chdir(judge_chroot) != 0) {
            printf("chdir() failed: %s\n", strerror(errno));
            exit(1);
        }
        if (chroot(judge_chroot) != 0) {
            printf("chroot() failed: %s\n", strerror(errno));
            exit(1);
        }
        uid_t uid = getuid();
        gid_t gid = getgid();
        if (setreuid(uid, uid) != 0 || setregid(gid, gid) != 0) {
            printf("setreuid/setregid() failed: %s\n", strerror(errno));
            exit(1);
        }
        execl("/bin/echo", "/bin/echo", "Hello", NULL);
        exit(0);
    } else {
        int wstatus;
        waitpid(ch, &wstatus, 0);
    }
}
