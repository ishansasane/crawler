package crawler

import (
	"io"
	"net/url"

	"github.com/PuerkitoBio/goquery"
)

func ExtractTitle(body io.Reader) string {
    doc, err := goquery.NewDocumentFromReader(body)
    if err != nil {
        return ""
    }
    return doc.Find("title").Text()
}

func ExtractLinks(body io.Reader, base string) []string {
    var links []string

    doc, err := goquery.NewDocumentFromReader(body)
    if err != nil {
        return links
    }

    baseURL, _ := url.Parse(base)

    doc.Find("a[href]").Each(func(_ int, s *goquery.Selection) {
        href, _ := s.Attr("href")
        absURL, err := baseURL.Parse(href)
        if err == nil {
            links = append(links, absURL.String())
        }
    })

    return links
}
