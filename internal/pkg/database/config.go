// Copyright (c) Kasefuchs
// SPDX-License-Identifier: MPL-2.0

package database

type Config struct {
	Type string `koanf:"type"`
	DSN  string `koanf:"dsn"`
}
