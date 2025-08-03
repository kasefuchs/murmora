// Copyright (c) Kasefuchs
// SPDX-License-Identifier: MPL-2.0

package channel

import (
	"github.com/google/uuid"
	"github.com/kasefuchs/murmora/api/proto/murmora/channel/v1"
	"github.com/kasefuchs/murmora/internal/pkg/bitflag"
)

type Channel struct {
	ID    uuid.UUID                             `gorm:"primaryKey"`
	Type  channel.ChannelType                   `gorm:"not null"`
	Name  string                                `gorm:"not null"`
	Flags *bitflag.FlagSet[channel.ChannelFlag] `gorm:"not null"`
}
