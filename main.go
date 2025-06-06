package main

import (
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const (
	maxDepth    = 3   // max levels deep to crawl
	maxPages    = 100 // max total pages to crawl
	workerCount = 5   // concurrent workers
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

func worker(id int, client *http.Client, jobs <-chan CrawlJob, results chan<- []string, visited *sync.Map, wg *sync.WaitGroup, pageCount *int, pageCountMu *sync.Mutex) {
	defer wg.Done()

	for job := range jobs {
		if job.depth > maxDepth {
			continue
		}

		// Check maxPages limit
		pageCountMu.Lock()
		if *pageCount >= maxPages {
			pageCountMu.Unlock()
			return
		}
		pageCountMu.Unlock()

		// Skip if already visited
		if _, loaded := visited.LoadOrStore(job.url, true); loaded {
			continue
		}

		fmt.Printf("[Worker %d] Crawling: %s (depth %d)\n", id, job.url, job.depth)

		links, err := fetchLinks(client, job.url)
		if err != nil {
			fmt.Printf("[Worker %d] Error fetching %s: %s\n", id, job.url, err)
			continue
		}

		// Increment page count
		pageCountMu.Lock()
		*pageCount++
		pageCountMu.Unlock()

		results <- links
	}
}

func main() {
	startURL := "https://news.ycombinator.com/"

	client := &http.Client{
		Timeout: timeout,
	}

	jobs := make(chan CrawlJob, 100)
	results := make(chan []string, 100)

	var visited sync.Map
	var wg sync.WaitGroup

	pageCount := 0
	var pageCountMu sync.Mutex

	// Start workers
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go worker(i+1, client, jobs, results, &visited, &wg, &pageCount, &pageCountMu)
	}

	// Start with initial URL
	jobs <- CrawlJob{url: startURL, depth: 0}

	go func() {
		for links := range results {
			// For each discovered link, send new crawl job with increased depth
			pageCountMu.Lock()
			if pageCount >= maxPages {
				pageCountMu.Unlock()
				close(jobs)
				return
			}
			pageCountMu.Unlock()

			for _, link := range links {
				jobs <- CrawlJob{url: link, depth: 1} // increase depth, or adapt to actual depth tracking
			}
		}
	}()

	wg.Wait()
	close(results)

	fmt.Println("Crawling completed.")
}
