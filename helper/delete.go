package helper

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Delete  get one for mongo
func Delete(ctx context.Context, collection *mongo.Collection, filter interface{}, ops ...*options.DeleteOptions) (int64, error) {
	if collection == nil {
		return 0, errors.New("connection not available")
	}
	cur, err := collection.DeleteMany(ctx, filter, ops...)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return 0, nil
		}
		return 0, err
	}
	return cur.DeletedCount, nil
}
