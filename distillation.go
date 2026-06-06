package onledgemem

import "context"

// DistillationService handles memory distillation from threads.
type DistillationService struct {
	client *Client
}

// Triage performs a lightweight check: does this conversation have save-worthy content?
func (s *DistillationService) Triage(ctx context.Context, req *TriageRequest) (*TriageResponse, error) {
	var resp TriageResponse
	if err := s.client.do(ctx, "POST", "/memories/distill/triage", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Preview previews distillation results without creating memories.
func (s *DistillationService) Preview(ctx context.Context, req *DistillPreviewRequest) (*DistillPreviewResponse, error) {
	var resp DistillPreviewResponse
	if err := s.client.do(ctx, "POST", "/memories/distill/preview", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Distill creates memories from thread content after distillation.
func (s *DistillationService) Distill(ctx context.Context, req *DistillRequest) (*DistillResponse, error) {
	var resp DistillResponse
	if err := s.client.do(ctx, "POST", "/memories/distill", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Plan creates a distillation plan.
func (s *DistillationService) Plan(ctx context.Context, req *DistillPlanRequest) (*DistillPlanResponse, error) {
	var resp DistillPlanResponse
	if err := s.client.do(ctx, "POST", "/memories/distill/plan", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Schedule schedules a distillation job.
func (s *DistillationService) Schedule(ctx context.Context, req *DistillScheduleRequest) (*DistillScheduleResponse, error) {
	var resp DistillScheduleResponse
	if err := s.client.do(ctx, "POST", "/memories/distill/schedule", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// --- Distillation types ---

// TriageRequest is the request for POST /memories/distill/triage.
type TriageRequest struct {
	ThreadID      string  `json:"thread_id"`
	ThreadContent *string `json:"thread_content,omitempty"`
	SpaceID       *string `json:"space_id,omitempty"`
}

// TriageResponse is the response for POST /memories/distill/triage.
type TriageResponse struct {
	Worthy    bool    `json:"worthy"`
	Reason    string  `json:"reason,omitempty"`
	Confidence float64 `json:"confidence,omitempty"`
}

// DistillPreviewRequest is the request for POST /memories/distill/preview.
type DistillPreviewRequest struct {
	ThreadID              string   `json:"thread_id"`
	ThreadTitle           *string  `json:"thread_title,omitempty"`
	ThreadContent         *string  `json:"thread_content,omitempty"`
	DistillationType      *string  `json:"distillation_type,omitempty"`
	ExtractionLevel       *string  `json:"extraction_level,omitempty"`
	SelectedMessageIndices []int   `json:"selected_message_indices,omitempty"`
	PreferredLanguage     *string  `json:"preferred_language,omitempty"`
	SpaceID               *string  `json:"space_id,omitempty"`
}

// DistillPreviewResponse is the response for POST /memories/distill/preview.
type DistillPreviewResponse struct {
	ProposedMemories []ProposedMemory `json:"proposed_memories"`
	CacheKey         string           `json:"cache_key,omitempty"`
}

// ProposedMemory is a memory proposed by distillation preview.
type ProposedMemory struct {
	Content    string   `json:"content"`
	Title      string   `json:"title,omitempty"`
	Importance float64  `json:"importance,omitempty"`
	Confidence float64  `json:"confidence,omitempty"`
	Labels     []string `json:"labels,omitempty"`
	UnitType   string   `json:"unit_type,omitempty"`
}

// DistillRequest is the request for POST /memories/distill.
type DistillRequest struct {
	ThreadID              string   `json:"thread_id"`
	ThreadTitle           *string  `json:"thread_title,omitempty"`
	ThreadContent         *string  `json:"thread_content,omitempty"`
	DistillationType      *string  `json:"distillation_type,omitempty"`
	ExtractionLevel       *string  `json:"extraction_level,omitempty"`
	CacheKey              *string  `json:"cache_key,omitempty"`
	SelectedMessageIndices []int   `json:"selected_message_indices,omitempty"`
	PreferredLanguage     *string  `json:"preferred_language,omitempty"`
	ForceDistill          bool     `json:"force_distill,omitempty"`
	SpaceID               *string  `json:"space_id,omitempty"`
}

// DistillResponse is the response for POST /memories/distill.
type DistillResponse struct {
	CreatedMemories  []Memory `json:"created_memories"`
	Skipped          bool     `json:"skipped,omitempty"`
	SkipReason       string   `json:"skip_reason,omitempty"`
}

// DistillPlanRequest is the request for POST /memories/distill/plan.
type DistillPlanRequest struct {
	ThreadID string  `json:"thread_id"`
	SpaceID  *string `json:"space_id,omitempty"`
}

// DistillPlanResponse is the response for POST /memories/distill/plan.
type DistillPlanResponse struct {
	Plan string `json:"plan"`
}

// DistillScheduleRequest is the request for POST /memories/distill/schedule.
type DistillScheduleRequest struct {
	ThreadID string  `json:"thread_id"`
	SpaceID  *string `json:"space_id,omitempty"`
}

// DistillScheduleResponse is the response for POST /memories/distill/schedule.
type DistillScheduleResponse struct {
	Scheduled bool   `json:"scheduled"`
	JobID     string `json:"job_id,omitempty"`
}
