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
	var k int64
	err := s.db.QueryRow("SELECT id FROM shorturl WHERE url = $1 LIMIT 1", originalURL).Scan(&k)

	if err != nil && err != sql.ErrNoRows {
		s.Log.Println("Error checking if key exists: ", err)
		return "", ErrKeyAlreadyExists{err}
	}

	if k != 0 {
		s.Log.Println("Key already exists")
		return generator.ConvertRadix62(k), nil
	}

	err = s.db.QueryRow("INSERT INTO shorturl (url, short) VALUES ($1,$2) RETURNING id", originalURL).Scan(&k)
	if err != nil {
		s.Log.Println("Error inserting into database: ", err)
		return "", ErrKeyAlreadyExists{err}
	}
	s.Log.Println("Inserted into database: ", k)

	return generator.ConvertRadix62(k), nil
}
