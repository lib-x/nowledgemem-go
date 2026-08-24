package openapi

import (
	"encoding/json"
	"os"
	"testing"
)

func TestSnapshotContainsDataTransferContract(t *testing.T) {
	raw, err := os.ReadFile("openapi.json")
	if err != nil {
		t.Fatalf("read OpenAPI snapshot: %v", err)
	}
	var document struct {
		OpenAPI string `json:"openapi"`
		Info    struct {
			Title   string `json:"title"`
			Version string `json:"version"`
		} `json:"info"`
		Paths map[string]json.RawMessage `json:"paths"`
	}
	if err := json.Unmarshal(raw, &document); err != nil {
		t.Fatalf("decode OpenAPI snapshot: %v", err)
	}
	if document.OpenAPI != "3.1.0" || document.Info.Title != "Nowledge Mem API" || document.Info.Version == "" {
		t.Fatalf("unexpected OpenAPI metadata: %#v", document)
	}
	for _, path := range []string{"/data/checkpoint", "/data/export/download", "/data/import/upload"} {
		if _, ok := document.Paths[path]; !ok {
			t.Errorf("OpenAPI snapshot is missing %s", path)
		}
	}
}

func TestGeneratedDataExportIncludesFeedEvents(t *testing.T) {
	value := true
	body, err := json.Marshal(DataExportDownloadRequest{IncludeFeedEvents: &value})
	if err != nil {
		t.Fatalf("marshal request: %v", err)
	}
	if string(body) != `{"include_feed_events":true}` {
		t.Fatalf("body = %s, want include_feed_events", body)
	}
}
