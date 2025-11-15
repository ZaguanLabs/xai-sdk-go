// Package models provides model information functionality for xAI SDK.
package models

import (
	"context"
	"fmt"

	xaiv1 "github.com/ZaguanLabs/xai-sdk-go/proto/gen/go/xai/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Model represents an xAI model.
type Model struct {
	id          string
	name        string
	description string
	maxTokens   int32
}

// ModelServiceClient is an interface for the models service client.
type ModelServiceClient interface {
	ListModels(ctx context.Context, req *xaiv1.ListModelsRequest, opts ...grpc.CallOption) (*xaiv1.ListModelsResponse, error)
	GetModel(ctx context.Context, req *xaiv1.GetModelRequest, opts ...grpc.CallOption) (*xaiv1.Model, error)
}

// Client provides model information functionality.
type Client struct {
	grpcClient ModelServiceClient
}

// NewClient creates a new models client.
func NewClient(grpcClient ModelServiceClient) *Client {
	return &Client{
		grpcClient: grpcClient,
	}
}

// List lists all available models.
func (c *Client) List(ctx context.Context) ([]*Model, error) {
	req := &xaiv1.ListModelsRequest{}
	
	resp, err := c.grpcClient.ListModels(ctx, req)
	if err != nil {
		// Handle specific gRPC errors
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.Unauthenticated:
				return nil, fmt.Errorf("authentication failed: %s", st.Message())
			case codes.PermissionDenied:
				return nil, fmt.Errorf("permission denied: %s", st.Message())
			case codes.Unavailable:
				return nil, fmt.Errorf("service unavailable: %s", st.Message())
			default:
				return nil, fmt.Errorf("list models failed (%s): %s", st.Code().String(), st.Message())
			}
		}
		return nil, fmt.Errorf("list models failed: %w", err)
	}

	if resp == nil {
		return nil, fmt.Errorf("received nil response")
	}

	// Convert response models
	models := make([]*Model, 0, len(resp.Models))
	for _, modelProto := range resp.Models {
		models = append(models, &Model{
			id:          modelProto.Id,
			name:        modelProto.Name,
			description: modelProto.Description,
			maxTokens:   modelProto.MaxTokens,
		})
	}

	return models, nil
}

// Get retrieves information about a specific model.
func (c *Client) Get(ctx context.Context, modelID string) (*Model, error) {
	if modelID == "" {
		return nil, fmt.Errorf("model ID is required")
	}

	req := &xaiv1.GetModelRequest{
		ModelId: modelID,
	}
	
	resp, err := c.grpcClient.GetModel(ctx, req)
	if err != nil {
		// Handle specific gRPC errors
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.Unauthenticated:
				return nil, fmt.Errorf("authentication failed: %s", st.Message())
			case codes.PermissionDenied:
				return nil, fmt.Errorf("permission denied: %s", st.Message())
			case codes.NotFound:
				return nil, fmt.Errorf("model not found: %s", st.Message())
			case codes.Unavailable:
				return nil, fmt.Errorf("service unavailable: %s", st.Message())
			default:
				return nil, fmt.Errorf("get model failed (%s): %s", st.Code().String(), st.Message())
			}
		}
		return nil, fmt.Errorf("get model failed: %w", err)
	}

	if resp == nil {
		return nil, fmt.Errorf("received nil response")
	}

	return &Model{
		id:          resp.Id,
		name:        resp.Name,
		description: resp.Description,
		maxTokens:   resp.MaxTokens,
	}, nil
}

// Model methods

// ID returns the model ID.
func (m *Model) ID() string {
	return m.id
}

// Name returns the model name.
func (m *Model) Name() string {
	return m.name
}

// Description returns the model description.
func (m *Model) Description() string {
	return m.description
}

// MaxTokens returns the maximum number of tokens this model supports.
func (m *Model) MaxTokens() int32 {
	return m.maxTokens
}

// String returns a string representation of the model.
func (m *Model) String() string {
	return fmt.Sprintf("Model{ID: %s, Name: %s, MaxTokens: %d}", m.id, m.name, m.maxTokens)
}