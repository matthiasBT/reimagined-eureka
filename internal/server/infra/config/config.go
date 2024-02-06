package config

import (
	"flag"
	"fmt"
	"time"

	"github.com/caarlos0/env/v9"
)

const SessionTTL = 1 * time.Hour

type Config struct {
	ServerAddr  string `env:"RUN_ADDRESS"`
	DatabaseDSN string `env:"DATABASE_URI"`
}

func Read() (*Config, error) {
	conf := new(Config)
	flag.StringVar(&conf.ServerAddr, "a", "", "Server address. Usage: -a=host:port")
	flag.StringVar(&conf.DatabaseDSN, "d", "", "PostgreSQL database DSN")
	flag.Parse()
	err := env.Parse(conf)
	if err != nil {
		return nil, err
	}
	if conf.ServerAddr == "" || conf.DatabaseDSN == "" {
		return nil, fmt.Errorf("empty fields in the config")
	}
	return conf, nil
}
