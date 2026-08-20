package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3" //  есть wiper, goenv
)

type Config struct {
	Database DatabaseConfig
	Redis    RedisConfig
}

type DatabaseConfig struct {
	Host         string
	Port         uint16
	User         string
	Password     string
	DatabaseName string `yaml:"databaseName"`
}

type RedisConfig struct {
	Host     string
	User     string
	Password string
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config error: %w", err)
	}
	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("unmarshal config error: %w", err)
	}
	return &config, nil
}
