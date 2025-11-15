// Package chat provides chat completion functionality for the xAI SDK.
package chat

import (
	"context"
	"fmt"

	xaiv1 "github.com/ZaguanLabs/xai-sdk-go/proto/gen/go/xai/v1"
	"google.golang.org/grpc"
)

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
	resp, err := client.CreateChatCompletion(ctx, r.proto)
	if err != nil {
		return nil, fmt.Errorf("chat completion failed: %w", err)
	}
	return &Response{proto: resp}, nil
}

// Stream performs a streaming chat completion request.
func (r *Request) Stream(ctx context.Context, client ChatServiceClient) (*Stream, error) {
	stream, err := client.StreamChatCompletion(ctx, r.proto)
	if err != nil {
		return nil, fmt.Errorf("chat completion stream failed: %w", err)
	}
	return &Stream{stream: stream}, nil
}

// Response methods

// Content returns the content of the response.
func (r *Response) Content() string {
	if r.proto == nil || len(r.proto.Choices) == 0 {
		return ""
	}
	return r.proto.Choices[0].Message.Content
}

// ToolCalls returns any tool calls in the response.
// Note: Tool functionality is not yet implemented in the proto definitions.
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
	return r.proto.Choices[0].Message.Role
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

// Proto returns the underlying protobuf usage.
func (u *TokenUsage) Proto() *xaiv1.Usage {
	return u.proto
}

// Next reads the next chunk from the stream.
func (s *Stream) Next() bool {
	if s.Err != nil {
		return false
	}

	chunk, err := s.stream.Recv()
	if err != nil {
		s.Err = err
		return false
	}

	s.current = &Chunk{proto: chunk}
	return true
}

// Current returns the current chunk.
func (s *Stream) Current() *Chunk {
	return s.current
}

// Close closes the stream.
func (s *Stream) Close() error {
	// For streaming RPCs, closing is handled by the client
	return nil
}
