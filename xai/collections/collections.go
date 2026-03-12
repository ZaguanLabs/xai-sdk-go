// Package collections provides a client for the xAI Collections API.
package collections

import (
	"bytes"
	"context"
	"fmt"
	"time"

	xaiv1 "github.com/ZaguanLabs/xai-sdk-go/proto/gen/go/xai/api/v1"
	"github.com/ZaguanLabs/xai-sdk-go/xai/documents"
	"github.com/ZaguanLabs/xai-sdk-go/xai/files"
	"github.com/ZaguanLabs/xai-sdk-go/xai/internal/rest"
	"google.golang.org/protobuf/encoding/protojson"
)

// Client provides access to the xAI Collections API.
type Client struct {
	restClient *rest.Client
}

func stringPtr(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}

func int32Ptr(value int32) *int32 {
	return &value
}

func orderingPtr(value xaiv1.Ordering) *xaiv1.Ordering {
	return &value
}

func collectionsSortByPtr(value xaiv1.CollectionsSortBy) *xaiv1.CollectionsSortBy {
	return &value
}

func documentsSortByPtr(value xaiv1.DocumentsSortBy) *xaiv1.DocumentsSortBy {
	return &value
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

type SearchOptions struct {
	Limit         int32
	Instructions  string
	RetrievalMode string
}

type UploadDocumentOptions struct {
	Fields          map[string]string
	TeamID          string
	WaitForIndexing bool
	PollInterval    time.Duration
	Timeout         time.Duration
	MaxFileSize     int64
	FilePurpose     string
}

// CreateCollection creates a new collection.
func (c *Client) CreateCollection(ctx context.Context, opts CreateCollectionOptions) (*Collection, error) {
	if c.restClient == nil {
		return nil, ErrClientNotInitialized
	}

	req := &xaiv1.CreateCollectionRequest{
		TeamId:         stringPtr(opts.TeamID),
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

func (c *Client) Create(ctx context.Context, opts CreateCollectionOptions) (*Collection, error) {
	return c.CreateCollection(ctx, opts)
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

func (c *Client) Get(ctx context.Context, collectionID, teamID string) (*Collection, error) {
	return c.GetCollection(ctx, collectionID, teamID)
}

// ListCollections lists collections with optional filtering and pagination.
func (c *Client) ListCollections(ctx context.Context, opts *ListCollectionsOptions) ([]*Collection, string, error) {
	if c.restClient == nil {
		return nil, "", ErrClientNotInitialized
	}

	req := &xaiv1.ListCollectionsRequest{}
	if opts != nil {
		req.TeamId = stringPtr(opts.TeamID)
		req.Limit = int32Ptr(opts.Limit)
		req.Order = orderingPtr(opts.Order)
		req.PaginationToken = stringPtr(opts.PaginationToken)
		req.SortBy = collectionsSortByPtr(opts.SortBy)
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

	if listResp.PaginationToken == nil {
		return collections, "", nil
	}

	return collections, *listResp.PaginationToken, nil
}

func (c *Client) List(ctx context.Context, opts *ListCollectionsOptions) ([]*Collection, string, error) {
	return c.ListCollections(ctx, opts)
}

// UpdateCollection updates a collection's configuration.
func (c *Client) UpdateCollection(ctx context.Context, collectionID, teamID string, opts CreateCollectionOptions) (*Collection, error) {
	if c.restClient == nil {
		return nil, ErrClientNotInitialized
	}

	req := &xaiv1.UpdateCollectionRequest{
		CollectionId:   collectionID,
		TeamId:         stringPtr(teamID),
		CollectionName: stringPtr(opts.Name),
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

func (c *Client) Update(ctx context.Context, collectionID, teamID string, opts CreateCollectionOptions) (*Collection, error) {
	return c.UpdateCollection(ctx, collectionID, teamID, opts)
}

// DeleteCollection deletes a collection.
func (c *Client) DeleteCollection(ctx context.Context, collectionID, teamID string) error {
	if c.restClient == nil {
		return ErrClientNotInitialized
	}

	_, err := c.restClient.Delete(ctx, fmt.Sprintf("/collections/%s?team_id=%s", collectionID, teamID))
	return err
}

func (c *Client) Delete(ctx context.Context, collectionID, teamID string) error {
	return c.DeleteCollection(ctx, collectionID, teamID)
}

func (c *Client) Search(ctx context.Context, query string, collectionIDs []string, opts *SearchOptions) (*documents.SearchResponse, error) {
	if c.restClient == nil {
		return nil, ErrClientNotInitialized
	}

	req := documents.NewSearchRequest(query, collectionIDs...)
	if opts != nil {
		if opts.Limit > 0 {
			req.WithLimit(opts.Limit)
		}
		if opts.Instructions != "" {
			req.WithInstructions(opts.Instructions)
		}
		switch opts.RetrievalMode {
		case "hybrid", "":
			if opts.RetrievalMode == "hybrid" {
				req.WithHybridRetrieval()
			}
		case "semantic":
			req.WithSemanticRetrieval()
		case "keyword":
			req.WithKeywordRetrieval()
		default:
			return nil, fmt.Errorf("unsupported retrieval mode %q", opts.RetrievalMode)
		}
	}

	return documents.NewClient(c.restClient).Search(ctx, req)
}

func (c *Client) UploadDocument(ctx context.Context, collectionID, name string, data []byte, opts *UploadDocumentOptions) (*Document, error) {
	if c.restClient == nil {
		return nil, ErrClientNotInitialized
	}

	fileOpts := files.UploadOptions{Name: name}
	addOpts := AddDocumentOptions{CollectionID: collectionID}
	pollInterval := 2 * time.Second
	timeout := 10 * time.Minute

	if opts != nil {
		fileOpts.Purpose = opts.FilePurpose
		fileOpts.MaxSize = opts.MaxFileSize
		addOpts.Fields = opts.Fields
		addOpts.TeamID = opts.TeamID
		if opts.PollInterval > 0 {
			pollInterval = opts.PollInterval
		}
		if opts.Timeout > 0 {
			timeout = opts.Timeout
		}
	}

	uploaded, err := files.NewClient(c.restClient).Upload(ctx, bytes.NewReader(data), fileOpts)
	if err != nil {
		return nil, err
	}

	addOpts.FileID = uploaded.ID
	document, err := c.AddDocument(ctx, addOpts)
	if err != nil {
		return nil, err
	}

	if opts == nil || !opts.WaitForIndexing {
		return document, nil
	}

	return c.waitForIndexing(ctx, collectionID, uploaded.ID, addOpts.TeamID, pollInterval, timeout)
}

func (c *Client) waitForIndexing(ctx context.Context, collectionID, fileID, teamID string, pollInterval, timeout time.Duration) (*Document, error) {
	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	for {
		document, err := c.GetDocument(timeoutCtx, collectionID, fileID, teamID)
		if err != nil {
			return nil, err
		}

		switch document.Status {
		case xaiv1.DocumentStatus_DOCUMENT_STATUS_PROCESSED:
			return document, nil
		case xaiv1.DocumentStatus_DOCUMENT_STATUS_PROCESSING:
		case xaiv1.DocumentStatus_DOCUMENT_STATUS_FAILED:
			return nil, fmt.Errorf("document indexing failed: %s", document.ErrorMsg)
		default:
			return nil, fmt.Errorf("unknown document status: %s", document.Status.String())
		}

		select {
		case <-timeoutCtx.Done():
			return nil, timeoutCtx.Err()
		case <-ticker.C:
		}
	}
}

// AddDocument adds a document to a collection.
func (c *Client) AddDocument(ctx context.Context, opts AddDocumentOptions) (*Document, error) {
	if c.restClient == nil {
		return nil, ErrClientNotInitialized
	}

	req := &xaiv1.AddDocumentToCollectionRequest{
		FileId:       opts.FileID,
		TeamId:       stringPtr(opts.TeamID),
		CollectionId: opts.CollectionID,
		Fields:       opts.Fields,
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
	var collectionID string
	if opts != nil {
		req.CollectionId = opts.CollectionID
		req.TeamId = stringPtr(opts.TeamID)
		req.Limit = int32Ptr(opts.Limit)
		req.Order = orderingPtr(opts.Order)
		req.PaginationToken = stringPtr(opts.PaginationToken)
		req.SortBy = documentsSortByPtr(opts.SortBy)
		collectionID = opts.CollectionID
	}

	jsonData, err := protojson.Marshal(req)
	if err != nil {
		return nil, "", err
	}

	resp, err := c.restClient.Post(ctx, fmt.Sprintf("/collections/%s/documents/list", collectionID), jsonData)
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

	if listResp.PaginationToken == nil {
		return docs, "", nil
	}

	return docs, *listResp.PaginationToken, nil
}

// UpdateDocument updates a document's fields.
func (c *Client) UpdateDocument(ctx context.Context, collectionID, fileID, teamID string, fields map[string]string) (*Document, error) {
	if c.restClient == nil {
		return nil, ErrClientNotInitialized
	}

	req := &xaiv1.UpdateDocumentRequest{
		CollectionId: collectionID,
		FileId:       fileID,
		TeamId:       stringPtr(teamID),
		Fields:       fields,
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
		TeamId:       stringPtr(teamID),
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
