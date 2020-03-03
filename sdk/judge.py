import bot_judge_py as bot_judge
import sys

if bot_judge.init(sys.argv) != 2:
    print('Incorrect number of players!')
    sys.exit()

# Inform each player of their side
bot_judge.send(0, '0')
bot_judge.send(1, '1')

board = [[-1] * 3 for _ in range(3)]
win = -1
row, col = -1, -1   # Last move

for count in range(9):
    move = count % 2    # Current player
    bot_judge.send(move, '%d %d' % (row, col))
    resp, err = bot_judge.recv(move, 1000)

    if resp == None:
        sys.stderr.write(
            'Side #%d errors with %d, considered resignation\n' %
            (move, err))
        win = move ^ 1
        break

    try:
        row, col = map(int, resp.split())
    except:
        row, col = -1, -1

    if row not in range(3) or col not in range(3):
        sys.stderr.write(
            'Side #%d generates invalid response (%s), considered resignation\n' %
            (move, resp if len(resp) < 10 else resp[:7] + '...'))
        win = move ^ 1
        break

    if board[row][col] != -1:
        sys.stderr.write(
            'Side #%d invalid move at (%d, %d), considered resignation\n' %
            (move, row, col))
        win = move ^ 1
        break

    board[row][col] = move

    # Check winning condition
    for i in range(3):
        if ((board[i][0] == move and board[i][1] == move and board[i][2] == move)
         or (board[0][i] == move and board[1][i] == move and board[2][i] == move)):
            win = move
            break

    if ((board[0][0] == move and board[1][1] == move and board[2][2] == move)
     or (board[2][0] == move and board[1][1] == move and board[0][2] == move)):
        win = move

    if win != -1:
        sys.stderr.write('Side #%d wins!\n' % win)
        break

    count += 1

sys.stderr.write('{\n  "winner": %d,\n  "board": "' % win)
sys.stderr.write(
    '\\n'.join(list(
        map(lambda r: ''.join(list(map(lambda x: ".ox"[x + 1], r))), board)
    ))
)
sys.stderr.write('"\n}\n')

bot_judge.finish()
