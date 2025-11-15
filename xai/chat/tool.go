// Package chat provides tool calling functionality for xAI SDK.
package chat

import (
	"encoding/json"
	"fmt"

	xaiv1 "github.com/ZaguanLabs/xai-sdk-go/proto/gen/go/xai/v1"
)

// Tool represents a function tool that can be called by the model.
type Tool struct {
	name        string
	description string
	parameters  map[string]interface{}
}

// NewTool creates a new tool with the given name and description.
func NewTool(name, description string) *Tool {
	return &Tool{
		name:        name,
		description: description,
		parameters:  make(map[string]interface{}),
	}
}

// WithParameter adds a parameter to the tool.
func (t *Tool) WithParameter(name string, paramType string, description string, required bool) *Tool {
	if t.parameters == nil {
		t.parameters = make(map[string]interface{})
	}

	t.parameters[name] = map[string]interface{}{
		"type":        paramType,
		"description": description,
		"required":     required,
	}
	return t
}

// WithParameters adds multiple parameters to the tool.
func (t *Tool) WithParameters(params map[string]interface{}) *Tool {
	if t.parameters == nil {
		t.parameters = make(map[string]interface{})
	}

	for name, param := range params {
		t.parameters[name] = param
	}
	return t
}

// Name returns the tool name.
func (t *Tool) Name() string {
	return t.name
}

// Description returns the tool description.
func (t *Tool) Description() string {
	return t.description
}

// Parameters returns the tool parameters.
func (t *Tool) Parameters() map[string]interface{} {
	if t.parameters == nil {
		return make(map[string]interface{})
	}
	return t.parameters
}

// ToJSONSchema converts the tool to JSON schema format.
func (t *Tool) ToJSONSchema() map[string]interface{} {
	properties := make(map[string]interface{})
	for name, param := range t.parameters {
		properties[name] = param
	}

	required := make([]string, 0)
	for name, param := range t.parameters {
		if paramMap, ok := param.(map[string]interface{}); ok {
			if req, exists := paramMap["required"]; exists && req.(bool) {
				required = append(required, name)
			}
		}
	}

	return map[string]interface{}{
		"type":        "object",
		"properties":   properties,
		"required":    required,
	}
}

// ToJSON converts the tool to JSON representation.
func (t *Tool) ToJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type":        "function",
		"function": map[string]interface{}{
			"name":        t.name,
			"description": t.description,
			"parameters":  t.ToJSONSchema(),
		},
	})
}

// Validate validates the tool definition.
func (t *Tool) Validate() error {
	if t.name == "" {
		return fmt.Errorf("tool name is required")
	}
	if t.description == "" {
		return fmt.Errorf("tool description is required")
	}
	if len(t.parameters) == 0 {
		return nil // Tools without parameters are valid
	}

	for name, param := range t.parameters {
		if name == "" {
			return fmt.Errorf("parameter name cannot be empty")
		}

		paramMap, ok := param.(map[string]interface{})
		if !ok {
			return fmt.Errorf("parameter '%s' must be a map", name)
		}

		// Validate required fields
		if paramType, exists := paramMap["type"]; !exists || paramType == "" {
			return fmt.Errorf("parameter '%s' must have a type", name)
		}
		if description, exists := paramMap["description"]; !exists || description == "" {
			return fmt.Errorf("parameter '%s' must have a description", name)
		}
	}

	return nil
}

// ToolChoice represents how tools should be chosen.
type ToolChoice string

const (
	// ToolChoiceNone means no tools will be used.
	ToolChoiceNone ToolChoice = "none"

	// ToolChoiceAuto means the model will decide whether to use tools.
	ToolChoiceAuto ToolChoice = "auto"

	// ToolChoiceRequired means tools must be used.
	ToolChoiceRequired ToolChoice = "required"

	// ToolChoiceSpecific means a specific tool must be used.
	ToolChoiceSpecific ToolChoice = "specific"
)

// NewToolChoiceNone creates a "none" tool choice.
func NewToolChoiceNone() *ToolChoiceOption {
	return &ToolChoiceOption{
		choice: ToolChoiceNone,
	}
}

// NewToolChoiceAuto creates an "auto" tool choice.
func NewToolChoiceAuto() *ToolChoiceOption {
	return &ToolChoiceOption{
		choice: ToolChoiceAuto,
	}
}

