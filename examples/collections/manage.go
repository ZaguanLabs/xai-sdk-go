// Package main demonstrates collections management usage.
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

	// Create a collections client
	collectionsClient := client.Collections()

	ctx := context.Background()

	// List all collections
	fmt.Println("=== Listing Collections ===")
	collections, err := collectionsClient.ListCollections(ctx, "", 10, "name", "asc")
	if err != nil {
		log.Fatalf("Failed to list collections: %v", err)
	}

	fmt.Printf("Found %d collections:\n", len(collections))
	for _, collection := range collections {
		fmt.Printf("  - %s: %s (%d documents, %d bytes)\n",
			collection.ID(), collection.Name(), collection.DocumentCount(), collection.TotalSize())
		if collection.Description() != "" {
			fmt.Printf("    %s\n", collection.Description())
		}
	}

	// Create a new collection
	fmt.Println("\n=== Creating Collection ===")
	newCollection, err := collectionsClient.CreateCollection(ctx, "My Documents", "A collection of important documents", false)
	if err != nil {
		log.Fatalf("Failed to create collection: %v", err)
	}

	fmt.Printf("Created collection: %s (%s)\n", newCollection.Name(), newCollection.ID())

	// Add documents to the collection
	fmt.Println("\n=== Adding Documents ===")
	doc1, err := collectionsClient.AddDocument(ctx, newCollection.ID(), "Introduction", "This is the introduction to my document collection.", "text/plain")
	if err != nil {
		log.Fatalf("Failed to add document 1: %v", err)
	}
	fmt.Printf("Added document: %s (%s)\n", doc1.Title(), doc1.ID())

	doc2, err := collectionsClient.AddDocument(ctx, newCollection.ID(), "Summary", "This is a summary of all the documents in the collection.", "text/plain")
	if err != nil {
		log.Fatalf("Failed to add document 2: %v", err)
	}
	fmt.Printf("Added document: %s (%s)\n", doc2.Title(), doc2.ID())

	// List documents in the collection
	fmt.Println("\n=== Listing Documents ===")
	documents, err := collectionsClient.ListDocuments(ctx, newCollection.ID(), 10, "title", "asc")
	if err != nil {
		log.Fatalf("Failed to list documents: %v", err)
	}

	fmt.Printf("Collection contains %d documents:\n", len(documents))
	for _, doc := range documents {
		fmt.Printf("  - %s: %s (%d bytes, %s)\n",
			doc.Title(), doc.ID(), doc.Size(), doc.ContentType())
	}

	// Get a specific document
	fmt.Println("\n=== Getting Document ===")
	document, err := collectionsClient.GetDocument(ctx, newCollection.ID(), doc1.ID())
	if err != nil {
		log.Fatalf("Failed to get document: %v", err)
	}

	fmt.Printf("Retrieved document: %s\n", document.Title())
	fmt.Printf("  Content: %s...\n", truncateString(document.Content(), 50))

	// Clean up - delete the collection
	fmt.Println("\n=== Cleaning Up ===")
	err = collectionsClient.DeleteCollection(ctx, newCollection.ID())
	if err != nil {
		log.Fatalf("Failed to delete collection: %v", err)
	}
	fmt.Printf("Deleted collection: %s\n", newCollection.ID())

	fmt.Println("\n✓ Collections demo completed successfully!")
}

// truncateString truncates a string to the specified length
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}