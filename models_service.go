package onledgemem

import (
	"context"
	"fmt"
	"net/url"
)

// ModelsService handles embedding model operations.
type ModelsService struct {
	client *Client
}

// GetEmbeddingModelStatus checks the search embedding model status.
func (s *ModelsService) GetEmbeddingModelStatus(ctx context.Context) (*EmbeddingModelStatus, error) {
	var resp EmbeddingModelStatus
	if err := s.client.do(ctx, "GET", "/models/bge-m3/status", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// InstallEmbeddingModel downloads and installs the search embedding model.
func (s *ModelsService) InstallEmbeddingModel(ctx context.Context) error {
	return s.client.do(ctx, "POST", "/models/bge-m3/install", nil, nil)
}

// GetMemoryStatus returns which models are currently loaded in memory.
func (s *ModelsService) GetMemoryStatus(ctx context.Context) (*ModelMemoryStatus, error) {
	var resp ModelMemoryStatus
	if err := s.client.do(ctx, "GET", "/models/memory-status", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// UnloadModel manually unloads a model from memory.
func (s *ModelsService) UnloadModel(ctx context.Context, modelType string) error {
	path := fmt.Sprintf("/models/unload/%s", url.PathEscape(modelType))
	return s.client.do(ctx, "POST", path, nil, nil)
}

// --- Models types ---

// EmbeddingModelStatus is the response for GET /models/bge-m3/status.
type EmbeddingModelStatus struct {
	Installed bool   `json:"installed"`
	Ready     bool   `json:"ready"`
	ModelPath string `json:"model_path,omitempty"`
	Version   string `json:"version,omitempty"`
}

// ModelMemoryStatus is the response for GET /models/memory-status.
type ModelMemoryStatus struct {
	Loaded []LoadedModel `json:"loaded"`
}

// LoadedModel represents a model loaded in memory.
type LoadedModel struct {
	ModelType  string `json:"model_type"`
	Device     string `json:"device"`
	MemoryMB   int    `json:"memory_mb"`
}