// NewToolChoiceRequired creates a "required" tool choice.
func NewToolChoiceRequired() *ToolChoiceOption {
	return &ToolChoiceOption{
		choice: ToolChoiceRequired,
	}
}

// NewToolChoiceSpecific creates a "specific" tool choice.
func NewToolChoiceSpecific(toolName string) *ToolChoiceOption {
	return &ToolChoiceOption{
		choice:  ToolChoiceSpecific,
		tool:   &toolName,
	}
}

// ToolChoiceOption represents a tool choice configuration.
type ToolChoiceOption struct {
	choice ToolChoice
	tool   *string // for ToolChoiceSpecific
}

// ToJSON converts the tool choice to JSON representation.
func (tco *ToolChoiceOption) ToJSON() map[string]interface{} {
	result := map[string]interface{}{
		"type": string(tco.choice),
	}

	if tco.choice == ToolChoiceSpecific && tco.tool != nil {
		result["function"] = map[string]interface{}{
			"name": *tco.tool,
		}
	}

	return result
}

// Validate validates the tool choice option.
func (tco *ToolChoiceOption) Validate() error {
	switch tco.choice {
	case ToolChoiceNone, ToolChoiceAuto, ToolChoiceRequired:
		return nil
	case ToolChoiceSpecific:
		if tco.tool == nil || *tco.tool == "" {
			return fmt.Errorf("tool name is required for specific tool choice")
		}
		return nil
	default:
		return fmt.Errorf("unknown tool choice: %s", tco.choice)
	}
}

// ToolCall represents a call to a tool.
type ToolCall struct {
	id       string
	name     string
	arguments map[string]interface{}
}

// NewToolCall creates a new tool call.
func NewToolCall(id, name string, arguments map[string]interface{}) *ToolCall {
	return &ToolCall{
		id:        id,
		name:      name,
		arguments: arguments,
	}
}

// ID returns the tool call ID.
func (tc *ToolCall) ID() string {
	return tc.id
}

// Name returns the tool call name.
func (tc *ToolCall) Name() string {
	return tc.name
}

// Arguments returns the tool call arguments.
func (tc *ToolCall) Arguments() map[string]interface{} {
	if tc.arguments == nil {
		return make(map[string]interface{})
	}
	return tc.arguments
}

// ToJSON converts the tool call to JSON representation.
func (tc *ToolCall) ToJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"id":       tc.id,
		"type":      "function",
		"function": map[string]interface{}{
			"name":      tc.name,
			"arguments": tc.arguments,
		},
	})
}

// Validate validates the tool call.
func (tc *ToolCall) Validate() error {
	if tc.id == "" {
		return fmt.Errorf("tool call ID is required")
	}
	if tc.name == "" {
		return fmt.Errorf("tool call name is required")
	}
	if tc.arguments == nil {
		return nil // Tool calls without arguments are valid
	}
	return nil
}

// ToolResult represents the result of a tool call.
type ToolResult struct {
	toolCallID string
	result     interface{}
	error      *string
}

// NewToolResult creates a new tool result.
func NewToolResult(toolCallID string, result interface{}) *ToolResult {
	return &ToolResult{
		toolCallID: toolCallID,
		result:     result,
	}
}

// NewToolResultError creates a new tool result with an error.
func NewToolResultError(toolCallID string, errorMsg string) *ToolResult {
	return &ToolResult{
		toolCallID: toolCallID,
		error:      &errorMsg,
	}
}

// ToolCallID returns the tool call ID this result is for.
func (tr *ToolResult) ToolCallID() string {
	return tr.toolCallID
}

// Result returns the result of the tool call.
func (tr *ToolResult) Result() interface{} {
	return tr.result
}

// Error returns the error from the tool call.
func (tr *ToolResult) Error() *string {
	return tr.error
}

// ToJSON converts the tool result to JSON representation.
func (tr *ToolResult) ToJSON() ([]byte, error) {
	if tr.error != nil {
		return json.Marshal(map[string]interface{}{
			"tool_call_id": tr.toolCallID,
			"error":        *tr.error,
		})
	}

	return json.Marshal(map[string]interface{}{
		"tool_call_id": tr.toolCallID,
		"result":       tr.result,
	})
}

// Validate validates the tool result.
func (tr *ToolResult) Validate() error {
	if tr.toolCallID == "" {
		return fmt.Errorf("tool call ID is required")
	}
	if tr.error != nil && tr.result != nil {
		return fmt.Errorf("tool result cannot have both error and result")
	}
	return nil
}