package chat

import (
	"context"

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
