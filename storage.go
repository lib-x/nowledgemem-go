package onledgemem

import "context"

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

// --- Storage types ---

// StorageInfo is the response for GET /storage/info.
type StorageInfo struct {
	GraphDBBytes     int64 `json:"graph_db_bytes"`
	SearchIndexBytes int64 `json:"search_index_bytes"`
	TotalBytes       int64 `json:"total_bytes"`
}
