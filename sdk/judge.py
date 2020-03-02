import bot_jury_py as BotJury
import sys

print(BotJury.init(sys.argv))

BotJury.send(0, '0')
BotJury.send(1, '1')

BotJury.send(0, '-1 -1')
s, err = BotJury.recv(0, 1000)
print(s)
print(err)

BotJury.finish()
