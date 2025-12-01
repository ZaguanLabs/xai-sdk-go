package chat

import (
	"testing"

	xaiv1 "github.com/ZaguanLabs/xai-sdk-go/proto/gen/go/xai/api/v1"
)

func TestIncludeOptionConversion(t *testing.T) {
	tests := []struct {
		name     string
		option   IncludeOption
		expected xaiv1.IncludeOption
	}{
		{"WebSearchCallOutput", IncludeWebSearchCallOutput, xaiv1.IncludeOption_INCLUDE_OPTION_WEB_SEARCH_CALL_OUTPUT},
		{"XSearchCallOutput", IncludeXSearchCallOutput, xaiv1.IncludeOption_INCLUDE_OPTION_X_SEARCH_CALL_OUTPUT},
		{"CodeExecutionCallOutput", IncludeCodeExecutionCallOutput, xaiv1.IncludeOption_INCLUDE_OPTION_CODE_EXECUTION_CALL_OUTPUT},
		{"CollectionsSearchCallOutput", IncludeCollectionsSearchCallOutput, xaiv1.IncludeOption_INCLUDE_OPTION_COLLECTIONS_SEARCH_CALL_OUTPUT},
		{"DocumentSearchCallOutput", IncludeDocumentSearchCallOutput, xaiv1.IncludeOption_INCLUDE_OPTION_DOCUMENT_SEARCH_CALL_OUTPUT},
		{"MCPCallOutput", IncludeMCPCallOutput, xaiv1.IncludeOption_INCLUDE_OPTION_MCP_CALL_OUTPUT},
		{"InlineCitations", IncludeInlineCitations, xaiv1.IncludeOption_INCLUDE_OPTION_INLINE_CITATIONS},
		{"Invalid", IncludeOption("invalid"), xaiv1.IncludeOption_INCLUDE_OPTION_INVALID},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := includeOptionToProto(tt.option)
			if result != tt.expected {
				t.Errorf("includeOptionToProto(%s) = %v, want %v", tt.option, result, tt.expected)
			}
		})
	}
}

func TestIncludeOptionFromProtoConversion(t *testing.T) {
	tests := []struct {
		name     string
		proto    xaiv1.IncludeOption
		expected IncludeOption
	}{
		{"WebSearchCallOutput", xaiv1.IncludeOption_INCLUDE_OPTION_WEB_SEARCH_CALL_OUTPUT, IncludeWebSearchCallOutput},
		{"XSearchCallOutput", xaiv1.IncludeOption_INCLUDE_OPTION_X_SEARCH_CALL_OUTPUT, IncludeXSearchCallOutput},
		{"CodeExecutionCallOutput", xaiv1.IncludeOption_INCLUDE_OPTION_CODE_EXECUTION_CALL_OUTPUT, IncludeCodeExecutionCallOutput},
		{"CollectionsSearchCallOutput", xaiv1.IncludeOption_INCLUDE_OPTION_COLLECTIONS_SEARCH_CALL_OUTPUT, IncludeCollectionsSearchCallOutput},
		{"DocumentSearchCallOutput", xaiv1.IncludeOption_INCLUDE_OPTION_DOCUMENT_SEARCH_CALL_OUTPUT, IncludeDocumentSearchCallOutput},
		{"MCPCallOutput", xaiv1.IncludeOption_INCLUDE_OPTION_MCP_CALL_OUTPUT, IncludeMCPCallOutput},
		{"InlineCitations", xaiv1.IncludeOption_INCLUDE_OPTION_INLINE_CITATIONS, IncludeInlineCitations},
		{"Invalid", xaiv1.IncludeOption_INCLUDE_OPTION_INVALID, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := includeOptionFromProto(tt.proto)
			if result != tt.expected {
				t.Errorf("includeOptionFromProto(%v) = %s, want %s", tt.proto, result, tt.expected)
			}
		})
	}
}

func TestWithInclude(t *testing.T) {
	req := NewRequest("grok-3",
		WithInclude(IncludeInlineCitations, IncludeWebSearchCallOutput),
	)

	if len(req.proto.Include) != 2 {
		t.Errorf("Expected 2 include options, got %d", len(req.proto.Include))
	}

	if req.proto.Include[0] != xaiv1.IncludeOption_INCLUDE_OPTION_INLINE_CITATIONS {
		t.Errorf("Expected first include option to be INLINE_CITATIONS, got %v", req.proto.Include[0])
	}

	if req.proto.Include[1] != xaiv1.IncludeOption_INCLUDE_OPTION_WEB_SEARCH_CALL_OUTPUT {
		t.Errorf("Expected second include option to be WEB_SEARCH_CALL_OUTPUT, got %v", req.proto.Include[1])
	}
}

