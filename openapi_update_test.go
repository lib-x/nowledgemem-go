package nowledgemem

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestOpenAPIServicesAreInitialized(t *testing.T) {
	t.Parallel()

	client := NewClient()
	if client.Skills == nil {
		t.Fatal("Skills service is nil")
	}
	if client.Context == nil {
		t.Fatal("Context service is nil")
	}
}

func TestThreadsDeleteUsesDeleteMethod(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Fatalf("method = %s, want DELETE", r.Method)
		}
		if r.URL.Path != "/threads/thread-1" {
			t.Fatalf("path = %q, want /threads/thread-1", r.URL.Path)
		}
		if got := r.URL.Query().Get("cascade_delete_memories"); got != "true" {
			t.Fatalf("cascade_delete_memories = %q, want true", got)
		}
		_ = json.NewEncoder(w).Encode(DeleteThreadResponse{Message: "deleted"})
	}))
	defer srv.Close()

	client := NewClient(WithBaseURL(srv.URL))
	resp, err := client.Threads.Delete(context.Background(), "thread-1", &DeleteThreadParams{
		CascadeDeleteMemories: true,
	})
	if err != nil {
		t.Fatalf("Threads.Delete returned error: %v", err)
	}
	if resp.Message != "deleted" {
		t.Fatalf("message = %q, want deleted", resp.Message)
	}
}

func TestThreadsDeleteBulkUsesOpenAPIEndpoint(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Fatalf("method = %s, want DELETE", r.Method)
		}
		if r.URL.Path != "/threads/bulk" {
			t.Fatalf("path = %q, want /threads/bulk", r.URL.Path)
		}
		if got := r.URL.Query().Get("space_id"); got != "work" {
			t.Fatalf("space_id = %q, want work", got)
		}
		var body BulkDeleteThreadsRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if len(body.ThreadIDs) != 1 || body.ThreadIDs[0] != "thread-1" {
			t.Fatalf("thread_ids = %#v, want [thread-1]", body.ThreadIDs)
		}
		_ = json.NewEncoder(w).Encode(ThreadBulkDeleteResponse{Message: "deleted", DeletedCount: 1})
	}))
	defer srv.Close()

	client := NewClient(WithBaseURL(srv.URL))
	resp, err := client.Threads.DeleteBulk(
		context.Background(),
		&BulkDeleteThreadsRequest{ThreadIDs: []string{"thread-1"}},
		&BulkDeleteThreadsParams{SpaceID: "work"},
	)
	if err != nil {
		t.Fatalf("Threads.DeleteBulk returned error: %v", err)
	}
	if resp.DeletedCount != 1 {
		t.Fatalf("deleted_count = %d, want 1", resp.DeletedCount)
	}
}

func TestMemoryRelationRoutes(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("method = %s, want POST", r.Method)
		}
		if r.URL.Path != "/memories/mem-1/relations" {
			t.Fatalf("path = %q, want /memories/mem-1/relations", r.URL.Path)
		}
		var body MemoryRelationCreateRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if body.TargetMemoryID != "mem-2" || body.RelationType != "supports" {
			t.Fatalf("body = %#v, want target mem-2 and supports", body)
		}
		_ = json.NewEncoder(w).Encode(MemoryRelation{
			ID:             "rel-1",
			SourceMemoryID: "mem-1",
			SourceSpaceID:  "default",
			TargetMemoryID: "mem-2",
			TargetSpaceID:  "default",
			RelationType:   "supports",
		})
	}))
	defer srv.Close()

	client := NewClient(WithBaseURL(srv.URL))
	resp, err := client.Memories.CreateRelation(context.Background(), "mem-1", &MemoryRelationCreateRequest{
		TargetMemoryID: "mem-2",
		RelationType:   "supports",
	})
	if err != nil {
		t.Fatalf("Memories.CreateRelation returned error: %v", err)
	}
	if resp.ID != "rel-1" {
		t.Fatalf("relation ID = %q, want rel-1", resp.ID)
	}
}

