// Copyright (c) Kasefuchs
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"fmt"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"

	"github.com/kasefuchs/murmora/api/proto/murmora/token/v1"
	"github.com/kasefuchs/murmora/api/proto/murmora/user/v1"

	"github.com/kasefuchs/murmora/api/proto/murmora/session/v1"
	"github.com/kasefuchs/murmora/internal/app/session/config"
	"github.com/kasefuchs/murmora/internal/app/session/data"
	"github.com/kasefuchs/murmora/internal/app/session/service"
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

	if err := db.Migrate(&data.Session{}); err != nil {
		log.Fatalf("Error migrating data: %v", err)
	}

	userGrpcClient, err := grpc.NewClient(cfg.UserServiceUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Error creating user grpc client: %v", err)
	}

	userServiceClient := user.NewUserServiceClient(userGrpcClient)

	tokenGrpcClient, err := grpc.NewClient(cfg.TokenServiceUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Error creating token grpc client: %v", err)
	}

	tokenServiceClient := token.NewTokenServiceClient(tokenGrpcClient)

	sessionRepository := data.NewSessionRepository(db)
	sessionServer := service.NewSessionServiceServer(sessionRepository, userServiceClient, tokenServiceClient)

	grpcServer := grpc.NewServer()

	session.RegisterSessionServiceServer(grpcServer, sessionServer)

	if cfg.EnableReflection {
		reflection.Register(grpcServer)
	}

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Error serving grpc: %v", err)
	}
}
