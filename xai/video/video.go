package video

import (
	"context"
	"fmt"
	"time"

	xaiv1 "github.com/ZaguanLabs/xai-sdk-go/proto/gen/go/xai/api/v1"
)

const (
	DefaultPollTimeout  = 10 * time.Minute
	DefaultPollInterval = 2 * time.Second
)

type Client struct {
	grpcClient xaiv1.VideoClient
}

type GenerateOptions struct {
	ImageURL    string
	VideoURL    string
	Duration    *int32
	AspectRatio *xaiv1.VideoAspectRatio
	Resolution  *xaiv1.VideoResolution
	Timeout     time.Duration
	Interval    time.Duration
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

func (c *Client) Generate(ctx context.Context, req *xaiv1.GenerateVideoRequest) (*xaiv1.StartDeferredResponse, error) {
	if c.grpcClient == nil {
		return nil, fmt.Errorf("video client not initialized")
	}
	return c.grpcClient.GenerateVideo(ctx, req)
}

func (c *Client) Start(ctx context.Context, prompt, model string, opts *GenerateOptions) (*xaiv1.StartDeferredResponse, error) {
	return c.Generate(ctx, NewGenerateRequestWithOptions(prompt, model, opts))
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

func (c *Client) GenerateAndPoll(ctx context.Context, prompt, model string, opts *GenerateOptions) (*xaiv1.VideoResponse, error) {
	start, err := c.Start(ctx, prompt, model, opts)
	if err != nil {
		return nil, err
	}

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
		resp, err := c.Get(pollCtx, start.RequestId)
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
