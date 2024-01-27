package config

import (
	"errors"
	"flag"
)

var ErrEmptyDatabase = errors.New("database path must be set explicitly")
var ErrEmptyServerAddress = errors.New("server path must be set explicitly")

type Config struct {
	DatabasePath  string
	ServerAddress string
}

func InitConfig() (*Config, error) {
	conf := new(Config)
	flag.StringVar(&conf.DatabasePath, "d", "", "Path to a secrets storage") // TODO: SQLCipher
	flag.StringVar(&conf.ServerAddress, "a", "", "Server address")
	flag.Parse()
	if err := conf.validateDB(); err != nil {
		return nil, err
	}
	if err := conf.validateServerAddress(); err != nil {
		return nil, err
	}
	return conf, nil
}

func (c *Config) validateDB() error {
	if c.DatabasePath == "" {
		return ErrEmptyDatabase
	}
	return nil
}

func (c *Config) validateServerAddress() error {
	if c.ServerAddress == "" {
		return ErrEmptyServerAddress
	}
	return nil
}
