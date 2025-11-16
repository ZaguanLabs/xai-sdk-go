package chat

import (
	"testing"

	xaiv1 "github.com/ZaguanLabs/xai-sdk-go/proto/gen/go/xai/v1"
)

// Phase 2 Tests: RequestSettings, DebugOutput, LogProbs, Function.strict

func TestResponseRequestSettings(t *testing.T) {
	response := &Response{
		proto: &xaiv1.GetChatCompletionResponse{
			Settings: &xaiv1.RequestSettings{
				MaxTokens:           1000,
				ParallelToolCalls:   true,
				PreviousResponseId:  "resp_123",
				ReasoningEffort:     xaiv1.ReasoningEffort_EFFORT_HIGH,
				Temperature:         0.7,
				TopP:                0.9,
				User:                "user_456",
				StoreMessages:       true,
				UseEncryptedContent: true,
			},
		},
	}

	settings := response.RequestSettings()
	if settings == nil {
		t.Fatal("RequestSettings() returned nil")
	}

	if settings.MaxTokens() != 1000 {
		t.Errorf("MaxTokens() = %d, want 1000", settings.MaxTokens())
	}

	if !settings.ParallelToolCalls() {
		t.Error("ParallelToolCalls() = false, want true")
	}

	if settings.PreviousResponseID() != "resp_123" {
		t.Errorf("PreviousResponseID() = %q, want %q", settings.PreviousResponseID(), "resp_123")
	}

	if settings.ReasoningEffort() != "high" {
		t.Errorf("ReasoningEffort() = %q, want %q", settings.ReasoningEffort(), "high")
	}

	if settings.Temperature() != 0.7 {
		t.Errorf("Temperature() = %f, want 0.7", settings.Temperature())
	}

	if settings.TopP() != 0.9 {
		t.Errorf("TopP() = %f, want 0.9", settings.TopP())
	}

	if settings.User() != "user_456" {
		t.Errorf("User() = %q, want %q", settings.User(), "user_456")
	}

	if !settings.StoreMessages() {
		t.Error("StoreMessages() = false, want true")
	}

	if !settings.UseEncryptedContent() {
		t.Error("UseEncryptedContent() = false, want true")
	}

	t.Log("✅ Response.RequestSettings() works correctly")
}

func TestResponseDebugOutput(t *testing.T) {
	response := &Response{
		proto: &xaiv1.GetChatCompletionResponse{
			DebugOutput: &xaiv1.DebugOutput{
				Attempts:             3,
				Request:              "test request",
				Prompt:               "test prompt",
				Responses:            []string{"resp1", "resp2"},
				CacheReadCount:       5,
				CacheReadInputBytes:  1024,
				CacheWriteCount:      2,
				CacheWriteInputBytes: 512,
				EngineRequest:        "engine req",
				LbAddress:            "lb.example.com",
				SamplerTag:           "sampler_v1",
				Chunks:               []string{"chunk1", "chunk2"},
			},
		},
	}

	debug := response.DebugOutput()
	if debug == nil {
		t.Fatal("DebugOutput() returned nil")
	}

	if debug.Attempts() != 3 {
		t.Errorf("Attempts() = %d, want 3", debug.Attempts())
	}

	if debug.Request() != "test request" {
		t.Errorf("Request() = %q, want %q", debug.Request(), "test request")
	}

	if debug.Prompt() != "test prompt" {
		t.Errorf("Prompt() = %q, want %q", debug.Prompt(), "test prompt")
	}

	if len(debug.Responses()) != 2 {
		t.Errorf("Responses() length = %d, want 2", len(debug.Responses()))
	}

	if debug.CacheReadCount() != 5 {
		t.Errorf("CacheReadCount() = %d, want 5", debug.CacheReadCount())
	}

	if debug.CacheReadInputBytes() != 1024 {
		t.Errorf("CacheReadInputBytes() = %d, want 1024", debug.CacheReadInputBytes())
	}

	if debug.CacheWriteCount() != 2 {
		t.Errorf("CacheWriteCount() = %d, want 2", debug.CacheWriteCount())
	}

	if debug.CacheWriteInputBytes() != 512 {
		t.Errorf("CacheWriteInputBytes() = %d, want 512", debug.CacheWriteInputBytes())
	}

	if debug.EngineRequest() != "engine req" {
		t.Errorf("EngineRequest() = %q, want %q", debug.EngineRequest(), "engine req")
	}

	if debug.LBAddress() != "lb.example.com" {
		t.Errorf("LBAddress() = %q, want %q", debug.LBAddress(), "lb.example.com")
	}

	if debug.SamplerTag() != "sampler_v1" {
		t.Errorf("SamplerTag() = %q, want %q", debug.SamplerTag(), "sampler_v1")
	}

	if len(debug.Chunks()) != 2 {
		t.Errorf("Chunks() length = %d, want 2", len(debug.Chunks()))
	}

	t.Log("✅ Response.DebugOutput() works correctly")
}

