package nowledgemem

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func TestUpstreamFieldsRoundTrip(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		body  string
		value any
	}{
		{"memory provenance", `{"write_intent":{"version":1,"kind":"authored"},"claim_status":"observed","source_grounding":"supported"}`, &CreateMemoryRequest{}},
		{"source-derived memory", `{"write_intent":{"version":1,"kind":"source_derived","assessment":"human_ruling","source":{"kind":"thread","thread_id":"thread-1","message_ids":["msg-1"]}}}`, &CreateMemoryRequest{}},
		{"scoped FS delete", `{"space_id":"engineering"}`, &FSDeleteRequest{}},
		{"thread project", `{"project":"sdk","workspace":"/work/sdk"}`, &ThreadListItem{}},
		{"thread deletion", `{"index_repair_pending":true,"storage_cleanup_pending":true}`, &DeleteThreadResponse{}},
		{"explicit false deletion status", `{"index_repair_pending":false,"storage_cleanup_pending":false}`, &DeleteThreadResponse{}},
		{"bulk deletion", `{"partial_success":true,"index_repair_pending":true,"deletion_outcome_uncertain":true}`, &ThreadBulkDeleteResponse{}},
		{"agent identity", `{"identityReference":"opaque-reference"}`, &AgentProfileResponse{}},
		{"KG request", `{"entity_resolution":"exact"}`, &KGApplyRequest{}},
		{"KG response", `{"entity_resolution":"exact","index_repair_needed":true,"warnings":["repair pending"],"resolved_entities":[{"id":"entity-1","name":"SDK","entity_type":"Project"}]}`, &KGApplyResponse{}},
		{"distillation", `{"native_source_operation":{"version":1,"kind":"source_derived","assessment":"source_support","source":{"kind":"thread","thread_id":"thread-1","message_ids":["msg-1"]}}}`, &DistillRequest{}},
		{"distillation schedule", `{"distill_again_of_task_id":"task-1","memory_policy_override":{"priorities":[],"exclusions":[]}}`, &DistillScheduleRequest{}},
		{"distillation response", `{"success":true,"task_id":"task-1","distill_again":{"run_id":"run-1","status":"queued"}}`, &DistillScheduleResponse{}},
		{"export timeline", `{"include_feed_events":false}`, &DataExportRequest{}},
		{"download timeline", `{"include_feed_events":false}`, &DataExportDownloadRequest{}},
		{"import timeline", `{"include_feed_events":false}`, &DataImportRequest{}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			if err := json.Unmarshal([]byte(test.body), test.value); err != nil {
				t.Fatal(err)
			}
			encoded, err := json.Marshal(test.value)
			if err != nil {
				t.Fatal(err)
			}
			var want, got map[string]any
			if err := json.Unmarshal([]byte(test.body), &want); err != nil {
				t.Fatal(err)
			}
			if err := json.Unmarshal(encoded, &got); err != nil {
				t.Fatal(err)
			}
			for key, value := range want {
				if !reflect.DeepEqual(got[key], value) {
					t.Errorf("%s = %#v, want %#v", key, got[key], value)
				}
			}
		})
	}
}

