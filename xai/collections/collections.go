// Package collections provides document collection management functionality for xAI SDK.
package collections

import (
	"context"
	"fmt"
	"time"

	xaiv1 "github.com/ZaguanLabs/xai-sdk-go/proto/gen/go/xai/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Collection represents a document collection.
type Collection struct {
	id            string
	name          string
	description   string
	createdAt     time.Time
	updatedAt     time.Time
	documentCount int32
	totalSize     int64
}

// Document represents a document in a collection.
type Document struct {
	id           string
	collectionID string
	title        string
	content      string
	contentType  string
	size         int64
	createdAt    time.Time
	updatedAt    time.Time
}

// CollectionsServiceClient is an interface for the collections service client.
type CollectionsServiceClient interface {
	ListCollections(ctx context.Context, req *xaiv1.ListCollectionsRequest, opts ...grpc.CallOption) (*xaiv1.ListCollectionsResponse, error)
	GetCollection(ctx context.Context, req *xaiv1.GetCollectionRequest, opts ...grpc.CallOption) (*xaiv1.Collection, error)
	CreateCollection(ctx context.Context, req *xaiv1.CreateCollectionRequest, opts ...grpc.CallOption) (*xaiv1.Collection, error)
	DeleteCollection(ctx context.Context, req *xaiv1.DeleteCollectionRequest, opts ...grpc.CallOption) (*xaiv1.DeleteCollectionResponse, error)
	ListDocuments(ctx context.Context, req *xaiv1.ListDocumentsRequest, opts ...grpc.CallOption) (*xaiv1.ListDocumentsResponse, error)
	GetDocument(ctx context.Context, req *xaiv1.GetDocumentRequest, opts ...grpc.CallOption) (*xaiv1.Document, error)
	AddDocument(ctx context.Context, req *xaiv1.AddDocumentRequest, opts ...grpc.CallOption) (*xaiv1.Document, error)
	DeleteDocument(ctx context.Context, req *xaiv1.DeleteDocumentRequest, opts ...grpc.CallOption) (*xaiv1.DeleteDocumentResponse, error)
}

// Client provides collections management functionality.
type Client struct {
	grpcClient CollectionsServiceClient
}

// NewClient creates a new collections client.
func NewClient(grpcClient CollectionsServiceClient) *Client {
	return &Client{
		grpcClient: grpcClient,
	}
}

// ListCollections lists all collections with optional filtering.
func (c *Client) ListCollections(ctx context.Context, purpose string, limit int32, sortBy, sortOrder string) ([]*Collection, error) {
	req := &xaiv1.ListCollectionsRequest{
		Purpose:   purpose,
		Limit:     limit,
		SortBy:    sortBy,
		SortOrder: sortOrder,
	}

	resp, err := c.grpcClient.ListCollections(ctx, req)
	if err != nil {
		// Handle specific gRPC errors
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.Unauthenticated:
				return nil, fmt.Errorf("authentication failed: %s", st.Message())
			case codes.PermissionDenied:
				return nil, fmt.Errorf("permission denied: %s", st.Message())
			case codes.Unavailable:
				return nil, fmt.Errorf("service unavailable: %s", st.Message())
			default:
				return nil, fmt.Errorf("list collections failed (%s): %s", st.Code().String(), st.Message())
			}
		}
		return nil, fmt.Errorf("list collections failed: %w", err)
	}

	if resp == nil {
		return nil, fmt.Errorf("received nil response")
	}

	// Convert response collections
	collections := make([]*Collection, 0, len(resp.Collections))
	for _, collProto := range resp.Collections {
		collections = append(collections, &Collection{
			id:            collProto.Id,
			name:          collProto.Name,
			description:   collProto.Description,
			createdAt:     time.Unix(collProto.CreatedAt, 0),
			updatedAt:     time.Unix(collProto.UpdatedAt, 0),
			documentCount: collProto.DocumentCount,
			totalSize:     collProto.TotalSize,
		})
	}

	return collections, nil
}

