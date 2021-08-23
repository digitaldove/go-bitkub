package bitkub

import "context"

const CtxKeyCredentials CtxKey = "btk-creds"

type CtxKey string

type Credentials struct {
	Key    string
	Secret []byte
}

func NewCredentials(key, secret string) *Credentials {
	return &Credentials{
		Key:    key,
		Secret: []byte(secret),
	}
}

func OverrideCreds(ctx context.Context, creds *Credentials) context.Context {
	return context.WithValue(ctx, CtxKeyCredentials, creds)
}
