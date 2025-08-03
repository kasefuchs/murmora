// Copyright (c) Kasefuchs
// SPDX-License-Identifier: MPL-2.0

package channel

import "github.com/kasefuchs/murmora/internal/pkg/database"

type Repository struct {
	*database.Repository[Channel]
}

func NewRepository(db *database.Database) *Repository {
	return &Repository{
		Repository: &database.Repository[Channel]{
			Database: db,
		},
	}
}
