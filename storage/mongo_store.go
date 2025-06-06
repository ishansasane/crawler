package storage

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type Page struct {
	URL       string    `bson:"url"`
	Timestamp time.Time `bson:"timestamp"`
}

// SavePageToMongo inserts the page URL and timestamp into the MongoDB collection.
func SavePageToMongo(client *mongo.Client, pageURL string) error {
	collection := client.Database("crawler").Collection("pages")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, Page{
		URL:       pageURL,
		Timestamp: time.Now(),
	})
	return err
}
