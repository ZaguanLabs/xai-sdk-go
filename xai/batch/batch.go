package batch

import (
	"context"
	"fmt"

	xaiv1 "github.com/ZaguanLabs/xai-sdk-go/proto/gen/go/xai/api/v1"
	"github.com/ZaguanLabs/xai-sdk-go/xai/chat"
	"github.com/ZaguanLabs/xai-sdk-go/xai/image"
	"github.com/ZaguanLabs/xai-sdk-go/xai/video"
	"google.golang.org/protobuf/proto"
)

type Client struct {
	grpcClient xaiv1.BatchMgmtClient
}

type ListOptions struct {
	Limit           int32
	PaginationToken string
}

type ListBatchResultsResponse struct {
	proto *xaiv1.ListBatchResultsResponse
}

type Result struct {
	proto *xaiv1.BatchResult
}

func NewClient(grpcClient xaiv1.BatchMgmtClient) *Client {
	return &Client{grpcClient: grpcClient}
}

func (c *Client) Create(ctx context.Context, name string) (*xaiv1.Batch, error) {
	return c.CreateWithRequest(ctx, &xaiv1.CreateBatchRequest{Name: name})
}

func (c *Client) CreateFromFile(ctx context.Context, name, inputFileID string) (*xaiv1.Batch, error) {
	return c.CreateWithRequest(ctx, &xaiv1.CreateBatchRequest{Name: name, InputFileId: inputFileID})
}

func (c *Client) CreateWithRequest(ctx context.Context, req *xaiv1.CreateBatchRequest) (*xaiv1.Batch, error) {
	if c.grpcClient == nil {
		return nil, fmt.Errorf("batch client not initialized")
	}
	return c.grpcClient.CreateBatch(ctx, req)
}

func (c *Client) Get(ctx context.Context, batchID string) (*xaiv1.Batch, error) {
	if c.grpcClient == nil {
		return nil, fmt.Errorf("batch client not initialized")
	}
	return c.grpcClient.GetBatch(ctx, &xaiv1.GetBatchRequest{BatchId: batchID})
}

func (c *Client) List(ctx context.Context, opts *ListOptions) ([]*xaiv1.Batch, string, error) {
	if c.grpcClient == nil {
		return nil, "", fmt.Errorf("batch client not initialized")
	}

	req := &xaiv1.ListBatchesRequest{}
	if opts != nil {
		req.Limit = opts.Limit
		if opts.PaginationToken != "" {
			req.PaginationToken = &opts.PaginationToken
		}
	}

	resp, err := c.grpcClient.ListBatches(ctx, req)
	if err != nil {
		return nil, "", err
	}

	return resp.Batches, resp.GetPaginationToken(), nil
}

func (c *Client) Cancel(ctx context.Context, batchID string) (*xaiv1.Batch, error) {
	if c.grpcClient == nil {
		return nil, fmt.Errorf("batch client not initialized")
	}
	return c.grpcClient.CancelBatch(ctx, &xaiv1.CancelBatchRequest{BatchId: batchID})
}

func (c *Client) AddRequests(ctx context.Context, batchID string, requests ...*xaiv1.BatchRequest) error {
	if c.grpcClient == nil {
		return fmt.Errorf("batch client not initialized")
	}

	_, err := c.grpcClient.AddBatchRequests(ctx, &xaiv1.AddBatchRequestsRequest{
		BatchId:       batchID,
		BatchRequests: requests,
	})
	return err
}

func (c *Client) AddChatRequests(ctx context.Context, batchID string, requests ...*chat.Request) error {
	batchRequests := make([]*xaiv1.BatchRequest, 0, len(requests))
	for _, req := range requests {
		batchReq := RequestFromChatRequest(req)
		if batchReq != nil {
			batchRequests = append(batchRequests, batchReq)
		}
	}

	return c.AddRequests(ctx, batchID, batchRequests...)
}

func (c *Client) Add(ctx context.Context, batchID string, requests ...interface{}) error {
	batchRequests := make([]*xaiv1.BatchRequest, 0, len(requests))
	for _, request := range requests {
		switch value := request.(type) {
		case *xaiv1.BatchRequest:
			if value != nil {
				batchRequests = append(batchRequests, value)
			}
		case *chat.Request:
			batchReq := RequestFromChatRequest(value)
			if batchReq != nil {
				batchRequests = append(batchRequests, batchReq)
			}
		case *image.GenerateRequest:
			batchReq := RequestFromImageRequest(value, "")
			if batchReq != nil {
				batchRequests = append(batchRequests, batchReq)
			}
		case *xaiv1.GenerateVideoRequest:
			batchRequests = append(batchRequests, RequestFromVideoRequest(value, ""))
		default:
			return fmt.Errorf("unsupported batch request type: %T", request)
		}
	}
	return c.AddRequests(ctx, batchID, batchRequests...)
}

func (c *Client) ListRequestMetadata(ctx context.Context, batchID string, opts *ListOptions) ([]*xaiv1.BatchRequestMetadata, string, error) {
	if c.grpcClient == nil {
		return nil, "", fmt.Errorf("batch client not initialized")
	}

	req := &xaiv1.ListBatchRequestMetadataRequest{BatchId: batchID}
	if opts != nil {
		req.Limit = opts.Limit
		if opts.PaginationToken != "" {
			req.PaginationToken = &opts.PaginationToken
		}
	}

	resp, err := c.grpcClient.ListBatchRequestMetadata(ctx, req)
	if err != nil {
		return nil, "", err
	}

	return resp.BatchRequestMetadata, resp.GetPaginationToken(), nil
}

