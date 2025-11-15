// Package chat provides chat completion functionality for the xAI SDK.
package chat

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	xaiv1 "github.com/ZaguanLabs/xai-sdk-go/proto/gen/go/xai/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ChoiceDelta represents a choice delta in streaming response.
type ChoiceDelta struct {
	proto *xaiv1.ChoiceDelta
}

// Index returns the index of the choice delta.
func (c *ChoiceDelta) Index() int32 {
	if c.proto == nil {
		return 0
	}
	return c.proto.Index
}

// Delta returns the message delta.
func (c *ChoiceDelta) Delta() *MessageDelta {
	if c.proto == nil || c.proto.Delta == nil {
		return nil
	}
	return &MessageDelta{proto: c.proto.Delta}
}

// Proto returns the underlying protobuf choice delta.
func (c *ChoiceDelta) Proto() *xaiv1.ChoiceDelta {
	return c.proto
}

// MessageDelta represents a message delta in streaming response.
type MessageDelta struct {
	proto *xaiv1.MessageDelta
}

// Role returns the role from the message delta.
func (m *MessageDelta) Role() string {
	if m.proto == nil {
		return ""
	}
	return m.proto.GetRole()
}

// Content returns the content from the message delta.
func (m *MessageDelta) Content() string {
	if m.proto == nil {
		return ""
	}
	return m.proto.GetContent()
}

// HasRole returns whether the message delta has a role.
func (m *MessageDelta) HasRole() bool {
	return m.Role() != ""
}

// HasContent returns whether the message delta has content.
func (m *MessageDelta) HasContent() bool {
	return m.Content() != ""
}

// Proto returns the underlying protobuf message delta.
func (m *MessageDelta) Proto() *xaiv1.MessageDelta {
	return m.proto
}

// Choice represents a single choice in the response.
type Choice struct {
	proto *xaiv1.Choice
}

// Index returns the index of the choice.
func (c *Choice) Index() int32 {
	if c.proto == nil {
		return 0
	}
	return c.proto.Index
}

// Message returns the message content of the choice.
func (c *Choice) Message() *Message {
	if c.proto == nil || c.proto.Message == nil {
		return nil
	}
	return &Message{proto: c.proto.Message}
}

// FinishReason returns the finish reason of the choice.
func (c *Choice) FinishReason() string {
	if c.proto == nil {
		return ""
	}
	return c.proto.FinishReason
}

// Proto returns the underlying protobuf choice.
func (c *Choice) Proto() *xaiv1.Choice {
	return c.proto
}

// Request represents a chat completion request.
type Request struct {
	proto *xaiv1.CreateChatCompletionRequest
}

// Response represents a chat completion response.
type Response struct {
	proto *xaiv1.CreateChatCompletionResponse
}

// Chunk represents a streaming response chunk.
type Chunk struct {
	proto *xaiv1.ChatCompletionChunk
}

// RequestOption is a functional option for configuring a Request.
type RequestOption func(*Request)

// ChatServiceClient is an interface for the chat service client.
type ChatServiceClient interface {
	CreateChatCompletion(ctx context.Context, req *xaiv1.CreateChatCompletionRequest, opts ...grpc.CallOption) (*xaiv1.CreateChatCompletionResponse, error)
	StreamChatCompletion(ctx context.Context, req *xaiv1.CreateChatCompletionRequest, opts ...grpc.CallOption) (xaiv1.Chat_StreamChatCompletionClient, error)
}

// Stream represents a streaming chat completion response.
type Stream struct {
	stream xaiv1.Chat_StreamChatCompletionClient
	Err    error
	current *Chunk
}

// NewRequest creates a new chat completion request.
func NewRequest(model string, opts ...RequestOption) *Request {
	req := &Request{
		proto: &xaiv1.CreateChatCompletionRequest{
			Model: model,
		},
	}

	// Apply functional options
	for _, opt := range opts {
		opt(req)
	}

	return req
}

