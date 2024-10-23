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

func (mc *MongoClient) InsertOne(ctx context.Context, collection string, document interface{}) (interface{}, error) {
	coll := mc.client.Database(mc.dbName).Collection(collection)
	result, err := coll.InsertOne(ctx, document)
	if err != nil {
		return nil, err
	}
	return result.InsertedID, nil
}

func (mc *MongoClient) FindOne(ctx context.Context, collection string, filter interface{}) (interface{}, error) {
	coll := mc.client.Database(mc.dbName).Collection(collection)
	var result bson.M
	err := coll.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (mc *MongoClient) UpdateOne(ctx context.Context, collection string, filter interface{}, update interface{}) (interface{}, error) {
	coll := mc.client.Database(mc.dbName).Collection(collection)
	result, err := coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}
	return result.ModifiedCount, nil
}

func (mc *MongoClient) DeleteOne(ctx context.Context, collection string, filter interface{}) (interface{}, error) {
	coll := mc.client.Database(mc.dbName).Collection(collection)
	result, err := coll.DeleteOne(ctx, filter)
	if err != nil {
		return nil, err
	}
	return result.DeletedCount, nil
}

// This DOES NOT check if the user is unique
func (mc *MongoClient) AddUser(ctx context.Context, user map[string]interface{}) (interface{}, error) {
	collection := mc.client.Database(mc.dbName).Collection("users")
	result, err := collection.InsertOne(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to add user: %v", err)
	}
	return result.InsertedID, nil
}

func (mc *MongoClient) UpdateUser(ctx context.Context, username string, update map[string]interface{}) (interface{}, error) {
	collection := mc.client.Database(mc.dbName).Collection("users")
	filter := bson.M{"username": username}
	updateDoc := bson.M{"$set": update}

	result, err := collection.UpdateOne(ctx, filter, updateDoc)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %v", err)
	}
	return result.ModifiedCount, nil
}

func (mc *MongoClient) DeleteUser(ctx context.Context, username string) (interface{}, error) {
	collection := mc.client.Database(mc.dbName).Collection("users")
	filter := bson.M{"username": username}

	result, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to delete user: %v", err)
	}
	return result.DeletedCount, nil
}
func (mc *MongoClient) FindUser(ctx context.Context, username string) (interface{}, error) {
	collection := mc.client.Database(mc.dbName).Collection("users")
	filter := bson.M{"username": username}

	var user bson.M
	err := collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to find user: %v", err)
	}

	return user, nil
}

func (mc *MongoClient) GetUserInbox(ctx context.Context, username string) ([]map[string]string, error) {
	collection := mc.client.Database(mc.dbName).Collection("outboxes")
	filter := bson.M{"username": username}

	var result struct {
		Activities []map[string]string `bson:"activities"`
	}

	err := collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("user outbox not found")
		}
		return nil, fmt.Errorf("failed to find user outbox: %v", err)
	}

	return result.Activities, nil
}
func (mc *MongoClient) AddToUserInbox(ctx context.Context, username string, activity map[string]string) (interface{}, error) {
	collection := mc.client.Database(mc.dbName).Collection("inboxes")
	filter := bson.M{"username": username}

	update := bson.M{
		"$push": bson.M{"activities": bson.M{"$each": []map[string]string{activity}, "$position": 0}},
	}

	result, err := collection.UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))
	if err != nil {
		return nil, fmt.Errorf("failed to add to user outbox: %v", err)
	}

	return result.UpsertedID, nil
}

func (mc *MongoClient) GetUserOutbox(ctx context.Context, username string) ([]map[string]string, error) {
	collection := mc.client.Database(mc.dbName).Collection("outboxes")
	filter := bson.M{"username": username}

	var result struct {
		Activities []map[string]string `bson:"activities"`
	}

	err := collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("user outbox not found")
		}
		return nil, fmt.Errorf("failed to find user outbox: %v", err)
	}

	return result.Activities, nil
}
func (mc *MongoClient) AddToUserOutbox(ctx context.Context, username string, activity map[string]string) (interface{}, error) {
	collection := mc.client.Database(mc.dbName).Collection("outboxes")
	filter := bson.M{"username": username}

	update := bson.M{
		"$push": bson.M{"activities": bson.M{"$each": []map[string]string{activity}, "$position": 0}},
	}

	result, err := collection.UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))
	if err != nil {
		return nil, fmt.Errorf("failed to add to user outbox: %v", err)
	}

	return result.UpsertedID, nil
}

func (mc *MongoClient) GetSharedInbox(ctx context.Context) ([]interface{}, error) {
	collection := mc.client.Database(mc.dbName).Collection("outboxes")
	filter := bson.M{"username": "shared"}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find user inbox: %v", err)
	}
	defer cursor.Close(ctx)

	var inbox []interface{}
	for cursor.Next(ctx) {
		var message bson.M
		if err := cursor.Decode(&message); err != nil {
			return nil, fmt.Errorf("failed to decode message: %v", err)
		}
		inbox = append(inbox, message)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %v", err)
	}

	return inbox, nil
}

func (mc *MongoClient) GetImage(ctx context.Context, imageID string) (string, error) {
	// TODO: Implement the method
	return "", errors.New("GetImage not implemented")
}

func (mc *MongoClient) GetImagesByUser(ctx context.Context, username string) ([]string, error) {
	// TODO: Implement the method
	return nil, errors.New("GetImagesByUser not implemented")
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
