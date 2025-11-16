package chat

import (
	"testing"

	xaiv1 "github.com/ZaguanLabs/xai-sdk-go/proto/gen/go/xai/v1"
)

// Phase 1 Tests: Message.Name, Response.Citations, Response.SystemFingerprint, Chunk.Citations, Chunk.SystemFingerprint

func TestMessageName(t *testing.T) {
	msg := User(Text("Hello"))

	// Test default (empty)
	if msg.Name() != "" {
		t.Errorf("Name() = %q, want empty string", msg.Name())
	}

	// Test WithName
	msg.WithName("Alice")
	if msg.Name() != "Alice" {
		t.Errorf("Name() = %q, want %q", msg.Name(), "Alice")
	}

	// Verify proto was updated
	if msg.Proto().Name != "Alice" {
		t.Errorf("Proto().Name = %q, want %q", msg.Proto().Name, "Alice")
	}

	t.Log("✅ Message.Name() and WithName() work correctly")
}

func TestMessageNameChaining(t *testing.T) {
	msg := User(Text("Hello")).
		WithName("Bob").
		WithReasoningContent("thinking...")

	if msg.Name() != "Bob" {
		t.Errorf("Name() = %q, want %q", msg.Name(), "Bob")
	}

	if msg.ReasoningContent() != "thinking..." {
		t.Error("Chaining failed for ReasoningContent")
	}

	t.Log("✅ Message.WithName() chaining works correctly")
}

func TestResponseCitations(t *testing.T) {
	tests := []struct {
		name     string
		response *Response
		want     []string
	}{
		{
			name:     "nil response",
			response: &Response{},
			want:     nil,
		},
		{
			name: "response with citations",
			response: &Response{
				proto: &xaiv1.GetChatCompletionResponse{
					Citations: []string{
						"https://example.com/source1",
						"https://example.com/source2",
					},
				},
			},
			want: []string{
				"https://example.com/source1",
				"https://example.com/source2",
			},
		},
		{
			name: "response without citations",
			response: &Response{
				proto: &xaiv1.GetChatCompletionResponse{},
			},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.response.Citations()
			if len(got) != len(tt.want) {
				t.Errorf("Citations() length = %d, want %d", len(got), len(tt.want))
				return
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("Citations()[%d] = %q, want %q", i, got[i], tt.want[i])
				}
			}
		})
	}

	t.Log("✅ Response.Citations() works correctly")
}

func TestResponseSystemFingerprint(t *testing.T) {
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
			name: "response with fingerprint",
			response: &Response{
				proto: &xaiv1.GetChatCompletionResponse{
					SystemFingerprint: "fp_1234567890",
				},
			},
			want: "fp_1234567890",
		},
		{
			name: "response without fingerprint",
			response: &Response{
				proto: &xaiv1.GetChatCompletionResponse{},
			},
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.response.SystemFingerprint()
			if got != tt.want {
				t.Errorf("SystemFingerprint() = %q, want %q", got, tt.want)
			}
		})
	}

	t.Log("✅ Response.SystemFingerprint() works correctly")
}

func TestChunkCitations(t *testing.T) {
	tests := []struct {
		name  string
		chunk *Chunk
		want  []string
	}{
		{
			name:  "nil chunk",
			chunk: &Chunk{},
			want:  nil,
		},
		{
			name: "chunk with citations",
			chunk: &Chunk{
				proto: &xaiv1.GetChatCompletionChunk{
					Citations: []string{
						"https://example.com/stream1",
						"https://example.com/stream2",
					},
				},
			},
			want: []string{
				"https://example.com/stream1",
				"https://example.com/stream2",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.chunk.Citations()
			if len(got) != len(tt.want) {
				t.Errorf("Citations() length = %d, want %d", len(got), len(tt.want))
				return
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("Citations()[%d] = %q, want %q", i, got[i], tt.want[i])
				}
			}
		})
	}

	t.Log("✅ Chunk.Citations() works correctly")
}

func TestChunkSystemFingerprint(t *testing.T) {
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
			name: "chunk with fingerprint",
			chunk: &Chunk{
				proto: &xaiv1.GetChatCompletionChunk{
					SystemFingerprint: "fp_stream_123",
				},
			},
			want: "fp_stream_123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.chunk.SystemFingerprint()
			if got != tt.want {
				t.Errorf("SystemFingerprint() = %q, want %q", got, tt.want)
			}
		})
	}

	t.Log("✅ Chunk.SystemFingerprint() works correctly")
}

func TestChunkUsage(t *testing.T) {
	// Test that Chunk.Usage() now properly returns usage info
	chunk := &Chunk{
		proto: &xaiv1.GetChatCompletionChunk{
			Usage: &xaiv1.SamplingUsage{
				PromptTokens:     10,
				CompletionTokens: 20,
				TotalTokens:      30,
			},
		},
	}

	usage := chunk.Usage()
	if usage == nil {
		t.Fatal("Usage() returned nil, expected usage info")
	}

	if usage.PromptTokens() != 10 {
		t.Errorf("PromptTokens() = %d, want 10", usage.PromptTokens())
	}

	if usage.CompletionTokens() != 20 {
		t.Errorf("CompletionTokens() = %d, want 20", usage.CompletionTokens())
	}

	if usage.TotalTokens() != 30 {
		t.Errorf("TotalTokens() = %d, want 30", usage.TotalTokens())
	}

	t.Log("✅ Chunk.Usage() now works correctly (was previously always nil)")
}
