#!/bin/sh

sudo cp compile.sh match.sh $1/var/botany
sudo cp -r lib $1/var/botany/lib
find $1/var/botany/matches -type d -name '12*' -exec rm -rf {} \;
find $1/var/botany/submissions -type d -name '12*' -exec rm -rf {} \;
