package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/gorillamux"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"go-url-short/internal/server"
	"go-url-short/internal/store"
	"log"
	"net/http"
	"os"
)

func init() {

}

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Errorf("Error loading .env file")
	}
	host, port := os.Getenv("HOST"), os.Getenv("PORT")
	var dbConfig store.DatabaseConfig
	if err := envconfig.Process("DB", &dbConfig); err != nil {
		log.Fatalln("Error processing database config: ", err)
	}
	s := server.NewHTTPServer(&server.HTTPServerArgs{
		Host:     host,
		Port:     port,
		DbConfig: &dbConfig,
	})

	log.Println("Starting lambda server")
	if runtime, _ := os.LookupEnv("AWS_EXECUTION_ENV"); runtime != "" {
		log.Println("Running in AWS Lambda")
		r := s.Handler.(*mux.Router)
		adapter := gorillamux.New(r)
		lambda.Start(adapter.ProxyWithContext)
	} else {
		fmt.Printf("Running in Local :%s\n", port)
		if err = s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalln("Could not listen on ", port, " ", err)
		}
	}
}
