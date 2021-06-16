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
	req, err := s.client.reqGet("/api/status")
	if err != nil {
		return nil, err
	}
	var res []EndpointStatus
	_, err = s.client.Do(ctx, req, &res)
	return res, err
}

func (s *ServerService) Time(ctx context.Context) (time.Time, error) {
	req, err := s.client.reqGet("/api/servertime")
	if err != nil {
		return time.Time{}, err
	}
	var res int64
	_, err = s.client.Do(ctx, req, &res)
	return time.Unix(res, 0), err
}
