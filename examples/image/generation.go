// Package main demonstrates image generation usage.
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/ZaguanLabs/xai-sdk-go/xai"
	"github.com/ZaguanLabs/xai-sdk-go/xai/image"
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

	// Create an image generation request
	req := image.NewGenerateRequest("A beautiful sunset over mountains with a lake in the foreground", "dall-e-3").
		WithSize("1024x1024").
		WithQuality("hd").
		WithStyle("vivid").
		WithN(1)

	// Generate the image
	ctx := context.Background()
	images, err := client.Images().Generate(ctx, req)
	if err != nil {
		log.Fatalf("Failed to generate image: %v", err)
	}

	fmt.Printf("Generated %d image(s):\n", len(images))
	for i, img := range images {
		fmt.Printf("\nImage %d:\n", i+1)
		fmt.Printf("  URL: %s\n", img.URL())
		if img.RevisedPrompt() != "" {
			fmt.Printf("  Revised Prompt: %s\n", img.RevisedPrompt())
		}

		// Save the image to a file
		fileName := fmt.Sprintf("generated_image_%d.png", i+1)
		if err := img.Save(fileName, nil); err != nil { // Pass nil for default HTTP client
			fmt.Printf("  Warning: Failed to save image: %v\n", err)
		} else {
			fmt.Printf("  Saved to: %s\n", fileName)
		}
	}
}