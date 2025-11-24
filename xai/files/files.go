// Package files provides a client for the xAI Files API.
package files

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"sync"
	"time"

	xaiv1 "github.com/ZaguanLabs/xai-sdk-go/proto/gen/go/xai/api/v1"
	"github.com/ZaguanLabs/xai-sdk-go/xai/internal/constants"
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
	// MaxSize is the maximum file size in bytes. If 0, defaults to DefaultMaxFileSize (100MB).
	MaxSize int64
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

	// Determine max size (use default if not specified)
	maxSize := opts.MaxSize
	if maxSize <= 0 {
		maxSize = constants.DefaultMaxFileSize
	}

	// Limit file size to prevent memory exhaustion
	limitedReader := io.LimitReader(reader, maxSize+1)

	// Read file content
	content, err := io.ReadAll(limitedReader)
	if err != nil {
		return nil, fmt.Errorf("failed to read file content: %w", err)
	}

	// Check if file size exceeds the limit
	if int64(len(content)) > maxSize {
		return nil, fmt.Errorf("%w: file size %d exceeds limit %d", ErrFileTooLarge, len(content), maxSize)
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

// BatchUploadCallback is called after each file upload completes (success or failure).
// The callback receives the file index, reader, and result (File or error).
type BatchUploadCallback func(index int, reader io.Reader, result interface{})

// BatchUploadResult contains the result of a single file upload in a batch.
type BatchUploadResult struct {
	Index int
	File  *File
	Error error
}

// BatchUpload uploads multiple files concurrently with controlled concurrency.
// Returns a map of file indices to results (File or error).
// This method handles partial failures gracefully - successful uploads are returned
// even if some uploads fail.
//
// Parameters:
//   - ctx: Context for cancellation and timeouts
//   - readers: Slice of io.Readers containing file data
//   - opts: Slice of UploadOptions (must match length of readers)
//   - batchSize: Maximum number of concurrent uploads (default: 50 if <= 0)
//   - callback: Optional callback invoked after each file completes
//
// Returns:
//   - map[int]*BatchUploadResult: Map of file indices to results
//   - error: Only returns error for invalid parameters, not upload failures
//
// Example:
//
//	results, err := client.BatchUpload(ctx, readers, opts, 10, func(idx int, r io.Reader, result interface{}) {
//	    if res, ok := result.(*BatchUploadResult); ok {
//	        if res.Error != nil {
//	            fmt.Printf("File %d failed: %v\n", idx, res.Error)
//	        } else {
//	            fmt.Printf("File %d uploaded: %s\n", idx, res.File.ID)
//	        }
//	    }
//	})
func (c *Client) BatchUpload(
	ctx context.Context,
	readers []io.Reader,
	opts []UploadOptions,
	batchSize int,
	callback BatchUploadCallback,
) (map[int]*BatchUploadResult, error) {
	if len(readers) == 0 {
		return nil, fmt.Errorf("readers cannot be empty - please provide at least one file to upload")
	}
	if len(opts) != len(readers) {
		return nil, fmt.Errorf("opts length (%d) must match readers length (%d)", len(opts), len(readers))
	}
	if batchSize <= 0 {
		batchSize = 50 // default
	}

	results := make(map[int]*BatchUploadResult)
	var mu sync.Mutex

	// Use semaphore pattern for concurrency control
	sem := make(chan struct{}, batchSize)
	var wg sync.WaitGroup

	for i := range readers {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			sem <- struct{}{}        // Acquire
			defer func() { <-sem }() // Release

			file, err := c.Upload(ctx, readers[idx], opts[idx])

			result := &BatchUploadResult{
				Index: idx,
				File:  file,
				Error: err,
			}

			mu.Lock()
			results[idx] = result
			mu.Unlock()

			if callback != nil {
				callback(idx, readers[idx], result)
			}
		}(i)
	}

	wg.Wait()
	return results, nil
}
