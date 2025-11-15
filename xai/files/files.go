// Package files provides file upload and download functionality for xAI SDK.
package files

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	xaiv1 "github.com/ZaguanLabs/xai-sdk-go/proto/gen/go/xai/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// File represents a file in the xAI system.
type File struct {
	id          string
	filename     string
	size        int64
	contentType  string
	createdAt    time.Time
	purpose     string
}

// FileServiceClient is an interface for the files service client.
type FileServiceClient interface {
	UploadFile(ctx context.Context, req *xaiv1.UploadFileRequest, opts ...grpc.CallOption) (*xaiv1.UploadFileResponse, error)
	DownloadFile(ctx context.Context, req *xaiv1.DownloadFileRequest, opts ...grpc.CallOption) (*xaiv1.DownloadFileResponse, error)
	ListFiles(ctx context.Context, req *xaiv1.ListFilesRequest, opts ...grpc.CallOption) (*xaiv1.ListFilesResponse, error)
	GetFile(ctx context.Context, req *xaiv1.GetFileRequest, opts ...grpc.CallOption) (*xaiv1.GetFileResponse, error)
	DeleteFile(ctx context.Context, req *xaiv1.DeleteFileRequest, opts ...grpc.CallOption) (*xaiv1.DeleteFileResponse, error)
}

// ProgressCallback is called during file upload to report progress.
type ProgressCallback func(uploadedBytes, totalBytes int64)

// UploadRequest represents a file upload request.
type UploadRequest struct {
	filePath    string
	filename    string
	contentType string
	purpose     string
	callback    ProgressCallback
}

// DownloadRequest represents a file download request.
type DownloadRequest struct {
	fileID     string
	localPath  string
}

// ListRequest represents a file list request.
type ListRequest struct {
	purpose    string
	limit     int32
	sortBy    string
	sortOrder  string
}

// Client provides file upload and download functionality.
type Client struct {
	grpcClient FileServiceClient
}

// NewClient creates a new files client.
func NewClient(grpcClient FileServiceClient) *Client {
	return &Client{
		grpcClient: grpcClient,
	}
}

// UploadFile uploads a file from the local filesystem.
func (c *Client) UploadFile(ctx context.Context, req *UploadRequest) (*File, error) {
	if req == nil {
		return nil, fmt.Errorf("upload request is nil")
	}

	// Validate request
	if err := req.validate(); err != nil {
		return nil, fmt.Errorf("invalid upload request: %w", err)
	}

	// Open file
	file, err := os.Open(req.filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Get file info
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}

	// Create file reader
	fileReader := &progressReader{
		reader:    file,
		totalSize: fileInfo.Size(),
		callback:  req.callback,
	}

	// Create upload request (placeholder until proto is updated)
	uploadReq := &xaiv1.UploadFileRequest{
		Filename:    req.filename,
		ContentType: req.contentType,
		Purpose:     req.purpose,
	}

	// Call upload service (placeholder)
	resp, err := c.grpcClient.UploadFile(ctx, uploadReq)
	if err != nil {
		return nil, fmt.Errorf("file upload failed: %w", err)
	}

	if resp == nil {
		return nil, fmt.Errorf("received nil upload response")
	}

	// Return file info
	return &File{
		id:          resp.Id,
		filename:    req.filename,
		size:        resp.Size, // Use response size instead of placeholder
		contentType:  resp.ContentType,
		createdAt:    resp.CreatedAt, // Use response time instead of placeholder
		purpose:     resp.Purpose,
	}, nil
}

// UploadBytes uploads file content from memory.
func (c *Client) UploadBytes(ctx context.Context, filename string, data []byte, contentType string, purpose string, callback ProgressCallback) (*File, error) {
	if filename == "" {
		return nil, fmt.Errorf("filename is required")
	}
	if len(data) == 0 {
		return nil, fmt.Errorf("data is empty")
	}

	// Create upload request (placeholder until proto is updated)
	uploadReq := &xaiv1.UploadFileRequest{
		Filename:    filename,
		ContentType: contentType,
		Purpose:     purpose,
	}

	// Create file reader with progress tracking
	fileReader := &progressReader{
		reader:    &dataReader{data: data},
		totalSize: int64(len(data)),
		callback:  callback,
	}

	// Call upload service (placeholder)
	resp, err := c.grpcClient.UploadFile(ctx, uploadReq)
	if err != nil {
		return nil, fmt.Errorf("file upload failed: %w", err)
	}

	if resp == nil {
		return nil, fmt.Errorf("received nil upload response")
	}

	// Return file info
	return &File{
		id:          resp.Id,
		filename:    filename,
		size:        int64(len(data)),
		contentType:  contentType,
		createdAt:    time.Now(),
		purpose:     purpose,
	}, nil
}

