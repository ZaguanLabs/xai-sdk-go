package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/ZaguanLabs/xai-sdk-go/xai"
	"github.com/ZaguanLabs/xai-sdk-go/xai/files"
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

	// Example 1: Upload a text file
	fmt.Println("=== Example 1: Upload Text File ===")
	content := strings.NewReader("This is a sample document for testing.\nIt contains multiple lines of text.")

	file, err := client.Files().Upload(context.Background(), content, files.UploadOptions{
		Name:    "sample.txt",
		Purpose: "assistants",
	})
	if err != nil {
		log.Fatalf("File upload failed: %v", err)
	}

	fmt.Printf("Uploaded file: %s\n", file.ID)
	fmt.Printf("Filename: %s\n", file.Filename)
	fmt.Printf("Size: %d bytes\n", file.Size)
	fmt.Printf("Created: %s\n", file.CreatedAt)

	// Example 2: List all files
	fmt.Println("\n=== Example 2: List Files ===")
	listResult, err := client.Files().List(context.Background(), &files.ListOptions{
		Limit: 10,
	})
	if err != nil {
		log.Fatalf("List files failed: %v", err)
	}

	fmt.Printf("Found %d file(s)\n", len(listResult.Files))
	for i, f := range listResult.Files {
		fmt.Printf("%d. %s (%s, %d bytes)\n", i+1, f.Filename, f.ID, f.Size)
	}

	// Example 3: Get file metadata
	if file != nil {
		fmt.Println("\n=== Example 3: Get File Metadata ===")
		metadata, err := client.Files().Get(context.Background(), file.ID)
		if err != nil {
			log.Fatalf("Get file failed: %v", err)
		}

		fmt.Printf("File ID: %s\n", metadata.ID)
		fmt.Printf("Filename: %s\n", metadata.Filename)
		fmt.Printf("Size: %d bytes\n", metadata.Size)
		fmt.Printf("Created: %s\n", metadata.CreatedAt)
		if !metadata.ExpiresAt.IsZero() {
			fmt.Printf("Expires: %s\n", metadata.ExpiresAt)
		}
	}

	// Example 4: Get download URL
	if file != nil {
		fmt.Println("\n=== Example 4: Get Download URL ===")
		url, err := client.Files().GetURL(context.Background(), file.ID)
		if err != nil {
			log.Fatalf("Get URL failed: %v", err)
		}

		fmt.Printf("Download URL: %s\n", url)
	}

	// Example 5: Delete file
	if file != nil {
		fmt.Println("\n=== Example 5: Delete File ===")
		err = client.Files().Delete(context.Background(), file.ID)
		if err != nil {
			log.Fatalf("Delete file failed: %v", err)
		}

		fmt.Printf("File %s deleted successfully\n", file.ID)
	}

	fmt.Println("\n✅ File operations completed!")
}
