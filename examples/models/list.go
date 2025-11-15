// Package main demonstrates model listing and retrieval usage.
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

	// Create a models client
	modelsClient := client.Models()

	// List all available models
	ctx := context.Background()
	models, err := modelsClient.List(ctx)
	if err != nil {
		log.Fatalf("Failed to list models: %v", err)
	}

	fmt.Printf("Available models (%d):\n", len(models))
	for _, model := range models {
		fmt.Printf("  - %s: %s (max tokens: %d)\n", model.ID(), model.Name(), model.MaxTokens())
		if model.Description() != "" {
			fmt.Printf("    %s\n", model.Description())
		}
	}

	// Get information about a specific model
	if len(models) > 0 {
		modelID := models[0].ID()
		model, err := modelsClient.Get(ctx, modelID)
		if err != nil {
			log.Fatalf("Failed to get model %s: %v", modelID, err)
		}

		fmt.Printf("\nDetailed information for model '%s':\n", modelID)
		fmt.Printf("  ID: %s\n", model.ID())
		fmt.Printf("  Name: %s\n", model.Name())
		fmt.Printf("  Description: %s\n", model.Description())
		fmt.Printf("  Max Tokens: %d\n", model.MaxTokens())
	}
}