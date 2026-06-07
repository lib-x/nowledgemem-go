package nowledgemem

import (
	"fmt"
	"time"
)

// Memory represents a memory node in the knowledge base.
type Memory struct {
	ID                 string                 `json:"id"`
	NodeType           string                 `json:"node_type,omitempty"`
	CreatedAt          *time.Time             `json:"created_at,omitempty"`
	UpdatedAt          *time.Time             `json:"updated_at,omitempty"`
	Metadata           map[string]any         `json:"metadata,omitempty"`
	Content            string                 `json:"content"`
	Title              string                 `json:"title,omitempty"`
	Importance         float64                `json:"importance,omitempty"`
	Confidence         float64                `json:"confidence,omitempty"`
	PagerankScore      float64                `json:"pagerank_score,omitempty"`
	Embedding          []float64              `json:"embedding,omitempty"`
	SourceRange        map[string]int         `json:"source_range,omitempty"`
	Source             string                 `json:"source,omitempty"`
	SpaceID            string                 `json:"space_id,omitempty"`
	SemanticField      string                 `json:"semantic_field,omitempty"`
	ReindexNeeded      bool                   `json:"reindex_needed,omitempty"`
	LastReindexedAt    *time.Time             `json:"last_reindexed_at,omitempty"`
	LastAccessedAt     *time.Time             `json:"last_accessed_at,omitempty"`
	AccessCount        int                    `json:"access_count,omitempty"`
	Appearances        int                    `json:"appearances,omitempty"`
	Clicks             int                    `json:"clicks,omitempty"`
	TotalDwellTimeMs   int                    `json:"total_dwell_time_ms,omitempty"`
	LastClickedAt      *time.Time             `json:"last_clicked_at,omitempty"`
	DecayScoreCached   float64                `json:"decay_score_cached,omitempty"`
	TemporalContext    string                 `json:"temporal_context,omitempty"`
	TemporalType       string                 `json:"temporal_type,omitempty"`
	EventStart         string                 `json:"event_start,omitempty"`
	EventEnd           string                 `json:"event_end,omitempty"`
	TemporalPrecision  string                 `json:"temporal_precision,omitempty"`
	TemporalConfidence float64                `json:"temporal_confidence,omitempty"`
	UnitType           string                 `json:"unit_type,omitempty"`
	IsLatest           bool                   `json:"is_latest,omitempty"`
	Version            int                    `json:"version,omitempty"`
	IsCrystal          bool                   `json:"is_crystal,omitempty"`
	CrystalTitle       string                 `json:"crystal_title,omitempty"`
	SourceUnitCount    int                    `json:"source_unit_count,omitempty"`
	ExtractionMethod   string                 `json:"extraction_method,omitempty"`
	LastEvaluatedAt    *time.Time             `json:"last_evaluated_at,omitempty"`
	ReviewStatus       string                 `json:"review_status,omitempty"`
}

// MemoryListItem is the summary view returned by list endpoints.
type MemoryListItem struct {
	ID           string         `json:"id"`
	Title        string         `json:"title,omitempty"`
	Content      string         `json:"content,omitempty"`
	Source       string         `json:"source,omitempty"`
	Time         string         `json:"time,omitempty"`
	Rating       float64        `json:"rating,omitempty"`
	LabelIDs     []string       `json:"label_ids,omitempty"`
	IsFavorite   bool           `json:"is_favorite,omitempty"`
	SourceThread *SourceThread  `json:"source_thread,omitempty"`
	Confidence   float64        `json:"confidence,omitempty"`
	SpaceID      string         `json:"space_id,omitempty"`
	Metadata     map[string]any `json:"metadata,omitempty"`
}

// SourceThread is a lightweight thread reference attached to a memory.
type SourceThread struct {
	ID      string `json:"id"`
	Title   string `json:"title,omitempty"`
	SpaceID string `json:"space_id,omitempty"`
}

// Pagination holds pagination metadata.
type Pagination struct {
	Limit   int  `json:"limit"`
	Offset  int  `json:"offset"`
	Total   int  `json:"total"`
	HasMore bool `json:"has_more"`
}

