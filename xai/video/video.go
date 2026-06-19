package video

import (
	"context"
	"fmt"
	"time"

	xaiv1 "github.com/ZaguanLabs/xai-sdk-go/proto/gen/go/xai/api/v1"
	"github.com/ZaguanLabs/xai-sdk-go/xai/cost"
	"github.com/ZaguanLabs/xai-sdk-go/xai/files"
)

const (
	DefaultPollTimeout  = 10 * time.Minute
	DefaultPollInterval = 2 * time.Second
)

type Client struct {
	grpcClient xaiv1.VideoClient
}

type GenerateOptions struct {
	ImageURL              string
	ImageFileID           string
	VideoURL              string
	VideoFileID           string
	ReferenceImageURLs    []string
	ReferenceImageFileIDs []string
	Storage               *files.StorageOptions
	Duration              *int32
	AspectRatio           *xaiv1.VideoAspectRatio
	Resolution            *xaiv1.VideoResolution
	Timeout               time.Duration
	Interval              time.Duration
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

func (r *Response) FileOutput() *xaiv1.FileOutput {
	if r == nil || r.proto == nil || r.proto.Video == nil {
		return nil
	}
	return r.proto.Video.GetFileOutput()
}

func (r *Response) StorageError() string {
	if r == nil || r.proto == nil || r.proto.Video == nil {
		return ""
	}
	return r.proto.Video.GetStorageError()
}

func (r *Response) PublicURL() string {
	if output := r.FileOutput(); output != nil {
		return output.GetPublicUrl()
	}
	return ""
}

func (r *Response) PublicURLError() string {
	if output := r.FileOutput(); output != nil {
		return output.GetPublicUrlError()
	}
	return ""
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
		req.Image = imageURLContent(opts.ImageURL)
	} else if opts.ImageFileID != "" {
		req.Image = imageFileContent(opts.ImageFileID)
	}
	if opts.VideoURL != "" {
		req.Video = videoURLContent(opts.VideoURL)
	} else if opts.VideoFileID != "" {
		req.Video = videoFileContent(opts.VideoFileID)
	}
	for _, fileID := range opts.ReferenceImageFileIDs {
		req.ReferenceImages = append(req.ReferenceImages, imageFileContent(fileID))
	}
	for _, imageURL := range opts.ReferenceImageURLs {
		req.ReferenceImages = append(req.ReferenceImages, imageURLContent(imageURL))
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
	if opts.Storage != nil {
		req.StorageOptions = opts.Storage.Proto()
	}
	return req
}

func NewExtendRequest(prompt, model, videoURL string, duration *int32) *xaiv1.ExtendVideoRequest {
	return NewExtendRequestWithOptions(prompt, model, videoURL, duration, nil)
}

func NewExtendRequestWithOptions(prompt, model, videoURL string, duration *int32, opts *GenerateOptions) *xaiv1.ExtendVideoRequest {
	video := videoURLContent(videoURL)
	if opts != nil && opts.VideoFileID != "" {
		video = videoFileContent(opts.VideoFileID)
	}
	req := &xaiv1.ExtendVideoRequest{
		Prompt:   prompt,
		Model:    model,
		Video:    video,
		Duration: duration,
	}
	if opts != nil && opts.Storage != nil {
		req.StorageOptions = opts.Storage.Proto()
	}
	return req
}

func NewExtendRequestFromFileID(prompt, model, videoFileID string, duration *int32) *xaiv1.ExtendVideoRequest {
	return &xaiv1.ExtendVideoRequest{
		Prompt:   prompt,
		Model:    model,
		Video:    videoFileContent(videoFileID),
		Duration: duration,
	}
}

func imageURLContent(imageURL string) *xaiv1.ImageUrlContent {
	return &xaiv1.ImageUrlContent{
		ImageUrl: imageURL,
		Detail:   xaiv1.ImageDetail_DETAIL_AUTO,
	}
}

func imageFileContent(fileID string) *xaiv1.ImageUrlContent {
	return &xaiv1.ImageUrlContent{
		FileId: fileID,
		Detail: xaiv1.ImageDetail_DETAIL_AUTO,
	}
}

func videoURLContent(videoURL string) *xaiv1.VideoUrlContent {
	return &xaiv1.VideoUrlContent{Url: videoURL}
}

func videoFileContent(fileID string) *xaiv1.VideoUrlContent {
	return &xaiv1.VideoUrlContent{FileId: fileID}
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

func PrepareExtension(prompt, model, videoURL, batchRequestID string, duration *int32) *xaiv1.BatchRequest {
	return PrepareExtensionWithOptions(prompt, model, videoURL, batchRequestID, duration, nil)
}

func PrepareExtensionWithOptions(prompt, model, videoURL, batchRequestID string, duration *int32, opts *GenerateOptions) *xaiv1.BatchRequest {
	req := NewExtendRequestWithOptions(prompt, model, videoURL, duration, opts)
	batchReq := &xaiv1.BatchRequest{
		Request: &xaiv1.BatchRequest_VideoExtensionRequest{
			VideoExtensionRequest: req,
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

func (c *Client) PrepareExtension(prompt, model, videoURL, batchRequestID string, duration *int32) *xaiv1.BatchRequest {
	return PrepareExtension(prompt, model, videoURL, batchRequestID, duration)
}

func (c *Client) PrepareExtensionWithOptions(prompt, model, videoURL, batchRequestID string, duration *int32, opts *GenerateOptions) *xaiv1.BatchRequest {
	return PrepareExtensionWithOptions(prompt, model, videoURL, batchRequestID, duration, opts)
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
	return c.ExtendStartWithOptions(ctx, prompt, model, videoURL, duration, nil)
}

func (c *Client) ExtendStartWithOptions(ctx context.Context, prompt, model, videoURL string, duration *int32, opts *GenerateOptions) (*xaiv1.StartDeferredResponse, error) {
	if c.grpcClient == nil {
		return nil, fmt.Errorf("video client not initialized")
	}
	return c.grpcClient.ExtendVideo(ctx, NewExtendRequestWithOptions(prompt, model, videoURL, duration, opts))
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
	start, err := c.ExtendStartWithOptions(ctx, prompt, model, videoURL, duration, opts)
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