// WithTemperature sets the temperature for sampling.
func WithTemperature(temp float32) RequestOption {
	return func(r *Request) {
		r.proto.Temperature = &temp
	}
}

// WithMaxTokens sets the maximum number of tokens.
func WithMaxTokens(maxTokens int32) RequestOption {
	return func(r *Request) {
		r.proto.MaxTokens = &maxTokens
	}
}

// WithSearch adds search parameters to the request.
func WithSearch(params *SearchParameters) RequestOption {
	return func(r *Request) {
		// Note: Search functionality is not yet implemented in proto definitions
		// r.proto.Search = params.ToJSON()
	}
}

// WithReasoningEffort adds reasoning effort to the request.
func WithReasoningEffort(effort *ReasoningEffortOption) RequestOption {
	return func(r *Request) {
		// Note: Reasoning functionality is not yet implemented in proto definitions
		// r.proto.ReasoningEffort = effort.ToJSON()
	}
}

// SetModel sets the model for the request.
func (r *Request) SetModel(model string) *Request {
	r.proto.Model = model
	return r
}

// GetModel returns the model for the request.
func (r *Request) GetModel() string {
	return r.proto.Model
}

// AppendMessage appends a message to the request.
func (r *Request) AppendMessage(msg Message) *Request {
	r.proto.Messages = append(r.proto.Messages, msg.Proto())
	return r
}

// SetMessages sets all messages for the request.
func (r *Request) SetMessages(messages ...Message) *Request {
	r.proto.Messages = make([]*xaiv1.Message, 0, len(messages))
	for _, msg := range messages {
		r.proto.Messages = append(r.proto.Messages, msg.Proto())
	}
	return r
}

// AddMessage adds a message to the request (alias for AppendMessage).
func (r *Request) AddMessage(msg Message) *Request {
	return r.AppendMessage(msg)
}

// WithMessages sets all messages for the request as a functional option.
func WithMessages(messages ...Message) RequestOption {
	return func(r *Request) {
		r.SetMessages(messages...)
	}
}

// WithMessage appends a message to the request as a functional option.
func WithMessage(msg Message) RequestOption {
	return func(r *Request) {
		r.AppendMessage(msg)
	}
}

// SetTemperature sets the temperature for sampling.
func (r *Request) SetTemperature(temp float32) *Request {
	r.proto.Temperature = &temp
	return r
}

// SetMaxTokens sets the maximum number of tokens.
func (r *Request) SetMaxTokens(maxTokens int32) *Request {
	r.proto.MaxTokens = &maxTokens
	return r
}

// MaxTokens returns the max tokens setting.
func (r *Request) MaxTokens() int32 {
	if r.proto.MaxTokens != nil {
		return *r.proto.MaxTokens
	}
	return 0
}

// SetTools sets the tools for function calling.
// Note: Tool functionality is not yet implemented in the proto definitions.
func (r *Request) SetTools(tools ...Tool) *Request {
	// Placeholder implementation until tools are properly defined in proto
	// r.proto.Tools = make([]*xaiv1.Tool, 0, len(tools))
	// for _, tool := range tools {
	// 	if tool.Proto() != nil {
	// 		r.proto.Tools = append(r.proto.Tools, tool.Proto())
	// 	}
	// }
	return r
}

// SetToolChoice sets how tools are chosen.
// Note: Tool functionality is not yet implemented in the proto definitions.
func (r *Request) SetToolChoice(choice ToolChoice) *Request {
	// Placeholder implementation until tool choice is properly defined in proto
	// r.proto.ToolChoice = choice.Proto()
	return r
}

// Protos returns the underlying protobuf request.
func (r *Request) Proto() *xaiv1.CreateChatCompletionRequest {
	return r.proto
}

