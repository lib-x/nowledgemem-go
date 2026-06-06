package onledgemem

import "context"

// CapabilitiesService handles server capabilities.
type CapabilitiesService struct {
	client *Client
}

// Get returns server capabilities — unauthenticated, used by clients to adapt UI.
func (s *CapabilitiesService) Get(ctx context.Context) (*Capabilities, error) {
	var resp Capabilities
	if err := s.client.do(ctx, "GET", "/capabilities", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// --- Capabilities types ---

// Capabilities is the response for GET /capabilities.
type Capabilities struct {
	Version       string         `json:"version"`
	Features      map[string]bool `json:"features,omitempty"`
	SpacesEnabled bool           `json:"spaces_enabled,omitempty"`
}
