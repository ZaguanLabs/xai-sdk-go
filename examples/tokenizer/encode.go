// Package main demonstrates tokenizer usage for encoding, decoding, and counting tokens.
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

	// Create a tokenizer client
	tokenizerClient := client.Tokenizer()

	// Text to tokenize
	text := "The quick brown fox jumps over the lazy dog."
	model := "gpt-4"

	ctx := context.Background()

	// Count tokens
	tokenCount, err := tokenizerClient.Count(ctx, text, model)
	if err != nil {
		log.Fatalf("Failed to count tokens: %v", err)
	}

	fmt.Printf("Text: %s\n", text)
	fmt.Printf("Token count: %d\n", tokenCount)

	// Count tokens with details
	tokenCount, charCount, err := tokenizerClient.CountWithDetails(ctx, text, model)
	if err != nil {
		log.Fatalf("Failed to count tokens with details: %v", err)
	}

	fmt.Printf("Token count (detailed): %d\n", tokenCount)
	fmt.Printf("Character count: %d\n", charCount)

	// Encode text to tokens
	tokens, err := tokenizerClient.Encode(ctx, text, model)
	if err != nil {
		log.Fatalf("Failed to encode text: %v", err)
	}

	fmt.Printf("Encoded tokens (%d): %v\n", len(tokens), tokens)

	// Decode tokens back to text
	decodedText, err := tokenizerClient.Decode(ctx, tokens, model)
	if err != nil {
		log.Fatalf("Failed to decode tokens: %v", err)
	}

	fmt.Printf("Decoded text: %s\n", decodedText)

	// Verify round-trip encoding/decoding
	if decodedText == text {
		fmt.Println("✓ Round-trip encoding/decoding successful!")
	} else {
		fmt.Printf("✗ Round-trip failed. Original: '%s', Decoded: '%s'\n", text, decodedText)
	}
}