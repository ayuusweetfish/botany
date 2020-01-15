#!/bin/bash
sudo rm a.out
if [[ "$OSTYPE" == "darwin"* ]]; then
    gcc *.c /usr/local/lib/libhiredis.a /usr/local/lib/libb2.a -O2 -DBOT_POSIX_COMPLIANT
else
    gcc *.c -l:libhiredis.a -l:libb2.a -O2
fi
sudo chown root a.out
sudo chmod +s a.out
