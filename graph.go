package nowledgemem

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

// GraphService handles knowledge graph operations.
//
// It provides methods for graph analysis, health checks, augmentation jobs,
// and orphan management.
type GraphService struct {
	client *Client
}

// Analysis returns comprehensive graph analysis including community and centrality metrics.
//
// GET /graph/analysis
func (s *GraphService) Analysis(ctx context.Context) (*GraphAnalysis, error) {
	var resp GraphAnalysis
	if err := s.client.do(ctx, "GET", "/graph/analysis", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Health returns graph database health status.
//
// GET /graph/health
func (s *GraphService) Health(ctx context.Context) (*GraphHealthResponse, error) {
	var resp GraphHealthResponse
	if err := s.client.do(ctx, "GET", "/graph/health", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// FindOrphans finds entities with no relationships.
//
// GET /graph/orphans
func (s *GraphService) FindOrphans(ctx context.Context) ([]Entity, error) {
	var resp []Entity
	if err := s.client.do(ctx, "GET", "/graph/orphans", nil, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// CleanupOrphans removes orphaned entities from the graph.
//
// DELETE /graph/orphans
func (s *GraphService) CleanupOrphans(ctx context.Context) (*CleanupOrphansResponse, error) {
	var resp CleanupOrphansResponse
	if err := s.client.do(ctx, "DELETE", "/graph/orphans", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// StartAugmentation starts a graph augmentation job (community detection, PageRank).
//
// POST /graph/augmentation/start
func (s *GraphService) StartAugmentation(ctx context.Context, req *StartAugmentationRequest) (*AugmentationJob, error) {
	var resp AugmentationJob
	if err := s.client.do(ctx, "POST", "/graph/augmentation/start", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// AugmentationState returns current augmentation status and parameters.
//
// GET /graph/augmentation/state
func (s *GraphService) AugmentationState(ctx context.Context) (*AugmentationState, error) {
	var resp AugmentationState
	if err := s.client.do(ctx, "GET", "/graph/augmentation/state", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// JobStatus checks progress of a specific augmentation job.
//
// GET /graph/augmentation/status/{id}
func (s *GraphService) JobStatus(ctx context.Context, jobID string) (*AugmentationJob, error) {
	var resp AugmentationJob
	path := fmt.Sprintf("/graph/augmentation/status/%s", url.PathEscape(jobID))
	if err := s.client.do(ctx, "GET", path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ListJobs lists recent augmentation jobs.
//
// GET /graph/augmentation/jobs
func (s *GraphService) ListJobs(ctx context.Context, limit int) ([]AugmentationJob, error) {
	q := url.Values{}
	if limit > 0 {
		q.Set("limit", strconv.Itoa(limit))
	}
	var resp []AugmentationJob
	if err := s.client.doQuery(ctx, "/graph/augmentation/jobs", q, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// --- Graph types ---

// GraphCapabilities holds graph analysis feature flags.
type GraphCapabilities struct {
	CommunityDetection   bool `json:"community_detection"`
	PagerankCalculation  bool `json:"pagerank_calculation"`
	UnifiedGraphAnalysis bool `json:"unified_graph_analysis"`
	LLMSummarization     bool `json:"llm_summarization"`
}

// GraphHealthResponse is the response from Health (GET /graph/health).
type GraphHealthResponse struct {
	Status              string            `json:"status"`
	Error               string            `json:"error,omitempty"`
	AlgoExtensionLoaded bool              `json:"algo_extension_loaded"`
	DbConnection        string            `json:"db_connection,omitempty"`
	Capabilities        GraphCapabilities `json:"capabilities"`
	CheckedAt           string            `json:"checked_at,omitempty"`
}

// CleanupOrphansResponse is the response from CleanupOrphans (DELETE /graph/orphans).
type CleanupOrphansResponse struct {
	Removed int `json:"removed"`
}

// StartAugmentationRequest is the request body for StartAugmentation (POST /graph/augmentation/start).
type StartAugmentationRequest struct {
	JobType    string         `json:"job_type"`
	Parameters map[string]any `json:"parameters,omitempty"`
}

// AugmentationJob represents a graph augmentation job.
type AugmentationJob struct {
	ID        string         `json:"id"`
	JobType   string         `json:"job_type"`
	Status    string         `json:"status"`
	Progress  float64        `json:"progress"`
	StartedAt string         `json:"started_at,omitempty"`
	EndedAt   string         `json:"ended_at,omitempty"`
	Error     string         `json:"error,omitempty"`
	Params    map[string]any `json:"params,omitempty"`
}

// AugmentationState is the response from AugmentationState (GET /graph/augmentation/state).
type AugmentationState struct {
	Running    bool             `json:"running"`
	CurrentJob *AugmentationJob `json:"current_job,omitempty"`
	LastRun    string           `json:"last_run,omitempty"`
}
