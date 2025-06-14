package database

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DbCloseFunc func(context.Context) error

func clientOptions(connectionString string) *options.ClientOptions {
	clientOptions := options.Client().ApplyURI(connectionString)

	return clientOptions
}

func createClient(ctx context.Context, clientOptions *options.ClientOptions) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func ConnectDatabase(ctx context.Context, host string, port int, user string, password string) (*mongo.Database, DbCloseFunc, error) {
	connectionString := fmt.Sprintf("mongodb://%s:%s@%s:%d", user, password, host, port)
	clientOptions := clientOptions(connectionString)

	client, err := createClient(ctx, clientOptions)
	if err != nil {
		return nil, func(ctx context.Context) error { return nil }, err
	}

	return client.Database("yapper"), client.Disconnect, nil
}
