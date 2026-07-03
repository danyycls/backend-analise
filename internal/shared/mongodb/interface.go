package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Client interface {
	InsertOne(ctx context.Context, collection string, document interface{}) (interface{}, error)
	Find(ctx context.Context, collection string, filter interface{}, opts ...*options.FindOptions) ([]bson.M, error)
	Disconnect(ctx context.Context) error
}
