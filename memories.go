package nowledgemem

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// MemoriesService handles memory CRUD operations.
type MemoriesService struct {
	client *Client
}

// List returns memories with filtering and pagination.
func (s *MemoriesService) List(ctx context.Context, params *ListMemoriesParams) (*ListMemoriesResponse, error) {
	q := url.Values{}
	if params != nil {
		if params.Limit > 0 {
			q.Set("limit", strconv.Itoa(params.Limit))
		}
		if params.Offset > 0 {
			q.Set("offset", strconv.Itoa(params.Offset))
		}
		if params.State != "" {
			q.Set("state", params.State)
		}
		if params.ImportanceMin > 0 {
			q.Set("importance_min", strconv.FormatFloat(params.ImportanceMin, 'f', -1, 64))
		}
		if params.SpaceID != "" {
			q.Set("space_id", params.SpaceID)
		}
		if params.IsCrystal != nil {
			q.Set("is_crystal", strconv.FormatBool(*params.IsCrystal))
		}
	}
	var resp ListMemoriesResponse
	if err := s.client.doQuery(ctx, "/memories", q, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Create creates a new memory with automatic entity extraction.
func (s *MemoriesService) Create(ctx context.Context, req *CreateMemoryRequest) (*CreateMemoryResponse, error) {
	var resp CreateMemoryResponse
	if err := s.client.do(ctx, "POST", "/memories", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Get retrieves a specific memory by ID.
func (s *MemoriesService) Get(ctx context.Context, memoryID string, spaceID string) (*MemoryListItem, error) {
	q := url.Values{}
	if spaceID != "" {
		q.Set("space_id", spaceID)
	}
	var resp MemoryListItem
	path := fmt.Sprintf("/memories/%s", url.PathEscape(memoryID))
	if err := s.client.doQuery(ctx, path, q, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Update updates memory properties like importance, title, and content.
func (s *MemoriesService) Update(ctx context.Context, memoryID string, updates map[string]any) (*MemoryListItem, error) {
	var resp MemoryListItem
	path := fmt.Sprintf("/memories/%s", url.PathEscape(memoryID))
	if err := s.client.do(ctx, "PATCH", path, updates, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Delete deletes a memory and optionally its relationships.
func (s *MemoriesService) Delete(ctx context.Context, memoryID string, params *DeleteMemoryParams) (*DeleteMemoryResponse, error) {
	q := url.Values{}
	if params != nil {
		q.Set("cascade_delete", strconv.FormatBool(params.CascadeDelete))
		if params.SpaceID != "" {
			q.Set("space_id", params.SpaceID)
		}
	}
	var resp DeleteMemoryResponse
	path := fmt.Sprintf("/memories/%s", url.PathEscape(memoryID))
	if err := s.client.doWithQuery(ctx, http.MethodDelete, path, q, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Search performs a hybrid search across memories.
func (s *MemoriesService) Search(ctx context.Context, req *SearchMemoriesRequest) (*SearchMemoriesResponse, error) {
	var resp SearchMemoriesResponse
	if err := s.client.do(ctx, "POST", "/memories/search", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// BulkMovePreview previews a bulk move between spaces.
func (s *MemoriesService) BulkMovePreview(ctx context.Context, req *BulkMovePreviewRequest) (*BulkMovePreviewResponse, error) {
	var resp BulkMovePreviewResponse
	if err := s.client.do(ctx, "POST", "/memories/bulk/move/preview", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// BulkMove moves selected memories into another space.
func (s *MemoriesService) BulkMove(ctx context.Context, req *BulkMoveRequest) (*BulkMoveResponse, error) {
	var resp BulkMoveResponse
	if err := s.client.do(ctx, "POST", "/memories/bulk/move", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// BulkDelete deletes selected memories.
func (s *MemoriesService) BulkDelete(ctx context.Context, req *BulkDeleteRequest) (*BulkDeleteResponse, error) {
	var resp BulkDeleteResponse
	if err := s.client.do(ctx, "POST", "/memories/bulk/delete", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ToggleFavorite toggles favorite status for a memory.
func (s *MemoriesService) ToggleFavorite(ctx context.Context, memoryID string) (*ToggleFavoriteResponse, error) {
	var resp ToggleFavoriteResponse
	path := fmt.Sprintf("/memories/%s/favorite", url.PathEscape(memoryID))
	if err := s.client.do(ctx, "POST", path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetLabels returns labels assigned to a memory.
func (s *MemoriesService) GetLabels(ctx context.Context, memoryID string) ([]Label, error) {
	var resp []Label
	path := fmt.Sprintf("/memories/%s/labels", url.PathEscape(memoryID))
	if err := s.client.do(ctx, "GET", path, nil, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// AssignLabel assigns a label to a memory.
func (s *MemoriesService) AssignLabel(ctx context.Context, memoryID, labelID string) error {
	path := fmt.Sprintf("/memories/%s/labels/%s", url.PathEscape(memoryID), url.PathEscape(labelID))
	return s.client.do(ctx, "POST", path, nil, nil)
}

// RemoveLabel removes a label from a memory.
func (s *MemoriesService) RemoveLabel(ctx context.Context, memoryID, labelID string) error {
	path := fmt.Sprintf("/memories/%s/labels/%s", url.PathEscape(memoryID), url.PathEscape(labelID))
	return s.client.do(ctx, "DELETE", path, nil, nil)
}

// Export exports a memory in the specified format (markdown, json, etc).
func (s *MemoriesService) Export(ctx context.Context, memoryID, format string) ([]byte, error) {
	q := url.Values{}
	if format != "" {
		q.Set("format", format)
	}
	path := fmt.Sprintf("/memories/%s/export", url.PathEscape(memoryID))
	return s.client.doBytes(ctx, http.MethodGet, path, q, nil)
}

// --- Search types ---

// SearchMemoriesRequest is the request body for POST /memories/search.
type SearchMemoriesRequest struct {
	Query      string   `json:"query"`
	SpaceID    string   `json:"space_id,omitempty"`
	Limit      int      `json:"limit,omitempty"`
	Offset     int      `json:"offset,omitempty"`
	Labels     []string `json:"labels,omitempty"`
	SourceType string   `json:"source_type,omitempty"`
}

// SearchMemoriesResponse is the response body for POST /memories/search.
type SearchMemoriesResponse struct {
	Results []SearchResult `json:"results"`
	Total   int            `json:"total"`
	Query   string         `json:"query"`
}

// SearchResult is a single search result.
type SearchResult struct {
	ID      string  `json:"id"`
	Title   string  `json:"title,omitempty"`
	Content string  `json:"content,omitempty"`
	Score   float64 `json:"score"`
	SpaceID string  `json:"space_id,omitempty"`
	Source  string  `json:"source,omitempty"`
}

// --- Bulk types ---

// BulkMovePreviewRequest is the request for POST /memories/bulk/move/preview.
type BulkMovePreviewRequest struct {
	MemoryIDs []string `json:"memory_ids,omitempty"`
	FromSpace string   `json:"from_space,omitempty"`
	ToSpace   string   `json:"to_space"`
}

// BulkMovePreviewResponse is the response for POST /memories/bulk/move/preview.
type BulkMovePreviewResponse struct {
	WillMove  int      `json:"will_move"`
	Conflicts []string `json:"conflicts,omitempty"`
}

// BulkMoveRequest is the request for POST /memories/bulk/move.
type BulkMoveRequest struct {
	MemoryIDs []string `json:"memory_ids,omitempty"`
	FromSpace string   `json:"from_space,omitempty"`
	ToSpace   string   `json:"to_space"`
}

// BulkMoveResponse is the response for POST /memories/bulk/move.
type BulkMoveResponse struct {
	Moved  int `json:"moved"`
	Failed int `json:"failed"`
}

// BulkDeleteRequest is the request for POST /memories/bulk/delete.
type BulkDeleteRequest struct {
	MemoryIDs []string `json:"memory_ids,omitempty"`
	SpaceID   string   `json:"space_id,omitempty"`
}

// BulkDeleteResponse is the response for POST /memories/bulk/delete.
type BulkDeleteResponse struct {
	Deleted int `json:"deleted"`
	Failed  int `json:"failed"`
}

// ToggleFavoriteResponse is the response for POST /memories/{id}/favorite.
type ToggleFavoriteResponse struct {
	IsFavorite bool `json:"is_favorite"`
}

// Reindex reindexes multiple memories or all needing reindex.
func (s *MemoriesService) Reindex(ctx context.Context, req *ReindexRequest) (*ReindexResponse, error) {
	var resp ReindexResponse
	if err := s.client.do(ctx, "POST", "/memories/reindex", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetReindexStatus returns status of memories needing reindex.
func (s *MemoriesService) GetReindexStatus(ctx context.Context) (*MemoryReindexStatus, error) {
	var resp MemoryReindexStatus
	if err := s.client.do(ctx, "GET", "/memories/reindex/status", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// --- Reindex types ---

// ReindexRequest is the request for POST /memories/reindex.
type ReindexRequest struct {
	MemoryIDs []string `json:"memory_ids,omitempty"`
	All       bool     `json:"all,omitempty"`
}

// ReindexResponse is the response for POST /memories/reindex.
type ReindexResponse struct {
	Queued  int `json:"queued"`
	Skipped int `json:"skipped"`
}

// MemoryReindexStatus is the response for GET /memories/reindex/status.
type MemoryReindexStatus struct {
	Total        int `json:"total"`
	NeedsReindex int `json:"needs_reindex"`
}