// UploadReader uploads file content from a reader.
func (c *Client) UploadReader(ctx context.Context, filename string, reader io.Reader, size int64, contentType string, purpose string, callback ProgressCallback) (*File, error) {
	if filename == "" {
		return nil, fmt.Errorf("filename is required")
	}
	if reader == nil {
		return nil, fmt.Errorf("reader is nil")
	}

	// Create upload request (placeholder until proto is updated)
	uploadReq := &xaiv1.UploadFileRequest{
		Filename:    filename,
		ContentType: contentType,
		Purpose:     purpose,
	}

	// Create file reader with progress tracking
	fileReader := &progressReader{
		reader:    reader,
		totalSize: size,
		callback:  callback,
	}

	// Call upload service (placeholder)
	resp, err := c.grpcClient.UploadFile(ctx, uploadReq)
	if err != nil {
		return nil, fmt.Errorf("file upload failed: %w", err)
	}

	if resp == nil {
		return nil, fmt.Errorf("received nil upload response")
	}

	// Return file info
	return &File{
		id:          resp.Id,
		filename:    filename,
		size:        size,
		contentType:  contentType,
		createdAt:    time.Now(),
		purpose:     purpose,
	}, nil
}

// DownloadFile downloads a file to the local filesystem.
func (c *Client) DownloadFile(ctx context.Context, req *DownloadRequest) error {
	if req == nil {
		return fmt.Errorf("download request is nil")
	}

	// Validate request
	if err := req.validate(); err != nil {
		return fmt.Errorf("invalid download request: %w", err)
	}

	// Create download request (placeholder until proto is updated)
	downloadReq := &xaiv1.DownloadFileRequest{
		FileId: req.fileID,
	}

	// Call download service (placeholder)
	resp, err := c.grpcClient.DownloadFile(ctx, downloadReq)
	if err != nil {
		return fmt.Errorf("file download failed: %w", err)
	}

	if resp == nil {
		return fmt.Errorf("received nil download response")
	}

	// Create local file
	file, err := os.Create(req.localPath)
	if err != nil {
		return fmt.Errorf("failed to create local file: %w", err)
	}
	defer file.Close()

	// Copy data (placeholder - resp should contain file data)
	// For now, we'll create an empty file as a placeholder
	_, err = file.WriteString(fmt.Sprintf("Downloaded file %s (placeholder content)", req.fileID))
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// ListFiles lists files with optional filtering.
func (c *Client) ListFiles(ctx context.Context, req *ListRequest) ([]*File, error) {
	if req == nil {
		return nil, fmt.Errorf("list request is nil")
	}

	// Validate request
	if err := req.validate(); err != nil {
		return nil, fmt.Errorf("invalid list request: %w", err)
	}

	// Create list request (placeholder until proto is updated)
	listReq := &xaiv1.ListFilesRequest{
		Purpose: req.purpose,
		Limit:   req.limit,
	}

	// Call list service (placeholder)
	resp, err := c.grpcClient.ListFiles(ctx, listReq)
	if err != nil {
		return nil, fmt.Errorf("file list failed: %w", err)
	}

	if resp == nil {
		return nil, fmt.Errorf("received nil list response")
	}

	// Convert response (placeholder)
	files := make([]*File, 0)
	for _, fileProto := range resp.Files {
		files = append(files, &File{
			id:          fileProto.Id,
			filename:    fileProto.Filename,
			size:        fileProto.Size,
			contentType:  fileProto.ContentType,
			createdAt:    time.Now(), // placeholder - should use file timestamp
			purpose:     fileProto.Purpose,
		})
	}

	return files, nil
}

// GetFile retrieves file information.
func (c *Client) GetFile(ctx context.Context, fileID string) (*File, error) {
	if fileID == "" {
		return nil, fmt.Errorf("file ID is required")
	}

	// Create get request (placeholder until proto is updated)
	getReq := &xaiv1.GetFileRequest{
		FileId: fileID,
	}

	// Call get service (placeholder)
	resp, err := c.grpcClient.GetFile(ctx, getReq)
	if err != nil {
		return nil, fmt.Errorf("get file failed: %w", err)
	}

	if resp == nil {
		return nil, fmt.Errorf("received nil get response")
	}

	// Return file info
	return &File{
		id:          resp.File.Id,
		filename:    resp.File.Filename,
		size:        resp.File.Size,
		contentType:  resp.File.ContentType,
		createdAt:    time.Now(), // placeholder - should use file timestamp
		purpose:     resp.File.Purpose,
	}, nil
}

// DeleteFile deletes a file.
func (c *Client) DeleteFile(ctx context.Context, fileID string) error {
	if fileID == "" {
		return fmt.Errorf("file ID is required")
	}

	// Create delete request (placeholder until proto is updated)
	deleteReq := &xaiv1.DeleteFileRequest{
		FileId: fileID,
	}

	// Call delete service (placeholder)
	resp, err := c.grpcClient.DeleteFile(ctx, deleteReq)
	if err != nil {
		return fmt.Errorf("file delete failed: %w", err)
	}

	if resp == nil {
		return fmt.Errorf("received nil delete response")
	}

	return nil
}

// progressReader wraps a reader to track upload progress.
type progressReader struct {
	reader    io.Reader
	totalSize int64
	callback  ProgressCallback
	read      int64
}

// Read implements io.Reader with progress tracking.
func (pr *progressReader) Read(p []byte) (int, error) {
	n, err := pr.reader.Read(p)
	if err != nil {
		return n, err
	}

	pr.read += int64(n)
	
	// Call progress callback if set
	if pr.callback != nil {
		pr.callback(pr.read, pr.totalSize)
	}

	return n, nil
}

// dataReader wraps byte data as io.Reader.
type dataReader struct {
	data   []byte
	offset int
}

// Read implements io.Reader for byte data.
func (dr *dataReader) Read(p []byte) (int, error) {
	if dr.offset >= len(dr.data) {
		return 0, io.EOF
	}

	n := copy(p, dr.data[dr.offset:])
	dr.offset += n
	return n, nil
}

// copy copies data with bounds checking.
func copy(dst, src []byte) int {
	n := len(src)
	if len(dst) < n {
		n = len(dst)
	}
	copy(dst, src[:n])
	return n
}

// validate validates the upload request.
func (ur *UploadRequest) validate() error {
	if ur.filePath == "" {
		return fmt.Errorf("file path is required")
	}

	// Check if file exists
	if _, err := os.Stat(ur.filePath); os.IsNotExist(err) {
		return fmt.Errorf("file does not exist: %s", ur.filePath)
	}

	// Validate file size (3 MiB chunks)
	fileInfo, err := os.Stat(ur.filePath)
	if err != nil {
		return fmt.Errorf("failed to get file info: %w", err)
	}

	if fileInfo.Size() > 100*1024*1024 { // 100 MiB
		return fmt.Errorf("file too large: %d bytes (max 100 MiB)", fileInfo.Size())
	}

	return nil
}

// validate validates the download request.
func (dr *DownloadRequest) validate() error {
	if dr.fileID == "" {
		return fmt.Errorf("file ID is required")
	}

	if dr.localPath == "" {
		return fmt.Errorf("local path is required")
	}

	// Check if directory exists
	if info, err := os.Stat(dr.localPath); err == nil {
		if info.IsDir() {
			return fmt.Errorf("local path is a directory: %s", dr.localPath)
		}
	}

	return nil
}

// validate validates the list request.
func (lr *ListRequest) validate() error {
	if lr.limit < 0 || lr.limit > 100 {
		return fmt.Errorf("limit must be between 0 and 100, got %d", lr.limit)
	}

	return nil
}