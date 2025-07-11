// Copyright (c) Kasefuchs
// SPDX-License-Identifier: MPL-2.0

package session

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

func (r *Repository) Create(session *Session) (*Session, error) {
	if err := r.database.DB.Create(session).Error; err != nil {
		return nil, err
	}

	return session, nil
}

func (r *Repository) findOneByCondition(conds ...interface{}) (*Session, error) {
	var session Session
	if err := r.database.DB.First(&session, conds...).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &session, nil
}

func (r *Repository) FindByID(id string) (*Session, error) {
	return r.findOneByCondition("id = ?", id)
}
