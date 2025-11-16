package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/ZaguanLabs/xai-sdk-go/xai"
)

func main() {
	client, err := xai.NewClient(&xai.Config{
		APIKey: os.Getenv("XAI_API_KEY"),
	})
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	// Example 1: Validate API key
	fmt.Println("=== Example 1: Validate API Key ===")
	apiKey := os.Getenv("XAI_API_KEY")
	key, err := client.Auth().ValidateKey(context.Background(), apiKey)
	if err != nil {
		log.Printf("Validate key failed: %v", err)
	} else {
		fmt.Printf("API Key validated: %s\n", key.RedactedApiKey)
		fmt.Printf("User ID: %s\n", key.UserID)
		if key.Name != "" {
			fmt.Printf("Name: %s\n", key.Name)
		}
		if key.TeamID != "" {
			fmt.Printf("Team ID: %s\n", key.TeamID)
		}
	}

	// Example 2: List API keys
	fmt.Println("\n=== Example 2: List API Keys ===")
	keys, err := client.Auth().ListKeys(context.Background())
	if err != nil {
		log.Printf("List keys failed: %v", err)
	} else {
		fmt.Printf("Found %d API key(s)\n", len(keys))
		for i, k := range keys {
			fmt.Printf("%d. %s (User: %s)\n", i+1, k.RedactedApiKey, k.UserID)
			if k.Name != "" {
				fmt.Printf("   Name: %s\n", k.Name)
			}
		}
	}

	// Example 3: Get specific API key
	fmt.Println("\n=== Example 3: Get API Key by ID ===")
	// Note: You need to provide a valid key ID
	keyID := os.Getenv("XAI_KEY_ID")
	if keyID != "" {
		key, err := client.Auth().GetKey(context.Background(), keyID)
		if err != nil {
			log.Printf("Get key failed: %v", err)
		} else {
			fmt.Printf("API Key: %s\n", key.RedactedApiKey)
			fmt.Printf("User ID: %s\n", key.UserID)
			fmt.Printf("Created: %s\n", key.CreateTime)
			if !key.ModifyTime.IsZero() {
				fmt.Printf("Modified: %s\n", key.ModifyTime)
			}
		}
	} else {
		fmt.Println("Set XAI_KEY_ID environment variable to test GetKey")
	}

	fmt.Println("\n✅ Auth API examples completed!")
}
