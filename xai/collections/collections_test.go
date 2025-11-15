package collections

import (
	"context"
	"testing"
	"time"

	xaiv1 "github.com/ZaguanLabs/xai-sdk-go/proto/gen/go/xai/v1"
	"google.golang.org/grpc"
)

// mockCollectionsServiceClient implements CollectionsServiceClient for testing
type mockCollectionsServiceClient struct {
	collections []*xaiv1.Collection
	documents   []*xaiv1.Document
	err         error
}

func (m *mockCollectionsServiceClient) ListCollections(ctx context.Context, req *xaiv1.ListCollectionsRequest, opts ...grpc.CallOption) (*xaiv1.ListCollectionsResponse, error) {
	if m.err != nil {
		return nil, m.err
	}
	return &xaiv1.ListCollectionsResponse{Collections: m.collections}, nil
}

func (m *mockCollectionsServiceClient) GetCollection(ctx context.Context, req *xaiv1.GetCollectionRequest, opts ...grpc.CallOption) (*xaiv1.Collection, error) {
	if m.err != nil {
		return nil, m.err
	}
	for _, coll := range m.collections {
		if coll.Id == req.CollectionId {
			return coll, nil
		}
	}
	return nil, nil
}

func (m *mockCollectionsServiceClient) CreateCollection(ctx context.Context, req *xaiv1.CreateCollectionRequest, opts ...grpc.CallOption) (*xaiv1.Collection, error) {
	if m.err != nil {
		return nil, m.err
	}
	return &xaiv1.Collection{
		Id:          "new-collection-id",
		Name:        req.Name,
		Description: req.Description,
		CreatedAt:   time.Now().Unix(),
		UpdatedAt:   time.Now().Unix(),
	}, nil
}

func (m *mockCollectionsServiceClient) DeleteCollection(ctx context.Context, req *xaiv1.DeleteCollectionRequest, opts ...grpc.CallOption) (*xaiv1.DeleteCollectionResponse, error) {
	if m.err != nil {
		return nil, m.err
	}
	return &xaiv1.DeleteCollectionResponse{Success: true}, nil
}

func (m *mockCollectionsServiceClient) ListDocuments(ctx context.Context, req *xaiv1.ListDocumentsRequest, opts ...grpc.CallOption) (*xaiv1.ListDocumentsResponse, error) {
	if m.err != nil {
		return nil, m.err
	}
	// Filter documents by collection ID
	var filtered []*xaiv1.Document
	for _, doc := range m.documents {
		if doc.CollectionId == req.CollectionId {
			filtered = append(filtered, doc)
		}
	}
	return &xaiv1.ListDocumentsResponse{Documents: filtered}, nil
}

func (m *mockCollectionsServiceClient) GetDocument(ctx context.Context, req *xaiv1.GetDocumentRequest, opts ...grpc.CallOption) (*xaiv1.Document, error) {
	if m.err != nil {
		return nil, m.err
	}
	for _, doc := range m.documents {
		if doc.Id == req.DocumentId && doc.CollectionId == req.CollectionId {
			return doc, nil
		}
	}
	return nil, nil
}

func (m *mockCollectionsServiceClient) AddDocument(ctx context.Context, req *xaiv1.AddDocumentRequest, opts ...grpc.CallOption) (*xaiv1.Document, error) {
	if m.err != nil {
		return nil, m.err
	}
	return &xaiv1.Document{
		Id:           "new-doc-id",
		CollectionId: req.CollectionId,
		Title:        req.Title,
		Content:      req.Content,
		ContentType:  req.ContentType,
		Size:         int64(len(req.Content)),
		CreatedAt:    time.Now().Unix(),
		UpdatedAt:    time.Now().Unix(),
	}, nil
}

func (m *mockCollectionsServiceClient) DeleteDocument(ctx context.Context, req *xaiv1.DeleteDocumentRequest, opts ...grpc.CallOption) (*xaiv1.DeleteDocumentResponse, error) {
	if m.err != nil {
		return nil, m.err
	}
	return &xaiv1.DeleteDocumentResponse{Success: true}, nil
}

func TestListCollections(t *testing.T) {
	mockClient := &mockCollectionsServiceClient{
		collections: []*xaiv1.Collection{
			{Id: "coll1", Name: "Collection 1", Description: "First collection"},
			{Id: "coll2", Name: "Collection 2", Description: "Second collection"},
		},
	}

	client := NewClient(mockClient)
	collections, err := client.ListCollections(context.Background(), "", 10, "name", "asc")

	if err != nil {
		t.Fatalf("ListCollections() returned error: %v", err)
	}

	if len(collections) != 2 {
		t.Errorf("Expected 2 collections, got %d", len(collections))
	}

	if collections[0].ID() != "coll1" {
		t.Errorf("Expected first collection ID to be 'coll1', got '%s'", collections[0].ID())
	}

	if collections[0].Name() != "Collection 1" {
		t.Errorf("Expected first collection name to be 'Collection 1', got '%s'", collections[0].Name())
	}
}

