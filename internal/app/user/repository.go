// Copyright (c) Kasefuchs
// SPDX-License-Identifier: MPL-2.0

package user

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

func (r *Repository) Create(user *User) (*User, error) {
	if err := r.database.DB.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (r *Repository) findOneByCondition(conds ...interface{}) (*User, error) {
	var user User
	if err := r.database.DB.First(&user, conds...).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *Repository) FindByID(id string) (*User, error) {
	return r.findOneByCondition("id = ?", id)
}

func (r *Repository) FindByName(name string) (*User, error) {
	return r.findOneByCondition("name = ?", name)
}

func (r *Repository) FindByEmail(email string) (*User, error) {
	return r.findOneByCondition("email = ?", email)
}
