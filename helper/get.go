package helper

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Get  get one for mongo
func Get(ctx context.Context, collection *mongo.Collection, filter interface{}, result interface{}, ops ...*options.FindOneOptions) error {
	if result == nil {
		return nil
	}
	cur := collection.FindOne(ctx, filter, ops...)
	err := cur.Decode(filter)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil
		}
	}
	return nil
}
