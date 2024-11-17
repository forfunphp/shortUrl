package handlers

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore(dsn string) (*PostgresStore, error) {

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	_, err = db.Exec(`
  CREATE TABLE IF NOT EXISTS short_urls (
   id UUID PRIMARY KEY,
   short_url VARCHAR(255) UNIQUE NOT NULL,
   long_url TEXT NOT NULL
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
