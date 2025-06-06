package main

import (
	"crawler/crawler"
	"crawler/storage"
	"flag"
	"fmt"
	"math/rand"
	"time"
)

func main() {
    rand.Seed(time.Now().UnixNano())

    startURL := flag.String("url", "https://books.toscrape.com", "Start URL")
    depth := flag.Int("depth", 2, "Crawling depth")
    flag.Parse()

    err := storage.ConnectMongo("mongodb://localhost:27017")
    if err != nil {
        fmt.Println("MongoDB connection error:", err)
        return
    }

    results := crawler.Start(*startURL, *depth)
    fmt.Println("Crawled URLs:", results)
    fmt.Println("Crawling complete. Results saved.")
}

