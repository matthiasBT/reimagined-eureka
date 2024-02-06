package config

import (
	"errors"
	"flag"
	"net/url"
)

var ErrEmptyDatabase = errors.New("database path must be set explicitly")
var ErrEmptyServerAddress = errors.New("server path must be set explicitly")
var ErrInvalidServerAddress = errors.New("server URL scheme and host must be set explicitly")

type Config struct {
	DatabasePath  string
	ServerURL     *url.URL
	serverAddress string
}

func InitConfig() (*Config, error) {
	conf := new(Config)
	flag.StringVar(&conf.DatabasePath, "d", "", "Path to a secrets storage")
	flag.StringVar(&conf.serverAddress, "a", "", "Server address")
	flag.Parse()
	if err := conf.validateDB(); err != nil {
		return nil, err
	}
	serverURL, err := conf.validateServerAddress()
	if err != nil {
		return nil, err
	}
	conf.ServerURL = serverURL
	return conf, nil
}

func (c *Config) validateDB() error {
	if c.DatabasePath == "" {
		return ErrEmptyDatabase
	}
	return nil
}

func (c *Config) validateServerAddress() (*url.URL, error) {
	if c.serverAddress == "" {
		return nil, ErrEmptyServerAddress
	}
	parsedURL, err := url.Parse(c.serverAddress)
	if err != nil {
		return nil, err
	}
	if parsedURL.Scheme == "" || parsedURL.Host == "" {
		return nil, ErrInvalidServerAddress
	}
	return parsedURL, nil
}
