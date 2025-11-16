// Package collections provides a client for the xAI Collections API.
package collections

import (
	"context"
	"fmt"
	"time"

	xaiv1 "github.com/ZaguanLabs/xai-sdk-go/proto/gen/go/xai/v1"
	"github.com/ZaguanLabs/xai-sdk-go/xai/internal/rest"
	"google.golang.org/protobuf/encoding/protojson"
)

// Client provides access to the xAI Collections API.
type Client struct {
	restClient *rest.Client
}

// NewClient creates a new Collections API client.
func NewClient(restClient *rest.Client) *Client {
	return &Client{
		restClient: restClient,
	}
}

// Collection represents a document collection.
type Collection struct {
	ID             string
	Name           string
	CreatedAt      time.Time
	DocumentsCount int32
}

// Document represents a document in a collection.
type Document struct {
	FileID      string
	Name        string
	SizeBytes   int64
	ContentType string
	CreatedAt   time.Time
	ExpiresAt   time.Time
	Hash        string
	Status      xaiv1.DocumentStatus
	ErrorMsg    string
	Fields      map[string]string
}

// IndexConfiguration contains index settings.
type IndexConfiguration struct {
	// Configuration fields from types.proto
}

// ChunkConfiguration contains chunking settings.
type ChunkConfiguration struct {
	// Configuration fields from types.proto
}

// CreateCollectionOptions contains options for creating a collection.
type CreateCollectionOptions struct {
	Name               string
	TeamID             string
	IndexConfiguration *IndexConfiguration
	ChunkConfiguration *ChunkConfiguration
}

// ListCollectionsOptions contains options for listing collections.
type ListCollectionsOptions struct {
	TeamID          string
	Limit           int32
	Order           xaiv1.Ordering
	PaginationToken string
	SortBy          xaiv1.CollectionsSortBy
}

// ListDocumentsOptions contains options for listing documents.
type ListDocumentsOptions struct {
	CollectionID    string
	TeamID          string
	Limit           int32
	Order           xaiv1.Ordering
	PaginationToken string
	SortBy          xaiv1.DocumentsSortBy
}

// AddDocumentOptions contains options for adding a document.
type AddDocumentOptions struct {
	FileID       string
	TeamID       string
	CollectionID string
	Fields       map[string]string
}

// CreateCollection creates a new collection.
func (c *Client) CreateCollection(ctx context.Context, opts CreateCollectionOptions) (*Collection, error) {
	if c.restClient == nil {
		return nil, ErrClientNotInitialized
	}

	req := &xaiv1.CreateCollectionRequest{
		TeamId:         opts.TeamID,
		CollectionName: opts.Name,
	}

	jsonData, err := protojson.Marshal(req)
	if err != nil {
		return nil, err
	}

	resp, err := c.restClient.Post(ctx, "/collections", jsonData)
	if err != nil {
		return nil, err
	}

	var collection xaiv1.CollectionMetadata
	if err := protojson.Unmarshal(resp.Body, &collection); err != nil {
		return nil, err
	}

	return fromProtoCollection(&collection), nil
}

// GetCollection retrieves a collection by ID.
func (c *Client) GetCollection(ctx context.Context, collectionID, teamID string) (*Collection, error) {
	if c.restClient == nil {
		return nil, ErrClientNotInitialized
	}

	resp, err := c.restClient.Get(ctx, fmt.Sprintf("/collections/%s?team_id=%s", collectionID, teamID))
	if err != nil {
		return nil, err
	}

	var collection xaiv1.CollectionMetadata
	if err := protojson.Unmarshal(resp.Body, &collection); err != nil {
		return nil, err
	}

	return fromProtoCollection(&collection), nil
}

