package onledgemem

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

// CommunitiesService handles community operations.
type CommunitiesService struct {
	client *Client
}

// List returns knowledge communities with AI summaries.
func (s *CommunitiesService) List(ctx context.Context, limit int) ([]Community, error) {
	q := url.Values{}
	if limit > 0 {
		q.Set("limit", strconv.Itoa(limit))
	}
	var resp []Community
	if err := s.client.doQuery(ctx, "/communities", q, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// Get returns community details including entities and sample memories.
func (s *CommunitiesService) Get(ctx context.Context, communityID string) (*CommunityDetail, error) {
	var resp CommunityDetail
	path := fmt.Sprintf("/communities/%s", url.PathEscape(communityID))
	if err := s.client.do(ctx, "GET", path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetRecentMemories returns recent memories in a community.
func (s *CommunitiesService) GetRecentMemories(ctx context.Context, communityID string, limit int) ([]MemoryListItem, error) {
	q := url.Values{}
	if limit > 0 {
		q.Set("limit", strconv.Itoa(limit))
	}
	var resp []MemoryListItem
	path := fmt.Sprintf("/library/community/%s/recent-memories", url.PathEscape(communityID))
	if err := s.client.doQuery(ctx, path, q, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// GetRelated returns related communities.
func (s *CommunitiesService) GetRelated(ctx context.Context, communityID string, limit int) ([]Community, error) {
	q := url.Values{}
	if limit > 0 {
		q.Set("limit", strconv.Itoa(limit))
	}
	var resp []Community
	path := fmt.Sprintf("/library/community/%s/related", url.PathEscape(communityID))
	if err := s.client.doQuery(ctx, path, q, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// GetSubgraph returns the subgraph for a community.
func (s *CommunitiesService) GetSubgraph(ctx context.Context, communityID string) (*GraphData, error) {
	var resp GraphData
	path := fmt.Sprintf("/library/community/%s/subgraph", url.PathEscape(communityID))
	if err := s.client.do(ctx, "GET", path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// --- Community types ---

// CommunityDetail holds detailed community info.
type CommunityDetail struct {
	ID             string           `json:"id"`
	Name           string           `json:"name"`
	Description    string           `json:"description,omitempty"`
	Size           int              `json:"size"`
	SampleMemIDs   []string         `json:"sample_mem_ids,omitempty"`
	Entities       []Entity         `json:"entities,omitempty"`
	SampleMemories []MemoryListItem `json:"sample_memories,omitempty"`
}
