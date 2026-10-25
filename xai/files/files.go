// Package files provides a client for the xAI Files API.
package files

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"path/filepath"
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
	grpcClient xaiv1.FilesClient
}

// NewClient creates a new Files API client.
func NewClient(restClient *rest.Client) *Client {
	return &Client{
		restClient: restClient,
	}
}

// NewGRPCClient uses the native Files protocol and retains REST for legacy URL helpers.
func NewGRPCClient(client xaiv1.FilesClient, restClient *rest.Client) *Client {
	return &Client{grpcClient: client, restClient: restClient}
}

// File represents a file with metadata.
type File struct {
	ID                 string
	Filename           string
	Size               int64
	CreatedAt          time.Time
	ExpiresAt          time.Time
	TeamID             string
	PublicURL          string
	PublicURLExpiresAt time.Time
}

// ListOptions contains options for listing files.
type ListOptions struct {
	Limit           int32
	Order           xaiv1.FilesOrdering
	PaginationToken string
	SortBy          xaiv1.FilesSortBy
	Filter          string
}

// ListResult contains the result of a list operation.
type ListResult struct {
	Files           []*File
	PaginationToken string
}

// UploadOptions contains options for uploading a file.
type UploadOptions struct {
	Name string
	// Purpose is retained for legacy REST clients. Native uploads use expiry instead.
	Purpose      string
	ExpiresAfter *time.Duration
	OnProgress   func(uploaded, total int64)
	// MaxSize is the maximum file size in bytes. If 0, defaults to DefaultMaxFileSize (100MB).
	MaxSize int64
}

// PublicURLOptions configures public URL creation for stored generated assets.
type PublicURLOptions struct {
	ExpiresAfter time.Duration
}

// StorageOptions configures Files API persistence for generated images and videos.
type StorageOptions struct {
	Filename     string
	ExpiresAfter time.Duration
	PublicURL    *PublicURLOptions
}

// PublicURLWithDefaults requests a public URL using the API's default expiry behavior.
func PublicURLWithDefaults() *PublicURLOptions {
	return &PublicURLOptions{}
}

// PublicURLWithExpiry requests a public URL with an independent expiry.
func PublicURLWithExpiry(expiresAfter time.Duration) *PublicURLOptions {
	return &PublicURLOptions{ExpiresAfter: expiresAfter}
}

// Proto converts storage options into the generated protobuf form.
func (o *StorageOptions) Proto() *xaiv1.StorageOptions {
	if o == nil {
		return nil
	}
	pb := &xaiv1.StorageOptions{Filename: o.Filename}
	if o.ExpiresAfter != 0 {
		seconds := int64(o.ExpiresAfter.Seconds())
		pb.ExpiresAfter = &seconds
	}
	if o.PublicURL != nil {
		pb.PublicUrl = &xaiv1.PublicUrlOptions{}
		if o.PublicURL.ExpiresAfter != 0 {
			seconds := int64(o.PublicURL.ExpiresAfter.Seconds())
			pb.PublicUrl.ExpiresAfter = &seconds
		}
	}
	return pb
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
	}

	if pf.CreatedAt != nil {
		f.CreatedAt = pf.CreatedAt.AsTime()
	}

	if pf.ExpiresAt != nil {
		f.ExpiresAt = pf.ExpiresAt.AsTime()
	}
	if pf.GetPublicUrl() != "" {
		f.PublicURL = pf.GetPublicUrl()
	}
	if pf.PublicUrlExpiresAt != nil {
		f.PublicURLExpiresAt = pf.PublicUrlExpiresAt.AsTime()
	}

	return f
}

// Upload uploads a file.
func (c *Client) Upload(ctx context.Context, reader io.Reader, opts UploadOptions) (*File, error) {
	if c.grpcClient != nil {
		return c.uploadGRPC(ctx, reader, opts)
	}
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

	// Legacy JSON endpoint accepts both initialization and content in one request.
	jsonData, err := json.Marshal(map[string]any{"init": map[string]any{"name": opts.Name, "purpose": opts.Purpose}, "data": content})
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
	return c.DownloadWithFormat(ctx, fileID, nil)
}

