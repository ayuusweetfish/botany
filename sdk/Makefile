LOGS = p1.log p2.log

run: judge player
	@./judge ./player ./player $(LOGS)

run-lua: CMD = "LUA_PATH=lib/?.lua lua player.lua"
run-lua: judge
	@./judge $(CMD) $(CMD) $(LOGS)

run-py: CMD = "PYTHONPATH=lib python3 player.py"
run-py: judge
	@./judge $(CMD) $(CMD) $(LOGS)

ifneq (, $(wildcard judge.c))
judge: lib/bot.h lib/bot.c judge.c
	$(CC) judge.c lib/bot.c -Ilib -O2 -o judge
else
ifneq (, $(wildcard judge.cpp))
judge: lib/bot.h lib/bot.c judge.cpp
	$(CXX) judge.cpp -x c lib/bot.c -Ilib -O2 -o judge
endif
endif

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
	$(RM) judge player p1.log p2.log

.PHONY: run run-lua run-py
