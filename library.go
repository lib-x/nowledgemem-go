package nowledgemem

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

// LibraryService handles library/wiki operations.
type LibraryService struct {
	client *Client
}

// GetWikiIndex returns the wiki index.
func (s *LibraryService) GetWikiIndex(ctx context.Context) (*WikiIndex, error) {
	var resp WikiIndex
	if err := s.client.do(ctx, "GET", "/library/wiki-index", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ExportWiki exports wiki pages in the specified format.
func (s *LibraryService) ExportWiki(ctx context.Context, format string) ([]byte, error) {
	q := url.Values{}
	if format != "" {
		q.Set("format", format)
	}
	u := s.client.baseURL.ResolveReference(&url.URL{Path: "/library/wiki-export"})
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	resp, err := s.client.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		var apiErr APIError
		if err := json.NewDecoder(resp.Body).Decode(&apiErr); err != nil {
			return nil, fmt.Errorf("API error (status %d)", resp.StatusCode)
		}
		return nil, &apiErr
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}
	return data, nil
}

// GetWikiPageByEntity returns a wiki page for an entity.
func (s *LibraryService) GetWikiPageByEntity(ctx context.Context, idOrName string) (*WikiPage, error) {
	var resp WikiPage
	path := fmt.Sprintf("/library/wiki-page/entity/%s", url.PathEscape(idOrName))
	if err := s.client.do(ctx, "GET", path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetWikiPageByTopic returns a wiki page for a topic/community.
func (s *LibraryService) GetWikiPageByTopic(ctx context.Context, communityID string) (*WikiPage, error) {
	var resp WikiPage
	path := fmt.Sprintf("/library/wiki-page/topic/%s", url.PathEscape(communityID))
	if err := s.client.do(ctx, "GET", path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetWikiPageByCrystal returns a wiki page for a crystal.
func (s *LibraryService) GetWikiPageByCrystal(ctx context.Context, crystalID string) (*WikiPage, error) {
	var resp WikiPage
	path := fmt.Sprintf("/library/wiki-page/crystal/%s", url.PathEscape(crystalID))
	if err := s.client.do(ctx, "GET", path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetCrystalSourceMemories returns source memories for a crystal.
func (s *LibraryService) GetCrystalSourceMemories(ctx context.Context, crystalID string, limit int) ([]MemoryListItem, error) {
	q := url.Values{}
	if limit > 0 {
		q.Set("limit", strconv.Itoa(limit))
	}
	var resp []MemoryListItem
	path := fmt.Sprintf("/library/crystal/%s/source-memories", url.PathEscape(crystalID))
	if err := s.client.doQuery(ctx, path, q, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// --- Library types ---

// WikiIndex is the response for GET /library/wiki-index.
type WikiIndex struct {
	Pages []WikiIndexEntry `json:"pages"`
}

// WikiIndexEntry is a single entry in the wiki index.
type WikiIndexEntry struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Type  string `json:"type"`
}

// WikiPage represents a wiki page.
type WikiPage struct {
	ID       string         `json:"id"`
	Title    string         `json:"title"`
	Content  string         `json:"content"`
	Type     string         `json:"type"`
	Metadata map[string]any `json:"metadata,omitempty"`
}

// ExportWikiSummary exports wiki summary.
func (s *LibraryService) ExportWikiSummary(ctx context.Context) ([]byte, error) {
	u := s.client.baseURL.ResolveReference(&url.URL{Path: "/library/wiki-export-summary"})
	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	resp, err := s.client.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		var apiErr APIError
		if err := json.NewDecoder(resp.Body).Decode(&apiErr); err != nil {
			return nil, fmt.Errorf("API error (status %d)", resp.StatusCode)
		}
		return nil, &apiErr
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}
	return data, nil
}
