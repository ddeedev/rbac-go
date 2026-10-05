package client

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	appconfig "github.com/ddeedev/rbac-go/config"
)

func ConnectToMongoDB(cfc *appconfig.DatabaseConfig) (*mongo.Client, error) {
	uri, err := appconfig.BuildMongoURI(cfc)
	if err != nil {
		return nil, err
	}

	// bound the initial connect + ping so startup fails fast if Mongo is down.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("mongo connect: %w", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("mongo ping: %w", err)
	}

	log.Println("connected to mongo successfully")
	return client, nil
}
