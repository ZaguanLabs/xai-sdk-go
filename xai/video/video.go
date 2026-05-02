package video

import (
	"context"
	"fmt"
	"time"

	xaiv1 "github.com/ZaguanLabs/xai-sdk-go/proto/gen/go/xai/api/v1"
	"github.com/ZaguanLabs/xai-sdk-go/xai/cost"
)

const (
	DefaultPollTimeout  = 10 * time.Minute
	DefaultPollInterval = 2 * time.Second
)

type Client struct {
	grpcClient xaiv1.VideoClient
}

type GenerateOptions struct {
	ImageURL           string
	VideoURL           string
	ReferenceImageURLs []string
	Duration           *int32
	AspectRatio        *xaiv1.VideoAspectRatio
	Resolution         *xaiv1.VideoResolution
	Timeout            time.Duration
	Interval           time.Duration
}

type Response struct {
	proto *xaiv1.VideoResponse
}

func NewResponse(proto *xaiv1.VideoResponse) *Response {
	return &Response{proto: proto}
}

func (r *Response) Proto() *xaiv1.VideoResponse {
	if r == nil {
		return nil
	}
	return r.proto
}

func (r *Response) Model() string {
	if r == nil || r.proto == nil {
		return ""
	}
	return r.proto.Model
}

func (r *Response) Usage() *xaiv1.SamplingUsage {
	if r == nil || r.proto == nil {
		return nil
	}
	return r.proto.Usage
}

func (r *Response) CostUSD() (float64, bool) {
	return cost.USDFromUsage(r.Usage())
}

func (r *Response) URL() (string, error) {
	if r == nil || r.proto == nil || r.proto.Video == nil {
		return "", fmt.Errorf("video URL missing from response")
	}
	if r.proto.Video.Url == "" {
		if !r.RespectModeration() {
			return "", fmt.Errorf("video did not respect moderation rules; URL is not available")
		}
		return "", fmt.Errorf("video URL missing from response")
	}
	return r.proto.Video.Url, nil
}

func (r *Response) Duration() int32 {
	if r == nil || r.proto == nil || r.proto.Video == nil {
		return 0
	}
	return r.proto.Video.Duration
}

func (r *Response) RespectModeration() bool {
	if r == nil || r.proto == nil || r.proto.Video == nil {
		return true
	}
	return r.proto.Video.RespectModeration
}

func NewClient(grpcClient xaiv1.VideoClient) *Client {
	return &Client{grpcClient: grpcClient}
}

func NewGenerateRequest(prompt, model string) *xaiv1.GenerateVideoRequest {
	return &xaiv1.GenerateVideoRequest{
		Prompt: prompt,
		Model:  model,
	}
}

func NewGenerateRequestWithOptions(prompt, model string, opts *GenerateOptions) *xaiv1.GenerateVideoRequest {
	req := NewGenerateRequest(prompt, model)
	if opts == nil {
		return req
	}
	if opts.ImageURL != "" {
		req.Image = &xaiv1.ImageUrlContent{ImageUrl: opts.ImageURL}
	}
	if opts.VideoURL != "" {
		req.Video = &xaiv1.VideoUrlContent{Url: opts.VideoURL}
	}
	for _, imageURL := range opts.ReferenceImageURLs {
		req.ReferenceImages = append(req.ReferenceImages, &xaiv1.ImageUrlContent{ImageUrl: imageURL})
	}
	if opts.Duration != nil {
		req.Duration = opts.Duration
	}
	if opts.AspectRatio != nil {
		req.AspectRatio = opts.AspectRatio
	}
	if opts.Resolution != nil {
		req.Resolution = opts.Resolution
	}
	return req
}

func NewExtendRequest(prompt, model, videoURL string, duration *int32) *xaiv1.ExtendVideoRequest {
	return &xaiv1.ExtendVideoRequest{
		Prompt:   prompt,
		Model:    model,
		Video:    &xaiv1.VideoUrlContent{Url: videoURL},
		Duration: duration,
	}
}

func Prepare(prompt, model, batchRequestID string, opts *GenerateOptions) *xaiv1.BatchRequest {
	req := NewGenerateRequestWithOptions(prompt, model, opts)
	batchReq := &xaiv1.BatchRequest{
		Request: &xaiv1.BatchRequest_VideoRequest{
			VideoRequest: req,
		},
	}
	if batchRequestID != "" {
		batchReq.BatchRequestId = &batchRequestID
	}
	return batchReq
}

