package chat

import (
	"testing"

	xaiv1 "github.com/ZaguanLabs/xai-sdk-go/proto/gen/go/xai/api/v1"
)

func TestResponseReasoningContent(t *testing.T) {
	tests := []struct {
		name     string
		response *Response
		want     string
	}{
		{
			name:     "nil response",
			response: &Response{},
			want:     "",
		},
		{
			name: "response with reasoning content",
			response: &Response{
				proto: &xaiv1.GetChatCompletionResponse{
					Outputs: []*xaiv1.CompletionOutput{
						{
							Message: &xaiv1.CompletionMessage{
								Content:          "The answer is 42",
								ReasoningContent: "Let me think about this...",
							},
						},
					},
				},
			},
			want: "Let me think about this...",
		},
		{
			name: "response without reasoning content",
			response: &Response{
				proto: &xaiv1.GetChatCompletionResponse{
					Outputs: []*xaiv1.CompletionOutput{
						{
							Message: &xaiv1.CompletionMessage{
								Content: "The answer is 42",
							},
						},
					},
				},
			},
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.response.ReasoningContent()
			if got != tt.want {
				t.Errorf("ReasoningContent() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestResponseEncryptedContent(t *testing.T) {
	tests := []struct {
		name     string
		response *Response
		want     string
	}{
		{
			name:     "nil response",
			response: &Response{},
			want:     "",
		},
		{
			name: "response with encrypted content",
			response: &Response{
				proto: &xaiv1.GetChatCompletionResponse{
					Outputs: []*xaiv1.CompletionOutput{
						{
							Message: &xaiv1.CompletionMessage{
								Content:          "The answer is 42",
								EncryptedContent: "encrypted_data_here",
							},
						},
					},
				},
			},
			want: "encrypted_data_here",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.response.EncryptedContent()
			if got != tt.want {
				t.Errorf("EncryptedContent() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestChunkReasoningContent(t *testing.T) {
	tests := []struct {
		name  string
		chunk *Chunk
		want  string
	}{
		{
			name:  "nil chunk",
			chunk: &Chunk{},
			want:  "",
		},
		{
			name: "chunk with reasoning content",
			chunk: &Chunk{
				proto: &xaiv1.GetChatCompletionChunk{
					Outputs: []*xaiv1.CompletionOutputChunk{
						{
							Delta: &xaiv1.Delta{
								Content:          "thinking...",
								ReasoningContent: "step 1: analyze",
							},
						},
					},
				},
			},
			want: "step 1: analyze",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.chunk.ReasoningContent()
			if got != tt.want {
				t.Errorf("ReasoningContent() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestChunkEncryptedContent(t *testing.T) {
	tests := []struct {
		name  string
		chunk *Chunk
		want  string
	}{
		{
			name:  "nil chunk",
			chunk: &Chunk{},
			want:  "",
		},
		{
			name: "chunk with encrypted content",
			chunk: &Chunk{
				proto: &xaiv1.GetChatCompletionChunk{
					Outputs: []*xaiv1.CompletionOutputChunk{
						{
							Delta: &xaiv1.Delta{
								Content:          "hello",
								EncryptedContent: "encrypted_chunk",
							},
						},
					},
				},
			},
			want: "encrypted_chunk",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.chunk.EncryptedContent()
			if got != tt.want {
				t.Errorf("EncryptedContent() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestAppendResponse(t *testing.T) {
	// Create a mock response with tool calls, reasoning, and encrypted content
	response := &Response{
		proto: &xaiv1.GetChatCompletionResponse{
			Outputs: []*xaiv1.CompletionOutput{
				{
					Message: &xaiv1.CompletionMessage{
						Role:             xaiv1.MessageRole_ROLE_ASSISTANT,
						Content:          "I need to call a tool",
						ReasoningContent: "Let me think about this",
						EncryptedContent: "encrypted_reasoning",
						ToolCalls: []*xaiv1.ToolCall{
							{
								Id:   "call_123",
								Type: xaiv1.ToolCallType_TOOL_CALL_TYPE_CLIENT_SIDE_TOOL,
								Tool: &xaiv1.ToolCall_Function{
									Function: &xaiv1.FunctionCall{
										Name:      "get_weather",
										Arguments: `{"city": "SF"}`,
									},
								},
							},
						},
					},
				},
			},
		},
	}

	// Create a request and append the response
	req := NewRequest("grok-beta")
	req.AppendResponse(response)

	// Verify the message was appended
	if len(req.proto.Messages) != 1 {
		t.Fatalf("Expected 1 message, got %d", len(req.proto.Messages))
	}

	msg := req.proto.Messages[0]

	// Verify role
	if msg.Role != xaiv1.MessageRole_ROLE_ASSISTANT {
		t.Errorf("Expected role ASSISTANT, got %v", msg.Role)
	}

	// Verify content
	if len(msg.Content) != 1 || msg.Content[0].GetText() != "I need to call a tool" {
		t.Errorf("Content not preserved correctly")
	}

	// Verify reasoning content
	if msg.ReasoningContent == nil || *msg.ReasoningContent != "Let me think about this" {
		var val string
		if msg.ReasoningContent != nil {
			val = *msg.ReasoningContent
		}
		t.Errorf("ReasoningContent = %q, want %q", val, "Let me think about this")
	}

	// Verify encrypted content
	if msg.EncryptedContent != "encrypted_reasoning" {
		t.Errorf("EncryptedContent = %q, want %q", msg.EncryptedContent, "encrypted_reasoning")
	}

	// Verify tool calls
	if len(msg.ToolCalls) != 1 {
		t.Fatalf("Expected 1 tool call, got %d", len(msg.ToolCalls))
	}

	if msg.ToolCalls[0].Id != "call_123" {
		t.Errorf("ToolCall ID = %q, want %q", msg.ToolCalls[0].Id, "call_123")
	}

	t.Log("✅ AppendResponse correctly preserves all fields")
}

func TestAppendResponseMultipleOutputs(t *testing.T) {
	// Create a response with multiple outputs (N > 1)
	response := &Response{
		proto: &xaiv1.GetChatCompletionResponse{
			Outputs: []*xaiv1.CompletionOutput{
				{
					Message: &xaiv1.CompletionMessage{
						Role:    xaiv1.MessageRole_ROLE_ASSISTANT,
						Content: "First response",
					},
				},
				{
					Message: &xaiv1.CompletionMessage{
						Role:    xaiv1.MessageRole_ROLE_ASSISTANT,
						Content: "Second response",
					},
				},
			},
		},
	}

	req := NewRequest("grok-beta")
	req.AppendResponse(response)

	// Verify both outputs were appended
	if len(req.proto.Messages) != 2 {
		t.Fatalf("Expected 2 messages, got %d", len(req.proto.Messages))
	}

	if req.proto.Messages[0].Content[0].GetText() != "First response" {
		t.Errorf("First message content incorrect")
	}

	if req.proto.Messages[1].Content[0].GetText() != "Second response" {
		t.Errorf("Second message content incorrect")
	}

	t.Log("✅ AppendResponse correctly handles multiple outputs (N > 1)")
}
