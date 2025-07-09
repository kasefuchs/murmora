// Copyright (c) Kasefuchs
// SPDX-License-Identifier: MPL-2.0

package data

import (
	"errors"

	"github.com/kasefuchs/murmora/internal/pkg/database"
	"gorm.io/gorm"
)

type SessionRepository struct {
	database *database.Database
}

func NewSessionRepository(database *database.Database) *SessionRepository {
	return &SessionRepository{
		database: database,
	}
}

func (r *SessionRepository) Create(session *Session) (*Session, error) {
	if err := r.database.DB.Create(session).Error; err != nil {
		return nil, err
	}

	return session, nil
}

func (r *SessionRepository) findOneByCondition(conds ...interface{}) (*Session, error) {
	var session Session
	if err := r.database.DB.First(&session, conds...).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &session, nil
}

func (r *SessionRepository) FindByID(id string) (*Session, error) {
	return r.findOneByCondition("id = ?", id)
}
