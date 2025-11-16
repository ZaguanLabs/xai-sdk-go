package main

import (
	"context"
	"fmt"
	"log"
	"os"

	xaiv1 "github.com/ZaguanLabs/xai-sdk-go/proto/gen/go/xai/v1"
	"github.com/ZaguanLabs/xai-sdk-go/xai"
	"github.com/ZaguanLabs/xai-sdk-go/xai/image"
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

	// Example 1: Simple text-to-image
	fmt.Println("=== Example 1: Text-to-Image ===")
	req := image.NewRequest(
		"A serene mountain landscape at sunset with a lake reflection",
		"grok-vision-beta",
	)

	resp, err := client.Images().Generate(context.Background(), req)
	if err != nil {
		log.Fatalf("Image generation failed: %v", err)
	}

	fmt.Printf("Generated %d image(s)\n", len(resp.Images))
	for i, img := range resp.Images {
		if img.URL != "" {
			fmt.Printf("Image %d URL: %s\n", i+1, img.URL)
		}
		if img.UpsampledPrompt != "" {
			fmt.Printf("Upsampled prompt: %s\n", img.UpsampledPrompt)
		}
	}

	// Example 2: Multiple images with Base64 format
	fmt.Println("\n=== Example 2: Multiple Images (Base64) ===")
	req2 := image.NewRequest(
		"A futuristic city with flying cars",
		"grok-vision-beta",
	).WithCount(2).WithFormat(xaiv1.ImageFormat_IMG_FORMAT_BASE64)

	resp2, err := client.Images().Generate(context.Background(), req2)
	if err != nil {
		log.Fatalf("Image generation failed: %v", err)
	}

	fmt.Printf("Generated %d image(s) in Base64 format\n", len(resp2.Images))
	for i, img := range resp2.Images {
		if img.Base64 != "" {
			fmt.Printf("Image %d: Base64 data (%d bytes)\n", i+1, len(img.Base64))
		}
	}

	// Example 3: Image-to-image with input
	fmt.Println("\n=== Example 3: Image-to-Image ===")
	req3 := image.NewRequest(
		"Make this image more vibrant and colorful",
		"grok-vision-beta",
	).WithImage(
		"https://example.com/input-image.jpg",
		xaiv1.ImageDetail_DETAIL_HIGH,
	)

	resp3, err := client.Images().Generate(context.Background(), req3)
	if err != nil {
		log.Printf("Image-to-image generation failed (expected if URL invalid): %v", err)
	} else {
		fmt.Printf("Generated image URL: %s\n", resp3.Images[0].URL)
	}

	fmt.Println("\n✅ Image generation examples completed!")
}
