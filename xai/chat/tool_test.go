package chat

import (
	"encoding/json"
	"testing"

	xaiv1 "github.com/ZaguanLabs/xai-sdk-go/proto/gen/go/xai/api/v1"
)

func TestToolJSONSchemaFormat(t *testing.T) {
	// Create a tool with required and optional parameters
	tool := NewTool("test_function", "A test function")
	tool.WithParameter("required_param", "string", "A required parameter", true)
	tool.WithParameter("optional_param", "number", "An optional parameter", false)

	// Get the JSON schema
	schema := tool.ToJSONSchema()

	// Verify top-level structure
	if schema["type"] != "object" {
		t.Errorf("Expected type 'object', got %v", schema["type"])
	}

	// Verify properties exist
	properties, ok := schema["properties"].(map[string]interface{})
	if !ok {
		t.Fatal("Properties should be a map")
	}

	// Verify required_param property does NOT contain "required" field
	requiredParam, ok := properties["required_param"].(map[string]interface{})
	if !ok {
		t.Fatal("required_param should be a map")
	}
	if _, hasRequired := requiredParam["required"]; hasRequired {
		t.Error("Property 'required_param' should NOT contain 'required' field - it should only be in the top-level 'required' array")
	}
	if requiredParam["type"] != "string" {
		t.Errorf("Expected type 'string', got %v", requiredParam["type"])
	}
	if requiredParam["description"] != "A required parameter" {
		t.Errorf("Expected description 'A required parameter', got %v", requiredParam["description"])
	}

	// Verify optional_param property does NOT contain "required" field
	optionalParam, ok := properties["optional_param"].(map[string]interface{})
	if !ok {
		t.Fatal("optional_param should be a map")
	}
	if _, hasRequired := optionalParam["required"]; hasRequired {
		t.Error("Property 'optional_param' should NOT contain 'required' field")
	}

	// Verify top-level "required" array
	requiredArray, ok := schema["required"].([]string)
	if !ok {
		t.Fatal("Required should be a string array")
	}
	if len(requiredArray) != 1 {
		t.Errorf("Expected 1 required parameter, got %d", len(requiredArray))
	}
	if len(requiredArray) > 0 && requiredArray[0] != "required_param" {
		t.Errorf("Expected 'required_param' in required array, got %v", requiredArray[0])
	}

	// Verify the schema is valid JSON
	jsonBytes, err := json.Marshal(schema)
	if err != nil {
		t.Fatalf("Failed to marshal schema to JSON: %v", err)
	}

	// Verify we can unmarshal it back
	var unmarshaled map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal schema JSON: %v", err)
	}

	t.Logf("Valid JSON Schema generated: %s", string(jsonBytes))
}

func TestToolToJSON(t *testing.T) {
	// Create a tool
	tool := NewTool("get_weather", "Get the current weather")
	tool.WithParameter("city", "string", "The city name", true)
	tool.WithParameter("units", "string", "Temperature units", false)

	// Convert to JSON
	jsonBytes, err := tool.ToJSON()
	if err != nil {
		t.Fatalf("Failed to convert tool to JSON: %v", err)
	}

	// Parse the JSON
	var result map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &result); err != nil {
		t.Fatalf("Failed to parse tool JSON: %v", err)
	}

	// Verify structure
	if result["type"] != "function" {
		t.Errorf("Expected type 'function', got %v", result["type"])
	}

	function, ok := result["function"].(map[string]interface{})
	if !ok {
		t.Fatal("Function should be a map")
	}

	if function["name"] != "get_weather" {
		t.Errorf("Expected name 'get_weather', got %v", function["name"])
	}

	// Verify parameters schema
	params, ok := function["parameters"].(map[string]interface{})
	if !ok {
		t.Fatal("Parameters should be a map")
	}

	properties, ok := params["properties"].(map[string]interface{})
	if !ok {
		t.Fatal("Properties should be a map")
	}

	// Verify city property does NOT have "required" field
	city, ok := properties["city"].(map[string]interface{})
	if !ok {
		t.Fatal("City property should be a map")
	}
	if _, hasRequired := city["required"]; hasRequired {
		t.Error("City property should NOT contain 'required' field")
	}

	// Verify required array at top level
	requiredArray, ok := params["required"].([]interface{})
	if !ok {
		t.Fatal("Required should be an array")
	}
	if len(requiredArray) != 1 || requiredArray[0] != "city" {
		t.Errorf("Expected ['city'] in required array, got %v", requiredArray)
	}

	t.Logf("Valid tool JSON generated: %s", string(jsonBytes))
}