// CreateMemoryRequest is the request body for POST /memories.
type CreateMemoryRequest struct {
	ID               *string        `json:"id,omitempty"`
	Content          string         `json:"content"`
	Title            *string        `json:"title,omitempty"`
	SourceThreadID   *string        `json:"source_thread_id,omitempty"`
	SourceMsgRange   map[string]int `json:"source_message_range,omitempty"`
	Source           *string        `json:"source,omitempty"`
	Importance       *float64       `json:"importance,omitempty"`
	Confidence       *float64       `json:"confidence,omitempty"`
	Labels           []string       `json:"labels,omitempty"`
	SpaceID          *string        `json:"space_id,omitempty"`
	Metadata         map[string]any `json:"metadata,omitempty"`
	EventStart       *string        `json:"event_start,omitempty"`
	EventEnd         *string        `json:"event_end,omitempty"`
	TemporalContext  *string        `json:"temporal_context,omitempty"`
	UnitType         *string        `json:"unit_type,omitempty"`
}

// CreateMemoryResponse is the response body for POST /memories.
type CreateMemoryResponse struct {
	Memory               Memory     `json:"memory"`
	ExtractedEntities    []Entity   `json:"extracted_entities,omitempty"`
	AssignedLabels       []string   `json:"assigned_labels,omitempty"`
	CreatedRelationships int        `json:"created_relationships,omitempty"`
	Action               string     `json:"action,omitempty"`
	Warnings             []string   `json:"warnings,omitempty"`
}

// ListMemoriesResponse is the response body for GET /memories.
type ListMemoriesResponse struct {
	Memories   []MemoryListItem `json:"memories"`
	Pagination Pagination       `json:"pagination"`
}

// ListMemoriesParams are query parameters for GET /memories.
type ListMemoriesParams struct {
	Limit         int     `json:"limit,omitempty"`
	Offset        int     `json:"offset,omitempty"`
	State         string  `json:"state,omitempty"`
	ImportanceMin float64 `json:"importance_min,omitempty"`
	SpaceID       string  `json:"space_id,omitempty"`
	IsCrystal     *bool   `json:"is_crystal,omitempty"`
}

// DeleteMemoryResponse is the response body for DELETE /memories/{id}.
type DeleteMemoryResponse struct {
	Message             string `json:"message"`
	DeletedRelationships int   `json:"deleted_relationships,omitempty"`
	DeletedEntities      int   `json:"deleted_entities,omitempty"`
}

// DeleteMemoryParams are query parameters for DELETE /memories/{id}.
type DeleteMemoryParams struct {
	CascadeDelete bool   `json:"cascade_delete,omitempty"`
	SpaceID       string `json:"space_id,omitempty"`
}

// --- Threads ---

// Thread represents a conversation thread.
type Thread struct {
	ID           string         `json:"id"`
	NodeType     string         `json:"node_type,omitempty"`
	CreatedAt    *time.Time     `json:"created_at,omitempty"`
	UpdatedAt    *time.Time     `json:"updated_at,omitempty"`
	Metadata     map[string]any `json:"metadata,omitempty"`
	ThreadID     string         `json:"thread_id"`
	Title        string         `json:"title,omitempty"`
	Summary      string         `json:"summary,omitempty"`
	MessageCount int            `json:"message_count,omitempty"`
	Participants []string       `json:"participants,omitempty"`
	Source       string         `json:"source,omitempty"`
	SpaceID      string         `json:"space_id,omitempty"`
	Project      string         `json:"project,omitempty"`
	Workspace    string         `json:"workspace,omitempty"`
	ToolVersion  string         `json:"tool_version,omitempty"`
	ImportDate   *time.Time     `json:"import_date,omitempty"`
}

// ThreadMessage represents a single message in a thread.
type ThreadMessage struct {
	ID         string         `json:"id"`
	NodeType   string         `json:"node_type,omitempty"`
	CreatedAt  *time.Time     `json:"created_at,omitempty"`
	UpdatedAt  *time.Time     `json:"updated_at,omitempty"`
	Metadata   map[string]any `json:"metadata,omitempty"`
	Content    string         `json:"content"`
	Role       string         `json:"role"`
	OrderIndex int            `json:"order_index,omitempty"`
	Timestamp  *time.Time     `json:"timestamp,omitempty"`
	TokenCount int            `json:"token_count,omitempty"`
}

// ThreadListItem is the summary view returned by list endpoints.
type ThreadListItem struct {
	ID         string `json:"id"`
	Title      string `json:"title,omitempty"`
	Summary    string `json:"summary,omitempty"`
	Source     string `json:"source,omitempty"`
	Messages   int    `json:"messages,omitempty"`
	Date       string `json:"date,omitempty"`
	IsFavorite bool   `json:"is_favorite,omitempty"`
	SpaceID    string `json:"space_id,omitempty"`
}

// MessageCreateRequest is a message in a create-thread request.
type MessageCreateRequest struct {
	Content string `json:"content"`
	Role    string `json:"role"`
}