func TestSetInclude(t *testing.T) {
	req := NewRequest("grok-3")
	req.SetInclude(IncludeMCPCallOutput, IncludeCodeExecutionCallOutput)

	if len(req.proto.Include) != 2 {
		t.Errorf("Expected 2 include options, got %d", len(req.proto.Include))
	}

	if req.proto.Include[0] != xaiv1.IncludeOption_INCLUDE_OPTION_MCP_CALL_OUTPUT {
		t.Errorf("Expected first include option to be MCP_CALL_OUTPUT, got %v", req.proto.Include[0])
	}

	if req.proto.Include[1] != xaiv1.IncludeOption_INCLUDE_OPTION_CODE_EXECUTION_CALL_OUTPUT {
		t.Errorf("Expected second include option to be CODE_EXECUTION_CALL_OUTPUT, got %v", req.proto.Include[1])
	}
}

func TestInlineCitation(t *testing.T) {
	// Test with nil proto
	ic := &InlineCitation{proto: nil}
	if ic.ID() != "" {
		t.Error("Expected empty ID for nil proto")
	}
	if ic.StartIndex() != 0 {
		t.Error("Expected 0 StartIndex for nil proto")
	}
	if ic.WebCitation() != nil {
		t.Error("Expected nil WebCitation for nil proto")
	}
	if ic.XCitation() != nil {
		t.Error("Expected nil XCitation for nil proto")
	}
	if ic.CollectionsCitation() != nil {
		t.Error("Expected nil CollectionsCitation for nil proto")
	}

	// Test with web citation
	webCitation := &xaiv1.InlineCitation{
		Id:         "cite-1",
		StartIndex: 42,
		WebCitation: &xaiv1.WebCitation{
			Url: "https://example.com",
		},
	}
	ic = &InlineCitation{proto: webCitation}
	if ic.ID() != "cite-1" {
		t.Errorf("Expected ID 'cite-1', got '%s'", ic.ID())
	}
	if ic.StartIndex() != 42 {
		t.Errorf("Expected StartIndex 42, got %d", ic.StartIndex())
	}
	if ic.WebCitation() == nil {
		t.Error("Expected non-nil WebCitation")
	} else if ic.WebCitation().URL() != "https://example.com" {
		t.Errorf("Expected URL 'https://example.com', got '%s'", ic.WebCitation().URL())
	}

	// Test with X citation
	xCitation := &xaiv1.InlineCitation{
		Id:         "cite-2",
		StartIndex: 100,
		XCitation: &xaiv1.XCitation{
			Url: "https://x.com/user/status/123",
		},
	}
	ic = &InlineCitation{proto: xCitation}
	if ic.XCitation() == nil {
		t.Error("Expected non-nil XCitation")
	} else if ic.XCitation().URL() != "https://x.com/user/status/123" {
		t.Errorf("Expected URL 'https://x.com/user/status/123', got '%s'", ic.XCitation().URL())
	}

	// Test with collections citation
	collectionsCitation := &xaiv1.InlineCitation{
		Id:         "cite-3",
		StartIndex: 200,
		CollectionsCitation: &xaiv1.CollectionsCitation{
			FileId:        "file-123",
			ChunkId:       "chunk-456",
			ChunkContent:  "Some content from the document",
			Score:         0.95,
			CollectionIds: []string{"col-1", "col-2"},
		},
	}
	ic = &InlineCitation{proto: collectionsCitation}
	cc := ic.CollectionsCitation()
	if cc == nil {
		t.Error("Expected non-nil CollectionsCitation")
	} else {
		if cc.FileID() != "file-123" {
			t.Errorf("Expected FileID 'file-123', got '%s'", cc.FileID())
		}
		if cc.ChunkID() != "chunk-456" {
			t.Errorf("Expected ChunkID 'chunk-456', got '%s'", cc.ChunkID())
		}
		if cc.ChunkContent() != "Some content from the document" {
			t.Errorf("Expected ChunkContent 'Some content from the document', got '%s'", cc.ChunkContent())
		}
		if cc.Score() != 0.95 {
			t.Errorf("Expected Score 0.95, got %f", cc.Score())
		}
		if len(cc.CollectionIDs()) != 2 {
			t.Errorf("Expected 2 CollectionIDs, got %d", len(cc.CollectionIDs()))
		}
	}
}

