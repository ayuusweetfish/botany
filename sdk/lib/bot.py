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
        s = sys.stdin.read(3)
        l = ord(s[0]) + (ord(s[1]) << 8) + (ord(s[2]) << 16)
        return sys.stdin.read(l)

    @staticmethod
    def log(s):
        sys.stderr.write(s)
        sys.stderr.flush()
