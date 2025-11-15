// Package main demonstrates streaming chat completion usage.
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
		chat.WithMessages(
			chat.System(chat.Text("You are a helpful assistant.")),
			chat.User(chat.Text("Tell me a story about a brave knight in 100 words.")),
		),
	)

	// Perform the streaming chat completion
	ctx := context.Background()
	stream, err := req.Stream(ctx, client.Chat())
	if err != nil {
		log.Fatalf("Failed to get streaming chat completion: %v", err)
	}
	defer stream.Close()

	// Print the response as it streams
	fmt.Println("Streaming response:")
	for stream.Next() {
		chunk := stream.Current()
		if content := chunk.Content(); content != "" {
			fmt.Print(content)
		}
	}

	// Check for any errors during streaming
	if stream.Err() != nil {
		log.Fatalf("Error during streaming: %v", stream.Err())
	}

	fmt.Println("\nStreaming completed.")
}