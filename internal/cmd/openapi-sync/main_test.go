package main

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestPrepareAppliesReviewedPathParameterCorrections(t *testing.T) {
	raw := []byte(`{
  "openapi": "3.1.0",
  "info": {"title": "Nowledge Mem API", "version": "test"},
  "paths": {
    "/agent/ai-now/schedules": {
      "get": {
        "parameters": [{"name": "include_deleted", "in": "path", "required": true, "schema": {"type": "boolean"}}],
        "responses": {"200": {"description": "ok"}}
      }
    }
  }
}`)

	prepared, version, corrections, err := prepare(raw)
	if err != nil {
		t.Fatalf("prepare: %v", err)
	}
	if version != "test" {
		t.Fatalf("version = %q, want test", version)
	}
	if corrections != 1 {
		t.Fatalf("corrections = %d, want 1", corrections)
	}

	var document map[string]any
	if err := json.Unmarshal(prepared, &document); err != nil {
		t.Fatalf("decode prepared document: %v", err)
	}
	if document["openapi"] != "3.1.0" {
		t.Fatalf("openapi = %v, want 3.1.0", document["openapi"])
	}
	parameter := document["paths"].(map[string]any)["/agent/ai-now/schedules"].(map[string]any)["get"].(map[string]any)["parameters"].([]any)[0].(map[string]any)
	if parameter["in"] != "query" || parameter["required"] != false {
		t.Fatalf("parameter = %#v, want optional query parameter", parameter)
	}
}

func TestPrepareRejectsExternalReference(t *testing.T) {
	raw := []byte(`{
  "openapi": "3.1.0",
  "info": {"title": "Nowledge Mem API", "version": "test"},
  "paths": {},
  "components": {"schemas": {"Unsafe": {"$ref": "http://127.0.0.1/private"}}}
}`)
	_, _, _, err := prepare(raw)
	if err == nil || !strings.Contains(err.Error(), "external OpenAPI reference") {
		t.Fatalf("error = %v, want external reference rejection", err)
	}
}

func TestPrepareRejectsUnknownPathParameterMismatch(t *testing.T) {
	raw := []byte(`{
  "openapi": "3.1.0",
  "info": {"title": "Nowledge Mem API", "version": "test"},
  "paths": {"/items": {"get": {"parameters": [{"name": "limit", "in": "path", "required": true}], "responses": {"200": {"description": "ok"}}}}}
}`)
	_, _, _, err := prepare(raw)
	if err == nil || !strings.Contains(err.Error(), "unreviewed path parameter mismatch") {
		t.Fatalf("error = %v, want mismatch rejection", err)
	}
}

func TestPrepareRejectsUnexpectedDocument(t *testing.T) {
	_, _, _, err := prepare([]byte(`{"openapi":"3.1.0","info":{"title":"Other API","version":"1"},"paths":{}}`))
	if err == nil {
		t.Fatal("expected title validation error")
	}
}
