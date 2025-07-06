// Copyright (c) Kasefuchs
// SPDX-License-Identifier: MPL-2.0

package data

import "github.com/kasefuchs/murmora/internal/pkg/database"

type UserRepository struct {
	database *database.Database
}

func NewUserRepository(database *database.Database) *UserRepository {
	return &UserRepository{
		database: database,
	}
}

func (r *UserRepository) Create(user *User) (*User, error) {
	if err := r.database.DB.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}
