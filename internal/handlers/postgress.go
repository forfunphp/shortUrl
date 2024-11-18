package handlers

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore(dsn string) (*PostgresStore, error) {

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if dsn == "" {
		dsnPtr := flag.String("d", "", "MySQL DSN (database source name)")
		flag.Parse()
		dsn = *dsnPtr
	}

	logger5, _ := zap.NewDevelopment()
	defer logger5.Sync()

	logger5.Info("Request processed33ffdd-----d",
		zap.String("dsn", dsn),
	)

	defer db.Close()

	_, err = db.Exec(`
 CREATE TABLE movies (
  id SERIAL PRIMARY KEY,
  title VARCHAR(250) NOT NULL DEFAULT '',
  created TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  tags TEXT,
  views INTEGER NOT NULL DEFAULT 0
 )
`)
	if err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	return &PostgresStore{db: db}, nil
}

func (s *PostgresStore) Shorten(longURL string) (string, error) {
	shortURL := uuid.New().String()[:6] //Генерируем короткий URL (6 символов)

	_, err := s.db.Exec("INSERT INTO short_urls (id, short_url, long_url) VALUES ($1, $2, $3)", uuid.New(), shortURL, longURL)
	if err != nil {
		return "", fmt.Errorf("failed to insert short URL: %w", err)
	}
	return shortURL, nil
}

func (s *PostgresStore) GetLongURL(shortURL string) (string, error) {
	var longURL string
	err := s.db.QueryRow("SELECT long_url FROM short_urls WHERE short_url = $1", shortURL).Scan(&longURL)
	if err != nil {
		return "", fmt.Errorf("failed to get long URL: %w", err)
	}
	return longURL, nil
}

func (s *PostgresStore) Close() error {
	return s.db.Close()
}
