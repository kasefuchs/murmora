// Copyright (c) Kasefuchs
// SPDX-License-Identifier: MPL-2.0

package session

import "github.com/google/uuid"

type Session struct {
	ID      uuid.UUID `gorm:"primaryKey"`
	UserID  uuid.UUID `gorm:"not null"`
	TokenID uuid.UUID `gorm:"uniqueIndex"`
}
