package main

import (
	"Portfolio/db"
	"context"
	"github.com/joho/godotenv"
	"log"
	"os"
)

var DBClient *db.MongoClient
var CTX context.Context
var ExpectedUser string
var ExpectedKey string

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	ExpectedUser = os.Getenv("EXPECTED_USER")
	ExpectedKey = os.Getenv("EXPECTED_KEY")
	if ExpectedUser == "" || ExpectedKey == "" {
		log.Fatal("Credentials not set in env")
	}
	DBClient, err = db.NewMongoClient("mongodb://localhost:27017", "Portfolio")
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}
	Route()
}
