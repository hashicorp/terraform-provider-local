#!/bin/bash
# Copyright IBM Corp. 2017, 2025
# SPDX-License-Identifier: MPL-2.0


NAME=$(cat)
echo "Hello $NAME!"

echo "stdin: $NAME, args: $@" >> test_file.txt