// GetCollection retrieves information about a specific collection.
func (c *Client) GetCollection(ctx context.Context, collectionID string) (*Collection, error) {
	if collectionID == "" {
		return nil, fmt.Errorf("collection ID is required")
	}

	req := &xaiv1.GetCollectionRequest{
		CollectionId: collectionID,
	}

	resp, err := c.grpcClient.GetCollection(ctx, req)
	if err != nil {
		// Handle specific gRPC errors
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.Unauthenticated:
				return nil, fmt.Errorf("authentication failed: %s", st.Message())
			case codes.PermissionDenied:
				return nil, fmt.Errorf("permission denied: %s", st.Message())
			case codes.NotFound:
				return nil, fmt.Errorf("collection not found: %s", st.Message())
			case codes.Unavailable:
				return nil, fmt.Errorf("service unavailable: %s", st.Message())
			default:
				return nil, fmt.Errorf("get collection failed (%s): %s", st.Code().String(), st.Message())
			}
		}
		return nil, fmt.Errorf("get collection failed: %w", err)
	}

	if resp == nil {
		return nil, fmt.Errorf("received nil response")
	}

	return &Collection{
		id:            resp.Id,
		name:          resp.Name,
		description:   resp.Description,
		createdAt:     time.Unix(resp.CreatedAt, 0),
		updatedAt:     time.Unix(resp.UpdatedAt, 0),
		documentCount: resp.DocumentCount,
		totalSize:     resp.TotalSize,
	}, nil
}

// CreateCollection creates a new collection.
func (c *Client) CreateCollection(ctx context.Context, name, description string, encrypted bool) (*Collection, error) {
	if name == "" {
		return nil, fmt.Errorf("collection name is required")
	}

	req := &xaiv1.CreateCollectionRequest{
		Name:        name,
		Description: description,
		Encrypted:   encrypted,
	}

	resp, err := c.grpcClient.CreateCollection(ctx, req)
	if err != nil {
		// Handle specific gRPC errors
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.Unauthenticated:
				return nil, fmt.Errorf("authentication failed: %s", st.Message())
			case codes.PermissionDenied:
				return nil, fmt.Errorf("permission denied: %s", st.Message())
			case codes.InvalidArgument:
				return nil, fmt.Errorf("invalid request: %s", st.Message())
			case codes.ResourceExhausted:
				return nil, fmt.Errorf("quota exceeded: %s", st.Message())
			case codes.Unavailable:
				return nil, fmt.Errorf("service unavailable: %s", st.Message())
			default:
				return nil, fmt.Errorf("create collection failed (%s): %s", st.Code().String(), st.Message())
			}
		}
		return nil, fmt.Errorf("create collection failed: %w", err)
	}

	if resp == nil {
		return nil, fmt.Errorf("received nil response")
	}

	return &Collection{
		id:            resp.Id,
		name:          resp.Name,
		description:   resp.Description,
		createdAt:     time.Unix(resp.CreatedAt, 0),
		updatedAt:     time.Unix(resp.UpdatedAt, 0),
		documentCount: resp.DocumentCount,
		totalSize:     resp.TotalSize,
	}, nil
}

// DeleteCollection deletes a collection.
func (c *Client) DeleteCollection(ctx context.Context, collectionID string) error {
	if collectionID == "" {
		return fmt.Errorf("collection ID is required")
	}

	req := &xaiv1.DeleteCollectionRequest{
		CollectionId: collectionID,
	}

	resp, err := c.grpcClient.DeleteCollection(ctx, req)
	if err != nil {
		// Handle specific gRPC errors
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.Unauthenticated:
				return fmt.Errorf("authentication failed: %s", st.Message())
			case codes.PermissionDenied:
				return fmt.Errorf("permission denied: %s", st.Message())
			case codes.NotFound:
				return fmt.Errorf("collection not found: %s", st.Message())
			case codes.Unavailable:
				return fmt.Errorf("service unavailable: %s", st.Message())
			default:
				return fmt.Errorf("delete collection failed (%s): %s", st.Code().String(), st.Message())
			}
		}
		return fmt.Errorf("delete collection failed: %w", err)
	}

	if resp == nil {
		return fmt.Errorf("received nil response")
	}

	if !resp.Success {
		return fmt.Errorf("collection deletion failed")
	}

	return nil
}

