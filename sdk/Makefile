LOGS = p1.log p2.log

run: judge player
	@./judge ./player ./player $(LOGS)

run-lua: CMD = "LUA_PATH=lib/?.lua lua player.lua"
run-lua: judge
	@./judge $(CMD) $(CMD) $(LOGS)

run-py: CMD = "PYTHONPATH=lib python3 player.py"
run-py: judge
	@./judge $(CMD) $(CMD) $(LOGS)

run-pyjudge: lib/bot_jury_py.so player
	@PYTHONPATH=lib python3 judge.py ./player ./player $(LOGS)

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

.PHONY: run run-lua run-py
