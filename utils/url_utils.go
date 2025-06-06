package utils

import (
	"net/url"
)

func IsSameDomain(base, target string) bool {
    baseURL, err1 := url.Parse(base)
    targetURL, err2 := url.Parse(target)
    if err1 != nil || err2 != nil {
        return false
    }
    return baseURL.Hostname() == targetURL.Hostname()
}
