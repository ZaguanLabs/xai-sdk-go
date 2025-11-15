// Package chat provides structured output parsing functionality for xAI SDK.
package chat

import (
	"context"
	"encoding/json"
	"fmt"

	xaiv1 "github.com/ZaguanLabs/xai-sdk-go/proto/gen/go/xai/v1"
	"google.golang.org/grpc"
)

// ResponseFormat represents the desired response format.
type ResponseFormat string

const (
	// ResponseFormatText indicates plain text response.
	ResponseFormatText ResponseFormat = "text"
	
	// ResponseFormatJSONObject indicates JSON object response.
	ResponseFormatJSONObject ResponseFormat = "json_object"
	
	// ResponseFormatJSONSchema indicates JSON schema response.
	ResponseFormatJSONSchema ResponseFormat = "json_schema"
)

// Parse performs a chat completion request and parses the response into the provided type.
func (r *Request) Parse(ctx context.Context, client ChatServiceClient, v any) error {
	if client == nil {
		return fmt.Errorf("chat client is nil")
	}
	if r.proto == nil {
		return fmt.Errorf("request proto is nil")
	}
	if r.proto.Model == "" {
		return fmt.Errorf("model is required")
	}

	// Perform chat completion
	resp, err := client.CreateChatCompletion(ctx, r.proto)
	if err != nil {
		return fmt.Errorf("chat completion failed: %w", err)
	}

	if resp == nil || len(resp.Choices) == 0 {
		return fmt.Errorf("empty response received")
	}

	// Get the content from the first choice
	content := resp.Choices[0].Message.Content
	if content == "" {
		return fmt.Errorf("empty content in response")
	}

	// Parse based on the target type
	switch target := v.(type) {
	case *string:
		*target = content
		return nil
	case **string:
		*target = content
		return nil
	case map[string]interface{}:
		// Try to parse as JSON
		if err := json.Unmarshal([]byte(content), target); err != nil {
			return fmt.Errorf("failed to parse JSON response: %w", err)
		}
		return nil
	case *map[string]interface{}:
		// Try to parse as JSON
		if err := json.Unmarshal([]byte(content), *target); err != nil {
			return fmt.Errorf("failed to parse JSON response: %w", err)
		}
		return nil
	default:
		return fmt.Errorf("unsupported target type: %T", v)
	}
}

// ParseJSON performs a chat completion request and parses the response as JSON.
func (r *Request) ParseJSON(ctx context.Context, client ChatServiceClient, result interface{}) error {
	return r.Parse(ctx, client, result)
}

// ParseString performs a chat completion request and parses the response as a string.
func (r *Request) ParseString(ctx context.Context, client ChatServiceClient) (string, error) {
	var result string
	err := r.Parse(ctx, client, &result)
	if err != nil {
		return "", err
	}
	return result, nil
}

// WithResponseFormat sets the desired response format.
func WithResponseFormat(format ResponseFormat) RequestOption {
	return func(r *Request) {
		// Note: Response format is not yet implemented in proto definitions
		// r.proto.ResponseFormat = format
	}
}

// WithJSONSchema sets a JSON schema for structured output.
func WithJSONSchema(schema map[string]interface{}) RequestOption {
	return func(r *Request) {
		// Note: JSON schema is not yet implemented in proto definitions
		// r.proto.JSONSchema = schema
	}
}

// ResponseFormatOption represents a response format configuration.
type ResponseFormatOption struct {
	Type ResponseFormat `json:"type"`
	
	// Optional JSON schema for json_schema format
	Schema map[string]interface{} `json:"schema,omitempty"`
}

// NewResponseFormatOption creates a new response format option.
func NewResponseFormatOption(formatType ResponseFormat) *ResponseFormatOption {
	return &ResponseFormatOption{
		Type: formatType,
	}
}

// WithSchema sets the JSON schema for the response format option.
func (rfo *ResponseFormatOption) WithSchema(schema map[string]interface{}) *ResponseFormatOption {
	rfo.Schema = schema
	return rfo
}

// ToJSON converts the response format option to JSON.
func (rfo *ResponseFormatOption) ToJSON() ([]byte, error) {
	return json.Marshal(rfo)
}

// Validate validates the response format option.
func (rfo *ResponseFormatOption) Validate() error {
	switch rfo.Type {
	case ResponseFormatText, ResponseFormatJSONObject, ResponseFormatJSONSchema:
		// Valid formats
	default:
		return fmt.Errorf("unsupported response format: %s", rfo.Type)
	}
	
	// Validate schema if present for json_schema format
	if rfo.Type == ResponseFormatJSONSchema && rfo.Schema == nil {
		return fmt.Errorf("schema is required for json_schema format")
	}
	
	return nil
}