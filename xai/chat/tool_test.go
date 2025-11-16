package chat

import (
	"encoding/json"
	"testing"
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
	if protoTool.Function == nil {
		t.Fatal("Function should not be nil")
	}

	// Parse the parameters JSON
	var params map[string]interface{}
	if err := json.Unmarshal([]byte(protoTool.Function.Parameters), &params); err != nil {
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

	t.Logf("✅ WithTool generates valid JSON Schema: %s", protoTool.Function.Parameters)
}
