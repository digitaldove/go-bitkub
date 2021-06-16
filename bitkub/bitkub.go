package bitkub

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"time"
)

const (
	defaultBaseURL = "https://api.bitkub.com/"
	userAgent      = "go-bitkub"
)

type Client struct {
	ctx    context.Context
	client *http.Client

	BaseURL   *url.URL
	UserAgent string

	Server *ServerService

	common service
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
	return c
}

type service struct {
	client *Client
}

func (c *Client) reqGet(endpoint string) (*http.Request, error) {
	u, err := c.BaseURL.Parse(endpoint)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.UserAgent)
	return req, nil
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
