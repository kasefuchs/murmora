// Copyright (c) Kasefuchs
// SPDX-License-Identifier: MPL-2.0

package token

import (
	"errors"

	"github.com/kasefuchs/murmora/internal/pkg/database"
	"gorm.io/gorm"
)

type Repository struct {
	database *database.Database
}

func NewRepository(database *database.Database) *Repository {
	return &Repository{
		database: database,
	}
}

func (r *Repository) Create(user *Token) (*Token, error) {
	if err := r.database.DB.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (r *Repository) findOneByCondition(conds ...interface{}) (*Token, error) {
	var token Token
	if err := r.database.DB.First(&token, conds...).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &token, nil
}

func (r *Repository) FindByID(id string) (*Token, error) {
	return r.findOneByCondition("id = ?", id)
}
