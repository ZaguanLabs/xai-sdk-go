// Package deferred provides a client for the xAI Deferred Completions API.
package deferred

import (
	"context"

	xaiv1 "github.com/ZaguanLabs/xai-sdk-go/proto/gen/go/xai/api/v1"
	"github.com/ZaguanLabs/xai-sdk-go/xai/internal/rest"
	"google.golang.org/protobuf/encoding/protojson"
)

// Client provides access to the xAI Deferred Completions API.
type Client struct {
	restClient *rest.Client
}

// NewClient creates a new Deferred Completions API client.
func NewClient(restClient *rest.Client) *Client {
	return &Client{
		restClient: restClient,
	}
}

// StartResponse represents the response from starting a deferred completion.
type StartResponse struct {
	RequestID string
}

// Status represents the status of a deferred completion.
type Status struct {
	RequestID string
	Status    xaiv1.DeferredStatus
	// Result will contain the completion when status is DONE
	Result interface{}
}

// Start initiates a deferred completion.
// This is typically used for long-running completions that will be retrieved later.
func (c *Client) Start(ctx context.Context, chatRequest interface{}) (*StartResponse, error) {
	if c.restClient == nil {
		return nil, ErrClientNotInitialized
	}

	// The chatRequest should be a chat completion request
	// We'll send it to the deferred endpoint
	resp, err := c.restClient.Post(ctx, "/deferred/start", chatRequest)
	if err != nil {
		return nil, err
	}

	var startResp xaiv1.StartDeferredResponse
	if err := protojson.Unmarshal(resp.Body, &startResp); err != nil {
		return nil, err
	}

	return &StartResponse{
		RequestID: startResp.RequestId,
	}, nil
}

// Get retrieves the status and result of a deferred completion.
func (c *Client) Get(ctx context.Context, requestID string) (*Status, error) {
	if c.restClient == nil {
		return nil, ErrClientNotInitialized
	}

	req := &xaiv1.GetDeferredRequest{
		RequestId: requestID,
	}

	jsonData, err := protojson.Marshal(req)
	if err != nil {
		return nil, err
	}

	resp, err := c.restClient.Post(ctx, "/deferred/get", jsonData)
	if err != nil {
		return nil, err
	}

	// The response structure depends on the status
	// For now, we'll return a basic status
	// TODO: Parse the actual completion result when status is DONE
	var statusResp struct {
		RequestID string               `json:"request_id"`
		Status    xaiv1.DeferredStatus `json:"status"`
		Result    interface{}          `json:"result,omitempty"`
	}

	if err := resp.DecodeJSON(&statusResp); err != nil {
		return nil, err
	}

	return &Status{
		RequestID: statusResp.RequestID,
		Status:    statusResp.Status,
		Result:    statusResp.Result,
	}, nil
}
