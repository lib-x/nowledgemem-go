package nowledgemem

import "context"

// SearchIndexService handles search index operations.
type SearchIndexService struct {
	client *Client
}

// GetStatus returns status of LanceDB and hybrid search.
func (s *SearchIndexService) GetStatus(ctx context.Context) (*SearchIndexStatus, error) {
	var resp SearchIndexStatus
	if err := s.client.do(ctx, "GET", "/search-index/status", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Reindex rebuilds the search index from the database.
func (s *SearchIndexService) Reindex(ctx context.Context) (*SearchReindexResponse, error) {
	var resp SearchReindexResponse
	if err := s.client.do(ctx, "POST", "/search-index/reindex", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetReindexStatus returns status of memories needing reindex.
func (s *SearchIndexService) GetReindexStatus(ctx context.Context) (*ReindexStatus, error) {
	var resp ReindexStatus
	if err := s.client.do(ctx, "GET", "/search-index/reindex/status", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// --- Search Index types ---

// SearchReindexResponse is the response for POST /search-index/reindex.
type SearchReindexResponse struct {
	Success            bool     `json:"success"`
	Memories           int      `json:"memories,omitempty"`
	Messages           int      `json:"messages,omitempty"`
	Sources            int      `json:"sources,omitempty"`
	SourceChunks       int      `json:"source_chunks,omitempty"`
	Communities        int      `json:"communities,omitempty"`
	Entities           int      `json:"entities,omitempty"`
	Errors             []string `json:"errors,omitempty"`
	Message            string   `json:"message,omitempty"`
	RestartRecommended bool     `json:"restart_recommended,omitempty"`
}

// SearchIndexStatus is the response for GET /search-index/status.
type SearchIndexStatus struct {
	Ready       bool   `json:"ready"`
	IndexType   string `json:"index_type"`
	VectorCount int    `json:"vector_count"`
	IndexPath   string `json:"index_path,omitempty"`
}

// ReindexStatus is the response for GET /search-index/reindex/status.
type ReindexStatus struct {
	Total      int `json:"total"`
	NeedsIndex int `json:"needs_index"`
}
