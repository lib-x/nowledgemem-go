package nowledgemem

import (
	"context"
	"net/url"
	"strconv"
)

// WorkingMemoryService handles working memory operations.
//
// It provides methods for reading, writing, and browsing the Working Memory file.
type WorkingMemoryService struct {
	client *Client
}

// Get reads the Working Memory file (today's or an archived day).
//
// GET /agent/working-memory
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
//
// PUT /agent/working-memory
func (s *WorkingMemoryService) Update(ctx context.Context, req *UpdateWorkingMemoryRequest) error {
	return s.client.do(ctx, "PUT", "/agent/working-memory", req, nil)
}

// History lists dates with archived Working Memory files.
//
// GET /agent/working-memory/history
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

// WorkingMemory is the response from Get (GET /agent/working-memory).
type WorkingMemory struct {
	Date    string `json:"date"`
	Content string `json:"content"`
}

// UpdateWorkingMemoryRequest is the request for Update (PUT /agent/working-memory).
type UpdateWorkingMemoryRequest struct {
	Content string `json:"content"`
	SpaceID string `json:"space_id,omitempty"`
}
