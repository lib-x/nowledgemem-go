package nowledgemem

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	defaultBaseURL = "http://127.0.0.1:14242"
	defaultTimeout = 30 * time.Second
	envAPIURL      = "NMEM_API_URL"
	envAPIKey      = "NMEM_API_KEY"
)

// Client is the Nowledge Mem API client.
//
// Create a client with NewClient:
//
//	client := nowledgemem.NewClient()
//	client := nowledgemem.NewClient(nowledgemem.WithBaseURL("http://host:14242"))
type Client struct {
	baseURL      *url.URL
	httpClient   *http.Client
	bearerToken  string
	headerAPIKey string
	queryAPIKey  string

	// Services provides access to API resource operations.
	Memories      *MemoriesService
	Threads       *ThreadsService
	Spaces        *SpacesService
	Labels        *LabelsService
	Entities      *EntitiesService
	Sources       *SourcesService
	Health        *HealthService
	FS            *FSService
	Agent         *AgentService
	Graph         *GraphService
	GraphVis      *GraphVisService
	Distillation  *DistillationService
	KG            *KGService
	Communities   *CommunitiesService
	Events        *EventsService
	Data          *DataService
	Storage       *StorageService
	Settings      *SettingsService
	Models        *ModelsService
	SearchIndex   *SearchIndexService
	Embeddings    *EmbeddingsService
	Feed          *FeedService
	WorkingMemory *WorkingMemoryService
	Library       *LibraryService
	Capabilities  *CapabilitiesService
	Admin         *AdminService
	Favorites     *FavoritesService
	ContentStore  *ContentStoreService
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

// WithBearerToken sets an Authorization: Bearer token header for every request.
//
// Pass the raw token value, for example "nmem_xxxx". If the value already starts
// with "Bearer ", the prefix is stripped and normalized.
func WithBearerToken(token string) Option {
	return func(c *Client) {
		token = normalizeBearerToken(token)
		if token == "" {
			panic("bearer token cannot be empty")
		}
		c.bearerToken = token
	}
}

// WithAPIKey sets the Nowledge Mem remote API key for every request.
//
// This sends both supported header forms: Authorization: Bearer nmem_xxxx and
// X-NMEM-API-Key: nmem_xxxx.
func WithAPIKey(apiKey string) Option {
	return func(c *Client) {
		apiKey = normalizeBearerToken(apiKey)
		if apiKey == "" {
			panic("API key cannot be empty")
		}
		c.bearerToken = apiKey
		c.headerAPIKey = apiKey
	}
}

// WithAPIKeyQuery sends nmem_api_key=nmem_xxxx on every request.
//
// Prefer header authentication when possible. Use this for proxies or clients
// that strip custom headers.
func WithAPIKeyQuery(apiKey string) Option {
	return func(c *Client) {
		apiKey = normalizeBearerToken(apiKey)
		if apiKey == "" {
			panic("query API key cannot be empty")
		}
		c.queryAPIKey = apiKey
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
	c.Events = &EventsService{client: c}
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

// NewRemoteClient creates a client for a remote Nowledge Mem deployment.
//
// Remote deployments use the backend API URL, such as "https://mem.example.com",
// and an nmem API key. Do not append the web app's frontend-only /remote-api
// route. The key is sent as both Authorization: Bearer nmem_xxxx and
// X-NMEM-API-Key: nmem_xxxx.
func NewRemoteClient(rawURL, apiKey string, opts ...Option) *Client {
	remoteOpts := make([]Option, 0, len(opts)+2)
	remoteOpts = append(remoteOpts, WithBaseURL(rawURL), WithAPIKey(apiKey))
	remoteOpts = append(remoteOpts, opts...)
	return NewClient(remoteOpts...)
}

// NewClientFromEnv creates a client from NMEM_API_URL and NMEM_API_KEY.
//
// Explicit options are applied after environment-derived options.
func NewClientFromEnv(opts ...Option) *Client {
	envOpts := optionsFromClientConfig(os.Getenv(envAPIURL), os.Getenv(envAPIKey))
	envOpts = append(envOpts, opts...)
	return NewClient(envOpts...)
}

// ClientConfig is the shared local client configuration written by nmem.
type ClientConfig struct {
	APIURL string `json:"apiUrl"`
	APIKey string `json:"apiKey"`
}

// NewClientFromConfig creates a client from ~/.nowledge-mem/config.json, with
// NMEM_API_URL and NMEM_API_KEY overriding file values when present.
//
// Explicit options are applied last.
func NewClientFromConfig(opts ...Option) (*Client, error) {
	cfg, err := readClientConfig()
	if err != nil {
		return nil, err
	}

	apiURL := cfg.APIURL
	apiKey := cfg.APIKey
	if envURL := os.Getenv(envAPIURL); envURL != "" {
		apiURL = envURL
	}
	if envKey := os.Getenv(envAPIKey); envKey != "" {
		apiKey = envKey
	}

	configOpts := optionsFromClientConfig(apiURL, apiKey)
	configOpts = append(configOpts, opts...)
	return NewClient(configOpts...), nil
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
	return c.doWithQuery(ctx, method, path, nil, body, dst)
}

// doWithQuery executes an HTTP request with optional query parameters and JSON body.
func (c *Client) doWithQuery(ctx context.Context, method, path string, params url.Values, body any, dst any) error {
	var reqBody io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("marshal request body: %w", err)
		}
		reqBody = bytes.NewReader(data)
	}

	req, err := c.newRequest(ctx, method, path, params, reqBody, "application/json")
	if err != nil {
		return err
	}

	return c.doRequest(req, dst)
}

func (c *Client) newRequest(ctx context.Context, method, path string, params url.Values, body io.Reader, contentType string) (*http.Request, error) {
	u, err := c.requestURL(path, params)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), body)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	if contentType != "" && body != nil {
		req.Header.Set("Content-Type", contentType)
	}
	if c.bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.bearerToken)
	}
	if c.headerAPIKey != "" {
		req.Header.Set("X-NMEM-API-Key", c.headerAPIKey)
	}
	return req, nil
}

