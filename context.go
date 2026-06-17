package nowledgemem

import (
	"context"
	"net/url"
	"strconv"
)

// ContextService handles context bundle operations.
type ContextService struct {
	client *Client
}

// GetBundle returns a merged context bundle for an agent or space.
//
// GET /context/bundle
func (s *ContextService) GetBundle(ctx context.Context, params *ContextBundleParams) (map[string]any, error) {
	q := url.Values{}
	if params != nil {
		if params.AgentID != "" {
			q.Set("agent_id", params.AgentID)
		}
		if params.SourceApp != "" {
			q.Set("source_app", params.SourceApp)
		}
		if params.HostAgentID != "" {
			q.Set("host_agent_id", params.HostAgentID)
		}
		if params.SpaceID != "" {
			q.Set("space_id", params.SpaceID)
		}
		if params.IncludeWorkingMemory != nil {
			q.Set("include_working_memory", strconv.FormatBool(*params.IncludeWorkingMemory))
		}
	}

	var resp map[string]any
	if err := s.client.doQuery(ctx, "/context/bundle", q, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// ContextBundleParams are query parameters for GetBundle.
type ContextBundleParams struct {
	AgentID              string `json:"agent_id,omitempty"`
	SourceApp            string `json:"source_app,omitempty"`
	HostAgentID          string `json:"host_agent_id,omitempty"`
	SpaceID              string `json:"space_id,omitempty"`
	IncludeWorkingMemory *bool  `json:"include_working_memory,omitempty"`
}