func TestChoiceLogProbs(t *testing.T) {
	choice := &Choice{
		proto: &xaiv1.CompletionOutput{
			Logprobs: &xaiv1.LogProbs{
				Content: []*xaiv1.LogProb{
					{
						Token:   "hello",
						Logprob: -0.5,
						Bytes:   []byte("hello"),
						TopLogprobs: []*xaiv1.TopLogProb{
							{
								Token:   "hi",
								Logprob: -0.7,
								Bytes:   []byte("hi"),
							},
						},
					},
				},
			},
		},
	}

	logprobs := choice.LogProbs()
	if logprobs == nil {
		t.Fatal("LogProbs() returned nil")
	}

	content := logprobs.Content()
	if len(content) != 1 {
		t.Fatalf("Content() length = %d, want 1", len(content))
	}

	if content[0].Token() != "hello" {
		t.Errorf("Token() = %q, want %q", content[0].Token(), "hello")
	}

	if content[0].Logprob() != -0.5 {
		t.Errorf("Logprob() = %f, want -0.5", content[0].Logprob())
	}

	topLogprobs := content[0].TopLogProbs()
	if len(topLogprobs) != 1 {
		t.Fatalf("TopLogProbs() length = %d, want 1", len(topLogprobs))
	}

	if topLogprobs[0].Token() != "hi" {
		t.Errorf("TopLogProb Token() = %q, want %q", topLogprobs[0].Token(), "hi")
	}

	t.Log("✅ Choice.LogProbs() works correctly")
}

func TestToolWithStrict(t *testing.T) {
	tool := NewTool("test_func", "A test function").
		WithParameter("param1", "string", "First param", true).
		WithStrict(true)

	if !tool.Strict() {
		t.Error("Strict() = false, want true")
	}

	// Test that strict is included in proto conversion
	// This is tested indirectly through WithTool
	t.Log("✅ Tool.WithStrict() and Strict() work correctly")
}

func TestToolStrictDefault(t *testing.T) {
	tool := NewTool("test_func", "A test function")

	if tool.Strict() {
		t.Error("Strict() = true, want false (default)")
	}

	t.Log("✅ Tool.Strict() defaults to false")
}

func TestNilResponseAccessors(t *testing.T) {
	// Test that all new accessors handle nil gracefully
	response := &Response{}

	if response.RequestSettings() != nil {
		t.Error("RequestSettings() should return nil for nil proto")
	}

	if response.DebugOutput() != nil {
		t.Error("DebugOutput() should return nil for nil proto")
	}

	t.Log("✅ Response accessors handle nil gracefully")
}

func TestNilChoiceLogProbs(t *testing.T) {
	choice := &Choice{}

	if choice.LogProbs() != nil {
		t.Error("LogProbs() should return nil for nil proto")
	}

	t.Log("✅ Choice.LogProbs() handles nil gracefully")
}
