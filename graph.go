package onledgemem

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

// GraphService handles knowledge graph operations.
type GraphService struct {
	client *Client
}

// Analysis returns comprehensive graph analysis including community and centrality metrics.
func (s *GraphService) Analysis(ctx context.Context) (*GraphAnalysis, error) {
	var resp GraphAnalysis
	if err := s.client.do(ctx, "GET", "/graph/analysis", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Health checks graph analysis service and algo extensions.
func (s *GraphService) Health(ctx context.Context) (*GraphHealthResponse, error) {
	var resp GraphHealthResponse
	if err := s.client.do(ctx, "GET", "/graph/health", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// FindOrphans finds entities with no relationships.
func (s *GraphService) FindOrphans(ctx context.Context) ([]Entity, error) {
	var resp []Entity
	if err := s.client.do(ctx, "GET", "/graph/orphans", nil, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// CleanupOrphans removes orphaned entities from the graph.
func (s *GraphService) CleanupOrphans(ctx context.Context) (*CleanupOrphansResponse, error) {
	var resp CleanupOrphansResponse
	if err := s.client.do(ctx, "DELETE", "/graph/orphans", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// StartAugmentation starts a background job (community detection, PageRank).
func (s *GraphService) StartAugmentation(ctx context.Context, req *StartAugmentationRequest) (*AugmentationJob, error) {
	var resp AugmentationJob
	if err := s.client.do(ctx, "POST", "/graph/augmentation/start", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// AugmentationState returns current augmentation status and parameters.
func (s *GraphService) AugmentationState(ctx context.Context) (*AugmentationState, error) {
	var resp AugmentationState
	if err := s.client.do(ctx, "GET", "/graph/augmentation/state", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// JobStatus checks progress of a specific augmentation job.
func (s *GraphService) JobStatus(ctx context.Context, jobID string) (*AugmentationJob, error) {
	var resp AugmentationJob
	path := fmt.Sprintf("/graph/augmentation/status/%s", url.PathEscape(jobID))
	if err := s.client.do(ctx, "GET", path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ListJobs lists recent augmentation jobs.
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

// GraphHealthResponse is the response body for GET /graph/health.
type GraphHealthResponse struct {
	Status    string `json:"status"`
	AlgoReady bool   `json:"algo_ready"`
}

// CleanupOrphansResponse is the response body for DELETE /graph/orphans.
type CleanupOrphansResponse struct {
	Removed int `json:"removed"`
}

// StartAugmentationRequest is the request body for POST /graph/augmentation/start.
type StartAugmentationRequest struct {
	JobType string `json:"job_type"`
	Params  map[string]any `json:"params,omitempty"`
}

// AugmentationJob represents an augmentation job.
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

// AugmentationState is the response body for GET /graph/augmentation/state.
type AugmentationState struct {
	Running    bool             `json:"running"`
	CurrentJob *AugmentationJob `json:"current_job,omitempty"`
	LastRun    string           `json:"last_run,omitempty"`
}

// GetNodeDetails returns detailed information about a specific node.
func (s *GraphService) GetNodeDetails(ctx context.Context, nodeID string) (*GraphNode, error) {
	var resp GraphNode
	path := fmt.Sprintf("/graph/node-details/%s", url.PathEscape(nodeID))
	if err := s.client.do(ctx, "GET", path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetCommunityMembers returns members of a specific community.
func (s *GraphService) GetCommunityMembers(ctx context.Context, communityID string, limit int) (*GraphData, error) {
	q := url.Values{}
	if limit > 0 {
		q.Set("limit", strconv.Itoa(limit))
	}
	var resp GraphData
	path := fmt.Sprintf("/graph/community-members/%s", url.PathEscape(communityID))
	if err := s.client.doQuery(ctx, path, q, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
