package bitkub

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
