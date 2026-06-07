package nowledgemem

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// ThreadsService handles thread operations.
type ThreadsService struct {
	client *Client
}

// List returns threads with filtering and pagination.
func (s *ThreadsService) List(ctx context.Context, params *ListThreadsParams) (*ListThreadsResponse, error) {
	q := url.Values{}
	if params != nil {
		if params.Limit > 0 {
			q.Set("limit", strconv.Itoa(params.Limit))
		}
		if params.Offset > 0 {
			q.Set("offset", strconv.Itoa(params.Offset))
		}
		if params.Source != "" {
			q.Set("source", params.Source)
		}
		if params.SpaceID != "" {
			q.Set("space_id", params.SpaceID)
		}
	}
	var resp ListThreadsResponse
	if err := s.client.doQuery(ctx, "/threads", q, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Create creates a new thread with messages.
func (s *ThreadsService) Create(ctx context.Context, req *CreateThreadRequest) (*CreateThreadResponse, error) {
	var resp CreateThreadResponse
	if err := s.client.do(ctx, "POST", "/threads", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Get retrieves a thread with messages and pagination.
func (s *ThreadsService) Get(ctx context.Context, threadID string) (*Thread, error) {
	var resp Thread
	path := fmt.Sprintf("/threads/%s", url.PathEscape(threadID))
	if err := s.client.do(ctx, "GET", path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Delete deletes a thread and optionally its extracted memories.
func (s *ThreadsService) Delete(ctx context.Context, threadID string) error {
	path := fmt.Sprintf("/threads/%s", url.PathEscape(threadID))
	return s.client.do(ctx, "DELETE", path, nil, nil)
}

// Search performs full thread search with message matching.
func (s *ThreadsService) Search(ctx context.Context, query string, limit int) ([]ThreadListItem, error) {
	q := url.Values{}
	q.Set("q", query)
	if limit > 0 {
		q.Set("limit", strconv.Itoa(limit))
	}
	var resp []ThreadListItem
	if err := s.client.doQuery(ctx, "/threads/search", q, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// Summaries returns all thread titles and summaries.
func (s *ThreadsService) Summaries(ctx context.Context) ([]ThreadSummary, error) {
	var resp []ThreadSummary
	if err := s.client.do(ctx, "GET", "/threads/summaries", nil, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// AppendMessages appends messages to an existing thread.
func (s *ThreadsService) AppendMessages(ctx context.Context, threadID string, messages []MessageCreateRequest) (*AppendMessagesResponse, error) {
	var resp AppendMessagesResponse
	path := fmt.Sprintf("/threads/%s/append", url.PathEscape(threadID))
	if err := s.client.do(ctx, "POST", path, messages, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ToggleFavorite toggles favorite status for a thread.
func (s *ThreadsService) ToggleFavorite(ctx context.Context, threadID string) (*ToggleFavoriteResponse, error) {
	var resp ToggleFavoriteResponse
	path := fmt.Sprintf("/threads/%s/favorite", url.PathEscape(threadID))
	if err := s.client.do(ctx, "POST", path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// BulkDelete deletes multiple threads at once using POST (preferred over DELETE with body).
func (s *ThreadsService) BulkDelete(ctx context.Context, threadIDs []string) (*BulkDeleteResponse, error) {
	var resp BulkDeleteResponse
	body := map[string]any{"thread_ids": threadIDs}
	if err := s.client.do(ctx, "POST", "/threads/bulk/delete", body, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// --- Thread types ---

// ThreadSummary is a lightweight thread summary.
type ThreadSummary struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Summary string `json:"summary"`
}

// AppendMessagesResponse is the response for POST /threads/{id}/append.
type AppendMessagesResponse struct {
	Thread   Thread          `json:"thread"`
	Messages []ThreadMessage `json:"messages"`
}

// Import imports threads from JSON messages or conversation markdown.
func (s *ThreadsService) Import(ctx context.Context, req *ImportThreadsRequest) (*ImportThreadsResponse, error) {
	var resp ImportThreadsResponse
	if err := s.client.do(ctx, "POST", "/threads/import", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Parse parses thread content from various formats.
func (s *ThreadsService) Parse(ctx context.Context, req *ParseContentRequest) (*ParseContentResponse, error) {
	var resp ParseContentResponse
	if err := s.client.do(ctx, "POST", "/threads/parse", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetSources returns thread sources.
func (s *ThreadsService) GetSources(ctx context.Context) ([]ThreadSource, error) {
	var resp []ThreadSource
	if err := s.client.do(ctx, "GET", "/threads/sources", nil, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// GetImportConfig returns the current import configuration.
func (s *ThreadsService) GetImportConfig(ctx context.Context) (*ImportConfig, error) {
	var resp ImportConfig
	if err := s.client.do(ctx, "GET", "/threads/import-config", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// UpdateImportConfig updates import configuration.
func (s *ThreadsService) UpdateImportConfig(ctx context.Context, req *UpdateImportConfigRequest) error {
	return s.client.do(ctx, "PUT", "/threads/import-config", req, nil)
}

// GetWatcherStatus returns the status of the session watcher.
func (s *ThreadsService) GetWatcherStatus(ctx context.Context) (*WatcherStatus, error) {
	var resp WatcherStatus
	if err := s.client.do(ctx, "GET", "/threads/watcher/status", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// StartWatcher starts auto-importing sessions.
func (s *ThreadsService) StartWatcher(ctx context.Context) error {
	return s.client.do(ctx, "POST", "/threads/watcher/start", nil, nil)
}

// StopWatcher stops the session watcher.
func (s *ThreadsService) StopWatcher(ctx context.Context) error {
	return s.client.do(ctx, "POST", "/threads/watcher/stop", nil, nil)
}

// DiscoverSessions scans for conversation files from AI assistants.
func (s *ThreadsService) DiscoverSessions(ctx context.Context) ([]DiscoveredSession, error) {
	var resp []DiscoveredSession
	if err := s.client.do(ctx, "GET", "/threads/conversations/discover", nil, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// ImportConversation imports an external conversation file.
func (s *ThreadsService) ImportConversation(ctx context.Context, req *ImportConversationRequest) (*ImportConversationResponse, error) {
	var resp ImportConversationResponse
	if err := s.client.do(ctx, "POST", "/threads/conversations/import", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// SaveSession saves coding sessions as conversation threads.
func (s *ThreadsService) SaveSession(ctx context.Context, req *SaveSessionRequest) (*SaveSessionResponse, error) {
	var resp SaveSessionResponse
	if err := s.client.do(ctx, "POST", "/threads/sessions/save", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// BulkMovePreview previews a bulk move between spaces.
func (s *ThreadsService) BulkMovePreview(ctx context.Context, req *ThreadBulkMovePreviewRequest) (*ThreadBulkMovePreviewResponse, error) {
	var resp ThreadBulkMovePreviewResponse
	if err := s.client.do(ctx, "POST", "/threads/bulk/move/preview", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// BulkMove moves selected threads into another space.
func (s *ThreadsService) BulkMove(ctx context.Context, req *ThreadBulkMoveRequest) (*ThreadBulkMoveResponse, error) {
	var resp ThreadBulkMoveResponse
	if err := s.client.do(ctx, "POST", "/threads/bulk/move", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// BulkDeleteSelection deletes selected threads using a selector.
func (s *ThreadsService) BulkDeleteSelection(ctx context.Context, req *ThreadBulkDeleteSelectionRequest) (*BulkDeleteResponse, error) {
	var resp BulkDeleteResponse
	if err := s.client.do(ctx, "POST", "/threads/bulk/delete", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Export exports a thread in various formats.
func (s *ThreadsService) Export(ctx context.Context, threadID, format string) ([]byte, error) {
	q := url.Values{}
	if format != "" {
		q.Set("format", format)
	}
	path := fmt.Sprintf("/threads/%s/export", url.PathEscape(threadID))
	return s.client.doBytes(ctx, http.MethodGet, path, q, nil)
}

// GetCoverage returns a coverage report for debugging.
func (s *ThreadsService) GetCoverage(ctx context.Context, threadID string) (*ThreadCoverage, error) {
	var resp ThreadCoverage
	path := fmt.Sprintf("/threads/%s/coverage", url.PathEscape(threadID))
	if err := s.client.do(ctx, "GET", path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// --- Thread extended types ---

// ImportThreadsRequest is the request for POST /threads/import.
type ImportThreadsRequest struct {
	Content string `json:"content"`
	Format  string `json:"format,omitempty"`
	Source  string `json:"source,omitempty"`
	SpaceID string `json:"space_id,omitempty"`
}

// ImportThreadsResponse is the response for POST /threads/import.
type ImportThreadsResponse struct {
	Threads  []Thread `json:"threads"`
	Imported int      `json:"imported"`
	Skipped  int      `json:"skipped"`
}

// ParseContentRequest is the request for POST /threads/parse.
type ParseContentRequest struct {
	Content string `json:"content"`
	Format  string `json:"format,omitempty"`
}

// ParseContentResponse is the response for POST /threads/parse.
type ParseContentResponse struct {
	Messages []ThreadMessage `json:"messages"`
}

// ThreadSource represents a thread source.
type ThreadSource struct {
	Source string `json:"source"`
	Count  int    `json:"count"`
}

// ImportConfig is the response for GET /threads/import-config.
type ImportConfig struct {
	AutoImport   bool     `json:"auto_import"`
	Sources      []string `json:"sources,omitempty"`
	ExcludePaths []string `json:"exclude_paths,omitempty"`
}

// UpdateImportConfigRequest is the request for PUT /threads/import-config.
type UpdateImportConfigRequest struct {
	AutoImport   *bool    `json:"auto_import,omitempty"`
	Sources      []string `json:"sources,omitempty"`
	ExcludePaths []string `json:"exclude_paths,omitempty"`
}

// WatcherStatus is the response for GET /threads/watcher/status.
type WatcherStatus struct {
	Running   bool   `json:"running"`
	LastScan  string `json:"last_scan,omitempty"`
	ScanCount int    `json:"scan_count"`
}

// DiscoveredSession represents a discovered conversation session.
type DiscoveredSession struct {
	Path      string `json:"path"`
	Source    string `json:"source"`
	Project   string `json:"project,omitempty"`
	SessionID string `json:"session_id,omitempty"`
	Title     string `json:"title,omitempty"`
	Messages  int    `json:"messages"`
	Date      string `json:"date,omitempty"`
}

// ImportConversationRequest is the request for POST /threads/conversations/import.
type ImportConversationRequest struct {
	Path    string `json:"path"`
	SpaceID string `json:"space_id,omitempty"`
}

// ImportConversationResponse is the response for POST /threads/conversations/import.
type ImportConversationResponse struct {
	Thread   Thread `json:"thread"`
	Imported bool   `json:"imported"`
}

// SaveSessionRequest is the request for POST /threads/sessions/save.
type SaveSessionRequest struct {
	SessionID string                 `json:"session_id"`
	Title     string                 `json:"title,omitempty"`
	Source    string                 `json:"source,omitempty"`
	Project   string                 `json:"project,omitempty"`
	Messages  []MessageCreateRequest `json:"messages"`
	SpaceID   string                 `json:"space_id,omitempty"`
}

// SaveSessionResponse is the response for POST /threads/sessions/save.
type SaveSessionResponse struct {
	Thread  Thread `json:"thread"`
	Created bool   `json:"created"`
}

// ThreadBulkMovePreviewRequest is the request for POST /threads/bulk/move/preview.
type ThreadBulkMovePreviewRequest struct {
	ThreadIDs []string `json:"thread_ids,omitempty"`
	FromSpace string   `json:"from_space,omitempty"`
	ToSpace   string   `json:"to_space"`
}

// ThreadBulkMovePreviewResponse is the response for POST /threads/bulk/move/preview.
type ThreadBulkMovePreviewResponse struct {
	WillMove  int      `json:"will_move"`
	Conflicts []string `json:"conflicts,omitempty"`
}

// ThreadBulkMoveRequest is the request for POST /threads/bulk/move.
type ThreadBulkMoveRequest struct {
	ThreadIDs []string `json:"thread_ids,omitempty"`
	FromSpace string   `json:"from_space,omitempty"`
	ToSpace   string   `json:"to_space"`
}

// ThreadBulkMoveResponse is the response for POST /threads/bulk/move.
type ThreadBulkMoveResponse struct {
	Moved  int `json:"moved"`
	Failed int `json:"failed"`
}

// ThreadBulkDeleteSelectionRequest is the request for POST /threads/bulk/delete.
type ThreadBulkDeleteSelectionRequest struct {
	ThreadIDs []string `json:"thread_ids,omitempty"`
	SpaceID   string   `json:"space_id,omitempty"`
}

// ThreadCoverage is the response for GET /threads/{id}/coverage.
type ThreadCoverage struct {
	ThreadID      string  `json:"thread_id"`
	MessageCount  int     `json:"message_count"`
	MemoryCount   int     `json:"memory_count"`
	CoverageRatio float64 `json:"coverage_ratio"`
	UncoveredMsgs int     `json:"uncovered_msgs"`
}

// PreviewConversation loads a richer head-and-tail preview for one discovered conversation before import.
func (s *ThreadsService) PreviewConversation(ctx context.Context, req *PreviewConversationRequest) (*PreviewConversationResponse, error) {
	var resp PreviewConversationResponse
	if err := s.client.do(ctx, "POST", "/threads/conversations/preview", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ExportRaw exports a raw conversation file as markdown or JSON without importing.
func (s *ThreadsService) ExportRaw(ctx context.Context, req *ExportRawRequest) ([]byte, error) {
	return s.client.doBytes(ctx, http.MethodPost, "/threads/conversations/export-raw", nil, req)
}

// ReconcileTail reconciles the tail of a thread.
func (s *ThreadsService) ReconcileTail(ctx context.Context, threadID string) error {
	path := fmt.Sprintf("/threads/%s/reconcile-tail", url.PathEscape(threadID))
	return s.client.do(ctx, "POST", path, nil, nil)
}

// HideProject hides a project from the browse view.
func (s *ThreadsService) HideProject(ctx context.Context, project string) error {
	body := map[string]string{"project": project}
	return s.client.do(ctx, "POST", "/threads/import-config/hide-project", body, nil)
}

// UnhideProject unhides a project.
func (s *ThreadsService) UnhideProject(ctx context.Context, project string) error {
	body := map[string]string{"project": project}
	return s.client.do(ctx, "POST", "/threads/import-config/unhide-project", body, nil)
}

// HideSession hides a session from the browse view.
func (s *ThreadsService) HideSession(ctx context.Context, sessionID string) error {
	body := map[string]string{"session_id": sessionID}
	return s.client.do(ctx, "POST", "/threads/import-config/hide-session", body, nil)
}

// UnhideSession unhides a session.
func (s *ThreadsService) UnhideSession(ctx context.Context, sessionID string) error {
	body := map[string]string{"session_id": sessionID}
	return s.client.do(ctx, "POST", "/threads/import-config/unhide-session", body, nil)
}

// --- Thread preview/export types ---

// PreviewConversationRequest is the request for POST /threads/conversations/preview.
type PreviewConversationRequest struct {
	Path string `json:"path"`
}

// PreviewConversationResponse is the response for POST /threads/conversations/preview.
type PreviewConversationResponse struct {
	Preview  string `json:"preview"`
	Title    string `json:"title,omitempty"`
	Messages int    `json:"messages"`
}

// ExportRawRequest is the request for POST /threads/conversations/export-raw.
type ExportRawRequest struct {
	Path   string `json:"path"`
	Format string `json:"format,omitempty"`
}
