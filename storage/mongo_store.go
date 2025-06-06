package storage

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)
type PageData struct {
    URL       string    `bson:"url"`
    Title     string    `bson:"title"`
    Status    int       `bson:"status"`
    FetchedAt time.Time `bson:"fetched_at"`
}

var mongoClient *mongo.Client

func ConnectMongo(uri string) error {
    var err error
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    mongoClient, err = mongo.Connect(ctx, options.Client().ApplyURI(uri))
    return err
}

func SavePageToMongo(dbName string, data PageData) error {
    collection := mongoClient.Database(dbName).Collection("pages")
    filter := bson.M{"url": data.URL}
    update := bson.M{"$set": data}
    opts := options.Update().SetUpsert(true)
    _, err := collection.UpdateOne(context.Background(), filter, update, opts)
    return err
}