func (c *Client) requestURL(path string, params url.Values) (*url.URL, error) {
	rel, err := url.Parse(path)
	if err != nil {
		return nil, fmt.Errorf("parse request path: %w", err)
	}
	if rel.IsAbs() || rel.Host != "" {
		return nil, fmt.Errorf("request path must be relative: %s", path)
	}

	u := *c.baseURL
	u.Path = joinURLPath(c.baseURL.Path, rel.Path)
	u.RawPath = ""
	q := u.Query()
	for key, values := range rel.Query() {
		for _, value := range values {
			q.Add(key, value)
		}
	}
	if params != nil {
		for key, values := range params {
			for _, value := range values {
				q.Add(key, value)
			}
		}
	}
	if c.queryAPIKey != "" {
		q.Set("nmem_api_key", c.queryAPIKey)
	}
	u.RawQuery = q.Encode()
	return &u, nil
}

func normalizeBearerToken(token string) string {
	token = strings.TrimSpace(token)
	if strings.HasPrefix(strings.ToLower(token), "bearer ") {
		token = strings.TrimSpace(token[len("Bearer "):])
	}
	return token
}

func joinURLPath(basePath, requestPath string) string {
	if basePath == "" || basePath == "/" {
		if strings.HasPrefix(requestPath, "/") {
			return requestPath
		}
		return "/" + requestPath
	}
	if requestPath == "" || requestPath == "/" {
		return basePath
	}
	return strings.TrimRight(basePath, "/") + "/" + strings.TrimLeft(requestPath, "/")
}

func optionsFromClientConfig(apiURL, apiKey string) []Option {
	opts := make([]Option, 0, 2)
	if apiURL != "" {
		opts = append(opts, WithBaseURL(apiURL))
	}
	if apiKey != "" {
		opts = append(opts, WithAPIKey(apiKey))
	}
	return opts
}

func readClientConfig() (*ClientConfig, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return &ClientConfig{}, nil
	}
	path := filepath.Join(home, ".nowledge-mem", "config.json")
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &ClientConfig{}, nil
		}
		return nil, fmt.Errorf("read client config: %w", err)
	}
	if len(data) == 0 {
		return &ClientConfig{}, nil
	}

	var cfg ClientConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("decode client config: %w", err)
	}
	return &cfg, nil
}

func (c *Client) doBytes(ctx context.Context, method, path string, params url.Values, body any) ([]byte, error) {
	var reqBody io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshal request body: %w", err)
		}
		reqBody = bytes.NewReader(data)
	}

	req, err := c.newRequest(ctx, method, path, params, reqBody, "application/json")
	if err != nil {
		return nil, err
	}
	return c.doBytesRequest(req)
}

func (c *Client) doStream(ctx context.Context, method, path string, params url.Values, body any) (*http.Response, error) {
	var reqBody io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshal request body: %w", err)
		}
		reqBody = bytes.NewReader(data)
	}

	req, err := c.newRequest(ctx, method, path, params, reqBody, "application/json")
	if err != nil {
		return nil, err
	}
	return c.doStreamRequest(req)
}

// doQuery is like do but builds a URL with query parameters.
func (c *Client) doQuery(ctx context.Context, path string, params url.Values, dst any) error {
	return c.doWithQuery(ctx, http.MethodGet, path, params, nil, dst)
}

// doRequest executes an HTTP request and decodes the response.
func (c *Client) doRequest(req *http.Request, dst any) error {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close()

	if err := apiErrorFromResponse(resp); err != nil {
		return err
	}

	if dst != nil && resp.StatusCode != http.StatusNoContent {
		if err := json.NewDecoder(resp.Body).Decode(dst); err != nil {
			return fmt.Errorf("decode response: %w", err)
		}
	}
	return nil
}

func (c *Client) doBytesRequest(req *http.Request) ([]byte, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close()

	if err := apiErrorFromResponse(resp); err != nil {
		return nil, err
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}
	return data, nil
}

func (c *Client) doStreamRequest(req *http.Request) (*http.Response, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("execute request: %w", err)
	}

	if err := apiErrorFromResponse(resp); err != nil {
		resp.Body.Close()
		return nil, err
	}
	return resp, nil
}

func apiErrorFromResponse(resp *http.Response) error {
	if resp.StatusCode < 400 {
		return nil
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read API error (status %d): %w", resp.StatusCode, err)
	}

	apiErr := &APIError{
		StatusCode: resp.StatusCode,
		Status:     resp.Status,
		Body:       strings.TrimSpace(string(data)),
	}
	if len(data) > 0 {
		_ = json.Unmarshal(data, apiErr)
	}
	return apiErr
}
