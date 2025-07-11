// Copyright (c) Kasefuchs
// SPDX-License-Identifier: MPL-2.0

package server

type Config struct {
	Address            string `koanf:"address"`
	Reflection         bool   `koanf:"reflection"`
	DisableHealthCheck bool   `koanf:"disable_health_check"`
}
