package onledgemem

import (
	"context"
	"fmt"
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

// BulkDelete deletes multiple threads at once.
func (s *ThreadsService) BulkDelete(ctx context.Context, threadIDs []string) (*BulkDeleteResponse, error) {
	var resp BulkDeleteResponse
	body := map[string]any{"thread_ids": threadIDs}
	if err := s.client.do(ctx, "DELETE", "/threads/bulk/delete", body, &resp); err != nil {
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
