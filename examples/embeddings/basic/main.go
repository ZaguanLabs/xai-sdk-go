package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"os"

	xaiv1 "github.com/ZaguanLabs/xai-sdk-go/proto/gen/go/xai/api/v1"
	"github.com/ZaguanLabs/xai-sdk-go/xai"
	"github.com/ZaguanLabs/xai-sdk-go/xai/embed"
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

	// Example 1: Single text embedding
	fmt.Println("=== Example 1: Single Text Embedding ===")
	req := embed.NewRequest(
		"grok-embedding-1",
		embed.Text("The quick brown fox jumps over the lazy dog"),
	)

	resp, err := client.Embed().Generate(context.Background(), req)
	if err != nil {
		log.Fatalf("Embedding generation failed: %v", err)
	}

	fmt.Printf("Generated %d embedding(s)\n", len(resp.Embeddings()))
	if len(resp.Embeddings()) > 0 {
		emb := resp.Embeddings()[0]
		vectors := emb.Vectors()
		if len(vectors) > 0 {
			vec := vectors[0].FloatArray()
			fmt.Printf("Embedding dimensions: %d\n", len(vec))
			if len(vec) >= 5 {
				fmt.Printf("First 5 values: %v\n", vec[:5])
			}
		}
	}

	// Example 2: Multiple text embeddings
	fmt.Println("\n=== Example 2: Batch Text Embeddings ===")
	req2 := embed.NewRequest(
		"grok-embedding-1",
		embed.Text("Machine learning is fascinating"),
		embed.Text("Artificial intelligence is the future"),
		embed.Text("Deep learning models are powerful"),
	)

	resp2, err := client.Embed().Generate(context.Background(), req2)
	if err != nil {
		log.Fatalf("Batch embedding failed: %v", err)
	}

	fmt.Printf("Generated %d embeddings\n", len(resp2.Embeddings()))

	// Calculate cosine similarity between first two embeddings
	if len(resp2.Embeddings()) >= 2 {
		vec1 := resp2.Embeddings()[0].Vectors()[0].FloatArray()
		vec2 := resp2.Embeddings()[1].Vectors()[0].FloatArray()
		sim := cosineSimilarity(vec1, vec2)
		fmt.Printf("Cosine similarity between first two: %.4f\n", sim)
	}

	// Example 3: Image embedding
	fmt.Println("\n=== Example 3: Image Embedding ===")
	req3 := embed.NewRequest(
		"grok-embedding-1",
		embed.Image("https://example.com/image.jpg", xaiv1.ImageDetail_DETAIL_AUTO),
	)

	resp3, err := client.Embed().Generate(context.Background(), req3)
	if err != nil {
		log.Printf("Image embedding failed (expected if URL invalid): %v", err)
	} else if len(resp3.Embeddings()) > 0 && len(resp3.Embeddings()[0].Vectors()) > 0 {
		vec := resp3.Embeddings()[0].Vectors()[0].FloatArray()
		fmt.Printf("Image embedding dimensions: %d\n", len(vec))
	}

	// Example 4: Mixed text and image embeddings
	fmt.Println("\n=== Example 4: Mixed Embeddings ===")
	req4 := embed.NewRequest(
		"grok-embedding-1",
		embed.Text("A beautiful sunset"),
		embed.Image("https://example.com/sunset.jpg", xaiv1.ImageDetail_DETAIL_AUTO),
		embed.Text("Nature photography"),
	)

	resp4, err := client.Embed().Generate(context.Background(), req4)
	if err != nil {
		log.Printf("Mixed embedding failed: %v", err)
	} else {
		fmt.Printf("Generated %d mixed embeddings\n", len(resp4.Embeddings()))
	}

	fmt.Println("\n✅ Embedding examples completed!")
}

// cosineSimilarity calculates the cosine similarity between two vectors
func cosineSimilarity(a, b []float32) float64 {
	if len(a) != len(b) {
		return 0
	}

	var dotProduct, normA, normB float64
	for i := range a {
		dotProduct += float64(a[i]) * float64(b[i])
		normA += float64(a[i]) * float64(a[i])
		normB += float64(b[i]) * float64(b[i])
	}

	if normA == 0 || normB == 0 {
		return 0
	}

	return dotProduct / (math.Sqrt(normA) * math.Sqrt(normB))
}
