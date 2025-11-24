package documents

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	xaiv1 "github.com/ZaguanLabs/xai-sdk-go/proto/gen/go/xai/api/v1"
	"github.com/ZaguanLabs/xai-sdk-go/xai/internal/rest"
)

func TestNewClient(t *testing.T) {
	restClient := &rest.Client{}
	client := NewClient(restClient)

	if client == nil {
		t.Fatal("NewClient returned nil")
	}
	if client.restClient != restClient {
		t.Error("restClient not set correctly")
	}
}

func TestNewSearchRequest(t *testing.T) {
	req := NewSearchRequest("test query", "col1", "col2")

	if req == nil {
		t.Fatal("NewSearchRequest returned nil")
	}
	if req.Query != "test query" {
		t.Errorf("Query = %v, want test query", req.Query)
	}
	if len(req.CollectionIDs) != 2 {
		t.Errorf("len(CollectionIDs) = %v, want 2", len(req.CollectionIDs))
	}
	if req.Limit != 10 {
		t.Errorf("Limit = %v, want 10", req.Limit)
	}
}

func TestSearchRequestWithLimit(t *testing.T) {
	req := NewSearchRequest("test").WithLimit(20)

	if req.Limit != 20 {
		t.Errorf("Limit = %v, want 20", req.Limit)
	}
}

func TestSearch(t *testing.T) {
	tests := []struct {
		name       string
		req        *SearchRequest
		statusCode int
		response   *xaiv1.SearchResponse
		wantErr    bool
	}{
		{
			name:       "success",
			req:        NewSearchRequest("test query", "col1"),
			statusCode: http.StatusOK,
			response: &xaiv1.SearchResponse{
				Matches: []*xaiv1.SearchMatch{
					{
						FileId:        "file-123",
						ChunkId:       "chunk-456",
						ChunkContent:  "This is a test document",
						Score:         0.95,
						CollectionIds: []string{"col1"},
					},
					{
						FileId:       "file-789",
						ChunkId:      "chunk-012",
						ChunkContent: "Another document",
						Score:        0.85,
					},
				},
			},
			wantErr: false,
		},
		{
			name:       "empty results",
			req:        NewSearchRequest("no match"),
			statusCode: http.StatusOK,
			response: &xaiv1.SearchResponse{
				Matches: []*xaiv1.SearchMatch{},
			},
			wantErr: false,
		},
		{
			name:       "server error",
			req:        NewSearchRequest("test"),
			statusCode: http.StatusInternalServerError,
			response:   nil,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Errorf("Expected POST, got %s", r.Method)
				}
				if r.URL.Path != "/documents/search" {
					t.Errorf("Expected /documents/search, got %s", r.URL.Path)
				}

				w.WriteHeader(tt.statusCode)
				if tt.response != nil {
					json.NewEncoder(w).Encode(tt.response)
				}
			}))
			defer server.Close()

			restClient := rest.NewClient(rest.Config{
				BaseURL: server.URL,
				APIKey:  "test",
			})
			client := NewClient(restClient)

			resp, err := client.Search(context.Background(), tt.req)

			if (err != nil) != tt.wantErr {
				t.Errorf("Search() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if resp == nil {
					t.Fatal("Search() returned nil response")
				}
				if len(resp.Matches) != len(tt.response.Matches) {
					t.Errorf("Search() returned %d matches, want %d", len(resp.Matches), len(tt.response.Matches))
				}
				if len(resp.Matches) > 0 {
					if resp.Matches[0].FileID != tt.response.Matches[0].FileId {
						t.Errorf("FileID = %v, want %v", resp.Matches[0].FileID, tt.response.Matches[0].FileId)
					}
					if resp.Matches[0].Score != tt.response.Matches[0].Score {
						t.Errorf("Score = %v, want %v", resp.Matches[0].Score, tt.response.Matches[0].Score)
					}
				}
			}
		})
	}
}

func TestSearchClientNotInitialized(t *testing.T) {
	client := &Client{restClient: nil}
	req := NewSearchRequest("test")

	_, err := client.Search(context.Background(), req)
	if err != ErrClientNotInitialized {
		t.Errorf("Search() error = %v, want %v", err, ErrClientNotInitialized)
	}
}

func TestSearchMatch(t *testing.T) {
	match := &SearchMatch{
		FileID:        "file-123",
		ChunkID:       "chunk-456",
		ChunkContent:  "test content",
		Score:         0.95,
		CollectionIDs: []string{"col1", "col2"},
	}

	if match.FileID != "file-123" {
		t.Errorf("FileID = %v, want file-123", match.FileID)
	}
	if match.Score != 0.95 {
		t.Errorf("Score = %v, want 0.95", match.Score)
	}
	if len(match.CollectionIDs) != 2 {
		t.Errorf("len(CollectionIDs) = %v, want 2", len(match.CollectionIDs))
	}
}
