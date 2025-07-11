// Copyright (c) Kasefuchs
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"flag"

	"github.com/kasefuchs/murmora/api/proto/murmora/authentication/v1"
	"github.com/kasefuchs/murmora/api/proto/murmora/session/v1"
	"github.com/kasefuchs/murmora/api/proto/murmora/user/v1"
	"github.com/kasefuchs/murmora/internal/app/authentication/config"
	"github.com/kasefuchs/murmora/internal/app/authentication/service"
	conf "github.com/kasefuchs/murmora/internal/pkg/config"
	"github.com/kasefuchs/murmora/internal/pkg/grpc/client"
	"github.com/kasefuchs/murmora/internal/pkg/grpc/server"
	"google.golang.org/grpc"
)

func main() {
	configFile := flag.String("config-file", "configs/authentication.hcl", "Path to the config file")
	flag.Parse()

	cfg := conf.New[config.Config]()
	cfg.MustLoadConfigFile(*configFile)

	userServiceClient := client.MustNew(&cfg.Value.UserService, user.NewUserServiceClient)
	sessionServiceClient := client.MustNew(&cfg.Value.SessionService, session.NewSessionServiceClient)

	server.MustServe(&cfg.Value.Server, func(srv *grpc.Server) {
		authenticationServer := service.NewAuthenticationServiceServer(userServiceClient, sessionServiceClient)

		authentication.RegisterAuthenticationServiceServer(srv, authenticationServer)
	})
}
