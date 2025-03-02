package config

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/spf13/pflag"
	"log"
	"net/url"
	"os"
	"strconv"
	"time"
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

	if os.Getenv("FILE_STORAGE_PATH") != "" {
		c.EnvFilePath = os.Getenv("FILE_STORAGE_PATH")
	}

	if os.Getenv("DATABASE_DSN") != "" {
		c.Databes = os.Getenv("DATABASE_DSN")

		db, err := sql.Open("postgres", c.Databes) // Замените "postgres" именем вашего драйвера
		if err != nil {
			log.Printf("не удалось открыть базу данных: %w", err) // Записываем ошибку в лог. Используем Printf вместо Fatalf
			os.Exit(1)                                            // Явно завершаем работу после записи в лог
		}
		defer db.Close()

		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		if err := db.PingContext(ctx); err != nil {
			log.Fatalf("не удалось подключиться к базе данных с DSN '%s': %w", c.Databes, err) // Здесь по-прежнему используем Fatalf, так как это критическая ошибка для программы
		}
		fmt.Println("Подключение к базе данных успешно!")
	} else {
		log.Println("Переменная окружения DATABASE_DSN не установлена. Используется конфигурация по умолчанию (если есть).")
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
