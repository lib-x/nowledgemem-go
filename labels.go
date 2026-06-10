package nowledgemem

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

// LabelsService handles label operations.
type LabelsService struct {
	client *Client
}

// List returns all labels with usage counts.
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
func (s *LabelsService) Get(ctx context.Context, labelID string) (*Label, error) {
	var resp Label
	path := fmt.Sprintf("/labels/%s", url.PathEscape(labelID))
	if err := s.client.do(ctx, "GET", path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Update updates an existing label.
func (s *LabelsService) Update(ctx context.Context, labelID string, req *UpdateLabelRequest) (*Label, error) {
	var resp Label
	path := fmt.Sprintf("/labels/%s", url.PathEscape(labelID))
	if err := s.client.do(ctx, "PUT", path, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Delete deletes a label and all its relationships.
func (s *LabelsService) Delete(ctx context.Context, labelID string) error {
	path := fmt.Sprintf("/labels/%s", url.PathEscape(labelID))
	return s.client.do(ctx, "DELETE", path, nil, nil)
}

// --- Label request types ---

// ListLabelsParams are query parameters for GET /labels.
type ListLabelsParams struct {
	Limit     int    `json:"limit,omitempty"`
	OrderBy   string `json:"order_by,omitempty"`
	OrderDesc bool   `json:"order_desc,omitempty"`
}

// CreateLabelRequest is the request body for POST /labels.
type CreateLabelRequest struct {
	Name        string `json:"name"`
	Color       string `json:"color,omitempty"`
	Description string `json:"description,omitempty"`
}

// UpdateLabelRequest is the request body for PUT /labels/{id}.
type UpdateLabelRequest struct {
	Name        string `json:"name,omitempty"`
	Color       string `json:"color,omitempty"`
	Description string `json:"description,omitempty"`
}