func (c *Client) Prepare(prompt, model, batchRequestID string, opts *GenerateOptions) *xaiv1.BatchRequest {
	return Prepare(prompt, model, batchRequestID, opts)
}

func (c *Client) GenerateDeferred(ctx context.Context, req *xaiv1.GenerateVideoRequest) (*xaiv1.StartDeferredResponse, error) {
	if c.grpcClient == nil {
		return nil, fmt.Errorf("video client not initialized")
	}
	return c.grpcClient.GenerateVideo(ctx, req)
}

func (c *Client) Start(ctx context.Context, prompt, model string, opts *GenerateOptions) (*xaiv1.StartDeferredResponse, error) {
	return c.GenerateDeferred(ctx, NewGenerateRequestWithOptions(prompt, model, opts))
}

func (c *Client) ExtendStart(ctx context.Context, prompt, model, videoURL string, duration *int32) (*xaiv1.StartDeferredResponse, error) {
	if c.grpcClient == nil {
		return nil, fmt.Errorf("video client not initialized")
	}
	return c.grpcClient.ExtendVideo(ctx, NewExtendRequest(prompt, model, videoURL, duration))
}

func (c *Client) GetDeferred(ctx context.Context, requestID string) (*xaiv1.GetDeferredVideoResponse, error) {
	if c.grpcClient == nil {
		return nil, fmt.Errorf("video client not initialized")
	}
	return c.grpcClient.GetDeferredVideo(ctx, &xaiv1.GetDeferredVideoRequest{
		RequestId: requestID,
	})
}

func (c *Client) Get(ctx context.Context, requestID string) (*xaiv1.GetDeferredVideoResponse, error) {
	return c.GetDeferred(ctx, requestID)
}

func (c *Client) Generate(ctx context.Context, prompt, model string, opts *GenerateOptions) (*xaiv1.VideoResponse, error) {
	start, err := c.Start(ctx, prompt, model, opts)
	if err != nil {
		return nil, err
	}
	return c.poll(ctx, start.RequestId, opts)
}

func (c *Client) GenerateAndPoll(ctx context.Context, prompt, model string, opts *GenerateOptions) (*xaiv1.VideoResponse, error) {
	return c.Generate(ctx, prompt, model, opts)
}

func (c *Client) GenerateSync(ctx context.Context, prompt, model string, opts *GenerateOptions) (*xaiv1.VideoResponse, error) {
	return c.Generate(ctx, prompt, model, opts)
}

func (c *Client) Extend(ctx context.Context, prompt, model, videoURL string, duration *int32, opts *GenerateOptions) (*xaiv1.VideoResponse, error) {
	start, err := c.ExtendStart(ctx, prompt, model, videoURL, duration)
	if err != nil {
		return nil, err
	}
	return c.poll(ctx, start.RequestId, opts)
}

func (c *Client) poll(ctx context.Context, requestID string, opts *GenerateOptions) (*xaiv1.VideoResponse, error) {
	timeout := DefaultPollTimeout
	interval := DefaultPollInterval
	if opts != nil {
		if opts.Timeout > 0 {
			timeout = opts.Timeout
		}
		if opts.Interval > 0 {
			interval = opts.Interval
		}
	}

	pollCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		resp, err := c.Get(pollCtx, requestID)
		if err != nil {
			return nil, err
		}

		switch resp.Status {
		case xaiv1.DeferredStatus_DONE:
			if resp.Response == nil {
				return nil, fmt.Errorf("deferred video completed without a response")
			}
			return resp.Response, nil
		case xaiv1.DeferredStatus_EXPIRED:
			return nil, fmt.Errorf("deferred video request expired")
		case xaiv1.DeferredStatus_FAILED:
			if resp.Response != nil && resp.Response.Error != nil {
				return nil, fmt.Errorf("video generation failed (%s): %s", resp.Response.Error.Code, resp.Response.Error.Message)
			}
			return nil, fmt.Errorf("video generation failed")
		case xaiv1.DeferredStatus_PENDING:
		default:
			return nil, fmt.Errorf("unknown deferred video status: %s", resp.Status.String())
		}

		select {
		case <-pollCtx.Done():
			return nil, pollCtx.Err()
		case <-ticker.C:
		}
	}
}
