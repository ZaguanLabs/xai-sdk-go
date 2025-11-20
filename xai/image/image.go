// Package image provides a client for the xAI Image Generation API.
package image

import (
	"context"

	xaiv1 "github.com/ZaguanLabs/xai-sdk-go/proto/gen/go/xai/v1"
	"github.com/ZaguanLabs/xai-sdk-go/xai/internal/rest"
	"google.golang.org/protobuf/encoding/protojson"
)

// Client provides access to the xAI Image Generation API.
type Client struct {
	restClient *rest.Client
}

// NewClient creates a new Image Generation API client.
func NewClient(restClient *rest.Client) *Client {
	return &Client{
		restClient: restClient,
	}
}

// GenerateRequest represents an image generation request.
type GenerateRequest struct {
	Prompt string
	Model  string
	N      int32
	User   string
	Image  *Input
	Format xaiv1.ImageFormat
}

// Input represents an input image for image-to-image generation.
type Input struct {
	ImageURL string
	Detail   xaiv1.ImageDetail
}

// GeneratedImage represents a generated image.
type GeneratedImage struct {
	Base64            string
	URL               string
	UpsampledPrompt   string
	RespectModeration bool
}

// Response represents an image generation response.
type Response struct {
	Images []*GeneratedImage
	Model  string
}

// NewRequest creates a new image generation request.
func NewRequest(prompt, model string) *GenerateRequest {
	return &GenerateRequest{
		Prompt: prompt,
		Model:  model,
		N:      1,
		Format: xaiv1.ImageFormat_IMG_FORMAT_URL,
	}
}

// WithCount sets the number of images to generate.
func (r *GenerateRequest) WithCount(n int32) *GenerateRequest {
	r.N = n
	return r
}

// WithUser sets the user identifier.
func (r *GenerateRequest) WithUser(user string) *GenerateRequest {
	r.User = user
	return r
}

// WithFormat sets the image format (URL or Base64).
func (r *GenerateRequest) WithFormat(format xaiv1.ImageFormat) *GenerateRequest {
	r.Format = format
	return r
}

// WithImage sets an input image for image-to-image generation.
func (r *GenerateRequest) WithImage(imageURL string, detail xaiv1.ImageDetail) *GenerateRequest {
	r.Image = &Input{
		ImageURL: imageURL,
		Detail:   detail,
	}
	return r
}

// Generate generates images based on the request.
func (c *Client) Generate(ctx context.Context, req *GenerateRequest) (*Response, error) {
	if c.restClient == nil {
		return nil, ErrClientNotInitialized
	}

	protoReq := &xaiv1.GenerateImageRequest{
		Prompt: req.Prompt,
		Model:  req.Model,
		N:      req.N,
		User:   req.User,
		Format: req.Format,
	}

	if req.Image != nil {
		protoReq.Image = &xaiv1.ImageUrlContent{
			ImageUrl: req.Image.ImageURL,
			Detail:   req.Image.Detail,
		}
	}

	jsonData, err := protojson.Marshal(protoReq)
	if err != nil {
		return nil, err
	}

	resp, err := c.restClient.Post(ctx, "/images/generations", jsonData)
	if err != nil {
		return nil, err
	}

	var imageResp xaiv1.ImageResponse
	if err := protojson.Unmarshal(resp.Body, &imageResp); err != nil {
		return nil, err
	}

	images := make([]*GeneratedImage, len(imageResp.Images))
	for i, img := range imageResp.Images {
		images[i] = &GeneratedImage{
			Base64:            img.Base64,
			URL:               img.Url,
			UpsampledPrompt:   img.UpSampledPrompt,
			RespectModeration: img.RespectModeration,
		}
	}

	return &Response{
		Images: images,
		Model:  imageResp.Model,
	}, nil
}
