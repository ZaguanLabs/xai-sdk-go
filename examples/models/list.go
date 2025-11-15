// Package main demonstrates model listing and retrieval usage.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/ZaguanLabs/xai-sdk-go/xai"
	"google.golang.org/grpc/metadata"
)

func main() {
	log.SetFlags(log.Ltime | log.Lmicroseconds)

	// Get API key from environment variable
	apiKey := os.Getenv("XAI_API_KEY")
	if apiKey == "" {
		log.Fatal("XAI_API_KEY environment variable is required")
	}

	// Mask API key for logging (show first 10 chars)
	maskedKey := apiKey
	if len(apiKey) > 10 {
		maskedKey = apiKey[:10] + strings.Repeat("*", len(apiKey)-10)
	}
	log.Printf("[DEBUG] Using API key: %s", maskedKey)

	// Create a new client with the API key
	log.Printf("[DEBUG] Creating xAI client...")
	client, err := xai.NewClientWithAPIKey(apiKey)
	if err != nil {
		log.Fatalf("[ERROR] Failed to create client: %v", err)
	}
	defer func() {
		log.Printf("[DEBUG] Closing client...")
		client.Close()
	}()

	log.Printf("[DEBUG] Client created successfully")
	log.Printf("[DEBUG] Client config: %s", client.Config())
	log.Printf("[DEBUG] Client health: %+v", client.GetHealthStatus())

	// Create a models client
	log.Printf("[DEBUG] Creating models client...")
	modelsClient := client.Models()
	log.Printf("[DEBUG] Models client created")

	// Create a new context with the client's configuration and metadata
	log.Printf("[DEBUG] Creating context with client metadata...")
	ctx := client.NewContext(context.Background())

	// Log outgoing metadata
	if md, ok := metadata.FromOutgoingContext(ctx); ok {
		log.Printf("[DEBUG] Outgoing metadata keys: %v", md)
		for key, values := range md {
			// Don't log full API key
			if strings.Contains(strings.ToLower(key), "key") || strings.Contains(strings.ToLower(key), "auth") {
				log.Printf("[DEBUG]   %s: [REDACTED]", key)
			} else {
				log.Printf("[DEBUG]   %s: %v", key, values)
			}
		}
	} else {
		log.Printf("[DEBUG] No outgoing metadata found in context")
	}

	// List all available language models
	log.Printf("[DEBUG] Calling modelsClient.ListLanguageModels()...")
	models, err := modelsClient.ListLanguageModels(ctx)
	if err != nil {
		log.Printf("[ERROR] Failed to list language models: %v", err)
		log.Printf("[ERROR] Error type: %T", err)
		log.Fatalf("[ERROR] Exiting due to error")
	}

	log.Printf("[DEBUG] Successfully retrieved %d language models", len(models))

	fmt.Printf("\nAvailable Language Models (%d):\n", len(models))
	for _, model := range models {
		fmt.Printf("  - %s (v%s)\n", model.Name(), model.Version())
		if len(model.Aliases()) > 0 {
			fmt.Printf("    Aliases: %v\n", model.Aliases())
		}
		fmt.Printf("    Max Prompt Length: %d\n", model.MaxPromptLength())
		fmt.Printf("    Input Modalities: %v\n", model.InputModalities())
		fmt.Printf("    Output Modalities: %v\n", model.OutputModalities())
		if model.SystemFingerprint() != "" {
			fmt.Printf("    System Fingerprint: %s\n", model.SystemFingerprint())
		}
		fmt.Println()
	}

	// Get information about a specific model
	if len(models) > 0 {
		modelName := models[0].Name()
		log.Printf("[DEBUG] Fetching detailed info for model: %s", modelName)
		model, err := modelsClient.GetLanguageModel(ctx, modelName)
		if err != nil {
			log.Printf("[WARN] Failed to get model %s: %v", modelName, err)
		} else {
			log.Printf("[DEBUG] Successfully retrieved model details for: %s", modelName)
			fmt.Printf("\nDetailed information for model '%s':\n", modelName)
			fmt.Printf("  Name: %s\n", model.Name())
			fmt.Printf("  Version: %s\n", model.Version())
			fmt.Printf("  Aliases: %v\n", model.Aliases())
			fmt.Printf("  Max Prompt Length: %d\n", model.MaxPromptLength())
			fmt.Printf("  Prompt Text Token Price: %d\n", model.PromptTextTokenPrice())
			fmt.Printf("  Completion Text Token Price: %d\n", model.CompletionTextTokenPrice())
			fmt.Printf("  System Fingerprint: %s\n", model.SystemFingerprint())
		}
	}

	fmt.Printf("\nNote: Use these model names when making chat completion requests.\n")

	log.Printf("[DEBUG] Program completed successfully")
}
