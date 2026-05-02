package chat

import (
	"context"
	"fmt"
	"time"

	xaiv1 "github.com/ZaguanLabs/xai-sdk-go/proto/gen/go/xai/api/v1"
	"google.golang.org/grpc"
)

// Client provides a Python-style chat client wrapper on top of the generated gRPC client.
type Client struct {
	grpcClient xaiv1.ChatClient
}

// NewClient creates a new chat client wrapper.
func NewClient(grpcClient xaiv1.ChatClient) *Client {
	return &Client{grpcClient: grpcClient}
}

// Create creates a new chat request with the provided model and options.
func (c *Client) Create(model string, opts ...RequestOption) *Request {
	return NewRequest(model, opts...)
}

func (c *Client) Sample(ctx context.Context, req *Request) (*Response, error) {
	return req.Sample(ctx, c.grpcClient)
}

func (c *Client) SampleBatch(ctx context.Context, req *Request, n int32) ([]*Response, error) {
	return req.SampleBatch(ctx, c.grpcClient, n)
}

func (c *Client) Stream(ctx context.Context, req *Request) (*Stream, error) {
	return req.Stream(ctx, c.grpcClient)
}

func (c *Client) StreamBatch(ctx context.Context, req *Request, n int32) (*BatchStream, error) {
	return req.StreamBatch(ctx, c.grpcClient, n)
}

func (c *Client) Parse(ctx context.Context, req *Request, v any) error {
	return req.Parse(ctx, c.grpcClient, v)
}

func (c *Client) Defer(ctx context.Context, req *Request, timeout, interval time.Duration) (*DeferredResponse, error) {
	deferred := &DeferredRequest{proto: req.Proto()}
	result, err := deferred.Poll(ctx, c.grpcClient, interval, timeout)
	if err != nil {
		return nil, err
	}
	return result.Response, nil
}

func (c *Client) DeferBatch(ctx context.Context, req *Request, n int32, timeout, interval time.Duration) ([]*Response, error) {
	if n <= 0 {
		return nil, fmt.Errorf("n must be greater than 0")
	}
	batchReq, err := req.cloneWithN(n)
	if err != nil {
		return nil, err
	}
	deferred := &DeferredRequest{proto: batchReq}
	result, err := deferred.Poll(ctx, c.grpcClient, interval, timeout)
	if err != nil {
		return nil, err
	}
	if result == nil || result.Response == nil {
		return nil, nil
	}
	return splitResponses(result.Response.Proto()), nil
}

// GetCompletion forwards to the underlying gRPC chat client.
func (c *Client) GetCompletion(ctx context.Context, in *xaiv1.GetCompletionsRequest, opts ...grpc.CallOption) (*xaiv1.GetChatCompletionResponse, error) {
	return c.grpcClient.GetCompletion(ctx, in, opts...)
}

// GetCompletionChunk forwards to the underlying gRPC chat client.
func (c *Client) GetCompletionChunk(ctx context.Context, in *xaiv1.GetCompletionsRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[xaiv1.GetChatCompletionChunk], error) {
	return c.grpcClient.GetCompletionChunk(ctx, in, opts...)
}

// StartDeferredCompletion forwards to the underlying gRPC chat client.
func (c *Client) StartDeferredCompletion(ctx context.Context, in *xaiv1.GetCompletionsRequest, opts ...grpc.CallOption) (*xaiv1.StartDeferredResponse, error) {
	return c.grpcClient.StartDeferredCompletion(ctx, in, opts...)
}

// GetDeferredCompletion forwards to the underlying gRPC chat client.
func (c *Client) GetDeferredCompletion(ctx context.Context, in *xaiv1.GetDeferredRequest, opts ...grpc.CallOption) (*xaiv1.GetDeferredCompletionResponse, error) {
	return c.grpcClient.GetDeferredCompletion(ctx, in, opts...)
}

// GetStoredCompletion retrieves a stored completion and returns one response per output, matching Python semantics.
func (c *Client) GetStoredCompletion(ctx context.Context, responseID string, opts ...grpc.CallOption) ([]*Response, error) {
	response, err := c.grpcClient.GetStoredCompletion(ctx, &xaiv1.GetStoredCompletionRequest{ResponseId: responseID}, opts...)
	if err != nil {
		return nil, err
	}
	responses := make([]*Response, 0, len(response.Outputs))
	for i := range response.Outputs {
		outputProto := &xaiv1.GetChatCompletionResponse{
			Id:                response.Id,
			Model:             response.Model,
			Created:           response.Created,
			SystemFingerprint: response.SystemFingerprint,
			Usage:             response.Usage,
			Citations:         response.Citations,
			Outputs:           []*xaiv1.CompletionOutput{response.Outputs[i]},
		}
		responses = append(responses, &Response{proto: outputProto})
	}
	if len(response.Outputs) == 0 {
		responses = append(responses, &Response{proto: response})
	}
	return responses, nil
}

// DeleteStoredCompletion deletes a stored completion and returns its response ID, matching Python semantics.
func (c *Client) DeleteStoredCompletion(ctx context.Context, responseID string, opts ...grpc.CallOption) (string, error) {
	response, err := c.grpcClient.DeleteStoredCompletion(ctx, &xaiv1.DeleteStoredCompletionRequest{ResponseId: responseID}, opts...)
	if err != nil {
		return "", err
	}
	return response.ResponseId, nil
}

// GRPCClient returns the underlying generated gRPC chat client.
func (c *Client) GRPCClient() xaiv1.ChatClient {
	return c.grpcClient
}
