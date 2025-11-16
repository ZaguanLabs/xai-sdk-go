// Package auth provides a client for the xAI Auth API.
package auth

import (
	"context"
	"time"

	xaiv1 "github.com/ZaguanLabs/xai-sdk-go/proto/gen/go/xai/v1"
)

// Client provides access to the xAI Auth API.
type Client struct {
	// Note: Auth API is currently REST-based in the Python SDK
	// This wrapper is prepared for when gRPC support is added
}

// NewClient creates a new Auth API client.
func NewClient() *Client {
	return &Client{}
}

// ApiKey represents an API key with metadata.
type ApiKey struct {
	RedactedApiKey string
	UserID         string
	Name           string
	CreateTime     time.Time
	TeamID         string
	ACLs           []string
	ApiKeyID       string
	ModifyTime     time.Time
	ApiKeyBlocked  bool
	ModifiedBy     string
	Disabled       bool
	TeamBlocked    bool
}

// fromProto converts a proto ApiKey to an ApiKey.
func fromProto(pk *xaiv1.ApiKey) *ApiKey {
	if pk == nil {
		return nil
	}

	ak := &ApiKey{
		RedactedApiKey: pk.RedactedApiKey,
		UserID:         pk.UserId,
		Name:           pk.Name,
		TeamID:         pk.TeamId,
		ACLs:           pk.Acls,
		ApiKeyID:       pk.ApiKeyId,
		ApiKeyBlocked:  pk.ApiKeyBlocked,
		ModifiedBy:     pk.ModifiedBy,
		Disabled:       pk.Disabled,
		TeamBlocked:    pk.TeamBlocked,
	}

	if pk.CreateTime != nil {
		ak.CreateTime = pk.CreateTime.AsTime()
	}

	if pk.ModifyTime != nil {
		ak.ModifyTime = pk.ModifyTime.AsTime()
	}

	return ak
}

// ValidateKey validates an API key and returns its metadata.
// Note: This method is a placeholder until gRPC support is added.
func (c *Client) ValidateKey(ctx context.Context, apiKey string) (*ApiKey, error) {
	// TODO: Implement when gRPC service is available
	// For now, this would need to use REST API
	return nil, ErrNotImplemented
}

// GetKey retrieves API key metadata by ID.
// Note: This method is a placeholder until gRPC support is added.
func (c *Client) GetKey(ctx context.Context, apiKeyID string) (*ApiKey, error) {
	// TODO: Implement when gRPC service is available
	// For now, this would need to use REST API
	return nil, ErrNotImplemented
}

// ListKeys lists API keys for the authenticated user.
// Note: This method is a placeholder until gRPC support is added.
func (c *Client) ListKeys(ctx context.Context) ([]*ApiKey, error) {
	// TODO: Implement when gRPC service is available
	// For now, this would need to use REST API
	return nil, ErrNotImplemented
}
