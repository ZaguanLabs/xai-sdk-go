// Package main demonstrates basic chat completion usage.
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
	// Get API key from environment variable
	apiKey := os.Getenv("XAI_API_KEY")
	if apiKey == "" {
		log.Fatal("XAI_API_KEY environment variable is required")
	}

	// Create a new client with the API key
	client, err := xai.NewClient(&xai.Config{APIKey: apiKey})
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	// Create a chat request
	req := chat.NewRequest("grok-1.5-flash",
		chat.WithTemperature(0.7),
		chat.WithMaxTokens(1000),
		chat.WithMessages(
			chat.System(chat.Text("You are a helpful assistant.")),
			chat.User(chat.Text("What is the capital of France?")),
		),
	)

	// Perform the chat completion
	ctx := context.Background()
	resp, err := req.Sample(ctx, client.Chat())
	if err != nil {
		log.Fatalf("Failed to get chat completion: %v", err)
	}

	// Print the response
	fmt.Printf("Response: %s\n", resp.Content())
	fmt.Printf("Role: %s\n", resp.Role())
	
	// Print token usage if available
	if usage := resp.Usage(); usage != nil {
		fmt.Printf("Prompt tokens: %d\n", usage.PromptTokens())
		fmt.Printf("Completion tokens: %d\n", usage.CompletionTokens())
		fmt.Printf("Total tokens: %d\n", usage.TotalTokens())
	}
}