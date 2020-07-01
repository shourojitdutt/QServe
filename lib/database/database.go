package database

import (
	"context"
	"lib/misc"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GetMongoClient : Utility function to get new Mongo Client Connections
func GetMongoClient() *mongo.Client {
	mongoClient, clientError := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if misc.Must(clientError) {
		return nil
	}
	return mongoClient
}

// GetCollection : Pass in DB Name and Collection Name to get Collection
func GetCollection(dbName string, collectionName string) *mongo.Collection {
	var clientConnection = GetMongoClient()
	return clientConnection.Database(dbName).Collection(collectionName)
}

// SearchInCollection : Takes collection present in MongoDB and checks if the key-value pair exists
func SearchInCollection(collection *mongo.Collection, key string, value interface{}) bson.Raw {
	ctx, cancelCtx := context.WithTimeout(context.Background(), 10*time.Second)
	filter := bson.M{key: value}
	bsonRaw, err := collection.FindOne(ctx, filter).DecodeBytes()
	cancelCtx()
	if misc.Must(err) {
		return nil
	}
	return bsonRaw
}

// CheckExistsInCollection : A simple true-false check for the existence of a value such as an ID in DB
func CheckExistsInCollection(collection *mongo.Collection, key string, value interface{}) bool {
	ctx, cancelCtx := context.WithTimeout(context.Background(), 10*time.Second)
	filter := bson.M{key: value}
	_, err := collection.FindOne(ctx, filter).DecodeBytes()
	cancelCtx()
	if misc.Must(err) {
		return false
	}
	return true
}

// UpdateOneInCollection : Updates an entry in collection with the new value by ID, then returns bool based on
// whether it was successfully updated or not.
func UpdateOneInCollection(collection *mongo.Collection, id primitive.ObjectID, key string, value interface{}) bool {
	if CheckExistsInCollection(collection, "_id", id) {
		// If it exists in collection, Updating it
		ctx, cancelCtx := context.WithTimeout(context.Background(), 10*time.Second)
		filter := bson.M{"_id": id}
		update := bson.M{"$set": bson.M{key: value}}
		_, err := collection.UpdateOne(ctx, filter, update)
		cancelCtx()
		if misc.Must(err) {
			// Some error occurred
			return false
		}
		// Successfully Updated
		return true
	}
	// Doesn't exist in collection
	return false
}

// InsertOneIntoCollection : Inserts a new entry into the collection
func InsertOneIntoCollection(collection *mongo.Collection, data bson.M) bool {
	// If it exists in collection, Updating it
	ctx, cancelCtx := context.WithTimeout(context.Background(), 10*time.Second)
	_, err := collection.InsertOne(ctx, data)
	cancelCtx()
	if misc.Must(err) {
		// Some error occurred
		return false
	}
	// Successfully Inserted
	return true
}

// DeleteFromCollection : Deletes an entry from the collection
func DeleteFromCollection(collection *mongo.Collection, key string, value interface{}) bool {
	if CheckExistsInCollection(collection, key, value) {
		// Key-value pair exists. Deleting entry.
		ctx, cancelCtx := context.WithTimeout(context.Background(), 10*time.Second)
		filter := bson.M{key: value}
		_, err := collection.DeleteOne(ctx, filter)
		cancelCtx()
		if misc.Must(err) {
			// Some error occurred
			return false
		}
		// Successfully Deleted
		return true
	}
	// If it doesn't exist, returning false
	return false
}

// ClearCollection : Deletes every entry from a collection. Use carefully.
func ClearCollection(collection *mongo.Collection) bool {
	ctx, cancelCtx := context.WithTimeout(context.Background(), 10*time.Second)
	filter := bson.M{} // Empty filter since we're deleting everything.
	_, err := collection.DeleteMany(ctx, filter)
	cancelCtx()
	if misc.Must(err) {
		// Some error occurred
		return false
	}
	// Successfully cleared the entire collection
	return true
}