// CreateThreadRequest is the request body for POST /threads.
type CreateThreadRequest struct {
	ThreadID    string                `json:"thread_id"`
	Title       *string               `json:"title,omitempty"`
	Messages    []MessageCreateRequest `json:"messages"`
	Participants []string             `json:"participants,omitempty"`
	Source      *string               `json:"source,omitempty"`
	SpaceID     string                `json:"space_id,omitempty"`
	Project     *string               `json:"project,omitempty"`
	Workspace   *string               `json:"workspace,omitempty"`
	ToolVersion *string               `json:"tool_version,omitempty"`
	ImportDate  *string               `json:"import_date,omitempty"`
	Metadata    map[string]any        `json:"metadata,omitempty"`
}

// CreateThreadResponse is the response body for POST /threads.
type CreateThreadResponse struct {
	Thread                 Thread          `json:"thread"`
	Messages               []ThreadMessage `json:"messages,omitempty"`
	CreatedRelationships   int             `json:"created_relationships,omitempty"`
	AutoGeneratedSummary   string          `json:"auto_generated_summary,omitempty"`
	ExtractedMemories      []Memory        `json:"extracted_memories,omitempty"`
	AutoExtractionPerformed bool           `json:"auto_extraction_performed,omitempty"`
}

// ListThreadsResponse is the response body for GET /threads.
type ListThreadsResponse struct {
	Threads    []ThreadListItem `json:"threads"`
	Pagination Pagination       `json:"pagination"`
}

// ListThreadsParams are query parameters for GET /threads.
type ListThreadsParams struct {
	Limit   int    `json:"limit,omitempty"`
	Offset  int    `json:"offset,omitempty"`
	Source  string `json:"source,omitempty"`
	SpaceID string `json:"space_id,omitempty"`
}

// --- Labels ---

// Label represents a label/tag.
type Label struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Color       string `json:"color,omitempty"`
	Description string `json:"description,omitempty"`
	CreatedAt   string `json:"created_at,omitempty"`
	UpdatedAt   string `json:"updated_at,omitempty"`
	UsageCount  int    `json:"usage_count,omitempty"`
}

// --- Entities ---

// Entity represents an entity node in the knowledge graph.
type Entity struct {
	ID                 string         `json:"id"`
	NodeType           string         `json:"node_type,omitempty"`
	CreatedAt          *time.Time     `json:"created_at,omitempty"`
	UpdatedAt          *time.Time     `json:"updated_at,omitempty"`
	Metadata           map[string]any `json:"metadata,omitempty"`
	Name               string         `json:"name"`
	EntityType         string         `json:"entity_type,omitempty"`
	Description        string         `json:"description,omitempty"`
	Aliases            []string       `json:"aliases,omitempty"`
	Confidence         float64        `json:"confidence,omitempty"`
	EntityCreated      string         `json:"entity_created,omitempty"`
	EntityEnded        string         `json:"entity_ended,omitempty"`
	TemporalPrecision  string         `json:"temporal_precision,omitempty"`
	TemporalConfidence float64        `json:"temporal_confidence,omitempty"`
	TemporalContext    string         `json:"temporal_context,omitempty"`
}

// EntityWithStats wraps an Entity with mention count.
type EntityWithStats struct {
	Entity      Entity `json:"entity"`
	MentionCount int   `json:"mention_count"`
}

// --- Spaces ---

// Space represents a space profile.
type Space struct {
	ID                   string       `json:"id"`
	Key                  string       `json:"key,omitempty"`
	Name                 string       `json:"name,omitempty"`
	Aliases              []string     `json:"aliases,omitempty"`
	Description          string       `json:"description,omitempty"`
	Icon                 string       `json:"icon,omitempty"`
	Instructions         string       `json:"instructions,omitempty"`
	SharedSpaceIDs       []string     `json:"sharedSpaceIds,omitempty"`
	DefaultRetrievalMode string       `json:"defaultRetrievalMode,omitempty"`
	Usage                *SpaceUsage  `json:"usage,omitempty"`
	Observed             bool         `json:"observed,omitempty"`
	HasProfile           bool         `json:"hasProfile,omitempty"`
}

// SpaceUsage holds usage statistics for a space.
type SpaceUsage struct {
	Memories        int  `json:"memories"`
	Threads         int  `json:"threads"`
	Sources         int  `json:"sources"`
	HasWorkingMemory bool `json:"hasWorkingMemory"`
}

