// Package main demonstrates all available chat parameters including the newly added high-priority ones.
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/ZaguanLabs/xai-sdk-go/xai"
	"github.com/ZaguanLabs/xai-sdk-go/xai/chat"
)

func main() {
	// Get API key from environment
	apiKey := os.Getenv("XAI_API_KEY")
	if apiKey == "" {
		log.Fatal("XAI_API_KEY environment variable is required")
	}

	// Create client
	client, err := xai.NewClient(&xai.Config{
		APIKey: apiKey,
	})
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	ctx := context.Background()

	// Example 1: Deterministic outputs with Seed
	fmt.Println("=== Example 1: Deterministic Outputs (Seed) ===")
	req1 := chat.NewRequest("grok-beta",
		chat.WithMessage(chat.User(chat.Text("Generate a random number between 1 and 100"))),
		chat.WithSeed(12345), // Same seed = same output
		chat.WithTemperature(0.7),
	)
	req1.SetMaxTokens(50)

	resp1, err := req1.Sample(ctx, client.Chat())
	if err != nil {
		log.Fatalf("Seed example failed: %v", err)
	}
	fmt.Printf("Response 1: %s\n", resp1.Content())

	// Run again with same seed - should get same result
	resp1b, err := req1.Sample(ctx, client.Chat())
	if err != nil {
		log.Fatalf("Seed example (repeat) failed: %v", err)
	}
	fmt.Printf("Response 2 (same seed): %s\n", resp1b.Content())
	fmt.Println()

	// Example 2: Log Probabilities for confidence scoring
	fmt.Println("=== Example 2: Log Probabilities ===")
	req2 := chat.NewRequest("grok-beta",
		chat.WithMessage(chat.User(chat.Text("Is the sky blue? Answer yes or no."))),
		chat.WithLogprobs(true),
		chat.WithTopLogprobs(3), // Show top 3 alternative tokens
		chat.WithTemperature(0.3),
	)
	req2.SetMaxTokens(10)

	resp2, err := req2.Sample(ctx, client.Chat())
	if err != nil {
		log.Fatalf("Logprobs example failed: %v", err)
	}
	fmt.Printf("Response: %s\n", resp2.Content())
	fmt.Printf("Note: Log probabilities are available in the response for confidence analysis\n\n")

	// Example 3: Multiple completions with N
	fmt.Println("=== Example 3: Multiple Completions (N) ===")
	req3 := chat.NewRequest("grok-beta",
		chat.WithMessage(chat.System(chat.Text("You are a creative writer."))),
		chat.WithMessage(chat.User(chat.Text("Write a one-sentence story about a robot."))),
		chat.WithN(3), // Generate 3 different completions
		chat.WithTemperature(0.9),
	)
	req3.SetMaxTokens(50)

	resp3, err := req3.Sample(ctx, client.Chat())
	if err != nil {
		log.Fatalf("Multiple completions example failed: %v", err)
	}
	fmt.Printf("Generated %d completion(s):\n", resp3.ChoiceCount())
	for i, choice := range resp3.Choices() {
		fmt.Printf("  %d. %s\n", i+1, choice.Message().Content())
	}
	fmt.Println()

	// Example 4: User identifier for abuse monitoring
	fmt.Println("=== Example 4: User Identifier ===")
	req4 := chat.NewRequest("grok-beta",
		chat.WithMessage(chat.User(chat.Text("Hello, how are you?"))),
		chat.WithUser("user-12345"), // Track by user ID
	)
	req4.SetMaxTokens(50)

	resp4, err := req4.Sample(ctx, client.Chat())
	if err != nil {
		log.Fatalf("User identifier example failed: %v", err)
	}
	fmt.Printf("Response: %s\n", resp4.Content())
	fmt.Printf("Note: Request tracked with user ID 'user-12345' for abuse monitoring\n\n")

	// Example 5: Combining ALL parameters
	fmt.Println("=== Example 5: All Parameters Combined ===")
	req5 := chat.NewRequest("grok-beta")
	req5.SetMessages(
		*chat.System(chat.Text("You are a helpful assistant.")),
		*chat.User(chat.Text("Explain quantum computing in one sentence.")),
	)

	// Sampling parameters
	req5.SetTemperature(0.7)
	req5.SetTopP(0.9)
	req5.SetMaxTokens(100)

	// Penalty parameters
	req5.SetFrequencyPenalty(0.3)
	req5.SetPresencePenalty(0.3)

	// Stop sequences
	req5.SetStop("However", "Additionally")

	// Determinism & logging
	req5.SetSeed(42)
	req5.SetLogprobs(true)
	req5.SetTopLogprobs(2)

	// Multiple completions
	req5.SetN(2)

	// User tracking
	req5.SetUser("demo-user")

	resp5, err := req5.Sample(ctx, client.Chat())
	if err != nil {
		log.Fatalf("Combined parameters example failed: %v", err)
	}

	fmt.Printf("Generated %d completion(s):\n", resp5.ChoiceCount())
	for i, choice := range resp5.Choices() {
		fmt.Printf("  %d. %s\n", i+1, choice.Message().Content())
	}
	fmt.Printf("\nTokens Used: %d\n", resp5.Usage().TotalTokens())
	fmt.Printf("Finish Reason: %s\n", resp5.Choices()[0].FinishReason())

	fmt.Println("\n✅ All parameter examples completed successfully!")
	fmt.Println("\n📊 Parameters demonstrated:")
	fmt.Println("  - Seed (deterministic sampling)")
	fmt.Println("  - Logprobs & TopLogprobs (confidence scoring)")
	fmt.Println("  - N (multiple completions)")
	fmt.Println("  - User (abuse monitoring)")
	fmt.Println("  - Temperature, TopP, MaxTokens")
	fmt.Println("  - FrequencyPenalty, PresencePenalty")
	fmt.Println("  - Stop sequences")
}
