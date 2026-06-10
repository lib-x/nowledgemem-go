package nowledgemem

import (
	"context"
	"fmt"
	"net/url"
)

// SpacesService handles space profile operations.
//
// It provides methods for managing spaces, their retrieval defaults, and guidance.
type SpacesService struct {
	client *Client
}

// List returns the shared space roster, profile metadata, and usage.
//
// GET /spaces
func (s *SpacesService) List(ctx context.Context) (*ListSpacesResponse, error) {
	var resp ListSpacesResponse
	if err := s.client.do(ctx, "GET", "/spaces", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Roster returns the shared space roster.
//
// GET /spaces/roster
func (s *SpacesService) Roster(ctx context.Context) (*ListSpacesResponse, error) {
	var resp ListSpacesResponse
	if err := s.client.do(ctx, "GET", "/spaces/roster", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Create creates a space profile with retrieval defaults and guidance.
//
// POST /spaces
func (s *SpacesService) Create(ctx context.Context, req *CreateSpaceRequest) (*Space, error) {
	var resp Space
	if err := s.client.do(ctx, "POST", "/spaces", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Get reads one space profile by name, alias, or hidden key.
//
// GET /spaces/{id}
func (s *SpacesService) Get(ctx context.Context, spaceID string) (*Space, error) {
	var resp Space
	path := fmt.Sprintf("/spaces/%s", url.PathEscape(spaceID))
	if err := s.client.do(ctx, "GET", path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Update renames a space, changes retrieval defaults, or updates guidance.
//
// PATCH /spaces/{id}
func (s *SpacesService) Update(ctx context.Context, spaceID string, req *UpdateSpaceRequest) (*Space, error) {
	var resp Space
	path := fmt.Sprintf("/spaces/%s", url.PathEscape(spaceID))
	if err := s.client.do(ctx, "PATCH", path, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Delete removes an empty space profile.
//
// DELETE /spaces/{id}
func (s *SpacesService) Delete(ctx context.Context, spaceID string, params *DeleteSpaceParams) (*ListSpacesResponse, error) {
	path := fmt.Sprintf("/spaces/%s", url.PathEscape(spaceID))
	var q url.Values
	if params != nil {
		q = url.Values{}
		if params.PurgeWorkingMemory {
			q.Set("purge_working_memory", "true")
		}
	}
	var resp ListSpacesResponse
	if err := s.client.doWithQuery(ctx, "DELETE", path, q, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// UpdateConfig enables or disables spaces at the product level.
//
// POST /spaces/config
func (s *SpacesService) UpdateConfig(ctx context.Context, req *SpacesConfigRequest) error {
	return s.client.do(ctx, "POST", "/spaces/config", req, nil)
}

// --- Space request types ---

// DeleteSpaceParams are query parameters for Delete (DELETE /spaces/{id}).
type DeleteSpaceParams struct {
	PurgeWorkingMemory bool `json:"purge_working_memory,omitempty"`
}

// SpacesConfigRequest is the request for UpdateConfig (POST /spaces/config).
type SpacesConfigRequest struct {
	Enabled bool `json:"enabled"`
}

// CreateSpaceRequest is the request body for Create (POST /spaces).
type CreateSpaceRequest struct {
	Name                 string   `json:"name"`
	Description          string   `json:"description,omitempty"`
	Icon                 string   `json:"icon,omitempty"`
	Instructions         string   `json:"instructions,omitempty"`
	SharedSpaceIds       []string `json:"sharedSpaceIds,omitempty"`
	DefaultRetrievalMode string   `json:"defaultRetrievalMode,omitempty"`
}

// UpdateSpaceRequest is the request body for Update (PATCH /spaces/{id}).
type UpdateSpaceRequest struct {
	Name                 *string  `json:"name,omitempty"`
	Description          *string  `json:"description,omitempty"`
	Icon                 *string  `json:"icon,omitempty"`
	Instructions         *string  `json:"instructions,omitempty"`
	SharedSpaceIds       []string `json:"sharedSpaceIds,omitempty"`
	DefaultRetrievalMode *string  `json:"defaultRetrievalMode,omitempty"`
}
