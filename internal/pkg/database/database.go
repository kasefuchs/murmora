// Copyright (c) Kasefuchs
// SPDX-License-Identifier: MPL-2.0

package database

import (
	"errors"

	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Database struct {
	DB *gorm.DB
}

// MustNew opens the data connection.
func MustNew(config *Config) *Database {
	dial, err := openDialector(config)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to open dialector")
	}

	db, err := gorm.Open(dial)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to Database")
	}

	return &Database{DB: db}
}

// MustMigrate runs auto migration.
func (d *Database) MustMigrate(models ...interface{}) {
	if err := d.DB.AutoMigrate(models...); err != nil {
		log.Fatal().Err(err).Msg("Failed to migrate Database")
	}
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
