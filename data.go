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

// DataService handles data export and import operations.
//
// It provides methods for exporting data to archives, importing from previous exports,
// and forcing database checkpoints.
type DataService struct {
	client *Client
}

// Export exports data to a portable archive at a server-side path.
//
// POST /data/export
func (s *DataService) Export(ctx context.Context, req *DataExportRequest) (*DataExportResponse, error) {
	var resp DataExportResponse
	if err := s.client.do(ctx, "POST", "/data/export", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// DownloadExport exports data as a downloadable ZIP file.
//
// POST /data/export/download
func (s *DataService) DownloadExport(ctx context.Context, req *DataExportDownloadRequest) ([]byte, error) {
	return s.client.doBytes(ctx, http.MethodPost, "/data/export/download", nil, req)
}

// ExportStatus checks status of a background data export job.
//
// GET /data/export/status/{id}
func (s *DataService) ExportStatus(ctx context.Context, jobID string) (*DataTransferStatus, error) {
	var resp DataTransferStatus
	path := fmt.Sprintf("/data/export/status/%s", url.PathEscape(jobID))
	if err := s.client.do(ctx, "GET", path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Import imports data from a previous export at a server-side path.
//
// POST /data/import
func (s *DataService) Import(ctx context.Context, req *DataImportRequest) (*DataImportResponse, error) {
	var resp DataImportResponse
	if err := s.client.do(ctx, "POST", "/data/import", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ImportStatus checks status of a data import job.
//
// GET /data/import/status/{id}
func (s *DataService) ImportStatus(ctx context.Context, jobID string) (*DataImportStatus, error) {
	var resp DataImportStatus
	path := fmt.Sprintf("/data/import/status/%s", url.PathEscape(jobID))
	if err := s.client.do(ctx, "GET", path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Checkpoint forces a database checkpoint.
//
// POST /data/checkpoint
func (s *DataService) Checkpoint(ctx context.Context) error {
	return s.client.do(ctx, "POST", "/data/checkpoint", nil, nil)
}

// CheckpointResult forces a database checkpoint and returns the API response.
//
// POST /data/checkpoint
func (s *DataService) CheckpointResult(ctx context.Context) (*DataTransferCheckpointResponse, error) {
	var resp DataTransferCheckpointResponse
	if err := s.client.do(ctx, "POST", "/data/checkpoint", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// --- Data types ---

// DataExportRequest is the request for Export (POST /data/export).
type DataExportRequest struct {
	ExportPath                  string `json:"export_path"`
	Async                       bool   `json:"async,omitempty"`
	Compress                    bool   `json:"compress,omitempty"`
	Overwrite                   bool   `json:"overwrite,omitempty"`
	IncludeMemories             *bool  `json:"include_memories,omitempty"`
	IncludeThreads              *bool  `json:"include_threads,omitempty"`
	IncludeMessages             *bool  `json:"include_messages,omitempty"`
	IncludeEntities             *bool  `json:"include_entities,omitempty"`
	IncludeLabels               *bool  `json:"include_labels,omitempty"`
	IncludeSources              *bool  `json:"include_sources,omitempty"`
	IncludeCommunities          *bool  `json:"include_communities,omitempty"`
	IncludeSkills               *bool  `json:"include_skills,omitempty"`
	IncludeEdges                *bool  `json:"include_edges,omitempty"`
	IncludeWorkingMemory        *bool  `json:"include_working_memory,omitempty"`
	IncludeWorkingMemoryArchive *bool  `json:"include_working_memory_archive,omitempty"`
	IncludeSourceFiles          *bool  `json:"include_source_files,omitempty"`
}

// DataExportDownloadRequest is the request for DownloadExport (POST /data/export/download).
//
// This endpoint streams a ZIP directly to the client and does not require a server-side path.
type DataExportDownloadRequest struct {
	IncludeMemories             *bool `json:"include_memories,omitempty"`
	IncludeThreads              *bool `json:"include_threads,omitempty"`
	IncludeMessages             *bool `json:"include_messages,omitempty"`
	IncludeEntities             *bool `json:"include_entities,omitempty"`
	IncludeLabels               *bool `json:"include_labels,omitempty"`
	IncludeSources              *bool `json:"include_sources,omitempty"`
	IncludeCommunities          *bool `json:"include_communities,omitempty"`
	IncludeSkills               *bool `json:"include_skills,omitempty"`
	IncludeEdges                *bool `json:"include_edges,omitempty"`
	IncludeWorkingMemory        *bool `json:"include_working_memory,omitempty"`
	IncludeWorkingMemoryArchive *bool `json:"include_working_memory_archive,omitempty"`
	IncludeSourceFiles          *bool `json:"include_source_files,omitempty"`
}

// DataExportResponse is the response from Export (POST /data/export).
//
// The API returns sync export result fields when async is false and job start
// fields when async is true. Legacy path/size fields are kept for older servers.
type DataExportResponse struct {
	Path      string `json:"path,omitempty"`
	SizeBytes int64  `json:"size_bytes,omitempty"`
	ItemCount int    `json:"item_count,omitempty"`

	Success      bool           `json:"success,omitempty"`
	ExportID     string         `json:"export_id,omitempty"`
	Format       string         `json:"format,omitempty"`
	Version      int64          `json:"version,omitempty"`
	CreatedAt    string         `json:"created_at,omitempty"`
	Counts       map[string]any `json:"counts,omitempty"`
	Warnings     []string       `json:"warnings,omitempty"`
	ExportDir    *string        `json:"export_dir,omitempty"`
	ExportPath   *string        `json:"export_path,omitempty"`
	ManifestPath *string        `json:"manifest_path,omitempty"`

	JobID   string `json:"job_id,omitempty"`
	Status  string `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
}

// DataImportRequest is the request for Import (POST /data/import).
type DataImportRequest struct {
	ImportPath                  string `json:"import_path"`
	Mode                        string `json:"mode,omitempty"`
	IncludeMemories             *bool  `json:"include_memories,omitempty"`
	IncludeThreads              *bool  `json:"include_threads,omitempty"`
	IncludeMessages             *bool  `json:"include_messages,omitempty"`
	IncludeEntities             *bool  `json:"include_entities,omitempty"`
	IncludeLabels               *bool  `json:"include_labels,omitempty"`
	IncludeSources              *bool  `json:"include_sources,omitempty"`
	IncludeCommunities          *bool  `json:"include_communities,omitempty"`
	IncludeSkills               *bool  `json:"include_skills,omitempty"`
	IncludeEdges                *bool  `json:"include_edges,omitempty"`
	IncludeWorkingMemory        *bool  `json:"include_working_memory,omitempty"`
	IncludeWorkingMemoryArchive *bool  `json:"include_working_memory_archive,omitempty"`
	IncludeSourceFiles          *bool  `json:"include_source_files,omitempty"`
}

// DataImportResponse is the response from Import and UploadImport.
type DataImportResponse struct {
	JobID   string `json:"job_id"`
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
	Success bool   `json:"success,omitempty"`
}

// DataTransferStatus is the response from import/export status endpoints.
type DataTransferStatus struct {
	JobID       string         `json:"job_id"`
	Status      string         `json:"status"`
	Kind        *string        `json:"kind,omitempty"`
	Result      map[string]any `json:"result,omitempty"`
	StartedAt   *string        `json:"started_at,omitempty"`
	CompletedAt *string        `json:"completed_at,omitempty"`
	Error       *string        `json:"error,omitempty"`

	// Older servers returned progress counters at the top level.
	Progress float64 `json:"progress,omitempty"`
	Imported int     `json:"imported,omitempty"`
	Skipped  int     `json:"skipped,omitempty"`
	Failed   int     `json:"failed,omitempty"`
	Message  string  `json:"message,omitempty"`
}

// DataImportStatus is kept as an alias for ImportStatus callers.
type DataImportStatus = DataTransferStatus

// DataTransferCheckpointResponse is the response from CheckpointResult.
type DataTransferCheckpointResponse struct {
	Success      bool    `json:"success"`
	Checkpointed bool    `json:"checkpointed"`
	Detail       *string `json:"detail,omitempty"`
}

// UploadImportRequest is the multipart request for UploadImport (POST /data/import/upload).
type UploadImportRequest struct {
	File                        io.Reader `json:"-"`
	Filename                    string    `json:"-"`
	Mode                        string    `json:"mode,omitempty"`
	IncludeMemories             *bool     `json:"include_memories,omitempty"`
	IncludeThreads              *bool     `json:"include_threads,omitempty"`
	IncludeMessages             *bool     `json:"include_messages,omitempty"`
	IncludeEntities             *bool     `json:"include_entities,omitempty"`
	IncludeLabels               *bool     `json:"include_labels,omitempty"`
	IncludeSources              *bool     `json:"include_sources,omitempty"`
	IncludeCommunities          *bool     `json:"include_communities,omitempty"`
	IncludeSkills               *bool     `json:"include_skills,omitempty"`
	IncludeEdges                *bool     `json:"include_edges,omitempty"`
	IncludeWorkingMemory        *bool     `json:"include_working_memory,omitempty"`
	IncludeWorkingMemoryArchive *bool     `json:"include_working_memory_archive,omitempty"`
	IncludeSourceFiles          *bool     `json:"include_source_files,omitempty"`
}

// UploadImport imports data from an uploaded ZIP file.
//
// POST /data/import/upload
func (s *DataService) UploadImport(ctx context.Context, req *UploadImportRequest) (*DataImportResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("request is required")
	}
	if req.File == nil {
		return nil, fmt.Errorf("file is required")
	}

	// Build multipart form
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// Add file field
	filename := req.Filename
	if filename == "" {
		filename = "import.zip"
	}
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return nil, fmt.Errorf("create form file: %w", err)
	}
	if _, err := io.Copy(part, req.File); err != nil {
		return nil, fmt.Errorf("copy file: %w", err)
	}

	// Add optional fields
	if req.Mode != "" {
		writer.WriteField("mode", req.Mode)
	}
	if req.IncludeMemories != nil {
		writer.WriteField("include_memories", strconv.FormatBool(*req.IncludeMemories))
	}
	if req.IncludeThreads != nil {
		writer.WriteField("include_threads", strconv.FormatBool(*req.IncludeThreads))
	}
	if req.IncludeMessages != nil {
		writer.WriteField("include_messages", strconv.FormatBool(*req.IncludeMessages))
	}
	if req.IncludeEntities != nil {
		writer.WriteField("include_entities", strconv.FormatBool(*req.IncludeEntities))
	}
	if req.IncludeLabels != nil {
		writer.WriteField("include_labels", strconv.FormatBool(*req.IncludeLabels))
	}
	if req.IncludeSources != nil {
		writer.WriteField("include_sources", strconv.FormatBool(*req.IncludeSources))
	}
	if req.IncludeCommunities != nil {
		writer.WriteField("include_communities", strconv.FormatBool(*req.IncludeCommunities))
	}
	if req.IncludeEdges != nil {
		writer.WriteField("include_edges", strconv.FormatBool(*req.IncludeEdges))
	}
	if req.IncludeWorkingMemory != nil {
		writer.WriteField("include_working_memory", strconv.FormatBool(*req.IncludeWorkingMemory))
	}
	if req.IncludeSkills != nil {
		writer.WriteField("include_skills", strconv.FormatBool(*req.IncludeSkills))
	}
	if req.IncludeWorkingMemoryArchive != nil {
		writer.WriteField("include_working_memory_archive", strconv.FormatBool(*req.IncludeWorkingMemoryArchive))
	}
	if req.IncludeSourceFiles != nil {
		writer.WriteField("include_source_files", strconv.FormatBool(*req.IncludeSourceFiles))
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("close multipart writer: %w", err)
	}

	httpReq, err := s.client.newRequest(ctx, http.MethodPost, "/data/import/upload", nil, &buf, writer.FormDataContentType())
	if err != nil {
		return nil, err
	}

	var result DataImportResponse
	if err := s.client.doRequest(httpReq, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