// ListCollections lists collections with optional filtering and pagination.
func (c *Client) ListCollections(ctx context.Context, opts *ListCollectionsOptions) ([]*Collection, string, error) {
	if c.restClient == nil {
		return nil, "", ErrClientNotInitialized
	}

	req := &xaiv1.ListCollectionsRequest{}
	if opts != nil {
		req.TeamId = opts.TeamID
		req.Limit = opts.Limit
		req.Order = opts.Order
		req.PaginationToken = opts.PaginationToken
		req.SortBy = opts.SortBy
	}

	jsonData, err := protojson.Marshal(req)
	if err != nil {
		return nil, "", err
	}

	resp, err := c.restClient.Post(ctx, "/collections/list", jsonData)
	if err != nil {
		return nil, "", err
	}

	var listResp xaiv1.ListCollectionsResponse
	if err := protojson.Unmarshal(resp.Body, &listResp); err != nil {
		return nil, "", err
	}

	collections := make([]*Collection, len(listResp.Collections))
	for i, col := range listResp.Collections {
		collections[i] = fromProtoCollection(col)
	}

	return collections, listResp.PaginationToken, nil
}

// UpdateCollection updates a collection's configuration.
func (c *Client) UpdateCollection(ctx context.Context, collectionID, teamID string, opts CreateCollectionOptions) (*Collection, error) {
	if c.restClient == nil {
		return nil, ErrClientNotInitialized
	}

	req := &xaiv1.UpdateCollectionRequest{
		CollectionId:   collectionID,
		TeamId:         teamID,
		CollectionName: opts.Name,
	}

	jsonData, err := protojson.Marshal(req)
	if err != nil {
		return nil, err
	}

	resp, err := c.restClient.Put(ctx, fmt.Sprintf("/collections/%s", collectionID), jsonData)
	if err != nil {
		return nil, err
	}

	var collection xaiv1.CollectionMetadata
	if err := protojson.Unmarshal(resp.Body, &collection); err != nil {
		return nil, err
	}

	return fromProtoCollection(&collection), nil
}

// DeleteCollection deletes a collection.
func (c *Client) DeleteCollection(ctx context.Context, collectionID, teamID string) error {
	if c.restClient == nil {
		return ErrClientNotInitialized
	}

	_, err := c.restClient.Delete(ctx, fmt.Sprintf("/collections/%s?team_id=%s", collectionID, teamID))
	return err
}

// AddDocument adds a document to a collection.
func (c *Client) AddDocument(ctx context.Context, opts AddDocumentOptions) (*Document, error) {
	if c.restClient == nil {
		return nil, ErrClientNotInitialized
	}

	fields := make([]*xaiv1.FieldsEntry, 0, len(opts.Fields))
	for k, v := range opts.Fields {
		fields = append(fields, &xaiv1.FieldsEntry{Key: k, Value: v})
	}

	req := &xaiv1.AddDocumentToCollectionRequest{
		FileId:       opts.FileID,
		TeamId:       opts.TeamID,
		CollectionId: opts.CollectionID,
		Fields:       fields,
	}

	jsonData, err := protojson.Marshal(req)
	if err != nil {
		return nil, err
	}

	resp, err := c.restClient.Post(ctx, fmt.Sprintf("/collections/%s/documents", opts.CollectionID), jsonData)
	if err != nil {
		return nil, err
	}

	var doc xaiv1.DocumentMetadata
	if err := protojson.Unmarshal(resp.Body, &doc); err != nil {
		return nil, err
	}

	return fromProtoDocument(&doc), nil
}

// GetDocument retrieves a document by ID.
func (c *Client) GetDocument(ctx context.Context, collectionID, fileID, teamID string) (*Document, error) {
	if c.restClient == nil {
		return nil, ErrClientNotInitialized
	}

	resp, err := c.restClient.Get(ctx, fmt.Sprintf("/collections/%s/documents/%s?team_id=%s", collectionID, fileID, teamID))
	if err != nil {
		return nil, err
	}

	var doc xaiv1.DocumentMetadata
	if err := protojson.Unmarshal(resp.Body, &doc); err != nil {
		return nil, err
	}

	return fromProtoDocument(&doc), nil
}

