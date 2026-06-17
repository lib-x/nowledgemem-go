package nowledgemem

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// SkillsService handles skill discovery, authoring, lifecycle, and evaluation operations.
type SkillsService struct {
	client *Client
}

// List returns skills with optional stage, space, and pagination filters.
//
// GET /skills
func (s *SkillsService) List(ctx context.Context, params *ListSkillsParams) (map[string]any, error) {
	q := url.Values{}
	if params != nil {
		if params.Stage != "" {
			q.Set("stage", params.Stage)
		}
		if params.SpaceID != "" {
			q.Set("space_id", params.SpaceID)
		}
		if params.Limit > 0 {
			q.Set("limit", strconv.Itoa(params.Limit))
		}
		if params.Offset > 0 {
			q.Set("offset", strconv.Itoa(params.Offset))
		}
	}
	return s.getMap(ctx, "/skills", q)
}

// Match returns skills matching the given query.
//
// GET /skills/match
func (s *SkillsService) Match(ctx context.Context, query string, limit int) (map[string]any, error) {
	q := url.Values{"q": {query}}
	if limit > 0 {
		q.Set("limit", strconv.Itoa(limit))
	}
	return s.getMap(ctx, "/skills/match", q)
}

// Author creates an authored skill proposal from named sources.
//
// POST /skills/author
func (s *SkillsService) Author(ctx context.Context, req *AuthorSkillRequest) (map[string]any, error) {
	return s.postMap(ctx, "/skills/author", req)
}

// RegistrationStatus returns host registration state.
//
// GET /skills/registration
func (s *SkillsService) RegistrationStatus(ctx context.Context) (map[string]any, error) {
	return s.getMap(ctx, "/skills/registration", nil)
}

// RegisterHost registers a local skill host.
//
// POST /skills/registration/{host}
func (s *SkillsService) RegisterHost(ctx context.Context, host string) (map[string]any, error) {
	path := fmt.Sprintf("/skills/registration/%s", url.PathEscape(host))
	return s.postMap(ctx, path, nil)
}

