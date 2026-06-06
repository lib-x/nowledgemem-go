package nowledgemem

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

// AgentService handles Background Intelligence operations.
type AgentService struct {
	client *Client
}

// Status returns the agent's current status.
func (s *AgentService) Status(ctx context.Context) (*AgentStatus, error) {
	var resp AgentStatus
	if err := s.client.do(ctx, "GET", "/agent/status", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// TriggerDailyBriefing triggers a daily briefing.
func (s *AgentService) TriggerDailyBriefing(ctx context.Context) error {
	return s.client.do(ctx, "POST", "/agent/trigger/daily-briefing", nil, nil)
}

// TriggerCrystallization triggers a crystallization review.
func (s *AgentService) TriggerCrystallization(ctx context.Context) error {
	return s.client.do(ctx, "POST", "/agent/trigger/crystallization", nil, nil)
}

// TriggerInsightDetection triggers proactive insight detection.
func (s *AgentService) TriggerInsightDetection(ctx context.Context) error {
	return s.client.do(ctx, "POST", "/agent/trigger/insight-detection", nil, nil)
}

// TriggerDecayRefresh triggers a decay score refresh.
func (s *AgentService) TriggerDecayRefresh(ctx context.Context) error {
	return s.client.do(ctx, "POST", "/agent/trigger/decay-refresh", nil, nil)
}

// TriggerMemoryCompaction triggers a memory compaction review.
func (s *AgentService) TriggerMemoryCompaction(ctx context.Context) error {
	return s.client.do(ctx, "POST", "/agent/trigger/memory-compaction", nil, nil)
}

// TriggerCommunityDetection triggers community detection on the knowledge graph.
func (s *AgentService) TriggerCommunityDetection(ctx context.Context) error {
	return s.client.do(ctx, "POST", "/agent/trigger/community-detection", nil, nil)
}

// TriggerKGExtraction triggers knowledge graph extraction.
func (s *AgentService) TriggerKGExtraction(ctx context.Context, req *KGExtractionRequest) error {
	return s.client.do(ctx, "POST", "/agent/trigger/kg-extraction", req, nil)
}

// GetEvolves returns EVOLVES relationships between memories.
func (s *AgentService) GetEvolves(ctx context.Context) ([]EvolutionEdge, error) {
	var resp []EvolutionEdge
	if err := s.client.do(ctx, "GET", "/agent/evolves", nil, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// GetKnowledgeProcessingStatus returns knowledge processing settings and status.
func (s *AgentService) GetKnowledgeProcessingStatus(ctx context.Context) (*KnowledgeProcessingStatus, error) {
	var resp KnowledgeProcessingStatus
	if err := s.client.do(ctx, "GET", "/agent/knowledge-processing/status", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetGraphIntelligenceStatus returns graph intelligence status.
func (s *AgentService) GetGraphIntelligenceStatus(ctx context.Context) (*GraphIntelligenceStatus, error) {
	var resp GraphIntelligenceStatus
	if err := s.client.do(ctx, "GET", "/agent/graph-intelligence/status", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetGraphIntelligenceSession returns a graph intelligence session.
func (s *AgentService) GetGraphIntelligenceSession(ctx context.Context, sessionID string) (*GraphIntelligenceSession, error) {
	var resp GraphIntelligenceSession
	q := url.Values{}
	if sessionID != "" {
		q.Set("session_id", sessionID)
	}
	if err := s.client.doQuery(ctx, "/agent/graph-intelligence/session", q, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// CreateGraphIntelligenceSession creates a new graph intelligence session.
func (s *AgentService) CreateGraphIntelligenceSession(ctx context.Context) (*GraphIntelligenceSession, error) {
	var resp GraphIntelligenceSession
	if err := s.client.do(ctx, "POST", "/agent/graph-intelligence/session", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// SendGraphIntelligenceMessage sends a message to graph intelligence.
func (s *AgentService) SendGraphIntelligenceMessage(ctx context.Context, req *GraphIntelligenceMessageRequest) (*GraphIntelligenceMessageResponse, error) {
	var resp GraphIntelligenceMessageResponse
	if err := s.client.do(ctx, "POST", "/agent/graph-intelligence/message", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetAINowSessions returns AI Now sessions.
func (s *AgentService) GetAINowSessions(ctx context.Context, limit int) ([]AINowSession, error) {
	q := url.Values{}
	if limit > 0 {
		q.Set("limit", strconv.Itoa(limit))
	}
	var resp []AINowSession
	if err := s.client.doQuery(ctx, "/agent/ai-now/sessions", q, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// CreateAINowSession creates a new AI Now session.
func (s *AgentService) CreateAINowSession(ctx context.Context, req *CreateAINowSessionRequest) (*AINowSession, error) {
	var resp AINowSession
	if err := s.client.do(ctx, "POST", "/agent/ai-now/sessions", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// --- Agent request types ---

// KGExtractionRequest is the request body for POST /agent/trigger/kg-extraction.
type KGExtractionRequest struct {
	MemoryIDs []string `json:"memory_ids,omitempty"`
	Scope     string   `json:"scope,omitempty"`
}

// EvolutionEdge represents an EVOLVES relationship.
type EvolutionEdge struct {
	SourceID string `json:"source_id"`
	TargetID string `json:"target_id"`
	EdgeType string `json:"edge_type"`
}

// KnowledgeProcessingStatus is the response for GET /agent/knowledge-processing/status.
type KnowledgeProcessingStatus struct {
	Enabled  bool   `json:"enabled"`
	Status   string `json:"status"`
	LastRun  string `json:"last_run,omitempty"`
	NextRun  string `json:"next_run,omitempty"`
}

// GraphIntelligenceStatus is the response for GET /agent/graph-intelligence/status.
type GraphIntelligenceStatus struct {
	Running    bool   `json:"running"`
	SessionID  string `json:"session_id,omitempty"`
}

// GraphIntelligenceSession is a graph intelligence session.
type GraphIntelligenceSession struct {
	SessionID string `json:"session_id"`
	Status    string `json:"status"`
}

// GraphIntelligenceMessageRequest is the request for POST /agent/graph-intelligence/message.
type GraphIntelligenceMessageRequest struct {
	SessionID string `json:"session_id"`
	Message   string `json:"message"`
}

// GraphIntelligenceMessageResponse is the response for POST /agent/graph-intelligence/message.
type GraphIntelligenceMessageResponse struct {
	Response  string `json:"response"`
	SessionID string `json:"session_id"`
}

// AINowSession represents an AI Now session.
type AINowSession struct {
	ID        string `json:"id"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at,omitempty"`
}

// CreateAINowSessionRequest is the request for POST /agent/ai-now/sessions.
type CreateAINowSessionRequest struct {
	Prompt  string `json:"prompt,omitempty"`
	Context string `json:"context,omitempty"`
}

// GetAINowSession returns a specific AI Now session.
func (s *AgentService) GetAINowSession(ctx context.Context, sessionID string) (*AINowSession, error) {
	var resp AINowSession
	path := fmt.Sprintf("/agent/ai-now/sessions/%s", url.PathEscape(sessionID))
	if err := s.client.do(ctx, "GET", path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// UpdateAINowSession updates an AI Now session.
func (s *AgentService) UpdateAINowSession(ctx context.Context, sessionID string, req map[string]any) error {
	path := fmt.Sprintf("/agent/ai-now/sessions/%s", url.PathEscape(sessionID))
	return s.client.do(ctx, "PATCH", path, req, nil)
}

// DeleteAINowSession deletes an AI Now session.
func (s *AgentService) DeleteAINowSession(ctx context.Context, sessionID string) error {
	path := fmt.Sprintf("/agent/ai-now/sessions/%s", url.PathEscape(sessionID))
	return s.client.do(ctx, "DELETE", path, nil, nil)
}

// CloseAINowSession closes an AI Now session.
func (s *AgentService) CloseAINowSession(ctx context.Context, sessionID string) error {
	path := fmt.Sprintf("/agent/ai-now/sessions/%s/close", url.PathEscape(sessionID))
	return s.client.do(ctx, "POST", path, nil, nil)
}

// CancelAINowSession cancels an AI Now session.
func (s *AgentService) CancelAINowSession(ctx context.Context, sessionID string) error {
	path := fmt.Sprintf("/agent/ai-now/sessions/%s/cancel", url.PathEscape(sessionID))
	return s.client.do(ctx, "POST", path, nil, nil)
}

// GetAINowSessionEvents returns events for an AI Now session.
func (s *AgentService) GetAINowSessionEvents(ctx context.Context, sessionID string) ([]AINowEvent, error) {
	var resp []AINowEvent
	path := fmt.Sprintf("/agent/ai-now/sessions/%s/events", url.PathEscape(sessionID))
	if err := s.client.do(ctx, "GET", path, nil, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// GetAINowSessionMessages returns messages for an AI Now session.
func (s *AgentService) GetAINowSessionMessages(ctx context.Context, sessionID string) ([]AINowMessage, error) {
	var resp []AINowMessage
	path := fmt.Sprintf("/agent/ai-now/sessions/%s/messages", url.PathEscape(sessionID))
	if err := s.client.do(ctx, "GET", path, nil, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// SendAINowSessionMessage sends a message to an AI Now session.
func (s *AgentService) SendAINowSessionMessage(ctx context.Context, sessionID string, req *AINowMessageRequest) (*AINowMessage, error) {
	var resp AINowMessage
	path := fmt.Sprintf("/agent/ai-now/sessions/%s/messages", url.PathEscape(sessionID))
	if err := s.client.do(ctx, "POST", path, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// PromptAINowSession sends a prompt to an AI Now session.
func (s *AgentService) PromptAINowSession(ctx context.Context, sessionID string, req *AINowPromptRequest) (*AINowPromptResponse, error) {
	var resp AINowPromptResponse
	path := fmt.Sprintf("/agent/ai-now/sessions/%s/prompt", url.PathEscape(sessionID))
	if err := s.client.do(ctx, "POST", path, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetAINowAutoApprove returns auto-approve status for an AI Now session.
func (s *AgentService) GetAINowAutoApprove(ctx context.Context, sessionID string) (*AINowAutoApprove, error) {
	var resp AINowAutoApprove
	path := fmt.Sprintf("/agent/ai-now/sessions/%s/auto-approve", url.PathEscape(sessionID))
	if err := s.client.do(ctx, "GET", path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// SetAINowAutoApprove sets auto-approve for an AI Now session.
func (s *AgentService) SetAINowAutoApprove(ctx context.Context, sessionID string, req *AINowAutoApproveRequest) error {
	path := fmt.Sprintf("/agent/ai-now/sessions/%s/auto-approve", url.PathEscape(sessionID))
	return s.client.do(ctx, "POST", path, req, nil)
}

// ReadAINowSessionFile reads a file from an AI Now session.
func (s *AgentService) ReadAINowSessionFile(ctx context.Context, sessionID string, req *AINowFileReadRequest) (*AINowFileReadResponse, error) {
	var resp AINowFileReadResponse
	path := fmt.Sprintf("/agent/ai-now/sessions/%s/files/read", url.PathEscape(sessionID))
	if err := s.client.do(ctx, "POST", path, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// RequestAINowPermission requests permission in an AI Now session.
func (s *AgentService) RequestAINowPermission(ctx context.Context, sessionID, requestID string, req *AINowPermissionRequest) error {
	path := fmt.Sprintf("/agent/ai-now/sessions/%s/permissions/%s", url.PathEscape(sessionID), url.PathEscape(requestID))
	return s.client.do(ctx, "POST", path, req, nil)
}

// SendAINowSkillPrompt sends a skill prompt to AI Now.
func (s *AgentService) SendAINowSkillPrompt(ctx context.Context, req *AINowSkillPromptRequest) (*AINowSkillPromptResponse, error) {
	var resp AINowSkillPromptResponse
	if err := s.client.do(ctx, "POST", "/agent/ai-now/skill-prompts", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// --- AI Now types ---

// AINowEvent represents an AI Now session event.
type AINowEvent struct {
	ID        string `json:"id"`
	EventType string `json:"event_type"`
	Content   string `json:"content,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
}

// AINowMessage represents an AI Now session message.
type AINowMessage struct {
	ID        string `json:"id"`
	Role      string `json:"role"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at,omitempty"`
}

// AINowMessageRequest is the request for POST /agent/ai-now/sessions/{id}/messages.
type AINowMessageRequest struct {
	Content string `json:"content"`
}

// AINowPromptRequest is the request for POST /agent/ai-now/sessions/{id}/prompt.
type AINowPromptRequest struct {
	Prompt string `json:"prompt"`
}

// AINowPromptResponse is the response for POST /agent/ai-now/sessions/{id}/prompt.
type AINowPromptResponse struct {
	Response string `json:"response"`
}

// AINowAutoApprove is the response for GET /agent/ai-now/sessions/{id}/auto-approve.
type AINowAutoApprove struct {
	Enabled bool `json:"enabled"`
}

// AINowAutoApproveRequest is the request for POST /agent/ai-now/sessions/{id}/auto-approve.
type AINowAutoApproveRequest struct {
	Enabled bool `json:"enabled"`
}

// AINowFileReadRequest is the request for POST /agent/ai-now/sessions/{id}/files/read.
type AINowFileReadRequest struct {
	Path string `json:"path"`
}

// AINowFileReadResponse is the response for POST /agent/ai-now/sessions/{id}/files/read.
type AINowFileReadResponse struct {
	Content string `json:"content"`
}

// AINowPermissionRequest is the request for POST /agent/ai-now/sessions/{id}/permissions/{request_id}.
type AINowPermissionRequest struct {
	Approved bool   `json:"approved"`
	Reason   string `json:"reason,omitempty"`
}

// AINowSkillPromptRequest is the request for POST /agent/ai-now/skill-prompts.
type AINowSkillPromptRequest struct {
	SkillID string `json:"skill_id"`
	Prompt  string `json:"prompt,omitempty"`
}

// AINowSkillPromptResponse is the response for POST /agent/ai-now/skill-prompts.
type AINowSkillPromptResponse struct {
	Response string `json:"response"`
}
