package db

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoClient struct {
	client *mongo.Client
	dbName string
}

// NewMongoClient creates a new MongoClient
func NewMongoClient(uri, dbName string) (*MongoClient, error) {
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, err
	}

	return &MongoClient{
		client: client,
		dbName: dbName,
	}, nil
}

func (mc *MongoClient) GetImage(ctx context.Context, imageID string) (string, error) {
	// TODO: Implement the method
	return "", errors.New("GetImage not implemented")
}

func (mc *MongoClient) AddPortfolioEntry(ctx context.Context, entry interface{}) (interface{}, error) {
	collection := mc.client.Database("portfolio").Collection("entries")

	entryBSON, err := bson.Marshal(entry)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal entry: %v", err)
	}
	result, err := collection.InsertOne(ctx, entryBSON)
	if err != nil {
		return nil, fmt.Errorf("failed to add entry: %v", err)
	}
	return result.InsertedID, nil
}

func (mc *MongoClient) GetPortfolioEntries(ctx context.Context) (interface{}, error) {
	collection := mc.client.Database("portfolio").Collection("entries")
	filter := bson.M{}
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find user entries: %v", err)
	}
	defer cursor.Close(ctx)

	var entries []interface{}
	for cursor.Next(ctx) {
		var message bson.M
		if err := cursor.Decode(&message); err != nil {
			return nil, fmt.Errorf("failed to decode message: %v", err)
		}
		entries = append(entries, message)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %v", err)
	}

	return entries, nil
}
