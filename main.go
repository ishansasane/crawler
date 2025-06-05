package main

import (
	"crawler/config"
	"crawler/crawler"
	"crawler/storage"
	"fmt"
)

func main() {
    cfg := config.ParseFlags()

    results := crawler.Start(cfg.StartURL, cfg.MaxDepth)
    storage.SaveToFile("results.txt", results)

    fmt.Println("Crawling complete. Results saved to results.txt")
}
