// Copyright (c) Kasefuchs
// SPDX-License-Identifier: MPL-2.0

package user

import (
	"github.com/google/uuid"
	"github.com/kasefuchs/murmora/api/proto/murmora/user/v1"
	"github.com/kasefuchs/murmora/internal/pkg/bitflag"
)

type User struct {
	ID           uuid.UUID                       `gorm:"primaryKey"`
	Name         string                          `gorm:"uniqueIndex"`
	Flags        *bitflag.FlagSet[user.UserFlag] `gorm:"not null"`
	Email        string                          `gorm:"uniqueIndex"`
	PasswordHash []byte                          `gorm:"not null"`
}
