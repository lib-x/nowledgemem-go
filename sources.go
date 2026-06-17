package nowledgemem

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
)

// SourcesService handles source/library operations.
//
// It provides methods for listing, creating, updating, and deleting sources,
// as well as ingesting files, URLs, and batches.
type SourcesService struct {
	client *Client
}

// List returns sources with optional filtering and pagination.
//
// GET /sources
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
//
// GET /sources/search
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
//
// GET /sources/{id}
func (s *SourcesService) Get(ctx context.Context, sourceID string) (*Source, error) {
	var resp Source
	path := fmt.Sprintf("/sources/%s", url.PathEscape(sourceID))
	if err := s.client.do(ctx, "GET", path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetContent reads the parsed markdown content of a source.
//
// GET /sources/{id}/content
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
//
// DELETE /sources/{id}
func (s *SourcesService) Delete(ctx context.Context, sourceID string, spaceID string) error {
	q := url.Values{}
	if spaceID != "" {
		q.Set("space_id", spaceID)
	}
	path := fmt.Sprintf("/sources/%s", url.PathEscape(sourceID))
	return s.client.doWithQuery(ctx, "DELETE", path, q, nil, nil)
}

// Update updates source processing state (reparse, ocr_reparse, or mark_stale).
//
// PATCH /sources/{id}
func (s *SourcesService) Update(ctx context.Context, sourceID string, req *UpdateSourceRequest, spaceID string) (*Source, error) {
	var resp Source
	q := url.Values{}
	if spaceID != "" {
		q.Set("space_id", spaceID)
	}
	path := fmt.Sprintf("/sources/%s", url.PathEscape(sourceID))
	if err := s.client.doWithQuery(ctx, "PATCH", path, q, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// IngestFile ingests a file as a new source via multipart upload.
//
// POST /sources/ingest/file
func (s *SourcesService) IngestFile(ctx context.Context, req *IngestFileRequest) (*IngestSourceResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("request is required")
	}
	if req.File == nil {
		return nil, fmt.Errorf("file is required")
	}

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	filename := req.Filename
	if filename == "" {
		filename = "upload"
	}
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return nil, fmt.Errorf("create form file: %w", err)
	}
	if _, err := io.Copy(part, req.File); err != nil {
		return nil, fmt.Errorf("copy file: %w", err)
	}
	if err := writeFormFieldIfSet(writer, "user_comment", req.UserComment); err != nil {
		return nil, err
	}
	if err := writeFormFieldIfSet(writer, "labels", req.Labels); err != nil {
		return nil, err
	}
	if err := writeFormFieldIfSet(writer, "metadata", req.Metadata); err != nil {
		return nil, err
	}
	if err := writeFormFieldIfSet(writer, "space_id", req.SpaceID); err != nil {
		return nil, err
	}
	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("close multipart writer: %w", err)
	}

	httpReq, err := s.client.newRequest(ctx, http.MethodPost, "/sources/ingest/file", nil, &buf, writer.FormDataContentType())
	if err != nil {
		return nil, err
	}
	var resp IngestSourceResponse
	if err := s.client.doRequest(httpReq, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// IngestURL ingests content from a URL.
//
// POST /sources/ingest/url
func (s *SourcesService) IngestURL(ctx context.Context, req *IngestURLRequest) (*IngestSourceResponse, error) {
	var resp IngestSourceResponse
	if err := s.client.do(ctx, "POST", "/sources/ingest/url", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// IngestByPath ingests a file from a server-side path.
//
// POST /sources/ingest/file-path
func (s *SourcesService) IngestByPath(ctx context.Context, req *IngestByPathRequest) (*IngestSourceResponse, error) {
	var resp IngestSourceResponse
	if err := s.client.do(ctx, "POST", "/sources/ingest/file-path", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// BatchIngest ingests a batch of files (folder import).
//
// POST /sources/ingest/batch
func (s *SourcesService) BatchIngest(ctx context.Context, req *BatchIngestRequest) (*BatchIngestResponse, error) {
	var resp BatchIngestResponse
	if err := s.client.do(ctx, "POST", "/sources/ingest/batch", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// IngestContent ingests raw text content as a library source.
//
// POST /sources/ingest/content
func (s *SourcesService) IngestContent(ctx context.Context, req *IngestContentRequest) (*IngestSourceResponse, error) {
	var resp IngestSourceResponse
	if err := s.client.do(ctx, "POST", "/sources/ingest/content", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// --- Source request types ---

// UpdateSourceRequest is the request body for Update (PATCH /sources/{id}).
type UpdateSourceRequest struct {
	Action string `json:"action"` // "reparse", "ocr_reparse", "mark_stale"
}

// IngestFileRequest is the multipart request body for IngestFile (POST /sources/ingest/file).
type IngestFileRequest struct {
	File        io.Reader `json:"-"`
	Filename    string    `json:"-"`
	UserComment string    `json:"user_comment,omitempty"`
	Labels      string    `json:"labels,omitempty"`
	Metadata    string    `json:"metadata,omitempty"`
	SpaceID     string    `json:"space_id,omitempty"`
}

// IngestURLRequest is the request body for IngestURL (POST /sources/ingest/url).
type IngestURLRequest struct {
	URL         string   `json:"url"`
	UserComment string   `json:"user_comment,omitempty"`
	Labels      []string `json:"labels,omitempty"`
	SpaceID     string   `json:"space_id,omitempty"`
}

// IngestByPathRequest is the request body for IngestByPath (POST /sources/ingest/file-path).
type IngestByPathRequest struct {
	FilePath string `json:"file_path"`
	SpaceID  string `json:"space_id,omitempty"`
}

// BatchIngestFile represents a single file entry in a BatchIngestRequest.
type BatchIngestFile struct {
	FilePath string `json:"file_path"`
}

// BatchIngestRequest is the request body for BatchIngest (POST /sources/ingest/batch).
type BatchIngestRequest struct {
	Files             []BatchIngestFile `json:"files"`
	FolderName        string            `json:"folder_name,omitempty"`
	UserComment       string            `json:"user_comment,omitempty"`
	Labels            []string          `json:"labels,omitempty"`
	SpaceID           string            `json:"space_id,omitempty"`
	EmitFeedEvent     *bool             `json:"emit_feed_event,omitempty"`
	AccumulatedTotals map[string]any    `json:"accumulated_totals,omitempty"`
}

// BatchIngestResponse is the response from BatchIngest (POST /sources/ingest/batch).
type BatchIngestResponse struct {
	FolderName      string                 `json:"folder_name,omitempty"`
	TotalIngested   int                    `json:"total_ingested"`
	TotalDuplicates int                    `json:"total_duplicates"`
	TotalErrors     int                    `json:"total_errors"`
	Results         []IngestSourceResponse `json:"results,omitempty"`
	Message         string                 `json:"message,omitempty"`
}

// IngestContentRequest is the request for IngestContent (POST /sources/ingest/content).
type IngestContentRequest struct {
	Name        string   `json:"name"`
	Content     string   `json:"content"`
	UserComment *string  `json:"user_comment,omitempty"`
	Labels      []string `json:"labels,omitempty"`
	SpaceID     string   `json:"space_id,omitempty"`
}

// UpdateContent updates the parsed markdown content of a source.
//
// PUT /sources/{id}/content
func (s *SourcesService) UpdateContent(ctx context.Context, sourceID, content string) error {
	path := fmt.Sprintf("/sources/%s/content", url.PathEscape(sourceID))
	body := map[string]string{"content": content}
	return s.client.do(ctx, "PUT", path, body, nil)
}

// Extract triggers knowledge extraction from a source.
//
// POST /sources/{id}/extract
func (s *SourcesService) Extract(ctx context.Context, sourceID string) error {
	path := fmt.Sprintf("/sources/%s/extract", url.PathEscape(sourceID))
	return s.client.do(ctx, "POST", path, nil, nil)
}

// Refetch re-fetches a URL source's content and re-parses it.
//
// POST /sources/{id}/refetch
func (s *SourcesService) Refetch(ctx context.Context, sourceID string) (*Source, error) {
	var resp Source
	path := fmt.Sprintf("/sources/%s/refetch", url.PathEscape(sourceID))
	if err := s.client.do(ctx, "POST", path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetLabels returns labels assigned to a source.
//
// GET /sources/{id}/labels
func (s *SourcesService) GetLabels(ctx context.Context, sourceID string) ([]Label, error) {
	var resp []Label
	path := fmt.Sprintf("/sources/%s/labels", url.PathEscape(sourceID))
	if err := s.client.do(ctx, "GET", path, nil, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// AssignLabel assigns a label to a source.
//
// POST /sources/{id}/labels/{label_id}
func (s *SourcesService) AssignLabel(ctx context.Context, sourceID, labelID string) error {
	path := fmt.Sprintf("/sources/%s/labels/%s", url.PathEscape(sourceID), url.PathEscape(labelID))
	return s.client.do(ctx, "POST", path, nil, nil)
}

// RemoveLabel removes a label from a source.
//
// DELETE /sources/{id}/labels/{label_id}
func (s *SourcesService) RemoveLabel(ctx context.Context, sourceID, labelID string) error {
	path := fmt.Sprintf("/sources/%s/labels/%s", url.PathEscape(sourceID), url.PathEscape(labelID))
	return s.client.do(ctx, "DELETE", path, nil, nil)
}

// IngestFolderUpload uploads a folder preserving relative paths.
//
// POST /sources/ingest/folder-upload
func (s *SourcesService) IngestFolderUpload(ctx context.Context, req *IngestFolderUploadRequest) (*FolderIngestResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("request is required")
	}
	if req.FolderName == "" {
		return nil, fmt.Errorf("folder_name is required")
	}
	if len(req.Files) == 0 {
		return nil, fmt.Errorf("at least one file is required")
	}

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	for i, file := range req.Files {
		if file.File == nil {
			return nil, fmt.Errorf("files[%d].file is required", i)
		}
		filename := file.Filename
		if filename == "" {
			filename = file.RelativePath
		}
		if filename == "" {
			filename = fmt.Sprintf("file-%d", i+1)
		}
		part, err := writer.CreateFormFile("files", filename)
		if err != nil {
			return nil, fmt.Errorf("create form file %q: %w", filename, err)
		}
		if _, err := io.Copy(part, file.File); err != nil {
			return nil, fmt.Errorf("copy file %q: %w", filename, err)
		}
	}
	if err := writer.WriteField("folder_name", req.FolderName); err != nil {
		return nil, fmt.Errorf("write folder_name field: %w", err)
	}
	if err := writeFormFieldIfSet(writer, "file_manifest", req.FileManifest); err != nil {
		return nil, err
	}
	if err := writeFormFieldIfSet(writer, "user_comment", req.UserComment); err != nil {
		return nil, err
	}
	if err := writeFormFieldIfSet(writer, "labels", req.Labels); err != nil {
		return nil, err
	}
	if err := writeFormFieldIfSet(writer, "space_id", req.SpaceID); err != nil {
		return nil, err
	}
	if req.EmitFeedEvent != nil {
		if err := writer.WriteField("emit_feed_event", strconv.FormatBool(*req.EmitFeedEvent)); err != nil {
			return nil, fmt.Errorf("write emit_feed_event field: %w", err)
		}
	}
	if err := writeFormFieldIfSet(writer, "accumulated_totals", req.AccumulatedTotals); err != nil {
		return nil, err
	}
	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("close multipart writer: %w", err)
	}

	httpReq, err := s.client.newRequest(ctx, http.MethodPost, "/sources/ingest/folder-upload", nil, &buf, writer.FormDataContentType())
	if err != nil {
		return nil, err
	}
	var resp FolderIngestResponse
	if err := s.client.doRequest(httpReq, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// IngestFolderSummary returns a summary of a folder before ingestion.
//
// POST /sources/ingest/folder-summary
func (s *SourcesService) IngestFolderSummary(ctx context.Context, req *IngestFolderSummaryRequest) (*FolderIngestResponse, error) {
	var resp FolderIngestResponse
	if err := s.client.do(ctx, "POST", "/sources/ingest/folder-summary", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// --- Source extended types ---

// IngestFolderUploadRequest is the multipart request for IngestFolderUpload (POST /sources/ingest/folder-upload).
type IngestFolderUploadRequest struct {
	Files             []FolderUploadFile `json:"-"`
	FolderName        string             `json:"folder_name"`
	FileManifest      string             `json:"file_manifest,omitempty"`
	UserComment       string             `json:"user_comment,omitempty"`
	Labels            string             `json:"labels,omitempty"`
	SpaceID           string             `json:"space_id,omitempty"`
	EmitFeedEvent     *bool              `json:"emit_feed_event,omitempty"`
	AccumulatedTotals string             `json:"accumulated_totals,omitempty"`
}

// FolderUploadFile represents a single file in a folder upload.
type FolderUploadFile struct {
	File         io.Reader `json:"-"`
	Filename     string    `json:"filename,omitempty"`
	RelativePath string    `json:"relative_path,omitempty"`
}

// IngestFolderSummaryRequest is the request for IngestFolderSummary (POST /sources/ingest/folder-summary).
type IngestFolderSummaryRequest struct {
	FolderName        string         `json:"folder_name"`
	AccumulatedTotals map[string]any `json:"accumulated_totals"`
	SpaceID           string         `json:"space_id,omitempty"`
}

// IngestSourceResponse is the result of a single source ingestion.
type IngestSourceResponse struct {
	SourceID       string `json:"source_id"`
	OriginalName   string `json:"original_name"`
	LifecycleState string `json:"lifecycle_state"`
	IsDuplicate    bool   `json:"is_duplicate"`
	Message        string `json:"message,omitempty"`
}

// FolderIngestResponse is the response from IngestFolderUpload and IngestFolderSummary.
type FolderIngestResponse struct {
	FolderName      string                 `json:"folder_name"`
	TotalIngested   int                    `json:"total_ingested"`
	TotalDuplicates int                    `json:"total_duplicates"`
	TotalErrors     int                    `json:"total_errors"`
	Results         []IngestSourceResponse `json:"results,omitempty"`
	Message         string                 `json:"message,omitempty"`
}

// GetRawFile serves the raw source file for native preview.
//
// GET /sources/{id}/raw
func (s *SourcesService) GetRawFile(ctx context.Context, sourceID string) ([]byte, error) {
	path := fmt.Sprintf("/sources/%s/raw", url.PathEscape(sourceID))
	return s.client.doBytes(ctx, http.MethodGet, path, nil, nil)
}

// GetImage serves an extracted image from a source.
//
// GET /sources/{id}/images/{filename}
func (s *SourcesService) GetImage(ctx context.Context, sourceID, filename string) ([]byte, error) {
	path := fmt.Sprintf("/sources/%s/images/%s", url.PathEscape(sourceID), url.PathEscape(filename))
	return s.client.doBytes(ctx, http.MethodGet, path, nil, nil)
}

func writeFormFieldIfSet(writer *multipart.Writer, name, value string) error {
	if value == "" {
		return nil
	}
	if err := writer.WriteField(name, value); err != nil {
		return fmt.Errorf("write %s field: %w", name, err)
	}
	return nil
}
