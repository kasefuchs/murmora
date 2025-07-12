// Copyright (c) Kasefuchs
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"flag"

	"github.com/goccy/go-json"
	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2"
	"github.com/kasefuchs/murmora/api/proto/murmora/authentication/v1"
	"github.com/kasefuchs/murmora/internal/app/gateway"
	"github.com/kasefuchs/murmora/internal/pkg/config"
	"github.com/kasefuchs/murmora/internal/pkg/grpc/client"
	"github.com/rs/zerolog/log"
)

var app = fiber.New(fiber.Config{
	JSONEncoder:           json.Marshal,
	JSONDecoder:           json.Unmarshal,
	DisableStartupMessage: true,
})

func main() {
	configFile := flag.String("config-file", "configs/gateway.hcl", "Path to the config file")
	flag.Parse()

	cfg := config.New[gateway.Config]()
	cfg.MustLoadConfigFile(*configFile)

	authenticationClient := client.MustNew(&cfg.Value.AuthenticationService, authentication.NewAuthenticationServiceClient)

	app.Use(fiberzerolog.New())
	gateway.SetupRoutes(app.Group(cfg.Value.Server.Prefix), authenticationClient)

	if err := app.Listen(cfg.Value.Server.Address); err != nil {
		log.Fatal().Err(err).Str("address", cfg.Value.Server.Address).Msg("Failed to listen")
	}
}
