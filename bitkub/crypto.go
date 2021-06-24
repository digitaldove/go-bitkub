package bitkub

import "context"

type CryptoService service

type CryptoDeposit struct {
	Hash          string      `json:"hash"`
	Currency      string      `json:"currency"`
	Amount        float64     `json:"amount"`
	Address       interface{} `json:"address"`
	Confirmations int         `json:"confirmations"`
	Status        string      `json:"status"`
	Time          Timestamp   `json:"time"`
}

type CryptoDepositHistoryRequest struct {
	Pagination Pagination
}

func (s *CryptoService) DepositHistory(ctx context.Context, req *CryptoDepositHistoryRequest) ([]*CryptoDeposit, error) {
	var output []*CryptoDeposit
	if err := s.client.fetchSecureList("/api/crypto/deposit-history", ctx, &req.Pagination, nil, &output); err != nil {
		return nil, err
	}
	return output, nil
}
