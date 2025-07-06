// Copyright (c) Kasefuchs
// SPDX-License-Identifier: MPL-2.0

package database

import (
	"errors"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Database struct {
	DB *gorm.DB
}

// New opens the data connection.
func New(config *Config) (*Database, error) {
	dial, err := openDialector(config)
	if err != nil {
		return nil, err
	}

	db, err := gorm.Open(dial)
	if err != nil {
		return nil, err
	}

	return &Database{DB: db}, nil
}

// Migrate runs auto migration.
func (d *Database) Migrate(models ...interface{}) error {
	return d.DB.AutoMigrate(models...)
}

// Close closes the data connection.
func (d *Database) Close() error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// openDialector returns dialector for specified configuration.
func openDialector(config *Config) (gorm.Dialector, error) {
	switch config.Type {
	case "postgres":
		return postgres.Open(config.DSN), nil
	case "sqlite":
		return sqlite.Open(config.DSN), nil
	default:
		return nil, errors.New("unsupported data type")
	}
}
