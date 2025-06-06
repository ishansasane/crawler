package utils

import (
	"net/url"
	"strings"
)

func NormalizeURL(rawurl string) string {
    u, err := url.Parse(rawurl)
    if err != nil {
        return rawurl
    }

    u.Fragment = ""
    u.RawQuery = ""

    if len(u.Path) > 1 && u.Path[len(u.Path)-1] == '/' {
        u.Path = strings.TrimSuffix(u.Path, "/")
    }
    return u.String()
}
