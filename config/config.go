package config

import (
	"flag"
)

type Config struct {
    StartURL string
    MaxDepth int
}

func ParseFlags() Config {
    var url string
    var depth int

    flag.StringVar(&url, "url", "https://old.reddit.com/", "Starting URL")
    flag.IntVar(&depth, "depth", 5, "Maximum crawl depth")
    flag.Parse()

    return Config{StartURL: url, MaxDepth: depth}
}
