// Copyright (c) Kasefuchs
// SPDX-License-Identifier: MPL-2.0

package data

import "github.com/google/uuid"

type User struct {
	ID           uuid.UUID `gorm:"primaryKey"`
	Name         string    `gorm:"uniqueIndex"`
	Email        string    `gorm:"uniqueIndex"`
	PasswordHash string    `gorm:"not null"`
}
