#!/bin/bash

MID=$1
JUDGE=$2
shift 2

mkdir -p matches/$MID

argv=()

i=0
for s in $@; do
    argv+=("isolate --run -b $i --dir=box=/var/botany/submissions/$s --dir=tmp= --dir=proc= --silent -- bin")
    i=$((i + 1))
done

i=0
for s in $@; do
    argv+=(matches/$MID/$i.log)
    i=$((i + 1))
done

submissions/$JUDGE/bin "${argv[@]}"
