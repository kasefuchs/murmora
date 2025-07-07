// Copyright (c) Kasefuchs
// SPDX-License-Identifier: MPL-2.0

package data

import (
	"errors"

	"github.com/kasefuchs/murmora/internal/pkg/database"
	"gorm.io/gorm"
)

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

func (r *UserRepository) findOneByCondition(conds ...interface{}) (*User, error) {
	var user User
	if err := r.database.DB.First(&user, conds...).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByID(id string) (*User, error) {
	return r.findOneByCondition("id = ?", id)
}

func (r *UserRepository) FindByName(name string) (*User, error) {
	return r.findOneByCondition("name = ?", name)
}

func (r *UserRepository) FindByEmail(email string) (*User, error) {
	return r.findOneByCondition("email = ?", email)
}
