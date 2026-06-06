package onledgemem

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

// FeedService handles feed event operations.
type FeedService struct {
	client *Client
}

// GetEvents returns feed events with filtering.
func (s *FeedService) GetEvents(ctx context.Context, params *FeedEventsParams) ([]FeedEvent, error) {
	q := url.Values{}
	if params != nil {
		if params.Limit > 0 {
			q.Set("limit", strconv.Itoa(params.Limit))
		}
		if params.Offset > 0 {
			q.Set("offset", strconv.Itoa(params.Offset))
		}
		if params.Severity != "" {
			q.Set("severity", params.Severity)
		}
		if params.EventType != "" {
			q.Set("event_type", params.EventType)
		}
		if params.UnresolvedOnly {
			q.Set("unresolved_only", "true")
		}
		if params.LastNDays > 0 {
			q.Set("last_n_days", strconv.Itoa(params.LastNDays))
		}
		if params.DateFrom != "" {
			q.Set("date_from", params.DateFrom)
		}
		if params.DateTo != "" {
			q.Set("date_to", params.DateTo)
		}
		if params.Source != "" {
			q.Set("source", params.Source)
		}
	}
	var resp []FeedEvent
	if err := s.client.doQuery(ctx, "/agent/feed/events", q, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// ResolveEvent resolves an action-required event.
func (s *FeedService) ResolveEvent(ctx context.Context, eventID string, req *ResolveEventRequest) error {
	path := fmt.Sprintf("/agent/feed/events/%s/resolve", url.PathEscape(eventID))
	return s.client.do(ctx, "POST", path, req, nil)
}

// RetryEvent retries a failed background task.
func (s *FeedService) RetryEvent(ctx context.Context, eventID string) error {
	path := fmt.Sprintf("/agent/feed/events/%s/retry", url.PathEscape(eventID))
	return s.client.do(ctx, "POST", path, nil, nil)
}

// DeleteEvent soft-deletes a feed event.
func (s *FeedService) DeleteEvent(ctx context.Context, eventID string) error {
	path := fmt.Sprintf("/agent/feed/events/%s", url.PathEscape(eventID))
	return s.client.do(ctx, "DELETE", path, nil, nil)
}

// PersistQuestion persists a question and agent response as a feed event.
func (s *FeedService) PersistQuestion(ctx context.Context, req *PersistQuestionRequest) error {
	return s.client.do(ctx, "POST", "/agent/feed/input/persist-question", req, nil)
}

// --- Feed types ---

// FeedEventsParams are parameters for GET /agent/feed/events.
type FeedEventsParams struct {
	Limit          int    `json:"limit,omitempty"`
	Offset         int    `json:"offset,omitempty"`
	Severity       string `json:"severity,omitempty"`
	EventType      string `json:"event_type,omitempty"`
	UnresolvedOnly bool   `json:"unresolved_only,omitempty"`
	LastNDays      int    `json:"last_n_days,omitempty"`
	DateFrom       string `json:"date_from,omitempty"`
	DateTo         string `json:"date_to,omitempty"`
	Source         string `json:"source,omitempty"`
}

// FeedEvent represents a feed event.
type FeedEvent struct {
	ID          string         `json:"id"`
	EventType   string         `json:"event_type"`
	Severity    string         `json:"severity"`
	Title       string         `json:"title"`
	Body        string         `json:"body,omitempty"`
	Resolved    bool           `json:"resolved"`
	CreatedAt   string         `json:"created_at,omitempty"`
	Metadata    map[string]any `json:"metadata,omitempty"`
}

// ResolveEventRequest is the request for POST /agent/feed/events/{id}/resolve.
type ResolveEventRequest struct {
	Resolution string         `json:"resolution,omitempty"`
	GraphMutations map[string]any `json:"graph_mutations,omitempty"`
}

// PersistQuestionRequest is the request for POST /agent/feed/input/persist-question.
type PersistQuestionRequest struct {
	Question string `json:"question"`
	Response string `json:"response"`
	Source   string `json:"source,omitempty"`
}
