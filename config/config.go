package config

import (
	"database/sql"
	"fmt"
	"github.com/spf13/pflag"
	"net/url"
	"os"
	"strconv"
)

type Config struct {
	HTTPAddr    string
	BaseURL     string
	EnvFilePath string
	Databes     string
}

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) Init() error {

	c.HTTPAddr = os.Getenv("SERVER_ADDRESS")
	c.BaseURL = os.Getenv("BASE_URL")
	c.EnvFilePath = os.Getenv("FILE_STORAGE_PATH")
	c.Databes = os.Getenv("DATABASE_DSN")

	pflag.StringVarP(&c.HTTPAddr, "http-addr", "a", "localhost:8080", "адрес прослушивания HTTP-сервера")
	pflag.StringVarP(&c.BaseURL, "base-url", "b", "http://localhost:8080", "базовый адрес для сокращенных URL")
	pflag.StringVarP(&c.EnvFilePath, "f", "f", "urls.json", "Путь к файлу для хранения URL")
	pflag.StringVarP(&c.Databes, "d", "d", "", "Строка с адресом подключения к БД")
	pflag.Parse()

	fmt.Println("34t4tg4g444g4g")
	fmt.Println(c.Databes)

	if c.Databes != "" {

		// создаём соединение с СУБД PostgreSQL с помощью аргумента командной строки
		conn, err := sql.Open("pgx", c.Databes)
		if err != nil {
			fmt.Println("conn222")

			return err
		}

		_, err = conn.Exec(`
		 CREATE TABLE movies (
		  id SERIAL PRIMARY KEY,
		  title VARCHAR(250) NOT NULL DEFAULT '',
		  created TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		  tags TEXT,
		  views INTEGER NOT NULL DEFAULT 0
		 )
		`)
		if err != nil {
			fmt.Errorf("failed to create table: %w", err)
			return nil

		}

		conn.Close()

	}

	fmt.Println("conn2")

	if os.Getenv("FILE_STORAGE_PATH") != "" {
		c.EnvFilePath = os.Getenv("FILE_STORAGE_PATH")
	}

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
