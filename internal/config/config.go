package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type AppConfig struct {
	Env       string `yaml:"env"`
	Port      int    `yaml:"port"`
	URLLength int    `yaml:"url_length"`
}
type PostgresConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
	SSLMode  string `yaml:"sslmode"`
}
type Config struct {
	App      AppConfig
	Postgres PostgresConfig
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения кофнига:", err)
	}
	cfg := Config{}
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("ошибка парсинга конфига:", err)
	}
	return &cfg, nil
}
