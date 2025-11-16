// Package files provides a client for the xAI Files API.
package files

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"time"

	xaiv1 "github.com/ZaguanLabs/xai-sdk-go/proto/gen/go/xai/v1"
	"github.com/ZaguanLabs/xai-sdk-go/xai/internal/rest"
	"google.golang.org/protobuf/encoding/protojson"
)

// Client provides access to the xAI Files API.
type Client struct {
	restClient *rest.Client
}

// NewClient creates a new Files API client.
func NewClient(restClient *rest.Client) *Client {
	return &Client{
		restClient: restClient,
	}
}

// File represents a file with metadata.
type File struct {
	ID        string
	Filename  string
	Size      int64
	CreatedAt time.Time
	ExpiresAt time.Time
	TeamID    string
}

// ListOptions contains options for listing files.
type ListOptions struct {
	Limit           int32
	Order           xaiv1.FilesOrdering
	PaginationToken string
	SortBy          xaiv1.FilesSortBy
}

// ListResult contains the result of a list operation.
type ListResult struct {
	Files           []*File
	PaginationToken string
}

// UploadOptions contains options for uploading a file.
type UploadOptions struct {
	Name    string
	Purpose string
}

// fromProto converts a proto File to a File.
func fromProto(pf *xaiv1.File) *File {
	if pf == nil {
		return nil
	}

	f := &File{
		ID:       pf.Id,
		Filename: pf.Filename,
		Size:     pf.Size,
		TeamID:   pf.TeamId,
	}

	if pf.CreatedAt != nil {
		f.CreatedAt = pf.CreatedAt.AsTime()
	}

	if pf.ExpiresAt != nil {
		f.ExpiresAt = pf.ExpiresAt.AsTime()
	}

	return f
}

// Upload uploads a file.
func (c *Client) Upload(ctx context.Context, reader io.Reader, opts UploadOptions) (*File, error) {
	if c.restClient == nil {
		return nil, ErrClientNotInitialized
	}

	// Read file content
	content, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read file content: %w", err)
	}

	// Create upload request
	req := &xaiv1.UploadFileChunk{
		Init: &xaiv1.UploadFileInit{
			Name:    opts.Name,
			Purpose: opts.Purpose,
		},
		Data: content,
	}

	jsonData, err := protojson.Marshal(req)
	if err != nil {
		return nil, err
	}

	resp, err := c.restClient.Post(ctx, "/files", jsonData)
	if err != nil {
		return nil, err
	}

	var file xaiv1.File
	if err := protojson.Unmarshal(resp.Body, &file); err != nil {
		return nil, err
	}

	return fromProto(&file), nil
}

// Download downloads a file's content.
func (c *Client) Download(ctx context.Context, fileID string) (io.ReadCloser, error) {
	if c.restClient == nil {
		return nil, ErrClientNotInitialized
	}

	resp, err := c.restClient.Get(ctx, fmt.Sprintf("/files/%s/content", fileID))
	if err != nil {
		return nil, err
	}

	// Parse the response to get file chunks
	var chunks xaiv1.FileContentChunk
	if err := protojson.Unmarshal(resp.Body, &chunks); err != nil {
		return nil, err
	}

	// Return the data as a ReadCloser
	return io.NopCloser(bytes.NewReader(chunks.Data)), nil
}

// List lists files with optional filtering and pagination.
func (c *Client) List(ctx context.Context, opts *ListOptions) (*ListResult, error) {
	if c.restClient == nil {
		return nil, ErrClientNotInitialized
	}

	req := &xaiv1.ListFilesRequest{}
	if opts != nil {
		req.Limit = opts.Limit
		req.Order = opts.Order
		req.PaginationToken = opts.PaginationToken
		req.SortBy = opts.SortBy
	}

	jsonData, err := protojson.Marshal(req)
	if err != nil {
		return nil, err
	}

	resp, err := c.restClient.Post(ctx, "/files/list", jsonData)
	if err != nil {
		return nil, err
	}

	var listResp xaiv1.ListFilesResponse
	if err := protojson.Unmarshal(resp.Body, &listResp); err != nil {
		return nil, err
	}

	files := make([]*File, len(listResp.Data))
	for i, f := range listResp.Data {
		files[i] = fromProto(f)
	}

	return &ListResult{
		Files:           files,
		PaginationToken: listResp.PaginationToken,
	}, nil
}

// Get retrieves file metadata by ID.
func (c *Client) Get(ctx context.Context, fileID string) (*File, error) {
	if c.restClient == nil {
		return nil, ErrClientNotInitialized
	}

	resp, err := c.restClient.Get(ctx, fmt.Sprintf("/files/%s", fileID))
	if err != nil {
		return nil, err
	}

	var file xaiv1.File
	if err := protojson.Unmarshal(resp.Body, &file); err != nil {
		return nil, err
	}

	return fromProto(&file), nil
}

// GetURL retrieves a temporary URL for downloading a file.
func (c *Client) GetURL(ctx context.Context, fileID string) (string, error) {
	if c.restClient == nil {
		return "", ErrClientNotInitialized
	}

	resp, err := c.restClient.Get(ctx, fmt.Sprintf("/files/%s/url", fileID))
	if err != nil {
		return "", err
	}

	var urlResp xaiv1.RetrieveFileURLResponse
	if err := protojson.Unmarshal(resp.Body, &urlResp); err != nil {
		return "", err
	}

	return urlResp.Url, nil
}

// Delete deletes a file by ID.
func (c *Client) Delete(ctx context.Context, fileID string) error {
	if c.restClient == nil {
		return ErrClientNotInitialized
	}

	_, err := c.restClient.Delete(ctx, fmt.Sprintf("/files/%s", fileID))
	return err
}
