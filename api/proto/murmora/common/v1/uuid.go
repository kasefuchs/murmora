// Copyright (c) Kasefuchs
// SPDX-License-Identifier: MPL-2.0

package common

import "github.com/google/uuid"

func NewUUID(id uuid.UUID) *UUID {
	return &UUID{
		Value: id.String(),
	}
}

func (u *UUID) ToUUID() (uuid.UUID, error) {
	return uuid.Parse(u.Value)
}

func (u *UUID) MustToUUID() uuid.UUID {
	return uuid.MustParse(u.Value)
}
