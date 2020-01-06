#ifndef __BOTANY_CHILD_H__
#define __BOTANY_CHILD_H__

#include <errno.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/resource.h>
#include <sys/time.h>
#include <unistd.h>

static inline void child_enter_box()
{
    if (chdir(judge_chroot) != 0) {
        printf("chdir(%s) failed: %s\n", judge_chroot, strerror(errno));
        exit(1);
    }
    if (chroot(judge_chroot) != 0) {
        printf("chroot(%s) failed: %s\n", judge_chroot, strerror(errno));
        exit(1);
    }
    uid_t uid = getuid();
    gid_t gid = getgid();
    if (setreuid(uid, uid) != 0 || setregid(gid, gid) != 0) {
        printf("setreuid/setregid() failed: %s\n", strerror(errno));
        exit(1);
    }

    // Change to working directory
    if (chdir("/var/botany") != 0) {
        printf("chdir(%s) failed: %s\n", "/var/botany", strerror(errno));
        exit(1);
    }
}

static inline void write_file(const char *path, const char *contents)
{
    FILE *fp = fopen(path, "w");
    if (fp == NULL) {
        printf("Cannot open file %s for writing: %s\n", path, strerror(errno));
        exit(1);
    }
    if (fwrite(contents, strlen(contents), 1, fp) != 1) {
        printf("Writing to file %s failed: %s", path, strerror(errno));
        exit(1);
    }
    fclose(fp);
}

static inline void impose_rlimits()
{
    struct rlimit l;
    // Memory limit (address space)
    l.rlim_cur = 1 << 30;
    l.rlim_max = 1 << 30;
    setrlimit(RLIMIT_AS, &l);
    // Memory limit (data segment)
    l.rlim_cur = 1 << 30;
    l.rlim_max = 1 << 30;
    setrlimit(RLIMIT_DATA, &l);
    // CPU limit
    l.rlim_cur = 15;
    l.rlim_max = 20;
    setrlimit(RLIMIT_CPU, &l);
    // Wall clock limit
    // Compilers usually do not override SIGALRM behaviour
    struct itimerval t;
    t.it_interval.tv_sec = 20;
    t.it_interval.tv_usec = 0;
    t.it_value = t.it_interval;
    setitimer(ITIMER_REAL, &t, NULL);
}

#endif
