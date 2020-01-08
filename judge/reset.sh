#!/bin/sh

sudo cp compile.sh match.sh $1/var/botany
sudo cp lib/* $1/var/botany/lib
sudo rm -rf $1/var/botany/matches/*
sudo rm -rf $1/var/botany/submissions/*
