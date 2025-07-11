// Copyright (c) Kasefuchs
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"github.com/kasefuchs/murmora/api/proto/murmora/user/v1"
	"github.com/kasefuchs/murmora/internal/app/user/config"
	"github.com/kasefuchs/murmora/internal/app/user/data"
	"github.com/kasefuchs/murmora/internal/app/user/service"
	conf "github.com/kasefuchs/murmora/internal/pkg/config"
	"github.com/kasefuchs/murmora/internal/pkg/database"
	"github.com/kasefuchs/murmora/internal/pkg/grpc/server"
	"google.golang.org/grpc"
)

func main() {
	cfg := conf.New[config.Config]()
	cfg.MustLoadConfigFile("configs/user.hcl")

	db := database.MustNew(&cfg.Value.Database)
	db.MustMigrate(&data.User{})

	userRepository := data.NewUserRepository(db)

	server.MustServe(&cfg.Value.Server, func(srv *grpc.Server) {
		userServer := service.NewUserServiceServer(userRepository)

		user.RegisterUserServiceServer(srv, userServer)
	})
}
