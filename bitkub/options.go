package bitkub

import (
	"net/http"
	"net/url"
)

type Options struct {
	Client    *http.Client
	BaseURL   *url.URL
	UserAgent string
}
