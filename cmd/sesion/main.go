// Copyright (c) Kasefuchs
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"github.com/kasefuchs/murmora/api/proto/murmora/session/v1"
	"github.com/kasefuchs/murmora/api/proto/murmora/token/v1"
	"github.com/kasefuchs/murmora/api/proto/murmora/user/v1"
	"github.com/kasefuchs/murmora/internal/app/session/config"
	"github.com/kasefuchs/murmora/internal/app/session/data"
	"github.com/kasefuchs/murmora/internal/app/session/service"
	conf "github.com/kasefuchs/murmora/internal/pkg/config"
	"github.com/kasefuchs/murmora/internal/pkg/database"
	"github.com/kasefuchs/murmora/internal/pkg/grpc/client"
	"github.com/kasefuchs/murmora/internal/pkg/grpc/server"
	"google.golang.org/grpc"
)

func main() {
	cfg := conf.New[config.Config]()
	cfg.MustLoadConfigFile("configs/session.hcl")

	db := database.MustNew(&cfg.Value.Database)
	db.MustMigrate(&data.Session{})

	sessionRepository := data.NewSessionRepository(db)

	userServiceClient := client.MustNew(&cfg.Value.UserService, user.NewUserServiceClient)
	tokenServiceClient := client.MustNew(&cfg.Value.TokenService, token.NewTokenServiceClient)

	server.MustServe(&cfg.Value.Server, func(srv *grpc.Server) {
		sessionServer := service.NewSessionServiceServer(sessionRepository, userServiceClient, tokenServiceClient)

		session.RegisterSessionServiceServer(srv, sessionServer)
	})
}
