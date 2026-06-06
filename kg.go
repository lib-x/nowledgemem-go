package nowledgemem

import (
	"context"
	"fmt"
	"net/url"
)

// KGService handles knowledge graph extraction from memories.
type KGService struct {
	client *Client
}

// PreviewExtraction previews KG extraction for a memory before applying.
func (s *KGService) PreviewExtraction(ctx context.Context, memoryID string) (*KGPreviewResponse, error) {
	var resp KGPreviewResponse
	path := fmt.Sprintf("/memories/%s/extract-kg/preview", url.PathEscape(memoryID))
	if err := s.client.do(ctx, "POST", path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ApplyExtraction saves extracted entities and relationships to the graph database.
func (s *KGService) ApplyExtraction(ctx context.Context, memoryID string, req *KGApplyRequest) (*KGApplyResponse, error) {
	var resp KGApplyResponse
	path := fmt.Sprintf("/memories/%s/extract-kg/apply", url.PathEscape(memoryID))
	if err := s.client.do(ctx, "POST", path, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// --- KG types ---

// KGPreviewResponse is the response for POST /memories/{id}/extract-kg/preview.
type KGPreviewResponse struct {
	Entities      []Entity     `json:"entities"`
	Relationships []KGRelation `json:"relationships"`
}

// KGRelation represents a relationship between entities.
type KGRelation struct {
	SourceID   string  `json:"source_id"`
	TargetID   string  `json:"target_id"`
	RelType    string  `json:"rel_type"`
	Confidence float64 `json:"confidence,omitempty"`
}

// KGApplyRequest is the request for POST /memories/{id}/extract-kg/apply.
type KGApplyRequest struct {
	Entities      []Entity     `json:"entities,omitempty"`
	Relationships []KGRelation `json:"relationships,omitempty"`
}

// KGApplyResponse is the response for POST /memories/{id}/extract-kg/apply.
type KGApplyResponse struct {
	CreatedEntities      int `json:"created_entities"`
	CreatedRelationships int `json:"created_relationships"`
}