// ListDocuments lists documents in a collection.
func (c *Client) ListDocuments(ctx context.Context, opts *ListDocumentsOptions) ([]*Document, string, error) {
	if c.restClient == nil {
		return nil, "", ErrClientNotInitialized
	}

	req := &xaiv1.ListDocumentsRequest{}
	if opts != nil {
		req.CollectionId = opts.CollectionID
		req.TeamId = opts.TeamID
		req.Limit = opts.Limit
		req.Order = opts.Order
		req.PaginationToken = opts.PaginationToken
		req.SortBy = opts.SortBy
	}

	jsonData, err := protojson.Marshal(req)
	if err != nil {
		return nil, "", err
	}

	resp, err := c.restClient.Post(ctx, fmt.Sprintf("/collections/%s/documents/list", opts.CollectionID), jsonData)
	if err != nil {
		return nil, "", err
	}

	var listResp xaiv1.ListDocumentsResponse
	if err := protojson.Unmarshal(resp.Body, &listResp); err != nil {
		return nil, "", err
	}

	docs := make([]*Document, len(listResp.Documents))
	for i, d := range listResp.Documents {
		docs[i] = fromProtoDocument(d)
	}

	return docs, listResp.PaginationToken, nil
}

// UpdateDocument updates a document's fields.
func (c *Client) UpdateDocument(ctx context.Context, collectionID, fileID, teamID string, fields map[string]string) (*Document, error) {
	if c.restClient == nil {
		return nil, ErrClientNotInitialized
	}

	fieldsProto := make([]*xaiv1.FieldsEntry, 0, len(fields))
	for k, v := range fields {
		fieldsProto = append(fieldsProto, &xaiv1.FieldsEntry{Key: k, Value: v})
	}

	req := &xaiv1.UpdateDocumentRequest{
		CollectionId: collectionID,
		FileId:       fileID,
		TeamId:       teamID,
		Fields:       fieldsProto,
	}

	jsonData, err := protojson.Marshal(req)
	if err != nil {
		return nil, err
	}

	resp, err := c.restClient.Put(ctx, fmt.Sprintf("/collections/%s/documents/%s", collectionID, fileID), jsonData)
	if err != nil {
		return nil, err
	}

	var doc xaiv1.DocumentMetadata
	if err := protojson.Unmarshal(resp.Body, &doc); err != nil {
		return nil, err
	}

	return fromProtoDocument(&doc), nil
}

// DeleteDocument removes a document from a collection.
func (c *Client) DeleteDocument(ctx context.Context, collectionID, fileID, teamID string) error {
	if c.restClient == nil {
		return ErrClientNotInitialized
	}

	_, err := c.restClient.Delete(ctx, fmt.Sprintf("/collections/%s/documents/%s?team_id=%s", collectionID, fileID, teamID))
	return err
}

// BatchGetDocuments retrieves multiple documents at once.
func (c *Client) BatchGetDocuments(ctx context.Context, collectionID, teamID string, fileIDs []string) ([]*Document, error) {
	if c.restClient == nil {
		return nil, ErrClientNotInitialized
	}

	req := &xaiv1.BatchGetDocumentsRequest{
		TeamId:       teamID,
		CollectionId: collectionID,
		FileIds:      fileIDs,
	}

	jsonData, err := protojson.Marshal(req)
	if err != nil {
		return nil, err
	}

	resp, err := c.restClient.Post(ctx, fmt.Sprintf("/collections/%s/documents/batch", collectionID), jsonData)
	if err != nil {
		return nil, err
	}

	var batchResp xaiv1.BatchGetDocumentsResponse
	if err := protojson.Unmarshal(resp.Body, &batchResp); err != nil {
		return nil, err
	}

	docs := make([]*Document, len(batchResp.Documents))
	for i, d := range batchResp.Documents {
		docs[i] = fromProtoDocument(d)
	}

	return docs, nil
}
