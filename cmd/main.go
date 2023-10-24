package main

import (
	"fmt"
	"log"
	"os"
	"user_service/internal/api"
	"user_service/internal/logging"
	"user_service/internal/repository"

	"github.com/joho/godotenv"
)

func main() {
	logger := logging.GetLogger()
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	user := os.Getenv("USER")
	dbname := os.Getenv("DBNAME")
	password := os.Getenv("PASSWORD")

	connStr := fmt.Sprintf("host= %s port=%s user=%s dbname=%s password=%s sslmode=disable", host, port, user, dbname, password)

	store, err := repository.NewPostgresStore(connStr, logger)
	if err != nil {
		log.Fatal(err)
	}

	server := api.NewAPIServer("localhost:8080", store, logger)
	logger.Info("Server is running on htttp://localhost:8080")
	server.Run()

}
