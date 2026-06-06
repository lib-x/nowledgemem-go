package onledgemem

import (
	"context"
	"fmt"
	"net/url"
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
func (s *DataService) DownloadExport(ctx context.Context, req *DataExportRequest) ([]byte, error) {
	var resp struct{}
	if err := s.client.do(ctx, "POST", "/data/export/download", req, &resp); err != nil {
		return nil, err
	}
	return nil, nil // TODO: handle binary download
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
	ExportPath        string `json:"export_path"`
	Compress          bool   `json:"compress,omitempty"`
	Overwrite         bool   `json:"overwrite,omitempty"`
	IncludeMemories   *bool  `json:"include_memories,omitempty"`
	IncludeThreads    *bool  `json:"include_threads,omitempty"`
	IncludeMessages   *bool  `json:"include_messages,omitempty"`
	IncludeEntities   *bool  `json:"include_entities,omitempty"`
	IncludeLabels     *bool  `json:"include_labels,omitempty"`
	IncludeSources    *bool  `json:"include_sources,omitempty"`
	IncludeCommunities *bool `json:"include_communities,omitempty"`
}

// DataExportResponse is the response for POST /data/export.
type DataExportResponse struct {
	Path       string `json:"path"`
	SizeBytes  int64  `json:"size_bytes"`
	ItemCount  int    `json:"item_count"`
}

// DataImportRequest is the request for POST /data/import.
type DataImportRequest struct {
	ImportPath string `json:"import_path"`
	Overwrite  bool   `json:"overwrite,omitempty"`
}

// DataImportResponse is the response for POST /data/import.
type DataImportResponse struct {
	JobID   string `json:"job_id"`
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

// DataImportStatus is the response for GET /data/import/status/{job_id}.
type DataImportStatus struct {
	JobID       string  `json:"job_id"`
	Status      string  `json:"status"`
	Progress    float64 `json:"progress"`
	Imported    int     `json:"imported"`
	Skipped     int     `json:"skipped"`
	Failed      int     `json:"failed"`
	Message     string  `json:"message,omitempty"`
}
