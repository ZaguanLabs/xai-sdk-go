package chat

import (
	"testing"

	xaiv1 "github.com/ZaguanLabs/xai-sdk-go/proto/gen/go/xai/api/v1"
)

func TestToolCallStatus(t *testing.T) {
	tests := []struct {
		name           string
		protoCall      *xaiv1.ToolCall
		expectedStatus string
		expectedError  string
	}{
		{
			name: "in_progress_status",
			protoCall: &xaiv1.ToolCall{
				Id:     "call_123",
				Type:   xaiv1.ToolCallType_TOOL_CALL_TYPE_WEB_SEARCH_TOOL,
				Status: xaiv1.ToolCallStatus_TOOL_CALL_STATUS_IN_PROGRESS,
				Tool: &xaiv1.ToolCall_Function{
					Function: &xaiv1.FunctionCall{
						Name:      "web_search",
						Arguments: `{"query":"test"}`,
					},
				},
			},
			expectedStatus: "TOOL_CALL_STATUS_IN_PROGRESS",
			expectedError:  "",
		},
		{
			name: "completed_status",
			protoCall: &xaiv1.ToolCall{
				Id:     "call_456",
				Type:   xaiv1.ToolCallType_TOOL_CALL_TYPE_CODE_EXECUTION_TOOL,
				Status: xaiv1.ToolCallStatus_TOOL_CALL_STATUS_COMPLETED,
				Tool: &xaiv1.ToolCall_Function{
					Function: &xaiv1.FunctionCall{
						Name:      "execute_code",
						Arguments: `{"code":"print('hello')"}`,
					},
				},
			},
			expectedStatus: "TOOL_CALL_STATUS_COMPLETED",
			expectedError:  "",
		},
		{
			name: "failed_status_with_error",
			protoCall: &xaiv1.ToolCall{
				Id:           "call_789",
				Type:         xaiv1.ToolCallType_TOOL_CALL_TYPE_CLIENT_SIDE_TOOL,
				Status:       xaiv1.ToolCallStatus_TOOL_CALL_STATUS_FAILED,
				ErrorMessage: func() *string { s := "Connection timeout"; return &s }(),
				Tool: &xaiv1.ToolCall_Function{
					Function: &xaiv1.FunctionCall{
						Name:      "get_weather",
						Arguments: `{"city":"London"}`,
					},
				},
			},
			expectedStatus: "TOOL_CALL_STATUS_FAILED",
			expectedError:  "Connection timeout",
		},
		{
			name: "incomplete_status",
			protoCall: &xaiv1.ToolCall{
				Id:     "call_abc",
				Type:   xaiv1.ToolCallType_TOOL_CALL_TYPE_X_SEARCH_TOOL,
				Status: xaiv1.ToolCallStatus_TOOL_CALL_STATUS_INCOMPLETE,
				Tool: &xaiv1.ToolCall_Function{
					Function: &xaiv1.FunctionCall{
						Name:      "x_search",
						Arguments: `{"query":"#ai"}`,
					},
				},
			},
			expectedStatus: "TOOL_CALL_STATUS_INCOMPLETE",
			expectedError:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			toolCall := parseToolCall(tt.protoCall)

			if toolCall == nil {
				t.Fatal("parseToolCall returned nil")
			}

			if toolCall.Status() != tt.expectedStatus {
				t.Errorf("Status() = %q, want %q", toolCall.Status(), tt.expectedStatus)
			}

			if toolCall.ErrorMessage() != tt.expectedError {
				t.Errorf("ErrorMessage() = %q, want %q", toolCall.ErrorMessage(), tt.expectedError)
			}

			// Verify other fields are still parsed correctly
			if toolCall.ID() != tt.protoCall.Id {
				t.Errorf("ID() = %q, want %q", toolCall.ID(), tt.protoCall.Id)
			}

			if toolCall.Name() != tt.protoCall.GetFunction().Name {
				t.Errorf("Name() = %q, want %q", toolCall.Name(), tt.protoCall.GetFunction().Name)
			}
		})
	}
}

