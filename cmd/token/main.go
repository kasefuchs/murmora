// Copyright (c) Kasefuchs
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"flag"

	"github.com/kasefuchs/murmora/api/proto/murmora/token/v1"
	"github.com/kasefuchs/murmora/internal/app/token/config"
	"github.com/kasefuchs/murmora/internal/app/token/data"
	"github.com/kasefuchs/murmora/internal/app/token/service"
	conf "github.com/kasefuchs/murmora/internal/pkg/config"
	"github.com/kasefuchs/murmora/internal/pkg/database"
	"github.com/kasefuchs/murmora/internal/pkg/grpc/server"
	"google.golang.org/grpc"
)

func main() {
	configFile := flag.String("config-file", "configs/token.hcl", "Path to the config file")
	flag.Parse()

	cfg := conf.New[config.Config]()
	cfg.MustLoadConfigFile(*configFile)

	db := database.MustNew(&cfg.Value.Database)
	db.MustMigrate(&data.Token{})

	tokenRepository := data.NewTokenRepository(db)

	server.MustServe(&cfg.Value.Server, func(srv *grpc.Server) {
		tokenServer := service.NewTokenServiceServer(tokenRepository)

		token.RegisterTokenServiceServer(srv, tokenServer)
	})
}
