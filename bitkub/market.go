package bitkub

import (
	"context"
	"time"
)

type MarketService service

func (s *MarketService) Wallet(ctx context.Context) (map[string]float32, error) {
	res := make(map[string]float32)
	if err := s.client.fetchSecure("/api/market/wallet", ctx, nil, &res); err != nil {
		return nil, err
	}
	return res, nil
}

type Symbol struct {
	ID     int    `json:"id"`
	Symbol string `json:"symbol"`
	Info   string `json:"info"`
}

func (s *MarketService) ListSymbols(ctx context.Context) ([]*Symbol, error) {
	var res []*Symbol
	if err := s.client.fetch("/api/market/symbols", ctx, nil, &res); err != nil {
		return nil, err
	}
	return res, nil
}

type OrderHistory struct {
	TransactionID   string    `json:"txn_id"`
	OrderID         int64     `json:"order_id"`
	Hash            string    `json:"hash"`
	ParentOrderId   int64     `json:"parent_order_id"`
	ParentOrderHash string    `json:"parent_order_hash"` // undocumented
	SuperOrderId    int64     `json:"super_order_id"`
	SuperOrderHash  string    `json:"super_order_hash"` // undocumented
	TakenByMe       bool      `json:"taken_by_me"`
	IsMaker         bool      `json:"is_maker"`
	Side            string    `json:"side"`
	Type            string    `json:"type"`
	Rate            float64   `json:"rate"`
	Fee             float64   `json:"fee"`
	Credit          float64   `json:"credit"`
	Amount          float64   `json:"amount"`
	Timestamp       Timestamp `json:"ts"`
	// Date            time.Time `json:"date"` // undocumented // ignore because using a non-standard format
}

type MyOrderHistoryRequest struct {
	Symbol     string
	From       time.Time
	To         time.Time
	Pagination *Pagination
}

func (s *MarketService) MyOrderHistory(ctx context.Context, req *MyOrderHistoryRequest) ([]*OrderHistory, error) {
	// TODO err if pagination nil
	req.Pagination.InBody = true
	input := make(map[string]interface{})
	input["sym"] = req.Symbol
	if !req.From.IsZero() {
		input["start"] = req.From.Unix()
	}
	if !req.To.IsZero() {
		input["end"] = req.To.Unix()
	}
	if req.Pagination.Page > 0 {
		input["p"] = req.Pagination.Page
	}
	if req.Pagination.Limit > 0 {
		input["lmt"] = req.Pagination.Limit
	}
	var output []*OrderHistory
	if err := s.client.fetchSecureList("/api/market/my-order-history", ctx, req.Pagination, input, &output); err != nil {
		return nil, err
	}
	return output, nil
}
