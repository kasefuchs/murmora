// Copyright (c) Kasefuchs
// SPDX-License-Identifier: MPL-2.0

package database

type Config struct {
	Type string `mapstructure:"type"`
	DSN  string `mapstructure:"dsn"`
}
