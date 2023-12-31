package main

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"go-url-short/internal/server"
	"go-url-short/internal/store"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln("Error loading .env file")
	}

	host, port := os.Getenv("HOST"), os.Getenv("PORT")
	var dbConfig store.DatabaseConfig
	if err := envconfig.Process("DB", &dbConfig); err != nil {
		log.Fatalln("Error processing database config: ", err)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	s := server.NewHTTPServer(&server.HTTPServerArgs{
		Host:     host,
		Port:     port,
		DbConfig: &dbConfig,
	})
	go func() {
		fmt.Println("Server is listening on port: ", port)
		if err = s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalln("Could not listen on ", port, " ", err)
		}
	}()

	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		log.Fatalf("Could not gracefully shutdown the server: %v\n", err)
	}

	log.Println("Server Exited")
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
