package nowledgemem

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

// LabelsService handles label operations.
//
// It provides methods for creating, listing, updating, and deleting labels.
type LabelsService struct {
	client *Client
}

// List returns all labels with usage counts.
//
// GET /labels
func (s *LabelsService) List(ctx context.Context, params *ListLabelsParams) ([]Label, error) {
	q := url.Values{}
	if params != nil {
		if params.Limit > 0 {
			q.Set("limit", strconv.Itoa(params.Limit))
		}
		if params.OrderBy != "" {
			q.Set("order_by", params.OrderBy)
		}
		if params.OrderDesc {
			q.Set("order_desc", "true")
		}
	}
	var resp []Label
	if err := s.client.doQuery(ctx, "/labels", q, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// Create creates a new label.
//
// POST /labels
func (s *LabelsService) Create(ctx context.Context, req *CreateLabelRequest) (*Label, error) {
	q := url.Values{}
	if req != nil {
		if req.Name != "" {
			q.Set("name", req.Name)
		}
		if req.Color != "" {
			q.Set("color", req.Color)
		}
		if req.Description != "" {
			q.Set("description", req.Description)
		}
	}
	var resp Label
	if err := s.client.doWithQuery(ctx, "POST", "/labels", q, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Get returns a specific label by ID.
//
// GET /labels/{id}
func (s *LabelsService) Get(ctx context.Context, labelID string) (*Label, error) {
	var resp Label
	path := fmt.Sprintf("/labels/%s", url.PathEscape(labelID))
	if err := s.client.do(ctx, "GET", path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Update updates an existing label.
//
// PUT /labels/{id}
func (s *LabelsService) Update(ctx context.Context, labelID string, req *UpdateLabelRequest) (*Label, error) {
	var resp Label
	path := fmt.Sprintf("/labels/%s", url.PathEscape(labelID))
	if err := s.client.do(ctx, "PUT", path, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Delete deletes a label and all its relationships.
//
// DELETE /labels/{id}
func (s *LabelsService) Delete(ctx context.Context, labelID string) error {
	path := fmt.Sprintf("/labels/%s", url.PathEscape(labelID))
	return s.client.do(ctx, "DELETE", path, nil, nil)
}

// Health returns label system health and maintenance status.
//
// GET /labels/health
func (s *LabelsService) Health(ctx context.Context) (map[string]any, error) {
	var resp map[string]any
	if err := s.client.do(ctx, "GET", "/labels/health", nil, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// MergeCandidates returns likely duplicate or overlapping labels.
//
// GET /labels/merge-candidates
func (s *LabelsService) MergeCandidates(ctx context.Context, params *LabelMergeCandidatesParams) (map[string]any, error) {
	q := url.Values{}
	if params != nil {
		if params.SimFloor > 0 {
			q.Set("sim_floor", strconv.FormatFloat(params.SimFloor, 'f', -1, 64))
		}
		if params.MaxLabels > 0 {
			q.Set("max_labels", strconv.Itoa(params.MaxLabels))
		}
		if params.MaxPairs > 0 {
			q.Set("max_pairs", strconv.Itoa(params.MaxPairs))
		}
	}
	var resp map[string]any
	if err := s.client.doQuery(ctx, "/labels/merge-candidates", q, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// Merge merges one label into another.
//
// POST /labels/merge
func (s *LabelsService) Merge(ctx context.Context, sourceID, targetID string) (map[string]any, error) {
	q := url.Values{
		"source_id": {sourceID},
		"target_id": {targetID},
	}
	var resp map[string]any
	if err := s.client.doWithQuery(ctx, "POST", "/labels/merge", q, nil, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// ConsolidationPreview previews label consolidation candidates.
//
// POST /labels/consolidation-preview
func (s *LabelsService) ConsolidationPreview(ctx context.Context, params *LabelConsolidationPreviewParams) (map[string]any, error) {
	q := url.Values{}
	if params != nil {
		if params.SimFloor > 0 {
			q.Set("sim_floor", strconv.FormatFloat(params.SimFloor, 'f', -1, 64))
		}
		if params.MaxPairs > 0 {
			q.Set("max_pairs", strconv.Itoa(params.MaxPairs))
		}
	}
	var resp map[string]any
	if err := s.client.doWithQuery(ctx, "POST", "/labels/consolidation-preview", q, nil, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// Consolidate runs label consolidation.
//
// POST /labels/consolidate
func (s *LabelsService) Consolidate(ctx context.Context, params *LabelConsolidateParams) (map[string]any, error) {
	q := url.Values{}
	if params != nil {
		if params.DryRun != nil {
			q.Set("dry_run", strconv.FormatBool(*params.DryRun))
		}
		if params.MaxGroups > 0 {
			q.Set("max_groups", strconv.Itoa(params.MaxGroups))
		}
	}
	var resp map[string]any
	if err := s.client.doWithQuery(ctx, "POST", "/labels/consolidate", q, nil, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// --- Label request types ---

// ListLabelsParams are query parameters for List (GET /labels).
type ListLabelsParams struct {
	Limit     int    `json:"limit,omitempty"`
	OrderBy   string `json:"order_by,omitempty"`
	OrderDesc bool   `json:"order_desc,omitempty"`
}

// CreateLabelRequest is the request body for Create (POST /labels).
type CreateLabelRequest struct {
	Name        string `json:"name"`
	Color       string `json:"color,omitempty"`
	Description string `json:"description,omitempty"`
}

// UpdateLabelRequest is the request body for Update (PUT /labels/{id}).
type UpdateLabelRequest struct {
	Name        string `json:"name,omitempty"`
	Color       string `json:"color,omitempty"`
	Description string `json:"description,omitempty"`
}

// LabelMergeCandidatesParams are query parameters for MergeCandidates.
type LabelMergeCandidatesParams struct {
	SimFloor  float64 `json:"sim_floor,omitempty"`
	MaxLabels int     `json:"max_labels,omitempty"`
	MaxPairs  int     `json:"max_pairs,omitempty"`
}

// LabelConsolidationPreviewParams are query parameters for ConsolidationPreview.
type LabelConsolidationPreviewParams struct {
	SimFloor float64 `json:"sim_floor,omitempty"`
	MaxPairs int     `json:"max_pairs,omitempty"`
}

// LabelConsolidateParams are query parameters for Consolidate.
type LabelConsolidateParams struct {
	DryRun    *bool `json:"dry_run,omitempty"`
	MaxGroups int   `json:"max_groups,omitempty"`
}
