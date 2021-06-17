package bitkub

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const (
	defaultBaseURL = "https://api.bitkub.com/"
	userAgent      = "go-bitkub"
	headerApiKey   = "X-BTK-APIKEY"
)

type Client struct {
	client *http.Client

	BaseURL   *url.URL
	UserAgent string

	Credentials *Credentials

	Server *ServerService
	Market *MarketService

	common service
	nonce  uint64
}

func NewClient(options ...*Options) *Client {
	var opts Options
	if len(options) > 0 {
		opts = *options[0]
	}

	c := new(Client)
	if opts.Client != nil {
		c.client = opts.Client
	} else {
		c.client = new(http.Client)
	}
	if opts.BaseURL != nil {
		c.BaseURL = opts.BaseURL
	} else {
		c.BaseURL, _ = url.Parse(defaultBaseURL)
	}
	if opts.UserAgent != "" {
		c.UserAgent = opts.UserAgent
	} else {
		c.UserAgent = userAgent
	}
	if opts.Credentials != nil {
		c.Credentials = opts.Credentials
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

func (c *Client) do(ctx context.Context, req *http.Request, output interface{}) (*Response, error) {
	req = req.WithContext(ctx)
	// TODO do something with http.Response?
	httpRes, err := c.client.Do(req)
	if err != nil {
		// TODO return ctx?.err
		return nil, err
	}

	res := new(Response)
	res.Result = output
	if err := json.NewDecoder(httpRes.Body).Decode(res); err != nil {
		return nil, err
	}

	if res.Error != 0 {
		return res, newBtkError(res.Error)
	}
	return res, nil
}

func (c *Client) fetch(endpoint string, ctx context.Context, input map[string]interface{}, output interface{}) error {
	u, err := url.Parse(endpoint)
	if err != nil {
		return err
	}
	if input != nil {
		q := url.Values{}
		for k, v := range input {
			if vs, ok := v.(string); ok {
				q.Set(k, vs)
			} else if vs, ok := v.(fmt.Stringer); ok {
				q.Set(k, vs.String())
			} else {
				q.Set(k, fmt.Sprintf("%v", v))
			}
		}
		u.RawQuery = q.Encode()
	}
	req, err := c.request(http.MethodGet, u.String(), nil)
	if err != nil {
		return err
	}
	_, err = c.do(ctx, req, output)
	return err
}
