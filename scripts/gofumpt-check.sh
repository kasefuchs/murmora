#!/bin/bash
# Copyright (c) Kasefuchs
# SPDX-License-Identifier: MPL-2.0


output=$(go run mvdan.cc/gofumpt -l "$@")

if [ -n "$output" ]; then
    echo "$output"
    exit 1
fi

exit 0
