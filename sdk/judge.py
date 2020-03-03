import bot_judge_py as bot_judge
import sys

print(bot_judge.init(sys.argv))

bot_judge.send(0, '0')
bot_judge.send(1, '1')

bot_judge.send(0, '-1 -1')
s, err = bot_judge.recv(0, 1000)
print(s)
print(err)

bot_judge.finish()