// ListDocuments lists documents in a collection.
func (c *Client) ListDocuments(ctx context.Context, collectionID string, limit int32, sortBy, sortOrder string) ([]*Document, error) {
	if collectionID == "" {
		return nil, fmt.Errorf("collection ID is required")
	}

	req := &xaiv1.ListDocumentsRequest{
		CollectionId: collectionID,
		Limit:        limit,
		SortBy:       sortBy,
		SortOrder:    sortOrder,
	}

	resp, err := c.grpcClient.ListDocuments(ctx, req)
	if err != nil {
		// Handle specific gRPC errors
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.Unauthenticated:
				return nil, fmt.Errorf("authentication failed: %s", st.Message())
			case codes.PermissionDenied:
				return nil, fmt.Errorf("permission denied: %s", st.Message())
			case codes.NotFound:
				return nil, fmt.Errorf("collection not found: %s", st.Message())
			case codes.Unavailable:
				return nil, fmt.Errorf("service unavailable: %s", st.Message())
			default:
				return nil, fmt.Errorf("list documents failed (%s): %s", st.Code().String(), st.Message())
			}
		}
		return nil, fmt.Errorf("list documents failed: %w", err)
	}

	if resp == nil {
		return nil, fmt.Errorf("received nil response")
	}

	// Convert response documents
	documents := make([]*Document, 0, len(resp.Documents))
	for _, docProto := range resp.Documents {
		documents = append(documents, &Document{
			id:           docProto.Id,
			collectionID: docProto.CollectionId,
			title:        docProto.Title,
			content:      docProto.Content,
			contentType:  docProto.ContentType,
			size:         docProto.Size,
			createdAt:    time.Unix(docProto.CreatedAt, 0),
			updatedAt:    time.Unix(docProto.UpdatedAt, 0),
		})
	}

	return documents, nil
}

// GetDocument retrieves a specific document.
func (c *Client) GetDocument(ctx context.Context, collectionID, documentID string) (*Document, error) {
	if collectionID == "" {
		return nil, fmt.Errorf("collection ID is required")
	}
	if documentID == "" {
		return nil, fmt.Errorf("document ID is required")
	}

	req := &xaiv1.GetDocumentRequest{
		CollectionId: collectionID,
		DocumentId:   documentID,
	}

	resp, err := c.grpcClient.GetDocument(ctx, req)
	if err != nil {
		// Handle specific gRPC errors
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.Unauthenticated:
				return nil, fmt.Errorf("authentication failed: %s", st.Message())
			case codes.PermissionDenied:
				return nil, fmt.Errorf("permission denied: %s", st.Message())
			case codes.NotFound:
				return nil, fmt.Errorf("document not found: %s", st.Message())
			case codes.Unavailable:
				return nil, fmt.Errorf("service unavailable: %s", st.Message())
			default:
				return nil, fmt.Errorf("get document failed (%s): %s", st.Code().String(), st.Message())
			}
		}
		return nil, fmt.Errorf("get document failed: %w", err)
	}

	if resp == nil {
		return nil, fmt.Errorf("received nil response")
	}

	return &Document{
		id:           resp.Id,
		collectionID: resp.CollectionId,
		title:        resp.Title,
		content:      resp.Content,
		contentType:  resp.ContentType,
		size:         resp.Size,
		createdAt:    time.Unix(resp.CreatedAt, 0),
		updatedAt:    time.Unix(resp.UpdatedAt, 0),
	}, nil
}

// AddDocument adds a document to a collection.
func (c *Client) AddDocument(ctx context.Context, collectionID, title, content, contentType string) (*Document, error) {
	if collectionID == "" {
		return nil, fmt.Errorf("collection ID is required")
	}
	if title == "" {
		return nil, fmt.Errorf("document title is required")
	}
	if content == "" {
		return nil, fmt.Errorf("document content is required")
	}

	req := &xaiv1.AddDocumentRequest{
		CollectionId: collectionID,
		Title:        title,
		Content:      content,
		ContentType:  contentType,
	}

	resp, err := c.grpcClient.AddDocument(ctx, req)
	if err != nil {
		// Handle specific gRPC errors
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.Unauthenticated:
				return nil, fmt.Errorf("authentication failed: %s", st.Message())
			case codes.PermissionDenied:
				return nil, fmt.Errorf("permission denied: %s", st.Message())
			case codes.NotFound:
				return nil, fmt.Errorf("collection not found: %s", st.Message())
			case codes.ResourceExhausted:
				return nil, fmt.Errorf("quota exceeded: %s", st.Message())
			case codes.Unavailable:
				return nil, fmt.Errorf("service unavailable: %s", st.Message())
			default:
				return nil, fmt.Errorf("add document failed (%s): %s", st.Code().String(), st.Message())
			}
		}
		return nil, fmt.Errorf("add document failed: %w", err)
	}

	if resp == nil {
		return nil, fmt.Errorf("received nil response")
	}

	return &Document{
		id:           resp.Id,
		collectionID: resp.CollectionId,
		title:        resp.Title,
		content:      resp.Content,
		contentType:  resp.ContentType,
		size:         resp.Size,
		createdAt:    time.Unix(resp.CreatedAt, 0),
		updatedAt:    time.Unix(resp.UpdatedAt, 0),
	}, nil
}

