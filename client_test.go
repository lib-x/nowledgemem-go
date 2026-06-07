package nowledgemem

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
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
