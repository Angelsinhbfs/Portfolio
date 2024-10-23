package db

import "context"

// MongoDBInterface defines the methods for interacting with MongoDB
type MongoDBInterface interface {
	// single document functions
	InsertOne(ctx context.Context, collection string, document interface{}) (interface{}, error)
	FindOne(ctx context.Context, collection string, filter interface{}) (interface{}, error)
	UpdateOne(ctx context.Context, collection string, filter interface{}, update interface{}) (interface{}, error)
	DeleteOne(ctx context.Context, collection string, filter interface{}) (interface{}, error)

	// useful functions
	FindUser(ctx context.Context, username string) (interface{}, error)
	AddUser(ctx context.Context, user interface{}) (interface{}, error)
	UpdateUser(ctx context.Context, username string, update interface{}) (interface{}, error)
	DeleteUser(ctx context.Context, username string) (interface{}, error)
	GetUserInbox(ctx context.Context, username string) ([]map[string]string, error)
	GetUserOutbox(ctx context.Context, username string) ([]map[string]string, error)
	AddToUserOutbox(ctx context.Context, username string, message map[string]string) (interface{}, error)
	GetSharedInbox(ctx context.Context) ([]map[string]string, error)

	//todo:these should just get the url to get that specific image. maybe a low res thumbnail until main loads
	GetImage(ctx context.Context, imageID string)
	GetImagesByUser(ctx context.Context, username string)

	AddPortfolioEntry(ctx context.Context, entry string) (interface{}, error)
	GetPortfolioEntries(ctx context.Context) (interface{}, error)
}
