LOGS = p1.log p2.log

# Language selection

JUDGE_EXEC := ./judge
JUDGE_ARGS :=
PLAYER_CMD := ./player

run-default: judge-c player-c run

judge-c: judge
	$(eval JUDGE_EXEC := ./judge)
	$(eval JUDGE_ARGS := )
	@:

judge-py: lib/bot_jury_py.so
	$(eval JUDGE_EXEC := python3)
	$(eval JUDGE_ARGS := judge.py)
	@:

player-c: player
	$(eval PLAYER_EXEC := ./player)
	$(eval PLAYER_ARGS := )
	@:

player-lua: player.lua
	$(eval PLAYER_CMD := LUA_PATH=lib/?.lua lua player.lua)
	@:

player-py: player.py
	$(eval PLAYER_CMD := PYTHONPATH=lib python3 player.py)
	@:

run: export PYTHONPATH = lib
run: export LUA_PATH = lib/?.lua
run:
	@$(JUDGE_EXEC) $(JUDGE_ARGS) "$(PLAYER_CMD)" "$(PLAYER_CMD)" $(LOGS)

.PHONY: judge-py player-lua player-py run-default run

# Building

ifneq (, $(wildcard judge.c))
judge: lib/bot.h lib/bot.c judge.c
	$(CC) judge.c lib/bot.c -Ilib -O2 -o judge
else
ifneq (, $(wildcard judge.cpp))
judge: lib/bot.h lib/bot.c judge.cpp
	$(CXX) judge.cpp -x c lib/bot.c -Ilib -O2 -o judge
endif
endif

lib/bot_jury_py.so: lib/bot.h lib/bot.c lib/bot_jury_py.c
	$(CC) lib/bot.c -o lib/bot.o -c -fPIC
	$(CC) lib/bot_jury_py.c -o lib/bot_jury_py.o -c -fPIC `python3-config --cflags`
	$(CC) lib/bot.o lib/bot_jury_py.o -shared -o lib/bot_jury_py.so `python3-config --ldflags`

ifneq (, $(wildcard player.c))
player: lib/bot.h lib/bot.c player.c
	$(CC) player.c lib/bot.c -Ilib -O2 -o player
else
ifneq (, $(wildcard player.cpp))
player: lib/bot.h lib/bot.c player.cpp
	$(CXX) player.cpp -x c lib/bot.c -Ilib -O2 -o player
endif
endif

clean:
	$(RM) judge player p1.log p2.log lib/bot.o lib/bot_jury_py.o lib/bot_jury_py.so

.PHONY: clean
