#ifndef __BOTANY_JUDGE_H__
#define __BOTANY_JUDGE_H__

#include <stdbool.h>

extern const char *judge_chroot;

int compile(const char *sid, const char *lang, const char *contents, char **msg);
bool is_compiled(const char *sid);
int match(const char *mid, const char *judge, int num_parties, const char *parties[], char **msg, char ***logs);

#define COMPILE_MSG_LEN     4096
#define MATCH_REPORT_LEN    65536
#define MATCH_LOG_LEN       65536

#endif