// DeleteDocument deletes a document from a collection.
func (c *Client) DeleteDocument(ctx context.Context, collectionID, documentID string) error {
	if collectionID == "" {
		return fmt.Errorf("collection ID is required")
	}
	if documentID == "" {
		return fmt.Errorf("document ID is required")
	}

	req := &xaiv1.DeleteDocumentRequest{
		CollectionId: collectionID,
		DocumentId:   documentID,
	}

	resp, err := c.grpcClient.DeleteDocument(ctx, req)
	if err != nil {
		// Handle specific gRPC errors
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.Unauthenticated:
				return fmt.Errorf("authentication failed: %s", st.Message())
			case codes.PermissionDenied:
				return fmt.Errorf("permission denied: %s", st.Message())
			case codes.NotFound:
				return fmt.Errorf("document not found: %s", st.Message())
			case codes.Unavailable:
				return fmt.Errorf("service unavailable: %s", st.Message())
			default:
				return fmt.Errorf("delete document failed (%s): %s", st.Code().String(), st.Message())
			}
		}
		return fmt.Errorf("delete document failed: %w", err)
	}

	if resp == nil {
		return fmt.Errorf("received nil response")
	}

	if !resp.Success {
		return fmt.Errorf("document deletion failed")
	}

	return nil
}

// Collection methods

// ID returns the collection ID.
func (c *Collection) ID() string {
	return c.id
}

// Name returns the collection name.
func (c *Collection) Name() string {
	return c.name
}

// Description returns the collection description.
func (c *Collection) Description() string {
	return c.description
}

// CreatedAt returns the collection creation time.
func (c *Collection) CreatedAt() time.Time {
	return c.createdAt
}

// UpdatedAt returns the collection update time.
func (c *Collection) UpdatedAt() time.Time {
	return c.updatedAt
}

// DocumentCount returns the number of documents in the collection.
func (c *Collection) DocumentCount() int32 {
	return c.documentCount
}

// TotalSize returns the total size of documents in the collection.
func (c *Collection) TotalSize() int64 {
	return c.totalSize
}

// String returns a string representation of the collection.
func (c *Collection) String() string {
	return fmt.Sprintf("Collection{ID: %s, Name: %s, Documents: %d, Size: %d bytes}",
		c.id, c.name, c.documentCount, c.totalSize)
}

// Document methods

// ID returns the document ID.
func (d *Document) ID() string {
	return d.id
}

// CollectionID returns the collection ID.
func (d *Document) CollectionID() string {
	return d.collectionID
}

// Title returns the document title.
func (d *Document) Title() string {
	return d.title
}

// Content returns the document content.
func (d *Document) Content() string {
	return d.content
}

// ContentType returns the document content type.
func (d *Document) ContentType() string {
	return d.contentType
}

// Size returns the document size in bytes.
func (d *Document) Size() int64 {
	return d.size
}

// CreatedAt returns the document creation time.
func (d *Document) CreatedAt() time.Time {
	return d.createdAt
}

// UpdatedAt returns the document update time.
func (d *Document) UpdatedAt() time.Time {
	return d.updatedAt
}

// String returns a string representation of the document.
func (d *Document) String() string {
	return fmt.Sprintf("Document{ID: %s, Title: %s, Size: %d bytes, Type: %s}",
		d.id, d.title, d.size, d.contentType)
}