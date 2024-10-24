package main

import (
	"Portfolio/db"
	"context"
	"log"
	"os"
)

var DBClient *db.MongoClient
var CTX context.Context
var ExpectedUser string
var ExpectedKey string

func main() {
	ExpectedUser = os.Getenv("EXPECTED_USER")
	ExpectedKey = os.Getenv("EXPECTED_KEY")
	var err error
	DBClient, err = db.NewMongoClient("mongodb://localhost:27017", "Portfolio")
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}
	Route()
}
