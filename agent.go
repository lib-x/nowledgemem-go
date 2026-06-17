package nowledgemem

import (
	"context"
	"fmt"
	"net/http"
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

// GetTokenUsage returns token usage summary for background agents.
//
// GET /agent/token-usage
func (s *AgentService) GetTokenUsage(ctx context.Context) (map[string]any, error) {
	var resp map[string]any
	if err := s.client.do(ctx, "GET", "/agent/token-usage", nil, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// CancelKnowledgeProcessingTask cancels a queued or running processing task.
//
// POST /agent/knowledge-processing/tasks/cancel
func (s *AgentService) CancelKnowledgeProcessingTask(ctx context.Context, req *CancelKnowledgeProcessingTaskRequest) (map[string]any, error) {
	var resp map[string]any
	if err := s.client.do(ctx, "POST", "/agent/knowledge-processing/tasks/cancel", req, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// PlanCommunityDetection returns a read-only community detection plan.
//
// GET /agent/trigger/community-detection/plan
func (s *AgentService) PlanCommunityDetection(ctx context.Context, params *CommunityDetectionPlanParams) (map[string]any, error) {
	q := url.Values{}
	if params != nil {
		if params.Resolution > 0 {
			q.Set("resolution", strconv.FormatFloat(params.Resolution, 'f', -1, 64))
		}
		if params.GenerateAISummary != nil {
			q.Set("generate_ai_summary", strconv.FormatBool(*params.GenerateAISummary))
		}
		if params.ScanLimit > 0 {
			q.Set("scan_limit", strconv.Itoa(params.ScanLimit))
		}
	}
	return s.getAgentMap(ctx, "/agent/trigger/community-detection/plan", q)
}

// PlanKGExtraction returns a read-only KG extraction backfill plan.
//
// GET /agent/trigger/kg-extraction/plan
func (s *AgentService) PlanKGExtraction(ctx context.Context, params *KGExtractionPlanParams) (map[string]any, error) {
	q := url.Values{}
	if params != nil {
		if params.ScanLimit > 0 {
			q.Set("scan_limit", strconv.Itoa(params.ScanLimit))
		}
		if params.BatchSize > 0 {
			q.Set("batch_size", strconv.Itoa(params.BatchSize))
		}
	}
	return s.getAgentMap(ctx, "/agent/trigger/kg-extraction/plan", q)
}

// PlanMemoryCompaction returns a read-only memory compaction plan.
//
// GET /agent/trigger/memory-compaction/plan
func (s *AgentService) PlanMemoryCompaction(ctx context.Context, params *MemoryCompactionPlanParams) (map[string]any, error) {
	q := url.Values{}
	if params != nil {
		if params.Limit > 0 {
			q.Set("limit", strconv.Itoa(params.Limit))
		}
		if params.SpaceID != "" {
			q.Set("space_id", params.SpaceID)
		}
	}
	return s.getAgentMap(ctx, "/agent/trigger/memory-compaction/plan", q)
}

// PlanMemoryCreated returns a read-only plan for memory-created events.
//
// GET /agent/trigger/memory-created/plan
func (s *AgentService) PlanMemoryCreated(ctx context.Context) (map[string]any, error) {
	return s.getAgentMap(ctx, "/agent/trigger/memory-created/plan", nil)
}

// PlanThreadSynced returns a read-only plan for thread-synced events.
//
// GET /agent/trigger/thread-synced/plan
func (s *AgentService) PlanThreadSynced(ctx context.Context, maxTasks int) (map[string]any, error) {
	q := url.Values{}
	if maxTasks > 0 {
		q.Set("max_tasks", strconv.Itoa(maxTasks))
	}
	return s.getAgentMap(ctx, "/agent/trigger/thread-synced/plan", q)
}

// PlanScheduledContext returns planned context for a scheduled task type.
//
// GET /agent/trigger/{task_type}/context-plan
func (s *AgentService) PlanScheduledContext(ctx context.Context, taskType, spaceID string) (map[string]any, error) {
	q := url.Values{}
	if spaceID != "" {
		q.Set("space_id", spaceID)
	}
	path := fmt.Sprintf("/agent/trigger/%s/context-plan", url.PathEscape(taskType))
	return s.getAgentMap(ctx, path, q)
}

// DryRunRuleReview previews a guidance rule review without applying changes.
//
// GET /agent/trigger/rule-review/dry-run
func (s *AgentService) DryRunRuleReview(ctx context.Context, params *RuleReviewDryRunParams) (map[string]any, error) {
	q := url.Values{}
	if params != nil {
		for _, query := range params.Queries {
			q.Add("queries", query)
		}
		if params.MemoryID != "" {
			q.Set("memory_id", params.MemoryID)
		}
		for _, memoryID := range params.MemoryIDs {
			q.Add("memory_ids", memoryID)
		}
		if params.SpaceID != "" {
			q.Set("space_id", params.SpaceID)
		}
		if params.MaxEvidenceResults > 0 {
			q.Set("max_evidence_results", strconv.Itoa(params.MaxEvidenceResults))
		}
	}
	return s.getAgentMap(ctx, "/agent/trigger/rule-review/dry-run", q)
}

// TriggerLabelConsolidation starts label consolidation.
//
// POST /agent/trigger/label-consolidation
func (s *AgentService) TriggerLabelConsolidation(ctx context.Context, dryRun *bool) (map[string]any, error) {
	q := url.Values{}
	if dryRun != nil {
		q.Set("dry_run", strconv.FormatBool(*dryRun))
	}
	return s.postAgentMapWithQuery(ctx, "/agent/trigger/label-consolidation", q, nil)
}

// TriggerRuleReview starts guidance rule review.
//
// POST /agent/trigger/rule-review
func (s *AgentService) TriggerRuleReview(ctx context.Context) (map[string]any, error) {
	return s.postAgentMap(ctx, "/agent/trigger/rule-review", nil)
}

// TriggerSkillCompile starts skill compile.
//
// POST /agent/trigger/skill-compile
func (s *AgentService) TriggerSkillCompile(ctx context.Context, skillID string) (map[string]any, error) {
	q := url.Values{"skill_id": {skillID}}
	return s.postAgentMapWithQuery(ctx, "/agent/trigger/skill-compile", q, nil)
}

// TriggerSkillReview starts skill review.
//
// POST /agent/trigger/skill-review
func (s *AgentService) TriggerSkillReview(ctx context.Context) (map[string]any, error) {
	return s.postAgentMap(ctx, "/agent/trigger/skill-review", nil)
}

// TriggerUnitTypeReclassification starts or previews unit type reclassification.
//
// POST /agent/trigger/unit-type-reclassification
func (s *AgentService) TriggerUnitTypeReclassification(ctx context.Context, req *UnitTypeReclassificationRequest) (map[string]any, error) {
	return s.postAgentMap(ctx, "/agent/trigger/unit-type-reclassification", req)
}

// SkillBuilderPropose asks the skill builder to propose a skill draft.
//
// POST /agent/skill-builder/propose
func (s *AgentService) SkillBuilderPropose(ctx context.Context, req *SkillBuilderProposeRequest) (map[string]any, error) {
	return s.postAgentMap(ctx, "/agent/skill-builder/propose", req)
}

// SkillBuilderChat sends a skill builder chat turn.
//
// POST /agent/skill-builder/chat
func (s *AgentService) SkillBuilderChat(ctx context.Context, req *SkillBuilderChatRequest) (map[string]any, error) {
	return s.postAgentMap(ctx, "/agent/skill-builder/chat", req)
}

// SkillBuilderRefine refines a compiled skill.
//
// POST /agent/skill-builder/refine
func (s *AgentService) SkillBuilderRefine(ctx context.Context, req *SkillRefineRequest) (map[string]any, error) {
	return s.postAgentMap(ctx, "/agent/skill-builder/refine", req)
}

// SkillBuilderRefineStream streams skill refinement output.
//
// POST /agent/skill-builder/refine/stream
func (s *AgentService) SkillBuilderRefineStream(ctx context.Context, req *SkillRefineRequest) (*http.Response, error) {
	return s.client.doStream(ctx, http.MethodPost, "/agent/skill-builder/refine/stream", nil, req)
}

// SkillBuilderEditBody updates a skill body directly.
//
// POST /agent/skill-builder/edit-body
func (s *AgentService) SkillBuilderEditBody(ctx context.Context, req *SkillEditBodyRequest) (map[string]any, error) {
	return s.postAgentMap(ctx, "/agent/skill-builder/edit-body", req)
}

// SkillBuilderImport imports a skill into managed scope.
//
// POST /agent/skill-builder/import
func (s *AgentService) SkillBuilderImport(ctx context.Context, req *SkillImportRequest) (map[string]any, error) {
	return s.postAgentMap(ctx, "/agent/skill-builder/import", req)
}

// SkillBuilderDiscoverImportable returns importable skill folders.
//
// GET /agent/skill-builder/discover-importable
func (s *AgentService) SkillBuilderDiscoverImportable(ctx context.Context) (map[string]any, error) {
	return s.getAgentMap(ctx, "/agent/skill-builder/discover-importable", nil)
}

// SkillBuilderPreviewImportable previews an importable skill path.
//
// GET /agent/skill-builder/preview-importable
func (s *AgentService) SkillBuilderPreviewImportable(ctx context.Context, path string) (map[string]any, error) {
	q := url.Values{"path": {path}}
	return s.getAgentMap(ctx, "/agent/skill-builder/preview-importable", q)
}

func (s *AgentService) getAgentMap(ctx context.Context, path string, q url.Values) (map[string]any, error) {
	var resp map[string]any
	if err := s.client.doQuery(ctx, path, q, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *AgentService) postAgentMap(ctx context.Context, path string, body any) (map[string]any, error) {
	var resp map[string]any
	if err := s.client.do(ctx, http.MethodPost, path, body, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *AgentService) postAgentMapWithQuery(ctx context.Context, path string, q url.Values, body any) (map[string]any, error) {
	var resp map[string]any
	if err := s.client.doWithQuery(ctx, http.MethodPost, path, q, body, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// --- Agent request types ---

// KGExtractionRequest is the request body for POST /agent/trigger/kg-extraction.
type KGExtractionRequest struct {
	MemoryIDs []string `json:"memory_ids,omitempty"`
	Scope     string   `json:"scope,omitempty"`
}

// CancelKnowledgeProcessingTaskRequest is the request for POST /agent/knowledge-processing/tasks/cancel.
type CancelKnowledgeProcessingTaskRequest struct {
	TaskID   *string `json:"task_id,omitempty"`
	TaskType *string `json:"task_type,omitempty"`
}

// CommunityDetectionPlanParams are query parameters for PlanCommunityDetection.
type CommunityDetectionPlanParams struct {
	Resolution        float64 `json:"resolution,omitempty"`
	GenerateAISummary *bool   `json:"generate_ai_summary,omitempty"`
	ScanLimit         int     `json:"scan_limit,omitempty"`
}

// KGExtractionPlanParams are query parameters for PlanKGExtraction.
type KGExtractionPlanParams struct {
	ScanLimit int `json:"scan_limit,omitempty"`
	BatchSize int `json:"batch_size,omitempty"`
}

// MemoryCompactionPlanParams are query parameters for PlanMemoryCompaction.
type MemoryCompactionPlanParams struct {
	Limit   int    `json:"limit,omitempty"`
	SpaceID string `json:"space_id,omitempty"`
}

// RuleReviewDryRunParams are query parameters for DryRunRuleReview.
type RuleReviewDryRunParams struct {
	Queries            []string `json:"queries,omitempty"`
	MemoryID           string   `json:"memory_id,omitempty"`
	MemoryIDs          []string `json:"memory_ids,omitempty"`
	SpaceID            string   `json:"space_id,omitempty"`
	MaxEvidenceResults int      `json:"max_evidence_results,omitempty"`
}

// UnitTypeReclassificationRequest is the request for TriggerUnitTypeReclassification.
type UnitTypeReclassificationRequest struct {
	Limit           int      `json:"limit,omitempty"`
	ScanLimit       int      `json:"scan_limit,omitempty"`
	MinConfidence   float64  `json:"min_confidence,omitempty"`
	DryRun          bool     `json:"dry_run,omitempty"`
	UnitTypes       []string `json:"unit_types,omitempty"`
	TargetUnitTypes []string `json:"target_unit_types,omitempty"`
	SpaceID         *string  `json:"space_id,omitempty"`
}

// SkillBuilderMessage is one prior skill-builder conversation turn.
type SkillBuilderMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// SkillBuilderChatRequest is the request for POST /agent/skill-builder/chat.
type SkillBuilderChatRequest struct {
	Message string                `json:"message"`
	History []SkillBuilderMessage `json:"history,omitempty"`
	Context map[string]any        `json:"context,omitempty"`
	SpaceID *string               `json:"space_id,omitempty"`
}

// SkillBuilderProposeRequest is the request for POST /agent/skill-builder/propose.
type SkillBuilderProposeRequest struct {
	Goal    string                `json:"goal"`
	History []SkillBuilderMessage `json:"history,omitempty"`
	SpaceID *string               `json:"space_id,omitempty"`
	Deep    bool                  `json:"deep,omitempty"`
}

// SkillEditBodyRequest is the request for POST /agent/skill-builder/edit-body.
type SkillEditBodyRequest struct {
	SkillID   string  `json:"skill_id"`
	Body      string  `json:"body"`
	Rationale *string `json:"rationale,omitempty"`
}

// SkillImportRequest is the request for POST /agent/skill-builder/import.
type SkillImportRequest struct {
	SkillMD       *string           `json:"skill_md,omitempty"`
	Path          *string           `json:"path,omitempty"`
	URL           *string           `json:"url,omitempty"`
	Files         map[string]string `json:"files,omitempty"`
	Force         bool              `json:"force,omitempty"`
	TakeOwnership bool              `json:"take_ownership,omitempty"`
}

// SkillRefineRequest is the request for POST /agent/skill-builder/refine.
type SkillRefineRequest struct {
	SkillID     string  `json:"skill_id"`
	Instruction *string `json:"instruction,omitempty"`
}

// EvolutionEdge represents an EVOLVES relationship.
type EvolutionEdge struct {
	SourceID string `json:"source_id"`
	TargetID string `json:"target_id"`
	EdgeType string `json:"edge_type"`
}

// KnowledgeProcessingStatus is the response for GET /agent/knowledge-processing/status.
type KnowledgeProcessingStatus struct {
	Enabled bool   `json:"enabled"`
	Status  string `json:"status"`
	LastRun string `json:"last_run,omitempty"`
	NextRun string `json:"next_run,omitempty"`
}

// GraphIntelligenceStatus is the response for GET /agent/graph-intelligence/status.
type GraphIntelligenceStatus struct {
	Running   bool   `json:"running"`
	SessionID string `json:"session_id,omitempty"`
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