// Sample performs a synchronous chat completion request.
func (r *Request) Sample(ctx context.Context, client ChatServiceClient) (*Response, error) {
	if client == nil {
		return nil, fmt.Errorf("chat client is nil")
	}
	if r.proto == nil {
		return nil, fmt.Errorf("request proto is nil")
	}
	if r.proto.Model == "" {
		return nil, fmt.Errorf("model is required")
	}

	// Validate request
	if err := r.validate(); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	resp, err := client.CreateChatCompletion(ctx, r.proto)
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
				return nil, fmt.Errorf("chat completion failed (%s): %s", st.Code().String(), st.Message())
			}
		}
		return nil, fmt.Errorf("chat completion failed: %w", err)
	}

	if resp == nil {
		return nil, fmt.Errorf("received nil response")
	}

	return &Response{proto: resp}, nil
}

// Stream performs a streaming chat completion request.
func (r *Request) Stream(ctx context.Context, client ChatServiceClient) (*Stream, error) {
	if client == nil {
		return nil, fmt.Errorf("chat client is nil")
	}
	if r.proto == nil {
		return nil, fmt.Errorf("request proto is nil")
	}
	if r.proto.Model == "" {
		return nil, fmt.Errorf("model is required")
	}

	// Validate request
	if err := r.validate(); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	stream, err := client.StreamChatCompletion(ctx, r.proto)
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
				return nil, fmt.Errorf("chat completion stream failed (%s): %s", st.Code().String(), st.Message())
			}
		}
		return nil, fmt.Errorf("chat completion stream failed: %w", err)
	}

	if stream == nil {
		return nil, fmt.Errorf("received nil stream")
	}

	return &Stream{stream: stream}, nil
}

// SetTools sets the tools for the request.
func (r *Request) SetTools(tools ...Tool) *Request {
	// Convert tools to JSON format
	toolSchemas := make([]map[string]interface{}, len(tools))
	for i, tool := range tools {
		toolSchemas[i] = map[string]interface{}{
			"type":     "function",
			"function": tool.ToJSONSchema(),
		}
	}

	// Note: This is a placeholder until tools are properly defined in proto
	// r.proto.Tools = toolSchemas
	return r
}

// WithTools adds tools to the request as a functional option.
func WithTools(tools ...Tool) RequestOption {
	return func(r *Request) {
		r.SetTools(tools...)
	}
}

// SetToolChoice sets how tools should be chosen.
func (r *Request) SetToolChoice(choice *ToolChoiceOption) *Request {
	// Note: This is a placeholder until tool choice is properly defined in proto
	// if choice != nil {
	// 	r.proto.ToolChoice = choice.ToJSON()
	// }
	return r
}

// WithToolChoice adds tool choice to the request as a functional option.
func WithToolChoice(choice *ToolChoiceOption) RequestOption {
	return func(r *Request) {
		r.SetToolChoice(choice)
	}
}

// Response methods

// Content returns the content of the response.
func (r *Response) Content() string {
	if r.proto == nil || len(r.proto.Choices) == 0 {
		return ""
	}
	if r.proto.Choices[0] == nil || r.proto.Choices[0].Message == nil {
		return ""
	}
	return r.proto.Choices[0].Message.Content
}

// ToolCalls returns any tool calls in the response.
func (r *Response) ToolCalls() []ToolCall {
	// Placeholder implementation until tool calls are properly defined in proto
	// if r.proto == nil || len(r.proto.Choices) == 0 {
	// 	return nil
	// }

	// calls := r.proto.Choices[0].Message.ToolCalls
	// if len(calls) == 0 {
	// 	return nil
	// }

	// result := make([]ToolCall, len(calls))
	// for i, call := range calls {
	// 	result[i] = ToolCall{proto: call}
	// }
	// return result
	return nil
}

// Role returns the role of the response message.
func (r *Response) Role() string {
	if r.proto == nil || len(r.proto.Choices) == 0 {
		return ""
	}
	if r.proto.Choices[0] == nil || r.proto.Choices[0].Message == nil {
		return ""
	}
	return r.proto.Choices[0].Message.Role
}

// FinishReason returns the finish reason of the response.
func (r *Response) FinishReason() string {
	if r.proto == nil || len(r.proto.Choices) == 0 {
		return ""
	}
	if r.proto.Choices[0] == nil {
		return ""
	}
	return r.proto.Choices[0].FinishReason
}

