#!/bin/bash

NAME=$(</dev/stdin)
echo "Hello $NAME!"

echo "stdin: $NAME, args: $@" >> test_file.txt
