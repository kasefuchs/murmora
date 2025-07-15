// Copyright (c) Kasefuchs
// SPDX-License-Identifier: MPL-2.0

package database

import (
	"errors"

	"gorm.io/gorm"
)

type Repository[T any] struct {
	Database *Database
}

func (r *Repository[T]) Create(entity *T) (*T, error) {
	if err := r.Database.DB.Create(entity).Error; err != nil {
		return nil, err
	}
	return entity, nil
}

func (r *Repository[T]) FindByID(id string) (*T, error) {
	var entity T
	if err := r.Database.DB.First(&entity, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &entity, nil
}

func (r *Repository[T]) FindOneByCondition(args ...interface{}) (*T, error) {
	var entity T
	if err := r.Database.DB.First(&entity, args...).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &entity, nil
}
