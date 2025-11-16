//go:build integration
// +build integration

package image_test

import (
	"context"
	"os"
	"testing"

	xaiv1 "github.com/ZaguanLabs/xai-sdk-go/proto/gen/go/xai/v1"
	"github.com/ZaguanLabs/xai-sdk-go/xai"
	"github.com/ZaguanLabs/xai-sdk-go/xai/image"
)

func TestImageIntegration(t *testing.T) {
	apiKey := os.Getenv("XAI_API_KEY")
	if apiKey == "" {
		t.Skip("XAI_API_KEY not set, skipping integration test")
	}

	client, err := xai.NewClient(&xai.Config{APIKey: apiKey})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	ctx := context.Background()

	t.Run("GenerateImage", func(t *testing.T) {
		req := image.NewRequest(
			"A serene mountain landscape at sunset",
			"grok-vision-beta",
		)

		resp, err := client.Images().Generate(ctx, req)
		if err != nil {
			t.Fatalf("Generate failed: %v", err)
		}

		if len(resp.Images) == 0 {
			t.Fatal("Expected at least one image")
		}

		img := resp.Images[0]
		if img.URL == "" && img.Base64 == "" {
			t.Fatal("Expected either URL or Base64 data")
		}

		t.Logf("Generated image: URL=%v, Base64=%v", img.URL != "", img.Base64 != "")
	})

	t.Run("GenerateMultipleImages", func(t *testing.T) {
		req := image.NewRequest(
			"A futuristic city",
			"grok-vision-beta",
		).WithCount(2)

		resp, err := client.Images().Generate(ctx, req)
		if err != nil {
			t.Fatalf("Generate failed: %v", err)
		}

		if len(resp.Images) != 2 {
			t.Fatalf("Expected 2 images, got %d", len(resp.Images))
		}

		t.Logf("Generated %d images", len(resp.Images))
	})

	t.Run("GenerateBase64", func(t *testing.T) {
		req := image.NewRequest(
			"A colorful abstract pattern",
			"grok-vision-beta",
		).WithFormat(xaiv1.ImageFormat_IMG_FORMAT_BASE64)

		resp, err := client.Images().Generate(ctx, req)
		if err != nil {
			t.Fatalf("Generate failed: %v", err)
		}

		if len(resp.Images) == 0 {
			t.Fatal("Expected at least one image")
		}

		if resp.Images[0].Base64 == "" {
			t.Fatal("Expected Base64 data")
		}

		t.Logf("Generated Base64 image (%d bytes)", len(resp.Images[0].Base64))
	})
}