func TestToolCallStatusConversion(t *testing.T) {
	tests := []struct {
		name         string
		statusString string
		expectedEnum xaiv1.ToolCallStatus
	}{
		{
			name:         "in_progress",
			statusString: "TOOL_CALL_STATUS_IN_PROGRESS",
			expectedEnum: xaiv1.ToolCallStatus_TOOL_CALL_STATUS_IN_PROGRESS,
		},
		{
			name:         "completed",
			statusString: "TOOL_CALL_STATUS_COMPLETED",
			expectedEnum: xaiv1.ToolCallStatus_TOOL_CALL_STATUS_COMPLETED,
		},
		{
			name:         "incomplete",
			statusString: "TOOL_CALL_STATUS_INCOMPLETE",
			expectedEnum: xaiv1.ToolCallStatus_TOOL_CALL_STATUS_INCOMPLETE,
		},
		{
			name:         "failed",
			statusString: "TOOL_CALL_STATUS_FAILED",
			expectedEnum: xaiv1.ToolCallStatus_TOOL_CALL_STATUS_FAILED,
		},
		{
			name:         "unknown_defaults_to_in_progress",
			statusString: "UNKNOWN_STATUS",
			expectedEnum: xaiv1.ToolCallStatus_TOOL_CALL_STATUS_IN_PROGRESS,
		},
		{
			name:         "empty_defaults_to_in_progress",
			statusString: "",
			expectedEnum: xaiv1.ToolCallStatus_TOOL_CALL_STATUS_IN_PROGRESS,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseToolCallStatus(tt.statusString)
			if result != tt.expectedEnum {
				t.Errorf("parseToolCallStatus(%q) = %v, want %v", tt.statusString, result, tt.expectedEnum)
			}
		})
	}
}

func TestMessageWithToolCallsStatus(t *testing.T) {
	// Create a tool call with status and error message
	toolCall := &ToolCall{
		id:           "call_123",
		name:         "get_weather",
		arguments:    map[string]interface{}{"city": "Paris"},
		status:       "TOOL_CALL_STATUS_FAILED",
		errorMessage: "API rate limit exceeded",
	}

	msg := Assistant(Text("I'll check the weather"))
	msg.WithToolCalls([]*ToolCall{toolCall})

	// Verify the proto was set correctly
	if len(msg.proto.ToolCalls) != 1 {
		t.Fatalf("Expected 1 tool call, got %d", len(msg.proto.ToolCalls))
	}

	protoCall := msg.proto.ToolCalls[0]

	if protoCall.Status != xaiv1.ToolCallStatus_TOOL_CALL_STATUS_FAILED {
		t.Errorf("Status = %v, want TOOL_CALL_STATUS_FAILED", protoCall.Status)
	}

	if protoCall.ErrorMessage == nil || *protoCall.ErrorMessage != "API rate limit exceeded" {
		var val string
		if protoCall.ErrorMessage != nil {
			val = *protoCall.ErrorMessage
		}
		t.Errorf("ErrorMessage = %q, want %q", val, "API rate limit exceeded")
	}

	if protoCall.Id != "call_123" {
		t.Errorf("Id = %q, want %q", protoCall.Id, "call_123")
	}
}

