// Copyright (c) Kasefuchs
// SPDX-License-Identifier: MPL-2.0

package gateway

import "github.com/kasefuchs/murmora/internal/pkg/grpc/client"

type serverConfig struct {
	Address string `koanf:"address"`
	Prefix  string `koanf:"prefix"`
}

type Config struct {
	Server                serverConfig  `koanf:"server"`
	UserService           client.Config `koanf:"user_service"`
	TokenService          client.Config `koanf:"token_service"`
	SessionService        client.Config `koanf:"session_service"`
	AuthenticationService client.Config `koanf:"authentication_service"`
}
