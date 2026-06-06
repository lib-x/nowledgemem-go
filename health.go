package nowledgemem

import "context"

// HealthService handles health check operations.
type HealthService struct {
	client *Client
}

// Check performs a health check.
func (s *HealthService) Check(ctx context.Context) (*HealthCheck, error) {
	var resp HealthCheck
	if err := s.client.do(ctx, "GET", "/health", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ForceCheckpoint forces a database checkpoint to flush WAL to disk.
func (s *HealthService) ForceCheckpoint(ctx context.Context) error {
	return s.client.do(ctx, "POST", "/health/checkpoint", nil, nil)
}
