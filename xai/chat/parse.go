// Package chat provides structured output parsing functionality for xAI SDK.
package chat

import (
	"context"
	"encoding/json"
	"fmt"
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

// NewResponseFormatText creates a text response format.
func NewResponseFormatText() ResponseFormat {
	return ResponseFormatText
}

// NewResponseFormatJSONObject creates a JSON object response format.
func NewResponseFormatJSONObject() ResponseFormat {
	return ResponseFormatJSONObject
}

// NewResponseFormatJSONSchema creates a JSON schema response format with the given schema.
func NewResponseFormatJSONSchema(schema map[string]interface{}) *ResponseFormatOption {
	return &ResponseFormatOption{
		Type:   ResponseFormatJSONSchema,
		Schema: schema,
	}
}

func parseContent(content string, v any) error {
	switch target := v.(type) {
	case *string:
		*target = content
		return nil
	case **string:
		**target = content
		return nil
	case map[string]interface{}:
		// Try to parse as JSON - need to pass pointer to target
		targetPtr := &target
		if err := json.Unmarshal([]byte(content), targetPtr); err != nil {
			return fmt.Errorf("failed to parse JSON response: %w", err)
		}
		return nil
	case *map[string]interface{}:
		// Try to parse as JSON - target is already a pointer
		if err := json.Unmarshal([]byte(content), target); err != nil {
			return fmt.Errorf("failed to parse JSON response: %w", err)
		}
		return nil
	default:
		return fmt.Errorf("unsupported target type: %T", v)
	}
}

// ParseWithResponse performs a chat completion request, parses the response into the provided type,
// and returns the raw response alongside the parsed value.
func (r *Request) ParseWithResponse(ctx context.Context, client ServiceClient, v any) (*Response, error) {
	if client == nil {
		return nil, fmt.Errorf("chat client is nil")
	}
	if r.proto == nil {
		return nil, fmt.Errorf("request proto is nil")
	}
	if r.proto.Model == "" {
		return nil, fmt.Errorf("model is required")
	}

	resp, err := r.Sample(ctx, client)
	if err != nil {
		return nil, err
	}

	content := resp.Content()
	if content == "" {
		return nil, fmt.Errorf("empty content in response")
	}

	if err := parseContent(content, v); err != nil {
		return nil, err
	}

	return resp, nil
}

// Parse performs a chat completion request and parses the response into the provided type.
func (r *Request) Parse(ctx context.Context, client ServiceClient, v any) error {
	_, err := r.ParseWithResponse(ctx, client, v)
	return err
}

// ParseJSON performs a chat completion request and parses the response as JSON.
func (r *Request) ParseJSON(ctx context.Context, client ServiceClient, result interface{}) error {
	return r.Parse(ctx, client, result)
}

// ParseString performs a chat completion request and parses the response as a string.
func (r *Request) ParseString(ctx context.Context, client ServiceClient) (string, error) {
	var result string
	err := r.Parse(ctx, client, &result)
	if err != nil {
		return "", err
	}
	return result, nil
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
