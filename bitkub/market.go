package bitkub

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
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

type TradeInfoRequest struct {
	Symbol string `json:"sym"`
	Limit  int    `json:"lmt"`
}

func (t TradeInfoRequest) Map() map[string]interface{} {
	m := make(map[string]interface{})
	m["sym"] = t.Symbol
	m["lmt"] = t.Limit
	return m
}

type Trade struct {
	Timestamp Timestamp
	Rate      float64
	Amount    float64
	Side      string
}

func (t *Trade) UnmarshalJSON(b []byte) error {
	// Bitkub API for some reason chose to return it as an array, instead of an object
	return json.Unmarshal(b, &[]interface{}{&t.Timestamp, &t.Rate, &t.Amount, &t.Side})
}

func (s *MarketService) ListTrades(ctx context.Context, req *TradeInfoRequest) ([]*Trade, error) {
	var list []*Trade
	res := Response{
		Result: &list,
	}
	if err := s.client.fetch(ctx, "/api/market/trades", req.Map(), &res); err != nil {
		return nil, err
	}
	if res.Error != 0 {
		return nil, newBtkError(res.Error)
	}
	return list, nil
}

type BidAsk struct {
	OrderId   string
	Timestamp Timestamp
	Volume    float64
	Rate      float64
	Amount    float64
}

func (t *BidAsk) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, &[]interface{}{&t.OrderId, &t.Timestamp, &t.Volume, &t.Rate, &t.Amount})
}

func (s *MarketService) ListBids(ctx context.Context, req *TradeInfoRequest) ([]*BidAsk, error) {
	var list []*BidAsk
	res := Response{Result: &list}
	if err := s.client.fetch(ctx, "/api/market/bids", req.Map(), &res); err != nil {
		return nil, err
	}
	if res.Error != 0 {
		return nil, newBtkError(res.Error)
	}
	return list, nil
}

func (s *MarketService) ListAsks(ctx context.Context, req *TradeInfoRequest) ([]*BidAsk, error) {
	var list []*BidAsk
	res := Response{Result: &list}
	if err := s.client.fetch(ctx, "/api/market/asks", req.Map(), &res); err != nil {
		return nil, err
	}
	if res.Error != 0 {
		return nil, newBtkError(res.Error)
	}
	return list, nil
}

type Book struct {
	Bids []*BidAsk `json:"bids"`
	Asks []*BidAsk `json:"asks"`
}

func (s *MarketService) OrderBook(ctx context.Context, req *TradeInfoRequest) (*Book, error) {
	var book Book
	res := Response{Result: &book}
	if err := s.client.fetch(ctx, "/api/market/books", req.Map(), &res); err != nil {
		return nil, err
	}
	if res.Error != 0 {
		return nil, newBtkError(res.Error)
	}
	return &book, nil
}

/* TODO cannot seem to get the API to return useful information, need help
type TradingViewRequest struct {
	Symbol   string
	Interval time.Duration
	From     time.Time
	To       time.Time
}

func (t TradingViewRequest) Map() map[string]interface{} {
	m := make(map[string]interface{})
	m["sym"] = t.Symbol
	m["int"] = int(t.Interval / time.Second)
	m["frm"] = NewTimestamp(t.From)
	m["to"] = NewTimestamp(t.To)
	return m
}

type TradingViewResponse struct {
}
*/

type DepthResponse struct {
	Asks []*Depth
	Bids []*Depth
}

type Depth struct {
	Price  float64
	Volume float64
}

func (d *Depth) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, &[]interface{}{&d.Price, &d.Volume})
}

func (s *MarketService) GetDepth(ctx context.Context, req *TradeInfoRequest) (*DepthResponse, error) {
	var depth DepthResponse
	if err := s.client.fetch(ctx, "/api/market/depth", req.Map(), &depth); err != nil {
		return nil, err
	}
	return &depth, nil
}

type OrderHistory struct {
	TransactionID   string    `json:"txn_id"`
	OrderID         string    `json:"order_id"`
	Hash            string    `json:"hash"`
	ParentOrderID   string    `json:"parent_order_id"`
	ParentOrderHash string    `json:"parent_order_hash"` // undocumented
	SuperOrderID    string    `json:"super_order_id"`
	SuperOrderHash  string    `json:"super_order_hash"` // undocumented
	TakenByMe       bool      `json:"taken_by_me"`
	IsMaker         bool      `json:"is_maker"`
	Side            string    `json:"side"`
	Type            string    `json:"type"`
	Rate            Float64s  `json:"rate"`
	Fee             Float64s  `json:"fee"`
	Credit          Float64s  `json:"credit"`
	Amount          Float64s  `json:"amount"`
	Timestamp       Timestamp `json:"ts"`
	// Date            time.Time `json:"date"` // undocumented // ignore because using a non-standard format
}

// Float64s is just a float64, but dum Bitkub is not consistent in its API
// which sometimes quote and unquote numbers in its responses
type Float64s float64

func (v *Float64s) UnmarshalJSON(b []byte) error {
	var s interface{}
	var err error
	if err = json.Unmarshal(b, &s); err != nil {
		return err
	}
	if s == nil {
		// TODO is this an error?
		*v = 0
		return nil
	}
	switch s.(type) {
	case string:
		*(*float64)(v), err = strconv.ParseFloat(s.(string), 64)
		if err != nil {
			return err
		}
	case float64:
		*(*float64)(v) = s.(float64)
	default:
		return fmt.Errorf("%t is not a string or a number", s)
	}
	return nil
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
	ID     string `json:"id"`
	SD     string `json:"sd"`
	Hash   string `json:"hash"`
}

// OrderInfo response
type OrderInfo struct {
	ID      string             `json:"id"`
	First   string             `json:"first"`
	Parent  string             `json:"parent"`
	Last    string             `json:"last"`
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
	ID        string    `json:"id"`
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
