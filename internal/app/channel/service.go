// Copyright (c) Kasefuchs
// SPDX-License-Identifier: MPL-2.0

package channel

import (
	"context"

	"github.com/google/uuid"
	"github.com/kasefuchs/murmora/api/proto/murmora/channel/v1"
	"github.com/kasefuchs/murmora/api/proto/murmora/common/v1"
	"github.com/kasefuchs/murmora/internal/pkg/database"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	channel.UnimplementedChannelServiceServer

	repository *Repository
}

func NewServer(db *database.Database) *Server {
	repository := NewRepository(db)

	return &Server{
		repository: repository,
	}
}

func (s *Server) CreateChannel(_ context.Context, request *channel.CreateChannelRequest) (*channel.CreateChannelResponse, error) {
	if err := request.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	id, err := uuid.NewV7()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate UUID: %v", err)
	}

	entity, err := s.repository.Create(&Channel{
		ID:    id,
		Type:  request.Type,
		Name:  request.Name,
		Flags: common.NewTypedBitField[channel.ChannelFlag](request.Flags).ToFlagSet(),
	})
	if err != nil {
		return nil, status.Errorf(codes.AlreadyExists, "channel already exists: %v", err)
	}

	return &channel.CreateChannelResponse{
		Channel: &channel.Channel{
			Id:    common.NewUUID(entity.ID),
			Type:  request.Type,
			Name:  request.Name,
			Flags: request.Flags,
		},
	}, nil
}

func (s *Server) GetChannel(_ context.Context, request *channel.GetChannelRequest) (*channel.GetChannelResponse, error) {
	if err := request.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	entity, err := s.repository.FindByID(request.Id.Value)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error querying channel: %v", err)
	}
	if entity == nil {
		return nil, status.Errorf(codes.NotFound, "channel not found")
	}

	return &channel.GetChannelResponse{
		Channel: &channel.Channel{
			Id:    common.NewUUID(entity.ID),
			Type:  entity.Type,
			Name:  entity.Name,
			Flags: common.BitFieldFromFlagSet(entity.Flags),
		},
	}, nil
}
