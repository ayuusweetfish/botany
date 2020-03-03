import bot
import random

bot.log('Hello from the snake! Side #' + bot.recv() + '\n')
random.seed()

used = [[False] * 3 for _ in range(3)]
while True:
    r, c = map(int, bot.recv().split())
    if r != -1: used[r][c] = True

    while True:
        r = random.randrange(0, 3)
        c = random.randrange(0, 3)
        if not used[r][c]: break

    used[r][c] = True
    bot.send('%d %d' % (r, c))
