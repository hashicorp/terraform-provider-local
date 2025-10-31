#!/bin/bash

NAME=$(</dev/stdin)
echo "Hello $NAME!"

echo "$NAME - args: $@" >> test_file.txt
