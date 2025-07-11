// Copyright (c) Kasefuchs
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"flag"

	"github.com/kasefuchs/murmora/api/proto/murmora/token/v1"
	service "github.com/kasefuchs/murmora/internal/app/token"
	"github.com/kasefuchs/murmora/internal/pkg/config"
	"github.com/kasefuchs/murmora/internal/pkg/database"
	"github.com/kasefuchs/murmora/internal/pkg/grpc/server"
	"google.golang.org/grpc"
)

func main() {
	configFile := flag.String("config-file", "configs/token.hcl", "Path to the config file")
	flag.Parse()

	cfg := config.New[service.Config]()
	cfg.MustLoadConfigFile(*configFile)

	db := database.MustNew(&cfg.Value.Database)
	db.MustMigrate(&service.Token{})

	tokenRepository := service.NewRepository(db)

	server.MustServe(&cfg.Value.Server, func(srv *grpc.Server) {
		tokenServer := service.NewServer(tokenRepository)

		token.RegisterTokenServiceServer(srv, tokenServer)
	})
}
