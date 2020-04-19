package index

import (
	"context"

	"bytes"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// EnsureIndex if acs
func EnsureIndex(ctx context.Context, collection *mongo.Collection, keys map[string]int32) error {
	if collection == nil {
		return nil
	}
	if keys == nil {
		return nil
	}

	for key, expire := range keys {
		indexView := collection.Indexes()
		userIndex := mongo.IndexModel{
			Keys:    bson.D{{Key: key, Value: int32(1)}},
			Options: options.Index().SetName(key + "_index"),
		}

		if key == "create_time" && expire > 0 {
			userIndex.Options = options.Index().SetName("create_time_index").SetExpireAfterSeconds(expire)
		}
		err := createIndexIfNotExists(ctx, indexView, userIndex)
		if err != nil {
			return err
		}
	}

	/*
		expireIndex := mongo.IndexModel{
			Keys:    bson.D{{Key: "create_time", Value: int32(1)}},
			Options: options.Index().SetName("create_time_index").SetExpireAfterSeconds(60 * 60 * 24 * 30),
		}
		err = createIndexIfNotExists(ctx, indexView, expireIndex)
		if err != nil {
			log.Errorf("Create expire index error %v", err)
			return err
		}
	*/

	return nil
}

func createIndexIfNotExists(ctx context.Context, iv mongo.IndexView, model mongo.IndexModel) error {
	c, err := iv.List(ctx)
	if err != nil {
		return err
	}
	defer func() {
		_ = c.Close(ctx)
	}()

	var found bool
	for c.Next(ctx) {
		keyElem, err := c.Current.LookupErr("key")
		if err != nil {
			return err
		}

		keyElemDoc := keyElem.Document()
		modelKeysDoc, err := bson.Marshal(model.Keys)
		if err != nil {
			return err
		}

		if bytes.Equal(modelKeysDoc, keyElemDoc) {
			found = true
			break
		}
	}

	if !found {
		_, err = iv.CreateOne(ctx, model)
		if err != nil {
			return err
		}
	}

	return nil
}
