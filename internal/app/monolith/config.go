// Copyright (c) Kasefuchs
// SPDX-License-Identifier: MPL-2.0

package monolith

import (
	"github.com/kasefuchs/murmora/internal/pkg/database"
	"github.com/kasefuchs/murmora/internal/pkg/grpc/client"
	"github.com/kasefuchs/murmora/internal/pkg/grpc/server"
)

type Config struct {
	Server   server.Config   `koanf:"server"`
	Client   client.Config   `koanf:"client"`
	Database database.Config `koanf:"database"`
}
