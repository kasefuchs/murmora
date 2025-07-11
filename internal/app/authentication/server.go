// Copyright (c) Kasefuchs
// SPDX-License-Identifier: MPL-2.0

package authentication

import (
	"context"

	"github.com/kasefuchs/murmora/api/proto/murmora/authentication/v1"
	"github.com/kasefuchs/murmora/api/proto/murmora/session/v1"
	"github.com/kasefuchs/murmora/api/proto/murmora/user/v1"
	"github.com/matthewhartstonge/argon2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	authentication.UnimplementedAuthenticationServiceServer

	argon         *argon2.Config
	userClient    user.UserServiceClient
	sessionClient session.SessionServiceClient
}

func NewServer(
	userClient user.UserServiceClient,
	sessionClient session.SessionServiceClient,
) *Server {
	argon := argon2.DefaultConfig()

	return &Server{
		argon:         &argon,
		userClient:    userClient,
		sessionClient: sessionClient,
	}
}

func (s *Server) Register(ctx context.Context, request *authentication.RegisterRequest) (*authentication.TokenResponse, error) {
	if err := request.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	hashEncoded, err := s.argon.HashEncoded([]byte(request.Password))
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "hash encoding failed: %s", err.Error())
	}

	userDataResponse, err := s.userClient.CreateUser(ctx, &user.CreateUserRequest{
		Name:         request.Name,
		Email:        request.Email,
		PasswordHash: string(hashEncoded),
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create user: %s", err.Error())
	}

	createSessionResponse, err := s.sessionClient.CreateSession(ctx, &session.CreateSessionRequest{
		UserId: userDataResponse.Id,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create session: %s", err.Error())
	}

	return &authentication.TokenResponse{
		Token: createSessionResponse.Token,
	}, nil
}

func (s *Server) Login(ctx context.Context, request *authentication.LoginRequest) (*authentication.TokenResponse, error) {
	if err := request.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	userDataResponse, err := s.userClient.GetUser(ctx, &user.GetUserRequest{
		Query: &user.GetUserRequest_Email{
			Email: request.Email,
		},
	})
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Failed to get user: %v", err)
	}

	ok, err := argon2.VerifyEncoded([]byte(request.Password), []byte(userDataResponse.PasswordHash))
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Failed to verify password: %v", err)
	}

	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "Invalid password")
	}

	createSessionResponse, err := s.sessionClient.CreateSession(ctx, &session.CreateSessionRequest{
		UserId: userDataResponse.Id,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to create session: %s", err.Error())
	}

	return &authentication.TokenResponse{
		Token: createSessionResponse.Token,
	}, nil
}
