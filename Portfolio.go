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
	mongoURI := os.Getenv("MONGODB_URI")
	dbName := os.Getenv("MONGODB_DB_NAME")
	if mongoURI == "" || dbName == "" {
		log.Fatal("MongoDB connection information not set in env")
	}
	DBClient, err = db.NewMongoClient(mongoURI, dbName)
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}
	Route()
}
