#!/bin/sh
sudo rm a.out
gcc *.c /usr/local/lib/libhiredis.a -O2
sudo chown root a.out
sudo chmod +s a.out
