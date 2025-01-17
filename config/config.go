package config

import (
	"fmt"
	"github.com/spf13/pflag"
	"net/url"
	"os"
	"strconv"
)

type Config struct {
	HTTPAddr       string
	BaseURL        string
	EnvFilePath    string
	Databes        string
	UseFileStorage bool
}

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) Init() error {
	c.Databes = ""
	c.EnvFilePath = "urls.json"
	c.UseFileStorage = false

	c.HTTPAddr = os.Getenv("SERVER_ADDRESS")
	c.BaseURL = os.Getenv("BASE_URL")
	c.EnvFilePath = os.Getenv("FILE_STORAGE_PATH")

	if os.Getenv("DATABASE_DSN") != "" {
		c.Databes = os.Getenv("DATABASE_DSN")
		c.UseFileStorage = false // Отключаем использование файла, если есть БД
	}

	pflag.StringVarP(&c.HTTPAddr, "http-addr", "a", "localhost:8080", "адрес прослушивания HTTP-сервера")
	pflag.StringVarP(&c.BaseURL, "base-url", "b", "http://localhost:8080", "базовый адрес для сокращенных URL")
	pflag.StringVarP(&c.EnvFilePath, "f", "f", "urls.json", "Путь к файлу для хранения URL")
	pflag.StringVar(&c.Databes, "d", c.Databes, "Строка с адресом подключения к БД")
	pflag.Parse()

	// Логика для определения источника хранения
	if c.Databes == "" {
		if c.EnvFilePath != "" && c.EnvFilePath != "urls.json" {
			c.UseFileStorage = true
		} else {
			c.UseFileStorage = false // По умолчанию используем память
		}
	} else {
		c.UseFileStorage = false // Если есть БД, то не используем файл
	}

	fmt.Printf("HTTPAddr: %s\n", c.HTTPAddr)
	fmt.Printf("BaseURL: %s\n", c.BaseURL)
	fmt.Printf("EnvFilePath: %s\n", c.EnvFilePath)
	fmt.Printf("Databes: %s\n", c.Databes)
	fmt.Printf("UseFileStorage: %t\n", c.UseFileStorage)

	//if os.Getenv("FILE_STORAGE_PATH") != "" {
	//	c.EnvFilePath = os.Getenv("FILE_STORAGE_PATH")
	//}

	//if os.Getenv("DATABASE_DSN") != "" {
	//	c.Databes = os.Getenv("DATABASE_DSN")
	//}

	port, err := strconv.Atoi(c.HTTPAddr[len(c.HTTPAddr)-4:])
	if err != nil || port < 0 || port > 65535 {
		return fmt.Errorf("недопустимый порт HTTP-сервера: %s", c.HTTPAddr)
	}

	parsedURL, err := url.Parse(c.BaseURL)
	if err != nil {
		return fmt.Errorf("ошибка при парсинге базового URL: %w", err)
	}

	if parsedURL.Port() != "" {
		port, err = strconv.Atoi(parsedURL.Port())
	}

	if err != nil || port < 0 || port > 65535 {
		return fmt.Errorf("недопустимый порт базового сервера: %s", c.BaseURL)
	}

	return nil
}
