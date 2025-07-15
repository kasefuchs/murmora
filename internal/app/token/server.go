// Copyright (c) Kasefuchs
// SPDX-License-Identifier: MPL-2.0

package token

import (
	"context"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/kasefuchs/murmora/api/proto/murmora/common/v1"
	"github.com/kasefuchs/murmora/api/proto/murmora/token/v1"
	"github.com/kasefuchs/murmora/internal/pkg/database"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

type Server struct {
	token.UnimplementedTokenServiceServer

	jwtParser  *jwt.Parser
	repository *Repository
}

func NewServer(db *database.Database) *Server {
	jwtParser := jwt.NewParser()
	repository := NewRepository(db)

	return &Server{
		jwtParser:  jwtParser,
		repository: repository,
	}
}

func (s *Server) CreateToken(_ context.Context, request *token.CreateTokenRequest) (*token.CreateTokenResponse, error) {
	if err := request.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	id, err := uuid.NewV7()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate UUID: %v", err)
	}

	jwtToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{ID: id.String()}).SignedString(request.Secret)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to sign JWT: %v", err)
	}

	payload, err := proto.Marshal(request.Payload)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to marshal payload: %v", err)
	}

	if _, err := s.repository.Create(&Token{
		ID:      id,
		Payload: payload,
	}); err != nil {
		return nil, status.Errorf(codes.Unavailable, "failed to create token: %v", err)
	}

	return &token.CreateTokenResponse{
		Id:    common.NewUUID(id),
		Token: &common.JWT{Value: jwtToken},
	}, nil
}

func (s *Server) resolveTokenData(id string) (*token.TokenDataResponse, error) {
	entity, err := s.repository.FindByID(id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to find token: %v", err)
	}
	if entity == nil {
		return nil, status.Errorf(codes.NotFound, "token not found")
	}

	var payload anypb.Any
	if err := proto.Unmarshal(entity.Payload, &payload); err != nil {
		return nil, status.Errorf(codes.Aborted, "failed to unmarshal payload: %v", err)
	}

	return &token.TokenDataResponse{
		Id:      common.NewUUID(entity.ID),
		Payload: &payload,
	}, nil
}

func (s *Server) GetTokenData(_ context.Context, request *token.GetTokenDataRequest) (*token.TokenDataResponse, error) {
	if err := request.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	jwtToken, _, err := s.jwtParser.ParseUnverified(request.Token.Value, &jwt.RegisteredClaims{})
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to parse JWT: %v", err)
	}

	claims, ok := jwtToken.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "failed to parse JWT claims: %v", err)
	}

	return s.resolveTokenData(claims.ID)
}

func (s *Server) GetValidatedTokenData(_ context.Context, request *token.GetValidatedTokenDataRequest) (*token.TokenDataResponse, error) {
	if err := request.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	jwtToken, err := s.jwtParser.ParseWithClaims(request.Token.Value, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return request.Secret, nil
	})
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to parse JWT: %v", err)
	}

	claims, ok := jwtToken.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "failed to parse JWT claims: %v", err)
	}

	return s.resolveTokenData(claims.ID)
}
