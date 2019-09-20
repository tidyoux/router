package config

import (
	"github.com/tidyoux/goutils/viper"
)

type Config struct {
	DSN  string
	Port string
}

func NewConfig() *Config {
	return &Config{
		DSN:  viper.GetString("dsn", ""),
		Port: viper.GetString("port", ":8080"),
	}
}
