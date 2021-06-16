package bitkub

import "context"

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
