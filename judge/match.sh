#!/bin/sh

echo "Running match $1"
MID=$1
shift

sleep 1

for s in $@; do submissions/$s/bin; done
