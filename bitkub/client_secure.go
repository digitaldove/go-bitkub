package bitkub

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strconv"
	"sync/atomic"
	"time"
)

type Response struct {
	Error  int
	Result interface{}
}

func (c *Client) nextNonce() uint64 {
	return atomic.AddUint64(&c.nonce, 1)
}

func (c *Client) fetchSecure(endpoint string, ctx context.Context, creds *Credentials, input, output interface{}) error {
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	if input != nil {
		if err := enc.Encode(input); err != nil {
			return err
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

	h := hmac.New(sha256.New, creds.Secret)
	if _, err := bytes.NewReader(buf.Bytes()).WriteTo(h); err != nil {
		return err
	}
	buf.Truncate(buf.Len() - 1) // remove the closing }
	buf.WriteString(`,"sig":"`)
	buf.WriteString(hex.EncodeToString(h.Sum(nil)))
	buf.WriteString(`"}`)

	req, err := c.request(http.MethodPost, endpoint, buf)
	if err != nil {
		return err
	}

	req.Header.Set(headerApiKey, creds.Key)

	var res = Response{
		Result: output,
	}
	if _, err := c.do(ctx, req, &res); err != nil {
		return err
	}

	if res.Error != 0 {
		return newBtkError(res.Error)
	}

	return nil
}
