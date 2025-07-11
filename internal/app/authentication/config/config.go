// Copyright (c) Kasefuchs
// SPDX-License-Identifier: MPL-2.0

package config

import (
	"github.com/kasefuchs/murmora/internal/pkg/grpc/client"
	"github.com/kasefuchs/murmora/internal/pkg/grpc/server"
)

type Config struct {
	Server         server.Config `koanf:"server"`
	UserService    client.Config `koanf:"user_service"`
	SessionService client.Config `koanf:"session_service"`
}
