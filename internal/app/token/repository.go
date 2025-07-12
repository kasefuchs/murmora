// Copyright (c) Kasefuchs
// SPDX-License-Identifier: MPL-2.0

package token

import "github.com/kasefuchs/murmora/internal/pkg/database"

type Repository struct {
	*database.Repository[Token]
}

func NewRepository(db *database.Database) *Repository {
	return &Repository{
		Repository: &database.Repository[Token]{
			Database: db,
		},
	}
}