func TestUpdatedServiceContracts(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		method   string
		path     string
		query    map[string]string
		body     map[string]any
		form     map[string]string
		response string
		call     func(*Client) error
	}{
		{
			name: "thread project filter", method: "GET", path: "/threads",
			query:    map[string]string{"project_key": "opaque+/project=", "space_id": "work"},
			response: `{"threads":[],"pagination":{}}`,
			call: func(client *Client) error {
				_, err := client.Threads.List(t.Context(), &ListThreadsParams{ProjectKey: "opaque+/project=", SpaceID: "work"})
				return err
			},
		},
		{
			name: "thread projects", method: "GET", path: "/threads/projects",
			query:    map[string]string{"limit": "5", "offset": "2", "source": "codex", "space_id": "work"},
			response: `{"projects":[{"key":"opaque","kind":"project","value":"sdk","count":3}],"pagination":{"total":1}}`,
			call: func(client *Client) error {
				response, err := client.Threads.Projects(t.Context(), &ListThreadProjectsParams{Limit: 5, Offset: 2, Source: "codex", SpaceID: "work"})
				if err == nil && (len(response.Projects) != 1 || response.Projects[0].Key != "opaque" || response.Pagination.Total != 1) {
					t.Errorf("unexpected projects: %#v", response)
				}
				return err
			},
		},
		{
			name: "scoped FS delete", method: "POST", path: "/fs/delete",
			body: map[string]any{"path": "/memories/by-id/memory-1.memory.md", "space_id": "work"}, response: `{}`,
			call: func(client *Client) error {
				return client.FS.DeleteWithOptions(t.Context(), &FSDeleteRequest{Path: "/memories/by-id/memory-1.memory.md", SpaceID: new("work")})
			},
		},
		{
			name: "import preview", method: "POST", path: "/data/import/preview",
			body:     map[string]any{"import_path": "/tmp/export.zip", "include_feed_events": false},
			response: `{"success":true,"preview":{"inspection_only":true}}`,
			call: func(client *Client) error {
				response, err := client.Data.PreviewImport(t.Context(), &DataImportRequest{ImportPath: "/tmp/export.zip", IncludeFeedEvents: new(false)})
				if err == nil && (!response.Success || !response.Preview.InspectionOnly) {
					t.Errorf("unexpected preview: %#v", response)
				}
				return err
			},
		},
		{
			name: "upload preview", method: "POST", path: "/data/import/upload/preview",
			form:     map[string]string{"include_feed_events": "false", "mode": "merge"},
			response: `{"success":true,"preview":{"inspection_only":true}}`,
			call: func(client *Client) error {
				_, err := client.Data.PreviewUploadImport(t.Context(), &UploadImportRequest{File: strings.NewReader("zip"), Mode: "merge", IncludeFeedEvents: new(false)})
				return err
			},
		},
		{
			name: "upload import", method: "POST", path: "/data/import/upload",
			form: map[string]string{"include_feed_events": "false"}, response: `{"success":true}`,
			call: func(client *Client) error {
				_, err := client.Data.UploadImport(t.Context(), &UploadImportRequest{File: strings.NewReader("zip"), IncludeFeedEvents: new(false)})
				return err
			},
		},
		{
			name: "attachment provenance", method: "POST", path: "/sources/ingest/file",
			form: map[string]string{"attachment_reference": `{"thread_id":"thread-1"}`}, response: `{}`,
			call: func(client *Client) error {
				_, err := client.Sources.IngestFile(t.Context(), &IngestFileRequest{File: strings.NewReader("file"), AttachmentReference: `{"thread_id":"thread-1"}`})
				return err
			},
		},
		{
			name: "optimize outcome", method: "POST", path: "/storage/optimize",
			response: `{"success":false,"outcome":"partial","checkpoint_failed":true,"reclaim_pending":true}`,
			call: func(client *Client) error {
				response, err := client.Storage.OptimizeResult(t.Context())
				if err == nil && (response.Outcome != "partial" || !response.CheckpointFailed || !response.ReclaimPending) {
					t.Errorf("unexpected optimize outcome: %#v", response)
				}
				return err
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				if request.Method != test.method || request.URL.Path != test.path {
					t.Errorf("request = %s %s, want %s %s", request.Method, request.URL.Path, test.method, test.path)
				}
				for key, want := range test.query {
					if got := request.URL.Query().Get(key); got != want {
						t.Errorf("query %s = %q, want %q", key, got, want)
					}
				}
				if test.body != nil {
					var body map[string]any
					content, err := io.ReadAll(request.Body)
					if err != nil {
						t.Error(err)
						return
					}
					if err := json.Unmarshal(content, &body); err != nil {
						t.Error(err)
					}
					if !reflect.DeepEqual(body, test.body) {
						t.Errorf("body = %#v, want %#v", body, test.body)
					}
				}
				if test.form != nil {
					if err := request.ParseMultipartForm(1 << 20); err != nil {
						t.Error(err)
						return
					}
					defer request.MultipartForm.RemoveAll()
					for key, want := range test.form {
						if got := request.FormValue(key); got != want {
							t.Errorf("form %s = %q, want %q", key, got, want)
						}
					}
					file, _, err := request.FormFile("file")
					if err != nil {
						t.Error(err)
						return
					}
					defer file.Close()
					content, err := io.ReadAll(file)
					if err != nil || len(content) == 0 {
						t.Errorf("uploaded file = %q, error = %v", content, err)
					}
				}
				writer.Header().Set("Content-Type", "application/json")
				_, _ = io.WriteString(writer, test.response)
			}))
			defer server.Close()
			if err := test.call(NewClient(WithBaseURL(server.URL))); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestPreviewUploadImportRejectsMissingFile(t *testing.T) {
	t.Parallel()
	for _, request := range []*UploadImportRequest{nil, {}} {
		if _, err := NewClient().Data.PreviewUploadImport(t.Context(), request); err == nil {
			t.Error("expected missing request or file error")
		}
	}
}