// UnregisterHost removes a local skill host registration.
//
// DELETE /skills/registration/{host}
func (s *SkillsService) UnregisterHost(ctx context.Context, host string) (map[string]any, error) {
	var resp map[string]any
	path := fmt.Sprintf("/skills/registration/%s", url.PathEscape(host))
	if err := s.client.do(ctx, http.MethodDelete, path, nil, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// SetHostConfigDir sets or clears the config directory for one host.
//
// PUT /skills/registration/{host}/config-dir
func (s *SkillsService) SetHostConfigDir(ctx context.Context, host string, req *HostConfigDir) (map[string]any, error) {
	var resp map[string]any
	path := fmt.Sprintf("/skills/registration/%s/config-dir", url.PathEscape(host))
	if err := s.client.do(ctx, http.MethodPut, path, req, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// Activity returns recent skill activity.
//
// GET /skills/activity
func (s *SkillsService) Activity(ctx context.Context, params *SkillActivityParams) (map[string]any, error) {
	q := url.Values{}
	if params != nil {
		if params.SkillID != "" {
			q.Set("skill_id", params.SkillID)
		}
		if params.SpaceID != "" {
			q.Set("space_id", params.SpaceID)
		}
		if params.Limit > 0 {
			q.Set("limit", strconv.Itoa(params.Limit))
		}
	}
	return s.getMap(ctx, "/skills/activity", q)
}

// Get returns a skill by ID.
//
// GET /skills/{skill_id}
func (s *SkillsService) Get(ctx context.Context, skillID string, params *GetSkillParams) (map[string]any, error) {
	q := url.Values{}
	if params != nil && params.IncludeBody != nil {
		q.Set("include_body", strconv.FormatBool(*params.IncludeBody))
	}
	path := fmt.Sprintf("/skills/%s", url.PathEscape(skillID))
	return s.getMap(ctx, path, q)
}

// DuplicateOf returns duplicate information for a skill.
//
// GET /skills/{skill_id}/duplicate-of
func (s *SkillsService) DuplicateOf(ctx context.Context, skillID string) (map[string]any, error) {
	path := fmt.Sprintf("/skills/%s/duplicate-of", url.PathEscape(skillID))
	return s.getMap(ctx, path, nil)
}

// Eval returns evaluation summary for a skill.
//
// GET /skills/{skill_id}/eval
func (s *SkillsService) Eval(ctx context.Context, skillID string) (map[string]any, error) {
	path := fmt.Sprintf("/skills/%s/eval", url.PathEscape(skillID))
	return s.getMap(ctx, path, nil)
}

// RunEval runs evaluation for a skill.
//
// POST /skills/{skill_id}/eval/run
func (s *SkillsService) RunEval(ctx context.Context, skillID string, strict *bool) (map[string]any, error) {
	q := url.Values{}
	if strict != nil {
		q.Set("strict", strconv.FormatBool(*strict))
	}
	path := fmt.Sprintf("/skills/%s/eval/run", url.PathEscape(skillID))
	var resp map[string]any
	if err := s.client.doWithQuery(ctx, http.MethodPost, path, q, nil, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// StreamEval streams evaluation output for a skill.
//
// GET /skills/{skill_id}/eval/stream
func (s *SkillsService) StreamEval(ctx context.Context, skillID string, strict *bool) (*http.Response, error) {
	q := url.Values{}
	if strict != nil {
		q.Set("strict", strconv.FormatBool(*strict))
	}
	path := fmt.Sprintf("/skills/%s/eval/stream", url.PathEscape(skillID))
	return s.client.doStream(ctx, http.MethodGet, path, q, nil)
}

// ListEvalCases returns evaluation cases for a skill.
//
// GET /skills/{skill_id}/eval/cases
func (s *SkillsService) ListEvalCases(ctx context.Context, skillID string) (map[string]any, error) {
	path := fmt.Sprintf("/skills/%s/eval/cases", url.PathEscape(skillID))
	return s.getMap(ctx, path, nil)
}

// AddEvalCase adds an evaluation case to a skill.
//
// POST /skills/{skill_id}/eval/cases
func (s *SkillsService) AddEvalCase(ctx context.Context, skillID string, req map[string]any) (map[string]any, error) {
	path := fmt.Sprintf("/skills/%s/eval/cases", url.PathEscape(skillID))
	return s.postMap(ctx, path, req)
}

// UpdateEvalCase updates an evaluation case.
//
// PATCH /skills/{skill_id}/eval/cases/{task_id}
func (s *SkillsService) UpdateEvalCase(ctx context.Context, skillID, taskID string, req map[string]any) (map[string]any, error) {
	var resp map[string]any
	path := fmt.Sprintf("/skills/%s/eval/cases/%s", url.PathEscape(skillID), url.PathEscape(taskID))
	if err := s.client.do(ctx, http.MethodPatch, path, req, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// DeleteEvalCase deletes an evaluation case.
//
// DELETE /skills/{skill_id}/eval/cases/{task_id}
func (s *SkillsService) DeleteEvalCase(ctx context.Context, skillID, taskID string) (map[string]any, error) {
	var resp map[string]any
	path := fmt.Sprintf("/skills/%s/eval/cases/%s", url.PathEscape(skillID), url.PathEscape(taskID))
	if err := s.client.do(ctx, http.MethodDelete, path, nil, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// Harden hardens a skill for up to maxRounds.
//
// POST /skills/{skill_id}/harden
func (s *SkillsService) Harden(ctx context.Context, skillID string, maxRounds int) (map[string]any, error) {
	q := url.Values{}
	if maxRounds > 0 {
		q.Set("max_rounds", strconv.Itoa(maxRounds))
	}
	path := fmt.Sprintf("/skills/%s/harden", url.PathEscape(skillID))
	var resp map[string]any
	if err := s.client.doWithQuery(ctx, http.MethodPost, path, q, nil, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// Activate activates a skill.
//
// POST /skills/{skill_id}/activate
func (s *SkillsService) Activate(ctx context.Context, skillID string) (map[string]any, error) {
	return s.postSkillAction(ctx, skillID, "activate", nil)
}

// ApplyVersion applies the current proposed version of a skill.
//
// POST /skills/{skill_id}/apply-version
func (s *SkillsService) ApplyVersion(ctx context.Context, skillID string) (map[string]any, error) {
	return s.postSkillAction(ctx, skillID, "apply-version", nil)
}

// EditFile updates a bundled skill file.
//
// POST /skills/{skill_id}/edit-file
func (s *SkillsService) EditFile(ctx context.Context, skillID string, req *SkillEditFileRequest) (map[string]any, error) {
	return s.postSkillAction(ctx, skillID, "edit-file", req)
}

// ApplyEnrichment applies pending enrichment to a skill.
//
// POST /skills/{skill_id}/apply-enrichment
func (s *SkillsService) ApplyEnrichment(ctx context.Context, skillID string) (map[string]any, error) {
	return s.postSkillAction(ctx, skillID, "apply-enrichment", nil)
}

// DismissEnrichment dismisses pending enrichment for a skill.
//
// POST /skills/{skill_id}/dismiss-enrichment
func (s *SkillsService) DismissEnrichment(ctx context.Context, skillID string) (map[string]any, error) {
	return s.postSkillAction(ctx, skillID, "dismiss-enrichment", nil)
}

// Outcome reports the result of using a skill version.
//
// POST /skills/{skill_id}/outcome
func (s *SkillsService) Outcome(ctx context.Context, skillID string, req *SkillOutcomeRequest) (map[string]any, error) {
	return s.postSkillAction(ctx, skillID, "outcome", req)
}

// Deactivate deactivates a skill.
//
// POST /skills/{skill_id}/deactivate
func (s *SkillsService) Deactivate(ctx context.Context, skillID string) (map[string]any, error) {
	return s.postSkillAction(ctx, skillID, "deactivate", nil)
}

// Archive archives a skill.
//
// POST /skills/{skill_id}/archive
func (s *SkillsService) Archive(ctx context.Context, skillID string) (map[string]any, error) {
	return s.postSkillAction(ctx, skillID, "archive", nil)
}

// Dismiss dismisses a skill suggestion.
//
// POST /skills/{skill_id}/dismiss
func (s *SkillsService) Dismiss(ctx context.Context, skillID string, req *DismissSkillRequest) (map[string]any, error) {
	return s.postSkillAction(ctx, skillID, "dismiss", req)
}

// StrengthenInstead records that another skill should absorb this suggestion.
//
// POST /skills/{skill_id}/strengthen-instead
func (s *SkillsService) StrengthenInstead(ctx context.Context, skillID string, req *StrengthenInsteadRequest) (map[string]any, error) {
	return s.postSkillAction(ctx, skillID, "strengthen-instead", req)
}

// MergeInto merges a skill into another skill.
//
// POST /skills/{skill_id}/merge-into
func (s *SkillsService) MergeInto(ctx context.Context, skillID string, req *MergeIntoRequest) (map[string]any, error) {
	return s.postSkillAction(ctx, skillID, "merge-into", req)
}

// CuratorProposals returns pending curator proposals.
//
// GET /skills/curator/proposals
func (s *SkillsService) CuratorProposals(ctx context.Context) (map[string]any, error) {
	return s.getMap(ctx, "/skills/curator/proposals", nil)
}

// CuratorDryRun previews a curator run.
//
// GET /skills/curator/dry-run
func (s *SkillsService) CuratorDryRun(ctx context.Context, params *SkillCuratorDryRunParams) (map[string]any, error) {
	q := url.Values{}
	if params != nil {
		if params.Budget > 0 {
			q.Set("budget", strconv.Itoa(params.Budget))
		}
		if params.IncludeMerges != nil {
			q.Set("include_merges", strconv.FormatBool(*params.IncludeMerges))
		}
	}
	return s.getMap(ctx, "/skills/curator/dry-run", q)
}

// CuratorRun runs the skill curator.
//
// POST /skills/curator/run
func (s *SkillsService) CuratorRun(ctx context.Context, budget int) (map[string]any, error) {
	q := url.Values{}
	if budget > 0 {
		q.Set("budget", strconv.Itoa(budget))
	}
	var resp map[string]any
	if err := s.client.doWithQuery(ctx, http.MethodPost, "/skills/curator/run", q, nil, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *SkillsService) getMap(ctx context.Context, path string, q url.Values) (map[string]any, error) {
	var resp map[string]any
	if err := s.client.doQuery(ctx, path, q, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *SkillsService) postMap(ctx context.Context, path string, body any) (map[string]any, error) {
	var resp map[string]any
	if err := s.client.do(ctx, http.MethodPost, path, body, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *SkillsService) postSkillAction(ctx context.Context, skillID, action string, body any) (map[string]any, error) {
	path := fmt.Sprintf("/skills/%s/%s", url.PathEscape(skillID), action)
	return s.postMap(ctx, path, body)
}

// ListSkillsParams are query parameters for List.
type ListSkillsParams struct {
	Stage   string `json:"stage,omitempty"`
	SpaceID string `json:"space_id,omitempty"`
	Limit   int    `json:"limit,omitempty"`
	Offset  int    `json:"offset,omitempty"`
}

// GetSkillParams are query parameters for Get.
type GetSkillParams struct {
	IncludeBody *bool `json:"include_body,omitempty"`
}

// SkillActivityParams are query parameters for Activity.
type SkillActivityParams struct {
	SkillID string `json:"skill_id,omitempty"`
	SpaceID string `json:"space_id,omitempty"`
	Limit   int    `json:"limit,omitempty"`
}

// SkillCuratorDryRunParams are query parameters for CuratorDryRun.
type SkillCuratorDryRunParams struct {
	Budget        int   `json:"budget,omitempty"`
	IncludeMerges *bool `json:"include_merges,omitempty"`
}

// AuthorSkillRequest is the request for POST /skills/author.
type AuthorSkillRequest struct {
	Name      string   `json:"name"`
	Note      *string  `json:"note,omitempty"`
	MemoryIDs []string `json:"memory_ids,omitempty"`
	ThreadIDs []string `json:"thread_ids,omitempty"`
	SourceIDs []string `json:"source_ids,omitempty"`
}

// DismissSkillRequest is the request for POST /skills/{skill_id}/dismiss.
type DismissSkillRequest struct {
	RejectionReason *string `json:"rejection_reason,omitempty"`
	ResurfaceRule   *string `json:"resurface_rule,omitempty"`
}

// MergeIntoRequest is the request for POST /skills/{skill_id}/merge-into.
type MergeIntoRequest struct {
	IntoSkillID string `json:"into_skill_id"`
}

// SkillEditFileRequest is the request for POST /skills/{skill_id}/edit-file.
type SkillEditFileRequest struct {
	Path    string `json:"path"`
	Content string `json:"content"`
}

// SkillOutcomeRequest is the request for POST /skills/{skill_id}/outcome.
type SkillOutcomeRequest struct {
	SkillVersion int     `json:"skill_version"`
	Outcome      string  `json:"outcome"`
	Deviations   *string `json:"deviations,omitempty"`
	Missing      *string `json:"missing,omitempty"`
	Failure      *string `json:"failure,omitempty"`
	Source       string  `json:"source,omitempty"`
}

// StrengthenInsteadRequest is the request for POST /skills/{skill_id}/strengthen-instead.
type StrengthenInsteadRequest struct {
	TargetSkillID string `json:"target_skill_id"`
}

// HostConfigDir is the request for PUT /skills/registration/{host}/config-dir.
type HostConfigDir struct {
	ConfigDir *string `json:"config_dir,omitempty"`
}
