package bitkub

import "context"

type UserService service

type UserLimits struct {
	// Limits by KYC level
	Limits struct {
		// Crypto uses BTC value equivalent
		Crypto Limit `json:"crypto"`
		Fiat   Limit `json:"fiat"`
	} `json:"limits"`

	// Today's Usage
	Usage struct {
		// Crypto uses BTC value equivalent
		Crypto struct {
			Limit
			DepositPercentage     float64 `json:"deposit_percentage"`
			WithdrawPercentage    float64 `json:"withdraw_percentage"`
			DepositTHBEquivalent  float64 `json:"deposit_thb_equivalent"`
			WithdrawTHBEquivalent float64 `json:"withdraw_thb_equivalent"`
		} `json:"crypto"`
		Fiat struct {
			Limit
			DepositPercentage  float64 `json:"deposit_percentage"`
			WithdrawPercentage float64 `json:"withdraw_percentage"`
		} `json:"fiat"`
	} `json:"usage"`

	// Current THB Rate used to calculate BTC value equivalent
	Rate float64 `json:"rate"`
}

type Limit struct {
	Deposit  float64 `json:"deposit"`
	Withdraw float64 `json:"withdraw"`
}

func (s *UserService) GetLimits(ctx context.Context) (*UserLimits, error) {
	var res UserLimits
	if err := s.client.fetchSecureContext(ctx, "/api/user/limits", nil, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (s *UserService) GetTradingCredits(ctx context.Context) (float64, error) {
	var res float64
	if err := s.client.fetchSecureContext(ctx, "/api/user/trading-credits", nil, &res); err != nil {
		return 0, err
	}
	return res, nil
}
