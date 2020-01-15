#!/bin/sh

sudo cp compile.sh match.sh $1/var/botany
sudo rm -rf $1/var/botany/lib
sudo cp -r lib $1/var/botany/lib
find $1/var/botany/submissions -type d -name '*' -exec rm -rf {} \;
find $1/var/botany/matches -type d -name '*' -exec rm -rf {} \;
sudo mkdir -p $1/var/botany/submissions
sudo chown 1000 $1/var/botany/submissions
sudo chgrp 1000 $1/var/botany/submissions
sudo mkdir -p $1/var/botany/matches
sudo chown 1000 $1/var/botany/matches
sudo chgrp 1000 $1/var/botany/matches
