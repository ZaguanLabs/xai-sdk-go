package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/ZaguanLabs/xai-sdk-go/xai"
	"github.com/ZaguanLabs/xai-sdk-go/xai/documents"
)

func main() {
	// Create client
	client, err := xai.NewClient(&xai.Config{
		APIKey: os.Getenv("XAI_API_KEY"),
	})
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	// Example 1: Search across collections
	fmt.Println("=== Example 1: Basic Document Search ===")
	req := documents.NewSearchRequest(
		"machine learning algorithms",
		"collection-id-1",
		"collection-id-2",
	).WithLimit(5)

	resp, err := client.Documents().Search(context.Background(), req)
	if err != nil {
		log.Fatalf("Document search failed: %v", err)
	}

	fmt.Printf("Found %d matches\n", len(resp.Matches))
	for i, match := range resp.Matches {
		fmt.Printf("\nMatch %d:\n", i+1)
		fmt.Printf("  File ID: %s\n", match.FileID)
		fmt.Printf("  Chunk ID: %s\n", match.ChunkID)
		fmt.Printf("  Score: %.4f\n", match.Score)
		fmt.Printf("  Content: %s\n", truncate(match.ChunkContent, 100))
		fmt.Printf("  Collections: %v\n", match.CollectionIDs)
	}

	// Example 2: Search with higher limit
	fmt.Println("\n=== Example 2: Extended Search ===")
	req2 := documents.NewSearchRequest(
		"neural networks deep learning",
		"collection-id-1",
	).WithLimit(10)

	resp2, err := client.Documents().Search(context.Background(), req2)
	if err != nil {
		log.Fatalf("Extended search failed: %v", err)
	}

	fmt.Printf("Found %d matches with extended limit\n", len(resp2.Matches))

	// Show top 3 results
	for i := 0; i < min(3, len(resp2.Matches)); i++ {
		match := resp2.Matches[i]
		fmt.Printf("\nTop %d (Score: %.4f):\n", i+1, match.Score)
		fmt.Printf("  %s\n", truncate(match.ChunkContent, 150))
	}

	fmt.Println("\n✅ Document search examples completed!")
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
