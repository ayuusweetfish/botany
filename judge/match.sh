#!/bin/sh

MID=$1
JUDGE=$2
echo "Running match $MID (judge $JUDGE)"
shift 2

mkdir -p matches/$MID
sleep 1

i=0
for s in $@; do
    submissions/$s/bin 2>matches/$MID/$i.log
    i=$((i + 1))
done
