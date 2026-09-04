package app

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.uber.org/zap"
)

func newMongoDb() *mongo.Database {
	uri := MongoConn
	if uri == "" {
		uri = "mongodb://127.0.0.1:27017"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		zap.S().Panicf("failed to create MongoDB client: %+v", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		zap.S().Warnf("MongoDB ping warning (server will retry on demand): %+v", err)
	} else {
		zap.S().Info("Successfully connected to MongoDB")
	}

	return client.Database("wfrp5e")
}
