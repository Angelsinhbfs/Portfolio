package db

import "context"

// MongoDBInterface defines the methods for interacting with MongoDB
type MongoDBInterface interface {
	//todo:these should just get the url to get that specific image. maybe a low res thumbnail until main loads
	GetImage(ctx context.Context, imageID string)

	AddPortfolioEntry(ctx context.Context, entry string) (interface{}, error)
	GetPortfolioEntries(ctx context.Context) (interface{}, error)
}
