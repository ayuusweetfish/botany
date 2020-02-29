import sys

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
        s = sys.stdin.buffer.read(3)
        l = s[0] + (s[1] << 8) + (s[2] << 16)
        return sys.stdin.buffer.read(l).decode()

    @staticmethod
    def log(s):
        sys.stderr.write(s)
        sys.stderr.flush()