// ListSpacesResponse is the response body for GET /spaces.
type ListSpacesResponse struct {
	Enabled bool    `json:"enabled"`
	Spaces  []Space `json:"spaces"`
}

// --- Sources ---

// Source represents a library source (file, URL, etc.).
type Source struct {
	ID             string         `json:"id"`
	SourceType     string         `json:"source_type,omitempty"`
	OriginalName   string         `json:"original_name,omitempty"`
	MimeType       string         `json:"mime_type,omitempty"`
	FilePath       string         `json:"file_path,omitempty"`
	ParsedPath     string         `json:"parsed_path,omitempty"`
	SourceURL      string         `json:"source_url,omitempty"`
	SHA256         string         `json:"sha256,omitempty"`
	SizeBytes      int64          `json:"size_bytes,omitempty"`
	Version        int            `json:"version,omitempty"`
	SpaceID        string         `json:"space_id,omitempty"`
	LifecycleState string         `json:"lifecycle_state,omitempty"`
	ChunkCount     int            `json:"chunk_count,omitempty"`
	MemoryCount    int            `json:"memory_count,omitempty"`
	SectionTree    string         `json:"section_tree,omitempty"`
	Summary        string         `json:"summary,omitempty"`
	ErrorMessage   string         `json:"error_message,omitempty"`
	CreatedAt      string         `json:"created_at,omitempty"`
	UpdatedAt      string         `json:"updated_at,omitempty"`
	Metadata       map[string]any `json:"metadata,omitempty"`
	LabelIDs       []string       `json:"label_ids,omitempty"`
}

// ListSourcesResponse is the response body for GET /sources.
type ListSourcesResponse struct {
	Sources []Source `json:"sources"`
	Total   int      `json:"total"`
}

// ListSourcesParams are query parameters for GET /sources.
type ListSourcesParams struct {
	Limit          int    `json:"limit,omitempty"`
	Offset         int    `json:"offset,omitempty"`
	SourceType     string `json:"source_type,omitempty"`
	LifecycleState string `json:"lifecycle_state,omitempty"`
	SpaceID        string `json:"space_id,omitempty"`
}

// --- Health ---

// HealthCheck is the response body for GET /health.
type HealthCheck struct {
	Status                       string            `json:"status"`
	Version                      string            `json:"version"`
	Timestamp                    *time.Time        `json:"timestamp,omitempty"`
	DatabaseConnected            bool              `json:"database_connected"`
	ServicesReady                bool              `json:"services_ready"`
	BufferPoolExhausted          bool              `json:"buffer_pool_exhausted"`
	BufferPoolAutoEscalatedMB    int               `json:"buffer_pool_auto_escalated_mb"`
	BufferPoolCurrentMB          int               `json:"buffer_pool_current_mb"`
	BufferPoolNextStartMB        int               `json:"buffer_pool_next_start_mb"`
	BufferPoolRestartRequired    bool              `json:"buffer_pool_restart_required"`
	BufferPoolRestartReason      string            `json:"buffer_pool_restart_reason"`
	BufferPoolAutoMode           bool              `json:"buffer_pool_auto_mode"`
	BufferPoolSource             string            `json:"buffer_pool_source"`
	BufferPoolAutoFloorMB        int               `json:"buffer_pool_auto_floor_mb"`
	BufferPoolAutoCapMB          int               `json:"buffer_pool_auto_cap_mb"`
	PluginUpdates                []PluginUpdate    `json:"plugin_updates,omitempty"`
	ContentStore                 *ContentStoreInfo `json:"content_store,omitempty"`
}

// PluginUpdate represents an available plugin update.
type PluginUpdate struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	InstalledVersion string `json:"installed_version"`
	AvailableVersion string `json:"available_version"`
}

// ContentStoreInfo holds content store status.
type ContentStoreInfo struct {
	State                        string `json:"state"`
	DBPath                       string `json:"db_path"`
	SQLiteReady                  bool   `json:"sqlite_ready"`
	SchemaVersion                int    `json:"schema_version"`
	SchemaMinSupportedVersion    int    `json:"schema_min_supported_version"`
	SchemaMigrationCount         int    `json:"schema_migration_count"`
	ThreadMessageOwner           string `json:"thread_message_owner"`
	CutoverReady                 bool   `json:"cutover_ready"`
	CutoverCompleted             bool   `json:"cutover_completed"`
	LegacyGraphCleanupCompleted  bool   `json:"legacy_graph_cleanup_completed"`
	SQLiteMessageCount           int    `json:"sqlite_message_count"`
	SQLiteThreadCount            int    `json:"sqlite_thread_count"`
	SQLiteAnchorCount            int    `json:"sqlite_anchor_count"`
	LegacyKuzuMessageCount       int    `json:"legacy_kuzu_message_count"`
	LastError                    string `json:"last_error"`
	UpdatedAt                    string `json:"updated_at"`
}

