package nowledgemem

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestDoWithQueryKeepsQueryOutOfPath(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/memories/m1" {
			t.Fatalf("path = %q, want /memories/m1", r.URL.Path)
		}
		if got := r.URL.Query().Get("cascade_delete"); got != "true" {
			t.Fatalf("cascade_delete = %q, want true", got)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	client := NewClient(WithBaseURL(srv.URL))
	params := url.Values{"cascade_delete": {"true"}}
	if err := client.doWithQuery(context.Background(), http.MethodDelete, "/memories/m1", params, nil, nil); err != nil {
		t.Fatalf("doWithQuery returned error: %v", err)
	}
}

func TestAPIErrorPreservesPlainTextBody(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "plain failure", http.StatusInternalServerError)
	}))
	defer srv.Close()

	client := NewClient(WithBaseURL(srv.URL))
	err := client.do(context.Background(), http.MethodGet, "/health", nil, nil)
	if err == nil {
		t.Fatal("expected error")
	}

	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("error type = %T, want *APIError", err)
	}
	if apiErr.StatusCode != http.StatusInternalServerError {
		t.Fatalf("status = %d, want 500", apiErr.StatusCode)
	}
	if !strings.Contains(apiErr.Error(), "plain failure") {
		t.Fatalf("error message = %q, want plain body", apiErr.Error())
	}
}

func TestWithBearerTokenSetsAuthorizationHeaderOnly(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Authorization"); got != "Bearer nmem_header" {
			t.Fatalf("authorization = %q, want Bearer nmem_header", got)
		}
		if got := r.URL.Query().Get("nmem_api_key"); got != "" {
			t.Fatalf("nmem_api_key = %q, want empty", got)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	client := NewClient(WithBaseURL(srv.URL), WithBearerToken("Bearer nmem_header"))
	if err := client.do(context.Background(), http.MethodGet, "/health", nil, nil); err != nil {
		t.Fatalf("do returned error: %v", err)
	}
}

func TestWithAPIKeySetsBearerAndAPIKeyHeaders(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Authorization"); got != "Bearer nmem_key" {
			t.Fatalf("authorization = %q, want Bearer nmem_key", got)
		}
		if got := r.Header.Get("X-NMEM-API-Key"); got != "nmem_key" {
			t.Fatalf("X-NMEM-API-Key = %q, want nmem_key", got)
		}
		if got := r.URL.Query().Get("nmem_api_key"); got != "" {
			t.Fatalf("nmem_api_key = %q, want empty", got)
		}
		if got := r.URL.Query().Get("existing"); got != "1" {
			t.Fatalf("existing = %q, want 1", got)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	client := NewClient(WithBaseURL(srv.URL), WithAPIKey("nmem_key"))
	if err := client.doQuery(context.Background(), "/health", url.Values{"existing": {"1"}}, nil); err != nil {
		t.Fatalf("doQuery returned error: %v", err)
	}
}

func TestWithAPIKeyQuerySetsQueryParamOnly(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Authorization"); got != "" {
			t.Fatalf("authorization = %q, want empty", got)
		}
		if got := r.Header.Get("X-NMEM-API-Key"); got != "" {
			t.Fatalf("X-NMEM-API-Key = %q, want empty", got)
		}
		if got := r.URL.Query().Get("nmem_api_key"); got != "nmem_key" {
			t.Fatalf("nmem_api_key = %q, want nmem_key", got)
		}
		if got := r.URL.Query().Get("existing"); got != "1" {
			t.Fatalf("existing = %q, want 1", got)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	client := NewClient(WithBaseURL(srv.URL), WithAPIKeyQuery("Bearer nmem_key"))
	if err := client.doQuery(context.Background(), "/health", url.Values{"existing": {"1"}}, nil); err != nil {
		t.Fatalf("doQuery returned error: %v", err)
	}
}

