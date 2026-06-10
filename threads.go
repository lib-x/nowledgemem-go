package nowledgemem

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// ThreadsService provides methods for thread operations.
type ThreadsService struct {
	client *Client
}

// List returns threads with filtering and pagination.
//
// GET /threads
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
//
// POST /threads
func (s *ThreadsService) Create(ctx context.Context, req *CreateThreadRequest) (*CreateThreadResponse, error) {
	var resp CreateThreadResponse
	if err := s.client.do(ctx, "POST", "/threads", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Get retrieves a thread with messages, supports pagination.
//
// GET /threads/{id}
func (s *ThreadsService) Get(ctx context.Context, threadID string, params *GetThreadParams) (*GetThreadResponse, error) {
	q := url.Values{}
	if params != nil {
		if params.Limit > 0 {
			q.Set("limit", strconv.Itoa(params.Limit))
		}
		if params.Offset > 0 {
			q.Set("offset", strconv.Itoa(params.Offset))
		}
		if params.SpaceID != "" {
			q.Set("space_id", params.SpaceID)
		}
	}
	path := fmt.Sprintf("/threads/%s", url.PathEscape(threadID))
	var resp GetThreadResponse
	if err := s.client.doQuery(ctx, path, q, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Delete deletes a thread and optionally its extracted memories.
//
// DELETE /threads/{id}
func (s *ThreadsService) Delete(ctx context.Context, threadID string, params *DeleteThreadParams) (*DeleteThreadResponse, error) {
	q := url.Values{}
	if params != nil {
		if params.CascadeDeleteMemories {
			q.Set("cascade_delete_memories", "true")
		}
		if params.SpaceID != "" {
			q.Set("space_id", params.SpaceID)
		}
	}
	path := fmt.Sprintf("/threads/%s", url.PathEscape(threadID))
	var resp DeleteThreadResponse
	if err := s.client.doQuery(ctx, path, q, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Search performs full thread search with message matching.
//
// GET /threads/search
func (s *ThreadsService) Search(ctx context.Context, params *SearchThreadsParams) (*SearchThreadsResponse, error) {
	q := url.Values{}
	if params != nil {
		if params.Query != "" {
			q.Set("query", params.Query)
		}
		if params.Mode != "" {
			q.Set("mode", params.Mode)
		}
		if params.Limit > 0 {
			q.Set("limit", strconv.Itoa(params.Limit))
		}
		if params.Source != "" {
			q.Set("source", params.Source)
		}
		if params.SpaceID != "" {
			q.Set("space_id", params.SpaceID)
		}
	}
	var resp SearchThreadsResponse
	if err := s.client.doQuery(ctx, "/threads/search", q, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Summaries returns all thread titles and summaries.
//
// GET /threads/summaries
func (s *ThreadsService) Summaries(ctx context.Context, spaceID string) (*ThreadSummariesResponse, error) {
	q := url.Values{}
	if spaceID != "" {
		q.Set("space_id", spaceID)
	}
	var resp ThreadSummariesResponse
	if err := s.client.doQuery(ctx, "/threads/summaries", q, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// AppendMessages appends messages to an existing thread.
//
// POST /threads/{id}/append
func (s *ThreadsService) AppendMessages(ctx context.Context, threadID string, req *AppendMessagesRequest) (*AppendMessagesResponse, error) {
	var resp AppendMessagesResponse
	path := fmt.Sprintf("/threads/%s/append", url.PathEscape(threadID))
	if err := s.client.do(ctx, "POST", path, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ToggleFavorite toggles the favorite status for a thread.
//
// POST /threads/{id}/favorite
func (s *ThreadsService) ToggleFavorite(ctx context.Context, threadID string) (*ToggleFavoriteResponse, error) {
	var resp ToggleFavoriteResponse
	path := fmt.Sprintf("/threads/%s/favorite", url.PathEscape(threadID))
	if err := s.client.do(ctx, "POST", path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// BulkDelete deletes multiple threads at once using POST.
//
// POST /threads/bulk/delete
func (s *ThreadsService) BulkDelete(ctx context.Context, threadIDs []string) (*ThreadBulkDeleteResponse, error) {
	var resp ThreadBulkDeleteResponse
	body := map[string]any{"thread_ids": threadIDs}
	if err := s.client.do(ctx, "POST", "/threads/bulk/delete", body, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// --- Thread types ---

// ThreadSummary represents a lightweight thread summary with title and summary text.
type ThreadSummary struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Summary string `json:"summary"`
}

// ThreadSummariesResponse contains the list of thread summaries.
type ThreadSummariesResponse struct {
	Summaries []ThreadSummary `json:"summaries"`
}

// SearchThreadsParams holds query parameters for thread search.
type SearchThreadsParams struct {
	Query   string `json:"query"`
	Mode    string `json:"mode,omitempty"` // "suggestions" or "full"
	Limit   int    `json:"limit,omitempty"`
	Source  string `json:"source,omitempty"`
	SpaceID string `json:"space_id,omitempty"`
}

// SearchMetadata holds metadata about a thread search result.
type SearchMetadata struct {
	Query                string `json:"query,omitempty"`
	Mode                 string `json:"mode,omitempty"`
	MatchedMessagesCount int    `json:"matched_messages_count,omitempty"`
	Error                string `json:"error,omitempty"`
}

// SearchThreadsResponse contains the results of a thread search.
type SearchThreadsResponse struct {
	Threads        []ThreadListItem `json:"threads"`
	TotalFound     int              `json:"total_found"`
	SearchMetadata SearchMetadata   `json:"search_metadata,omitempty"`
}

// GetThreadParams holds query parameters for getting a single thread.
type GetThreadParams struct {
	Limit   int    `json:"limit,omitempty"`
	Offset  int    `json:"offset,omitempty"`
	SpaceID string `json:"space_id,omitempty"`
}

// GetThreadResponse contains a thread with its messages and metadata.
type GetThreadResponse struct {
	Thread            Thread          `json:"thread"`
	Messages          []ThreadMessage `json:"messages,omitempty"`
	RelatedMemories   []Memory        `json:"related_memories,omitempty"`
	Entities          []string        `json:"entities,omitempty"`
	TotalMessages     int             `json:"total_messages,omitempty"`
	TotalTokens       int             `json:"total_tokens,omitempty"`
	CoveredMessageIDs []string        `json:"covered_message_ids,omitempty"`
}

// DeleteThreadParams holds query parameters for deleting a thread.
type DeleteThreadParams struct {
	CascadeDeleteMemories bool   `json:"cascade_delete_memories,omitempty"`
	SpaceID               string `json:"space_id,omitempty"`
}

// DeleteThreadResponse contains the result of deleting a thread.
type DeleteThreadResponse struct {
	Message         string `json:"message"`
	DeletedMessages int    `json:"deleted_messages,omitempty"`
	DeletedMemories int    `json:"deleted_memories,omitempty"`
	CascadeDeletion bool   `json:"cascade_deletion,omitempty"`
}

// AppendMessagesRequest holds the body for appending messages to a thread.
type AppendMessagesRequest struct {
	Messages       []MessageCreateRequest `json:"messages,omitempty"`
	FilePath       string                 `json:"file_path,omitempty"`
	Format         string                 `json:"format,omitempty"`
	Deduplicate    bool                   `json:"deduplicate,omitempty"`
	IdempotencyKey string                 `json:"idempotency_key,omitempty"`
	SpaceID        string                 `json:"space_id,omitempty"`
}

// AppendMessagesResponse contains the result of appending messages to a thread.
type AppendMessagesResponse struct {
	Success       bool   `json:"success"`
	ThreadID      string `json:"thread_id"`
	MessagesAdded int    `json:"messages_added"`
	TotalMessages int    `json:"total_messages"`
}

// ImportThreadItem represents a single thread in a thread import request.
type ImportThreadItem struct {
	ThreadID        string                 `json:"thread_id,omitempty"`
	Title           string                 `json:"title,omitempty"`
	Messages        []MessageCreateRequest `json:"messages,omitempty"`
	MarkdownContent string                 `json:"markdown_content,omitempty"`
	Source          string                 `json:"source,omitempty"`
	Participants    []string               `json:"participants,omitempty"`
	Project         string                 `json:"project,omitempty"`
	Workspace       string                 `json:"workspace,omitempty"`
	ToolVersion     string                 `json:"tool_version,omitempty"`
	Metadata        map[string]any         `json:"metadata,omitempty"`
}

// ImportThreadsRequest holds the body for importing threads.
type ImportThreadsRequest struct {
	ImportThreadItem
	Threads []ImportThreadItem `json:"threads,omitempty"` // batch mode
}

// ImportThreadResult represents the result of importing a single thread.
type ImportThreadResult struct {
	Success      bool   `json:"success"`
	ThreadID     string `json:"thread_id,omitempty"`
	Title        string `json:"title,omitempty"`
	MessageCount int    `json:"message_count,omitempty"`
	Error        string `json:"error,omitempty"`
}

// ImportThreadsResponse contains the result of a thread import operation.
type ImportThreadsResponse struct {
	Success       bool                 `json:"success"`
	ImportedCount int                  `json:"imported_count"`
	FailedCount   int                  `json:"failed_count"`
	Results       []ImportThreadResult `json:"results,omitempty"`
}

// Import imports threads from JSON messages or conversation markdown.
//
// POST /threads/import
func (s *ThreadsService) Import(ctx context.Context, req *ImportThreadsRequest) (*ImportThreadsResponse, error) {
	var resp ImportThreadsResponse
	if err := s.client.do(ctx, "POST", "/threads/import", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ParseContentRequest holds the body for parsing thread content.
type ParseContentRequest struct {
	FileContent string `json:"file_content"`
	FileName    string `json:"file_name"`
}

// ParseContentResponse contains the result of parsing thread content.
type ParseContentResponse struct {
	Success        bool           `json:"success"`
	ParsedThread   map[string]any `json:"parsed_thread,omitempty"`
	FormatDetected string         `json:"format_detected,omitempty"`
	Error          string         `json:"error,omitempty"`
}

// Parse parses thread content from various formats.
//
// POST /threads/parse
func (s *ThreadsService) Parse(ctx context.Context, req *ParseContentRequest) (*ParseContentResponse, error) {
	var resp ParseContentResponse
	if err := s.client.do(ctx, "POST", "/threads/parse", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetSources returns available thread sources.
//
// GET /threads/sources
func (s *ThreadsService) GetSources(ctx context.Context) ([]ThreadSource, error) {
	var resp []ThreadSource
	if err := s.client.do(ctx, "GET", "/threads/sources", nil, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// ThreadSource represents a thread source with its count.
type ThreadSource struct {
	Source string `json:"source"`
	Count  int    `json:"count"`
}

// AutoImportRule represents a single auto-import rule configuration.
type AutoImportRule struct {
	ID        string `json:"id"`
	Type      string `json:"type"`
	Value     string `json:"value"`
	Enabled   bool   `json:"enabled"`
	CreatedAt int64  `json:"created_at,omitempty"`
}

// ImportConfig represents the import configuration settings.
type ImportConfig struct {
	HiddenProjects      []string         `json:"hidden_projects,omitempty"`
	HiddenSessions      []string         `json:"hidden_sessions,omitempty"`
	AutoImportRules     []AutoImportRule `json:"auto_import_rules,omitempty"`
	WatcherEnabled      bool             `json:"watcher_enabled"`
	ShowHiddenByDefault bool             `json:"show_hidden_by_default"`
	DedupWindowSeconds  int              `json:"dedup_window_seconds,omitempty"`
	WatchedPlatforms    []string         `json:"watched_platforms,omitempty"`
	WatchedProjects     []string         `json:"watched_projects,omitempty"`
	CursorPollInterval  int              `json:"cursor_poll_interval,omitempty"`
}

// UpdateImportConfigRequest holds the body for updating import configuration.
type UpdateImportConfigRequest struct {
	HiddenProjects      *[]string         `json:"hidden_projects,omitempty"`
	HiddenSessions      *[]string         `json:"hidden_sessions,omitempty"`
	AutoImportRules     *[]AutoImportRule `json:"auto_import_rules,omitempty"`
	WatcherEnabled      *bool             `json:"watcher_enabled,omitempty"`
	ShowHiddenByDefault *bool             `json:"show_hidden_by_default,omitempty"`
	DedupWindowSeconds  *int              `json:"dedup_window_seconds,omitempty"`
	WatchedPlatforms    *[]string         `json:"watched_platforms,omitempty"`
	WatchedProjects     *[]string         `json:"watched_projects,omitempty"`
	CursorPollInterval  *float64          `json:"cursor_poll_interval,omitempty"`
}

// GetImportConfig returns the current import configuration.
//
// GET /threads/import-config
func (s *ThreadsService) GetImportConfig(ctx context.Context) (*ImportConfig, error) {
	var resp ImportConfig
	if err := s.client.do(ctx, "GET", "/threads/import-config", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// UpdateImportConfig updates the import configuration.
//
// PUT /threads/import-config
func (s *ThreadsService) UpdateImportConfig(ctx context.Context, req *UpdateImportConfigRequest) error {
	return s.client.do(ctx, "PUT", "/threads/import-config", req, nil)
}

// WatcherStatus represents the current status of the session watcher.
type WatcherStatus struct {
	Running   bool   `json:"running"`
	LastScan  string `json:"last_scan,omitempty"`
	ScanCount int    `json:"scan_count"`
}

// GetWatcherStatus returns the status of the session watcher.
//
// GET /threads/watcher/status
func (s *ThreadsService) GetWatcherStatus(ctx context.Context) (*WatcherStatus, error) {
	var resp WatcherStatus
	if err := s.client.do(ctx, "GET", "/threads/watcher/status", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// StartWatcher starts the session watcher for auto-importing sessions.
//
// POST /threads/watcher/start
func (s *ThreadsService) StartWatcher(ctx context.Context) error {
	return s.client.do(ctx, "POST", "/threads/watcher/start", nil, nil)
}

// StopWatcher stops the session watcher.
//
// POST /threads/watcher/stop
func (s *ThreadsService) StopWatcher(ctx context.Context) error {
	return s.client.do(ctx, "POST", "/threads/watcher/stop", nil, nil)
}

// DiscoveredSession represents a conversation session discovered by the scanner.
type DiscoveredSession struct {
	Path      string `json:"path"`
	Source    string `json:"source"`
	Project   string `json:"project,omitempty"`
	SessionID string `json:"session_id,omitempty"`
	Title     string `json:"title,omitempty"`
	Messages  int    `json:"messages"`
	Date      string `json:"date,omitempty"`
}

// DiscoverSessionsResponse contains discovered conversation sessions grouped by source.
type DiscoverSessionsResponse struct {
	Conversations map[string][]DiscoveredSession `json:"conversations"`
}

// DiscoverSessions scans for conversation files from Claude Code, Codex, Cursor, and OpenCode.
//
// GET /threads/conversations/discover
func (s *ThreadsService) DiscoverSessions(ctx context.Context, source string) (*DiscoverSessionsResponse, error) {
	q := url.Values{}
	if source != "" {
		q.Set("source", source)
	}
	var resp DiscoverSessionsResponse
	if err := s.client.doQuery(ctx, "/threads/conversations/discover", q, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ImportConversationRequest holds the body for importing a conversation file.
type ImportConversationRequest struct {
	Path               string         `json:"path"`
	Source             string         `json:"source"` // "claude", "codex", "cursor", "opencode"
	SessionID          string         `json:"session_id,omitempty"`
	ThreadIDOverride   string         `json:"thread_id_override,omitempty"`
	Summary            string         `json:"summary,omitempty"`
	AutoCompact        bool           `json:"auto_compact,omitempty"`
	PreserveTimestamps *bool          `json:"preserve_timestamps,omitempty"`
	Workspace          string         `json:"workspace,omitempty"`
	Project            string         `json:"project,omitempty"`
	Metadata           map[string]any `json:"metadata,omitempty"`
}

// ImportConversationResponse contains the result of importing a conversation file.
type ImportConversationResponse struct {
	Thread          Thread          `json:"thread"`
	Messages        []ThreadMessage `json:"messages,omitempty"`
	ImportSummary   string          `json:"import_summary,omitempty"`
	Warnings        []string        `json:"warnings,omitempty"`
	CreatedMemories []Memory        `json:"created_memories,omitempty"`
	SkippedMessages int             `json:"skipped_messages,omitempty"`
}

// ImportConversation imports an external conversation file into Nowledge Mem.
//
// POST /threads/conversations/import
func (s *ThreadsService) ImportConversation(ctx context.Context, req *ImportConversationRequest) (*ImportConversationResponse, error) {
	var resp ImportConversationResponse
	if err := s.client.do(ctx, "POST", "/threads/conversations/import", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// SaveSessionRequest holds the body for saving coding sessions as threads.
type SaveSessionRequest struct {
	Client               string `json:"client"` // "claude-code", "codex", "gemini-cli"
	ProjectPath          string `json:"project_path"`
	PersistMode          string `json:"persist_mode,omitempty"` // "current" or "all"
	SessionID            string `json:"session_id,omitempty"`
	Summary              string `json:"summary,omitempty"`
	TruncateLargeContent bool   `json:"truncate_large_content,omitempty"`
}

// SaveSessionResult represents the result of saving a single session.
type SaveSessionResult struct {
	Action        string `json:"action"`
	SessionID     string `json:"session_id,omitempty"`
	ThreadID      string `json:"thread_id,omitempty"`
	MessageCount  int    `json:"message_count,omitempty"`
	MessagesAdded int    `json:"messages_added,omitempty"`
	File          string `json:"file,omitempty"`
}

// SaveSessionResponse contains the result of saving coding sessions.
type SaveSessionResponse struct {
	Status      string              `json:"status"`
	Client      string              `json:"client,omitempty"`
	ProjectPath string              `json:"project_path,omitempty"`
	PersistMode string              `json:"persist_mode,omitempty"`
	Results     []SaveSessionResult `json:"results,omitempty"`
	Error       string              `json:"error,omitempty"`
	Hint        string              `json:"hint,omitempty"`
}

// SaveSession saves coding sessions as conversation threads with deduplication.
//
// POST /threads/sessions/save
func (s *ThreadsService) SaveSession(ctx context.Context, req *SaveSessionRequest) (*SaveSessionResponse, error) {
	var resp SaveSessionResponse
	if err := s.client.do(ctx, "POST", "/threads/sessions/save", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// BulkThreadSelection describes a set of threads to operate on in bulk operations.
type BulkThreadSelection struct {
	ThreadIDs []string `json:"thread_ids,omitempty"`
	SpaceID   string   `json:"space_id,omitempty"`
	SelectAll bool     `json:"select_all,omitempty"`
}

// ThreadBulkMovePreviewRequest holds the body for previewing a bulk thread move.
type ThreadBulkMovePreviewRequest struct {
	Selection     BulkThreadSelection `json:"selection"`
	TargetSpaceID string              `json:"target_space_id"`
}

// ThreadBulkMovePreviewResponse contains the preview of a bulk thread move.
type ThreadBulkMovePreviewResponse struct {
	Count         int              `json:"count"`
	MaxAllowed    int              `json:"max_allowed,omitempty"`
	LimitExceeded bool             `json:"limit_exceeded,omitempty"`
	SourceSpaceID string           `json:"source_space_id,omitempty"`
	TargetSpaceID string           `json:"target_space_id,omitempty"`
	SelectionMode string           `json:"selection_mode,omitempty"`
	ExcludedCount int              `json:"excluded_count,omitempty"`
	Conflicts     []map[string]any `json:"conflicts,omitempty"`
	Message       string           `json:"message,omitempty"`
}

// BulkMovePreview previews a bulk move between spaces and detects legacy conflicts.
//
// POST /threads/bulk/move/preview
func (s *ThreadsService) BulkMovePreview(ctx context.Context, req *ThreadBulkMovePreviewRequest) (*ThreadBulkMovePreviewResponse, error) {
	var resp ThreadBulkMovePreviewResponse
	if err := s.client.do(ctx, "POST", "/threads/bulk/move/preview", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ThreadBulkMoveRequest holds the body for moving threads between spaces.
type ThreadBulkMoveRequest struct {
	Selection     BulkThreadSelection `json:"selection"`
	TargetSpaceID string              `json:"target_space_id"`
}

// ThreadBulkMoveResponse contains the result of a bulk thread move.
type ThreadBulkMoveResponse struct {
	MovedCount    int              `json:"moved_count"`
	FailedCount   int              `json:"failed_count"`
	SourceSpaceID string           `json:"source_space_id,omitempty"`
	TargetSpaceID string           `json:"target_space_id,omitempty"`
	Conflicts     []map[string]any `json:"conflicts,omitempty"`
	Message       string           `json:"message,omitempty"`
}

// BulkMove moves selected threads into another space.
//
// POST /threads/bulk/move
func (s *ThreadsService) BulkMove(ctx context.Context, req *ThreadBulkMoveRequest) (*ThreadBulkMoveResponse, error) {
	var resp ThreadBulkMoveResponse
	if err := s.client.do(ctx, "POST", "/threads/bulk/move", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ThreadBulkDeleteSelectionRequest holds the body for bulk deleting threads by selection.
type ThreadBulkDeleteSelectionRequest struct {
	Selection             BulkThreadSelection `json:"selection"`
	CascadeDeleteMemories bool                `json:"cascade_delete_memories,omitempty"`
}

// ThreadBulkDeleteResponse contains the result of a bulk thread deletion.
type ThreadBulkDeleteResponse struct {
	Message              string          `json:"message"`
	DeletedCount         int             `json:"deleted_count"`
	FailedCount          int             `json:"failed_count"`
	TotalDeletedMessages int             `json:"total_deleted_messages,omitempty"`
	TotalDeletedMemories int             `json:"total_deleted_memories,omitempty"`
	CascadeDeletion      bool            `json:"cascade_deletion,omitempty"`
	Results              []map[string]any `json:"results,omitempty"`
}

// BulkDeleteSelection deletes selected threads using a backend-resolved selector.
//
// POST /threads/bulk/delete
func (s *ThreadsService) BulkDeleteSelection(ctx context.Context, req *ThreadBulkDeleteSelectionRequest) (*ThreadBulkDeleteResponse, error) {
	var resp ThreadBulkDeleteResponse
	if err := s.client.do(ctx, "POST", "/threads/bulk/delete", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Export exports a thread in various formats.
//
// GET /threads/{id}/export
func (s *ThreadsService) Export(ctx context.Context, threadID, format string) ([]byte, error) {
	q := url.Values{}
	if format != "" {
		q.Set("format", format)
	}
	path := fmt.Sprintf("/threads/%s/export", url.PathEscape(threadID))
	return s.client.doBytes(ctx, http.MethodGet, path, q, nil)
}

// GetCoverage returns a coverage report for a thread.
//
// GET /threads/{id}/coverage
func (s *ThreadsService) GetCoverage(ctx context.Context, threadID string) (*ThreadCoverage, error) {
	var resp ThreadCoverage
	path := fmt.Sprintf("/threads/%s/coverage", url.PathEscape(threadID))
	if err := s.client.do(ctx, "GET", path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ThreadCoverage represents the coverage ratio of memories to messages in a thread.
type ThreadCoverage struct {
	ThreadID      string  `json:"thread_id"`
	MessageCount  int     `json:"message_count"`
	MemoryCount   int     `json:"memory_count"`
	CoverageRatio float64 `json:"coverage_ratio"`
	UncoveredMsgs int     `json:"uncovered_msgs"`
}

// PreviewConversationRequest holds the body for previewing a conversation.
type PreviewConversationRequest struct {
	Path      string `json:"path"`
	Source    string `json:"source"` // "claude", "codex", "cursor", "opencode"
	SessionID string `json:"session_id,omitempty"`
}

// PreviewMessage represents a single message in a conversation preview response.
type PreviewMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// PreviewConversationResponse contains the head-and-tail preview of a conversation.
type PreviewConversationResponse struct {
	MessageCount    int              `json:"message_count"`
	PreviewMessages []PreviewMessage `json:"preview_messages,omitempty"`
}

// PreviewConversation loads a richer head-and-tail preview for one discovered conversation before import.
//
// POST /threads/conversations/preview
func (s *ThreadsService) PreviewConversation(ctx context.Context, req *PreviewConversationRequest) (*PreviewConversationResponse, error) {
	var resp PreviewConversationResponse
	if err := s.client.do(ctx, "POST", "/threads/conversations/preview", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ExportRawRequest holds the body for exporting a raw conversation file.
type ExportRawRequest struct {
	Path   string `json:"path"`
	Source string `json:"source"` // required
	Format string `json:"format,omitempty"`
}

// ExportRaw exports a raw conversation file as markdown or JSON without importing.
//
// POST /threads/conversations/export-raw
func (s *ThreadsService) ExportRaw(ctx context.Context, req *ExportRawRequest) ([]byte, error) {
	return s.client.doBytes(ctx, http.MethodPost, "/threads/conversations/export-raw", nil, req)
}

// ReconcileTail reconciles the tail of a thread.
//
// POST /threads/{id}/reconcile-tail
func (s *ThreadsService) ReconcileTail(ctx context.Context, threadID string) error {
	path := fmt.Sprintf("/threads/%s/reconcile-tail", url.PathEscape(threadID))
	return s.client.do(ctx, "POST", path, nil, nil)
}

// HideProject hides a project from the browse view.
//
// POST /threads/import-config/hide-project
func (s *ThreadsService) HideProject(ctx context.Context, project string) error {
	body := map[string]string{"project": project}
	return s.client.do(ctx, "POST", "/threads/import-config/hide-project", body, nil)
}

// UnhideProject unhides a project.
//
// POST /threads/import-config/unhide-project
func (s *ThreadsService) UnhideProject(ctx context.Context, project string) error {
	body := map[string]string{"project": project}
	return s.client.do(ctx, "POST", "/threads/import-config/unhide-project", body, nil)
}

// HideSession hides a session from the browse view.
//
// POST /threads/import-config/hide-session
func (s *ThreadsService) HideSession(ctx context.Context, sessionID string) error {
	body := map[string]string{"session_id": sessionID}
	return s.client.do(ctx, "POST", "/threads/import-config/hide-session", body, nil)
}

// UnhideSession unhides a session.
//
// POST /threads/import-config/unhide-session
func (s *ThreadsService) UnhideSession(ctx context.Context, sessionID string) error {
	body := map[string]string{"session_id": sessionID}
	return s.client.do(ctx, "POST", "/threads/import-config/unhide-session", body, nil)
}
