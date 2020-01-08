#!/bin/bash

MID=$1
JUDGE=$2
echo "Running match $MID (judge $JUDGE)"
shift 2

mkdir -p matches/$MID

argv=()

for s in $@; do
    argv+=(submissions/$s/bin)
done

i=0
for s in $@; do
    argv+=(matches/$MID/$i.log)
    i=$((i + 1))
done

submissions/$JUDGE/bin ${argv[*]}
