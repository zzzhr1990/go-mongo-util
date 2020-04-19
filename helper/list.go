package helper

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GetList get list for mongo
func GetList(ctx context.Context, collection *mongo.Collection, filter interface{}, lists interface{}, ops ...*options.FindOptions) error {
	if lists == nil {
		return nil
	}
	cur, err := collection.Find(ctx, filter, ops...)
	if err != nil {
		return err
	}
	err = cur.All(ctx, lists)
	if err != nil {
		return err
	}
	return nil
}
