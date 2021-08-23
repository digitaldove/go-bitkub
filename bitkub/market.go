package bitkub

import (
	"context"
	"strings"
	"time"
)

type MarketService service

func (s *MarketService) Wallet(ctx context.Context) (map[string]float32, error) {
	res := make(map[string]float32)
	if err := s.client.fetchSecureContext(ctx, "/api/market/wallet", nil, &res); err != nil {
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
	out := new(Response)
	var res []*Symbol
	out.Result = &res
	if err := s.client.fetch(ctx, "/api/market/symbols", nil, out); err != nil {
		return nil, err
	}
	if out.Error != 0 {
		return nil, btkError{Code: out.Error}
	}
	return res, nil
}

type Ticker struct {
	ID            int     `json:"id"`
	Last          float64 `json:"last"`
	LowestAsk     float64 `json:"lowestAsk"`
	HighestBid    float64 `json:"highestBid"`
	PercentChange float64 `json:"percentChange"`
	BaseVolume    float64 `json:"baseVolume"`
	QuoteVolume   float64 `json:"quoteVolume"`
	IsFrozen      int     `json:"isFrozen"`
	High24Hr      float64 `json:"high24Hr"`
	Low24Hr       float64 `json:"low24Hr"`
}

func (s *MarketService) GetTicker(ctx context.Context, symbols ...string) (map[string]*Ticker, error) {
	res := make(map[string]*Ticker)
	input := make(map[string]interface{})
	input["sym"] = strings.Join(symbols, ",")
	if err := s.client.fetch(ctx, "/api/market/ticker", input, &res); err != nil {
		return nil, err
	}
	return res, nil
}

type OrderHistory struct {
	TransactionID   string    `json:"txn_id"`
	OrderID         int64     `json:"order_id"`
	Hash            string    `json:"hash"`
	ParentOrderID   int64     `json:"parent_order_id"`
	ParentOrderHash string    `json:"parent_order_hash"` // undocumented
	SuperOrderID    int64     `json:"super_order_id"`
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
	Pagination Pagination
}

// MyOrderHistory lists all orders that have already matched.
func (s *MarketService) MyOrderHistory(ctx context.Context, req *MyOrderHistoryRequest) ([]*OrderHistory, error) {
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
	if err := s.client.fetchSecureList(ctx, "/api/market/my-order-history", &req.Pagination, input, &output); err != nil {
		return nil, err
	}
	return output, nil
}

// OrderInfoRequest ..
type OrderInfoRequest struct {
	Symbol string `json:"sym"`
	ID     int    `json:"id"`
	SD     string `json:"sd"`
	Hash   string `json:"hash"`
}

// OrderInfo response
type OrderInfo struct {
	ID      int                `json:"id"`
	First   int                `json:"first"`
	Parent  int                `json:"parent"`
	Last    int                `json:"last"`
	Amount  int                `json:"amount"`
	Rate    int                `json:"rate"`
	Fee     int                `json:"fee"`
	Credit  int                `json:"credit"`
	Filled  float64            `json:"filled"`
	Total   int                `json:"total"`
	Status  string             `json:"status"` // can only be "filled" or "unfilled"\
	History []OrderInfoHistory `json:"history"`
}

// OrderInfoHistory shows historical data of the order
type OrderInfoHistory struct {
	Amount    float64   `json:"amount"`
	Credit    float64   `json:"credit"`
	Fee       float64   `json:"fee"`
	ID        int       `json:"id"`
	Rate      int       `json:"rate"`
	Timestamp Timestamp `json:"timestamp"`
}

// OrderInfo calls /api/market/order-info
func (s *MarketService) OrderInfo(ctx context.Context, req *OrderInfoRequest) ([]*OrderInfo, error) {
	return s.OrderInfoContext(ctx, req)
}

// OrderInfoContext calls /api/market/order-info with context deadline
func (s *MarketService) OrderInfoContext(ctx context.Context, req *OrderInfoRequest) ([]*OrderInfo, error) {
	// Since OrderInfo API required all fields per OrderInfoRequest
	var output []*OrderInfo
	if err := s.client.fetchSecureContext(ctx, "/api/market/order-info", req, &output); err != nil {
		return nil, err
	}
	return output, nil
}
