package mongodb

import (
	"context"
	"os"
	"time"

	"github.com/danyele/podp/internal/shared/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoClientImpl struct {
	client   *mongo.Client
	database string
}

func NovoMongoClient(ctx context.Context) (Client, error) {
	log := logger.New("MongoDB: NovoMongoClient")

	url := os.Getenv("MONGO_URL")
	database := os.Getenv("MONGO_DATABASE")

	ctxConn, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()

	client, err := mongo.Connect(ctxConn, options.Client().ApplyURI(url))
	if err != nil {
		log.Error("erro ao conectar no MongoDB", "erro", err)
		return nil, err
	}

	if err := client.Ping(ctxConn, nil); err != nil {
		log.Error("erro ao pingar MongoDB", "erro", err)
		return nil, err
	}

	log.Info("conectado ao MongoDB", "database", database)
	return &mongoClientImpl{client: client, database: database}, nil
}

func (m *mongoClientImpl) InsertOne(ctx context.Context, collection string, document interface{}) (interface{}, error) {
	coll := m.client.Database(m.database).Collection(collection)
	result, err := coll.InsertOne(ctx, document)
	if err != nil {
		return nil, err
	}
	return result.InsertedID, nil
}

func (m *mongoClientImpl) Find(ctx context.Context, collection string, filter interface{}, opts ...*options.FindOptions) ([]bson.M, error) {
	coll := m.client.Database(m.database).Collection(collection)

	cursor, err := coll.Find(ctx, filter, opts...)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}

func (m *mongoClientImpl) Disconnect(ctx context.Context) error {
	return m.client.Disconnect(ctx)
}