func TestSettingsAgentProfileRouteAndJSON(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Fatalf("method = %s, want PUT", r.Method)
		}
		if r.URL.Path != "/settings/agent-profiles/agent-1" {
			t.Fatalf("path = %q, want /settings/agent-profiles/agent-1", r.URL.Path)
		}
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if got := body["displayName"]; got != "Agent One" {
			t.Fatalf("displayName = %#v, want Agent One", got)
		}
		_ = json.NewEncoder(w).Encode(AgentProfileResponse{ID: "agent-1", DisplayName: "Agent One"})
	}))
	defer srv.Close()

	name := "Agent One"
	client := NewClient(WithBaseURL(srv.URL))
	resp, err := client.Settings.UpdateAgentProfile(context.Background(), "agent-1", &AgentProfilePayload{
		DisplayName: &name,
	})
	if err != nil {
		t.Fatalf("Settings.UpdateAgentProfile returned error: %v", err)
	}
	if resp.DisplayName != "Agent One" {
		t.Fatalf("displayName = %q, want Agent One", resp.DisplayName)
	}
}

func TestSkillsRoutes(t *testing.T) {
	t.Parallel()

	var calls int
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		switch calls {
		case 1:
			if r.Method != http.MethodGet || r.URL.Path != "/skills/match" {
				t.Fatalf("first request = %s %s, want GET /skills/match", r.Method, r.URL.Path)
			}
			if got := r.URL.Query().Get("q"); got != "golang" {
				t.Fatalf("q = %q, want golang", got)
			}
			if got := r.URL.Query().Get("limit"); got != "3" {
				t.Fatalf("limit = %q, want 3", got)
			}
		case 2:
			if r.Method != http.MethodPost || r.URL.Path != "/skills/skill-1/activate" {
				t.Fatalf("second request = %s %s, want POST /skills/skill-1/activate", r.Method, r.URL.Path)
			}
		default:
			t.Fatalf("unexpected request %d: %s %s", calls, r.Method, r.URL.Path)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"ok": true})
	}))
	defer srv.Close()

	client := NewClient(WithBaseURL(srv.URL))
	if _, err := client.Skills.Match(context.Background(), "golang", 3); err != nil {
		t.Fatalf("Skills.Match returned error: %v", err)
	}
	if _, err := client.Skills.Activate(context.Background(), "skill-1"); err != nil {
		t.Fatalf("Skills.Activate returned error: %v", err)
	}
	if calls != 2 {
		t.Fatalf("calls = %d, want 2", calls)
	}
}

func TestAgentPlanAndSkillBuilderStreamRoutes(t *testing.T) {
	t.Parallel()

	var calls int
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		switch calls {
		case 1:
			if r.Method != http.MethodGet || r.URL.Path != "/agent/trigger/community-detection/plan" {
				t.Fatalf("first request = %s %s, want GET /agent/trigger/community-detection/plan", r.Method, r.URL.Path)
			}
			if got := r.URL.Query().Get("scan_limit"); got != "50" {
				t.Fatalf("scan_limit = %q, want 50", got)
			}
			_ = json.NewEncoder(w).Encode(map[string]any{"planned": true})
		case 2:
			if r.Method != http.MethodPost || r.URL.Path != "/agent/skill-builder/refine/stream" {
				t.Fatalf("second request = %s %s, want POST /agent/skill-builder/refine/stream", r.Method, r.URL.Path)
			}
			body, err := io.ReadAll(r.Body)
			if err != nil {
				t.Fatalf("read body: %v", err)
			}
			if !strings.Contains(string(body), `"skill_id":"skill-1"`) {
				t.Fatalf("body = %q, want skill_id", string(body))
			}
			w.Header().Set("Content-Type", "text/event-stream")
			_, _ = w.Write([]byte("data: ok\n\n"))
		default:
			t.Fatalf("unexpected request %d: %s %s", calls, r.Method, r.URL.Path)
		}
	}))
	defer srv.Close()

	client := NewClient(WithBaseURL(srv.URL))
	if _, err := client.Agent.PlanCommunityDetection(context.Background(), &CommunityDetectionPlanParams{ScanLimit: 50}); err != nil {
		t.Fatalf("Agent.PlanCommunityDetection returned error: %v", err)
	}
	resp, err := client.Agent.SkillBuilderRefineStream(context.Background(), &SkillRefineRequest{SkillID: "skill-1"})
	if err != nil {
		t.Fatalf("Agent.SkillBuilderRefineStream returned error: %v", err)
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

func TestSourcesIngestContentRoute(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("method = %s, want POST", r.Method)
		}
		if r.URL.Path != "/sources/ingest/content" {
			t.Fatalf("path = %q, want /sources/ingest/content", r.URL.Path)
		}
		var body IngestContentRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if body.Name != "note.md" || body.Content != "hello" {
			t.Fatalf("body = %#v, want note.md hello", body)
		}
		_ = json.NewEncoder(w).Encode(IngestSourceResponse{
			SourceID:       "src-1",
			OriginalName:   "note.md",
			LifecycleState: "ready",
		})
	}))
	defer srv.Close()

	client := NewClient(WithBaseURL(srv.URL))
	resp, err := client.Sources.IngestContent(context.Background(), &IngestContentRequest{
		Name:    "note.md",
		Content: "hello",
	})
	if err != nil {
		t.Fatalf("Sources.IngestContent returned error: %v", err)
	}
	if resp.SourceID != "src-1" {
		t.Fatalf("source_id = %q, want src-1", resp.SourceID)
	}
}

