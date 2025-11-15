// Package image provides image generation functionality for xAI SDK.
package image

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	xaiv1 "github.com/ZaguanLabs/xai-sdk-go/proto/gen/go/xai/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Image represents a generated image.
type Image struct {
	url           string
	data          []byte
	revisedPrompt string
}

// ImageServiceClient is an interface for the image service client.
type ImageServiceClient interface {
	GenerateImage(ctx context.Context, req *xaiv1.GenerateImageRequest, opts ...grpc.CallOption) (*xaiv1.GenerateImageResponse, error)
}

// GenerateRequest represents an image generation request.
type GenerateRequest struct {
	prompt      string
	model       string
	size        string
	quality     string
	style       string
	n           int32
}

// Client provides image generation functionality.
type Client struct {
	grpcClient ImageServiceClient
	httpClient *http.Client
}

// NewClient creates a new image generation client.
func NewClient(grpcClient ImageServiceClient) *Client {
	return &Client{
		grpcClient: grpcClient,
		httpClient: &http.Client{},
	}
}

// Generate generates an image from a text prompt.
func (c *Client) Generate(ctx context.Context, req *GenerateRequest) ([]*Image, error) {
	if req == nil {
		return nil, fmt.Errorf("generate request is nil")
	}

	// Validate request
	if err := req.validate(); err != nil {
		return nil, fmt.Errorf("invalid generate request: %w", err)
	}

	// Create proto request
	protoReq := &xaiv1.GenerateImageRequest{
		Prompt:  req.prompt,
		Model:   req.model,
		Size:    req.size,
		Quality: req.quality,
		Style:   req.style,
		N:       req.n,
	}

	// Call image generation service
	resp, err := c.grpcClient.GenerateImage(ctx, protoReq)
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
				return nil, fmt.Errorf("image generation failed (%s): %s", st.Code().String(), st.Message())
			}
		}
		return nil, fmt.Errorf("image generation failed: %w", err)
	}

	if resp == nil {
		return nil, fmt.Errorf("received nil response")
	}

	// Convert response images
	images := make([]*Image, 0, len(resp.Images))
	for _, imgProto := range resp.Images {
		images = append(images, &Image{
			url:           imgProto.Url,
			data:          imgProto.Data,
			revisedPrompt: imgProto.RevisedPrompt,
		})
	}

	return images, nil
}

// GenerateRequest methods

// NewGenerateRequest creates a new image generation request.
func NewGenerateRequest(prompt, model string) *GenerateRequest {
	return &GenerateRequest{
		prompt:  prompt,
		model:   model,
		size:    "1024x1024",
		quality: "standard",
		style:   "vivid",
		n:       1,
	}
}

// WithSize sets the image size.
func (r *GenerateRequest) WithSize(size string) *GenerateRequest {
	r.size = size
	return r
}

// WithQuality sets the image quality.
func (r *GenerateRequest) WithQuality(quality string) *GenerateRequest {
	r.quality = quality
	return r
}

// WithStyle sets the image style.
func (r *GenerateRequest) WithStyle(style string) *GenerateRequest {
	r.style = style
	return r
}

// WithN sets the number of images to generate.
func (r *GenerateRequest) WithN(n int32) *GenerateRequest {
	r.n = n
	return r
}

// validate validates the generate request.
func (r *GenerateRequest) validate() error {
	if r.prompt == "" {
		return fmt.Errorf("prompt is required")
	}
	if r.model == "" {
		return fmt.Errorf("model is required")
	}
	if r.n < 1 || r.n > 10 {
		return fmt.Errorf("n must be between 1 and 10, got %d", r.n)
	}
	
	// Validate size
	validSizes := map[string]bool{
		"256x256":   true,
		"512x512":   true,
		"1024x1024": true,
		"1792x1024": true,
		"1024x1792": true,
	}
	if !validSizes[r.size] {
		return fmt.Errorf("invalid size '%s', must be one of: 256x256, 512x512, 1024x1024, 1792x1024, 1024x1792", r.size)
	}
	
	// Validate quality
	validQualities := map[string]bool{
		"standard": true,
		"hd":       true,
	}
	if !validQualities[r.quality] {
		return fmt.Errorf("invalid quality '%s', must be 'standard' or 'hd'", r.quality)
	}
	
	// Validate style
	validStyles := map[string]bool{
		"vivid":    true,
		"natural":  true,
	}
	if !validStyles[r.style] {
		return fmt.Errorf("invalid style '%s', must be 'vivid' or 'natural'", r.style)
	}
	
	return nil
}

// Image methods

// URL returns the image URL.
func (i *Image) URL() string {
	return i.url
}

// Data returns the image data.
func (i *Image) Data() []byte {
	return i.data
}

// RevisedPrompt returns the revised prompt used for generation.
func (i *Image) RevisedPrompt() string {
	return i.revisedPrompt
}

// Save saves the image to a local file.
func (i *Image) Save(filePath string, client *http.Client) error {
	if i.data == nil && i.url == "" {
		return fmt.Errorf("no image data or URL available")
	}

	// If we have data, save it directly
	if len(i.data) > 0 {
		return os.WriteFile(filePath, i.data, 0644)
	}

	// If we have a URL, download the image
	if i.url != "" {
		return i.downloadAndSave(filePath, client)
	}

	return fmt.Errorf("no image content available")
}

// downloadAndSave downloads the image from URL and saves it.
func (i *Image) downloadAndSave(filePath string, client *http.Client) error {
	// Parse URL
	u, err := url.Parse(i.url)
	if err != nil {
		return fmt.Errorf("invalid image URL: %w", err)
	}

	// Download image
	resp, err := client.Get(u.String())
	if err != nil {
		return fmt.Errorf("failed to download image: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download image: status %d", resp.StatusCode)
	}

	// Create directory if needed
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Save to file
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to save image: %w", err)
	}

	return nil
}