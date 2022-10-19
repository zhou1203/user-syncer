package db

import (
	"fmt"
	"os"
)

type Config struct {
	Host         string
	DatabaseName string
	Password     string
	User         string
}

const (
	envHost     = "DATABASE_HOST"
	envName     = "DATABASE_NAME"
	envPassword = "DATABASE_PASSWORD"
	envUser     = "DATABASE_USER"
)

func NewConfigFromEnv() (*Config, error) {
	config := &Config{
		Host:         os.Getenv(envHost),
		DatabaseName: os.Getenv(envName),
		Password:     os.Getenv(envPassword),
		User:         os.Getenv(envUser),
	}

	if config.User == "" {
		config.User = "root"
	}

	if config.Host == "" || config.DatabaseName == "" || config.Password == "" {
		return nil, fmt.Errorf("not enough parameter")
	}
	return config, nil
}
