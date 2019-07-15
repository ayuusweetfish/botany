#include "ipc.h"

#include <errno.h>
#include <poll.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#define quq(__syscall, ...) _quq(#__syscall, __syscall(__VA_ARGS__))

static int _quq(const char *name, int ret)
{
    if (ret == -1) {
        fprintf(stderr, "%s() failed > < [errno %d]\n", name, errno);
        exit(1);
    }
    return ret;
}

void ipc_send(int pipe, size_t len, char *payload)
{
}

char *ipc_recv(int pipe, size_t *o_len, int timeout)
{
    struct pollfd pfd = (struct pollfd){pipe, POLLIN, 0};
    char *ret = NULL;
    size_t len = 0, ptr = 0;
    char buf[8];

    while (1) {
        pfd.revents = 0;
        int poll_ret = quq(poll, &pfd, 1, timeout);

        if (poll_ret == 1 && (pfd.revents & POLLIN)) {
            /* Got data */
            size_t size_to_read;
            if (ret == NULL) {
                size_to_read = 3;
            } else {
                size_to_read = len - ptr;
                if (size_to_read > sizeof buf)
                    size_to_read = sizeof buf;
            }
            ssize_t read_len = quq(read, pipe, buf, size_to_read);
            if (ret == NULL) {
                if (read_len < 3) {
                    /* Invalid */
                }
                /* Length */
                len = ((unsigned char)buf[2] << 16) |
                    ((unsigned char)buf[1] << 8) |
                    (unsigned char)buf[0];
                ret = (char *)malloc(len == 0 ? 1 : len);
            } else {
                if (read_len == 0) break;
                memcpy(ret + ptr, buf, read_len);
                ptr += read_len;
            }
        } else {
            puts("Unexpected");
            if (ret) free(ret);
            *o_len = 0;
            return NULL;
        }

    }

    *o_len = len;
    return ret;
}

int main()
{
    size_t len;
    char *s = ipc_recv(STDIN_FILENO, &len, 100000);
    puts(s);
    return 0;
}
