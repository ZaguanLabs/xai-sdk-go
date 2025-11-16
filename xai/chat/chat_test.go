package chat

import (
	"testing"
)

func TestMessageBuilders(t *testing.T) {
	// Test System message
	systemMsg := System(Text("You are a helpful assistant."))
	if systemMsg.Role() != "system" {
		t.Errorf("Expected role 'system', got '%s'", systemMsg.Role())
	}
	if systemMsg.Content() != "You are a helpful assistant." {
		t.Errorf("Expected content 'You are a helpful assistant.', got '%s'", systemMsg.Content())
	}

	// Test User message
	userMsg := User(Text("Hello, how are you?"))
	if userMsg.Role() != "user" {
		t.Errorf("Expected role 'user', got '%s'", userMsg.Role())
	}
	if userMsg.Content() != "Hello, how are you?" {
		t.Errorf("Expected content 'Hello, how are you?', got '%s'", userMsg.Content())
	}

	// Test Assistant message
	assistantMsg := Assistant(Text("I'm doing well, thank you!"))
	if assistantMsg.Role() != "assistant" {
		t.Errorf("Expected role 'assistant', got '%s'", assistantMsg.Role())
	}
	if assistantMsg.Content() != "I'm doing well, thank you!" {
		t.Errorf("Expected content 'I'm doing well, thank you!', got '%s'", assistantMsg.Content())
	}
}

func TestRequestBuilders(t *testing.T) {
	// Test NewRequest with model
	req := NewRequest("test-model")
	if req.GetModel() != "test-model" {
		t.Errorf("Expected model 'test-model', got '%s'", req.GetModel())
	}

	// Test WithTemperature option
	req = NewRequest("test-model", WithTemperature(0.8))
	if req.Proto().GetTemperature() != 0.8 {
		t.Errorf("Expected temperature 0.8, got %f", req.Proto().GetTemperature())
	}

	// Test WithMaxTokens option
	req = NewRequest("test-model", WithMaxTokens(100))
	if req.Proto().GetMaxTokens() != 100 {
		t.Errorf("Expected max tokens 100, got %d", req.Proto().GetMaxTokens())
	}

	// Test WithMessages option
	req = NewRequest("test-model",
		WithMessages(
			System(Text("You are a helpful assistant.")),
			User(Text("Hello!")),
		))
	if len(req.Proto().GetMessages()) != 2 {
		t.Errorf("Expected 2 messages, got %d", len(req.Proto().GetMessages()))
	}
}
