package utils

import (
	"net/url"
)

func IsSameDomain(base string, other string) bool {
    u1, err1 := url.Parse(base)
    u2, err2 := url.Parse(other)
    if err1 != nil || err2 != nil {
        return false
    }
    return u1.Hostname() == u2.Hostname()
}
