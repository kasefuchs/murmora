// Copyright (c) Kasefuchs
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"flag"

	"github.com/kasefuchs/murmora/api/proto/murmora/authentication/v1"
	"github.com/kasefuchs/murmora/api/proto/murmora/session/v1"
	"github.com/kasefuchs/murmora/api/proto/murmora/token/v1"
	"github.com/kasefuchs/murmora/api/proto/murmora/user/v1"
	authenticationService "github.com/kasefuchs/murmora/internal/app/authentication"
	service "github.com/kasefuchs/murmora/internal/app/monolith"
	sessionService "github.com/kasefuchs/murmora/internal/app/session"
	tokenService "github.com/kasefuchs/murmora/internal/app/token"
	userService "github.com/kasefuchs/murmora/internal/app/user"
	"github.com/kasefuchs/murmora/internal/pkg/config"
	"github.com/kasefuchs/murmora/internal/pkg/database"
	"github.com/kasefuchs/murmora/internal/pkg/grpc/client"
	"github.com/kasefuchs/murmora/internal/pkg/grpc/server"
	"google.golang.org/grpc"
)

func main() {
	configFile := flag.String("config-file", "configs/monolith.hcl", "Path to the config file")
	flag.Parse()

	cfg := config.New[service.Config]()
	cfg.MustLoadConfigFile(*configFile)

	db := database.MustNew(&cfg.Value.Database)
	db.MustMigrate(&userService.User{}, &tokenService.Token{}, &sessionService.Session{})

	userClient := client.MustNew(&cfg.Value.Client, user.NewUserServiceClient)
	tokenClient := client.MustNew(&cfg.Value.Client, token.NewTokenServiceClient)
	sessionClient := client.MustNew(&cfg.Value.Client, session.NewSessionServiceClient)

	server.MustServe(&cfg.Value.Server, func(srv *grpc.Server) {
		userServer := userService.NewServer(db)
		tokenServer := tokenService.NewServer(db)
		sessionServer := sessionService.NewServer(db, userClient, tokenClient)
		authenticationServer := authenticationService.NewServer(userClient, sessionClient)

		user.RegisterUserServiceServer(srv, userServer)
		token.RegisterTokenServiceServer(srv, tokenServer)
		session.RegisterSessionServiceServer(srv, sessionServer)
		authentication.RegisterAuthenticationServiceServer(srv, authenticationServer)
	})
}