func (c *Client) DownloadWithFormat(ctx context.Context, fileID string, format *xaiv1.DownloadFormat) (io.ReadCloser, error) {
	if c.grpcClient != nil {
		streamCtx, cancel := context.WithCancel(ctx)
		stream, err := c.grpcClient.RetrieveFileContent(streamCtx, &xaiv1.RetrieveFileContentRequest{FileId: fileID, Format: format})
		if err != nil {
			cancel()
			return nil, err
		}
		return &contentReader{recv: stream.Recv, cancel: cancel}, nil
	}
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

func (c *Client) Content(ctx context.Context, fileID string) ([]byte, error) {
	reader, err := c.Download(ctx, fileID)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = reader.Close()
	}()
	return io.ReadAll(reader)
}

// List lists files with optional filtering and pagination.
func (c *Client) List(ctx context.Context, opts *ListOptions) (*ListResult, error) {

	req := &xaiv1.ListFilesRequest{}
	if opts != nil {
		req.Limit = opts.Limit
		req.Order = opts.Order
		if opts.PaginationToken != "" {
			req.PaginationToken = &opts.PaginationToken
		}
		req.SortBy = &opts.SortBy
		if opts.Filter != "" {
			req.Filter = &opts.Filter
		}
	}

	var listResp *xaiv1.ListFilesResponse
	if c.grpcClient != nil {
		var err error
		listResp, err = c.grpcClient.ListFiles(ctx, req)
		if err != nil {
			return nil, err
		}
	} else {
		if c.restClient == nil {
			return nil, ErrClientNotInitialized
		}
		jsonData, err := protojson.Marshal(req)
		if err != nil {
			return nil, err
		}
		resp, err := c.restClient.Post(ctx, "/files/list", jsonData)
		if err != nil {
			return nil, err
		}
		listResp = &xaiv1.ListFilesResponse{}
		if err := protojson.Unmarshal(resp.Body, listResp); err != nil {
			return nil, err
		}
	}

	files := make([]*File, len(listResp.Data))
	for i, f := range listResp.Data {
		files[i] = fromProto(f)
	}

	return &ListResult{
		Files:           files,
		PaginationToken: listResp.GetPaginationToken(),
	}, nil
}

// CreatePublicURL creates a publicly shareable URL for a stored file.
func (c *Client) CreatePublicURL(ctx context.Context, fileID string, expiresAfter time.Duration) (*xaiv1.CreatePublicUrlResponse, error) {
	if c.grpcClient != nil {
		req := &xaiv1.CreatePublicUrlRequest{FileId: fileID}
		if expiresAfter != 0 {
			seconds := int64(expiresAfter.Seconds())
			req.ExpiresAfter = &seconds
		}
		return c.grpcClient.CreatePublicUrl(ctx, req)
	}
	if c.restClient == nil {
		return nil, ErrClientNotInitialized
	}

	req := &xaiv1.CreatePublicUrlRequest{FileId: fileID}
	if expiresAfter != 0 {
		seconds := int64(expiresAfter.Seconds())
		req.ExpiresAfter = &seconds
	}

	jsonData, err := protojson.Marshal(req)
	if err != nil {
		return nil, err
	}

	resp, err := c.restClient.Post(ctx, fmt.Sprintf("/files/%s/public-url", fileID), jsonData)
	if err != nil {
		return nil, err
	}

	var publicURLResp xaiv1.CreatePublicUrlResponse
	if err := protojson.Unmarshal(resp.Body, &publicURLResp); err != nil {
		return nil, err
	}
	return &publicURLResp, nil
}

// RevokePublicURL revokes the public URL for a stored file.
func (c *Client) RevokePublicURL(ctx context.Context, fileID string) (*xaiv1.RevokePublicUrlResponse, error) {
	if c.grpcClient != nil {
		return c.grpcClient.RevokePublicUrl(ctx, &xaiv1.RevokePublicUrlRequest{FileId: fileID})
	}
	if c.restClient == nil {
		return nil, ErrClientNotInitialized
	}

	resp, err := c.restClient.Delete(ctx, fmt.Sprintf("/files/%s/public-url", fileID))
	if err != nil {
		return nil, err
	}

	var revokeResp xaiv1.RevokePublicUrlResponse
	if err := protojson.Unmarshal(resp.Body, &revokeResp); err != nil {
		return nil, err
	}
	return &revokeResp, nil
}

