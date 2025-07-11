// Copyright (c) Kasefuchs
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"flag"

	"github.com/kasefuchs/murmora/api/proto/murmora/session/v1"
	"github.com/kasefuchs/murmora/api/proto/murmora/token/v1"
	"github.com/kasefuchs/murmora/api/proto/murmora/user/v1"
	service "github.com/kasefuchs/murmora/internal/app/session"
	"github.com/kasefuchs/murmora/internal/pkg/config"
	"github.com/kasefuchs/murmora/internal/pkg/database"
	"github.com/kasefuchs/murmora/internal/pkg/grpc/client"
	"github.com/kasefuchs/murmora/internal/pkg/grpc/server"
	"google.golang.org/grpc"
)

func main() {
	configFile := flag.String("config-file", "configs/session.hcl", "Path to the config file")
	flag.Parse()

	cfg := config.New[service.Config]()
	cfg.MustLoadConfigFile(*configFile)

	db := database.MustNew(&cfg.Value.Database)
	db.MustMigrate(&service.Session{})

	repository := service.NewRepository(db)
	userClient := client.MustNew(&cfg.Value.UserService, user.NewUserServiceClient)
	tokenClient := client.MustNew(&cfg.Value.TokenService, token.NewTokenServiceClient)

	server.MustServe(&cfg.Value.Server, func(srv *grpc.Server) {
		sessionServer := service.NewServer(repository, userClient, tokenClient)

		session.RegisterSessionServiceServer(srv, sessionServer)
	})
}
