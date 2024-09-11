package config

import (
	"fmt"
	"github.com/spf13/pflag"
	"strconv"
	"strings"
)

type Config struct {
	HTTPAddr string
	BaseURL  string
}

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) Init() error {
	pflag.StringVarP(&c.HTTPAddr, "http-addr", "a", "localhost:8888", "адрес прослушивания HTTP-сервера")
	pflag.StringVarP(&c.BaseURL, "base-url", "b", "http://localhost:8000", "базовый адрес для сокращенных URL")
	pflag.Parse()

	port, err := strconv.Atoi(c.HTTPAddr[len(c.HTTPAddr)-4:])
	if err != nil || port < 0 || port > 65535 {
		return fmt.Errorf("недопустимый порт HTTP-сервера: %s", c.HTTPAddr)
	}

	portStr := strings.Split(c.BaseURL, ":")[2]
	port, err = strconv.Atoi(portStr)

	if err != nil || port < 0 || port > 65535 {
		return fmt.Errorf("недопустимый порт URL-адреса сервера: %s", c.BaseURL)
	}

	return nil
}
