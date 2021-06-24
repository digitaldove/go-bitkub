package bitkub

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/url"
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

func (c *Client) nextNonce() uint64 {
	return atomic.AddUint64(&c.nonce, 1)
}

func (c *Client) fetchSecure2(endpoint string, ctx context.Context, input, output interface{}) (*Response, error) {
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

	req.Header.Set(headerApiKey, creds.Key)

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

func (c *Client) fetchSecure(endpoint string, ctx context.Context, input, output interface{}) error {
	_, err := c.fetchSecure2(endpoint, ctx, input, output)
	return err
}

func (c *Client) fetchSecureList(endpoint string, ctx context.Context,
	pagination *Pagination, input, output interface{}) error {
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
	raw, err := c.fetchSecure2(endpoint, ctx, input, output)
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
