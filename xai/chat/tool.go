// Package chat provides tool calling functionality for xAI SDK.
package chat

import (
	"encoding/json"
	"fmt"
	"time"

	xaiv1 "github.com/ZaguanLabs/xai-sdk-go/proto/gen/go/xai/api/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Tool represents a function tool that can be called by the model.
type Tool struct {
	name        string
	description string
	parameters  map[string]interface{}
	strict      bool // Enable strict schema validation
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
		"required":    required,
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

// WithStrict enables strict schema validation for this tool.
// When enabled, the model will strictly validate function arguments against the schema.
func (t *Tool) WithStrict(strict bool) *Tool {
	t.strict = strict
	return t
}

// Strict returns whether strict schema validation is enabled.
func (t *Tool) Strict() bool {
	return t.strict
}

// ToJSONSchema converts the tool to JSON schema format.
func (t *Tool) ToJSONSchema() map[string]interface{} {
	properties := make(map[string]interface{})
	required := make([]string, 0)

	for name, param := range t.parameters {
		if paramMap, ok := param.(map[string]interface{}); ok {
			// Create a copy of the parameter without the "required" field
			// The "required" field should only be in the top-level "required" array
			propCopy := make(map[string]interface{})
			for key, value := range paramMap {
				if key != "required" {
					propCopy[key] = value
				} else if req, ok := value.(bool); ok && req {
					// Add to required array if true
					required = append(required, name)
				}
			}
			properties[name] = propCopy
		} else {
			// If it's not a map, just copy it as-is
			properties[name] = param
		}
	}

	return map[string]interface{}{
		"type":       "object",
		"properties": properties,
		"required":   required,
	}
}

// ToJSON converts the tool to JSON representation.
func (t *Tool) ToJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type": "function",
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

// ToolCall represents a call to a tool.
type ToolCall struct {
	id           string
	name         string
	arguments    map[string]interface{}
	status       string
	errorMessage string
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

// Status returns the tool call status.
func (tc *ToolCall) Status() string {
	return tc.status
}

// ErrorMessage returns the tool call error message.
func (tc *ToolCall) ErrorMessage() string {
	return tc.errorMessage
}

// Function returns a function representation of the tool call.
func (tc *ToolCall) Function() *ToolCallFunction {
	return &ToolCallFunction{
		name:      tc.name,
		arguments: tc.arguments,
	}
}

// ToolCallFunction represents the function part of a tool call.
type ToolCallFunction struct {
	name      string
	arguments map[string]interface{}
}

// Name returns the function name.
func (f *ToolCallFunction) Name() string {
	return f.name
}

// Arguments returns the function arguments.
func (f *ToolCallFunction) Arguments() map[string]interface{} {
	if f.arguments == nil {
		return make(map[string]interface{})
	}
	return f.arguments
}

// ToJSON converts the tool call to JSON representation.
func (tc *ToolCall) ToJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"id":   tc.id,
		"type": "function",
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

// ============================================================================
// Server-Side Tools
// ============================================================================

// ServerTool represents a server-side tool (web search, code execution, etc.)
type ServerTool struct {
	proto *xaiv1.Tool
}

// Proto returns the underlying proto representation.
func (st *ServerTool) Proto() *xaiv1.Tool {
	return st.proto
}

// WebSearchTool creates a web search server-side tool.
// This enables the model to search the web for information.
func WebSearchTool(opts ...WebSearchOption) *ServerTool {
	ws := &xaiv1.WebSearch{}
	for _, opt := range opts {
		opt(ws)
	}
	return &ServerTool{
		proto: &xaiv1.Tool{
			Tool: &xaiv1.Tool_WebSearch{
				WebSearch: ws,
			},
		},
	}
}

// WebSearchOption configures web search tool.
type WebSearchOption func(*xaiv1.WebSearch)

// WithExcludedDomains excludes specific domains from web search.
func WithExcludedDomains(domains ...string) WebSearchOption {
	return func(ws *xaiv1.WebSearch) {
		ws.ExcludedDomains = domains
	}
}

// WithAllowedDomains restricts web search to specific domains.
func WithAllowedDomains(domains ...string) WebSearchOption {
	return func(ws *xaiv1.WebSearch) {
		ws.AllowedDomains = domains
	}
}

// WithImageUnderstanding enables image understanding in web search results.
func WithImageUnderstanding(enable bool) WebSearchOption {
	return func(ws *xaiv1.WebSearch) {
		ws.EnableImageUnderstanding = &enable
	}
}

// XSearchTool creates an X (Twitter) search server-side tool.
// This enables the model to search X/Twitter for information.
func XSearchTool(opts ...XSearchOption) *ServerTool {
	xs := &xaiv1.XSearch{}
	for _, opt := range opts {
		opt(xs)
	}
	return &ServerTool{
		proto: &xaiv1.Tool{
			Tool: &xaiv1.Tool_XSearch{
				XSearch: xs,
			},
		},
	}
}

// XSearchOption configures X search tool.
type XSearchOption func(*xaiv1.XSearch)

// WithXDateRange sets the date range for X search.
func WithXDateRange(from, to time.Time) XSearchOption {
	return func(xs *xaiv1.XSearch) {
		if !from.IsZero() {
			xs.FromDate = timestamppb.New(from)
		}
		if !to.IsZero() {
			xs.ToDate = timestamppb.New(to)
		}
	}
}

// WithAllowedXHandles restricts X search to specific handles.
func WithAllowedXHandles(handles ...string) XSearchOption {
	return func(xs *xaiv1.XSearch) {
		xs.AllowedXHandles = handles
	}
}

// WithExcludedXHandles excludes specific X handles from search.
func WithExcludedXHandles(handles ...string) XSearchOption {
	return func(xs *xaiv1.XSearch) {
		xs.ExcludedXHandles = handles
	}
}

// WithXImageUnderstanding enables image understanding in X search results.
func WithXImageUnderstanding(enable bool) XSearchOption {
	return func(xs *xaiv1.XSearch) {
		xs.EnableImageUnderstanding = &enable
	}
}

// WithXVideoUnderstanding enables video understanding in X search results.
func WithXVideoUnderstanding(enable bool) XSearchOption {
	return func(xs *xaiv1.XSearch) {
		xs.EnableVideoUnderstanding = &enable
	}
}

// CodeExecutionTool creates a code execution server-side tool.
// This enables the model to execute code.
func CodeExecutionTool() *ServerTool {
	return &ServerTool{
		proto: &xaiv1.Tool{
			Tool: &xaiv1.Tool_CodeExecution{
				CodeExecution: &xaiv1.CodeExecution{},
			},
		},
	}
}

// CollectionsSearchTool creates a collections search server-side tool.
// This enables the model to search within document collections.
func CollectionsSearchTool(collectionIDs []string, opts ...CollectionsSearchOption) *ServerTool {
	cs := &xaiv1.CollectionsSearch{
		CollectionIds: collectionIDs,
	}
	for _, opt := range opts {
		opt(cs)
	}
	return &ServerTool{
		proto: &xaiv1.Tool{
			Tool: &xaiv1.Tool_CollectionsSearch{
				CollectionsSearch: cs,
			},
		},
	}
}

// CollectionsSearchOption configures collections search tool.
type CollectionsSearchOption func(*xaiv1.CollectionsSearch)

// WithCollectionsLimit sets the maximum number of results.
func WithCollectionsLimit(limit int32) CollectionsSearchOption {
	return func(cs *xaiv1.CollectionsSearch) {
		cs.Limit = &limit
	}
}

// DocumentSearchTool creates a document search server-side tool.
// This enables the model to search within uploaded documents.
func DocumentSearchTool(opts ...DocumentSearchOption) *ServerTool {
	ds := &xaiv1.DocumentSearch{}
	for _, opt := range opts {
		opt(ds)
	}
	return &ServerTool{
		proto: &xaiv1.Tool{
			Tool: &xaiv1.Tool_DocumentSearch{
				DocumentSearch: ds,
			},
		},
	}
}

// DocumentSearchOption configures document search tool.
type DocumentSearchOption func(*xaiv1.DocumentSearch)

// WithDocumentLimit sets the maximum number of document results.
func WithDocumentLimit(limit int32) DocumentSearchOption {
	return func(ds *xaiv1.DocumentSearch) {
		ds.Limit = &limit
	}
}

// MCPTool creates an MCP (Model Context Protocol) server-side tool.
// This enables the model to interact with MCP servers.
func MCPTool(serverLabel, serverURL string, opts ...MCPOption) *ServerTool {
	mcp := &xaiv1.MCP{
		ServerLabel: serverLabel,
		ServerUrl:   serverURL,
	}
	for _, opt := range opts {
		opt(mcp)
	}
	return &ServerTool{
		proto: &xaiv1.Tool{
			Tool: &xaiv1.Tool_Mcp{
				Mcp: mcp,
			},
		},
	}
}

// MCPOption configures MCP tool.
type MCPOption func(*xaiv1.MCP)

// WithMCPDescription sets the MCP server description.
func WithMCPDescription(description string) MCPOption {
	return func(mcp *xaiv1.MCP) {
		mcp.ServerDescription = description
	}
}

// WithMCPAllowedTools restricts which MCP tools can be called.
func WithMCPAllowedTools(toolNames ...string) MCPOption {
	return func(mcp *xaiv1.MCP) {
		mcp.AllowedToolNames = toolNames
	}
}

// WithMCPAuthorization sets the authorization header for MCP server.
func WithMCPAuthorization(auth string) MCPOption {
	return func(mcp *xaiv1.MCP) {
		mcp.Authorization = &auth
	}
}

// WithMCPExtraHeaders adds extra headers for MCP server requests.
func WithMCPExtraHeaders(headers map[string]string) MCPOption {
	return func(mcp *xaiv1.MCP) {
		if mcp.ExtraHeaders == nil {
			mcp.ExtraHeaders = make(map[string]string)
		}
		for key, value := range headers {
			mcp.ExtraHeaders[key] = value
		}
	}
}
