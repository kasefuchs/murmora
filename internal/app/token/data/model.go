// Copyright (c) Kasefuchs
// SPDX-License-Identifier: MPL-2.0

package data

import "github.com/google/uuid"

type Token struct {
	ID      uuid.UUID `gorm:"primaryKey"`
	Payload []byte    `gorm:"not null"`
}
