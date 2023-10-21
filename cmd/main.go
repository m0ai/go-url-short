package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"go-url-short/internal/server"
	"go-url-short/internal/store"
	"log"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln("Error loading .env file")
	}
	host, port := os.Getenv("HOST"), os.Getenv("PORT")
	s := server.NewHTTPServer(host, port, *createDBConfig())

	fmt.Println("Server is listening on port: ", port)
	s.ListenAndServe()
}

func createDBConfig() *store.DatabaseConfig {
	var dbConfig *store.DatabaseConfig

	if host := os.Getenv("DB_HOST"); host != "" {
		dbConfig = &store.DatabaseConfig{
			Name:     os.Getenv("DB_NAME"),
			Host:     host,
			Port:     5432,
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
		}
	} else {
		dbConfig = &store.DatabaseConfig{}
	}

	return dbConfig
}