func TestDataExportAsyncAndStatusRoute(t *testing.T) {
	t.Parallel()

	var calls int
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		switch calls {
		case 1:
			if r.Method != http.MethodPost || r.URL.Path != "/data/export" {
				t.Fatalf("first request = %s %s, want POST /data/export", r.Method, r.URL.Path)
			}
			var body DataExportRequest
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Fatalf("decode body: %v", err)
			}
			if body.ExportPath != "/tmp/export.zip" || !body.Async {
				t.Fatalf("body = %#v, want async export path", body)
			}
			_ = json.NewEncoder(w).Encode(DataExportResponse{
				Success: true,
				JobID:   "job-1",
				Status:  "queued",
				Message: "Export job started",
			})
		case 2:
			if r.Method != http.MethodGet || r.URL.Path != "/data/export/status/job-1" {
				t.Fatalf("second request = %s %s, want GET /data/export/status/job-1", r.Method, r.URL.Path)
			}
			kind := "export"
			_ = json.NewEncoder(w).Encode(DataTransferStatus{
				JobID:  "job-1",
				Status: "completed",
				Kind:   &kind,
				Result: map[string]any{"export_path": "/tmp/export.zip"},
			})
		default:
			t.Fatalf("unexpected request %d: %s %s", calls, r.Method, r.URL.Path)
		}
	}))
	defer srv.Close()

	client := NewClient(WithBaseURL(srv.URL))
	start, err := client.Data.Export(context.Background(), &DataExportRequest{
		ExportPath: "/tmp/export.zip",
		Async:      true,
	})
	if err != nil {
		t.Fatalf("Data.Export returned error: %v", err)
	}
	if !start.Success || start.JobID != "job-1" || start.Status != "queued" {
		t.Fatalf("start = %#v, want async job response", start)
	}

	status, err := client.Data.ExportStatus(context.Background(), start.JobID)
	if err != nil {
		t.Fatalf("Data.ExportStatus returned error: %v", err)
	}
	if status.Status != "completed" || status.Kind == nil || *status.Kind != "export" {
		t.Fatalf("status = %#v, want completed export status", status)
	}
	if got := status.Result["export_path"]; got != "/tmp/export.zip" {
		t.Fatalf("result export_path = %#v, want /tmp/export.zip", got)
	}
	if calls != 2 {
		t.Fatalf("calls = %d, want 2", calls)
	}
}

func TestDataCheckpointResultRoute(t *testing.T) {
	t.Parallel()

	detail := "ok"
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/data/checkpoint" {
			t.Fatalf("request = %s %s, want POST /data/checkpoint", r.Method, r.URL.Path)
		}
		_ = json.NewEncoder(w).Encode(DataTransferCheckpointResponse{
			Success:      true,
			Checkpointed: true,
			Detail:       &detail,
		})
	}))
	defer srv.Close()

	client := NewClient(WithBaseURL(srv.URL))
	resp, err := client.Data.CheckpointResult(context.Background())
	if err != nil {
		t.Fatalf("Data.CheckpointResult returned error: %v", err)
	}
	if !resp.Success || !resp.Checkpointed || resp.Detail == nil || *resp.Detail != "ok" {
		t.Fatalf("checkpoint = %#v, want successful checkpoint", resp)
	}
}
