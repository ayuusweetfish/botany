#ifndef __BOTANY_JUDGE_H__
#define __BOTANY_JUDGE_H__

#include <stdbool.h>

extern const char *judge_chroot;

int compile(const char *sid, const char *lang, const char *contents, char **msg);
bool is_compiled(const char *sid);
int match(const char *mid, int num_parties, const char *parties[], char **msg);

#endif
