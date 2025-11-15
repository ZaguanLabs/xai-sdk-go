// Package auth provides authentication validation functionality for xAI SDK.
package auth

import (
	"context"
	"fmt"

	xaiv1 "github.com/ZaguanLabs/xai-sdk-go/proto/gen/go/xai/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ValidationResult represents the result of API key validation.
type ValidationResult struct {
	valid        bool
	message      string
	organization string
	project      string
}

// AuthServiceClient is an interface for the auth service client.
type AuthServiceClient interface {
	ValidateKey(ctx context.Context, req *xaiv1.ValidateKeyRequest, opts ...grpc.CallOption) (*xaiv1.ValidateKeyResponse, error)
}

// Client provides authentication validation functionality.
type Client struct {
	grpcClient AuthServiceClient
}

// NewClient creates a new auth validation client.
func NewClient(grpcClient AuthServiceClient) *Client {
	return &Client{
		grpcClient: grpcClient,
	}
}

// Validate validates an API key.
func (c *Client) Validate(ctx context.Context, apiKey string) (*ValidationResult, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("API key is required")
	}

	req := &xaiv1.ValidateKeyRequest{
		ApiKey: apiKey,
	}

	resp, err := c.grpcClient.ValidateKey(ctx, req)
	if err != nil {
		// Handle specific gRPC errors
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.Unauthenticated:
				return &ValidationResult{
					valid:   false,
					message: fmt.Sprintf("authentication failed: %s", st.Message()),
				}, nil
			case codes.PermissionDenied:
				return &ValidationResult{
					valid:   false,
					message: fmt.Sprintf("permission denied: %s", st.Message()),
				}, nil
			case codes.InvalidArgument:
				return &ValidationResult{
					valid:   false,
					message: fmt.Sprintf("invalid request: %s", st.Message()),
				}, nil
			case codes.Unavailable:
				return &ValidationResult{
					valid:   false,
					message: fmt.Sprintf("service unavailable: %s", st.Message()),
				}, nil
			default:
				return nil, fmt.Errorf("validate key failed (%s): %s", st.Code().String(), st.Message())
			}
		}
		return nil, fmt.Errorf("validate key failed: %w", err)
	}

	if resp == nil {
		return nil, fmt.Errorf("received nil response")
	}

	return &ValidationResult{
		valid:        resp.Valid,
		message:      resp.Message,
		organization: resp.Organization,
		project:      resp.Project,
	}, nil
}

// ValidationResult methods

// IsValid returns whether the API key is valid.
func (v *ValidationResult) IsValid() bool {
	return v.valid
}

// Message returns the validation message.
func (v *ValidationResult) Message() string {
	return v.message
}

// Organization returns the organization name.
func (v *ValidationResult) Organization() string {
	return v.organization
}

// Project returns the project name.
func (v *ValidationResult) Project() string {
	return v.project
}

// String returns a string representation of the validation result.
func (v *ValidationResult) String() string {
	if v.valid {
		return fmt.Sprintf("ValidationResult{Valid: true, Organization: %s, Project: %s}", v.organization, v.project)
	}
	return fmt.Sprintf("ValidationResult{Valid: false, Message: %s}", v.message)
}