package store

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	generator "go-url-short/internal/shorten"
	"log"
)

type DatabaseConfig struct {
	Name     string `default:"shorturl"`
	Host     string `default:""`
	Port     int    `default:""`
	User     string `default:""`
	Password string `default:""`
}

type PostgresStore struct {
	db  *sql.DB
	Log *log.Logger
}

func NewPostgresStore(config *DatabaseConfig) *PostgresStore {
	l := log.New(log.Writer(), "POSTGRESSTORE:", log.LstdFlags)
	l.Print("Creating new postgres store")
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s",
		config.Host, config.Port, config.User, config.Password, config.Name)
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		panic(err)
	}

	return &PostgresStore{
		Log: l,
		db:  db,
	}
}

func (s PostgresStore) DbClose() {
	s.Log.Println("Closing database connection")
	s.db.Close()
}

func (s PostgresStore) Get(shortKey string) (string, error) {
	k, err := generator.ConvertRadix10(shortKey)
	if err != nil {
		s.Log.Println("Error converting radix62 to radix10: ", err)
	}

	var url string
	row := s.db.QueryRow("SELECT url FROM shorturl WHERE id = $1", k)
	if err := row.Scan(&url); err != nil {
		s.Log.Println("Error querying database: ", err, k)
		return "", ErrKeyNotFound
	}

	if err != nil && err == sql.ErrNoRows {
		return "", ErrKeyNotFound
	}

	return url, nil
}

func (s PostgresStore) Set(originalURL string) (string, error) {
	// Check if the key already exists using orignalURL
	var k int64
	err := s.db.QueryRow("SELECT id FROM shorturl WHERE url = $1 LIMIT 1", originalURL).Scan(&k)
	if err != nil && err != sql.ErrNoRows {
		s.Log.Println("Error checking if key exists: ", err)
		return "", ErrKeyAlreadyExists
	}

	if k != 0 {
		return generator.ConvertRadix62(k), nil
	}

	newId, err := generator.GenerateSnowFlakeKey()
	if err != nil {
		s.Log.Println("Error generating snowflake key: ", err)
		return "", err
	}

	err = s.db.QueryRow("INSERT INTO shorturl (id, url) VALUES ($1,$2) RETURNING id", newId, originalURL).Scan(&k)
	if err != nil {
		s.Log.Println("Error inserting into database: ", err)
		return "", ErrKeyAlreadyExists
	}

	s.Log.Println("Inserted into database: ", k)
	return generator.ConvertRadix62(k), nil
}