func (c *Client) ListBatchRequests(ctx context.Context, batchID string, opts *ListOptions) ([]*xaiv1.BatchRequestMetadata, string, error) {
	return c.ListRequestMetadata(ctx, batchID, opts)
}

func (c *Client) ListResults(ctx context.Context, batchID string, opts *ListOptions) ([]*xaiv1.BatchResult, string, error) {
	if c.grpcClient == nil {
		return nil, "", fmt.Errorf("batch client not initialized")
	}

	req := &xaiv1.ListBatchResultsRequest{BatchId: batchID}
	if opts != nil {
		req.Limit = opts.Limit
		if opts.PaginationToken != "" {
			req.PaginationToken = &opts.PaginationToken
		}
	}

	resp, err := c.grpcClient.ListBatchResults(ctx, req)
	if err != nil {
		return nil, "", err
	}

	return resp.Results, resp.GetPaginationToken(), nil
}

func (c *Client) ListBatchResults(ctx context.Context, batchID string, opts *ListOptions) ([]*xaiv1.BatchResult, string, error) {
	return c.ListResults(ctx, batchID, opts)
}

func (c *Client) GetRequestResult(ctx context.Context, batchID, batchRequestID string) (*xaiv1.GetBatchRequestResultResponse, error) {
	if c.grpcClient == nil {
		return nil, fmt.Errorf("batch client not initialized")
	}
	return c.grpcClient.GetBatchRequestResult(ctx, &xaiv1.GetBatchRequestResultRequest{
		BatchId:        batchID,
		BatchRequestId: batchRequestID,
	})
}

func RequestFromChatRequest(req *chat.Request) *xaiv1.BatchRequest {
	if req == nil || req.Proto() == nil {
		return nil
	}

	completionRequest, ok := proto.Clone(req.Proto()).(*xaiv1.GetCompletionsRequest)
	if !ok {
		return nil
	}

	batchReq := &xaiv1.BatchRequest{
		Request: &xaiv1.BatchRequest_CompletionRequest{
			CompletionRequest: completionRequest,
		},
	}

	if batchRequestID := req.BatchRequestID(); batchRequestID != "" {
		batchReq.BatchRequestId = &batchRequestID
	}

	return batchReq
}

func RequestFromImageRequest(req *image.GenerateRequest, batchRequestID string) *xaiv1.BatchRequest {
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

func RequestFromVideoRequest(req *xaiv1.GenerateVideoRequest, batchRequestID string) *xaiv1.BatchRequest {
	if req == nil {
		return nil
	}
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

func PrepareVideoRequest(prompt, model, batchRequestID string, opts *video.GenerateOptions) *xaiv1.BatchRequest {
	return video.Prepare(prompt, model, batchRequestID, opts)
}

func NewListBatchResultsResponse(proto *xaiv1.ListBatchResultsResponse) *ListBatchResultsResponse {
	return &ListBatchResultsResponse{proto: proto}
}

func (r *ListBatchResultsResponse) Results() []*Result {
	if r == nil || r.proto == nil {
		return nil
	}
	results := make([]*Result, len(r.proto.Results))
	for i, result := range r.proto.Results {
		results[i] = &Result{proto: result}
	}
	return results
}

func (r *ListBatchResultsResponse) Succeeded() []*Result {
	results := r.Results()
	succeeded := make([]*Result, 0, len(results))
	for _, result := range results {
		if result.IsSuccess() {
			succeeded = append(succeeded, result)
		}
	}
	return succeeded
}

func (r *ListBatchResultsResponse) Failed() []*Result {
	results := r.Results()
	failed := make([]*Result, 0, len(results))
	for _, result := range results {
		if result.HasError() {
			failed = append(failed, result)
		}
	}
	return failed
}

func (r *ListBatchResultsResponse) PaginationToken() string {
	if r == nil || r.proto == nil {
		return ""
	}
	return r.proto.GetPaginationToken()
}

func (r *ListBatchResultsResponse) Proto() *xaiv1.ListBatchResultsResponse {
	if r == nil {
		return nil
	}
	return r.proto
}

func NewResult(proto *xaiv1.BatchResult) *Result {
	return &Result{proto: proto}
}

func (r *Result) BatchRequestID() string {
	if r == nil || r.proto == nil {
		return ""
	}
	return r.proto.GetBatchRequestId()
}

func (r *Result) Response() *xaiv1.BatchResultData {
	if r == nil || r.proto == nil {
		return nil
	}
	return r.proto.GetResponse()
}

func (r *Result) ImageResponse() *xaiv1.ImageResponse {
	if response := r.Response(); response != nil {
		return response.GetImageResponse()
	}
	return nil
}

func (r *Result) VideoResponse() *xaiv1.VideoResponse {
	if response := r.Response(); response != nil {
		return response.GetVideoResponse()
	}
	return nil
}

func (r *Result) HasError() bool {
	return r != nil && r.proto != nil && r.proto.GetError() != nil
}

func (r *Result) IsSuccess() bool {
	return r != nil && r.proto != nil && r.proto.GetResponse() != nil
}

func (r *Result) ErrorMessage() string {
	if r == nil || r.proto == nil || r.proto.GetError() == nil {
		return ""
	}
	return r.proto.GetError().Message
}

func (r *Result) Proto() *xaiv1.BatchResult {
	if r == nil {
		return nil
	}
	return r.proto
}