// --- FS ---

// FSEntry represents a file/directory in the Nowledge FS tree.
type FSEntry struct {
	Path  string `json:"path"`
	Name  string `json:"name"`
	Type  string `json:"type"`
	Size  int64  `json:"size,omitempty"`
	IsDir bool   `json:"is_dir,omitempty"`
}

// FSListResponse is the response body for GET /fs/ls.
type FSListResponse struct {
	Entries  []FSEntry `json:"entries"`
	Cursor   string    `json:"cursor,omitempty"`
	HasMore  bool      `json:"has_more,omitempty"`
}

// FSCatResponse is the response body for GET /fs/cat.
type FSCatResponse struct {
	Path       string         `json:"path"`
	Body       string         `json:"body"`
	Frontmatter map[string]any `json:"frontmatter,omitempty"`
	TotalLines int            `json:"total_lines,omitempty"`
}

// FSStatResponse is the response body for GET /fs/stat.
type FSStatResponse struct {
	Path         string         `json:"path"`
	Type         string         `json:"type"`
	Size         int64          `json:"size"`
	CreatedAt    string         `json:"created_at,omitempty"`
	UpdatedAt    string         `json:"updated_at,omitempty"`
	Metadata     map[string]any `json:"metadata,omitempty"`
}

// FSWriteRequest is the request body for POST /fs/write.
type FSWriteRequest struct {
	Path string `json:"path"`
	Body string `json:"body"`
}

// FSDeleteRequest is the request body for POST /fs/delete.
type FSDeleteRequest struct {
	Path string `json:"path"`
}

// FSSearchResponse is the response body for GET /fs/find, /fs/recall.
type FSSearchResponse struct {
	Paths []string `json:"paths"`
}

// FSGrepMatch represents a single grep match.
type FSGrepMatch struct {
	Path string `json:"path"`
	Line int    `json:"line"`
	Text string `json:"text"`
}

// FSGrepResponse is the response body for GET /fs/grep.
type FSGrepResponse struct {
	Matches []FSGrepMatch `json:"matches"`
}

// --- Agent ---

// AgentStatus is the response body for GET /agent/status.
type AgentStatus struct {
	Running         bool   `json:"running"`
	CurrentTask     string `json:"current_task,omitempty"`
	LastRunAt       string `json:"last_run_at,omitempty"`
	NextScheduledAt string `json:"next_scheduled_at,omitempty"`
}

// --- Graph ---

// GraphAnalysis is the response body for GET /graph/analysis.
type GraphAnalysis struct {
	NodeCount       int              `json:"node_count"`
	EdgeCount       int              `json:"edge_count"`
	CommunityCount  int              `json:"community_count"`
	Communities     []Community      `json:"communities,omitempty"`
	CentralityScores map[string]float64 `json:"centrality_scores,omitempty"`
}

// Community represents a knowledge community.
type Community struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description,omitempty"`
	Size        int      `json:"size"`
	SampleMemIDs []string `json:"sample_mem_ids,omitempty"`
}

// --- Error ---

// APIError represents an error response from the API.
type APIError struct {
	StatusCode int           `json:"-"`
	Status     string        `json:"-"`
	Body       string        `json:"-"`
	Detail     []ErrorDetail `json:"detail"`
}

// ErrorDetail is a single validation error.
type ErrorDetail struct {
	Loc  []string       `json:"loc"`
	Msg  string         `json:"msg"`
	Type string         `json:"type"`
	Input any           `json:"input,omitempty"`
	Ctx  map[string]any `json:"ctx,omitempty"`
}

func (e *APIError) Error() string {
	if len(e.Detail) == 0 {
		if e.Body != "" {
			return fmt.Sprintf("API error (status %d): %s", e.StatusCode, e.Body)
		}
		if e.Status != "" {
			return fmt.Sprintf("API error: %s", e.Status)
		}
		return "unknown API error"
	}
	if len(e.Detail) == 1 {
		return e.Detail[0].Msg
	}
	msgs := make([]string, len(e.Detail))
	for i, d := range e.Detail {
		msgs[i] = d.Msg
	}
	return fmt.Sprintf("%s (and %d more errors)", msgs[0], len(msgs)-1)
}
