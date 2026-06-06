package onledgemem

import "context"

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

// --- Agent request types ---

// KGExtractionRequest is the request body for POST /agent/trigger/kg-extraction.
type KGExtractionRequest struct {
	MemoryIDs []string `json:"memory_ids,omitempty"`
	Scope     string   `json:"scope,omitempty"`
}
