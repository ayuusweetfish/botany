#ifndef _BOT_H_
#define _BOT_H_

#ifdef __cplusplus
extern "C" {
#endif

#include <unistd.h>

#define BOT_ERR_NONE    0   /* No error */
#define BOT_ERR_FMT     1   /* Incorrect message format, does not happen if player uses this library */
#define BOT_ERR_SYSCALL 2   /* Failure during system calls */
#define BOT_ERR_TOOLONG 3   /* Message too long */
#define BOT_ERR_CLOSED  4   /* Pipe closed, usually caused by program exiting */
#define BOT_ERR_TIMEOUT 5   /* Time out */

const char *bot_strerr(size_t code);

/* Judge side interfaces */

/* Initializes players from command line arguments */
int bot_player_init(int argc, char *const argv[]);

/* Terminates the player and flushes all output */
void bot_player_finish();

/* Sends to and receives from a player */
void bot_player_send(int id, const char *str);

/*
  Returns a string on success, and NULL on failure.
  The returned string, if non-null, should be free()'d.
  `o_len` holds the length on success, or an error code on failure.
  See constant definitions at top of the file, or use `bot_strerr()`
  to get the error description.
 */
char *bot_player_recv(int id, size_t *o_len, int timeout);

/* Player side interfaces */

/* Sends to stdout and receives from stdin */
void bot_send(const char *s);

/* Receives a string. The returned string should be free()'d. */
char *bot_recv();

#ifdef __cplusplus
}
#endif

#endif
