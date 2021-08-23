package bitkub

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"sync/atomic"
	"time"
)

type Response struct {
	Error      int         `json:"error"`
	Result     interface{} `json:"result"`
	Pagination Pagination  `json:"pagination,omitempty"`
}

type Pagination struct {
	InBody bool `json:"-"`
	Page   int  `json:"page"`
	Limit  int  `json:"limit,omitempty"`
	Last   int  `json:"last"`
	Next   int  `json:"next,omitempty"`
	Prev   int  `json:"prev,omitempty"`
	Done   bool `json:"-"`
}

// UnmarshalJSON override required because API does not return correct data types as documented (and inconsistent)
func (p *Pagination) UnmarshalJSON(b []byte) error {
	raw := make(map[string]interface{})
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	intHelper := func(key string, omitEmpty bool, target *int) error {
		val, ok := raw[key]
		if !ok {
			if !omitEmpty {
				// TODO error?
				return nil
			}
			return nil
		}
		switch val.(type) {
		case float64:
			// +0.5 because https://stackoverflow.com/questions/8022389/convert-a-float64-to-an-int-in-go
			*target = (int)(val.(float64) + 0.5)
		case string:
			var err error
			if *target, err = strconv.Atoi(val.(string)); err != nil {
				return err
			}
		default:
			return &json.UnmarshalTypeError{
				Value:  fmt.Sprintf("%v", val),
				Type:   reflect.TypeOf(*target),
				Struct: "Pagination",
				Field:  key,
			}
		}
		return nil
	}
	if err := intHelper("page", false, &p.Page); err != nil {
		return err
	}
	if err := intHelper("limit", true, &p.Limit); err != nil {
		return err
	}
	if err := intHelper("last", false, &p.Last); err != nil {
		return err
	}
	if err := intHelper("next", true, &p.Next); err != nil {
		return err
	}
	if err := intHelper("prev", true, &p.Prev); err != nil {
		return err
	}
	return nil
}

func (c *Client) nextNonce() uint64 {
	return atomic.AddUint64(&c.nonce, 1)
}

// TODO: context.Context should be the first parameter of the function
func (c *Client) internalSecureFetch(ctx context.Context, endpoint string, input, output interface{}) (*Response, error) {
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	if input != nil {
		if err := enc.Encode(input); err != nil {
			return nil, err
		}
		buf.Truncate(buf.Len() - 2) // remove the newline and closing }
		buf.WriteString(`,"ts":`)
	} else {
		buf.WriteString(`{"ts":`)
	}
	buf.WriteString(strconv.FormatInt(time.Now().Unix(), 10))
	/* seems like nonce is tied to the api key, so no way of knowing which nonce we have to start with right now
	buf.WriteString(`,"non":`)
	buf.WriteString(strconv.FormatUint(c.nextNonce(), 10))
	*/
	buf.WriteRune('}')

	creds := c.Credentials
	if overrideCreds, ok := ctx.Value(CtxKeyCredentials).(*Credentials); ok && overrideCreds != nil {
		creds = overrideCreds
	}
	if creds == nil {
		return nil, ErrUnauthenticated
	}
	h := hmac.New(sha256.New, creds.Secret)
	if _, err := bytes.NewReader(buf.Bytes()).WriteTo(h); err != nil {
		return nil, err
	}
	buf.Truncate(buf.Len() - 1) // remove the closing }
	buf.WriteString(`,"sig":"`)
	buf.WriteString(hex.EncodeToString(h.Sum(nil)))
	buf.WriteString(`"}`)

	req, err := c.request(http.MethodPost, endpoint, buf)
	if err != nil {
		return nil, err
	}

	req.Header.Set(headerAPIKey, creds.Key)

	res := new(Response)
	res.Result = output
	if err := c.do(ctx, req, res); err != nil {
		return nil, err
	}
	if res.Error != 0 {
		return res, newBtkError(res.Error)
	}
	return res, nil
}

func (c *Client) fetchSecure(endpoint string, input, output interface{}) error {
	return c.fetchSecureContext(context.Background(), endpoint, input, output)
}

func (c *Client) fetchSecureContext(ctx context.Context, endpoint string, input, output interface{}) error {
	_, err := c.internalSecureFetch(ctx, endpoint, input, output)
	return err
}

func (c *Client) fetchSecureList(ctx context.Context, endpoint string, pagination *Pagination, input, output interface{}) error {
	if (pagination.Page > 0 || pagination.Limit > 0) && !pagination.InBody {
		u, err := url.Parse(endpoint)
		if err != nil {
			pagination.Done = true
			return err
		}
		q := make(url.Values)
		if pagination.Page > 0 {
			q.Set("p", strconv.Itoa(pagination.Page))
		}
		if pagination.Limit > 0 {
			q.Set("lmt", strconv.Itoa(pagination.Limit))
		}
		u.RawQuery = q.Encode()
		endpoint = u.String()
	}
	raw, err := c.internalSecureFetch(ctx, endpoint, input, output)
	if err != nil {
		pagination.Done = true
		return err
	}
	pagination.Page = raw.Pagination.Page
	pagination.Last = raw.Pagination.Last
	if pagination.Page == pagination.Last {
		pagination.Done = true
	} else {
		pagination.Page++
	}
	return nil
}
