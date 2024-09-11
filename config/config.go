package config

import (
	"github.com/spf13/pflag"
)

type Config struct {
	HTTPAddr string
	BaseURL  string
}

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) Init() error {
	pflag.StringVarP(&c.HTTPAddr, "http-addr", "a", "localhost:8080", "адрес прослушивания HTTP-сервера")
	pflag.StringVarP(&c.BaseURL, "base-url", "b", "http://localhost:8080/", "базовый адрес для сокращенных URL")
	pflag.Parse()

	return nil
}
