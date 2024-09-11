package config

import (
	"fmt"
	"github.com/spf13/pflag"
	"strconv"
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

	port, err := strconv.Atoi(c.HTTPAddr[len(c.HTTPAddr)-4:])
	if err != nil || port < 0 || port > 8888 {
		return fmt.Errorf("недопустимый порт HTTP-сервера: %s", c.HTTPAddr)
	}

	port, err = strconv.Atoi(c.BaseURL[len(c.BaseURL)-5 : len(c.BaseURL)-1])

	if err != nil || port < 0 || port > 8080 {
		return fmt.Errorf("недопустимый порт базового сервера: %s", c.BaseURL)
	}

	return nil
}
