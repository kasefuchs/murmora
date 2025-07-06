// Copyright (c) Kasefuchs
// SPDX-License-Identifier: MPL-2.0

package config

import (
	"github.com/kasefuchs/murmora/internal/pkg/database"
	"github.com/spf13/viper"
)

type Config struct {
	Database         *database.Config `mapstructure:"database"`
	ListenAddress    string           `mapstructure:"listen_address"`
	EnableReflection bool             `mapstructure:"enable_reflection"`
}

func Load() (*Config, error) {
	viper.SetConfigName("user")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("configs")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	cfg := &Config{}
	err := viper.Unmarshal(cfg)

	return cfg, err
}
