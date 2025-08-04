// Copyright (c) Kasefuchs
// SPDX-License-Identifier: MPL-2.0

//go:build tools
// +build tools

package tools

import (
	_ "github.com/bufbuild/buf/cmd/buf"
	_ "github.com/envoyproxy/protoc-gen-validate"
	_ "github.com/fullstorydev/grpcui/cmd/grpcui"
	_ "github.com/fullstorydev/grpcurl/cmd/grpcurl"
	_ "github.com/hashicorp/copywrite"
	_ "google.golang.org/grpc/cmd/protoc-gen-go-grpc"
	_ "google.golang.org/protobuf/cmd/protoc-gen-go"
	_ "mvdan.cc/gofumpt"
)
