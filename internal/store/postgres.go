package store

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	generator "go-url-short/internal/shorten"
	"log"
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

func (s PostgresStore) DbClose() {
	s.Log.Println("Closing database connection")
	s.db.Close()
}

func (s PostgresStore) Get(shortKey string) (string, error) {
	row := s.db.QueryRow("SELECT * FROM shorturl WHERE short = $1", shortKey)
	var id int
	var originalURL string
	err := row.Scan(&id, &shortKey, &originalURL)
	s.Log.Println("originalURL: ", originalURL)

	if err != nil {
		return "", ErrKeyNotFound{err}
	}

	return originalURL, nil
}

func (s PostgresStore) Set(originalURL string) (string, error) {
	// Check if the key already exists using orignalURL
	var shortkey string
	err := s.db.QueryRow("SELECT short FROM shorturl WHERE url = $1 LIMIT 1", originalURL).Scan(&shortkey)

	if err != nil && err != sql.ErrNoRows {
		s.Log.Println("Error checking if key exists: ", err)
		return "", ErrKeyAlreadyExists{err}
	}

	if shortkey != "" {
		s.Log.Println("Key already exists")
		return shortkey, nil
	}

	shortKey := generator.GenerateShortKey()
	err = s.db.QueryRow("INSERT INTO shorturl (url, short) VALUES ($1,$2) RETURNING short", originalURL, shortKey).Scan(&shortKey)
	if err != nil {
		s.Log.Println("Error inserting into database: ", err)
		return "", ErrKeyAlreadyExists{err}
	}

	return shortKey, nil
}