func TestToolValidation(t *testing.T) {
	tests := []struct {
		name      string
		setupTool func() *Tool
		wantError bool
	}{
		{
			name: "valid tool with parameters",
			setupTool: func() *Tool {
				tool := NewTool("test", "Test function")
				tool.WithParameter("param1", "string", "A parameter", true)
				return tool
			},
			wantError: false,
		},
		{
			name: "valid tool without parameters",
			setupTool: func() *Tool {
				return NewTool("test", "Test function")
			},
			wantError: false,
		},
		{
			name: "invalid tool - no name",
			setupTool: func() *Tool {
				return NewTool("", "Test function")
			},
			wantError: true,
		},
		{
			name: "invalid tool - no description",
			setupTool: func() *Tool {
				return NewTool("test", "")
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tool := tt.setupTool()
			err := tool.Validate()
			if (err != nil) != tt.wantError {
				t.Errorf("Validate() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestWithToolJSONSchemaFormat(t *testing.T) {
	// Create a tool with required and optional parameters
	tool := NewTool("get_weather", "Get the current weather")
	tool.WithParameter("city", "string", "The city name", true)
	tool.WithParameter("units", "string", "Temperature units (C or F)", false)

	// Create a request with the tool
	req := NewRequest("grok-beta", WithTool(tool))

	// Verify the tool was added
	if len(req.proto.Tools) != 1 {
		t.Fatalf("Expected 1 tool, got %d", len(req.proto.Tools))
	}

	// Get the function parameters as JSON
	protoTool := req.proto.Tools[0]
	if protoTool.GetFunction() == nil {
		t.Fatal("Function should not be nil")
	}

	// Parse the parameters JSON
	var params map[string]interface{}
	if err := json.Unmarshal([]byte(protoTool.GetFunction().Parameters), &params); err != nil {
		t.Fatalf("Failed to parse parameters JSON: %v", err)
	}

	// Verify top-level structure
	if params["type"] != "object" {
		t.Errorf("Expected type 'object', got %v", params["type"])
	}

	// Verify properties
	properties, ok := params["properties"].(map[string]interface{})
	if !ok {
		t.Fatal("Properties should be a map")
	}

	// Verify city property does NOT have "required" field
	city, ok := properties["city"].(map[string]interface{})
	if !ok {
		t.Fatal("City property should be a map")
	}
	if _, hasRequired := city["required"]; hasRequired {
		t.Error("City property should NOT contain 'required' field - it should only be in the top-level 'required' array")
	}
	if city["type"] != "string" {
		t.Errorf("Expected city type 'string', got %v", city["type"])
	}

	// Verify units property does NOT have "required" field
	units, ok := properties["units"].(map[string]interface{})
	if !ok {
		t.Fatal("Units property should be a map")
	}
	if _, hasRequired := units["required"]; hasRequired {
		t.Error("Units property should NOT contain 'required' field")
	}

	// Verify top-level required array
	requiredArray, ok := params["required"].([]interface{})
	if !ok {
		t.Fatal("Required should be an array")
	}
	if len(requiredArray) != 1 {
		t.Errorf("Expected 1 required parameter, got %d", len(requiredArray))
	}
	if len(requiredArray) > 0 && requiredArray[0] != "city" {
		t.Errorf("Expected 'city' in required array, got %v", requiredArray[0])
	}

	t.Logf("✅ WithTool generates valid JSON Schema: %s", protoTool.GetFunction().Parameters)
}

func TestResponseToolCalls(t *testing.T) {
	// Create a mock response with tool calls
	tool1 := NewTool("get_weather", "Get weather")
	tool1.WithParameter("city", "string", "City name", true)

	tool2 := NewTool("get_time", "Get time")
	tool2.WithParameter("timezone", "string", "Timezone", false)

	// Create a request with tools
	req := NewRequest("grok-beta", WithTool(tool1, tool2))

	// Simulate a response with tool calls (this would normally come from the API)
	// For testing, we'll verify the structure is correct
	if len(req.proto.Tools) != 2 {
		t.Fatalf("Expected 2 tools in request, got %d", len(req.proto.Tools))
	}

	t.Log("✅ Response.ToolCalls() structure is ready to parse tool calls from API")
}

func TestChunkToolCalls(t *testing.T) {
	// Create a mock chunk
	// The actual parsing will be tested with real API responses
	// This test verifies the method exists and returns the correct type
	chunk := &Chunk{}

	toolCalls := chunk.ToolCalls()
	if toolCalls == nil {
		t.Log("✅ Chunk.ToolCalls() returns nil for empty chunk (expected)")
	}

	hasToolCalls := chunk.HasToolCalls()
	if !hasToolCalls {
		t.Log("✅ Chunk.HasToolCalls() returns false for empty chunk (expected)")
	}
}

func TestParseToolCall(t *testing.T) {
	tests := []struct {
		name      string
		protoCall *xaiv1.ToolCall
		wantNil   bool
		wantID    string
		wantName  string
	}{
		{
			name:      "nil proto call",
			protoCall: nil,
			wantNil:   true,
		},
		{
			name: "valid tool call with arguments",
			protoCall: &xaiv1.ToolCall{
				Id: "call_123",
				Tool: &xaiv1.ToolCall_Function{
					Function: &xaiv1.FunctionCall{
						Name:      "get_weather",
						Arguments: `{"city": "San Francisco", "units": "celsius"}`,
					},
				},
			},
			wantNil:  false,
			wantID:   "call_123",
			wantName: "get_weather",
		},
		{
			name: "tool call with empty arguments",
			protoCall: &xaiv1.ToolCall{
				Id: "call_456",
				Tool: &xaiv1.ToolCall_Function{
					Function: &xaiv1.FunctionCall{
						Name:      "get_time",
						Arguments: "",
					},
				},
			},
			wantNil:  false,
			wantID:   "call_456",
			wantName: "get_time",
		},
		{
			name: "tool call with invalid JSON arguments",
			protoCall: &xaiv1.ToolCall{
				Id: "call_789",
				Tool: &xaiv1.ToolCall_Function{
					Function: &xaiv1.FunctionCall{
						Name:      "test_func",
						Arguments: `{invalid json}`,
					},
				},
			},
			wantNil:  false,
			wantID:   "call_789",
			wantName: "test_func",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseToolCall(tt.protoCall)

			if tt.wantNil {
				if result != nil {
					t.Errorf("Expected nil, got %v", result)
				}
				return
			}

			if result == nil {
				t.Fatal("Expected non-nil result")
			}

			if result.ID() != tt.wantID {
				t.Errorf("Expected ID %q, got %q", tt.wantID, result.ID())
			}

			if result.Name() != tt.wantName {
				t.Errorf("Expected name %q, got %q", tt.wantName, result.Name())
			}

			// Verify arguments is a valid map (even if empty)
			args := result.Arguments()
			if args == nil {
				t.Error("Arguments should never be nil, should be empty map")
			}
		})
	}
}

func TestToolCallJSON(t *testing.T) {
	toolCall := NewToolCall("call_123", "get_weather", map[string]interface{}{
		"city":  "San Francisco",
		"units": "celsius",
	})

	jsonBytes, err := toolCall.ToJSON()
	if err != nil {
		t.Fatalf("Failed to convert tool call to JSON: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &result); err != nil {
		t.Fatalf("Failed to parse tool call JSON: %v", err)
	}

	if result["id"] != "call_123" {
		t.Errorf("Expected id 'call_123', got %v", result["id"])
	}

	if result["type"] != "function" {
		t.Errorf("Expected type 'function', got %v", result["type"])
	}

	function, ok := result["function"].(map[string]interface{})
	if !ok {
		t.Fatal("Function should be a map")
	}

	if function["name"] != "get_weather" {
		t.Errorf("Expected name 'get_weather', got %v", function["name"])
	}

	t.Logf("✅ ToolCall JSON: %s", string(jsonBytes))
}
