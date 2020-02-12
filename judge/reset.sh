#!/bin/sh

sudo rm -rf $1/var/botany
sudo mkdir $1/var/botany
sudo cp compile.sh match.sh $1/var/botany
sudo cp -r lib $1/var/botany/lib
sudo mkdir -p $1/var/botany/submissions
sudo chown `whoami` $1/var/botany/submissions
sudo chgrp `whoami` $1/var/botany/submissions
sudo mkdir -p $1/var/botany/matches
sudo chown `whoami` $1/var/botany/matches
sudo chgrp `whoami` $1/var/botany/matches
