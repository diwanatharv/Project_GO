package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Dbinstance() *mongo.Client {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error in loading env file")
	}
	Mongodb := os.Getenv("MONGO")
	client, err := mongo.NewClient(options.Client().ApplyURI(Mongodb))
	if err != nil {
		log.Fatal("Error in loading mongo")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal("Error in loading mongo ctx")
	}
	fmt.Println("connecting to mongo")
	return client
}

var Client *mongo.Client = Dbinstance()

func OpenCollection(client *mongo.Client, collectionname string) *mongo.Collection {
	var collection *mongo.Collection = client.Database("cluster").Collection(collectionname)
	return collection
}
