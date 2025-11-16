// Package main demonstrates image understanding with vision models.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/ZaguanLabs/xai-sdk-go/xai"
	"github.com/ZaguanLabs/xai-sdk-go/xai/chat"
)

func main() {
	// Get API key from environment
	apiKey := os.Getenv("XAI_API_KEY")
	if apiKey == "" {
		log.Fatal("XAI_API_KEY environment variable is required")
	}

	// Create client
	client, err := xai.NewClient(&xai.Config{APIKey: apiKey})
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	ctx := context.Background()

	// Example 1: Single image analysis
	fmt.Println("=== Example 1: Single Image Analysis ===")
	singleImageExample(ctx, client)

	fmt.Println("\n" + strings.Repeat("=", 60) + "\n")

	// Example 2: Multiple images comparison
	fmt.Println("=== Example 2: Multiple Images Comparison ===")
	multipleImagesExample(ctx, client)

	fmt.Println("\n" + strings.Repeat("=", 60) + "\n")

	// Example 3: Image with different detail levels
	fmt.Println("=== Example 3: Detail Levels ===")
	detailLevelsExample(ctx, client)

	fmt.Println("\n" + strings.Repeat("=", 60) + "\n")

	// Example 4: Mixed content (text and images)
	fmt.Println("=== Example 4: Mixed Content ===")
	mixedContentExample(ctx, client)
}

func singleImageExample(ctx context.Context, client *xai.Client) {
	req := chat.NewRequest("grok-2-vision",
		chat.WithMessage(
			chat.User(
				chat.Text("What's in this image? Describe it in detail."),
				chat.Image("https://upload.wikimedia.org/wikipedia/commons/a/a7/Camponotus_flavomarginatus_ant.jpg"),
			),
		),
	)

	response, err := req.Sample(ctx, client.Chat())
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	fmt.Printf("Model: %s\n", response.Model())
	fmt.Printf("Response: %s\n", response.Content())
}

func multipleImagesExample(ctx context.Context, client *xai.Client) {
	req := chat.NewRequest("grok-2-vision",
		chat.WithMessage(
			chat.User(
				chat.Text("Compare these two images. What are the similarities and differences?"),
				chat.Image("https://upload.wikimedia.org/wikipedia/commons/a/a7/Camponotus_flavomarginatus_ant.jpg"),
				chat.Image("https://upload.wikimedia.org/wikipedia/commons/9/9f/Atta_cephalotes-pjt.jpg"),
			),
		),
	)

	response, err := req.Sample(ctx, client.Chat())
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	fmt.Printf("Response: %s\n", response.Content())
}

func detailLevelsExample(ctx context.Context, client *xai.Client) {
	// Low detail - faster, uses fewer tokens
	fmt.Println("--- Low Detail (faster, cheaper) ---")
	reqLow := chat.NewRequest("grok-2-vision",
		chat.WithMessage(
			chat.User(
				chat.Text("What's the main subject of this image?"),
				chat.Image("https://upload.wikimedia.org/wikipedia/commons/a/a7/Camponotus_flavomarginatus_ant.jpg", chat.ImageDetailLow),
			),
		),
	)

	responseLow, err := reqLow.Sample(ctx, client.Chat())
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	fmt.Printf("Response: %s\n", responseLow.Content())

	fmt.Println("\n--- High Detail (slower, more detailed) ---")
	// High detail - slower, uses more tokens, captures more detail
	reqHigh := chat.NewRequest("grok-2-vision",
		chat.WithMessage(
			chat.User(
				chat.Text("Describe this image in great detail, including colors, textures, and fine features."),
				chat.Image("https://upload.wikimedia.org/wikipedia/commons/a/a7/Camponotus_flavomarginatus_ant.jpg", chat.ImageDetailHigh),
			),
		),
	)

	responseHigh, err := reqHigh.Sample(ctx, client.Chat())
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	fmt.Printf("Response: %s\n", responseHigh.Content())
}

func mixedContentExample(ctx context.Context, client *xai.Client) {
	req := chat.NewRequest("grok-2-vision",
		chat.WithMessage(
			chat.User(
				chat.Text("I have a question about this insect:"),
				chat.Image("https://upload.wikimedia.org/wikipedia/commons/a/a7/Camponotus_flavomarginatus_ant.jpg", chat.ImageDetailHigh),
				chat.Text("Is this a carpenter ant? What are its distinguishing features?"),
			),
		),
	)

	response, err := req.Sample(ctx, client.Chat())
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	fmt.Printf("Response: %s\n", response.Content())
}
