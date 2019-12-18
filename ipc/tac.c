#include "ipc.h"

#include <stdbool.h>
#include <stdlib.h>
#include <stdio.h>
#include <time.h>

int main()
{
    char *s = ipc_recv_str();
    int side;
    sscanf(s, "%d", &side);
    free(s);
    fprintf(stderr, "Hello from side #%d\n", side);

    srand(((unsigned)time(NULL) << 1) | side);
    bool board[3][3] = {{ false }};

    while (1) {
        // Board state change
        int row, col;
        s = ipc_recv_str();
        sscanf(s, "%d%d", &row, &col);
        free(s);
        if (row != -1) board[row][col] = true;

        // Pick a random cell
        int u, v;
        do {
            u = rand() % 3;
            v = rand() % 3;
        } while (board[u][v]);
        board[u][v] = true;

        // Send
        char t[8];
        sprintf(t, "%d %d", u, v);
        ipc_send_str(t);
    }

    return 0;
}
