package main

import (
	"Portfolio/db"
	"context"
	"log"
)

var DBClient *db.MongoClient
var CTX context.Context

func main() {
	var err error
	DBClient, err = db.NewMongoClient("mongodb://localhost:27017", "Portfolio")
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}
	Route()
}
