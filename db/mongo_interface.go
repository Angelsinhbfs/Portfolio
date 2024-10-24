package db

import "context"

// MongoDBInterface defines the methods for interacting with MongoDB
type MongoDBInterface interface {
	AddPortfolioEntry(ctx context.Context, entry string) (interface{}, error)
	GetPortfolioEntries(ctx context.Context) (interface{}, error)
}