func TestToolCallStatusRoundTrip(t *testing.T) {
	// Test that we can convert from proto -> SDK -> proto without losing data
	originalProto := &xaiv1.ToolCall{
		Id:           "call_xyz",
		Type:         xaiv1.ToolCallType_TOOL_CALL_TYPE_WEB_SEARCH_TOOL,
		Status:       xaiv1.ToolCallStatus_TOOL_CALL_STATUS_COMPLETED,
		ErrorMessage: nil,
		Tool: &xaiv1.ToolCall_Function{
			Function: &xaiv1.FunctionCall{
				Name:      "web_search",
				Arguments: `{"query":"golang"}`,
			},
		},
	}

	// Parse to SDK type
	sdkCall := parseToolCall(originalProto)

	// Convert back to proto via Message
	msg := Assistant(Text("Result"))
	msg.WithToolCalls([]*ToolCall{sdkCall})

	// Verify round-trip
	if len(msg.proto.ToolCalls) != 1 {
		t.Fatalf("Expected 1 tool call, got %d", len(msg.proto.ToolCalls))
	}

	roundTripProto := msg.proto.ToolCalls[0]

	if roundTripProto.Id != originalProto.Id {
		t.Errorf("Id mismatch: got %q, want %q", roundTripProto.Id, originalProto.Id)
	}

	if roundTripProto.Status != originalProto.Status {
		t.Errorf("Status mismatch: got %v, want %v", roundTripProto.Status, originalProto.Status)
	}

	if (roundTripProto.ErrorMessage == nil) != (originalProto.ErrorMessage == nil) {
		t.Errorf("ErrorMessage nil mismatch")
	} else if roundTripProto.ErrorMessage != nil && *roundTripProto.ErrorMessage != *originalProto.ErrorMessage {
		t.Errorf("ErrorMessage mismatch: got %q, want %q", *roundTripProto.ErrorMessage, *originalProto.ErrorMessage)
	}

	if roundTripProto.GetFunction().Name != originalProto.GetFunction().Name {
		t.Errorf("Function.Name mismatch: got %q, want %q", roundTripProto.GetFunction().Name, originalProto.GetFunction().Name)
	}
}

func TestMultipleToolCallsWithDifferentStatuses(t *testing.T) {
	// Simulate multiple tool calls with different statuses (e.g., from streaming)
	toolCalls := []*ToolCall{
		{
			id:           "call_1",
			name:         "search_web",
			arguments:    map[string]interface{}{"query": "test"},
			status:       "TOOL_CALL_STATUS_IN_PROGRESS",
			errorMessage: "",
		},
		{
			id:           "call_1", // Same ID, different status (progress update)
			name:         "search_web",
			arguments:    map[string]interface{}{"query": "test"},
			status:       "TOOL_CALL_STATUS_COMPLETED",
			errorMessage: "",
		},
		{
			id:           "call_2",
			name:         "execute_code",
			arguments:    map[string]interface{}{"code": "print('hi')"},
			status:       "TOOL_CALL_STATUS_FAILED",
			errorMessage: "Syntax error",
		},
	}

	msg := Assistant(Text("Processing..."))
	msg.WithToolCalls(toolCalls)

	if len(msg.proto.ToolCalls) != 3 {
		t.Fatalf("Expected 3 tool calls, got %d", len(msg.proto.ToolCalls))
	}

	// Verify first call (in progress)
	if msg.proto.ToolCalls[0].Status != xaiv1.ToolCallStatus_TOOL_CALL_STATUS_IN_PROGRESS {
		t.Errorf("First call status = %v, want IN_PROGRESS", msg.proto.ToolCalls[0].Status)
	}

	// Verify second call (completed, same ID as first)
	if msg.proto.ToolCalls[1].Status != xaiv1.ToolCallStatus_TOOL_CALL_STATUS_COMPLETED {
		t.Errorf("Second call status = %v, want COMPLETED", msg.proto.ToolCalls[1].Status)
	}

	if msg.proto.ToolCalls[1].Id != msg.proto.ToolCalls[0].Id {
		t.Errorf("Second call should have same ID as first")
	}

	// Verify third call (failed with error)
	if msg.proto.ToolCalls[2].Status != xaiv1.ToolCallStatus_TOOL_CALL_STATUS_FAILED {
		t.Errorf("Third call status = %v, want FAILED", msg.proto.ToolCalls[2].Status)
	}

	if msg.proto.ToolCalls[2].ErrorMessage == nil || *msg.proto.ToolCalls[2].ErrorMessage != "Syntax error" {
		var val string
		if msg.proto.ToolCalls[2].ErrorMessage != nil {
			val = *msg.proto.ToolCalls[2].ErrorMessage
		}
		t.Errorf("Third call error = %q, want %q", val, "Syntax error")
	}
}
