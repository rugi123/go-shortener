package config

import (
	"fmt"
	"net/url"
	"os"

	"gopkg.in/yaml.v3"
)

type AppConfig struct {
	Env       string `yaml:"env"`
	Port      string `yaml:"port"`
	URLLength int    `yaml:"url_length"`
}
type PostgresConfig struct {
	Host      string `yaml:"host"`
	Port      int    `yaml:"port"`
	User      string `yaml:"user"`
	Password  string `yaml:"password"`
	DBName    string `yaml:"dbname"`
	TableName string `yaml:"tablename"`
	SSLMode   string `yaml:"sslmode"`
}
type Config struct {
	App      AppConfig
	Postgres PostgresConfig
}

func (c *PostgresConfig) DSN() string {
	user := url.QueryEscape(c.User)
	password := url.QueryEscape(c.Password)
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%v/%s?sslmode=%s",
		user,
		password,
		c.Host,
		c.Port,
		c.DBName,
		c.SSLMode,
	)
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения кофнига: %s", err)
	}
	cfg := Config{}
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("ошибка парсинга конфига: %s", err)
	}
	return &cfg, nil
}
