// Copyright (c) Kasefuchs
// SPDX-License-Identifier: MPL-2.0

package session

import (
	"context"

	"github.com/google/uuid"
	"github.com/kasefuchs/murmora/api/proto/murmora/common/v1"
	"github.com/kasefuchs/murmora/api/proto/murmora/session/v1"
	"github.com/kasefuchs/murmora/api/proto/murmora/token/v1"
	"github.com/kasefuchs/murmora/api/proto/murmora/user/v1"
	"github.com/kasefuchs/murmora/internal/pkg/database"
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
	db *database.Database,
	userClient user.UserServiceClient,
	tokenClient token.TokenServiceClient,
) *Server {
	repository := NewRepository(db)

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

	userData, err := s.userClient.GetUser(ctx, &user.GetUserRequest{
		Query: &user.GetUserRequest_Id{Id: request.UserId},
	})
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "failed to get user: %v", err)
	}

	userId, err := userData.User.Id.ToUUID()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to parse user id: %v", err)
	}

	id, err := uuid.NewV7()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate session ID: %v", err)
	}

	payload, err := anypb.New(&session.TokenPayload{Id: common.NewUUID(id)})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create payload: %v", err)
	}

	tokenData, err := s.tokenClient.CreateToken(ctx, &token.CreateTokenRequest{
		Secret:  userData.User.PasswordHash,
		Payload: payload,
	})
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to create token: %v", err)
	}

	tokenId, err := tokenData.Result.Id.ToUUID()
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to parse token id: %v", err)
	}

	entity, err := s.repository.Create(&Session{
		ID:      id,
		UserID:  userId,
		TokenID: tokenId,
	})
	if err != nil {
		return nil, status.Errorf(codes.Unavailable, "failed to create session: %v", err)
	}

	return &session.CreateSessionResponse{
		Id:    common.NewUUID(entity.ID),
		Token: tokenData.Result.Token,
	}, nil
}

func (s *Server) GetSession(_ context.Context, request *session.GetSessionRequest) (*session.GetSessionResponse, error) {
	if err := request.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	var err error
	var entity *Session

	switch query := request.GetQuery().(type) {
	case *session.GetSessionRequest_Id:
		entity, err = s.repository.FindByID(query.Id.Value)
	case *session.GetSessionRequest_TokenId:
		entity, err = s.repository.FindByTokenId(query.TokenId.Value)
	default:
		return nil, status.Error(codes.InvalidArgument, "invalid query type")
	}

	if err != nil {
		return nil, status.Errorf(codes.Internal, "error querying session: %v", err)
	}
	if entity == nil {
		return nil, status.Error(codes.NotFound, "session not found")
	}

	return &session.GetSessionResponse{
		Id:      common.NewUUID(entity.ID),
		UserId:  common.NewUUID(entity.UserID),
		TokenId: common.NewUUID(entity.TokenID),
	}, nil
}
