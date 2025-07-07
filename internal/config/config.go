package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type AppConfig struct {
	Env       string `yaml:"env"`
	Port      string `yaml:"port"`
	URLLength int    `yaml:"url_length"`
}

type Config struct {
	App      AppConfig
	Postgres PostgresConfig
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %s", err)
	}
	cfg := Config{}
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %s", err)
	}
	return &cfg, nil
}
