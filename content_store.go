package onledgemem

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
)

// ContentStoreService handles content store migration operations.
type ContentStoreService struct {
	client *Client
}

// GetMigrationStatus returns the current migration status.
func (s *ContentStoreService) GetMigrationStatus(ctx context.Context) (*ContentStoreMigrationStatus, error) {
	var resp ContentStoreMigrationStatus
	if err := s.client.do(ctx, "GET", "/content-store/migration/status", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// CopyAll copies all legacy thread messages into SQLite.
func (s *ContentStoreService) CopyAll(ctx context.Context) error {
	return s.client.do(ctx, "POST", "/content-store/thread-messages/copy-all", nil, nil)
}

// CopyBatch copies one bounded legacy Kuzu Message page into SQLite.
func (s *ContentStoreService) CopyBatch(ctx context.Context, batchSize int) error {
	q := url.Values{}
	if batchSize > 0 {
		q.Set("batch_size", strconv.Itoa(batchSize))
	}
	u := s.client.baseURL.ResolveReference(&url.URL{Path: "/content-store/thread-messages/copy-batch"})
	u.RawQuery = q.Encode()
	req, err := http.NewRequestWithContext(ctx, "POST", u.String(), nil)
	if err != nil {
		return err
	}
	return s.client.doRequest(req, nil)
}

// Cutover performs the content store cutover.
func (s *ContentStoreService) Cutover(ctx context.Context) error {
	return s.client.do(ctx, "POST", "/content-store/thread-messages/cutover", nil, nil)
}

// MigrateAnchors migrates thread message anchors.
func (s *ContentStoreService) MigrateAnchors(ctx context.Context) error {
	return s.client.do(ctx, "POST", "/content-store/thread-messages/migrate-anchors", nil, nil)
}

// MigrateThroughCutover migrates thread messages through cutover.
func (s *ContentStoreService) MigrateThroughCutover(ctx context.Context) error {
	return s.client.do(ctx, "POST", "/content-store/thread-messages/migrate-through-cutover", nil, nil)
}

// Verify verifies the content store migration.
func (s *ContentStoreService) Verify(ctx context.Context) error {
	return s.client.do(ctx, "POST", "/content-store/thread-messages/verify", nil, nil)
}

// CleanupLegacyGraph cleans up the legacy thread message graph.
func (s *ContentStoreService) CleanupLegacyGraph(ctx context.Context) error {
	return s.client.do(ctx, "POST", "/content-store/thread-messages/cleanup-legacy-graph", nil, nil)
}

// --- Content Store types ---

// ContentStoreMigrationStatus is the response for GET /content-store/migration/status.
type ContentStoreMigrationStatus struct {
	State                       string `json:"state"`
	DBPath                      string `json:"db_path"`
	SQLiteReady                 bool   `json:"sqlite_ready"`
	SchemaVersion               int    `json:"schema_version"`
	SchemaMinSupportedVersion   int    `json:"schema_min_supported_version"`
	SchemaMigrationCount        int    `json:"schema_migration_count"`
	ThreadMessageOwner          string `json:"thread_message_owner"`
	CutoverReady                bool   `json:"cutover_ready"`
	CutoverCompleted            bool   `json:"cutover_completed"`
	LegacyGraphCleanupCompleted bool   `json:"legacy_graph_cleanup_completed"`
	SQLiteMessageCount          int    `json:"sqlite_message_count"`
	SQLiteThreadCount           int    `json:"sqlite_thread_count"`
	SQLiteAnchorCount           int    `json:"sqlite_anchor_count"`
	LegacyKuzuMessageCount      int    `json:"legacy_kuzu_message_count"`
	LastError                   string `json:"last_error"`
	UpdatedAt                   string `json:"updated_at"`
}
