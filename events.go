package nowledgemem

import (
	"context"
	"net/http"
)

// EventsService handles server-sent event streams.
type EventsService struct {
	client *Client
}

// Stream opens the real-time server-sent events stream.
//
// The caller owns the returned response body and must close it.
func (s *EventsService) Stream(ctx context.Context) (*http.Response, error) {
	return s.client.doStream(ctx, http.MethodGet, "/events/stream", nil, nil)
}
