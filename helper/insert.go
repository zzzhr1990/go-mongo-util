package helper

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Insert  get one for mongo
func Insert(ctx context.Context, collection *mongo.Collection, filter interface{}, data interface{}) error {
	if data == nil {
		return nil
	}
	bs := bson.D{bson.E{Key: "$set", Value: data}}
	_, err := collection.UpdateOne(ctx, filter, bs, options.Update().SetUpsert(true))

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil
		}
		return err
	}
	return nil
}

// Create  get one for mongo
func Create(ctx context.Context, collection *mongo.Collection, filter interface{}, data interface{}) error {
	return Insert(ctx, collection, filter, data)
}
