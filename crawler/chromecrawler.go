package crawler

import (
	"context"
	"crawler/utils"
	"time"

	"github.com/chromedp/chromedp"
)



func CrawlPage(ctx context.Context, rawurl string) (string, []string, error) {
    url := utils.NormalizeURL(rawurl)

    ctx, cancel := chromedp.NewContext(ctx)
    defer cancel()

    ctx, cancel = context.WithTimeout(ctx, 15*time.Second)
    defer cancel()

    var title string
    var links []string

    err := chromedp.Run(ctx,
        chromedp.Navigate(url),
        chromedp.Sleep(2*time.Second),
        chromedp.WaitReady("body", chromedp.ByQuery),
        chromedp.Title(&title),
        chromedp.Evaluate(`Array.from(document.querySelectorAll("a")).map(a=>a.href);`, &links),
    )
    return title, links, err
}
