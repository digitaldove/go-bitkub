package bitkub

import "context"

type FiatService service

type FiatDeposit struct {
	TransactionID string    `json:"txn_id"`
	Currency      string    `json:"currency"`
	Amount        float64   `json:"amount"`
	Status        string    `json:"status"`
	Time          Timestamp `json:"time"`
}

type FiatDepositHistoryRequest struct {
	Pagination Pagination
}

// DepositHistory lists the fiat deposit history. It uses pagination.
func (s *FiatService) DepositHistory(ctx context.Context, req *FiatDepositHistoryRequest) ([]*FiatDeposit, error) {
	var output []*FiatDeposit
	if err := s.client.fetchSecureList("/api/fiat/deposit-history", ctx, &req.Pagination, nil, &output); err != nil {
		return nil, err
	}
	return output, nil
}
