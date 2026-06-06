package nowledgemem

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

// EntitiesService handles entity operations.
type EntitiesService struct {
	client *Client
}

// List returns entities with optional filtering.
func (s *EntitiesService) List(ctx context.Context, params *ListEntitiesParams) ([]Entity, error) {
	q := url.Values{}
	if params != nil {
		if params.Limit > 0 {
			q.Set("limit", strconv.Itoa(params.Limit))
		}
		if params.EntityType != "" {
			q.Set("entity_type", params.EntityType)
		}
		if params.IncludeStats {
			q.Set("include_stats", "true")
		}
	}
	var resp []Entity
	if err := s.client.doQuery(ctx, "/entities", q, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// ListWithStats returns entities sorted by mention count with stats.
func (s *EntitiesService) ListWithStats(ctx context.Context, params *ListEntitiesParams) ([]EntityWithStats, error) {
	if params == nil {
		params = &ListEntitiesParams{}
	}
	params.IncludeStats = true

	q := url.Values{}
	if params.Limit > 0 {
		q.Set("limit", strconv.Itoa(params.Limit))
	}
	if params.EntityType != "" {
		q.Set("entity_type", params.EntityType)
	}
	q.Set("include_stats", "true")

	var resp []EntityWithStats
	if err := s.client.doQuery(ctx, "/entities", q, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// GetRelationships returns all connected entities and memories for an entity.
func (s *EntitiesService) GetRelationships(ctx context.Context, entityID string) (*EntityRelationships, error) {
	var resp EntityRelationships
	path := fmt.Sprintf("/entities/%s/relationships", url.PathEscape(entityID))
	if err := s.client.do(ctx, "GET", path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// --- Entity types ---

// ListEntitiesParams are query parameters for GET /entities.
type ListEntitiesParams struct {
	Limit         int    `json:"limit,omitempty"`
	EntityType    string `json:"entity_type,omitempty"`
	IncludeStats  bool   `json:"include_stats,omitempty"`
}

// EntityRelationships holds an entity's connected nodes.
type EntityRelationships struct {
	Entity          Entity           `json:"entity"`
	RelatedEntities []Entity         `json:"related_entities,omitempty"`
	RelatedMemories []MemoryListItem `json:"related_memories,omitempty"`
}
