package nowledgemem

import (
	"context"
	"fmt"
	"net/url"
)

// KGService provides methods for knowledge graph extraction from memories.
type KGService struct {
	client *Client
}

// PreviewExtraction previews knowledge graph extraction for a memory before applying. (POST /memories/{id}/extract-kg/preview)
func (s *KGService) PreviewExtraction(ctx context.Context, memoryID string, req *KGPreviewRequest) (*KGPreviewResponse, error) {
	var resp KGPreviewResponse
	path := fmt.Sprintf("/memories/%s/extract-kg/preview", url.PathEscape(memoryID))
	if err := s.client.do(ctx, "POST", path, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ApplyExtraction saves extracted entities and relationships to the graph database. (POST /memories/{id}/extract-kg/apply)
func (s *KGService) ApplyExtraction(ctx context.Context, memoryID string, req *KGApplyRequest) (*KGApplyResponse, error) {
	var resp KGApplyResponse
	path := fmt.Sprintf("/memories/%s/extract-kg/apply", url.PathEscape(memoryID))
	if err := s.client.do(ctx, "POST", path, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// --- KG types ---

// KGPreviewRequest is the request for POST /memories/{id}/extract-kg/preview.
type KGPreviewRequest struct {
	ForceReextraction bool   `json:"force_reextraction,omitempty"`
	ExtractionLevel   string `json:"extraction_level,omitempty"`
	UseRemoteLLM      bool   `json:"use_remote_llm,omitempty"`
	PreferredLanguage string `json:"preferred_language,omitempty"`
}

// KGPreviewResponse is the response for POST /memories/{id}/extract-kg/preview.
type KGPreviewResponse struct {
	MemoryID            string       `json:"memory_id"`
	MemoryTitle         string       `json:"memory_title"`
	MemoryContent       string       `json:"memory_content"`
	Entities            []Entity     `json:"entities"`
	Relationships       []KGRelation `json:"relationships"`
	ExtractionConfidence float64     `json:"extraction_confidence"`
	EntitiesCount       int          `json:"entities_count"`
	RelationshipsCount  int          `json:"relationships_count"`
	KGAlreadyExtracted  bool         `json:"kg_already_extracted"`
	CanExtract          bool         `json:"can_extract"`
}

// KGRelation represents a relationship between two entities in the knowledge graph.
type KGRelation struct {
	SourceID   string  `json:"source_id"`
	TargetID   string  `json:"target_id"`
	RelType    string  `json:"rel_type"`
	Confidence float64 `json:"confidence,omitempty"`
}

// KGApplyRequest is the request for POST /memories/{id}/extract-kg/apply.
type KGApplyRequest struct {
	Entities              []Entity     `json:"entities,omitempty"`
	Relationships         []KGRelation `json:"relationships,omitempty"`
	ExtractionConfidence  *float64     `json:"extraction_confidence,omitempty"`
}

// KGApplyResponse is the response for POST /memories/{id}/extract-kg/apply.
type KGApplyResponse struct {
	Success             bool   `json:"success"`
	MemoryID            string `json:"memory_id,omitempty"`
	EntitiesCreated     int    `json:"entities_created"`
	RelationshipsCreated int   `json:"relationships_created"`
	MetadataUpdated     bool   `json:"metadata_updated,omitempty"`
	Error               string `json:"error,omitempty"`
}
