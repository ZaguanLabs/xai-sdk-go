// Package auth provides a client for the xAI Auth API.
package auth

import (
	"context"
	"fmt"
	"time"

	xaiv1 "github.com/ZaguanLabs/xai-sdk-go/proto/gen/go/xai/api/v1"
	"github.com/ZaguanLabs/xai-sdk-go/xai/internal/rest"
	"google.golang.org/protobuf/encoding/protojson"
)

// Client provides access to the xAI Auth API.
type Client struct {
	restClient *rest.Client
}

// NewClient creates a new Auth API client.
func NewClient(restClient *rest.Client) *Client {
	return &Client{
		restClient: restClient,
	}
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
func (c *Client) ValidateKey(ctx context.Context, apiKey string) (*ApiKey, error) {
	if c.restClient == nil {
		return nil, ErrClientNotInitialized
	}

	// Create validation request with the API key
	req := map[string]string{"api_key": apiKey}

	resp, err := c.restClient.Post(ctx, "/auth/validate", req)
	if err != nil {
		return nil, err
	}

	var key xaiv1.ApiKey
	if err := protojson.Unmarshal(resp.Body, &key); err != nil {
		return nil, err
	}

	return fromProto(&key), nil
}

// GetKey retrieves API key metadata by ID.
func (c *Client) GetKey(ctx context.Context, apiKeyID string) (*ApiKey, error) {
	if c.restClient == nil {
		return nil, ErrClientNotInitialized
	}

	resp, err := c.restClient.Get(ctx, fmt.Sprintf("/auth/keys/%s", apiKeyID))
	if err != nil {
		return nil, err
	}

	var key xaiv1.ApiKey
	if err := protojson.Unmarshal(resp.Body, &key); err != nil {
		return nil, err
	}

	return fromProto(&key), nil
}

// ListKeys lists API keys for the authenticated user.
func (c *Client) ListKeys(ctx context.Context) ([]*ApiKey, error) {
	if c.restClient == nil {
		return nil, ErrClientNotInitialized
	}

	resp, err := c.restClient.Get(ctx, "/auth/keys")
	if err != nil {
		return nil, err
	}

	// Parse response as array of API keys
	var keysResponse struct {
		Keys []*xaiv1.ApiKey `json:"keys"`
	}
	if err := resp.DecodeJSON(&keysResponse); err != nil {
		return nil, err
	}

	keys := make([]*ApiKey, len(keysResponse.Keys))
	for i, k := range keysResponse.Keys {
		keys[i] = fromProto(k)
	}

	return keys, nil
}
