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

	// Example 1: Tokenize simple text
	fmt.Println("=== Example 1: Simple Tokenization ===")
	text1 := "Hello, world! This is a test."
	resp1, err := client.Tokenizer().Tokenize(context.Background(), text1, "grok-1", "")
	if err != nil {
		log.Fatalf("Tokenization failed: %v", err)
	}

	fmt.Printf("Text: %s\n", text1)
	fmt.Printf("Token count: %d\n", len(resp1.Tokens))
	fmt.Println("Tokens:")
	for i, token := range resp1.Tokens {
		fmt.Printf("  %d. ID=%d, Text=%q\n", i+1, token.TokenID, token.StringToken)
	}

	// Example 2: Tokenize longer text
	fmt.Println("\n=== Example 2: Longer Text ===")
	text2 := `Machine learning is a subset of artificial intelligence that focuses on 
	developing algorithms and statistical models that enable computers to improve their 
	performance on tasks through experience.`

	resp2, err := client.Tokenizer().Tokenize(context.Background(), text2, "grok-1", "")
	if err != nil {
		log.Fatalf("Tokenization failed: %v", err)
	}

	fmt.Printf("Text length: %d characters\n", len(text2))
	fmt.Printf("Token count: %d\n", len(resp2.Tokens))
	fmt.Printf("Tokens per character: %.2f\n", float64(len(resp2.Tokens))/float64(len(text2)))

	// Example 3: Compare different texts
	fmt.Println("\n=== Example 3: Token Count Comparison ===")
	texts := []string{
		"Short text",
		"This is a medium length text with more words",
		"This is a significantly longer text that contains many more words and should result in a higher token count when tokenized by the model",
	}

	for i, text := range texts {
		resp, err := client.Tokenizer().Tokenize(context.Background(), text, "grok-1", "")
		if err != nil {
			log.Printf("Tokenization %d failed: %v", i+1, err)
			continue
		}
		fmt.Printf("%d. %d chars → %d tokens: %q\n", i+1, len(text), len(resp.Tokens), truncate(text, 50))
	}

	fmt.Println("\n✅ Tokenizer examples completed!")
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
