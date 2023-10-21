package store

import (
	"database/sql"
	"fmt"
	generator "go-url-short/internal/shorten"
	"log"
	"testing"

	_ "github.com/lib/pq"
)

type DatabaseConfig struct {
	Name     string
	Host     string
	Port     int
	User     string
	Password string
}

type PostgresStore struct {
	db  *sql.DB
	Log *log.Logger
}

func NewPostgresStore(config DatabaseConfig) *PostgresStore {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s",
		config.Host, config.Port, config.User, config.Password, config.Name)
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		panic(err)
	}

	return &PostgresStore{
		Log: log.New(log.Writer(), "POSTGRESSTORE:", log.LstdFlags),
		db:  db,
	}
}

func (s PostgresStore) Get(shortKey string) (string, error) {
	defer s.db.Close()
	row := s.db.QueryRow("SELECT * FROM shorturl WHERE short = $1", shortKey)
	var id int
	var originalURL string
	err := row.Scan(&id, &shortKey, &originalURL)
	if err != nil {
		return "", ErrKeyNotFound{err}
	}

	return originalURL, nil
}

func (s PostgresStore) Set(originalURL string) (string, error) {
	defer s.db.Close()
	shortKey := generator.GenerateShortKey()
	err := s.db.QueryRow("INSERT INTO shorturl (url,short) VALUES ($1,$2) RETURNING short_key", originalURL, testing.Short()).Scan(&shortKey)
	if err != nil {
		return "", ErrKeyAlreadyExists{err}
	}

	return shortKey, nil
}
