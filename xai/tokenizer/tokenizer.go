// Package tokenizer provides text tokenization functionality for xAI SDK.
package tokenizer

import (
	"context"
	"fmt"

	xaiv1 "github.com/ZaguanLabs/xai-sdk-go/proto/gen/go/xai/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// TokenizerServiceClient is an interface for the tokenizer service client.
type TokenizerServiceClient interface {
	EncodeText(ctx context.Context, req *xaiv1.EncodeTextRequest, opts ...grpc.CallOption) (*xaiv1.EncodeTextResponse, error)
	DecodeTokens(ctx context.Context, req *xaiv1.DecodeTokensRequest, opts ...grpc.CallOption) (*xaiv1.DecodeTokensResponse, error)
	CountTokens(ctx context.Context, req *xaiv1.CountTokensRequest, opts ...grpc.CallOption) (*xaiv1.CountTokensResponse, error)
}

// Client provides tokenization functionality.
type Client struct {
	grpcClient TokenizerServiceClient
}

// NewClient creates a new tokenizer client.
func NewClient(grpcClient TokenizerServiceClient) *Client {
	return &Client{
		grpcClient: grpcClient,
	}
}

// Encode converts text into tokens.
func (c *Client) Encode(ctx context.Context, text, model string) ([]int32, error) {
	if text == "" {
		return nil, fmt.Errorf("text is required")
	}
	if model == "" {
		return nil, fmt.Errorf("model is required")
	}

	req := &xaiv1.EncodeTextRequest{
		Text:  text,
		Model: model,
	}

	resp, err := c.grpcClient.EncodeText(ctx, req)
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
			case codes.NotFound:
				return nil, fmt.Errorf("model not found: %s", st.Message())
			case codes.Unavailable:
				return nil, fmt.Errorf("service unavailable: %s", st.Message())
			default:
				return nil, fmt.Errorf("encode text failed (%s): %s", st.Code().String(), st.Message())
			}
		}
		return nil, fmt.Errorf("encode text failed: %w", err)
	}

	if resp == nil {
		return nil, fmt.Errorf("received nil response")
	}

	return resp.Tokens, nil
}

// Decode converts tokens back into text.
func (c *Client) Decode(ctx context.Context, tokens []int32, model string) (string, error) {
	if len(tokens) == 0 {
		return "", fmt.Errorf("tokens are required")
	}
	if model == "" {
		return "", fmt.Errorf("model is required")
	}

	req := &xaiv1.DecodeTokensRequest{
		Tokens: tokens,
		Model:  model,
	}

	resp, err := c.grpcClient.DecodeTokens(ctx, req)
	if err != nil {
		// Handle specific gRPC errors
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.Unauthenticated:
				return "", fmt.Errorf("authentication failed: %s", st.Message())
			case codes.PermissionDenied:
				return "", fmt.Errorf("permission denied: %s", st.Message())
			case codes.InvalidArgument:
				return "", fmt.Errorf("invalid request: %s", st.Message())
			case codes.NotFound:
				return "", fmt.Errorf("model not found: %s", st.Message())
			case codes.Unavailable:
				return "", fmt.Errorf("service unavailable: %s", st.Message())
			default:
				return "", fmt.Errorf("decode tokens failed (%s): %s", st.Code().String(), st.Message())
			}
		}
		return "", fmt.Errorf("decode tokens failed: %w", err)
	}

	if resp == nil {
		return "", fmt.Errorf("received nil response")
	}

	return resp.Text, nil
}

// Count counts the number of tokens in text.
func (c *Client) Count(ctx context.Context, text, model string) (int32, error) {
	if text == "" {
		return 0, fmt.Errorf("text is required")
	}
	if model == "" {
		return 0, fmt.Errorf("model is required")
	}

	req := &xaiv1.CountTokensRequest{
		Text:  text,
		Model: model,
	}

	resp, err := c.grpcClient.CountTokens(ctx, req)
	if err != nil {
		// Handle specific gRPC errors
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.Unauthenticated:
				return 0, fmt.Errorf("authentication failed: %s", st.Message())
			case codes.PermissionDenied:
				return 0, fmt.Errorf("permission denied: %s", st.Message())
			case codes.InvalidArgument:
				return 0, fmt.Errorf("invalid request: %s", st.Message())
			case codes.NotFound:
				return 0, fmt.Errorf("model not found: %s", st.Message())
			case codes.Unavailable:
				return 0, fmt.Errorf("service unavailable: %s", st.Message())
			default:
				return 0, fmt.Errorf("count tokens failed (%s): %s", st.Code().String(), st.Message())
			}
		}
		return 0, fmt.Errorf("count tokens failed: %w", err)
	}

	if resp == nil {
		return 0, fmt.Errorf("received nil response")
	}

	return resp.TokenCount, nil
}

// CountWithDetails counts tokens and returns additional details.
func (c *Client) CountWithDetails(ctx context.Context, text, model string) (tokenCount int32, characterCount int32, err error) {
	if text == "" {
		return 0, 0, fmt.Errorf("text is required")
	}
	if model == "" {
		return 0, 0, fmt.Errorf("model is required")
	}

	req := &xaiv1.CountTokensRequest{
		Text:  text,
		Model: model,
	}

	resp, err := c.grpcClient.CountTokens(ctx, req)
	if err != nil {
		// Handle specific gRPC errors
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.Unauthenticated:
				return 0, 0, fmt.Errorf("authentication failed: %s", st.Message())
			case codes.PermissionDenied:
				return 0, 0, fmt.Errorf("permission denied: %s", st.Message())
			case codes.InvalidArgument:
				return 0, 0, fmt.Errorf("invalid request: %s", st.Message())
			case codes.NotFound:
				return 0, 0, fmt.Errorf("model not found: %s", st.Message())
			case codes.Unavailable:
				return 0, 0, fmt.Errorf("service unavailable: %s", st.Message())
			default:
				return 0, 0, fmt.Errorf("count tokens failed (%s): %s", st.Code().String(), st.Message())
			}
		}
		return 0, 0, fmt.Errorf("count tokens failed: %w", err)
	}

	if resp == nil {
		return 0, 0, fmt.Errorf("received nil response")
	}

	return resp.TokenCount, resp.CharacterCount, nil
}