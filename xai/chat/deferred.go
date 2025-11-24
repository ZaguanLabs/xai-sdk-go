// Package chat provides deferred and stored chat functionality for xAI SDK.
package chat

import (
	"context"
	"fmt"
	"time"

	xaiv1 "github.com/ZaguanLabs/xai-sdk-go/proto/gen/go/xai/api/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// DeferredRequest represents a deferred chat completion request.
type DeferredRequest struct {
	proto *xaiv1.GetCompletionsRequest
}

// DeferredRequestOption represents a functional option for DeferredRequest.
type DeferredRequestOption func(*DeferredRequest)

// NewDeferredRequest creates a new deferred chat completion request.
func NewDeferredRequest(model string, opts ...DeferredRequestOption) *DeferredRequest {
	req := &DeferredRequest{
		proto: &xaiv1.GetCompletionsRequest{
			Model: model,
		},
	}

	// Apply functional options
	for _, opt := range opts {
		opt(req)
	}

	return req
}

// WithStoreMessages enables message storage for the request.
func (r *DeferredRequest) WithStoreMessages(store bool) *DeferredRequest {
	r.proto.StoreMessages = store
	return r
}

// WithPreviousResponseID sets the previous response ID for conversation continuation.
func (r *DeferredRequest) WithPreviousResponseID(responseID string) *DeferredRequest {
	r.proto.PreviousResponseId = &responseID
	return r
}

// WithEncryptedContent enables encrypted content for the request.
func (r *DeferredRequest) WithEncryptedContent(encrypted bool) *DeferredRequest {
	r.proto.UseEncryptedContent = encrypted
	return r
}

// Model returns the model for the request.
func (r *DeferredRequest) Model() string {
	return r.proto.Model
}

// SetModel sets the model for the request.
func (r *DeferredRequest) SetModel(model string) *DeferredRequest {
	r.proto.Model = model
	return r
}

// AppendMessage appends a message to the request.
func (r *DeferredRequest) AppendMessage(msg Message) *DeferredRequest {
	r.proto.Messages = append(r.proto.Messages, msg.Proto())
	return r
}

// SetMessages sets all messages for the request.
func (r *DeferredRequest) SetMessages(messages ...Message) *DeferredRequest {
	r.proto.Messages = make([]*xaiv1.Message, 0, len(messages))
	for _, msg := range messages {
		r.proto.Messages = append(r.proto.Messages, msg.Proto())
	}
	return r
}

// WithTemperature sets the temperature for the request.
func (r *DeferredRequest) WithTemperature(temp float32) *DeferredRequest {
	r.proto.Temperature = &temp
	return r
}

// WithMaxTokens sets the maximum number of tokens for the request.
func (r *DeferredRequest) WithMaxTokens(maxTokens int32) *DeferredRequest {
	r.proto.MaxTokens = &maxTokens
	return r
}

// Proto returns the underlying protobuf request.
func (r *DeferredRequest) Proto() *xaiv1.GetCompletionsRequest {
	return r.proto
}

// Submit submits a deferred chat completion request.
func (r *DeferredRequest) Submit(ctx context.Context, client ServiceClient) (*DeferredResponse, error) {
	if client == nil {
		return nil, fmt.Errorf("chat client is nil")
	}
	if r.proto == nil {
		return nil, fmt.Errorf("request proto is nil")
	}
	if r.proto.Model == "" {
		return nil, fmt.Errorf("model is required")
	}

	// Submit the request
	resp, err := client.GetCompletion(ctx, r.proto)
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
			case codes.ResourceExhausted:
				return nil, fmt.Errorf("quota exceeded: %s", st.Message())
			case codes.Unavailable:
				return nil, fmt.Errorf("service unavailable: %s", st.Message())
			case codes.DeadlineExceeded:
				return nil, fmt.Errorf("request timeout: %s", st.Message())
			default:
				return nil, fmt.Errorf("deferred request failed (%s): %s", st.Code().String(), st.Message())
			}
		}
		return nil, fmt.Errorf("deferred request failed: %w", err)
	}

	if resp == nil {
		return nil, fmt.Errorf("received nil response")
	}

	return &DeferredResponse{proto: resp}, nil
}

// DeferredResponse represents a deferred chat completion response.
type DeferredResponse struct {
	proto *xaiv1.GetChatCompletionResponse
}

// ID returns the response ID.
func (r *DeferredResponse) ID() string {
	if r.proto == nil {
		return ""
	}
	return r.proto.Id
}

// Status returns the status of the deferred request.
// Returns "pending" if status is not available.
func (r *DeferredResponse) Status() string {
	// Status information is provided by the GetDeferredCompletionResponse
	// which wraps this response. This method returns a default value.
	return "pending"
}

