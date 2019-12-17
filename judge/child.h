#ifndef __BOTANY_CHILD_H__
#define __BOTANY_CHILD_H__

#include <errno.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
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

#endif
