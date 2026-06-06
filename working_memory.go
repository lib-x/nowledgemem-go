package nowledgemem

import (
	"context"
	"net/url"
)

// WorkingMemoryService handles working memory operations.
type WorkingMemoryService struct {
	client *Client
}

// Get reads the Working Memory file (today's or an archived day).
func (s *WorkingMemoryService) Get(ctx context.Context, date string) (*WorkingMemory, error) {
	q := url.Values{}
	if date != "" {
		q.Set("date", date)
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
func (s *WorkingMemoryService) History(ctx context.Context) ([]string, error) {
	var resp []string
	if err := s.client.do(ctx, "GET", "/agent/working-memory/history", nil, &resp); err != nil {
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
}
