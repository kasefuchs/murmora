// Copyright (c) Kasefuchs
// SPDX-License-Identifier: MPL-2.0

package config

import "github.com/spf13/viper"

type Config struct {
	ListenAddress     string `mapstructure:"listen_address"`
	EnableReflection  bool   `mapstructure:"enable_reflection"`
	UserServiceUrl    string `mapstructure:"user_service_url"`
	SessionServiceUrl string `mapstructure:"session_service_url"`
}

func Load() (*Config, error) {
	viper.SetConfigName("authentication")
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
