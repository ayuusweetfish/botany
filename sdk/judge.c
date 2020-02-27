#include "bot.h"

#include <errno.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

/*
    Usage:
    - ./judge <prog-1> <prog-2> <log-1> <log-2>

    Player protocol:
    - Startup: input "<side>"
        side - 0: first to move, 1: next to move
    - Repeat
        - Update: input "<row> <col>"
            row, col - 0..2
            for the first move, row = col = -1
        - Move: output "<row> <col>"
            row, col - 0..2
*/

static inline char *str_head(char *s, size_t n);

int main(int argc, char *argv[])
{
    /* Initialize players and make sure there are exactly 2 of them */
    int n = bot_player_init(argc, argv);
    if (n != 2) {
        fprintf(stderr, "Expected 2 players, got %d\n", n);
        exit(1);
    }

    /* Inform each player which side it is on */
    bot_player_send(0, "0");
    bot_player_send(1, "1");

    /* Set up board state */
    int move = 0;
    int win = -1, count = 0;
    int row = -1, col = -1;
    int board[3][3];
    memset(board, -1, sizeof board);

    /* Temporary variables and buffers */
    char buf[8];
    char *resp;
    size_t err;

    for (; win == -1 && count < 9; free(resp), move ^= 1) {
        /* Inform the current player of the last move */
        snprintf(buf, sizeof buf, "%d %d", row, col);
        bot_player_send(move, buf);
        resp = bot_player_recv(move, &err, 1000);

        /* An invalid response causes immediate defeat; the same below */
        if (resp == NULL) {
            fprintf(stderr, "Side #%d errors with %d (%s), considered resignation\n",
                move, (int)err, bot_strerr(err));
            win = move ^ 1;
            continue;
        }

        /* Parse the player's response */
        if (sscanf(resp, "%d%d", &row, &col) != 2 ||
            (row < 0 || row >= 3) ||
            (col < 0 || col >= 3))
        {
            fprintf(stderr, "Side #%d format incorrect (%s), considered resignation\n",
                move, str_head(resp, 10));
            win = move ^ 1;
            continue;
        }

        if (board[row][col] != -1) {
            fprintf(stderr, "Side #%d invalid move at (%d, %d), considered resignation\n",
                move, row, col);
            win = move ^ 1;
            continue;
        }

        fprintf(stderr, "Side #%d moves at (%d, %d)\n", move, row, col);
        board[row][col] = move;

        /* Check winning condition */
        for (int i = 0; i < 3; i++)
            if ((board[i][0] == move && board[i][1] == move && board[i][2] == move) ||
                (board[0][i] == move && board[1][i] == move && board[2][i] == move))
            {
                win = move;
                break;
            }
        if ((board[0][0] == move && board[1][1] == move && board[2][2] == move) ||
            (board[0][2] == move && board[1][1] == move && board[2][0] == move))
        {
            win = move;
        }
        if (win != -1) {
            fprintf(stderr, "Side #%d wins!\n", move);
        }
    }

    /* Write the report to stdout */
    printf("{\n  \"winner\": %d,\n  \"board\": \"", win);
    for (int i = 0; i < 3; i++) {
        for (int j = 0; j < 3; j++)
            putchar(".ox"[board[i][j] + 1]);
        printf("\\n");
    }
    printf("\"\n}\n");

    /* Stop all players to ensure their logs are written */
    bot_player_finish();

    return 0;
}

static inline char *str_head(char *s, size_t n)
{
    size_t m = strlen(s);
    if (m > n) {
        if (n >= 3) s[n - 3] = '.';
        if (n >= 2) s[n - 2] = '.';
        if (n >= 1) s[n - 1] = '.';
        s[n] = '\0';
    }
    return s;
}
