package nowledgemem

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

// SourcesService handles source/library operations.
type SourcesService struct {
	client *Client
}

// List returns sources with optional filtering and pagination.
func (s *SourcesService) List(ctx context.Context, params *ListSourcesParams) (*ListSourcesResponse, error) {
	q := url.Values{}
	if params != nil {
		if params.Limit > 0 {
			q.Set("limit", strconv.Itoa(params.Limit))
		}
		if params.Offset > 0 {
			q.Set("offset", strconv.Itoa(params.Offset))
		}
		if params.SourceType != "" {
			q.Set("source_type", params.SourceType)
		}
		if params.LifecycleState != "" {
			q.Set("lifecycle_state", params.LifecycleState)
		}
		if params.SpaceID != "" {
			q.Set("space_id", params.SpaceID)
		}
	}
	var resp ListSourcesResponse
	if err := s.client.doQuery(ctx, "/sources", q, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Search performs full-text search across source names and content.
func (s *SourcesService) Search(ctx context.Context, query string, limit int) ([]Source, error) {
	q := url.Values{}
	q.Set("q", query)
	if limit > 0 {
		q.Set("limit", strconv.Itoa(limit))
	}
	var resp []Source
	if err := s.client.doQuery(ctx, "/sources/search", q, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// Get returns source detail with related memories and revision chain.
func (s *SourcesService) Get(ctx context.Context, sourceID string) (*Source, error) {
	var resp Source
	path := fmt.Sprintf("/sources/%s", url.PathEscape(sourceID))
	if err := s.client.do(ctx, "GET", path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetContent reads the parsed markdown content of a source.
func (s *SourcesService) GetContent(ctx context.Context, sourceID string) (string, error) {
	var resp struct {
		Content string `json:"content"`
	}
	path := fmt.Sprintf("/sources/%s/content", url.PathEscape(sourceID))
	if err := s.client.do(ctx, "GET", path, nil, &resp); err != nil {
		return "", err
	}
	return resp.Content, nil
}

// Delete deletes a source and its search index records.
func (s *SourcesService) Delete(ctx context.Context, sourceID string) error {
	path := fmt.Sprintf("/sources/%s", url.PathEscape(sourceID))
	return s.client.do(ctx, "DELETE", path, nil, nil)
}

// Update updates source lifecycle state.
func (s *SourcesService) Update(ctx context.Context, sourceID string, req *UpdateSourceRequest) (*Source, error) {
	var resp Source
	path := fmt.Sprintf("/sources/%s", url.PathEscape(sourceID))
	if err := s.client.do(ctx, "PATCH", path, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// IngestFile ingests a file through the full source pipeline.
func (s *SourcesService) IngestFile(ctx context.Context, req *IngestFileRequest) (*Source, error) {
	var resp Source
	if err := s.client.do(ctx, "POST", "/sources/ingest/file", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// IngestURL fetches a URL and ingests through the source pipeline.
func (s *SourcesService) IngestURL(ctx context.Context, req *IngestURLRequest) (*Source, error) {
	var resp Source
	if err := s.client.do(ctx, "POST", "/sources/ingest/url", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// IngestByPath ingests a file by local filesystem path.
func (s *SourcesService) IngestByPath(ctx context.Context, req *IngestByPathRequest) (*Source, error) {
	var resp Source
	if err := s.client.do(ctx, "POST", "/sources/ingest/file-path", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// BatchIngest ingests a batch of files.
func (s *SourcesService) BatchIngest(ctx context.Context, req *BatchIngestRequest) (*BatchIngestResponse, error) {
	var resp BatchIngestResponse
	if err := s.client.do(ctx, "POST", "/sources/ingest/batch", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// --- Source request types ---

// UpdateSourceRequest is the request body for PATCH /sources/{id}.
type UpdateSourceRequest struct {
	LifecycleState string `json:"lifecycle_state,omitempty"`
}

// IngestFileRequest is the request body for POST /sources/ingest/file.
type IngestFileRequest struct {
	SpaceID string `json:"space_id,omitempty"`
	// File upload requires multipart form - use a custom HTTP client for this.
}

// IngestURLRequest is the request body for POST /sources/ingest/url.
type IngestURLRequest struct {
	URL     string `json:"url"`
	SpaceID string `json:"space_id,omitempty"`
}

// IngestByPathRequest is the request body for POST /sources/ingest/file-path.
type IngestByPathRequest struct {
	FilePath string `json:"file_path"`
	SpaceID  string `json:"space_id,omitempty"`
}

// BatchIngestRequest is the request body for POST /sources/ingest/batch.
type BatchIngestRequest struct {
	FilePaths []string `json:"file_paths"`
	SpaceID   string   `json:"space_id,omitempty"`
}

// BatchIngestResponse is the response body for POST /sources/ingest/batch.
type BatchIngestResponse struct {
	Sources []Source `json:"sources"`
	Failed  []string `json:"failed,omitempty"`
}

// UpdateContent updates the parsed markdown content of a source.
func (s *SourcesService) UpdateContent(ctx context.Context, sourceID, content string) error {
	path := fmt.Sprintf("/sources/%s/content", url.PathEscape(sourceID))
	body := map[string]string{"content": content}
	return s.client.do(ctx, "PUT", path, body, nil)
}

// Extract triggers knowledge extraction from a source.
func (s *SourcesService) Extract(ctx context.Context, sourceID string) error {
	path := fmt.Sprintf("/sources/%s/extract", url.PathEscape(sourceID))
	return s.client.do(ctx, "POST", path, nil, nil)
}

// Refetch re-fetches a URL source's content and re-parse.
func (s *SourcesService) Refetch(ctx context.Context, sourceID string) (*Source, error) {
	var resp Source
	path := fmt.Sprintf("/sources/%s/refetch", url.PathEscape(sourceID))
	if err := s.client.do(ctx, "POST", path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetLabels returns labels assigned to a source.
func (s *SourcesService) GetLabels(ctx context.Context, sourceID string) ([]Label, error) {
	var resp []Label
	path := fmt.Sprintf("/sources/%s/labels", url.PathEscape(sourceID))
	if err := s.client.do(ctx, "GET", path, nil, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// AssignLabel assigns a label to a source.
func (s *SourcesService) AssignLabel(ctx context.Context, sourceID, labelID string) error {
	path := fmt.Sprintf("/sources/%s/labels/%s", url.PathEscape(sourceID), url.PathEscape(labelID))
	return s.client.do(ctx, "POST", path, nil, nil)
}

// RemoveLabel removes a label from a source.
func (s *SourcesService) RemoveLabel(ctx context.Context, sourceID, labelID string) error {
	path := fmt.Sprintf("/sources/%s/labels/%s", url.PathEscape(sourceID), url.PathEscape(labelID))
	return s.client.do(ctx, "DELETE", path, nil, nil)
}

// IngestFolderUpload uploads a folder preserving relative paths.
func (s *SourcesService) IngestFolderUpload(ctx context.Context, req *IngestFolderUploadRequest) (*BatchIngestResponse, error) {
	var resp BatchIngestResponse
	if err := s.client.do(ctx, "POST", "/sources/ingest/folder-upload", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// IngestFolderSummary returns a summary of a folder before ingestion.
func (s *SourcesService) IngestFolderSummary(ctx context.Context, req *IngestFolderSummaryRequest) (*IngestFolderSummaryResponse, error) {
	var resp IngestFolderSummaryResponse
	if err := s.client.do(ctx, "POST", "/sources/ingest/folder-summary", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// --- Source extended types ---

// IngestFolderUploadRequest is the request for POST /sources/ingest/folder-upload.
type IngestFolderUploadRequest struct {
	Files   []string `json:"files"`
	SpaceID string   `json:"space_id,omitempty"`
}

// IngestFolderSummaryRequest is the request for POST /sources/ingest/folder-summary.
type IngestFolderSummaryRequest struct {
	Paths []string `json:"paths"`
}

// IngestFolderSummaryResponse is the response for POST /sources/ingest/folder-summary.
type IngestFolderSummaryResponse struct {
	TotalFiles int      `json:"total_files"`
	TotalSize  int64    `json:"total_size"`
	FileTypes  map[string]int `json:"file_types"`
}

// GetRawFile serves the raw source file for native preview.
func (s *SourcesService) GetRawFile(ctx context.Context, sourceID string) ([]byte, error) {
	u := s.client.baseURL.ResolveReference(&url.URL{Path: fmt.Sprintf("/sources/%s/raw", url.PathEscape(sourceID))})
	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	resp, err := s.client.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		var apiErr APIError
		if err := json.NewDecoder(resp.Body).Decode(&apiErr); err != nil {
			return nil, fmt.Errorf("API error (status %d)", resp.StatusCode)
		}
		return nil, &apiErr
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}
	return data, nil
}

// GetImage serves an extracted image from a source.
func (s *SourcesService) GetImage(ctx context.Context, sourceID, filename string) ([]byte, error) {
	u := s.client.baseURL.ResolveReference(&url.URL{Path: fmt.Sprintf("/sources/%s/images/%s", url.PathEscape(sourceID), url.PathEscape(filename))})
	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	resp, err := s.client.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		var apiErr APIError
		if err := json.NewDecoder(resp.Body).Decode(&apiErr); err != nil {
			return nil, fmt.Errorf("API error (status %d)", resp.StatusCode)
		}
		return nil, &apiErr
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}
	return data, nil
}
