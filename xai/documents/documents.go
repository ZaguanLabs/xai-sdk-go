// Package documents provides a client for the xAI Document Search API.
package documents

import (
	"context"

	xaiv1 "github.com/ZaguanLabs/xai-sdk-go/proto/gen/go/xai/v1"
	"github.com/ZaguanLabs/xai-sdk-go/xai/internal/rest"
	"google.golang.org/protobuf/encoding/protojson"
)

// Client provides access to the xAI Document Search API.
type Client struct {
	restClient *rest.Client
}

// NewClient creates a new Document Search API client.
func NewClient(restClient *rest.Client) *Client {
	return &Client{
		restClient: restClient,
	}
}

// SearchRequest represents a document search request.
type SearchRequest struct {
	Query         string
	CollectionIDs []string
	Limit         int32
}

// SearchMatch represents a single search result.
type SearchMatch struct {
	FileID        string
	ChunkID       string
	ChunkContent  string
	Score         float32
	CollectionIDs []string
}

// SearchResponse represents the search results.
type SearchResponse struct {
	Matches []*SearchMatch
}

// NewSearchRequest creates a new document search request.
func NewSearchRequest(query string, collectionIDs ...string) *SearchRequest {
	return &SearchRequest{
		Query:         query,
		CollectionIDs: collectionIDs,
		Limit:         10, // default limit
	}
}

// WithLimit sets the maximum number of results to return.
func (r *SearchRequest) WithLimit(limit int32) *SearchRequest {
	r.Limit = limit
	return r
}

// Search searches across document collections.
func (c *Client) Search(ctx context.Context, req *SearchRequest) (*SearchResponse, error) {
	if c.restClient == nil {
		return nil, ErrClientNotInitialized
	}

	protoReq := &xaiv1.SearchRequest{
		Query: req.Query,
		Source: &xaiv1.DocumentsSource{
			CollectionIds: req.CollectionIDs,
		},
		Limit: req.Limit,
	}

	jsonData, err := protojson.Marshal(protoReq)
	if err != nil {
		return nil, err
	}

	resp, err := c.restClient.Post(ctx, "/documents/search", jsonData)
	if err != nil {
		return nil, err
	}

	var searchResp xaiv1.SearchResponse
	if err := protojson.Unmarshal(resp.Body, &searchResp); err != nil {
		return nil, err
	}

	matches := make([]*SearchMatch, len(searchResp.Matches))
	for i, m := range searchResp.Matches {
		matches[i] = &SearchMatch{
			FileID:        m.FileId,
			ChunkID:       m.ChunkId,
			ChunkContent:  m.ChunkContent,
			Score:         m.Score,
			CollectionIDs: m.CollectionIds,
		}
	}

	return &SearchResponse{
		Matches: matches,
	}, nil
}
