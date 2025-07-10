// Copyright (c) Kasefuchs
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func MustNew[T any](cfg *Config, newClient func(grpc.ClientConnInterface) T) T {
	conn, err := grpc.NewClient(cfg.Address, grpc.WithTransportCredentials(
		insecure.NewCredentials(),
	))
	if err != nil {
		log.Fatal().Err(err).Msg("Error creating gRPC client")
	}

	return newClient(conn)
}
