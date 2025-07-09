// Copyright (c) Kasefuchs
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"fmt"
	"log"
	"net"

	"github.com/kasefuchs/murmora/api/proto/murmora/authentication/v1"
	"github.com/kasefuchs/murmora/api/proto/murmora/session/v1"
	"github.com/kasefuchs/murmora/api/proto/murmora/user/v1"
	"github.com/kasefuchs/murmora/internal/app/authentication/config"
	"github.com/kasefuchs/murmora/internal/app/authentication/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

	userGrpcClient, err := grpc.NewClient(cfg.UserServiceUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Error creating user grpc client: %v", err)
	}

	userServiceClient := user.NewUserServiceClient(userGrpcClient)

	sessionGrpcClient, err := grpc.NewClient(cfg.SessionServiceUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Error creating token grpc client: %v", err)
	}

	sessionServiceClient := session.NewSessionServiceClient(sessionGrpcClient)

	authenticationServer := service.NewAuthenticationServiceServer(userServiceClient, sessionServiceClient)

	grpcServer := grpc.NewServer()

	authentication.RegisterAuthenticationServiceServer(grpcServer, authenticationServer)

	if cfg.EnableReflection {
		reflection.Register(grpcServer)
	}

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Error serving grpc: %v", err)
	}
}
