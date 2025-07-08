// Copyright (c) Kasefuchs
// SPDX-License-Identifier: MPL-2.0

package service

import (
	"context"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/kasefuchs/murmora/api/proto/murmora/token/v1"
	"github.com/kasefuchs/murmora/internal/app/token/data"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

type TokenServiceServer struct {
	token.UnimplementedTokenServiceServer

	tokenRepository *data.TokenRepository
}

func (s *TokenServiceServer) CreateToken(_ context.Context, request *token.CreateTokenRequest) (*token.CreateTokenResponse, error) {
	if err := request.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	id, err := uuid.NewV7()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to generate UUID: %v", err)
	}

	jwtToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		ID: id.String(),
	}).SignedString(request.Secret)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Failed to sign JWT: %v", err)
	}

	payload, err := proto.Marshal(request.Payload)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Failed to marshal payload: %v", err)
	}

	entity, err := s.tokenRepository.Create(&data.Token{
		ID:      id,
		Payload: payload,
	})
	if err != nil {
		return nil, status.Errorf(codes.AlreadyExists, "User already exists: %v", err)
	}

	return &token.CreateTokenResponse{
		Id:    entity.ID.String(),
		Token: jwtToken,
	}, nil
}

func (s *TokenServiceServer) ValidateToken(_ context.Context, request *token.ValidateTokenRequest) (*token.ValidateTokenResponse, error) {
	if err := request.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	jwtToken, err := jwt.ParseWithClaims(request.Token, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return request.Secret, nil
	})
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Failed to parse JWT: %v", err)
	}

	claims, ok := jwtToken.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "Failed to parse JWT claims: %v", err)
	}

	entity, err := s.tokenRepository.FindByID(claims.ID)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Failed to find token: %v", err)
	}

	var payload anypb.Any
	if err := proto.Unmarshal(entity.Payload, &payload); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Failed to unmarshal payload: %v", err)
	}

	return &token.ValidateTokenResponse{
		Id:      entity.ID.String(),
		Payload: &payload,
	}, nil
}

func NewTokenServiceServer(tokenRepository *data.TokenRepository) *TokenServiceServer {
	return &TokenServiceServer{
		tokenRepository: tokenRepository,
	}
}
