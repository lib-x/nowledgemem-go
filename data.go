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

// DataService handles data export/import operations.
type DataService struct {
	client *Client
}

// Export exports a portable data bundle to a server-side path.
func (s *DataService) Export(ctx context.Context, req *DataExportRequest) (*DataExportResponse, error) {
	var resp DataExportResponse
	if err := s.client.do(ctx, "POST", "/data/export", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// DownloadExport creates and downloads a ZIP export.
func (s *DataService) DownloadExport(ctx context.Context, req *DataExportDownloadRequest) ([]byte, error) {
	return s.client.doBytes(ctx, http.MethodPost, "/data/export/download", nil, req)
}

// Import imports data from a server-side export path.
func (s *DataService) Import(ctx context.Context, req *DataImportRequest) (*DataImportResponse, error) {
	var resp DataImportResponse
	if err := s.client.do(ctx, "POST", "/data/import", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ImportStatus checks status of a data import job.
func (s *DataService) ImportStatus(ctx context.Context, jobID string) (*DataImportStatus, error) {
	var resp DataImportStatus
	path := fmt.Sprintf("/data/import/status/%s", url.PathEscape(jobID))
	if err := s.client.do(ctx, "GET", path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Checkpoint forces a database checkpoint.
func (s *DataService) Checkpoint(ctx context.Context) error {
	return s.client.do(ctx, "POST", "/data/checkpoint", nil, nil)
}

// --- Data types ---

// DataExportRequest is the request for POST /data/export.
type DataExportRequest struct {
	ExportPath                   string `json:"export_path"`
	Compress                     bool   `json:"compress,omitempty"`
	Overwrite                    bool   `json:"overwrite,omitempty"`
	IncludeMemories              *bool  `json:"include_memories,omitempty"`
	IncludeThreads               *bool  `json:"include_threads,omitempty"`
	IncludeMessages              *bool  `json:"include_messages,omitempty"`
	IncludeEntities              *bool  `json:"include_entities,omitempty"`
	IncludeLabels                *bool  `json:"include_labels,omitempty"`
	IncludeSources               *bool  `json:"include_sources,omitempty"`
	IncludeCommunities           *bool  `json:"include_communities,omitempty"`
	IncludeSkills                *bool  `json:"include_skills,omitempty"`
	IncludeEdges                 *bool  `json:"include_edges,omitempty"`
	IncludeWorkingMemory         *bool  `json:"include_working_memory,omitempty"`
	IncludeWorkingMemoryArchive  *bool  `json:"include_working_memory_archive,omitempty"`
	IncludeSourceFiles           *bool  `json:"include_source_files,omitempty"`
}

// DataExportDownloadRequest is the request for POST /data/export/download.
// This endpoint streams a ZIP directly to the client and does not require
// a server-side path.
type DataExportDownloadRequest struct {
	IncludeMemories              *bool `json:"include_memories,omitempty"`
	IncludeThreads               *bool `json:"include_threads,omitempty"`
	IncludeMessages              *bool `json:"include_messages,omitempty"`
	IncludeEntities              *bool `json:"include_entities,omitempty"`
	IncludeLabels                *bool `json:"include_labels,omitempty"`
	IncludeSources               *bool `json:"include_sources,omitempty"`
	IncludeCommunities           *bool `json:"include_communities,omitempty"`
	IncludeSkills                *bool `json:"include_skills,omitempty"`
	IncludeEdges                 *bool `json:"include_edges,omitempty"`
	IncludeWorkingMemory         *bool `json:"include_working_memory,omitempty"`
	IncludeWorkingMemoryArchive  *bool `json:"include_working_memory_archive,omitempty"`
	IncludeSourceFiles           *bool `json:"include_source_files,omitempty"`
}

// DataExportResponse is the response for POST /data/export.
type DataExportResponse struct {
	Path      string `json:"path"`
	SizeBytes int64  `json:"size_bytes"`
	ItemCount int    `json:"item_count"`
}

// DataImportRequest is the request for POST /data/import.
type DataImportRequest struct {
	ImportPath                   string `json:"import_path"`
	Mode                         string `json:"mode,omitempty"`
	IncludeMemories              *bool  `json:"include_memories,omitempty"`
	IncludeThreads               *bool  `json:"include_threads,omitempty"`
	IncludeMessages              *bool  `json:"include_messages,omitempty"`
	IncludeEntities              *bool  `json:"include_entities,omitempty"`
	IncludeLabels                *bool  `json:"include_labels,omitempty"`
	IncludeSources               *bool  `json:"include_sources,omitempty"`
	IncludeCommunities           *bool  `json:"include_communities,omitempty"`
	IncludeSkills                *bool  `json:"include_skills,omitempty"`
	IncludeEdges                 *bool  `json:"include_edges,omitempty"`
	IncludeWorkingMemory         *bool  `json:"include_working_memory,omitempty"`
	IncludeWorkingMemoryArchive  *bool  `json:"include_working_memory_archive,omitempty"`
	IncludeSourceFiles           *bool  `json:"include_source_files,omitempty"`
}

// DataImportResponse is the response for POST /data/import.
type DataImportResponse struct {
	JobID   string `json:"job_id"`
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

// DataImportStatus is the response for GET /data/import/status/{job_id}.
type DataImportStatus struct {
	JobID    string  `json:"job_id"`
	Status   string  `json:"status"`
	Progress float64 `json:"progress"`
	Imported int     `json:"imported"`
	Skipped  int     `json:"skipped"`
	Failed   int     `json:"failed"`
	Message  string  `json:"message,omitempty"`
}

// UploadImportRequest is the request for POST /data/import/upload.
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

// UploadImport uploads a ZIP export from web and remote clients.
// file is the ZIP file content, filename is the name of the file.
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
