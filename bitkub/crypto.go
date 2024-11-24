package bitkub

import (
	"context"
	"net/http"
)

type CryptoService service

type CryptoDeposit struct {
	Hash          string      `json:"hash"`
	Currency      string      `json:"currency"`
	Amount        float64     `json:"amount"`
	FromAddress   string      `json:"from_address"`
	ToAddress     string      `json:"to_address"`
	Confirmations int         `json:"confirmations"`
	Status        string      `json:"status"`
	Time          TimestampV2 `json:"time"`
}

type CryptoDepositHistoryRequest struct {
	Pagination Pagination
}

// DepositHistory lists all the crypto deposit history. It uses pagination.
func (s *CryptoService) DepositHistory(ctx context.Context, req *CryptoDepositHistoryRequest) ([]*CryptoDeposit, error) {
	var output []*CryptoDeposit
	if err := s.client.fetchSecureList(ctx, http.MethodPost, "/api/v3/crypto/deposit-history", &req.Pagination, nil, &output); err != nil {
		return nil, err
	}
	return output, nil
}

type CryptoWithdraw struct {
	TransactionID     string      `json:"txn_id"`
	ExternalReference string      `json:"ext_ref"`
	Hash              string      `json:"hash"`
	Currency          string      `json:"currency"`
	Amount            float64     `json:"amount,string"`
	Fee               float64     `json:"fee"`
	Address           interface{} `json:"address"`
	Memo              string      `json:"memo"`
	Status            string      `json:"status"`
	Note              string      `json:"note"`
	Time              TimestampV2 `json:"time"`
}

type CryptoWithdrawHistoryRequest struct {
	Pagination Pagination
}

// WithdrawHistory lists all the crypto withdrawal history. It uses pagination.
func (s *CryptoService) WithdrawHistory(ctx context.Context, req *CryptoWithdrawHistoryRequest) ([]*CryptoWithdraw, error) {
	var output []*CryptoWithdraw
	if err := s.client.fetchSecureList(ctx, http.MethodPost, "/api/v3/crypto/withdraw-history", &req.Pagination, nil, &output); err != nil {
		return nil, err
	}
	return output, nil
}

type ListAddressesRequest struct {
	Pagination
}

type Address struct {
	Currency string      `json:"currency"`
	Address  string      `json:"address"`
	Tag      string      `json:"tag"`
	Time     TimestampV2 `json:"time"`
}

func (s *CryptoService) ListAddresses(ctx context.Context, req *ListAddressesRequest) ([]*Address, error) {
	var res []*Address
	if err := s.client.fetchSecureList(ctx, http.MethodPost, "/api/v3/crypto/addresses", &req.Pagination, nil, &res); err != nil {
		return nil, err
	}
	return res, nil
}

type CryptoWithdrawRequest struct {
	// Currency for withdrawal (e.g. BTC, ETH)
	Currency string `json:"cur"`

	// Amount you want to withdraw
	Amount float64 `json:"amt"`

	// Address to which you want to withdraw
	Address string `json:"adr"`

	// (Optional) Memo or destination tag to which you want to withdraw
	Memo string `json:"mem,omitempty"`

	// Cryptocurrency network to withdraw. No default value of this field. Please find the available network from the
	// link as follows. https://www.bitkub.com/fee/cryptocurrency
	Network string `json:"net"`
}

type CryptoWithdrawResult struct {
	TransactionId string      `json:"txn"`
	Address       string      `json:"adr"`
	Memo          string      `json:"mem,omitempty"`
	Currency      string      `json:"cur"`
	Amount        float64     `json:"amt"`
	Fee           float64     `json:"fee"`
	Timestamp     TimestampV2 `json:"ts"`
}

// Withdraw crypto to a trusted address.
func (s *CryptoService) Withdraw(ctx context.Context, req *CryptoWithdrawRequest) (*CryptoWithdrawResult, error) {
	var res CryptoWithdrawResult
	if err := s.client.fetchSecure(ctx, http.MethodPost, "/api/v3/crypto/withdraw", req, &res); err != nil {
		return nil, err
	}
	return &res, nil
}
