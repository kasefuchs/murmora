// Copyright (c) Kasefuchs
// SPDX-License-Identifier: MPL-2.0

package user

import "github.com/kasefuchs/murmora/internal/pkg/database"

type Repository struct {
	*database.Repository[User]
}

func NewRepository(db *database.Database) *Repository {
	return &Repository{
		Repository: &database.Repository[User]{
			Database: db,
		},
	}
}

func (r *Repository) FindByName(name string) (*User, error) {
	return r.FindOneByCondition("name = ?", name)
}

func (r *Repository) FindByEmail(email string) (*User, error) {
	return r.FindOneByCondition("email = ?", email)
}
