#!/bin/bash

MID=$1
JUDGE=$2
shift 2

mkdir -p matches/$MID

argv=()

i=0
for s in $@; do
    CODE=`find submissions/$s/code.*`
    LANG=${CODE##*.}
    exec="bin"
    if [[ "$LANG" == "lua" ]]; then
        exec="/usr/bin/lua code.lua"
    fi
    argv+=("isolate --run -b $i --dir=box=/var/botany/submissions/$s --dir=tmp= --dir=proc= --silent -- $exec")
    i=$((i + 1))
done

i=0
for s in $@; do
    argv+=(matches/$MID/$i.log)
    i=$((i + 1))
done

submissions/$JUDGE/bin "${argv[@]}"
