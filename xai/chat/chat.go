// Package chat provides chat completion functionality for the xAI SDK.
package chat

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	xaiv1 "github.com/ZaguanLabs/xai-sdk-go/proto/gen/go/xai/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// OutputChunk represents a completion output chunk in streaming response.
type OutputChunk struct {
	proto *xaiv1.CompletionOutputChunk
}

// Index returns the index of the output chunk.
func (c *OutputChunk) Index() int32 {
	if c.proto == nil {
		return 0
	}
	return c.proto.Index
}

// Delta returns the message delta.
func (c *OutputChunk) Delta() *Delta {
	if c.proto == nil || c.proto.Delta == nil {
		return nil
	}
	return &Delta{proto: c.proto.Delta}
}

// FinishReason returns the finish reason.
func (c *OutputChunk) FinishReason() string {
	if c.proto == nil {
		return ""
	}
	return c.proto.FinishReason.String()
}

// Proto returns the underlying protobuf output chunk.
func (c *OutputChunk) Proto() *xaiv1.CompletionOutputChunk {
	return c.proto
}

// Delta represents a message delta in streaming response.
type Delta struct {
	proto *xaiv1.Delta
}

// Role returns the role from the delta.
func (d *Delta) Role() string {
	if d.proto == nil {
		return ""
	}
	return roleFromProto(d.proto.GetRole())
}

// Content returns the content from the delta.
func (d *Delta) Content() string {
	if d.proto == nil {
		return ""
	}
	return d.proto.GetContent()
}

// HasRole returns whether the delta has a role.
func (d *Delta) HasRole() bool {
	return d.proto != nil && d.proto.Role != xaiv1.MessageRole_INVALID_ROLE
}

// HasContent returns whether the delta has content.
func (d *Delta) HasContent() bool {
	return d.Content() != ""
}

// Proto returns the underlying protobuf delta.
func (d *Delta) Proto() *xaiv1.Delta {
	return d.proto
}

// Choice represents a single choice in the response (for compatibility).
type Choice struct {
	proto *xaiv1.CompletionOutput
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
	// Convert CompletionMessage (string content) to Message (Content array)
	contents := make([]*xaiv1.Content, 0)
	if c.proto.Message.Content != "" {
		contents = append(contents, &xaiv1.Content{Text: c.proto.Message.Content})
	}
	return &Message{
		proto: &xaiv1.Message{
			Role:    c.proto.Message.Role,
			Content: contents,
		},
	}
}

// FinishReason returns the finish reason of the choice.
func (c *Choice) FinishReason() string {
	if c.proto == nil {
		return ""
	}
	return c.proto.FinishReason.String()
}

// Proto returns the underlying protobuf output.
func (c *Choice) Proto() *xaiv1.CompletionOutput {
	return c.proto
}

// Request represents a chat completion request.
type Request struct {
	proto *xaiv1.GetCompletionsRequest
}

// Response represents a chat completion response.
type Response struct {
	proto *xaiv1.GetChatCompletionResponse
}

// Chunk represents a streaming response chunk.
type Chunk struct {
	proto *xaiv1.GetChatCompletionChunk
}

// RequestOption is a functional option for configuring a Request.
type RequestOption func(*Request)

// ChatServiceClient is an interface for the chat service client.
type ChatServiceClient = xaiv1.ChatClient

// Stream represents a streaming chat completion response.
type Stream struct {
	stream  xaiv1.Chat_GetCompletionChunkClient
	err     error
	current *Chunk
}

