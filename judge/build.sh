#!/bin/sh
sudo rm a.out
gcc *.c /usr/local/lib/libhiredis.a /usr/local/lib/libb2.a -O2
sudo chown root a.out
sudo chmod +s a.out
