package nowledgemem

import "context"

// EmbeddingsService handles OpenAI-compatible embedding operations.
type EmbeddingsService struct {
	client *Client
}

// ListModels lists available models.
func (s *EmbeddingsService) ListModels(ctx context.Context) ([]EmbeddingModel, error) {
	var resp struct {
		Data []EmbeddingModel `json:"data"`
	}
	if err := s.client.do(ctx, "GET", "/v1/models", nil, &resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// CreateEmbeddings generates embeddings using the local model.
func (s *EmbeddingsService) CreateEmbeddings(ctx context.Context, req *CreateEmbeddingsRequest) (*CreateEmbeddingsResponse, error) {
	var resp CreateEmbeddingsResponse
	if err := s.client.do(ctx, "POST", "/v1/embeddings", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// --- Embeddings types ---

// EmbeddingModel represents an available embedding model.
type EmbeddingModel struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	OwnedBy string `json:"owned_by"`
}

// CreateEmbeddingsRequest is the request for POST /v1/embeddings.
type CreateEmbeddingsRequest struct {
	Model          string   `json:"model,omitempty"`
	Input          any      `json:"input"`
	EncodingFormat string   `json:"encoding_format,omitempty"`
	Dimensions     *int     `json:"dimensions,omitempty"`
	User           string   `json:"user,omitempty"`
	InputType      string   `json:"input_type,omitempty"`
}

// CreateEmbeddingsResponse is the response for POST /v1/embeddings.
type CreateEmbeddingsResponse struct {
	Object string          `json:"object"`
	Data   []EmbeddingData `json:"data"`
	Model  string          `json:"model"`
	Usage  EmbeddingUsage  `json:"usage"`
}

// EmbeddingData is a single embedding result.
type EmbeddingData struct {
	Object    string    `json:"object"`
	Embedding []float64 `json:"embedding"`
	Index     int       `json:"index"`
}

// EmbeddingUsage tracks token usage.
type EmbeddingUsage struct {
	PromptTokens int `json:"prompt_tokens"`
	TotalTokens  int `json:"total_tokens"`
}
