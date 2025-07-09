// Copyright (c) Kasefuchs
// SPDX-License-Identifier: MPL-2.0

package service

import (
	"context"

	"github.com/kasefuchs/murmora/api/proto/murmora/authentication/v1"
	"github.com/kasefuchs/murmora/api/proto/murmora/session/v1"
	"github.com/kasefuchs/murmora/api/proto/murmora/token/v1"
	"github.com/kasefuchs/murmora/api/proto/murmora/user/v1"
	"github.com/matthewhartstonge/argon2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthenticationServiceServer struct {
	authentication.UnimplementedAuthenticationServiceServer

	argon                *argon2.Config
	userServiceClient    user.UserServiceClient
	tokenServiceClient   token.TokenServiceClient
	sessionServiceClient session.SessionServiceClient
}

func (s *AuthenticationServiceServer) Register(ctx context.Context, request *authentication.AuthenticationRegisterRequest) (*authentication.AuthenticationTokenResponse, error) {
	if err := request.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	hashEncoded, err := s.argon.HashEncoded([]byte(request.Password))
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "hash encoding failed: %s", err.Error())
	}

	userDataResponse, err := s.userServiceClient.CreateUser(ctx, &user.CreateUserRequest{
		Name:         request.Name,
		Email:        request.Email,
		PasswordHash: string(hashEncoded),
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create user: %s", err.Error())
	}

	createSessionResponse, err := s.sessionServiceClient.CreateSession(ctx, &session.CreateSessionRequest{
		UserId: userDataResponse.Id,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create session: %s", err.Error())
	}

	return &authentication.AuthenticationTokenResponse{
		Token: createSessionResponse.Token,
	}, nil
}

func (s *AuthenticationServiceServer) Login(ctx context.Context, request *authentication.AuthenticationLoginRequest) (*authentication.AuthenticationTokenResponse, error) {
	if err := request.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	userDataResponse, err := s.userServiceClient.GetUser(ctx, &user.GetUserRequest{
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

	createSessionResponse, err := s.sessionServiceClient.CreateSession(ctx, &session.CreateSessionRequest{
		UserId: userDataResponse.Id,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to create session: %s", err.Error())
	}

	return &authentication.AuthenticationTokenResponse{
		Token: createSessionResponse.Token,
	}, nil
}

func NewAuthenticationServiceServer(
	userServiceClient user.UserServiceClient,
	sessionServiceClient session.SessionServiceClient,
) *AuthenticationServiceServer {
	argon := argon2.DefaultConfig()

	return &AuthenticationServiceServer{
		argon:                &argon,
		userServiceClient:    userServiceClient,
		sessionServiceClient: sessionServiceClient,
	}
}
