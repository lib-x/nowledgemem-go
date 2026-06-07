package nowledgemem

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// GraphVisService handles graph visualization operations.
type GraphVisService struct {
	client *Client
}

// Explorer returns the interactive graph explorer HTML.
func (s *GraphVisService) Explorer(ctx context.Context) ([]byte, error) {
	return s.client.doBytes(ctx, http.MethodGet, "/graph/vis", nil, nil)
}

// SearchGraph finds relevant content and builds visualization-ready graph data.
func (s *GraphVisService) SearchGraph(ctx context.Context, params *GraphSearchParams) (*GraphData, error) {
	if params == nil {
		return nil, fmt.Errorf("params is required")
	}
	q := url.Values{}
	q.Set("query", params.Query)
	if params.Limit > 0 {
		q.Set("limit", strconv.Itoa(params.Limit))
	}
	if params.Depth > 0 {
		q.Set("depth", strconv.Itoa(params.Depth))
	}
	if len(params.NodeTypes) > 0 {
		for _, t := range params.NodeTypes {
			q.Add("node_types", t)
		}
	}
	if len(params.EdgeTypes) > 0 {
		for _, t := range params.EdgeTypes {
			q.Add("edge_types", t)
		}
	}
	if params.IncludeMetadata != nil {
		q.Set("include_metadata", strconv.FormatBool(*params.IncludeMetadata))
	}
	if params.SpaceID != "" {
		q.Set("space_id", params.SpaceID)
	}
	var resp GraphData
	if err := s.client.doQuery(ctx, "/graph/search", q, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ExploreGraph builds a neighborhood around one or more memory IDs with depth traversal.
func (s *GraphVisService) ExploreGraph(ctx context.Context, params *GraphExploreParams) (*GraphData, error) {
	q := url.Values{}
	q.Set("memory_ids", params.MemoryIDs)
	if params.Depth > 0 {
		q.Set("depth", strconv.Itoa(params.Depth))
	}
	if params.Limit > 0 {
		q.Set("limit", strconv.Itoa(params.Limit))
	}
	if params.SpaceID != "" {
		q.Set("space_id", params.SpaceID)
	}
	var resp GraphData
	if err := s.client.doQuery(ctx, "/graph/explore", q, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// SampleGraph gets a representative sample of graph data for visualization.
func (s *GraphVisService) SampleGraph(ctx context.Context, params *GraphSampleParams) (*GraphData, error) {
	q := url.Values{}
	if params.Limit > 0 {
		q.Set("limit", strconv.Itoa(params.Limit))
	}
	if params.SpaceID != "" {
		q.Set("space_id", params.SpaceID)
	}
	var resp GraphData
	if err := s.client.doQuery(ctx, "/graph/sample", q, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ExpandNode expands neighbors of a specific node.
func (s *GraphVisService) ExpandNode(ctx context.Context, nodeID string, params *GraphExpandParams) (*GraphData, error) {
	q := url.Values{}
	if params != nil {
		if params.Depth > 0 {
			q.Set("depth", strconv.Itoa(params.Depth))
		}
		if params.Limit > 0 {
			q.Set("limit", strconv.Itoa(params.Limit))
		}
		if params.SpaceID != "" {
			q.Set("space_id", params.SpaceID)
		}
	}
	var resp GraphData
	path := fmt.Sprintf("/graph/expand/%s", url.PathEscape(nodeID))
	if err := s.client.doQuery(ctx, path, q, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ShortestPath finds the shortest path between two nodes.
func (s *GraphVisService) ShortestPath(ctx context.Context, sourceID, targetID string) (*GraphData, error) {
	q := url.Values{}
	q.Set("source_id", sourceID)
	q.Set("target_id", targetID)
	var resp GraphData
	if err := s.client.doQuery(ctx, "/graph/shortest-path", q, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetNodeDetails returns detailed information about a specific node.
func (s *GraphVisService) GetNodeDetails(ctx context.Context, nodeID string) (*GraphNode, error) {
	var resp GraphNode
	path := fmt.Sprintf("/graph/node-details/%s", url.PathEscape(nodeID))
	if err := s.client.do(ctx, "GET", path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetCommunityMembers returns members of a specific community.
func (s *GraphVisService) GetCommunityMembers(ctx context.Context, communityID string, limit int) (*GraphData, error) {
	q := url.Values{}
	if limit > 0 {
		q.Set("limit", strconv.Itoa(limit))
	}
	var resp GraphData
	path := fmt.Sprintf("/graph/community-members/%s", url.PathEscape(communityID))
	if err := s.client.doQuery(ctx, path, q, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetOverview returns a high-level graph overview.
func (s *GraphVisService) GetOverview(ctx context.Context) (*GraphOverview, error) {
	var resp GraphOverview
	if err := s.client.do(ctx, "GET", "/graph/overview", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetLivePreviewGraph gets a compact merged graph for one or more seed nodes.
func (s *GraphVisService) GetLivePreviewGraph(ctx context.Context, params *GraphLivePreviewGraphParams) (*GraphData, error) {
	if params == nil || len(params.NodeIDs) == 0 {
		return nil, fmt.Errorf("node_ids is required")
	}
	q := url.Values{}
	q.Set("node_ids", strings.Join(params.NodeIDs, ","))
	if params.LimitPerSeed > 0 {
		q.Set("limit_per_seed", strconv.Itoa(params.LimitPerSeed))
	}
	if params.SpaceID != "" {
		q.Set("space_id", params.SpaceID)
	}
	var resp GraphData
	if err := s.client.doQuery(ctx, "/graph/live-preview", q, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetLivePreview returns a live preview for a node.
func (s *GraphVisService) GetLivePreview(ctx context.Context, nodeID string) (*GraphLivePreview, error) {
	var resp GraphLivePreview
	path := fmt.Sprintf("/graph/live-preview/%s", url.PathEscape(nodeID))
	if err := s.client.do(ctx, "GET", path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// --- Graph Vis types ---

// GraphSearchParams are parameters for graph search.
type GraphSearchParams struct {
	Query           string   `json:"query"`
	Limit           int      `json:"limit,omitempty"`
	Depth           int      `json:"depth,omitempty"`
	NodeTypes       []string `json:"node_types,omitempty"`
	EdgeTypes       []string `json:"edge_types,omitempty"`
	IncludeMetadata *bool    `json:"include_metadata,omitempty"`
	SpaceID         string   `json:"space_id,omitempty"`
}

// GraphExploreParams are parameters for graph explore.
type GraphExploreParams struct {
	MemoryIDs string `json:"memory_ids"`
	Depth     int    `json:"depth,omitempty"`
	Limit     int    `json:"limit,omitempty"`
	SpaceID   string `json:"space_id,omitempty"`
}

// GraphSampleParams are parameters for graph sample.
type GraphSampleParams struct {
	Limit   int    `json:"limit,omitempty"`
	SpaceID string `json:"space_id,omitempty"`
}

// GraphExpandParams are parameters for graph expand.
type GraphExpandParams struct {
	Depth   int    `json:"depth,omitempty"`
	Limit   int    `json:"limit,omitempty"`
	SpaceID string `json:"space_id,omitempty"`
}

// GraphLivePreviewGraphParams are parameters for GET /graph/live-preview.
type GraphLivePreviewGraphParams struct {
	NodeIDs      []string `json:"node_ids"`
	LimitPerSeed int      `json:"limit_per_seed,omitempty"`
	SpaceID      string   `json:"space_id,omitempty"`
}

// GraphData is the visualization-ready graph response.
type GraphData struct {
	Nodes               []GraphNode      `json:"nodes"`
	Edges               []GraphEdge      `json:"edges"`
	Communities         []map[string]any `json:"communities,omitempty"`
	CommunityHulls      []map[string]any `json:"community_hulls,omitempty"`
	VisualizationConfig map[string]any   `json:"visualization_config,omitempty"`
	Metadata            map[string]any   `json:"metadata,omitempty"`
}

// GraphNode represents a node in the graph visualization.
type GraphNode struct {
	ID            string         `json:"id"`
	Label         string         `json:"label"`
	NodeType      string         `json:"node_type"`
	NodeSubtype   string         `json:"node_subtype,omitempty"`
	Size          float64        `json:"size,omitempty"`
	Color         string         `json:"color,omitempty"`
	Community     string         `json:"community,omitempty"`
	Importance    float64        `json:"importance,omitempty"`
	HopCount      int            `json:"hop_count,omitempty"`
	PagerankScore float64        `json:"pagerank_score,omitempty"`
	ThreadID      string         `json:"thread_id,omitempty"`
	Metadata      map[string]any `json:"metadata,omitempty"`
}

// GraphEdge represents an edge in the graph visualization.
type GraphEdge struct {
	ID             string         `json:"id"`
	Source         string         `json:"source"`
	Target         string         `json:"target"`
	EdgeType       string         `json:"edge_type"`
	Weight         float64        `json:"weight,omitempty"`
	Label          string         `json:"label,omitempty"`
	RelevanceScore float64        `json:"relevance_score,omitempty"`
	Metadata       map[string]any `json:"metadata,omitempty"`
}

// GraphOverview is the response for GET /graph/overview.
type GraphOverview struct {
	NodeCount      int     `json:"node_count"`
	EdgeCount      int     `json:"edge_count"`
	CommunityCount int     `json:"community_count"`
	MemoryCount    int     `json:"memory_count"`
	EntityCount    int     `json:"entity_count"`
	ThreadCount    int     `json:"thread_count"`
	AvgDegree      float64 `json:"avg_degree"`
}

// GraphLivePreview is a live preview for a node.
type GraphLivePreview struct {
	Node  GraphNode   `json:"node"`
	Edges []GraphEdge `json:"edges"`
}
