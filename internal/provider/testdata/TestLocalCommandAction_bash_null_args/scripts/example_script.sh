#!/bin/bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0


if [ "$#" -ne 3 ]; then
  echo "You provided $# arguments, expected exactly 3 random number arguments (the 5 null arguments should be removed)." >&2
  exit 1
fi

NAME=$(</dev/stdin)
echo "Hello $NAME!"

echo "stdin: $NAME, args: $@" >> test_file.txt