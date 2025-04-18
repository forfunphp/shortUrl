package config

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/spf13/pflag"
	"log"
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

	if os.Getenv("FILE_STORAGE_PATH") != "" {
		c.EnvFilePath = os.Getenv("FILE_STORAGE_PATH")
	}

	if c.Databes != "" {
		//	c.Databes = os.Getenv("DATABASE_DSN")
		log.Println("3333333333333333")
		log.Println("3333333333333333")

		db, err := sql.Open("postgres", c.Databes) // Замените "postgres" именем вашего драйвера
		if err != nil {
			log.Printf("не удалось открыть базу данных: %v", err)
		}

		var tableExists bool
		err = db.QueryRow(`
  SELECT EXISTS (
   SELECT 1
   FROM   pg_catalog.pg_tables
   WHERE  schemaname = 'public'
   AND    tablename = 'short_urls'
  );
 `).Scan(&tableExists)
		if err != nil {
			log.Fatalf("Failed to check table existence: %v", err)
		}

		if tableExists {
			log.Println("Table 'short_urls' already exists")
		} else {
			log.Println("Table 'short_urls' does not exist, creating it")
			_, err = db.Exec(`
   CREATE TABLE short_urls (
    shortURL VARCHAR(255) PRIMARY KEY,
    parsedURL TEXT NOT NULL
   )
  `)
			if err != nil {
				log.Fatalf("Failed to create table: %v", err)
			}
		}

		fmt.Println("Подключение к базе данных успешно!")
	} else {
		log.Println(c.Databes)
		log.Println(c.EnvFilePath)
		log.Println("Переменная окружения DATABASE_DSN не установлена. Используется конфигурация по умолчанию (если есть).")
	}

	port, err := strconv.Atoi(c.HTTPAddr[len(c.HTTPAddr)-4:])
	if err != nil || port < 0 || port > 65535 {
		return fmt.Errorf("недопустимый порт HTTP-сервера: %s", c.HTTPAddr)
	}

	BaseURL, err := url.Parse(c.BaseURL)
	if err != nil {
		return fmt.Errorf("ошибка при парсинге базового URL: %w", err)
	}

	if BaseURL.Port() != "" {
		port, err = strconv.Atoi(BaseURL.Port())
	}

	if err != nil || port < 0 || port > 65535 {
		return fmt.Errorf("недопустимый порт базового сервера: %s", c.BaseURL)
	}

	return nil
}
