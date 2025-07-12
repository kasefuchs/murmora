// Copyright (c) Kasefuchs
// SPDX-License-Identifier: MPL-2.0

package user

import (
	"context"

	"github.com/google/uuid"
	"github.com/kasefuchs/murmora/api/proto/murmora/common/v1"
	"github.com/kasefuchs/murmora/api/proto/murmora/user/v1"
	"github.com/kasefuchs/murmora/internal/pkg/database"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	user.UnimplementedUserServiceServer

	repository *Repository
}

func NewServer(db *database.Database) *Server {
	repository := NewRepository(db)

	return &Server{
		repository: repository,
	}
}

func (s *Server) CreateUser(_ context.Context, request *user.CreateUserRequest) (*user.UserDataResponse, error) {
	if err := request.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	id, err := uuid.NewV7()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate UUID: %v", err)
	}

	entity, err := s.repository.Create(&User{
		ID:           id,
		Name:         request.Name,
		Email:        request.Email,
		PasswordHash: request.PasswordHash,
	})
	if err != nil {
		return nil, status.Errorf(codes.AlreadyExists, "user already exists: %v", err)
	}

	return &user.UserDataResponse{
		Id:           common.NewUUID(entity.ID),
		Name:         entity.Name,
		Email:        entity.Email,
		PasswordHash: entity.PasswordHash,
	}, nil
}

func (s *Server) GetUser(_ context.Context, request *user.GetUserRequest) (*user.UserDataResponse, error) {
	if err := request.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	var err error
	var entity *User

	switch query := request.GetQuery().(type) {
	case *user.GetUserRequest_Id:
		entity, err = s.repository.FindByID(query.Id.Value)
	case *user.GetUserRequest_Name:
		entity, err = s.repository.FindByName(query.Name)
	case *user.GetUserRequest_Email:
		entity, err = s.repository.FindByEmail(query.Email)
	default:
		return nil, status.Error(codes.InvalidArgument, "invalid query type")
	}

	if err != nil {
		return nil, status.Errorf(codes.Internal, "error querying user: %v", err)
	}
	if entity == nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	return &user.UserDataResponse{
		Id:           common.NewUUID(entity.ID),
		Name:         entity.Name,
		Email:        entity.Email,
		PasswordHash: entity.PasswordHash,
	}, nil
}