// ID returns the response ID.
func (r *Response) ID() string {
	if r.proto == nil {
		return ""
	}
	return r.proto.Id
}

// Model returns the model used for the response.
func (r *Response) Model() string {
	if r.proto == nil {
		return ""
	}
	return r.proto.Model
}

// ChoiceCount returns the number of choices in the response.
func (r *Response) ChoiceCount() int {
	if r.proto == nil {
		return 0
	}
	return len(r.proto.Choices)
}

// Choice returns the choice at the given index.
func (r *Response) Choice(index int) *Choice {
	if r.proto == nil || index < 0 || index >= len(r.proto.Choices) {
		return nil
	}
	if r.proto.Choices[index] == nil {
		return nil
	}
	return &Choice{proto: r.proto.Choices[index]}
}

// Choices returns all choices in the response.
func (r *Response) Choices() []*Choice {
	if r.proto == nil {
		return nil
	}
	choices := make([]*Choice, len(r.proto.Choices))
	for i, choice := range r.proto.Choices {
		choices[i] = &Choice{proto: choice}
	}
	return choices
}

// Proto returns the underlying protobuf response.
func (r *Response) Proto() *xaiv1.CreateChatCompletionResponse {
	return r.proto
}

// Usage returns the token usage information.
func (r *Response) Usage() *TokenUsage {
	if r.proto == nil || r.proto.Usage == nil {
		return nil
	}
	return &TokenUsage{proto: r.proto.Usage}
}

// Chunk methods

// Content returns the content of the chunk.
func (c *Chunk) Content() string {
	if c.proto == nil || len(c.proto.Choices) == 0 {
		return ""
	}
	// For streaming, return the delta content
	if c.proto.Choices[0].Delta != nil {
		return c.proto.Choices[0].Delta.GetContent()
	}
	return ""
}

// ToolCalls returns any tool calls in the chunk.
// Note: Tool functionality is not yet implemented in the proto definitions.
func (c *Chunk) ToolCalls() []ToolCall {
	// Placeholder implementation until tool calls are properly defined in proto
	return nil
}

// HasToolCalls returns whether the chunk has tool calls.
func (c *Chunk) HasToolCalls() bool {
	return len(c.ToolCalls()) > 0
}

// Role returns the role of the chunk message.
func (c *Chunk) Role() string {
	if c.proto == nil || len(c.proto.Choices) == 0 {
		return ""
	}

	// For streaming, check delta first
	if c.proto.Choices[0].Delta != nil {
		return c.proto.Choices[0].Delta.GetRole()
	}
	return ""
}

// Proto returns the underlying protobuf chunk.
func (c *Chunk) Proto() *xaiv1.ChatCompletionChunk {
	return c.proto
}

// FinishReason returns the finish reason for the chunk.
func (c *Chunk) FinishReason() string {
	// ChoiceDelta doesn't have a FinishReason field in the current proto
	return ""
}

// Usage returns the token usage information in the chunk.
func (c *Chunk) Usage() *TokenUsage {
	// Streaming chunks don't typically include usage information
	return nil
}

// TokenUsage represents token usage information.
type TokenUsage struct {
	proto *xaiv1.Usage
}

// PromptTokens returns the number of tokens in the prompt.
func (u *TokenUsage) PromptTokens() int32 {
	if u.proto == nil {
		return 0
	}
	return u.proto.GetPromptTokens()
}

// CompletionTokens returns the number of tokens in the completion.
func (u *TokenUsage) CompletionTokens() int32 {
	if u.proto == nil {
		return 0
	}
	return u.proto.GetCompletionTokens()
}

// TotalTokens returns the total number of tokens.
func (u *TokenUsage) TotalTokens() int32 {
	if u.proto == nil {
		return 0
	}
	return u.proto.GetTotalTokens()
}

