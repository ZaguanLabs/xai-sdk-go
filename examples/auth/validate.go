// Package main demonstrates API key validation usage.
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/ZaguanLabs/xai-sdk-go/xai"
)

func main() {
	// Get API key from environment variable
	apiKey := os.Getenv("XAI_API_KEY")
	if apiKey == "" {
		log.Fatal("XAI_API_KEY environment variable is required")
	}

	// Create a new client with the API key
	client, err := xai.NewClientWithAPIKey(apiKey)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	// Create an auth client
	authClient := client.Auth()

	// Validate the API key
	ctx := context.Background()
	result, err := authClient.Validate(ctx, apiKey)
	if err != nil {
		log.Fatalf("Failed to validate API key: %v", err)
	}

	fmt.Printf("API Key Validation Result:\n")
	fmt.Printf("  Valid: %t\n", result.IsValid())
	fmt.Printf("  Message: %s\n", result.Message())

	if result.IsValid() {
		fmt.Printf("  Organization: %s\n", result.Organization())
		fmt.Printf("  Project: %s\n", result.Project())
		fmt.Println("✓ API key is valid!")
	} else {
		fmt.Println("✗ API key is invalid!")
	}
}