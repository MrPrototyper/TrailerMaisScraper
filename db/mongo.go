package db

import (
	"context"
	"log"
	"os"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	clientInstance *mongo.Client
	clientOnce     sync.Once
	clientError    error
	ctx            context.Context
	cancel         context.CancelFunc
)

// ConnectMongoDB initializes the MongoDB client as a singleton
func Connect() (*mongo.Client, context.Context, context.CancelFunc, error) {
	uri := os.Getenv("MONGODB_URI")
	username := os.Getenv("MONGODB_USERNAME")
	password := os.Getenv("MONGODB_PASSWORD")

	clientOnce.Do(func() {
		clientOptions := options.Client().ApplyURI(uri)

		// Add authentication options only if username and password are provided
		if username != "" && password != "" {
			clientOptions.SetAuth(options.Credential{
				Username: username,
				Password: password,
			})
		}

		clientInstance, clientError = mongo.NewClient(clientOptions)
		if clientError != nil {
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		clientError = clientInstance.Connect(ctx)
		if clientError != nil {
			cancel()
			return
		}

		// Ensure we disconnect in case of an error during Ping
		defer func() {
			if clientError != nil {
				cancel()
				clientInstance.Disconnect(ctx)
			}
		}()

		clientError = clientInstance.Ping(ctx, nil)
		if clientError != nil {
			return
		}

		log.Println("Connected to MongoDB!")
	})

	if clientError != nil {
		return nil, nil, nil, clientError
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	return clientInstance, ctx, cancel, nil
}

// GetClient provides the MongoDB client instance
func GetClient() (*mongo.Client, context.Context) {
	return clientInstance, ctx
}
