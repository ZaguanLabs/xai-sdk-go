package batch

import (
	"context"
	"fmt"

	xaiv1 "github.com/ZaguanLabs/xai-sdk-go/proto/gen/go/xai/api/v1"
	"github.com/ZaguanLabs/xai-sdk-go/xai/chat"
	"google.golang.org/protobuf/proto"
)

type Client struct {
	grpcClient xaiv1.BatchMgmtClient
}

type ListOptions struct {
	Limit           int32
	PaginationToken string
}

func NewClient(grpcClient xaiv1.BatchMgmtClient) *Client {
	return &Client{grpcClient: grpcClient}
}

func (c *Client) Create(ctx context.Context, name string) (*xaiv1.Batch, error) {
	return c.CreateWithRequest(ctx, &xaiv1.CreateBatchRequest{Name: name})
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
		batchReq := BatchRequestFromChatRequest(req)
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
			batchReq := BatchRequestFromChatRequest(value)
			if batchReq != nil {
				batchRequests = append(batchRequests, batchReq)
			}
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

func BatchRequestFromChatRequest(req *chat.Request) *xaiv1.BatchRequest {
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