func TestResponseInlineCitations(t *testing.T) {
	// Test with nil response
	resp := &Response{proto: nil}
	if resp.InlineCitations() != nil {
		t.Error("Expected nil InlineCitations for nil proto")
	}

	// Test with response containing inline citations
	resp = &Response{
		proto: &xaiv1.GetChatCompletionResponse{
			Outputs: []*xaiv1.CompletionOutput{
				{
					Message: &xaiv1.CompletionMessage{
						Content: "Here is some content with citations",
						Citations: []*xaiv1.InlineCitation{
							{
								Id:         "cite-1",
								StartIndex: 10,
								WebCitation: &xaiv1.WebCitation{
									Url: "https://example.com",
								},
							},
							{
								Id:         "cite-2",
								StartIndex: 25,
								XCitation: &xaiv1.XCitation{
									Url: "https://x.com/post",
								},
							},
						},
					},
				},
			},
		},
	}

	citations := resp.InlineCitations()
	if len(citations) != 2 {
		t.Errorf("Expected 2 citations, got %d", len(citations))
	}

	if citations[0].ID() != "cite-1" {
		t.Errorf("Expected first citation ID 'cite-1', got '%s'", citations[0].ID())
	}
	if citations[1].ID() != "cite-2" {
		t.Errorf("Expected second citation ID 'cite-2', got '%s'", citations[1].ID())
	}
}

func TestToolCallType(t *testing.T) {
	tests := []struct {
		name     string
		proto    xaiv1.ToolCallType
		expected ToolCallType
	}{
		{"ClientSide", xaiv1.ToolCallType_TOOL_CALL_TYPE_CLIENT_SIDE_TOOL, ToolCallTypeClientSide},
		{"WebSearch", xaiv1.ToolCallType_TOOL_CALL_TYPE_WEB_SEARCH_TOOL, ToolCallTypeWebSearch},
		{"XSearch", xaiv1.ToolCallType_TOOL_CALL_TYPE_X_SEARCH_TOOL, ToolCallTypeXSearch},
		{"CodeExecution", xaiv1.ToolCallType_TOOL_CALL_TYPE_CODE_EXECUTION_TOOL, ToolCallTypeCodeExecution},
		{"CollectionsSearch", xaiv1.ToolCallType_TOOL_CALL_TYPE_COLLECTIONS_SEARCH_TOOL, ToolCallTypeCollectionsSearch},
		{"MCP", xaiv1.ToolCallType_TOOL_CALL_TYPE_MCP_TOOL, ToolCallTypeMCP},
		{"DocumentSearch", xaiv1.ToolCallType_TOOL_CALL_TYPE_DOCUMENT_SEARCH_TOOL, ToolCallTypeDocumentSearch},
		{"Invalid", xaiv1.ToolCallType_TOOL_CALL_TYPE_INVALID, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := toolCallTypeFromProto(tt.proto)
			if result != tt.expected {
				t.Errorf("toolCallTypeFromProto(%v) = %s, want %s", tt.proto, result, tt.expected)
			}
		})
	}
}

func TestToolCallIsClientSide(t *testing.T) {
	// Client-side tool call
	tc := &ToolCall{toolType: ToolCallTypeClientSide}
	if !tc.IsClientSide() {
		t.Error("Expected IsClientSide() to return true for client-side tool")
	}
	if tc.IsServerSide() {
		t.Error("Expected IsServerSide() to return false for client-side tool")
	}

	// Empty tool type (defaults to client-side)
	tc = &ToolCall{toolType: ""}
	if !tc.IsClientSide() {
		t.Error("Expected IsClientSide() to return true for empty tool type")
	}
	if tc.IsServerSide() {
		t.Error("Expected IsServerSide() to return false for empty tool type")
	}

	// Server-side tool call
	tc = &ToolCall{toolType: ToolCallTypeWebSearch}
	if tc.IsClientSide() {
		t.Error("Expected IsClientSide() to return false for server-side tool")
	}
	if !tc.IsServerSide() {
		t.Error("Expected IsServerSide() to return true for server-side tool")
	}
}

func TestToolCallTypeAccessor(t *testing.T) {
	tc := &ToolCall{
		id:       "call-123",
		name:     "web_search",
		toolType: ToolCallTypeWebSearch,
	}

	if tc.Type() != ToolCallTypeWebSearch {
		t.Errorf("Expected Type() to return ToolCallTypeWebSearch, got %s", tc.Type())
	}
}

func TestParseToolCallWithType(t *testing.T) {
	protoCall := &xaiv1.ToolCall{
		Id:   "call-123",
		Type: xaiv1.ToolCallType_TOOL_CALL_TYPE_WEB_SEARCH_TOOL,
		Tool: &xaiv1.ToolCall_Function{
			Function: &xaiv1.FunctionCall{
				Name:      "search",
				Arguments: `{"query": "test"}`,
			},
		},
	}

	tc := parseToolCall(protoCall)
	if tc == nil {
		t.Fatal("Expected non-nil ToolCall")
	}

	if tc.Type() != ToolCallTypeWebSearch {
		t.Errorf("Expected Type() to return ToolCallTypeWebSearch, got %s", tc.Type())
	}

	if !tc.IsServerSide() {
		t.Error("Expected IsServerSide() to return true")
	}
}
