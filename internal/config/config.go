package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	"os"
)

type Config struct {
	HTTPConfig     HTTPConfig     `yaml:"http" validate:"required"`
	RedisConfig    RedisConfig    `yaml:"redis" validate:"required"`
	PostgresConfig PostgresConfig `yaml:"postgres" validate:"required"`
}

type HTTPConfig struct {
	Addr string `yaml:"addr" validate:"required"`
}

type RedisConfig struct {
	Addr           string `yaml:"addr" validate:"required"`
	DataBase       int    `yaml:"database" validate:"required"`
	PasswordEnvVar string `yaml:"password_env_var" validate:"required"`
	CacheLifetime  string `yaml:"cache_lifetime" validate:"required"`
}

type PostgresConfig struct {
	Addr           string `yaml:"addr" validate:"required"`
	UserEnvVar     string `yaml:"user_env_var" validate:"required"`
	PasswordEnvVar string `yaml:"password_env_var" validate:"required"`
	DataBaseName   string `yaml:"db_name" validate:"required"`
	SSLMode        string `yaml:"ssl_mode" validate:"required"`
	ConnLifetime   string `yaml:"conn_lifetime" validate:"required"`
	MaxOpenConn    int    `yaml:"max_open_conn" validate:"required"`
	MaxIdleConn    int    `yaml:"max_idle_conn" validate:"required"`
	MigrationPath  string `yaml:"migration_path" validate:"required"`
}

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) Load(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	if err := yaml.Unmarshal(data, c); err != nil {
		return fmt.Errorf("failed to parse config file: %w", err)
	}
	return nil
}
