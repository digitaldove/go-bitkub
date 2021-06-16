package bitkub

import "context"

type MarketService service

func (s *MarketService) Wallet(ctx context.Context, creds *Credentials) (map[string]float32, error) {
	res := make(map[string]float32)
	if err := s.client.fetchSecure("/api/market/wallet", ctx, creds, nil, &res); err != nil {
		return nil, err
	}
	return res, nil
}
