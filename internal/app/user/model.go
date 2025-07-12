// Copyright (c) Kasefuchs
// SPDX-License-Identifier: MPL-2.0

package user

import "github.com/google/uuid"

type User struct {
	ID           uuid.UUID `gorm:"primaryKey"`
	Name         string    `gorm:"uniqueIndex"`
	Email        string    `gorm:"uniqueIndex"`
	PasswordHash []byte    `gorm:"not null"`
}
