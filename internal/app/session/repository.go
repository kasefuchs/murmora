// Copyright (c) Kasefuchs
// SPDX-License-Identifier: MPL-2.0

package session

import "github.com/kasefuchs/murmora/internal/pkg/database"

type Repository struct {
	*database.Repository[Session]
}

func NewRepository(db *database.Database) *Repository {
	return &Repository{
		Repository: &database.Repository[Session]{
			Database: db,
		},
	}
}

func (r *Repository) FindByTokenId(tokenId string) (*Session, error) {
	return r.FindOneByCondition("token_id = ?", tokenId)
}
