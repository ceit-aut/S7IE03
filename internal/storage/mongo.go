package storage

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// NewConnection
// opens a new connection to mongodb database.
func NewConnection(cfg Config) (*mongo.Database, error) {
	// mongodb server options
	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	mongoURI := fmt.Sprintf("mongodb://%s:%d", cfg.Host, cfg.Port)

	clientOptions := options.Client().
		ApplyURI(mongoURI).
		SetServerAPIOptions(serverAPIOptions)

	// creating mongodb client
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, fmt.Errorf("mongoDB connection failed: %w", err)
	}

	// ping mongodb
	if er := client.Ping(context.TODO(), nil); er != nil {
		return nil, fmt.Errorf("mongoDB ping failed: %w", er)
	}

	return client.Database(cfg.Database), nil
}
