package db

import (
	"fmt"
	"os"

	"github.com/spf13/pflag"
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

func (c *Config) Flags() *pflag.FlagSet {
	fs := pflag.NewFlagSet("database config", pflag.ContinueOnError)
	fs.StringVar(&c.Host, "database-host", c.Host, "database host address")
	fs.StringVar(&c.DatabaseName, "database-name", c.DatabaseName, "database name")
	fs.StringVar(&c.Password, "database-password", c.Password, "database password")
	fs.StringVar(&c.User, "database-user", c.User, "database user")
	return fs
}
