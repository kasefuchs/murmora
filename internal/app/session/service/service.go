// Copyright (c) Kasefuchs
// SPDX-License-Identifier: MPL-2.0

package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/kasefuchs/murmora/api/proto/murmora/session/v1"
	"github.com/kasefuchs/murmora/api/proto/murmora/token/v1"
	"github.com/kasefuchs/murmora/api/proto/murmora/user/v1"
	"github.com/kasefuchs/murmora/internal/app/session/data"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/anypb"
)

type SessionServiceServer struct {
	session.UnimplementedSessionServiceServer

	userServiceClient  user.UserServiceClient
	tokenServiceClient token.TokenServiceClient

	sessionRepository *data.SessionRepository
}

func (s *SessionServiceServer) CreateSession(ctx context.Context, request *session.CreateSessionRequest) (*session.CreateSessionResponse, error) {
	if err := request.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	userId, err := uuid.Parse(request.UserId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Failed to parse user id: %v", err)
	}

	userDataResponse, err := s.userServiceClient.GetUser(ctx, &user.GetUserRequest{
		Query: &user.GetUserRequest_Id{
			Id: userId.String(),
		},
	})
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Failed to get user: %v", err)
	}

	id, err := uuid.NewV7()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to generate session ID: %v", err)
	}

	tokenPayload, err := anypb.New(&session.SessionTokenPayload{
		SessionId: id.String(),
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to create session payload: %v", err)
	}

	tokenResponse, err := s.tokenServiceClient.CreateToken(ctx, &token.CreateTokenRequest{
		Secret:  []byte(userDataResponse.PasswordHash),
		Payload: tokenPayload,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to create token: %v", err)
	}

	tokenId, err := uuid.Parse(tokenResponse.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Failed to parse token id: %v", err)
	}

	entity, err := s.sessionRepository.Create(&data.Session{
		ID:      id,
		UserID:  userId,
		TokenID: tokenId,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to create session: %v", err)
	}

	return &session.CreateSessionResponse{
		Id:    entity.ID.String(),
		Token: tokenResponse.Token,
	}, nil
}

func NewSessionServiceServer(
	sessionRepository *data.SessionRepository,
	userServiceClient user.UserServiceClient,
	tokenServiceClient token.TokenServiceClient,
) *SessionServiceServer {
	return &SessionServiceServer{
		sessionRepository:  sessionRepository,
		userServiceClient:  userServiceClient,
		tokenServiceClient: tokenServiceClient,
	}
}
