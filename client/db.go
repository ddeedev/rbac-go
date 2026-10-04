package client

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	appconfig "github.com/ddeedev/rbac-go/config"
)

var collection *mongo.Collection

func ConnectToMongoDB(cfc *appconfig.DatabaseConfig) (*mongo.Client, error) {
	uri, err := appconfig.BuildMongoURI(cfc)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	log.Println("Conntect to mongo successfully....")

	return client, nil
}

func GetCollectionPointer() *mongo.Collection {
	return collection
}
