package main

import (
	"context"
	"crawler/crawler"
	"crawler/storage"
	"crawler/utils"
	"fmt"
	"log"
	"time"
)



func main() {
    err := storage.ConnectMongo("mongodb://localhost:27017")
    if err != nil {
        log.Fatal("MongoDB connect error:", err)
    }
    defer storage.DisconnectMongo()

    seed := []string{
        "https://news.ycombinator.com/",
        "https://www.bbc.com/",
        "https://www.geeksforgeeks.org/",
    }

    queue := make(chan string, 1000)
    visited := make(map[string]bool)

    for _, u := range seed {
        queue <- utils.NormalizeURL(u)
    }

    ctx := context.Background()

    for rawurl := range queue {
        url := utils.NormalizeURL(rawurl)

        if visited[url] || storage.AlreadyCrawled(url) {
            continue
        }

        fmt.Println("Crawling:", url)
        title, links, err := crawler.CrawlPage(ctx, url)
        if err != nil {
            fmt.Println("Error crawling:", url, err)
            visited[url] = true
            continue
        }

        err = storage.SavePage(storage.PageData{URL: url, Title: title, Links: links, FetchedAt: time.Now()})
        if err != nil {
            fmt.Println("DB save error:", err)
        }

        fmt.Printf("Saved: %-50s | found %d links\n", url, len(links))

        for _, l := range links {
            n := utils.NormalizeURL(l)
            if _, seen := visited[n]; !seen {
                queue <- n
            }
        }

        visited[url] = true
        time.Sleep(2 * time.Second)
    }
}