func TestGetCollection(t *testing.T) {
	mockClient := &mockCollectionsServiceClient{
		collections: []*xaiv1.Collection{
			{Id: "coll1", Name: "Collection 1", Description: "First collection"},
		},
	}

	client := NewClient(mockClient)
	collection, err := client.GetCollection(context.Background(), "coll1")

	if err != nil {
		t.Fatalf("GetCollection() returned error: %v", err)
	}

	if collection.ID() != "coll1" {
		t.Errorf("Expected collection ID to be 'coll1', got '%s'", collection.ID())
	}

	if collection.Name() != "Collection 1" {
		t.Errorf("Expected collection name to be 'Collection 1', got '%s'", collection.Name())
	}
}

func TestCreateCollection(t *testing.T) {
	mockClient := &mockCollectionsServiceClient{}
	client := NewClient(mockClient)

	collection, err := client.CreateCollection(context.Background(), "New Collection", "A new collection", false)

	if err != nil {
		t.Fatalf("CreateCollection() returned error: %v", err)
	}

	if collection.Name() != "New Collection" {
		t.Errorf("Expected collection name to be 'New Collection', got '%s'", collection.Name())
	}

	if collection.Description() != "A new collection" {
		t.Errorf("Expected collection description to be 'A new collection', got '%s'", collection.Description())
	}
}

func TestDeleteCollection(t *testing.T) {
	mockClient := &mockCollectionsServiceClient{}
	client := NewClient(mockClient)

	err := client.DeleteCollection(context.Background(), "coll1")

	if err != nil {
		t.Fatalf("DeleteCollection() returned error: %v", err)
	}
}

func TestListDocuments(t *testing.T) {
	mockClient := &mockCollectionsServiceClient{
		documents: []*xaiv1.Document{
			{Id: "doc1", CollectionId: "coll1", Title: "Document 1"},
			{Id: "doc2", CollectionId: "coll1", Title: "Document 2"},
			{Id: "doc3", CollectionId: "coll2", Title: "Document 3"},
		},
	}

	client := NewClient(mockClient)
	documents, err := client.ListDocuments(context.Background(), "coll1", 10, "title", "asc")

	if err != nil {
		t.Fatalf("ListDocuments() returned error: %v", err)
	}

	if len(documents) != 2 {
		t.Errorf("Expected 2 documents for collection coll1, got %d", len(documents))
	}

	if documents[0].Title() != "Document 1" {
		t.Errorf("Expected first document title to be 'Document 1', got '%s'", documents[0].Title())
	}
}

func TestAddDocument(t *testing.T) {
	mockClient := &mockCollectionsServiceClient{}
	client := NewClient(mockClient)

	document, err := client.AddDocument(context.Background(), "coll1", "New Document", "Document content", "text/plain")

	if err != nil {
		t.Fatalf("AddDocument() returned error: %v", err)
	}

	if document.Title() != "New Document" {
		t.Errorf("Expected document title to be 'New Document', got '%s'", document.Title())
	}

	if document.Content() != "Document content" {
		t.Errorf("Expected document content to be 'Document content', got '%s'", document.Content())
	}

	if document.ContentType() != "text/plain" {
		t.Errorf("Expected document content type to be 'text/plain', got '%s'", document.ContentType())
	}
}

func TestCollectionString(t *testing.T) {
	collection := &Collection{
		id:            "test-coll",
		name:          "Test Collection",
		documentCount: 5,
		totalSize:     1024,
	}

	expected := "Collection{ID: test-coll, Name: Test Collection, Documents: 5, Size: 1024 bytes}"
	if collection.String() != expected {
		t.Errorf("Expected String() to return '%s', got '%s'", expected, collection.String())
	}
}

func TestDocumentString(t *testing.T) {
	document := &Document{
		id:          "test-doc",
		title:       "Test Document",
		size:        512,
		contentType: "text/plain",
	}

	expected := "Document{ID: test-doc, Title: Test Document, Size: 512 bytes, Type: text/plain}"
	if document.String() != expected {
		t.Errorf("Expected String() to return '%s', got '%s'", expected, document.String())
	}
}