package nowledgemem

import (
	"context"
	"fmt"
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
	return s.client.doBytes(ctx, http.MethodGet, "/library/wiki-export", q, nil)
}

// ExportOKF exports the library as Open Knowledge Format JSON.
func (s *LibraryService) ExportOKF(ctx context.Context, params *OKFExportParams) (map[string]any, error) {
	q := url.Values{}
	if params != nil {
		if params.EntityLimit > 0 {
			q.Set("entity_limit", strconv.Itoa(params.EntityLimit))
		}
		if params.TopPerCommunity > 0 {
			q.Set("top_per_community", strconv.Itoa(params.TopPerCommunity))
		}
		if params.MaxMentionsPerEntity > 0 {
			q.Set("max_mentions_per_entity", strconv.Itoa(params.MaxMentionsPerEntity))
		}
	}
	var resp map[string]any
	if err := s.client.doQuery(ctx, "/library/okf-export", q, &resp); err != nil {
		return nil, err
	}
	return resp, nil
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

// OKFExportParams are query parameters for ExportOKF.
type OKFExportParams struct {
	EntityLimit          int `json:"entity_limit,omitempty"`
	TopPerCommunity      int `json:"top_per_community,omitempty"`
	MaxMentionsPerEntity int `json:"max_mentions_per_entity,omitempty"`
}

// ExportWikiSummary exports wiki summary.
func (s *LibraryService) ExportWikiSummary(ctx context.Context) ([]byte, error) {
	return s.client.doBytes(ctx, http.MethodGet, "/library/wiki-export-summary", nil, nil)
}
