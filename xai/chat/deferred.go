// Package chat provides deferred and stored chat functionality for xAI SDK.
package chat

import (
	"context"
	"fmt"
	"time"

	xaiv1 "github.com/ZaguanLabs/xai-sdk-go/proto/gen/go/xai/api/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// DeferredRequest represents a deferred chat completion request.
type DeferredRequest struct {
	proto *xaiv1.GetCompletionsRequest
}

type storedCompletionClient interface {
	GetStoredCompletion(ctx context.Context, in *xaiv1.GetStoredCompletionRequest, opts ...grpc.CallOption) (*xaiv1.GetChatCompletionResponse, error)
	DeleteStoredCompletion(ctx context.Context, in *xaiv1.DeleteStoredCompletionRequest, opts ...grpc.CallOption) (*xaiv1.DeleteStoredCompletionResponse, error)
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

	// Submit the deferred request
	start, err := client.StartDeferredCompletion(ctx, r.proto)
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

	if start == nil {
		return nil, fmt.Errorf("received nil response")
	}

	result, err := client.GetDeferredCompletion(ctx, &xaiv1.GetDeferredRequest{RequestId: start.RequestId})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch deferred response: %w", err)
	}
	if result == nil || result.Response == nil {
		return nil, fmt.Errorf("received nil deferred response")
	}

	return &DeferredResponse{proto: result.Response}, nil
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
	if interval <= 0 {
		interval = 100 * time.Millisecond
	}
	if timeout <= 0 {
		timeout = 10 * time.Minute
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	start, err := client.StartDeferredCompletion(timeoutCtx, r.proto)
	if err != nil {
		return nil, fmt.Errorf("failed to submit deferred request: %w", err)
	}
	if start == nil {
		return nil, fmt.Errorf("received nil deferred start response")
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-timeoutCtx.Done():
			return &PollResult{Done: false}, timeoutCtx.Err()
		case <-ticker.C:
			result, err := client.GetDeferredCompletion(timeoutCtx, &xaiv1.GetDeferredRequest{RequestId: start.RequestId})
			if err != nil {
				return &PollResult{Done: false}, err
			}
			if result == nil {
				return &PollResult{Done: false}, fmt.Errorf("received nil deferred completion response")
			}

			switch result.Status {
			case xaiv1.DeferredStatus_DONE:
				if result.Response == nil {
					return &PollResult{Done: false}, fmt.Errorf("deferred request completed without a response")
				}
				return &PollResult{
					Response: &DeferredResponse{proto: result.Response},
					Done:     true,
				}, nil
			case xaiv1.DeferredStatus_PENDING:
			case xaiv1.DeferredStatus_EXPIRED:
				return &PollResult{Done: false}, fmt.Errorf("deferred request expired")
			default:
				return &PollResult{Done: false}, fmt.Errorf("unknown deferred status: %s", result.Status.String())
			}
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

	storageClient, ok := client.(storedCompletionClient)
	if !ok {
		return nil, fmt.Errorf("chat client does not support stored completions")
	}

	response, err := storageClient.GetStoredCompletion(ctx, &xaiv1.GetStoredCompletionRequest{ResponseId: completionID})
	if err != nil {
		return nil, err
	}

	stored := &StoredCompletion{id: completionID}
	if response != nil {
		if response.Created != nil {
			stored.createdAt = response.Created.AsTime()
		}
		if len(response.Outputs) > 0 && response.Outputs[0] != nil && response.Outputs[0].Message != nil {
			stored.content = response.Outputs[0].Message.Content
		}
	}

	return stored, nil
}

// DeleteStoredCompletion deletes a stored completion by ID.
func DeleteStoredCompletion(ctx context.Context, client ServiceClient, completionID string) error {
	if client == nil {
		return fmt.Errorf("chat client is nil")
	}
	if completionID == "" {
		return fmt.Errorf("completion ID is required")
	}

	storageClient, ok := client.(storedCompletionClient)
	if !ok {
		return fmt.Errorf("chat client does not support stored completions")
	}

	_, err := storageClient.DeleteStoredCompletion(ctx, &xaiv1.DeleteStoredCompletionRequest{ResponseId: completionID})
	return err
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
