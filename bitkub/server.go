package bitkub

import (
	"context"
	"time"
)

type ServerService service

type EndpointStatus struct {
	Name    string `json:"name"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

// Status returns server status
func (s *ServerService) Status(ctx context.Context) ([]EndpointStatus, error) {
	var res []EndpointStatus
	if err := s.client.fetch(ctx, "/api/status", nil, &res); err != nil {
		return nil, err
	}
	return res, nil
}

// Time returns the server time in UNIX format
func (s *ServerService) Time(ctx context.Context) (time.Time, error) {
	var res int64
	if err := s.client.fetch(ctx, "/api/servertime", nil, &res); err != nil {
		return time.Time{}, err
	}
	return time.Unix(res, 0), nil
}
