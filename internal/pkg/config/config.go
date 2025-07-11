// Copyright (c) Kasefuchs
// SPDX-License-Identifier: MPL-2.0

package config

import (
	"os"
	"path/filepath"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/hcl"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/providers/file"
	"github.com/rs/zerolog/log"
)

const delim = "."

type Config[T any] struct {
	Value       *T
	configPath  string
	globalKoanf *koanf.Koanf
}

func New[T any]() *Config[T] {
	return &Config[T]{
		Value:       new(T),
		globalKoanf: koanf.New(delim),
	}
}

func (c *Config[T]) mustUnmarshal() {
	if err := c.globalKoanf.Unmarshal("", c.Value); err != nil {
		log.Fatal().Err(err).Msg("Failed to unmarshal config")
	}
}

func (c *Config[T]) MustLoadDefaults(defaults map[string]any) {
	if err := c.globalKoanf.Load(confmap.Provider(defaults, delim), nil); err != nil {
		log.Fatal().Err(err).Msg("Failed to load default config values")
	}

	c.mustUnmarshal()
}

func (c *Config[T]) MustLoadConfigFile(path string) {
	var err error
	if c.configPath, err = filepath.Abs(path); err != nil {
		log.Fatal().Err(err).Str("path", path).Msg("Invalid config path")
	}

	if _, err := os.Stat(c.configPath); err != nil {
		log.Fatal().Err(err).Str("path", c.configPath).Msg("Config file not found")
	}

	if err := c.globalKoanf.Load(file.Provider(c.configPath), hcl.Parser(true)); err != nil {
		log.Fatal().Err(err).Str("path", c.configPath).Msg("Failed to load config file")
	}

	c.mustUnmarshal()
}

func (c *Config[T]) Sprint() string {
	return c.globalKoanf.Sprint()
}

func (c *Config[T]) ConfigPath() string {
	return c.configPath
}
