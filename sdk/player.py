import sys
import random

class Bot:
    @staticmethod
    def send(s):
        l = len(s)
        if l >= (1 << 24):
            sys.stderr.write('Message too long!')
            return
        sys.stdout.write(chr(l & 0xff))
        sys.stdout.write(chr((l >> 8) & 0xff))
        sys.stdout.write(chr(l >> 16))
        sys.stdout.write(s)
        sys.stdout.flush()

    @staticmethod
    def recv():
        s = sys.stdin.read(3)
        l = ord(s[0]) + (ord(s[1]) << 8) + (ord(s[2]) << 16)
        return sys.stdin.read(l)

    @staticmethod
    def log(s):
        sys.stderr.write(s)
        sys.stderr.flush()

Bot.log('Hello from the snake! Side #' + Bot.recv() + '\n')
random.seed()

used = [[False] * 3 for _ in range(3)]
while True:
    r, c = map(int, Bot.recv().split())
    if r != -1: used[r][c] = True

    while True:
        r = random.randrange(0, 3)
        c = random.randrange(0, 3)
        if not used[r][c]: break

    used[r][c] = True
    Bot.send('%d %d' % (r, c))
