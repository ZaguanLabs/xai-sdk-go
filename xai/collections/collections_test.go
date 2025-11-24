package collections

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	xaiv1 "github.com/ZaguanLabs/xai-sdk-go/proto/gen/go/xai/api/v1"
	"github.com/ZaguanLabs/xai-sdk-go/xai/internal/rest"
	"google.golang.org/protobuf/types/known/timestamppb"
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

func TestCreateCollection(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/collections" {
			t.Errorf("Expected POST /collections, got %s %s", r.Method, r.URL.Path)
		}

		resp := &xaiv1.CollectionMetadata{
			CollectionId:   "col-123",
			CollectionName: "Test Collection",
			DocumentsCount: 0,
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	restClient := rest.NewClient(rest.Config{BaseURL: server.URL, APIKey: "test"})
	client := NewClient(restClient)

	col, err := client.CreateCollection(context.Background(), CreateCollectionOptions{
		Name:   "Test Collection",
		TeamID: "team-1",
	})

	if err != nil {
		t.Fatalf("CreateCollection() error = %v", err)
	}
	if col.ID != "col-123" {
		t.Errorf("ID = %v, want col-123", col.ID)
	}
}

func TestGetCollection(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || !strings.HasPrefix(r.URL.Path, "/collections/") {
			t.Errorf("Expected GET /collections/*, got %s %s", r.Method, r.URL.Path)
		}

		resp := &xaiv1.CollectionMetadata{
			CollectionId:   "col-123",
			CollectionName: "Test",
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	restClient := rest.NewClient(rest.Config{BaseURL: server.URL, APIKey: "test"})
	client := NewClient(restClient)

	col, err := client.GetCollection(context.Background(), "col-123", "team-1")

	if err != nil {
		t.Fatalf("GetCollection() error = %v", err)
	}
	if col.ID != "col-123" {
		t.Errorf("ID = %v, want col-123", col.ID)
	}
}

func TestListCollections(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := &xaiv1.ListCollectionsResponse{
			Collections: []*xaiv1.CollectionMetadata{
				{CollectionId: "col-1", CollectionName: "Collection 1"},
				{CollectionId: "col-2", CollectionName: "Collection 2"},
			},
			PaginationToken: "next-page",
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	restClient := rest.NewClient(rest.Config{BaseURL: server.URL, APIKey: "test"})
	client := NewClient(restClient)

	cols, token, err := client.ListCollections(context.Background(), &ListCollectionsOptions{
		TeamID: "team-1",
		Limit:  10,
	})

	if err != nil {
		t.Fatalf("ListCollections() error = %v", err)
	}
	if len(cols) != 2 {
		t.Errorf("len(cols) = %v, want 2", len(cols))
	}
	if token != "next-page" {
		t.Errorf("token = %v, want next-page", token)
	}
}

func TestUpdateCollection(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("Expected PUT, got %s", r.Method)
		}

		resp := &xaiv1.CollectionMetadata{
			CollectionId:   "col-123",
			CollectionName: "Updated",
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	restClient := rest.NewClient(rest.Config{BaseURL: server.URL, APIKey: "test"})
	client := NewClient(restClient)

	col, err := client.UpdateCollection(context.Background(), "col-123", "team-1", CreateCollectionOptions{
		Name: "Updated",
	})

	if err != nil {
		t.Fatalf("UpdateCollection() error = %v", err)
	}
	if col.Name != "Updated" {
		t.Errorf("Name = %v, want Updated", col.Name)
	}
}

func TestDeleteCollection(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("Expected DELETE, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	restClient := rest.NewClient(rest.Config{BaseURL: server.URL, APIKey: "test"})
	client := NewClient(restClient)

	err := client.DeleteCollection(context.Background(), "col-123", "team-1")
	if err != nil {
		t.Errorf("DeleteCollection() error = %v", err)
	}
}

func TestAddDocument(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := &xaiv1.DocumentMetadata{
			FileMetadata: &xaiv1.FileMetadata{
				FileId: "file-123",
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	restClient := rest.NewClient(rest.Config{BaseURL: server.URL, APIKey: "test"})
	client := NewClient(restClient)

	doc, err := client.AddDocument(context.Background(), AddDocumentOptions{
		FileID:       "file-123",
		TeamID:       "team-1",
		CollectionID: "col-123",
		Fields:       map[string]string{"key": "value"},
	})

	if err != nil {
		t.Fatalf("AddDocument() error = %v", err)
	}
	if doc.FileID != "file-123" {
		t.Errorf("FileID = %v, want file-123", doc.FileID)
	}
}

func TestGetDocument(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := &xaiv1.DocumentMetadata{
			FileMetadata: &xaiv1.FileMetadata{
				FileId: "file-123",
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	restClient := rest.NewClient(rest.Config{BaseURL: server.URL, APIKey: "test"})
	client := NewClient(restClient)

	doc, err := client.GetDocument(context.Background(), "col-123", "file-123", "team-1")

	if err != nil {
		t.Fatalf("GetDocument() error = %v", err)
	}
	if doc.FileID != "file-123" {
		t.Errorf("FileID = %v, want file-123", doc.FileID)
	}
}

func TestListDocuments(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := &xaiv1.ListDocumentsResponse{
			Documents: []*xaiv1.DocumentMetadata{
				{FileMetadata: &xaiv1.FileMetadata{FileId: "file-1"}},
				{FileMetadata: &xaiv1.FileMetadata{FileId: "file-2"}},
			},
			PaginationToken: "next",
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	restClient := rest.NewClient(rest.Config{BaseURL: server.URL, APIKey: "test"})
	client := NewClient(restClient)

	docs, token, err := client.ListDocuments(context.Background(), &ListDocumentsOptions{
		CollectionID: "col-123",
		TeamID:       "team-1",
	})

	if err != nil {
		t.Fatalf("ListDocuments() error = %v", err)
	}
	if len(docs) != 2 {
		t.Errorf("len(docs) = %v, want 2", len(docs))
	}
	if token != "next" {
		t.Errorf("token = %v, want next", token)
	}
}

func TestUpdateDocument(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := &xaiv1.DocumentMetadata{
			FileMetadata: &xaiv1.FileMetadata{
				FileId: "file-123",
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	restClient := rest.NewClient(rest.Config{BaseURL: server.URL, APIKey: "test"})
	client := NewClient(restClient)

	doc, err := client.UpdateDocument(context.Background(), "col-123", "file-123", "team-1", map[string]string{"key": "value"})

	if err != nil {
		t.Fatalf("UpdateDocument() error = %v", err)
	}
	if doc.FileID != "file-123" {
		t.Errorf("FileID = %v, want file-123", doc.FileID)
	}
}

func TestDeleteDocument(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	restClient := rest.NewClient(rest.Config{BaseURL: server.URL, APIKey: "test"})
	client := NewClient(restClient)

	err := client.DeleteDocument(context.Background(), "col-123", "file-123", "team-1")
	if err != nil {
		t.Errorf("DeleteDocument() error = %v", err)
	}
}

func TestBatchGetDocuments(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := &xaiv1.BatchGetDocumentsResponse{
			Documents: []*xaiv1.DocumentMetadata{
				{FileMetadata: &xaiv1.FileMetadata{FileId: "file-1"}},
				{FileMetadata: &xaiv1.FileMetadata{FileId: "file-2"}},
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	restClient := rest.NewClient(rest.Config{BaseURL: server.URL, APIKey: "test"})
	client := NewClient(restClient)

	docs, err := client.BatchGetDocuments(context.Background(), "col-123", "team-1", []string{"file-1", "file-2"})

	if err != nil {
		t.Fatalf("BatchGetDocuments() error = %v", err)
	}
	if len(docs) != 2 {
		t.Errorf("len(docs) = %v, want 2", len(docs))
	}
}

func TestClientNotInitialized(t *testing.T) {
	client := &Client{restClient: nil}

	_, err := client.CreateCollection(context.Background(), CreateCollectionOptions{})
	if err != ErrClientNotInitialized {
		t.Errorf("CreateCollection() error = %v, want %v", err, ErrClientNotInitialized)
	}

	_, err = client.GetCollection(context.Background(), "col", "team")
	if err != ErrClientNotInitialized {
		t.Errorf("GetCollection() error = %v, want %v", err, ErrClientNotInitialized)
	}

	_, _, err = client.ListCollections(context.Background(), nil)
	if err != ErrClientNotInitialized {
		t.Errorf("ListCollections() error = %v, want %v", err, ErrClientNotInitialized)
	}
}

func TestFromProtoHelpers(t *testing.T) {
	now := timestamppb.Now()

	// Test fromProtoCollection
	protoCol := &xaiv1.CollectionMetadata{
		CollectionId:   "col-123",
		CollectionName: "Test",
		CreatedAt:      now,
		DocumentsCount: 5,
	}
	col := fromProtoCollection(protoCol)
	if col.ID != "col-123" || col.Name != "Test" || col.DocumentsCount != 5 {
		t.Error("fromProtoCollection conversion failed")
	}

	// Test fromProtoDocument
	protoDoc := &xaiv1.DocumentMetadata{
		FileMetadata: &xaiv1.FileMetadata{
			FileId: "file-123",
		},
	}
	doc := fromProtoDocument(protoDoc)
	if doc.FileID != "file-123" {
		t.Error("fromProtoDocument conversion failed")
	}
}
