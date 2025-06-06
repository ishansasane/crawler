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
    Links     []string  `bson:"links"`
    FetchedAt time.Time `bson:"fetched_at"`
}

var client *mongo.Client

func ConnectMongo(uri string) error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    var err error
    client, err = mongo.Connect(ctx, options.Client().ApplyURI(uri))
    return err
}

func DisconnectMongo() {
    client.Disconnect(context.Background())
}

func SavePage(data PageData) error {
    coll := client.Database("webcrawler").Collection("pages")
    filter := bson.M{"url": data.URL}
    update := bson.M{"$setOnInsert": data}
    opts := options.Update().SetUpsert(true)
    _, err := coll.UpdateOne(context.Background(), filter, update, opts)
    return err
}

func AlreadyCrawled(url string) bool {
    coll := client.Database("webcrawler").Collection("pages")
    count, _ := coll.CountDocuments(context.Background(), bson.M{"url": url})
    return count > 0
}
