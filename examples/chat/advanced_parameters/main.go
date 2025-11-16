// Package main demonstrates advanced chat parameters including TopP, Stop, FrequencyPenalty, and PresencePenalty.
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

	// Example 1: Using TopP for nucleus sampling
	fmt.Println("=== Example 1: TopP (Nucleus Sampling) ===")
	req1 := chat.NewRequest("grok-beta",
		chat.WithMessage(chat.System(chat.Text("You are a creative writer."))),
		chat.WithMessage(chat.User(chat.Text("Write a short creative sentence about space."))),
		chat.WithTopP(0.9), // Consider tokens with top 90% probability mass
	)
	req1.SetMaxTokens(50)

	resp1, err := req1.Sample(ctx, client.Chat())
	if err != nil {
		log.Fatalf("TopP example failed: %v", err)
	}
	fmt.Printf("Response: %s\n\n", resp1.Content())

	// Example 2: Using Stop sequences
	fmt.Println("=== Example 2: Stop Sequences ===")
	req2 := chat.NewRequest("grok-beta",
		chat.WithMessage(chat.User(chat.Text("Count from 1 to 10: 1, 2, 3,"))),
		chat.WithStop("5", "10"), // Stop at 5 or 10
	)
	req2.SetMaxTokens(100)

	resp2, err := req2.Sample(ctx, client.Chat())
	if err != nil {
		log.Fatalf("Stop example failed: %v", err)
	}
	fmt.Printf("Response: %s\n", resp2.Content())
	fmt.Printf("Finish Reason: %s\n\n", resp2.FinishReason())

	// Example 3: Using FrequencyPenalty to reduce repetition
	fmt.Println("=== Example 3: Frequency Penalty (Reduce Repetition) ===")
	req3 := chat.NewRequest("grok-beta",
		chat.WithMessage(chat.User(chat.Text("Write a sentence using the word 'amazing' multiple times."))),
		chat.WithFrequencyPenalty(1.5), // Penalize repeated tokens
	)
	req3.SetMaxTokens(100)

	resp3, err := req3.Sample(ctx, client.Chat())
	if err != nil {
		log.Fatalf("Frequency penalty example failed: %v", err)
	}
	fmt.Printf("Response: %s\n\n", resp3.Content())

	// Example 4: Using PresencePenalty to encourage new topics
	fmt.Println("=== Example 4: Presence Penalty (Encourage New Topics) ===")
	req4 := chat.NewRequest("grok-beta",
		chat.WithMessage(chat.System(chat.Text("You are a helpful assistant."))),
		chat.WithMessage(chat.User(chat.Text("Tell me about cats."))),
		chat.WithPresencePenalty(1.0), // Encourage talking about new aspects
	)
	req4.SetMaxTokens(150)

	resp4, err := req4.Sample(ctx, client.Chat())
	if err != nil {
		log.Fatalf("Presence penalty example failed: %v", err)
	}
	fmt.Printf("Response: %s\n\n", resp4.Content())

	// Example 5: Combining multiple parameters
	fmt.Println("=== Example 5: Combined Parameters ===")
	req5 := chat.NewRequest("grok-beta")
	req5.SetMessages(
		*chat.System(chat.Text("You are a concise technical writer.")),
		*chat.User(chat.Text("Explain what an API is.")),
	)
	req5.SetTemperature(0.7)
	req5.SetTopP(0.95)
	req5.SetMaxTokens(100)
	req5.SetFrequencyPenalty(0.5)
	req5.SetPresencePenalty(0.3)
	req5.SetStop("However", "Additionally") // Stop at these transition words

	resp5, err := req5.Sample(ctx, client.Chat())
	if err != nil {
		log.Fatalf("Combined parameters example failed: %v", err)
	}
	fmt.Printf("Response: %s\n", resp5.Content())
	fmt.Printf("Tokens Used: %d\n", resp5.Usage().TotalTokens())

	fmt.Println("\n✅ All advanced parameter examples completed successfully!")
}
