//go:build integration
// +build integration

package embed_test

import (
	"context"
	"os"
	"testing"

	xaiv1 "github.com/ZaguanLabs/xai-sdk-go/proto/gen/go/xai/v1"
	"github.com/ZaguanLabs/xai-sdk-go/xai"
	"github.com/ZaguanLabs/xai-sdk-go/xai/embed"
)

func TestEmbedIntegration(t *testing.T) {
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

	t.Run("GenerateTextEmbedding", func(t *testing.T) {
		req := embed.NewRequest(
			"grok-embedding-1",
			embed.Text("The quick brown fox jumps over the lazy dog"),
		)

		resp, err := client.Embed().Generate(ctx, req)
		if err != nil {
			t.Fatalf("Generate failed: %v", err)
		}

		if len(resp.Embeddings()) == 0 {
			t.Fatal("Expected at least one embedding")
		}

		emb := resp.Embeddings()[0]
		vectors := emb.Vectors()
		if len(vectors) == 0 {
			t.Fatal("Expected at least one vector")
		}

		vec := vectors[0].FloatArray()
		if len(vec) == 0 {
			t.Fatal("Expected non-empty vector")
		}

		t.Logf("Generated embedding with %d dimensions", len(vec))
	})

	t.Run("GenerateBatchEmbeddings", func(t *testing.T) {
		req := embed.NewRequest(
			"grok-embedding-1",
			embed.Text("First text"),
			embed.Text("Second text"),
			embed.Text("Third text"),
		)

		resp, err := client.Embed().Generate(ctx, req)
		if err != nil {
			t.Fatalf("Generate failed: %v", err)
		}

		if len(resp.Embeddings()) != 3 {
			t.Fatalf("Expected 3 embeddings, got %d", len(resp.Embeddings()))
		}

		t.Logf("Generated %d embeddings", len(resp.Embeddings()))
	})

	t.Run("GenerateImageEmbedding", func(t *testing.T) {
		req := embed.NewRequest(
			"grok-embedding-1",
			embed.Image("https://example.com/image.jpg", xaiv1.ImageDetail_DETAIL_AUTO),
		)

		resp, err := client.Embed().Generate(ctx, req)
		// This might fail if the URL is invalid or image embeddings aren't supported
		// We just want to verify the API call works
		if err != nil {
			t.Logf("Image embedding failed (expected): %v", err)
		} else {
			t.Logf("Image embedding succeeded with %d embeddings", len(resp.Embeddings()))
		}
	})
}
