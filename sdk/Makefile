run: judge player
	@./judge ./player ./player

run-lua: CMD = "LUA_PATH=lib/?.lua lua player.lua"
run-lua: judge
	@./judge $(CMD) $(CMD)

run-py: CMD = "PYTHONPATH=lib python player.py"
run-py: judge
	@./judge $(CMD) $(CMD)

judge: lib/bot.h lib/bot.c judge.c
	$(CC) lib/bot.c judge.c -Ilib -O2 -o judge

player: lib/bot.h lib/bot.c player.c
	$(CC) lib/bot.c player.c -Ilib -O2 -o player

clean:
	$(RM) judge player

.PHONY: run run-py run-lua
