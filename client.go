package onledgemem

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	defaultBaseURL = "http://127.0.0.1:14242"
	defaultTimeout = 30 * time.Second
)

// Client is the Nowledge Mem API client.
//
// Create a client with NewClient:
//
//	client := onledgemem.NewClient()
//	client := onledgemem.NewClient(onledgemem.WithBaseURL("http://host:14242"))
type Client struct {
	baseURL    *url.URL
	httpClient *http.Client

	// Services provides access to API resource operations.
	Memories       *MemoriesService
	Threads        *ThreadsService
	Spaces         *SpacesService
	Labels         *LabelsService
	Entities       *EntitiesService
	Sources        *SourcesService
	Health         *HealthService
	FS             *FSService
	Agent          *AgentService
	Graph          *GraphService
	GraphVis       *GraphVisService
	Distillation   *DistillationService
	KG             *KGService
	Communities    *CommunitiesService
	Data           *DataService
	Storage        *StorageService
	Settings       *SettingsService
	Models         *ModelsService
	SearchIndex    *SearchIndexService
	Embeddings     *EmbeddingsService
	Feed           *FeedService
	WorkingMemory  *WorkingMemoryService
	Library        *LibraryService
	Capabilities   *CapabilitiesService
	Admin          *AdminService
	Favorites      *FavoritesService
	ContentStore   *ContentStoreService
}

// Option configures the client.
type Option func(*Client)

// WithBaseURL overrides the default base URL.
func WithBaseURL(rawURL string) Option {
	return func(c *Client) {
		u, err := url.Parse(rawURL)
		if err != nil {
			panic(fmt.Sprintf("invalid base URL %q: %v", rawURL, err))
		}
		c.baseURL = u
	}
}

// WithHTTPClient overrides the default HTTP client.
func WithHTTPClient(hc *http.Client) Option {
	return func(c *Client) {
		c.httpClient = hc
	}
}

// WithTimeout overrides the default HTTP timeout.
func WithTimeout(d time.Duration) Option {
	return func(c *Client) {
		c.httpClient.Timeout = d
	}
}

// NewClient creates a new Nowledge Mem API client.
func NewClient(opts ...Option) *Client {
	u, _ := url.Parse(defaultBaseURL)
	c := &Client{
		baseURL: u,
		httpClient: &http.Client{
			Timeout: defaultTimeout,
		},
	}
	for _, opt := range opts {
		opt(c)
	}
	c.Memories = &MemoriesService{client: c}
	c.Threads = &ThreadsService{client: c}
	c.Spaces = &SpacesService{client: c}
	c.Labels = &LabelsService{client: c}
	c.Entities = &EntitiesService{client: c}
	c.Sources = &SourcesService{client: c}
	c.Health = &HealthService{client: c}
	c.FS = &FSService{client: c}
	c.Agent = &AgentService{client: c}
	c.Graph = &GraphService{client: c}
	c.GraphVis = &GraphVisService{client: c}
	c.Distillation = &DistillationService{client: c}
	c.KG = &KGService{client: c}
	c.Communities = &CommunitiesService{client: c}
	c.Data = &DataService{client: c}
	c.Storage = &StorageService{client: c}
	c.Settings = &SettingsService{client: c}
	c.Models = &ModelsService{client: c}
	c.SearchIndex = &SearchIndexService{client: c}
	c.Embeddings = &EmbeddingsService{client: c}
	c.Feed = &FeedService{client: c}
	c.WorkingMemory = &WorkingMemoryService{client: c}
	c.Library = &LibraryService{client: c}
	c.Capabilities = &CapabilitiesService{client: c}
	c.Admin = &AdminService{client: c}
	c.Favorites = &FavoritesService{client: c}
	c.ContentStore = &ContentStoreService{client: c}
	return c
}

// Close closes idle HTTP connections. Call this when done with the client.
func (c *Client) Close() {
	c.httpClient.CloseIdleConnections()
}

// BaseURL returns the client's base URL.
func (c *Client) BaseURL() *url.URL {
	return c.baseURL
}

// do executes an HTTP request and decodes the response into dst.
// If dst is nil, the response body is discarded.
func (c *Client) do(ctx context.Context, method, path string, body any, dst any) error {
	u := c.baseURL.ResolveReference(&url.URL{Path: path})

	var reqBody io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("marshal request body: %w", err)
		}
		reqBody = bytes.NewReader(data)
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), reqBody)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return c.doRequest(req, dst)
}

// doQuery is like do but builds a URL with query parameters.
func (c *Client) doQuery(ctx context.Context, path string, params url.Values, dst any) error {
	u := c.baseURL.ResolveReference(&url.URL{Path: path})
	u.RawQuery = params.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	return c.doRequest(req, dst)
}

// doRequest executes an HTTP request and decodes the response.
func (c *Client) doRequest(req *http.Request, dst any) error {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		var apiErr APIError
		if err := json.NewDecoder(resp.Body).Decode(&apiErr); err != nil {
			return fmt.Errorf("API error (status %d)", resp.StatusCode)
		}
		return &apiErr
	}

	if dst != nil {
		if err := json.NewDecoder(resp.Body).Decode(dst); err != nil {
			return fmt.Errorf("decode response: %w", err)
		}
	}
	return nil
}
