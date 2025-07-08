// Copyright (c) Kasefuchs
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"fmt"
	"log"
	"net"

	"github.com/kasefuchs/murmora/api/proto/murmora/token/v1"
	"github.com/kasefuchs/murmora/internal/app/token/config"
	"github.com/kasefuchs/murmora/internal/app/token/data"
	"github.com/kasefuchs/murmora/internal/app/token/service"
	"github.com/kasefuchs/murmora/internal/pkg/database"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	lis, err := net.Listen("tcp", cfg.ListenAddress)
	if err != nil {
		log.Fatalf("Error listening on %s: %v", cfg.ListenAddress, err)
	}

	fmt.Printf("Listening on: %v", lis.Addr().String())

	db, err := database.New(cfg.Database)
	if err != nil {
		log.Fatalf("Error connecting to data: %v", err)
	}

	if err := db.Migrate(&data.Token{}); err != nil {
		log.Fatalf("Error migrating data: %v", err)
	}

	tokenRepository := data.NewTokenRepository(db)
	tokenServer := service.NewTokenServiceServer(tokenRepository)

	grpcServer := grpc.NewServer()

	token.RegisterTokenServiceServer(grpcServer, tokenServer)

	if cfg.EnableReflection {
		reflection.Register(grpcServer)
	}

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Error serving grpc: %v", err)
	}
}
