#!/bin/bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0


NAME=$(cat)
echo "Hello $NAME!"

echo "stdin: $NAME, args: $@" >> test_file.txt
