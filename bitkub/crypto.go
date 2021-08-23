package bitkub

import "context"

type CryptoService service

type CryptoDeposit struct {
	Hash          string    `json:"hash"`
	Currency      string    `json:"currency"`
	Amount        float64   `json:"amount"`
	FromAddress   string    `json:"from_address"`
	ToAddress     string    `json:"to_address"`
	Confirmations int       `json:"confirmations"`
	Status        string    `json:"status"`
	Time          Timestamp `json:"time"`
}

type CryptoDepositHistoryRequest struct {
	Pagination Pagination
}

// DepositHistory lists all the crypto deposit history. It uses pagination.
func (s *CryptoService) DepositHistory(ctx context.Context, req *CryptoDepositHistoryRequest) ([]*CryptoDeposit, error) {
	var output []*CryptoDeposit
	if err := s.client.fetchSecureList(ctx, "/api/crypto/deposit-history", &req.Pagination, nil, &output); err != nil {
		return nil, err
	}
	return output, nil
}

type CryptoWithdraw struct {
	TransactionID string      `json:"txn_id"`
	Hash          string      `json:"hash"`
	Currency      string      `json:"currency"`
	Amount        float64     `json:"amount,string"`
	Fee           float64     `json:"fee"`
	Address       interface{} `json:"address"`
	Status        string      `json:"status"`
	Time          Timestamp   `json:"time"`
}

type CryptoWithdrawHistoryRequest struct {
	Pagination Pagination
}

// WithdrawHistory lists all the crypto withdrawal history. It uses pagination.
func (s *CryptoService) WithdrawHistory(ctx context.Context, req *CryptoWithdrawHistoryRequest) ([]*CryptoWithdraw, error) {
	var output []*CryptoWithdraw
	if err := s.client.fetchSecureList(ctx, "/api/crypto/withdraw-history", &req.Pagination, nil, &output); err != nil {
		return nil, err
	}
	return output, nil
}
