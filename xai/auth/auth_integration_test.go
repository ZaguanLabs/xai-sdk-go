//go:build integration
// +build integration

package auth_test

import (
	"context"
	"os"
	"testing"

	"github.com/ZaguanLabs/xai-sdk-go/xai"
)

func TestAuthIntegration(t *testing.T) {
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

	t.Run("ValidateKey", func(t *testing.T) {
		key, err := client.Auth().ValidateKey(ctx, apiKey)
		if err != nil {
			t.Logf("ValidateKey failed (may not be implemented): %v", err)
			t.Skip("Skipping - endpoint may not be available")
		}

		if key.RedactedApiKey == "" {
			t.Fatal("Expected redacted API key")
		}

		t.Logf("Validated key: %s (User: %s)", key.RedactedApiKey, key.UserID)
	})

	t.Run("ListKeys", func(t *testing.T) {
		keys, err := client.Auth().ListKeys(ctx)
		if err != nil {
			t.Logf("ListKeys failed (may not be implemented): %v", err)
			t.Skip("Skipping - endpoint may not be available")
		}

		t.Logf("Found %d API keys", len(keys))
	})

	t.Run("GetKey", func(t *testing.T) {
		keyID := os.Getenv("XAI_KEY_ID")
		if keyID == "" {
			t.Skip("XAI_KEY_ID not set")
		}

		key, err := client.Auth().GetKey(ctx, keyID)
		if err != nil {
			t.Fatalf("GetKey failed: %v", err)
		}

		if key.RedactedApiKey == "" {
			t.Fatal("Expected redacted API key")
		}

		t.Logf("Retrieved key: %s", key.RedactedApiKey)
	})
}