func TestNewRemoteClientUsesBackendAPIPath(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/models/bge-m3/status" {
			t.Fatalf("path = %q, want /models/bge-m3/status", r.URL.Path)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer nmem_key" {
			t.Fatalf("authorization = %q, want Bearer nmem_key", got)
		}
		if got := r.Header.Get("X-NMEM-API-Key"); got != "nmem_key" {
			t.Fatalf("X-NMEM-API-Key = %q, want nmem_key", got)
		}
		if got := r.URL.Query().Get("nmem_api_key"); got != "" {
			t.Fatalf("nmem_api_key = %q, want empty", got)
		}
		_ = json.NewEncoder(w).Encode(EmbeddingModelStatus{Installed: true, Ready: true})
	}))
	defer srv.Close()

	client := NewRemoteClient(srv.URL, "Bearer nmem_key")
	status, err := client.Models.GetEmbeddingModelStatus(context.Background())
	if err != nil {
		t.Fatalf("GetEmbeddingModelStatus returned error: %v", err)
	}
	if !status.Ready {
		t.Fatal("ready = false, want true")
	}
}

func TestNewClientFromEnvReadsURLAndAPIKey(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/health" {
			t.Fatalf("path = %q, want /health", r.URL.Path)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer nmem_env" {
			t.Fatalf("authorization = %q, want Bearer nmem_env", got)
		}
		if got := r.Header.Get("X-NMEM-API-Key"); got != "nmem_env" {
			t.Fatalf("X-NMEM-API-Key = %q, want nmem_env", got)
		}
		_ = json.NewEncoder(w).Encode(HealthCheck{Status: "ok"})
	}))
	defer srv.Close()

	t.Setenv("NMEM_API_URL", srv.URL)
	t.Setenv("NMEM_API_KEY", "nmem_env")

	client := NewClientFromEnv()
	health, err := client.Health.Check(context.Background())
	if err != nil {
		t.Fatalf("Health.Check returned error: %v", err)
	}
	if health.Status != "ok" {
		t.Fatalf("status = %q, want ok", health.Status)
	}
}

func TestNewClientFromConfigReadsSharedClientConfig(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/spaces/roster" {
			t.Fatalf("path = %q, want /spaces/roster", r.URL.Path)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer nmem_config" {
			t.Fatalf("authorization = %q, want Bearer nmem_config", got)
		}
		if got := r.Header.Get("X-NMEM-API-Key"); got != "nmem_config" {
			t.Fatalf("X-NMEM-API-Key = %q, want nmem_config", got)
		}
		_ = json.NewEncoder(w).Encode(ListSpacesResponse{})
	}))
	defer srv.Close()

	home := t.TempDir()
	configDir := filepath.Join(home, ".nowledge-mem")
	if err := os.MkdirAll(configDir, 0o700); err != nil {
		t.Fatalf("mkdir config dir: %v", err)
	}
	config := []byte(`{"apiUrl":"` + srv.URL + `","apiKey":"nmem_config"}`)
	if err := os.WriteFile(filepath.Join(configDir, "config.json"), config, 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}

	t.Setenv("HOME", home)
	t.Setenv("NMEM_API_URL", "")
	t.Setenv("NMEM_API_KEY", "")

	client, err := NewClientFromConfig()
	if err != nil {
		t.Fatalf("NewClientFromConfig returned error: %v", err)
	}
	if _, err := client.Spaces.Roster(context.Background()); err != nil {
		t.Fatalf("Spaces.Roster returned error: %v", err)
	}
}

func TestNewClientFromConfigEnvOverridesFile(t *testing.T) {
	fileSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("file config server should not be called when env overrides are set")
	}))
	defer fileSrv.Close()

	envSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/health" {
			t.Fatalf("path = %q, want /health", r.URL.Path)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer nmem_env" {
			t.Fatalf("authorization = %q, want Bearer nmem_env", got)
		}
		if got := r.Header.Get("X-NMEM-API-Key"); got != "nmem_env" {
			t.Fatalf("X-NMEM-API-Key = %q, want nmem_env", got)
		}
		_ = json.NewEncoder(w).Encode(HealthCheck{Status: "ok"})
	}))
	defer envSrv.Close()

	home := t.TempDir()
	configDir := filepath.Join(home, ".nowledge-mem")
	if err := os.MkdirAll(configDir, 0o700); err != nil {
		t.Fatalf("mkdir config dir: %v", err)
	}
	config := []byte(`{"apiUrl":"` + fileSrv.URL + `","apiKey":"nmem_file"}`)
	if err := os.WriteFile(filepath.Join(configDir, "config.json"), config, 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}

	t.Setenv("HOME", home)
	t.Setenv("NMEM_API_URL", envSrv.URL)
	t.Setenv("NMEM_API_KEY", "nmem_env")

	client, err := NewClientFromConfig()
	if err != nil {
		t.Fatalf("NewClientFromConfig returned error: %v", err)
	}
	health, err := client.Health.Check(context.Background())
	if err != nil {
		t.Fatalf("Health.Check returned error: %v", err)
	}
	if health.Status != "ok" {
		t.Fatalf("status = %q, want ok", health.Status)
	}
}

