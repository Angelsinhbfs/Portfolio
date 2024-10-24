package db

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func (mc *MongoClient) DeletePortfolioEntry(ctx context.Context, entry string) (interface{}, error) {
	collection := mc.client.Database("portfolio").Collection("entries")
	// Convert entry string to ObjectID
	objID, err := primitive.ObjectIDFromHex(entry)
	if err != nil {
		return nil, fmt.Errorf("failed to convert entry to ObjectID: %v", err)
	}

	filter := bson.M{"_id": objID}

	result, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to delete entry: %v", err)
	}

	return result.DeletedCount, nil
}

func (mc *MongoClient) EditPortfolioEntry(ctx context.Context, entry interface{}) (interface{}, error) {
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

func (mc *MongoClient) GetPortfolioEntryByID(ctx context.Context, id string) (interface{}, error) {
	if id == "" {
		return nil, nil
	}
	collection := mc.client.Database("portfolio").Collection("entries")
	// Convert entry string to ObjectID
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("failed to convert entry to ObjectID: %v", err)
	}

	filter := bson.M{"_id": objID}
	var tile bson.M
	err = collection.FindOne(ctx, filter).Decode(&tile)
	if err != nil {
		return nil, err
	} else {
		return tile, nil
	}
}

func (mc *MongoClient) UpdatePortfolioEntry(ctx context.Context, id string, tile interface{}) (interface{}, error) {
	collection := mc.client.Database("portfolio").Collection("entries")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("failed to convert entry to ObjectID: %v", err)
	}

	update := bson.M{"$set": tile} // Use $set operator to update all fields of the tile object

	result, err := collection.UpdateByID(ctx, objID, update)
	if err != nil {
		return nil, fmt.Errorf("failed to update entry: %v", err)
	}

	return result.ModifiedCount, nil
}
