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

func (s *UserServiceServer) CreateUser(_ context.Context, request *user.CreateUserRequest) (*user.UserResponse, error) {
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

	return &user.UserResponse{
		Id:           entity.ID.String(),
		Name:         entity.Name,
		Email:        entity.Email,
		PasswordHash: entity.PasswordHash,
	}, nil
}

func (s *UserServiceServer) GetUser(_ context.Context, request *user.GetUserRequest) (*user.UserResponse, error) {
	if err := request.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	var err error
	var entity *data.User

	switch query := request.GetQuery().(type) {
	case *user.GetUserRequest_Id:
		entity, err = s.userRepository.FindByID(query.Id)

	case *user.GetUserRequest_Name:
		entity, err = s.userRepository.FindByName(query.Name)

	case *user.GetUserRequest_Email:
		entity, err = s.userRepository.FindByEmail(query.Email)

	default:
		return nil, status.Error(codes.InvalidArgument, "invalid query type")
	}

	if err != nil {
		return nil, status.Errorf(codes.Internal, "error querying user: %v", err)
	}

	if entity == nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	return &user.UserResponse{
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
