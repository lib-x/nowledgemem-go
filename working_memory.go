package nowledgemem

import (
	"context"
	"net/url"
	"strconv"
)

// WorkingMemoryService handles working memory operations.
type WorkingMemoryService struct {
	client *Client
}

// Get reads the Working Memory file (today's or an archived day).
func (s *WorkingMemoryService) Get(ctx context.Context, date string, spaceID string) (*WorkingMemory, error) {
	q := url.Values{}
	if date != "" {
		q.Set("date", date)
	}
	if spaceID != "" {
		q.Set("space_id", spaceID)
	}
	var resp WorkingMemory
	if err := s.client.doQuery(ctx, "/agent/working-memory", q, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Update writes the Working Memory file from user edits.
func (s *WorkingMemoryService) Update(ctx context.Context, req *UpdateWorkingMemoryRequest) error {
	return s.client.do(ctx, "PUT", "/agent/working-memory", req, nil)
}

// History lists dates with archived Working Memory files.
func (s *WorkingMemoryService) History(ctx context.Context, limit int, spaceID string) ([]string, error) {
	q := url.Values{}
	if limit > 0 {
		q.Set("limit", strconv.Itoa(limit))
	}
	if spaceID != "" {
		q.Set("space_id", spaceID)
	}
	var resp []string
	if err := s.client.doQuery(ctx, "/agent/working-memory/history", q, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// --- Working Memory types ---

// WorkingMemory is the response for GET /agent/working-memory.
type WorkingMemory struct {
	Date    string `json:"date"`
	Content string `json:"content"`
}

// UpdateWorkingMemoryRequest is the request for PUT /agent/working-memory.
type UpdateWorkingMemoryRequest struct {
	Content string `json:"content"`
	SpaceID string `json:"space_id,omitempty"`
}
