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
func (s *MemoriesService) Search(ctx context.Context, req *SearchMemoriesRequest) ([]SearchResult, error) {
	var resp []SearchResult
	if err := s.client.do(ctx, "POST", "/memories/search", req, &resp); err != nil {
		return nil, err
	}
	return resp, nil
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

// ExportOptions holds optional parameters for the Export method.
type ExportOptions struct {
	Format         string `json:"format,omitempty"`
	IncludeMetadata *bool `json:"include_metadata,omitempty"`
}

// Export exports a memory in the specified format (markdown, json, etc).
func (s *MemoriesService) Export(ctx context.Context, memoryID string, opts *ExportOptions) ([]byte, error) {
	q := url.Values{}
	if opts != nil {
		if opts.Format != "" {
			q.Set("format", opts.Format)
		}
		if opts.IncludeMetadata != nil {
			q.Set("include_metadata", strconv.FormatBool(*opts.IncludeMetadata))
		}
	}
	path := fmt.Sprintf("/memories/%s/export", url.PathEscape(memoryID))
	return s.client.doBytes(ctx, http.MethodGet, path, q, nil)
}

// --- Search types ---

// SearchMemoriesRequest is the request body for POST /memories/search.
type SearchMemoriesRequest struct {
	Query             string   `json:"query"`
	Mode              string   `json:"mode,omitempty"`
	Limit             int      `json:"limit,omitempty"`
	SpaceID           string   `json:"space_id,omitempty"`
	FilterLabels      []string `json:"filter_labels,omitempty"`
	UnitType          string   `json:"unit_type,omitempty"`
	IncludeEntities   *bool    `json:"include_entities,omitempty"`
	EventDateFrom     string   `json:"event_date_from,omitempty"`
	EventDateTo       string   `json:"event_date_to,omitempty"`
	TemporalContext   string   `json:"temporal_context,omitempty"`
	RecordedDateFrom  string   `json:"recorded_date_from,omitempty"`
	RecordedDateTo    string   `json:"recorded_date_to,omitempty"`
}

// SearchResult is a single search result.
type SearchResult struct {
	Memory             Memory           `json:"memory"`
	SimilarityScore    float64          `json:"similarity_score"`
	RelevanceReason    string           `json:"relevance_reason,omitempty"`
	RelatedEntities    []Entity         `json:"related_entities,omitempty"`
	EvolvesContext     map[string]any   `json:"evolves_context,omitempty"`
	RelatedMemoryLinks []map[string]any `json:"related_memory_links,omitempty"`
}

// --- Bulk types ---

// BulkMemorySelection describes a selection of memories for bulk operations.
type BulkMemorySelection struct {
	MemoryIDs []string `json:"memory_ids,omitempty"`
	SpaceID   string   `json:"space_id,omitempty"`
	SelectAll bool     `json:"select_all,omitempty"`
}

// BulkMovePreviewRequest is the request for POST /memories/bulk/move/preview.
type BulkMovePreviewRequest struct {
	Selection     BulkMemorySelection `json:"selection"`
	TargetSpaceID string              `json:"target_space_id"`
}

// BulkMovePreviewResponse is the response for POST /memories/bulk/move/preview.
type BulkMovePreviewResponse struct {
	Count          int    `json:"count"`
	MaxAllowed     int    `json:"max_allowed,omitempty"`
	LimitExceeded  bool   `json:"limit_exceeded,omitempty"`
	SourceSpaceID  string `json:"source_space_id,omitempty"`
	TargetSpaceID  string `json:"target_space_id,omitempty"`
	SelectionMode  string `json:"selection_mode,omitempty"`
	ExcludedCount  int    `json:"excluded_count,omitempty"`
	Message        string `json:"message,omitempty"`
}

// BulkMoveRequest is the request for POST /memories/bulk/move.
type BulkMoveRequest struct {
	Selection     BulkMemorySelection `json:"selection"`
	TargetSpaceID string              `json:"target_space_id"`
}

// BulkMoveResponse is the response for POST /memories/bulk/move.
type BulkMoveResponse struct {
	MovedCount        int    `json:"moved_count"`
	FailedCount       int    `json:"failed_count"`
	SourceSpaceID     string `json:"source_space_id,omitempty"`
	TargetSpaceID     string `json:"target_space_id,omitempty"`
	IndexUpdatedCount int    `json:"index_updated_count,omitempty"`
	IndexRepairNeeded bool   `json:"index_repair_needed,omitempty"`
	Message           string `json:"message,omitempty"`
}

// BulkDeleteRequest is the request for POST /memories/bulk/delete.
type BulkDeleteRequest struct {
	Selection     BulkMemorySelection `json:"selection"`
	CascadeDelete bool                `json:"cascade_delete,omitempty"`
}

// BulkDeleteResponse is the response for POST /memories/bulk/delete.
type BulkDeleteResponse struct {
	DeletedCount  int              `json:"deleted_count"`
	FailedCount   int              `json:"failed_count"`
	SourceSpaceID string           `json:"source_space_id,omitempty"`
	CascadeDelete bool             `json:"cascade_delete,omitempty"`
	Results       []map[string]any `json:"results,omitempty"`
	Message       string           `json:"message,omitempty"`
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

// --- Supersede / Deprecate ---

// SupersedeMemoryRequest is the request for POST /memories/{id}/supersede.
type SupersedeMemoryRequest struct {
	NewerMemoryID string `json:"newer_memory_id"`
	Reason        string `json:"reason,omitempty"`
	SpaceID       string `json:"space_id,omitempty"`
}

// DeprecateMemoryRequest is the request for POST /memories/{id}/deprecate.
type DeprecateMemoryRequest struct {
	Reason              string `json:"reason,omitempty"`
	ReplacementMemoryID string `json:"replacement_memory_id,omitempty"`
	SpaceID             string `json:"space_id,omitempty"`
}

// Supersede marks a memory as replaced by a newer memory.
func (s *MemoriesService) Supersede(ctx context.Context, memoryID string, req *SupersedeMemoryRequest) error {
	path := fmt.Sprintf("/memories/%s/supersede", url.PathEscape(memoryID))
	return s.client.do(ctx, "POST", path, req, nil)
}

// Deprecate marks a memory as deprecated while preserving it for graph history.
func (s *MemoriesService) Deprecate(ctx context.Context, memoryID string, req *DeprecateMemoryRequest) error {
	path := fmt.Sprintf("/memories/%s/deprecate", url.PathEscape(memoryID))
	return s.client.do(ctx, "POST", path, req, nil)
}
