// Copyright (c) Kasefuchs
// SPDX-License-Identifier: MPL-2.0

package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/kasefuchs/murmora/api/proto/murmora/user/v1"
	"github.com/kasefuchs/murmora/internal/app/user/data"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserServiceServer struct {
	user.UnimplementedUserServiceServer

	userRepository *data.UserRepository
}

func (s *UserServiceServer) CreateUser(_ context.Context, request *user.CreateUserRequest) (*user.CreateUserResponse, error) {
	if err := request.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	id, err := uuid.NewV7()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to generate UUID: %v", err)
	}

	entity, err := s.userRepository.Create(&data.User{
		ID:           id,
		Name:         request.Name,
		Email:        request.Email,
		PasswordHash: request.PasswordHash,
	})
	if err != nil {
		return nil, status.Errorf(codes.AlreadyExists, "User already exists: %v", err)
	}

	return &user.CreateUserResponse{
		Id:           entity.ID.String(),
		Name:         entity.Name,
		Email:        entity.Email,
		PasswordHash: entity.PasswordHash,
	}, nil
}

func NewUserServiceServer(userRepository *data.UserRepository) *UserServiceServer {
	return &UserServiceServer{
		userRepository: userRepository,
	}
}
