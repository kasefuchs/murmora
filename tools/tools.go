// Copyright (c) Kasefuchs
// SPDX-License-Identifier: MPL-2.0

//go:build tools
// +build tools

package tools

import (
	_ "github.com/bufbuild/buf/cmd/buf"
	_ "github.com/hashicorp/copywrite"
	_ "mvdan.cc/gofumpt"
)
