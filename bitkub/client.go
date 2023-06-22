package bitkub

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

const (
	defaultBaseURL = "https://api.bitkub.com/"
	userAgent      = "go-bitkub"
	headerAPIKey   = "X-BTK-APIKEY"
)

// TODO a zero-valued client should be usable?

// A Client is a Bitkub client. Its zero value is not usable. To get a usable client, bitkub.NewClient must be used.
type Client struct {
	client *http.Client

	BaseURL   *url.URL
	UserAgent string

	Credentials *Credentials

	Server *ServerService
	Market *MarketService
	Fiat   *FiatService
	Crypto *CryptoService
	User   *UserService

	common service
	nonce  uint64
}

// NewClient creates a new bitkub.Client with the specified options.
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
	c.Fiat = (*FiatService)(&c.common)
	c.Crypto = (*CryptoService)(&c.common)
	c.User = (*UserService)(&c.common)

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

func (c *Client) do(ctx context.Context, req *http.Request, output interface{}) error {
	req = req.WithContext(ctx)
	// TODO do something with http.Response?
	httpRes, err := c.client.Do(req)
	if err != nil {
		// TODO return ctx?.err
		return err
	}

	buf := bytes.Buffer{}
	if _, err = buf.ReadFrom(httpRes.Body); err != nil {
		return err
	}
	defer httpRes.Body.Close()

	err = json.Unmarshal(buf.Bytes(), output)
	if err != nil {
		fmt.Fprintf(os.Stderr, "DBG Client.do unmarshal error %T, data: %s\n", err, buf.String())
	}
	return err
}

func (c *Client) fetch(ctx context.Context, endpoint string, input map[string]interface{}, output interface{}) error {
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
	return c.do(ctx, req, output)
}
