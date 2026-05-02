// Package image provides a client for the xAI Image Generation API.
package image

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"

	xaiv1 "github.com/ZaguanLabs/xai-sdk-go/proto/gen/go/xai/api/v1"
	"github.com/ZaguanLabs/xai-sdk-go/xai/cost"
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
	Prompt      string
	Model       string
	N           int32
	User        string
	Image       *Input
	Images      []*Input
	Format      xaiv1.ImageFormat
	AspectRatio *xaiv1.ImageAspectRatio
	Resolution  *xaiv1.ImageResolution
}

// Input represents an input image for image-to-image generation.
type Input struct {
	ImageURL string
	Detail   xaiv1.ImageDetail
}

// GeneratedImage represents a generated image.
type GeneratedImage struct {
	proto *xaiv1.GeneratedImage
}

// Base64 returns the base64 encoded image data.
func (i *GeneratedImage) Base64() string {
	if i.proto == nil {
		return ""
	}
	return i.proto.GetBase64()
}

func (i *GeneratedImage) DecodeBase64() ([]byte, error) {
	value := i.Base64()
	if value == "" {
		if !i.RespectModeration() {
			return nil, fmt.Errorf("image did not respect moderation rules; base64 is not available")
		}
		return nil, fmt.Errorf("image was not returned via base64")
	}
	if comma := strings.Index(value, "base64,"); comma >= 0 {
		value = value[comma+len("base64,"):]
	}
	return base64.StdEncoding.DecodeString(value)
}

// URL returns the URL of the generated image.
func (i *GeneratedImage) URL() string {
	if i.proto == nil {
		return ""
	}
	return i.proto.GetUrl()
}

// UpsampledPrompt returns the upsampled prompt.
func (i *GeneratedImage) UpsampledPrompt() string {
	return ""
}

// RespectModeration returns whether the image respects moderation.
func (i *GeneratedImage) RespectModeration() bool {
	if i.proto == nil {
		return false
	}
	return i.proto.RespectModeration
}

// Response represents an image generation response.
type Response struct {
	Images []*GeneratedImage
	Model  string
	Usage  *xaiv1.SamplingUsage
}

func (r *Response) Image() *GeneratedImage {
	if r == nil || len(r.Images) == 0 {
		return nil
	}
	return r.Images[0]
}

func (r *Response) CostUSD() (float64, bool) {
	if r == nil {
		return 0, false
	}
	return cost.USDFromUsage(r.Usage)
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

func (r *GenerateRequest) WithImages(images ...*Input) *GenerateRequest {
	r.Images = append(r.Images, images...)
	return r
}

func (r *GenerateRequest) WithImageURL(imageURL string, detail xaiv1.ImageDetail) *GenerateRequest {
	return r.WithImage(imageURL, detail)
}

func (r *GenerateRequest) WithAspectRatio(aspectRatio xaiv1.ImageAspectRatio) *GenerateRequest {
	r.AspectRatio = &aspectRatio
	return r
}

func (r *GenerateRequest) WithResolution(resolution xaiv1.ImageResolution) *GenerateRequest {
	r.Resolution = &resolution
	return r
}

func (r *GenerateRequest) Proto() *xaiv1.GenerateImageRequest {
	protoReq := &xaiv1.GenerateImageRequest{
		Prompt: r.Prompt,
		Model:  r.Model,
		N:      &r.N,
		User:   r.User,
		Format: r.Format,
	}

	if r.Image != nil {
		protoReq.Image = &xaiv1.ImageUrlContent{
			ImageUrl: r.Image.ImageURL,
			Detail:   r.Image.Detail,
		}
	}
	if len(r.Images) > 0 {
		protoReq.Images = make([]*xaiv1.ImageUrlContent, 0, len(r.Images))
		for _, img := range r.Images {
			if img == nil {
				continue
			}
			protoReq.Images = append(protoReq.Images, &xaiv1.ImageUrlContent{
				ImageUrl: img.ImageURL,
				Detail:   img.Detail,
			})
		}
	}
	if r.AspectRatio != nil {
		protoReq.AspectRatio = r.AspectRatio
	}
	if r.Resolution != nil {
		protoReq.Resolution = r.Resolution
	}

	return protoReq
}

func (c *Client) Prepare(req *GenerateRequest, batchRequestID string) *xaiv1.BatchRequest {
	if req == nil {
		return nil
	}
	batchReq := &xaiv1.BatchRequest{
		Request: &xaiv1.BatchRequest_ImageRequest{
			ImageRequest: req.Proto(),
		},
	}
	if batchRequestID != "" {
		batchReq.BatchRequestId = &batchRequestID
	}
	return batchReq
}

func (c *Client) Sample(ctx context.Context, req *GenerateRequest) (*GeneratedImage, error) {
	resp, err := c.Generate(ctx, req)
	if err != nil {
		return nil, err
	}
	if len(resp.Images) == 0 {
		return nil, nil
	}
	return resp.Images[0], nil
}

func (c *Client) SampleBatch(ctx context.Context, req *GenerateRequest, n int32) ([]*GeneratedImage, error) {
	req.WithCount(n)
	resp, err := c.Generate(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.Images, nil
}

// Generate generates images based on the request.
func (c *Client) Generate(ctx context.Context, req *GenerateRequest) (*Response, error) {
	if c.restClient == nil {
		return nil, ErrClientNotInitialized
	}

	jsonData, err := protojson.Marshal(req.Proto())
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
			proto: img,
		}
	}

	return &Response{
		Images: images,
		Model:  imageResp.Model,
		Usage:  imageResp.Usage,
	}, nil
}
