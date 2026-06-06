package nowledgemem

import "context"

// FavoritesService handles favorites operations.
type FavoritesService struct {
	client *Client
}

// GetFavoriteMemories returns all favorite memories.
func (s *FavoritesService) GetFavoriteMemories(ctx context.Context) ([]MemoryListItem, error) {
	var resp []MemoryListItem
	if err := s.client.do(ctx, "GET", "/favorites/memories", nil, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// GetFavoriteThreads returns all favorite threads.
func (s *FavoritesService) GetFavoriteThreads(ctx context.Context) ([]ThreadListItem, error) {
	var resp []ThreadListItem
	if err := s.client.do(ctx, "GET", "/favorites/threads", nil, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}
