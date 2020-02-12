#!/bin/sh

echo "Compiling $1"
SID=$1
LANG=$2
CODE=submissions/$SID/code.$LANG
BIN=submissions/$SID/bin

if [[ "$LANG" == "c" ]]; then
    gcc $CODE -O2 -I/var/botany/lib /var/botany/lib/bot.c -o $BIN 2>&1
elif [[ "$LANG" == "cpp" ]]; then
    g++ $CODE -O2 -I/var/botany/lib /var/botany/lib/bot.c -o $BIN 2>&1
elif [[ "$LANG" == "lua" ]]; then
    :
fi