func TestRemoteAPIKeyIntegration(t *testing.T) {
	apiKey := os.Getenv("NMEM_TEST_API_KEY")
	if apiKey == "" {
		t.Skip("set NMEM_TEST_API_KEY to run remote API integration test")
	}
	baseURL := os.Getenv("NMEM_TEST_BASE_URL")
	if baseURL == "" {
		t.Skip("set NMEM_TEST_BASE_URL to the backend API URL, for example https://mem.example.com")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	client := NewRemoteClient(baseURL, apiKey, WithTimeout(20*time.Second))
	status, err := client.Models.GetEmbeddingModelStatus(ctx)
	if err != nil {
		t.Fatalf("GetEmbeddingModelStatus returned error: %v", err)
	}
	if status == nil {
		t.Fatal("status is nil")
	}
}

func TestSourcesIngestFileUsesMultipart(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("method = %s, want POST", r.Method)
		}
		if r.URL.Path != "/sources/ingest/file" {
			t.Fatalf("path = %q, want /sources/ingest/file", r.URL.Path)
		}
		if !strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data") {
			t.Fatalf("content-type = %q, want multipart/form-data", r.Header.Get("Content-Type"))
		}

		if err := r.ParseMultipartForm(1 << 20); err != nil {
			t.Fatalf("parse multipart: %v", err)
		}
		file, header, err := r.FormFile("file")
		if err != nil {
			t.Fatalf("form file: %v", err)
		}
		defer file.Close()
		data, err := io.ReadAll(file)
		if err != nil {
			t.Fatalf("read file: %v", err)
		}
		if header.Filename != "note.md" {
			t.Fatalf("filename = %q, want note.md", header.Filename)
		}
		if string(data) != "hello" {
			t.Fatalf("file body = %q, want hello", string(data))
		}
		if got := r.FormValue("space_id"); got != "work" {
			t.Fatalf("space_id = %q, want work", got)
		}
		if got := r.FormValue("labels"); got != "label_a,label_b" {
			t.Fatalf("labels = %q, want label_a,label_b", got)
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(IngestSourceResponse{
			SourceID:       "src1",
			OriginalName:   "note.md",
			LifecycleState: "ready",
		})
	}))
	defer srv.Close()

	client := NewClient(WithBaseURL(srv.URL))
	resp, err := client.Sources.IngestFile(context.Background(), &IngestFileRequest{
		File:     strings.NewReader("hello"),
		Filename: "note.md",
		Labels:   "label_a,label_b",
		SpaceID:  "work",
	})
	if err != nil {
		t.Fatalf("IngestFile returned error: %v", err)
	}
	if resp.SourceID != "src1" {
		t.Fatalf("source_id = %q, want src1", resp.SourceID)
	}
}

func TestEventsStreamReturnsOpenBody(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/events/stream" {
			t.Fatalf("path = %q, want /events/stream", r.URL.Path)
		}
		w.Header().Set("Content-Type", "text/event-stream")
		_, _ = w.Write([]byte("data: ok\n\n"))
	}))
	defer srv.Close()

	client := NewClient(WithBaseURL(srv.URL))
	resp, err := client.Events.Stream(context.Background())
	if err != nil {
		t.Fatalf("Stream returned error: %v", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read stream: %v", err)
	}
	if string(data) != "data: ok\n\n" {
		t.Fatalf("stream body = %q, want event payload", string(data))
	}
}
