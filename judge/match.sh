#!/bin/sh

MID=$1
JUDGE=$2
echo "Running match $MID (judge $JUDGE)"
shift 2

sleep 1

for s in $@; do submissions/$s/bin; done
