package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/ZaguanLabs/xai-sdk-go/xai"
	"github.com/ZaguanLabs/xai-sdk-go/xai/collections"
)

func main() {
	client, err := xai.NewClient(&xai.Config{
		APIKey: os.Getenv("XAI_API_KEY"),
	})
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	teamID := os.Getenv("XAI_TEAM_ID")
	if teamID == "" {
		log.Fatal("XAI_TEAM_ID environment variable required")
	}

	// Example 1: Create a collection
	fmt.Println("=== Example 1: Create Collection ===")
	collection, err := client.Collections().CreateCollection(context.Background(), collections.CreateCollectionOptions{
		Name:   "My Document Collection",
		TeamID: teamID,
	})
	if err != nil {
		log.Fatalf("Create collection failed: %v", err)
	}

	fmt.Printf("Created collection: %s\n", collection.ID)
	fmt.Printf("Name: %s\n", collection.Name)
	fmt.Printf("Created at: %s\n", collection.CreatedAt)

	// Example 2: List collections
	fmt.Println("\n=== Example 2: List Collections ===")
	cols, _, err := client.Collections().ListCollections(context.Background(), &collections.ListCollectionsOptions{
		TeamID: teamID,
		Limit:  10,
	})
	if err != nil {
		log.Fatalf("List collections failed: %v", err)
	}

	fmt.Printf("Found %d collection(s)\n", len(cols))
	for i, col := range cols {
		fmt.Printf("%d. %s (%s) - %d documents\n", i+1, col.Name, col.ID, col.DocumentsCount)
	}

	// Example 3: Get collection
	if collection != nil {
		fmt.Println("\n=== Example 3: Get Collection ===")
		col, err := client.Collections().GetCollection(context.Background(), collection.ID, teamID)
		if err != nil {
			log.Fatalf("Get collection failed: %v", err)
		}

		fmt.Printf("Collection ID: %s\n", col.ID)
		fmt.Printf("Name: %s\n", col.Name)
		fmt.Printf("Documents: %d\n", col.DocumentsCount)
	}

	// Example 4: Update collection
	if collection != nil {
		fmt.Println("\n=== Example 4: Update Collection ===")
		updated, err := client.Collections().UpdateCollection(context.Background(), collection.ID, teamID, collections.CreateCollectionOptions{
			Name: "Updated Collection Name",
		})
		if err != nil {
			log.Fatalf("Update collection failed: %v", err)
		}

		fmt.Printf("Updated collection name: %s\n", updated.Name)
	}

	// Example 5: Delete collection
	if collection != nil {
		fmt.Println("\n=== Example 5: Delete Collection ===")
		err = client.Collections().DeleteCollection(context.Background(), collection.ID, teamID)
		if err != nil {
			log.Fatalf("Delete collection failed: %v", err)
		}

		fmt.Printf("Collection %s deleted successfully\n", collection.ID)
	}

	fmt.Println("\n✅ Collection management examples completed!")
}
