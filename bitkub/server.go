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

func (s *ServerService) Status(ctx context.Context) ([]EndpointStatus, error) {
	var res []EndpointStatus
	if err := s.client.fetch("/api/status", ctx, nil, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (s *ServerService) Time(ctx context.Context) (time.Time, error) {
	var res int64
	if err := s.client.fetch("/api/servertime", ctx, nil, &res); err != nil {
		return time.Time{}, err
	}
	return time.Unix(res, 0), nil
}
