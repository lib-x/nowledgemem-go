package onledgemem

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

// FSService handles Nowledge FS operations (path-based tree browsing).
type FSService struct {
	client *Client
}

// List lists a directory in the FS tree.
func (s *FSService) List(ctx context.Context, path string, limit int, cursor string) (*FSListResponse, error) {
	q := url.Values{}
	q.Set("path", path)
	if limit > 0 {
		q.Set("limit", strconv.Itoa(limit))
	}
	if cursor != "" {
		q.Set("cursor", cursor)
	}
	var resp FSListResponse
	if err := s.client.doQuery(ctx, "/fs/ls", q, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Cat reads a rendered file body and frontmatter.
func (s *FSService) Cat(ctx context.Context, path string, line, lines int) (*FSCatResponse, error) {
	q := url.Values{}
	q.Set("path", path)
	if line > 0 {
		q.Set("line", strconv.Itoa(line))
	}
	if lines > 0 {
		q.Set("lines", strconv.Itoa(lines))
	}
	var resp FSCatResponse
	if err := s.client.doQuery(ctx, "/fs/cat", q, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Stat reads metadata without loading the body.
func (s *FSService) Stat(ctx context.Context, path string) (*FSStatResponse, error) {
	q := url.Values{}
	q.Set("path", path)
	var resp FSStatResponse
	if err := s.client.doQuery(ctx, "/fs/stat", q, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Find performs structural search (type, label, date, mention constraints).
func (s *FSService) Find(ctx context.Context, path, fileType, label string) (*FSSearchResponse, error) {
	q := url.Values{}
	q.Set("path", path)
	if fileType != "" {
		q.Set("type", fileType)
	}
	if label != "" {
		q.Set("label", label)
	}
	var resp FSSearchResponse
	if err := s.client.doQuery(ctx, "/fs/find", q, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Grep performs literal exact-string search.
func (s *FSService) Grep(ctx context.Context, path, query string) (*FSGrepResponse, error) {
	q := url.Values{}
	q.Set("path", path)
	q.Set("q", query)
	var resp FSGrepResponse
	if err := s.client.doQuery(ctx, "/fs/grep", q, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Recall performs semantic search that returns paths.
func (s *FSService) Recall(ctx context.Context, path, query string, k int) (*FSSearchResponse, error) {
	q := url.Values{}
	q.Set("path", path)
	q.Set("query", query)
	if k > 0 {
		q.Set("k", strconv.Itoa(k))
	}
	var resp FSSearchResponse
	if err := s.client.doQuery(ctx, "/fs/recall", q, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Write updates a file at the given path.
func (s *FSService) Write(ctx context.Context, path, body string) error {
	req := FSWriteRequest{Path: path, Body: body}
	return s.client.do(ctx, "POST", "/fs/write", req, nil)
}

// Delete deletes a file at the given path.
func (s *FSService) Delete(ctx context.Context, path string) error {
	req := FSDeleteRequest{Path: path}
	return s.client.do(ctx, "POST", "/fs/delete", req, nil)
}

// CatString is a convenience method that returns just the body text.
func (s *FSService) CatString(ctx context.Context, path string) (string, error) {
	resp, err := s.Cat(ctx, path, 0, 0)
	if err != nil {
		return "", err
	}
	return resp.Body, nil
}

// LsPaths is a convenience method that returns just the entry paths.
func (s *FSService) LsPaths(ctx context.Context, path string) ([]string, error) {
	resp, err := s.List(ctx, path, 0, "")
	if err != nil {
		return nil, err
	}
	paths := make([]string, len(resp.Entries))
	for i, e := range resp.Entries {
		paths[i] = fmt.Sprintf("%s/%s", path, e.Name)
	}
	return paths, nil
}
