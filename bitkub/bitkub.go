package bitkub

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	defaultBaseURL = "https://api.bitkub.com/"
	userAgent      = "go-bitkub"
	headerApiKey   = "X-BTK-APIKEY"
)

type Client struct {
	ctx    context.Context
	client *http.Client

	BaseURL   *url.URL
	UserAgent string

	Server *ServerService
	Market *MarketService

	common service
	nonce  uint64
}

func NewClient(options ...interface{}) *Client {
	baseURL, _ := url.Parse(defaultBaseURL)
	c := &Client{
		client: &http.Client{
			Timeout: 1 * time.Second,
		},
		BaseURL:   baseURL,
		UserAgent: userAgent,
	}
	c.common.client = c
	c.Server = (*ServerService)(&c.common)
	c.Market = (*MarketService)(&c.common)
	return c
}

type service struct {
	client *Client
}

func (c *Client) request(method, endpoint string, body io.ReadWriter) (*http.Request, error) {
	// https://github.com/google/go-github/blob/master/github/github.go
	u, err := c.BaseURL.Parse(endpoint)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.UserAgent)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return req, nil
}

func (c *Client) reqGet(endpoint string) (*http.Request, error) {
	return c.request(http.MethodGet, endpoint, nil)
}

func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) (*http.Response, error) {
	req = req.WithContext(ctx)
	res, err := c.client.Do(req)
	if err != nil {
		// TODO return ctx?.err
		return nil, err
	}
	return res, json.NewDecoder(res.Body).Decode(v)
}
