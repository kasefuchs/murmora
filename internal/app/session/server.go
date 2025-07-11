// Copyright (c) Kasefuchs
// SPDX-License-Identifier: MPL-2.0

package session

import (
	"context"

	"github.com/google/uuid"
	"github.com/kasefuchs/murmora/api/proto/murmora/session/v1"
	"github.com/kasefuchs/murmora/api/proto/murmora/token/v1"
	"github.com/kasefuchs/murmora/api/proto/murmora/user/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/anypb"
)

type Server struct {
	session.UnimplementedSessionServiceServer

	repository  *Repository
	userClient  user.UserServiceClient
	tokenClient token.TokenServiceClient
}

func NewServer(
	repository *Repository,
	userClient user.UserServiceClient,
	tokenClient token.TokenServiceClient,
) *Server {
	return &Server{
		repository:  repository,
		userClient:  userClient,
		tokenClient: tokenClient,
	}
}

func (s *Server) CreateSession(ctx context.Context, request *session.CreateSessionRequest) (*session.CreateSessionResponse, error) {
	if err := request.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	userId, err := uuid.Parse(request.UserId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Failed to parse userClient id: %v", err)
	}

	userDataResponse, err := s.userClient.GetUser(ctx, &user.GetUserRequest{
		Query: &user.GetUserRequest_Id{
			Id: userId.String(),
		},
	})
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Failed to get userClient: %v", err)
	}

	id, err := uuid.NewV7()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to generate session ID: %v", err)
	}

	tokenPayload, err := anypb.New(&session.TokenPayload{
		SessionId: id.String(),
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to create session payload: %v", err)
	}

	tokenResponse, err := s.tokenClient.CreateToken(ctx, &token.CreateTokenRequest{
		Secret:  []byte(userDataResponse.PasswordHash),
		Payload: tokenPayload,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to create tokenClient: %v", err)
	}

	tokenId, err := uuid.Parse(tokenResponse.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Failed to parse tokenClient id: %v", err)
	}

	entity, err := s.repository.Create(&Session{
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