// Get retrieves file metadata by ID.
func (c *Client) Get(ctx context.Context, fileID string) (*File, error) {
	if c.grpcClient != nil {
		f, err := c.grpcClient.RetrieveFile(ctx, &xaiv1.RetrieveFileRequest{FileId: fileID})
		return fromProto(f), err
	}
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
	if c.grpcClient != nil {
		_, err := c.grpcClient.DeleteFile(ctx, &xaiv1.DeleteFileRequest{FileId: fileID})
		return err
	}
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

// uploadGRPC streams bounded chunks without buffering the entire file.
func (c *Client) uploadGRPC(ctx context.Context, reader io.Reader, opts UploadOptions) (*File, error) {
	if reader == nil {
		return nil, fmt.Errorf("file reader is required")
	}
	name := opts.Name
	if name == "" {
		if named, ok := reader.(interface{ Name() string }); ok {
			name = filepath.Base(named.Name())
		}
	}
	if name == "" {
		return nil, fmt.Errorf("file name is required")
	}
	if opts.ExpiresAfter != nil && *opts.ExpiresAfter <= 0 {
		return nil, fmt.Errorf("expiry must be positive")
	}
	streamCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	stream, err := c.grpcClient.UploadFile(streamCtx)
	if err != nil {
		return nil, err
	}
	init := &xaiv1.UploadFileInit{Name: name}
	if opts.ExpiresAfter != nil {
		seconds := int64(opts.ExpiresAfter.Seconds())
		init.ExpiresAfter = &seconds
	}
	if sendErr := stream.Send(&xaiv1.UploadFileChunk{Chunk: &xaiv1.UploadFileChunk_Init{Init: init}}); sendErr != nil {
		return nil, sendErr
	}
	total, sizeErr := remainingUploadSize(reader)
	if sizeErr != nil {
		return nil, sizeErr
	}
	if sendErr := sendUploadChunks(ctx, reader, opts, stream, total); sendErr != nil {
		return nil, sendErr
	}

	file, err := stream.CloseAndRecv()
	return fromProto(file), err
}

type contentReader struct {
	recv    func() (*xaiv1.FileContentChunk, error)
	cancel  context.CancelFunc
	pending []byte
}

func (r *contentReader) Read(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}
	for len(r.pending) == 0 {
		chunk, err := r.recv()
		if err != nil {
			r.cancel()
			return 0, err
		}
		r.pending = chunk.Data
	}
	n := copy(p, r.pending)
	r.pending = r.pending[n:]
	return n, nil
}
func (r *contentReader) Close() error { r.cancel(); return nil }

func remainingUploadSize(reader io.Reader) (int64, error) {
	seeker, ok := reader.(io.Seeker)
	if !ok {
		return -1, nil
	}
	current, currentErr := seeker.Seek(0, io.SeekCurrent)
	if currentErr != nil {
		return -1, nil
	}
	end, endErr := seeker.Seek(0, io.SeekEnd)
	if _, restoreErr := seeker.Seek(current, io.SeekStart); restoreErr != nil {
		return -1, restoreErr
	}
	if endErr != nil {
		return -1, nil
	}
	return end - current, nil
}
func sendUploadChunks(ctx context.Context, reader io.Reader, opts UploadOptions, stream interface {
	Send(*xaiv1.UploadFileChunk) error
}, total int64) error {
	buffer := make([]byte, 3*1024*1024)
	var uploaded int64
	emptyReads := 0
	for {
		if cancelErr := ctx.Err(); cancelErr != nil {
			return cancelErr
		}
		n, readErr := reader.Read(buffer)
		if n > 0 {
			emptyReads = 0
			uploaded += int64(n)
			if opts.MaxSize > 0 && uploaded > opts.MaxSize {
				return ErrFileTooLarge
			}
			if sendErr := stream.Send(&xaiv1.UploadFileChunk{Chunk: &xaiv1.UploadFileChunk_Data{Data: buffer[:n]}}); sendErr != nil {
				return sendErr
			}
			if opts.OnProgress != nil {
				opts.OnProgress(uploaded, total)
			}
		} else {
			emptyReads++
			if emptyReads >= 100 && readErr == nil {
				return io.ErrNoProgress
			}
		}
		if readErr == io.EOF {
			return nil
		}
		if readErr != nil {
			return readErr
		}
	}
}
