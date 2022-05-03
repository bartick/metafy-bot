package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	collection *mongo.Collection
	ctx        context.Context
)

func connectDb() *mongo.Client {
	client, err := mongo.NewClient(options.Client().ApplyURI(os.Getenv("MONGODB_URI")))
	if err != nil {
		fmt.Println("Error creating mongo client")
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		fmt.Println("Error connecting to database")
		return nil
	} else {
		fmt.Println("Connected to database")
	}

	collection = client.Database("Metafy").Collection("modmail")

	var result bson.M

	cur, err := collection.Find(ctx, bson.D{})

	if err != nil {
		fmt.Println(err.Error())
	}

	for cur.Next(context.TODO()) {
		err := cur.Decode(&result)
		if err != nil {
			fmt.Println(err.Error())
			return nil
		}
		TicketChannelMapping[result["authorID"].(string)] = result["channelID"].(string)
		TicketChannelMapping[result["channelID"].(string)] = result["authorID"].(string)
	}

	return client
}
