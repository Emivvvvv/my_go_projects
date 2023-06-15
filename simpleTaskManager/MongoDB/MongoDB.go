package MongoDB

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

var AuthCollection *mongo.Collection = nil

func ConnectToMongo() (*options.ClientOptions, *mongo.Client, *mongo.Collection) {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	collection := client.Database("test").Collection("tasks")
	if collection == nil {
		log.Fatal("test collection is nil")
	}
	AuthCollection = client.Database("test").Collection("user_data")
	if AuthCollection == nil {
		log.Fatal("user_data collection is nil")
	}

	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!\n" +
		"Database: \"test\", Collection: \"tasks\" and \"user_data\"")

	return clientOptions, client, collection
}
