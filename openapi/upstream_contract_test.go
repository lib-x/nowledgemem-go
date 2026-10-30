package openapi

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestNewOperationRequests(t *testing.T) {
	t.Parallel()
	const server = "http://localhost:14242"
	tests := []struct {
		method string
		path   string
		build  func() (*http.Request, error)
	}{
		{"POST", "/agent/feed/events/event-1/review", func() (*http.Request, error) {
			return NewOpenapiGuardedTimelineReviewRequest(server, "event-1", OpenapiGuardedTimelineReviewJSONRequestBody{})
		}},
		{"POST", "/data/import/preview", func() (*http.Request, error) {
			return NewOpenapiImportPreviewRequest(server, OpenapiImportPreviewJSONRequestBody{})
		}},
		{"POST", "/data/import/upload/preview", func() (*http.Request, error) {
			return NewOpenapiImportUploadPreviewRequestWithBody(server, "multipart/form-data; boundary=test", strings.NewReader("--test--\r\n"))
		}},
		{"GET", "/entities/entity-1/aliases", func() (*http.Request, error) {
			return NewReadEntityAliasesHandlerRequest(server, "entity-1")
		}},
		{"POST", "/entities/entity-1/aliases", func() (*http.Request, error) {
			return NewChangeEntityAliasesHandlerRequest(server, "entity-1", ChangeEntityAliasesHandlerJSONRequestBody{})
		}},
		{"GET", "/identity/profiles?space_id=engineering", func() (*http.Request, error) {
			return NewGetIdentityProfilesRequest(server, &GetIdentityProfilesParams{SpaceId: "engineering"})
		}},
		{"POST", "/identity/profiles", func() (*http.Request, error) {
			return NewSaveIdentityProfileRequest(server, SaveIdentityProfileJSONRequestBody{})
		}},
		{"GET", "/memories/memory-1/identity?space_id=engineering", func() (*http.Request, error) {
			return NewReadMemoryIdentityRequest(server, "memory-1", &ReadMemoryIdentityParams{SpaceId: "engineering"})
		}},
		{"POST", "/memories/memory-1/identity/correction", func() (*http.Request, error) {
			return NewCorrectMemoryIdentityRequest(server, "memory-1", CorrectMemoryIdentityJSONRequestBody{})
		}},
		{"POST", "/memories/memory-1/identity/extract", func() (*http.Request, error) {
			return NewExtractMemoryIdentityRequest(server, "memory-1", ExtractMemoryIdentityJSONRequestBody{})
		}},
		{"GET", "/threads/projects?space_id=engineering", func() (*http.Request, error) {
			return NewListThreadProjectsThreadsProjectsGetRequest(server, &ListThreadProjectsThreadsProjectsGetParams{SpaceId: new("engineering")})
		}},
	}
	for _, test := range tests {
		t.Run(test.method+" "+test.path, func(t *testing.T) {
			t.Parallel()
			request, err := test.build()
			if err != nil {
				t.Fatal(err)
			}
			if request.Body != nil {
				defer request.Body.Close()
			}
			if request.Method != test.method || request.URL.RequestURI() != test.path {
				t.Errorf("request = %s %s, want %s %s", request.Method, request.URL.RequestURI(), test.method, test.path)
			}
		})
	}
}

func TestThreadDeletionErrorRetainsPartialOutcome(t *testing.T) {
	t.Parallel()
	for _, status := range []int{http.StatusForbidden, http.StatusConflict, http.StatusInternalServerError} {
		response, err := ParseDeleteThreadThreadsThreadIdDeleteResponse(&http.Response{
			StatusCode: status,
			Header:     http.Header{"Content-Type": {"application/json"}},
			Body:       io.NopCloser(strings.NewReader(`{"detail":"refresh before retrying","deleted_messages":3,"partial_success":true,"deletion_outcome_uncertain":true,"index_repair_pending":true}`)),
		})
		if err != nil {
			t.Fatal(err)
		}
		outcome := map[int]*ThreadDeletionErrorResponse{
			http.StatusForbidden:           response.JSON403,
			http.StatusConflict:            response.JSON409,
			http.StatusInternalServerError: response.JSON500,
		}[status]
		if outcome == nil || outcome.PartialSuccess == nil || !*outcome.PartialSuccess || outcome.DeletionOutcomeUncertain == nil || !*outcome.DeletionOutcomeUncertain || outcome.DeletedMessages == nil || *outcome.DeletedMessages != 3 {
			t.Errorf("status %d: partial outcome lost: %#v", status, outcome)
		}
	}
}
