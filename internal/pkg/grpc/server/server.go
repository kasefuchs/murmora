// Copyright (c) Kasefuchs
// SPDX-License-Identifier: MPL-2.0

package server

import (
	"net"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

func MustServe(cfg *Config, register func(*grpc.Server)) {
	lis, err := net.Listen("tcp", cfg.Address)
	if err != nil {
		log.Fatal().Err(err).Str("address", cfg.Address).Msg("Failed to listen")
	}

	srv := grpc.NewServer()

	register(srv)

	if cfg.Reflection {
		reflection.Register(srv)
	}

	if !cfg.DisableHealthCheck {
		grpc_health_v1.RegisterHealthServer(srv, health.NewServer())
	}

	log.Info().Str("address", cfg.Address).Msg("Server started")
	if err := srv.Serve(lis); err != nil {
		log.Fatal().Err(err).Str("address", cfg.Address).Msg("Failed to serve")
	}
}
