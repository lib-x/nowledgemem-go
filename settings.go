package nowledgemem

import (
	"context"
	"fmt"
	"net/url"
)

// SettingsService handles settings operations.
type SettingsService struct {
	client *Client
}

// GetProfile returns user profile, aliases, context, and preferred language.
func (s *SettingsService) GetProfile(ctx context.Context) (*UserProfile, error) {
	var resp UserProfile
	if err := s.client.do(ctx, "GET", "/settings/profile", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetAgentProfiles returns named long-running agent identities.
//
// GET /settings/agent-profiles
func (s *SettingsService) GetAgentProfiles(ctx context.Context) (*AgentProfilesResponse, error) {
	var resp AgentProfilesResponse
	if err := s.client.do(ctx, "GET", "/settings/agent-profiles", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// CreateAgentProfile creates an agent profile.
//
// POST /settings/agent-profiles
func (s *SettingsService) CreateAgentProfile(ctx context.Context, req *AgentProfilePayload) (*AgentProfileResponse, error) {
	var resp AgentProfileResponse
	if err := s.client.do(ctx, "POST", "/settings/agent-profiles", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// UpdateAgentProfile updates an agent profile.
//
// PUT /settings/agent-profiles/{agent_id}
func (s *SettingsService) UpdateAgentProfile(ctx context.Context, agentID string, req *AgentProfilePayload) (*AgentProfileResponse, error) {
	var resp AgentProfileResponse
	path := fmt.Sprintf("/settings/agent-profiles/%s", url.PathEscape(agentID))
	if err := s.client.do(ctx, "PUT", path, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// DeleteAgentProfile deletes an agent profile.
//
// DELETE /settings/agent-profiles/{agent_id}
func (s *SettingsService) DeleteAgentProfile(ctx context.Context, agentID string) (*AgentProfilesResponse, error) {
	var resp AgentProfilesResponse
	path := fmt.Sprintf("/settings/agent-profiles/%s", url.PathEscape(agentID))
	if err := s.client.do(ctx, "DELETE", path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetGuidanceRules returns owner-managed AI context rules.
//
// GET /settings/rules
func (s *SettingsService) GetGuidanceRules(ctx context.Context) (*GuidanceRulesResponse, error) {
	var resp GuidanceRulesResponse
	if err := s.client.do(ctx, "GET", "/settings/rules", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// CreateGuidanceRule creates an AI context guidance rule.
//
// POST /settings/rules
func (s *SettingsService) CreateGuidanceRule(ctx context.Context, req *GuidanceRulePayload) (*GuidanceRuleResponse, error) {
	var resp GuidanceRuleResponse
	if err := s.client.do(ctx, "POST", "/settings/rules", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// UpdateGuidanceRule updates an AI context guidance rule.
//
// PUT /settings/rules/{rule_id}
func (s *SettingsService) UpdateGuidanceRule(ctx context.Context, ruleID string, req *GuidanceRulePayload) (*GuidanceRuleResponse, error) {
	var resp GuidanceRuleResponse
	path := fmt.Sprintf("/settings/rules/%s", url.PathEscape(ruleID))
	if err := s.client.do(ctx, "PUT", path, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// DeleteGuidanceRule deletes an AI context guidance rule.
//
// DELETE /settings/rules/{rule_id}
func (s *SettingsService) DeleteGuidanceRule(ctx context.Context, ruleID string) (*GuidanceRulesResponse, error) {
	var resp GuidanceRulesResponse
	path := fmt.Sprintf("/settings/rules/%s", url.PathEscape(ruleID))
	if err := s.client.do(ctx, "DELETE", path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// --- Settings types ---

// UserProfile is the response for GET /settings/profile.
type UserProfile struct {
	Name               string `json:"name"`
	Aliases            string `json:"aliases"`
	Context            string `json:"context"`
	PreferredLanguage  string `json:"preferred_language"`
	CustomInstructions string `json:"custom_instructions"`
}

// AgentProfilesResponse is the response for GET /settings/agent-profiles.
type AgentProfilesResponse struct {
	AgentProfiles []AgentProfileResponse `json:"agentProfiles,omitempty"`
}

// AgentProfileResponse represents a named long-running agent identity.
type AgentProfileResponse struct {
	ID             string   `json:"id"`
	DisplayName    string   `json:"displayName"`
	Role           string   `json:"role,omitempty"`
	Description    string   `json:"description,omitempty"`
	Instructions   string   `json:"instructions,omitempty"`
	DefaultSpaceID string   `json:"defaultSpaceId,omitempty"`
	SourceApp      string   `json:"sourceApp,omitempty"`
	HostAgentID    string   `json:"hostAgentId,omitempty"`
	Tags           []string `json:"tags,omitempty"`
}

// AgentProfilePayload is the request body for agent profile create and update operations.
type AgentProfilePayload struct {
	ID             *string   `json:"id,omitempty"`
	DisplayName    *string   `json:"displayName,omitempty"`
	Role           *string   `json:"role,omitempty"`
	Description    *string   `json:"description,omitempty"`
	Instructions   *string   `json:"instructions,omitempty"`
	DefaultSpaceID *string   `json:"defaultSpaceId,omitempty"`
	SourceApp      *string   `json:"sourceApp,omitempty"`
	HostAgentID    *string   `json:"hostAgentId,omitempty"`
	Tags           *[]string `json:"tags,omitempty"`
}

// GuidanceRulesResponse is the response for GET /settings/rules.
type GuidanceRulesResponse struct {
	Rules []GuidanceRuleResponse `json:"rules,omitempty"`
}

// GuidanceRuleResponse represents an owner-managed AI context rule.
type GuidanceRuleResponse struct {
	ID                           string   `json:"id"`
	Title                        string   `json:"title"`
	Body                         string   `json:"body"`
	Scope                        string   `json:"scope,omitempty"`
	Status                       string   `json:"status,omitempty"`
	AgentProfileID               string   `json:"agentProfileId,omitempty"`
	SpaceID                      string   `json:"spaceId,omitempty"`
	Priority                     int      `json:"priority,omitempty"`
	Source                       string   `json:"source,omitempty"`
	Tags                         []string `json:"tags,omitempty"`
	EvidenceMemoryIDs            []string `json:"evidenceMemoryIds,omitempty"`
	SupportedEvidenceMemoryIDs   []string `json:"supportedEvidenceMemoryIds,omitempty"`
	UnsupportedEvidenceMemoryIDs []string `json:"unsupportedEvidenceMemoryIds,omitempty"`
	Confidence                   *float64 `json:"confidence,omitempty"`
	Rationale                    string   `json:"rationale,omitempty"`
	SupportCount                 int      `json:"supportCount,omitempty"`
	ArchivedAt                   string   `json:"archivedAt,omitempty"`
	ArchivedReason               string   `json:"archivedReason,omitempty"`
	ArchivedBy                   string   `json:"archivedBy,omitempty"`
	CreatedAt                    string   `json:"createdAt,omitempty"`
	UpdatedAt                    string   `json:"updatedAt,omitempty"`
}

// GuidanceRulePayload is the request body for guidance rule create and update operations.
type GuidanceRulePayload struct {
	ID                           *string   `json:"id,omitempty"`
	Title                        *string   `json:"title,omitempty"`
	Body                         *string   `json:"body,omitempty"`
	Scope                        *string   `json:"scope,omitempty"`
	Status                       *string   `json:"status,omitempty"`
	AgentProfileID               *string   `json:"agentProfileId,omitempty"`
	SpaceID                      *string   `json:"spaceId,omitempty"`
	Priority                     *int      `json:"priority,omitempty"`
	Source                       *string   `json:"source,omitempty"`
	Tags                         *[]string `json:"tags,omitempty"`
	EvidenceMemoryIDs            *[]string `json:"evidenceMemoryIds,omitempty"`
	SupportedEvidenceMemoryIDs   *[]string `json:"supportedEvidenceMemoryIds,omitempty"`
	UnsupportedEvidenceMemoryIDs *[]string `json:"unsupportedEvidenceMemoryIds,omitempty"`
	Confidence                   *float64  `json:"confidence,omitempty"`
	Rationale                    *string   `json:"rationale,omitempty"`
	SupportCount                 *int      `json:"supportCount,omitempty"`
	ArchivedAt                   *string   `json:"archivedAt,omitempty"`
	ArchivedReason               *string   `json:"archivedReason,omitempty"`
	ArchivedBy                   *string   `json:"archivedBy,omitempty"`
}
