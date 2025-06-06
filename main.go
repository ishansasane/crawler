package main

import (
	"context"
	"crawler/storage"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)



const (
	maxDepth    = 3
	maxPages    = 100
	workerCount = 5
	timeout     = 10 * time.Second
)

type CrawlJob struct {
	url   string
	depth int
}

func fetchLinks(client *http.Client, pageURL string) ([]string, error) {
	resp, err := client.Get(pageURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("non-200 status code: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	var links []string
	base, err := url.Parse(pageURL)
	if err != nil {
		return nil, err
	}

	doc.Find("a[href]").Each(func(_ int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if !exists {
			return
		}

		u, err := url.Parse(href)
		if err != nil {
			return
		}

		absURL := base.ResolveReference(u)
		links = append(links, absURL.String())
	})

	return links, nil
}

func worker(id int, client *http.Client, mongoClient *mongo.Client, jobs <-chan CrawlJob, results chan<- []string, visited *sync.Map, wg *sync.WaitGroup, pageCount *int, pageCountMu *sync.Mutex) {
	defer wg.Done()

	for job := range jobs {
		if job.depth > maxDepth {
			continue
		}

		pageCountMu.Lock()
		if *pageCount >= maxPages {
			pageCountMu.Unlock()
			return
		}
		pageCountMu.Unlock()

		if _, loaded := visited.LoadOrStore(job.url, true); loaded {
			continue
		}

		fmt.Printf("[Worker %d] Crawling: %s (depth %d)\n", id, job.url, job.depth)

		links, err := fetchLinks(client, job.url)
		if err != nil {
			fmt.Printf("[Worker %d] Error fetching %s: %s\n", id, job.url, err)
			continue
		}

		// Save page URL to MongoDB
		err = storage.SavePageToMongo(mongoClient, job.url)
		if err != nil {
			fmt.Printf("[Worker %d] Error saving to Mongo: %v\n", id, err)
		}

		pageCountMu.Lock()
		*pageCount++
		pageCountMu.Unlock()

		results <- links
	}
}

func main() {
	startURL := "https://old.reddit.com/"

	client := &http.Client{Timeout: timeout}

	mongoClient, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}
	defer mongoClient.Disconnect(context.Background())

	jobs := make(chan CrawlJob, 100)
	results := make(chan []string, 100)

	var visited sync.Map
	var wg sync.WaitGroup

	pageCount := 0
	var pageCountMu sync.Mutex

	// Start workers
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go worker(i+1, client, mongoClient, jobs, results, &visited, &wg, &pageCount, &pageCountMu)
	}

	jobs <- CrawlJob{url: startURL, depth: 0}

	go func() {
		for links := range results {
			pageCountMu.Lock()
			if pageCount >= maxPages {
				pageCountMu.Unlock()
				close(jobs)
				return
			}
			pageCountMu.Unlock()

			for _, link := range links {
				jobs <- CrawlJob{url: link, depth: 1}
			}
		}
	}()

	wg.Wait()
	close(results)

	fmt.Println("Crawling completed.")
}