// validate validates the request.
func (r *Request) validate() error {
	if r.proto == nil {
		return fmt.Errorf("request proto is nil")
	}

	// Validate model
	if r.proto.Model == "" {
		return fmt.Errorf("model is required")
	}

	// Validate temperature if set
	if r.proto.Temperature != nil {
		temp := *r.proto.Temperature
		if temp < 0.0 || temp > 2.0 {
			return fmt.Errorf("temperature must be between 0.0 and 2.0, got %f", temp)
		}
	}

	// Validate max_tokens if set
	if r.proto.MaxTokens != nil {
		maxTokens := *r.proto.MaxTokens
		if maxTokens < 1 || maxTokens > 8192 {
			return fmt.Errorf("max_tokens must be between 1 and 8192, got %d", maxTokens)
		}
	}

	// Validate messages
	if len(r.proto.Messages) == 0 {
		return fmt.Errorf("at least one message is required")
	}

	for i, msg := range r.proto.Messages {
		if msg == nil {
			return fmt.Errorf("message at index %d is nil", i)
		}
		if msg.Role == "" {
			return fmt.Errorf("message at index %d has empty role", i)
		}
		if msg.Content == "" {
			return fmt.Errorf("message at index %d has empty content", i)
		}

		// Validate role
		validRoles := map[string]bool{
			"system": true,
			"user": true,
			"assistant": true,
		}
		if !validRoles[msg.Role] {
			return fmt.Errorf("invalid role '%s' in message at index %d", msg.Role, i)
		}
	}

	return nil
}

// validateStream validates the stream state.
func (s *Stream) validateStream() error {
	if s == nil {
		return fmt.Errorf("stream is nil")
	}
	if s.stream == nil {
		return fmt.Errorf("underlying stream is nil")
	}
	if s.Err != nil {
		return fmt.Errorf("stream has error: %w", s.Err)
	}
	return nil
}

// Next reads the next chunk from the stream.
func (s *Stream) Next() bool {
	// Validate stream state
	if err := s.validateStream(); err != nil {
		s.Err = err
		return false
	}

	// Check if we already have an error
	if s.Err != nil {
		return false
	}

	// Receive next chunk
	chunk, err := s.stream.Recv()
	if err != nil {
		if err == io.EOF {
			// Normal stream termination
			s.Err = io.EOF
			return false
		}

		// Handle gRPC stream errors
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.Unauthenticated:
				s.Err = fmt.Errorf("authentication failed: %s", st.Message())
			case codes.PermissionDenied:
				s.Err = fmt.Errorf("permission denied: %s", st.Message())
			case codes.InvalidArgument:
				s.Err = fmt.Errorf("invalid request: %s", st.Message())
			case codes.NotFound:
				s.Err = fmt.Errorf("model not found: %s", st.Message())
			case codes.ResourceExhausted:
				s.Err = fmt.Errorf("quota exceeded: %s", st.Message())
			case codes.Unavailable:
				s.Err = fmt.Errorf("service unavailable: %s", st.Message())
			case codes.DeadlineExceeded:
				s.Err = fmt.Errorf("request timeout: %s", st.Message())
			default:
				s.Err = fmt.Errorf("stream error (%s): %s", st.Code().String(), st.Message())
			}
		} else {
			s.Err = fmt.Errorf("stream failed: %w", err)
		}
		return false
	}

	// Set current chunk
	s.current = &Chunk{proto: chunk}
	return true
}

// Current returns the current chunk.
func (s *Stream) Current() *Chunk {
	return s.current
}

// Close closes the stream and cleans up resources.
func (s *Stream) Close() error {
	if s == nil {
		return nil
	}

	// If we have a gRPC stream, close it
	if s.stream != nil {
		// Try to close the stream gracefully
		err := s.stream.CloseSend()
		if err != nil && err != io.EOF {
			return fmt.Errorf("failed to close stream: %w", err)
		}
	}

	return nil
}

// Proto returns the underlying protobuf usage.
func (u *TokenUsage) Proto() *xaiv1.Usage {
	return u.proto
}

// Err returns the error that occurred during streaming.
func (s *Stream) Err() error {
	return s.Err
}
