run: judge player
	./judge ./player ./player

judge: bot.h bot.c judge.c
	$(CC) bot.c judge.c -O2 -o judge

player: bot.h bot.c player.c
	$(CC) bot.c player.c -O2 -o player

clean:
	$(RM) judge player

.PHONY: run
