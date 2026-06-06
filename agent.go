package onledgemem

import (
	"context"
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
