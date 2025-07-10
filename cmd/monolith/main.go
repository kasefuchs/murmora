// Copyright (c) Kasefuchs
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"github.com/kasefuchs/murmora/api/proto/murmora/authentication/v1"
	"github.com/kasefuchs/murmora/api/proto/murmora/session/v1"
	"github.com/kasefuchs/murmora/api/proto/murmora/token/v1"
	"github.com/kasefuchs/murmora/api/proto/murmora/user/v1"
	authenticationservice "github.com/kasefuchs/murmora/internal/app/authentication/service"
	"github.com/kasefuchs/murmora/internal/app/monolith/config"
	sessiondata "github.com/kasefuchs/murmora/internal/app/session/data"
	sessionservice "github.com/kasefuchs/murmora/internal/app/session/service"
	tokendata "github.com/kasefuchs/murmora/internal/app/token/data"
	tokenservice "github.com/kasefuchs/murmora/internal/app/token/service"
	userdata "github.com/kasefuchs/murmora/internal/app/user/data"
	userservice "github.com/kasefuchs/murmora/internal/app/user/service"
	conf "github.com/kasefuchs/murmora/internal/pkg/config"
	"github.com/kasefuchs/murmora/internal/pkg/database"
	"github.com/kasefuchs/murmora/internal/pkg/grpc/client"
	"github.com/kasefuchs/murmora/internal/pkg/grpc/server"
	"google.golang.org/grpc"
)

func main() {
	cfg := conf.New[config.Config]()
	cfg.MustLoadConfigFile("configs/monolith.hcl")

	db := database.MustNew(cfg.Value.Database)
	db.MustMigrate(&userdata.User{}, &tokendata.Token{}, &sessiondata.Session{})

	userRepository := userdata.NewUserRepository(db)
	tokenRepository := tokendata.NewTokenRepository(db)
	sessionRepository := sessiondata.NewSessionRepository(db)

	userServiceClient := client.MustNew(cfg.Value.Client, user.NewUserServiceClient)
	tokenServiceClient := client.MustNew(cfg.Value.Client, token.NewTokenServiceClient)
	sessionServiceClient := client.MustNew(cfg.Value.Client, session.NewSessionServiceClient)

	server.MustServe(cfg.Value.Server, func(srv *grpc.Server) {
		userServer := userservice.NewUserServiceServer(userRepository)
		tokenServer := tokenservice.NewTokenServiceServer(tokenRepository)
		sessionServer := sessionservice.NewSessionServiceServer(sessionRepository, userServiceClient, tokenServiceClient)
		authenticationServer := authenticationservice.NewAuthenticationServiceServer(userServiceClient, sessionServiceClient)

		user.RegisterUserServiceServer(srv, userServer)
		token.RegisterTokenServiceServer(srv, tokenServer)
		session.RegisterSessionServiceServer(srv, sessionServer)
		authentication.RegisterAuthenticationServiceServer(srv, authenticationServer)
	})
}
