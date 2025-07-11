// Copyright (c) Kasefuchs
// SPDX-License-Identifier: MPL-2.0

package config

import (
	"github.com/kasefuchs/murmora/internal/pkg/database"
	"github.com/kasefuchs/murmora/internal/pkg/grpc/client"
	"github.com/kasefuchs/murmora/internal/pkg/grpc/server"
)

type Config struct {
	Server       server.Config   `koanf:"server"`
	Database     database.Config `koanf:"database"`
	UserService  client.Config   `koanf:"user_service"`
	TokenService client.Config   `koanf:"token_service"`
}
