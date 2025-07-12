// Copyright (c) Kasefuchs
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"flag"

	"github.com/kasefuchs/murmora/api/proto/murmora/authentication/v1"
	"github.com/kasefuchs/murmora/api/proto/murmora/session/v1"
	"github.com/kasefuchs/murmora/api/proto/murmora/user/v1"
	service "github.com/kasefuchs/murmora/internal/app/authentication"
	"github.com/kasefuchs/murmora/internal/pkg/config"
	"github.com/kasefuchs/murmora/internal/pkg/grpc/client"
	"github.com/kasefuchs/murmora/internal/pkg/grpc/server"
	"google.golang.org/grpc"
)

func main() {
	configFile := flag.String("config-file", "configs/authentication.hcl", "Path to the config file")
	flag.Parse()

	cfg := config.New[service.Config]()
	cfg.MustLoadConfigFile(*configFile)

	userClient := client.MustNew(&cfg.Value.UserService, user.NewUserServiceClient)
	sessionClient := client.MustNew(&cfg.Value.SessionService, session.NewSessionServiceClient)

	server.MustServe(&cfg.Value.Server, func(srv *grpc.Server) {
		authenticationServer := service.NewServer(userClient, sessionClient)

		authentication.RegisterAuthenticationServiceServer(srv, authenticationServer)
	})
}
