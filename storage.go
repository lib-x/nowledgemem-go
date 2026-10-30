package nowledgemem

import (
	"context"

	"github.com/lib-x/nowledgemem-go/openapi"
)

// StorageService handles storage operations.
type StorageService struct {
	client *Client
}

// Info returns on-disk sizes for the database and search index.
func (s *StorageService) Info(ctx context.Context) (*StorageInfo, error) {
	var resp StorageInfo
	if err := s.client.do(ctx, "GET", "/storage/info", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Optimize compacts search index and flushes database changes.
func (s *StorageService) Optimize(ctx context.Context) error {
	return s.client.do(ctx, "POST", "/storage/optimize", nil, nil)
}

// OptimizeResult returns the full optimization outcome, including partial results.
func (s *StorageService) OptimizeResult(ctx context.Context) (*openapi.OptimizeStorageResponse, error) {
	var response openapi.OptimizeStorageResponse
	if err := s.client.do(ctx, "POST", "/storage/optimize", nil, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

// --- Storage types ---

// StorageInfo is the response for GET /storage/info.
type StorageInfo struct {
	GraphDBBytes     int64 `json:"graph_db_bytes"`
	SearchIndexBytes int64 `json:"search_index_bytes"`
	TotalBytes       int64 `json:"total_bytes"`
}
