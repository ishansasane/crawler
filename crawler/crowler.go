package crawler

import (
	"bytes"
	"crawler/storage"
	"crawler/utils"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"sync"
	"time"
)
var visited = make(map[string]bool)
var mu sync.Mutex

func Start(url string, depth int) []string {
    var results []string
    var wg sync.WaitGroup

    var crawl func(string, int)
    crawl = func(link string, level int) {
        defer wg.Done()
        if level <= 0 {
            return
        }

        mu.Lock()
        if visited[link] {
            mu.Unlock()
            return
        }
        visited[link] = true
        mu.Unlock()

        fmt.Println("Crawling:", link)

        client := &http.Client{Timeout: 10 * time.Second}
        req, _ := http.NewRequest("GET", link, nil)
        req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; MyGoCrawler/1.0)")

        resp, err := client.Do(req)
        if err != nil || resp.StatusCode != 200 {
            return
        }
        defer resp.Body.Close()

        bodyBytes, err := io.ReadAll(resp.Body)
        if err != nil {
            return
        }
        bodyReader := bytes.NewReader(bodyBytes)

        title := ExtractTitle(bodyReader)
        bodyReader.Seek(0, io.SeekStart)
        links := ExtractLinks(bodyReader, link)

        page := storage.PageData{
            URL:       link,
            Title:     title,
            Status:    resp.StatusCode,
            FetchedAt: time.Now(),
        }
        storage.SavePageToMongo("webcrawler", page)

        mu.Lock()
        results = append(results, link)
        mu.Unlock()

        time.Sleep(time.Duration(1+rand.Intn(2)) * time.Second)

        for _, l := range links {
            if utils.IsSameDomain(link, l) {
                wg.Add(1)
                go crawl(l, level-1)
            }
        }
    }

    wg.Add(1)
    go crawl(url, depth)
    wg.Wait()
    return results
}