// NewRequest creates a new chat completion request.
func NewRequest(model string, opts ...RequestOption) *Request {
	req := &Request{
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

// WithTemperature sets the temperature for sampling.
func WithTemperature(temp float32) RequestOption {
	return func(r *Request) {
		r.proto.Temperature = temp
	}
}

// WithMaxTokens sets the maximum number of tokens.
func WithMaxTokens(maxTokens int32) RequestOption {
	return func(r *Request) {
		r.proto.MaxTokens = maxTokens
	}
}

// WithSearch adds search parameters to the request.
func WithSearch(params *SearchParameters) RequestOption {
	return func(r *Request) {
		r.proto.SearchParameters = params.Proto()
	}
}

// WithReasoningEffort adds reasoning effort to the request.
func WithReasoningEffort(effort ReasoningEffort) RequestOption {
	return func(r *Request) {
		r.proto.ReasoningEffort = reasoningEffortToProto(string(effort))
	}
}

// ReasoningEffort is the reasoning effort setting.
type ReasoningEffort string

const (
	// ReasoningEffortLow is the low reasoning effort.
	ReasoningEffortLow ReasoningEffort = "low"
	// ReasoningEffortMedium is the medium reasoning effort.
	ReasoningEffortMedium ReasoningEffort = "medium"
	// ReasoningEffortHigh is the high reasoning effort.
	ReasoningEffortHigh ReasoningEffort = "high"
)

// SearchParameters defines the search parameters.
type SearchParameters struct {
	pb *xaiv1.SearchParameters
}

// NewSearchParameters creates a new search parameter object.
func NewSearchParameters() *SearchParameters {
	return &SearchParameters{
		pb: &xaiv1.SearchParameters{},
	}
}

// WithMode sets the search mode.
func (p *SearchParameters) WithMode(mode string) *SearchParameters {
	p.pb.Mode = searchModeToProto(mode)
	return p
}

// WithReturnCitations sets whether to return citations.
func (p *SearchParameters) WithReturnCitations(returnCitations bool) *SearchParameters {
	p.pb.ReturnCitations = returnCitations
	return p
}

// WithMaxSearchResults sets the maximum number of search results.
func (p *SearchParameters) WithMaxSearchResults(maxResults int32) *SearchParameters {
	p.pb.MaxSearchResults = maxResults
	return p
}

// Proto returns the underlying protobuf message.
func (p *SearchParameters) Proto() *xaiv1.SearchParameters {
	if p == nil {
		return nil
	}
	return p.pb
}

// WithTool adds client-side function tools to the request.
func WithTool(tools ...*Tool) RequestOption {
	return func(r *Request) {
		// Convert tools to proto format
		protoTools := make([]*xaiv1.Tool, len(tools))
		for i, tool := range tools {
			// Convert parameters to proper JSON Schema format
			// ToJSONSchema() strips "required" from properties and builds top-level required array
			paramsJSON, _ := json.Marshal(tool.ToJSONSchema())
			protoTools[i] = &xaiv1.Tool{
				Function: &xaiv1.Function{
					Name:        tool.Name(),
					Description: tool.Description(),
					Parameters:  string(paramsJSON),
				},
			}
		}
		r.proto.Tools = append(r.proto.Tools, protoTools...)
	}
}

// WithServerTool adds server-side tools to the request.
// Server-side tools include web search, X search, code execution, etc.
func WithServerTool(tools ...*ServerTool) RequestOption {
	return func(r *Request) {
		for _, tool := range tools {
			if tool != nil && tool.Proto() != nil {
				r.proto.Tools = append(r.proto.Tools, tool.Proto())
			}
		}
	}
}

// ... (rest of the code remains the same)
// WithToolResults adds tool results to the request.
// Tool results are added as user messages with tool role.
func WithToolResults(results ...ToolResult) RequestOption {
	return func(r *Request) {
		for _, result := range results {
			var content string
			if result.Error() != nil {
				content = *result.Error()
			} else if str, ok := result.Result().(string); ok {
				content = str
			} else {
				// Convert to JSON string if not a string
				if jsonData, err := json.Marshal(result.Result()); err == nil {
					content = string(jsonData)
				} else {
					content = fmt.Sprintf("%v", result.Result())
				}
			}

			msg := &xaiv1.Message{
				Role:    xaiv1.MessageRole_ROLE_TOOL,
				Content: []*xaiv1.Content{{Text: content}},
				// Additional tool call info would be added here
			}
			r.proto.Messages = append(r.proto.Messages, msg)
		}
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

// AppendResponse appends a response as assistant message(s) to the request.
// This extracts the assistant's message including content, tool_calls,
// reasoning_content, and encrypted_content from the response.
// If the response has multiple outputs (N > 1), all outputs are appended.
func (r *Request) AppendResponse(resp *Response) *Request {
	if resp == nil || resp.proto == nil {
		return r
	}

	// Append all outputs from the response
	for _, output := range resp.proto.Outputs {
		if output == nil || output.Message == nil {
			continue
		}

		// Convert string content to Content array
		var contents []*xaiv1.Content
		if output.Message.Content != "" {
			contents = []*xaiv1.Content{
				{Text: output.Message.Content},
			}
		}

		// Create message with all fields from the response
		msg := &xaiv1.Message{
			Role:             output.Message.Role,
			Content:          contents,
			ToolCalls:        output.Message.ToolCalls,
			ReasoningContent: output.Message.ReasoningContent,
			EncryptedContent: output.Message.EncryptedContent,
		}

		r.proto.Messages = append(r.proto.Messages, msg)
	}

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
func WithMessages(messages ...*Message) RequestOption {
	return func(r *Request) {
		r.proto.Messages = make([]*xaiv1.Message, 0, len(messages))
		for _, msg := range messages {
			r.proto.Messages = append(r.proto.Messages, msg.Proto())
		}
	}
}

// WithMessage appends a message to the request as a functional option.
func WithMessage(msg *Message) RequestOption {
	return func(r *Request) {
		r.proto.Messages = append(r.proto.Messages, msg.Proto())
	}
}

// WithTopP sets the nucleus sampling parameter as a functional option.
func WithTopP(topP float32) RequestOption {
	return func(r *Request) {
		r.proto.TopP = topP
	}
}

// WithStop sets stop sequences as a functional option.
func WithStop(stop ...string) RequestOption {
	return func(r *Request) {
		r.proto.Stop = stop
	}
}

// WithFrequencyPenalty sets the frequency penalty as a functional option.
func WithFrequencyPenalty(penalty float32) RequestOption {
	return func(r *Request) {
		r.proto.FrequencyPenalty = penalty
	}
}

// WithPresencePenalty sets the presence penalty as a functional option.
func WithPresencePenalty(penalty float32) RequestOption {
	return func(r *Request) {
		r.proto.PresencePenalty = penalty
	}
}

// WithSeed sets the seed for deterministic sampling as a functional option.
func WithSeed(seed int32) RequestOption {
	return func(r *Request) {
		r.proto.Seed = seed
	}
}

// WithLogprobs enables returning log probabilities as a functional option.
func WithLogprobs(logprobs bool) RequestOption {
	return func(r *Request) {
		r.proto.Logprobs = logprobs
	}
}

// WithTopLogprobs sets the number of top log probabilities as a functional option.
func WithTopLogprobs(topLogprobs int32) RequestOption {
	return func(r *Request) {
		r.proto.TopLogprobs = topLogprobs
	}
}

// WithN sets the number of completions to generate as a functional option.
func WithN(n int32) RequestOption {
	return func(r *Request) {
		r.proto.N = n
	}
}

// WithUser sets the user identifier as a functional option.
func WithUser(user string) RequestOption {
	return func(r *Request) {
		r.proto.User = user
	}
}

// WithParallelToolCalls sets whether to execute tool calls in parallel as a functional option.
func WithParallelToolCalls(parallel bool) RequestOption {
	return func(r *Request) {
		r.proto.ParallelToolCalls = parallel
	}
}

// WithPreviousResponseID sets the previous response ID as a functional option.
func WithPreviousResponseID(responseID string) RequestOption {
	return func(r *Request) {
		r.proto.PreviousResponseId = responseID
	}
}

// WithStoreMessages sets whether to store messages as a functional option.
func WithStoreMessages(store bool) RequestOption {
	return func(r *Request) {
		r.proto.StoreMessages = store
	}
}

// WithUseEncryptedContent sets whether to use encrypted content as a functional option.
func WithUseEncryptedContent(encrypted bool) RequestOption {
	return func(r *Request) {
		r.proto.UseEncryptedContent = encrypted
	}
}

// SetTemperature sets the temperature for sampling.
func (r *Request) SetTemperature(temp float32) *Request {
	r.proto.Temperature = temp
	return r
}

// SetMaxTokens sets the maximum number of tokens.
func (r *Request) SetMaxTokens(maxTokens int32) *Request {
	r.proto.MaxTokens = maxTokens
	return r
}

// MaxTokens returns the max tokens setting.
func (r *Request) MaxTokens() int32 {
	return r.proto.MaxTokens
}

// SetTopP sets the nucleus sampling parameter (top_p).
// An alternative to sampling with temperature, where the model considers
// the results of the tokens with top_p probability mass.
// Values should be between 0 and 1.
func (r *Request) SetTopP(topP float32) *Request {
	r.proto.TopP = topP
	return r
}

// SetStop sets up to 4 sequences where the API will stop generating further tokens.
func (r *Request) SetStop(stop ...string) *Request {
	r.proto.Stop = stop
	return r
}

// SetFrequencyPenalty sets the frequency penalty parameter.
// Positive values penalize new tokens based on their existing frequency in the text so far,
// decreasing the model's likelihood to repeat the same line verbatim.
// Values typically range from 0 to 2.
func (r *Request) SetFrequencyPenalty(penalty float32) *Request {
	r.proto.FrequencyPenalty = penalty
	return r
}

// SetPresencePenalty sets the presence penalty parameter.
// Positive values penalize new tokens based on whether they appear in the text so far,
// increasing the model's likelihood to talk about new topics.
// Values typically range from 0 to 2.
func (r *Request) SetPresencePenalty(penalty float32) *Request {
	r.proto.PresencePenalty = penalty
	return r
}

// SetSeed sets the seed for deterministic sampling.
// If specified, the system will make a best effort to sample deterministically,
// such that repeated requests with the same seed and parameters should return the same result.
func (r *Request) SetSeed(seed int32) *Request {
	r.proto.Seed = seed
	return r
}

// SetLogprobs enables or disables returning log probabilities.
// If true, returns the log probabilities of each output token returned in the content of message.
func (r *Request) SetLogprobs(logprobs bool) *Request {
	r.proto.Logprobs = logprobs
	return r
}

// SetTopLogprobs sets the number of most likely tokens to return at each token position.
// An integer between 0 and 8 specifying the number of most likely tokens to return.
// logprobs must be set to true if this parameter is used.
func (r *Request) SetTopLogprobs(topLogprobs int32) *Request {
	r.proto.TopLogprobs = topLogprobs
	return r
}

// SetN sets the number of completions to generate.
// How many chat completion choices to generate for each input message.
// Note: Because this parameter generates many completions, it can quickly consume your token quota.
func (r *Request) SetN(n int32) *Request {
	r.proto.N = n
	return r
}

// SetUser sets a unique identifier representing your end-user.
// This can help xAI to monitor and detect abuse.
func (r *Request) SetUser(user string) *Request {
	r.proto.User = user
	return r
}

// SetParallelToolCalls sets whether to execute tool calls in parallel.
// If true, the model can call multiple tools simultaneously.
// If false, tool calls are executed sequentially.
func (r *Request) SetParallelToolCalls(parallel bool) *Request {
	r.proto.ParallelToolCalls = parallel
	return r
}

// SetPreviousResponseID sets the ID of a previous response for conversation continuity.
// This allows the model to reference and build upon previous interactions.
func (r *Request) SetPreviousResponseID(responseID string) *Request {
	r.proto.PreviousResponseId = responseID
	return r
}

// SetStoreMessages sets whether to store messages from this conversation.
// If true, messages will be stored for future reference.
func (r *Request) SetStoreMessages(store bool) *Request {
	r.proto.StoreMessages = store
	return r
}

// SetUseEncryptedContent sets whether to use encrypted content.
// If true, content will be encrypted for enhanced security.
func (r *Request) SetUseEncryptedContent(encrypted bool) *Request {
	r.proto.UseEncryptedContent = encrypted
	return r
}

// SetTools sets the tools for function calling.
// This method is deprecated. Use WithTool() or WithServerTool() instead.
func (r *Request) SetTools(tools ...Tool) *Request {
	// This method is kept for backwards compatibility but does nothing.
	// Use WithTool() for client-side function tools or WithServerTool() for server-side tools.
	// 		r.proto.Tools = append(r.proto.Tools, tool.Proto())
	// 	}
	// }
	return r
}

// SetToolChoice sets the tool choice for the request.
func (r *Request) SetToolChoice(choice ToolChoice) *Request {
	r.proto.ToolChoice = &xaiv1.ToolChoice{
		Mode: toolModeToProto(string(choice)),
	}
	return r
}

// Protos returns the underlying protobuf request.
func (r *Request) Proto() *xaiv1.GetCompletionsRequest {
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

	stream, err := client.GetCompletionChunk(ctx, r.proto)
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
	if r.proto.Temperature != 0 {
		temp := r.proto.Temperature
		if temp < 0.0 || temp > 2.0 {
			return fmt.Errorf("temperature must be between 0.0 and 2.0, got %f", temp)
		}
	}

	// Validate max_tokens if set
	if r.proto.MaxTokens != 0 {
		maxTokens := r.proto.MaxTokens
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
		if msg.Role == xaiv1.MessageRole_INVALID_ROLE {
			return fmt.Errorf("message at index %d has invalid role", i)
		}
		if len(msg.Content) == 0 {
			return fmt.Errorf("message at index %d has empty content", i)
		}

		// Validate role
		validRoles := map[xaiv1.MessageRole]bool{
			xaiv1.MessageRole_ROLE_SYSTEM:    true,
			xaiv1.MessageRole_ROLE_USER:      true,
			xaiv1.MessageRole_ROLE_ASSISTANT: true,
			xaiv1.MessageRole_ROLE_TOOL:      true,
		}
		if !validRoles[msg.Role] {
			return fmt.Errorf("invalid role '%s' in message at index %d", roleFromProto(msg.Role), i)
		}
	}

	return nil
}

// WithToolChoice adds tool choice to the request as a functional option.
func WithToolChoice(choice ToolChoice) RequestOption {
	return func(r *Request) {
		r.SetToolChoice(choice)
	}
}

// ToolChoice is the tool choice setting.
type ToolChoice string

const (
	// ToolChoiceAuto is the auto tool choice.
	ToolChoiceAuto ToolChoice = "auto"
	// ToolChoiceNone is the none tool choice.
	ToolChoiceNone ToolChoice = "none"
	// ToolChoiceRequired is the required tool choice.
	ToolChoiceRequired ToolChoice = "required"
)

// WithResponseFormat adds response format to the request.
func WithResponseFormat(format ResponseFormat) RequestOption {
	return func(r *Request) {
		// Convert to proto format
		switch format {
		case ResponseFormatText:
			r.proto.ResponseFormat = &xaiv1.ResponseFormat{
				FormatType: xaiv1.FormatType_FORMAT_TYPE_TEXT,
			}
		case ResponseFormatJSONObject:
			r.proto.ResponseFormat = &xaiv1.ResponseFormat{
				FormatType: xaiv1.FormatType_FORMAT_TYPE_JSON_OBJECT,
			}
		case ResponseFormatJSONSchema:
			// For JSON schema, we need to handle the ResponseFormatOption
			// This is a simplified implementation
			r.proto.ResponseFormat = &xaiv1.ResponseFormat{
				FormatType: xaiv1.FormatType_FORMAT_TYPE_JSON_SCHEMA,
			}
		}
	}
}

// WithResponseFormatOption adds response format with schema to the request.
func WithResponseFormatOption(option *ResponseFormatOption) RequestOption {
	return func(r *Request) {
		// Convert to proto format
		switch option.Type {
		case ResponseFormatText:
			r.proto.ResponseFormat = &xaiv1.ResponseFormat{
				FormatType: xaiv1.FormatType_FORMAT_TYPE_TEXT,
			}
		case ResponseFormatJSONObject:
			r.proto.ResponseFormat = &xaiv1.ResponseFormat{
				FormatType: xaiv1.FormatType_FORMAT_TYPE_JSON_OBJECT,
			}
		case ResponseFormatJSONSchema:
			// Convert schema to JSON string
			var schemaStr string
			if option.Schema != nil {
				if schemaData, err := json.Marshal(option.Schema); err == nil {
					schemaStr = string(schemaData)
				}
			}
			r.proto.ResponseFormat = &xaiv1.ResponseFormat{
				FormatType: xaiv1.FormatType_FORMAT_TYPE_JSON_SCHEMA,
				Schema:     schemaStr,
			}
		}
	}
}

// Response methods

// Content returns the content of the response.
func (r *Response) Content() string {
	if r.proto == nil || len(r.proto.Outputs) == 0 {
		return ""
	}
	if r.proto.Outputs[0] == nil || r.proto.Outputs[0].Message == nil {
		return ""
	}
	return r.proto.Outputs[0].Message.Content
}

// parseToolCall converts a proto ToolCall to a chat.ToolCall.
func parseToolCall(protoCall *xaiv1.ToolCall) *ToolCall {
	if protoCall == nil {
		return nil
	}

	// Parse function arguments from JSON string
	var arguments map[string]interface{}
	if protoCall.Function != nil && protoCall.Function.Arguments != "" {
		if err := json.Unmarshal([]byte(protoCall.Function.Arguments), &arguments); err != nil {
			// If parsing fails, create empty arguments map
			arguments = make(map[string]interface{})
		}
	} else {
		arguments = make(map[string]interface{})
	}

	// Get function name
	name := ""
	if protoCall.Function != nil {
		name = protoCall.Function.Name
	}

	return &ToolCall{
		id:        protoCall.Id,
		name:      name,
		arguments: arguments,
	}
}

// ToolCalls returns any tool calls in the response.
func (r *Response) ToolCalls() []*ToolCall {
	if r.proto == nil || len(r.proto.Outputs) == 0 {
		return nil
	}

	var toolCalls []*ToolCall
	for _, output := range r.proto.Outputs {
		if output.Message == nil {
			continue
		}
		// Only include tool calls from assistant messages
		if output.Message.Role == xaiv1.MessageRole_ROLE_ASSISTANT {
			for _, protoCall := range output.Message.ToolCalls {
				toolCall := parseToolCall(protoCall)
				if toolCall != nil {
					toolCalls = append(toolCalls, toolCall)
				}
			}
		}
	}

	return toolCalls
}

// Role returns the role of the response message.
func (r *Response) Role() string {
	if r.proto == nil || len(r.proto.Outputs) == 0 {
		return ""
	}
	if r.proto.Outputs[0] == nil || r.proto.Outputs[0].Message == nil {
		return ""
	}
	return roleFromProto(r.proto.Outputs[0].Message.Role)
}

// ReasoningContent returns the reasoning content from the response.
func (r *Response) ReasoningContent() string {
	if r.proto == nil || len(r.proto.Outputs) == 0 {
		return ""
	}
	if r.proto.Outputs[0] == nil || r.proto.Outputs[0].Message == nil {
		return ""
	}
	return r.proto.Outputs[0].Message.ReasoningContent
}

// EncryptedContent returns the encrypted content from the response.
func (r *Response) EncryptedContent() string {
	if r.proto == nil || len(r.proto.Outputs) == 0 {
		return ""
	}
	if r.proto.Outputs[0] == nil || r.proto.Outputs[0].Message == nil {
		return ""
	}
	return r.proto.Outputs[0].Message.EncryptedContent
}

// FinishReason returns the finish reason of the response.
func (r *Response) FinishReason() string {
	if r.proto == nil || len(r.proto.Outputs) == 0 {
		return ""
	}
	if r.proto.Outputs[0] == nil {
		return ""
	}
	return r.proto.Outputs[0].FinishReason.String()
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
	return len(r.proto.Outputs)
}

// Choice returns the choice at the given index.
func (r *Response) Choice(index int) *Choice {
	if r.proto == nil || index < 0 || index >= len(r.proto.Outputs) {
		return nil
	}
	if r.proto.Outputs[index] == nil {
		return nil
	}
	// Note: Choice type needs to be updated to work with CompletionOutput
	return nil // Temporarily return nil until Choice type is updated
}

// Choices returns all choices in the response.
func (r *Response) Choices() []*Choice {
	if r.proto == nil {
		return nil
	}
	// Note: Choice type needs to be updated to work with CompletionOutput
	return nil // Temporarily return nil until Choice type is updated
}

// Proto returns the underlying protobuf response.
func (r *Response) Proto() *xaiv1.GetChatCompletionResponse {
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
	if c.proto == nil || len(c.proto.Outputs) == 0 {
		return ""
	}
	// For streaming, return the delta content
	if c.proto.Outputs[0].Delta != nil {
		return c.proto.Outputs[0].Delta.GetContent()
	}
	return ""
}

// ToolCalls returns any tool calls in the chunk.
func (c *Chunk) ToolCalls() []*ToolCall {
	if c.proto == nil || len(c.proto.Outputs) == 0 {
		return nil
	}

	var toolCalls []*ToolCall
	for _, output := range c.proto.Outputs {
		if output.Delta == nil {
			continue
		}
		for _, protoCall := range output.Delta.ToolCalls {
			toolCall := parseToolCall(protoCall)
			if toolCall != nil {
				toolCalls = append(toolCalls, toolCall)
			}
		}
	}

	return toolCalls
}

// HasToolCalls returns whether the chunk has tool calls.
func (c *Chunk) HasToolCalls() bool {
	return len(c.ToolCalls()) > 0
}

// Role returns the role of the chunk message.
func (c *Chunk) Role() string {
	if c.proto == nil || len(c.proto.Outputs) == 0 {
		return ""
	}

	// For streaming, check delta first
	if c.proto.Outputs[0].Delta != nil {
		return roleFromProto(c.proto.Outputs[0].Delta.GetRole())
	}
	return ""
}

// ReasoningContent returns the reasoning content from the chunk.
func (c *Chunk) ReasoningContent() string {
	if c.proto == nil || len(c.proto.Outputs) == 0 {
		return ""
	}
	if c.proto.Outputs[0] == nil || c.proto.Outputs[0].Delta == nil {
		return ""
	}
	return c.proto.Outputs[0].Delta.ReasoningContent
}

// EncryptedContent returns the encrypted content from the chunk.
func (c *Chunk) EncryptedContent() string {
	if c.proto == nil || len(c.proto.Outputs) == 0 {
		return ""
	}
	if c.proto.Outputs[0] == nil || c.proto.Outputs[0].Delta == nil {
		return ""
	}
	return c.proto.Outputs[0].Delta.EncryptedContent
}

// Proto returns the underlying protobuf chunk.
func (c *Chunk) Proto() *xaiv1.GetChatCompletionChunk {
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
	proto *xaiv1.SamplingUsage
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

// validateStream validates the stream state.
func (s *Stream) validateStream() error {
	if s == nil {
		return fmt.Errorf("stream is nil")
	}
	if s.stream == nil {
		return fmt.Errorf("underlying stream is nil")
	}
	if s.err != nil {
		return fmt.Errorf("stream has error: %w", s.err)
	}
	return nil
}

// Next reads the next chunk from the stream.
func (s *Stream) Next() bool {
	// Validate stream state
	if err := s.validateStream(); err != nil {
		s.err = err
		return false
	}

	// Check if we already have an error
	if s.err != nil {
		return false
	}

	// Receive next chunk
	chunk, err := s.stream.Recv()
	if err != nil {
		if err == io.EOF {
			// Normal stream termination
			s.err = io.EOF
			return false
		}

		// Handle gRPC stream errors
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.Unauthenticated:
				s.err = fmt.Errorf("authentication failed: %s", st.Message())
			case codes.PermissionDenied:
				s.err = fmt.Errorf("permission denied: %s", st.Message())
			case codes.InvalidArgument:
				s.err = fmt.Errorf("invalid request: %s", st.Message())
			case codes.NotFound:
				s.err = fmt.Errorf("model not found: %s", st.Message())
			case codes.ResourceExhausted:
				s.err = fmt.Errorf("quota exceeded: %s", st.Message())
			case codes.Unavailable:
				s.err = fmt.Errorf("service unavailable: %s", st.Message())
			case codes.DeadlineExceeded:
				s.err = fmt.Errorf("request timeout: %s", st.Message())
			default:
				s.err = fmt.Errorf("stream error (%s): %s", st.Code().String(), st.Message())
			}
		} else {
			s.err = fmt.Errorf("stream failed: %w", err)
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
func (u *TokenUsage) Proto() *xaiv1.SamplingUsage {
	return u.proto
}

// Err returns the error that occurred during streaming.
func (s *Stream) Err() error {
	return s.err
}
