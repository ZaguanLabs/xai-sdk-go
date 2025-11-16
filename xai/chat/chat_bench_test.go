package chat

import (
	"testing"
)

// BenchmarkNewRequest benchmarks creating a new chat request
func BenchmarkNewRequest(b *testing.B) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = NewRequest("grok-1.5-flash",
			WithTemperature(0.7),
			WithMaxTokens(1000),
			WithMessages(
				System(Text("You are a helpful assistant.")),
				User(Text("Hello, world!")),
			),
		)
	}
}

// BenchmarkRequestValidate benchmarks request validation
func BenchmarkRequestValidate(b *testing.B) {
	req := NewRequest("grok-1.5-flash",
		WithTemperature(0.7),
		WithMaxTokens(1000),
		WithMessages(
			System(Text("You are a helpful assistant.")),
			User(Text("Hello, world!")),
		),
	)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = req.validate()
	}
}

// BenchmarkMessageBuilder benchmarks message construction
func BenchmarkMessageBuilder(b *testing.B) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = User(Text("Hello, world!"))
	}
}

// BenchmarkMultipleMessages benchmarks building multiple messages
func BenchmarkMultipleMessages(b *testing.B) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		messages := []*Message{
			System(Text("You are a helpful assistant.")),
			User(Text("What is the capital of France?")),
			Assistant(Text("The capital of France is Paris.")),
			User(Text("What about Germany?")),
		}
		_ = messages
	}
}
