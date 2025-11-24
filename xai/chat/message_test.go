package chat

import (
	"testing"

	xaiv1 "github.com/ZaguanLabs/xai-sdk-go/proto/gen/go/xai/api/v1"
)

func TestMessageWithToolCalls(t *testing.T) {
	msg := Assistant(Text("I'll help you with that"))

	toolCall := NewToolCall("call_123", "get_weather", map[string]interface{}{
		"city": "San Francisco",
	})

	msg.WithToolCalls([]*ToolCall{toolCall})

	// Verify tool calls were set
	toolCalls := msg.ToolCalls()
	if len(toolCalls) != 1 {
		t.Fatalf("Expected 1 tool call, got %d", len(toolCalls))
	}

	if toolCalls[0].ID() != "call_123" {
		t.Errorf("ToolCall ID = %q, want %q", toolCalls[0].ID(), "call_123")
	}

	if toolCalls[0].Name() != "get_weather" {
		t.Errorf("ToolCall Name = %q, want %q", toolCalls[0].Name(), "get_weather")
	}

	// Verify proto was updated
	if len(msg.proto.ToolCalls) != 1 {
		t.Fatalf("Expected 1 proto tool call, got %d", len(msg.proto.ToolCalls))
	}

	if msg.proto.ToolCalls[0].Id != "call_123" {
		t.Errorf("Proto ToolCall ID = %q, want %q", msg.proto.ToolCalls[0].Id, "call_123")
	}

	t.Log("✅ Message.WithToolCalls() works correctly")
}

func TestMessageWithReasoningContent(t *testing.T) {
	msg := Assistant(Text("The answer is 42"))
	msg.WithReasoningContent("Let me think about this...")

	reasoning := msg.ReasoningContent()
	if reasoning != "Let me think about this..." {
		t.Errorf("ReasoningContent = %q, want %q", reasoning, "Let me think about this...")
	}

	// Verify proto was updated
	if msg.proto.ReasoningContent == nil || *msg.proto.ReasoningContent != "Let me think about this..." {
		var val string
		if msg.proto.ReasoningContent != nil {
			val = *msg.proto.ReasoningContent
		}
		t.Errorf("Proto ReasoningContent = %q, want %q", val, "Let me think about this...")
	}

	t.Log("✅ Message.WithReasoningContent() works correctly")
}

func TestMessageWithEncryptedContent(t *testing.T) {
	msg := Assistant(Text("Hello"))
	msg.WithEncryptedContent("encrypted_data_here")

	encrypted := msg.EncryptedContent()
	if encrypted != "encrypted_data_here" {
		t.Errorf("EncryptedContent = %q, want %q", encrypted, "encrypted_data_here")
	}

	// Verify proto was updated
	if msg.proto.EncryptedContent != "encrypted_data_here" {
		t.Errorf("Proto EncryptedContent = %q, want %q", msg.proto.EncryptedContent, "encrypted_data_here")
	}

	t.Log("✅ Message.WithEncryptedContent() works correctly")
}

func TestMessageToolCallsAccessor(t *testing.T) {
	// Create a message with proto tool calls directly
	msg := &Message{
		proto: &xaiv1.Message{
			Role: xaiv1.MessageRole_ROLE_ASSISTANT,
			Content: []*xaiv1.Content{
				{
					Content: &xaiv1.Content_Text{
						Text: "I'll call a tool",
					},
				},
			},
			ToolCalls: []*xaiv1.ToolCall{
				{
					Id:   "call_456",
					Type: xaiv1.ToolCallType_TOOL_CALL_TYPE_CLIENT_SIDE_TOOL,
					Tool: &xaiv1.ToolCall_Function{
						Function: &xaiv1.FunctionCall{
							Name:      "get_time",
							Arguments: `{"timezone": "UTC"}`,
						},
					},
				},
			},
		},
	}

	toolCalls := msg.ToolCalls()
	if len(toolCalls) != 1 {
		t.Fatalf("Expected 1 tool call, got %d", len(toolCalls))
	}

	if toolCalls[0].ID() != "call_456" {
		t.Errorf("ToolCall ID = %q, want %q", toolCalls[0].ID(), "call_456")
	}

	if toolCalls[0].Name() != "get_time" {
		t.Errorf("ToolCall Name = %q, want %q", toolCalls[0].Name(), "get_time")
	}

	args := toolCalls[0].Arguments()
	if args["timezone"] != "UTC" {
		t.Errorf("ToolCall timezone argument = %v, want %q", args["timezone"], "UTC")
	}

	t.Log("✅ Message.ToolCalls() accessor works correctly")
}

func TestMessageChaining(t *testing.T) {
	// Test that all With* methods return the message for chaining
	msg := Assistant(Text("Hello")).
		WithReasoningContent("thinking...").
		WithEncryptedContent("encrypted").
		WithToolCalls([]*ToolCall{
			NewToolCall("call_789", "test_func", map[string]interface{}{}),
		})

	if msg.ReasoningContent() != "thinking..." {
		t.Error("Chaining failed for ReasoningContent")
	}

	if msg.EncryptedContent() != "encrypted" {
		t.Error("Chaining failed for EncryptedContent")
	}

	if len(msg.ToolCalls()) != 1 {
		t.Error("Chaining failed for ToolCalls")
	}

	t.Log("✅ Message method chaining works correctly")
}

func TestMessageEmptyToolCalls(t *testing.T) {
	msg := User(Text("Hello"))

	toolCalls := msg.ToolCalls()
	if toolCalls != nil {
		t.Errorf("Expected nil tool calls for message without tools, got %v", toolCalls)
	}

	t.Log("✅ Message.ToolCalls() returns nil for empty tool calls")
}

func TestMessageWithNilToolCalls(t *testing.T) {
	msg := Assistant(Text("Hello"))

	// Test with nil tool call in the array
	msg.WithToolCalls([]*ToolCall{nil, NewToolCall("call_1", "func1", nil), nil})

	toolCalls := msg.ToolCalls()
	if len(toolCalls) != 1 {
		t.Errorf("Expected 1 tool call (nil should be skipped), got %d", len(toolCalls))
	}

	t.Log("✅ Message.WithToolCalls() correctly skips nil tool calls")
}
