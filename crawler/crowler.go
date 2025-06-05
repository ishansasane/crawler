package crawler

import (
	"crawler/utils"
	"fmt"
	"net/http"
	"sync"
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
        resp, err := http.Get(link)
        if err != nil {
            return
        }
        defer resp.Body.Close()

        links := ExtractLinks(resp.Body, link)

        mu.Lock()
        results = append(results, link)
        mu.Unlock()

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