// CreatedAt returns the creation time of the deferred request.
// Returns zero time if creation time is not available in the response.
func (r *DeferredResponse) CreatedAt() time.Time {
	if r.proto != nil && r.proto.Created != nil {
		return r.proto.Created.AsTime()
	}
	return time.Time{}
}

// CompletedAt returns the completion time of the deferred request.
// Returns zero time if completion time is not available.
// Note: Completion time is typically tracked externally to the response.
func (r *DeferredResponse) CompletedAt() time.Time {
	// Completion time is not stored in GetChatCompletionResponse
	return time.Time{}
}

// Proto returns the underlying protobuf response.
func (r *DeferredResponse) Proto() *xaiv1.GetChatCompletionResponse {
	return r.proto
}

// PollResult represents the result of polling a deferred request.
type PollResult struct {
	Response *DeferredResponse
	Done     bool
}

// Poll polls a deferred request until completion or timeout.
func (r *DeferredRequest) Poll(ctx context.Context, client ServiceClient, interval time.Duration, timeout time.Duration) (*PollResult, error) {
	if client == nil {
		return nil, fmt.Errorf("chat client is nil")
	}
	if r.proto == nil {
		return nil, fmt.Errorf("request proto is nil")
	}

	// Create a context with timeout
	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Submit the request
	response, err := r.Submit(timeoutCtx, client)
	if err != nil {
		return nil, fmt.Errorf("failed to submit deferred request: %w", err)
	}

	if response == nil {
		return nil, fmt.Errorf("received nil response")
	}

	// Create a ticker for polling
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	startTime := time.Now()

	for {
		select {
		case <-timeoutCtx.Done():
			return &PollResult{
				Response: response,
				Done:     false,
			}, timeoutCtx.Err()

		case <-ticker.C:
			// Check if the request is complete
			if response.Status() == "completed" {
				return &PollResult{
					Response: response,
					Done:     true,
				}, nil
			}

			// Check if we've timed out
			if time.Since(startTime) > timeout {
				return &PollResult{
					Response: response,
					Done:     false,
				}, fmt.Errorf("polling timeout")
			}

		case <-ctx.Done():
			return &PollResult{
				Response: response,
				Done:     false,
			}, ctx.Err()
		}
	}
}

// GetStoredCompletion retrieves a stored completion by ID.
func GetStoredCompletion(ctx context.Context, client ServiceClient, completionID string) (*StoredCompletion, error) {
	if client == nil {
		return nil, fmt.Errorf("chat client is nil")
	}
	if completionID == "" {
		return nil, fmt.Errorf("completion ID is required")
	}

	// TODO: Implement actual gRPC call to GetStoredCompletion
	// This requires implementing the ChatServiceClient.GetStoredCompletion method
	return nil, fmt.Errorf("GetStoredCompletion not yet implemented")
}

// DeleteStoredCompletion deletes a stored completion by ID.
func DeleteStoredCompletion(ctx context.Context, client ServiceClient, completionID string) error {
	if client == nil {
		return fmt.Errorf("chat client is nil")
	}
	if completionID == "" {
		return fmt.Errorf("completion ID is required")
	}

	// TODO: Implement actual gRPC call to DeleteStoredCompletion
	// This requires implementing the ChatServiceClient.DeleteStoredCompletion method
	return fmt.Errorf("DeleteStoredCompletion not yet implemented")
}

// StoredCompletion represents a stored chat completion.
type StoredCompletion struct {
	id        string
	content   string
	createdAt time.Time
}

// ID returns the stored completion ID.
func (sc *StoredCompletion) ID() string {
	return sc.id
}

// Content returns the stored completion content.
func (sc *StoredCompletion) Content() string {
	return sc.content
}

// CreatedAt returns the creation time of the stored completion.
func (sc *StoredCompletion) CreatedAt() time.Time {
	return sc.createdAt
}

// ListStoredCompletions retrieves a list of stored completions.
func ListStoredCompletions(ctx context.Context, client ServiceClient, opts ...ListOption) ([]*StoredCompletion, error) {
	if client == nil {
		return nil, fmt.Errorf("chat client is nil")
	}

	// TODO: Implement actual gRPC call to ListStoredCompletions
	// This requires implementing the ChatServiceClient.ListStoredCompletions method
	return nil, fmt.Errorf("ListStoredCompletions not yet implemented")
}

// ListOption represents an option for listing stored completions.
type ListOption func(*listConfig)

// listConfig represents configuration for listing stored completions.
type listConfig struct {
	limit int32
}

// WithLimit sets the limit for listing stored completions.
func WithLimit(limit int32) ListOption {
	return func(config *listConfig) {
		config.limit = limit
	}
}

// Validate validates the deferred request.
func (r *DeferredRequest) Validate() error {
	if r.proto == nil {
		return fmt.Errorf("request proto is nil")
	}

	if r.proto.Model == "" {
		return fmt.Errorf("model is required")
	}

	return nil
}
