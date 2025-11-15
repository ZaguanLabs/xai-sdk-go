package xai

import (
	"context"
	"testing"
	"time"

	"github.com/ZaguanLabs/xai-sdk-go/xai/internal/metadata"
	"google.golang.org/grpc"
)

func TestNewClient(t *testing.T) {
	t.Run("ValidConfig", func(t *testing.T) {
		config := NewConfigWithAPIKey("test-api-key")
		client, err := NewClient(config)
		
		if err != nil {
			t.Errorf("Should not return error for valid config: %v", err)
		}
		if client == nil {
			t.Error("NewClient should not return nil")
		}
		
		// Verify client properties
		if client.Config() != config {
			t.Error("Client should return the same config")
		}
		
		if client.IsClosed() {
			t.Error("New client should not be closed")
		}
		
		if client.Metadata() == nil {
			t.Error("Client metadata should not be nil")
		}
	})

	t.Run("NilConfig", func(t *testing.T) {
		client, err := NewClient(nil)
		
		if err != nil {
			t.Errorf("Should not return error for nil config: %v", err)
		}
		if client == nil {
			t.Error("NewClient should not return nil for nil config")
		}
	})

	t.Run("InvalidConfig", func(t *testing.T) {
		config := &Config{
			APIKey: "", // Empty API key should cause validation error
		}
		
		client, err := NewClient(config)
		
		if err == nil {
			t.Error("Should return error for invalid config")
		}
		if client != nil {
			t.Error("Should return nil client for invalid config")
		}
	})
}

func TestNewClientWithAPIKey(t *testing.T) {
	apiKey := "test-api-key"
	client, err := NewClientWithAPIKey(apiKey)
	
	if err != nil {
		t.Errorf("Should not return error: %v", err)
	}
	if client == nil {
		t.Error("NewClientWithAPIKey should not return nil")
	}
	
	if client.Config().APIKey != apiKey {
		t.Errorf("Expected API key %s, got %s", apiKey, client.Config().APIKey)
	}
}

func TestNewClientFromEnvironment(t *testing.T) {
	// Set environment variable
	apiKey := "env-api-key"
	t.Setenv("XAI_API_KEY", apiKey)
	
	client, err := NewClientFromEnvironment()
	
	if err != nil {
		t.Errorf("Should not return error: %v", err)
	}
	if client == nil {
		t.Error("NewClientFromEnvironment should not return nil")
	}
	
	if client.Config().APIKey != apiKey {
		t.Errorf("Expected API key %s, got %s", apiKey, client.Config().APIKey)
	}
}

func TestClientConfig(t *testing.T) {
	config := NewConfigWithAPIKey("test-api-key")
	client, err := NewClient(config)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()
	
	returnedConfig := client.Config()
	if returnedConfig != config {
		t.Error("Client should return the same config instance")
	}
}

func TestClientMetadata(t *testing.T) {
	config := NewConfigWithAPIKey("test-api-key")
	client, err := NewClient(config)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()
	
	metadata := client.Metadata()
	if metadata == nil {
		t.Error("Client metadata should not be nil")
	}
	
	if metadata.APIKey != "test-api-key" {
		t.Errorf("Expected API key 'test-api-key', got %s", metadata.APIKey)
	}
}

func TestClientCreatedAt(t *testing.T) {
	before := time.Now()
	client, err := NewClient(NewConfigWithAPIKey("test-api-key"))
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	after := time.Now()
	
	createdAt := client.CreatedAt()
	if createdAt.Before(before) || createdAt.After(after) {
		t.Error("CreatedAt should be within the creation time window")
	}
}

func TestClientIsClosed(t *testing.T) {
	config := NewConfigWithAPIKey("test-api-key")
	client, err := NewClient(config)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	
	// Client should not be closed initially
	if client.IsClosed() {
		t.Error("New client should not be closed")
	}
	
	// Close the client
	err = client.Close()
	if err != nil {
		t.Errorf("Should not return error when closing: %v", err)
	}
	
	// Client should be closed now
	if !client.IsClosed() {
		t.Error("Client should be closed after calling Close()")
	}
	
	// Closing again should not return error
	err = client.Close()
	if err != nil {
		t.Errorf("Closing already closed client should not return error: %v", err)
	}
}

func TestClientGRPCConnection(t *testing.T) {
	config := NewConfigWithAPIKey("test-api-key")
	client, err := NewClient(config)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()
	
	grpcConn := client.GRPCConnection()
	if grpcConn == nil {
		t.Error("gRPC connection should not be nil")
	}
}

func TestClientNewContext(t *testing.T) {
	config := NewConfigWithAPIKey("test-api-key")
	client, err := NewClient(config)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()
	
	ctx := client.NewContext(nil)
	if ctx == nil {
		t.Error("NewContext should not return nil")
	}
	
	// Test with existing context
	existingCtx := context.Background()
	resultCtx := client.NewContext(existingCtx)
	if resultCtx == nil {
		t.Error("NewContext should not return nil for existing context")
	}
}

func TestClientNewContextWithTimeout(t *testing.T) {
	config := NewConfigWithAPIKey("test-api-key")
	client, err := NewClient(config)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()
	
	timeout := 5 * time.Second
	ctx, cancel := client.NewContextWithTimeout(nil, timeout)
	defer cancel()
	
	if ctx == nil {
		t.Error("NewContextWithTimeout should not return nil context")
	}
	if cancel == nil {
		t.Error("NewContextWithTimeout should not return nil cancel function")
	}
	
	// Check that context has timeout
	deadline, ok := ctx.Deadline()
	if !ok {
		t.Error("Context should have a deadline")
	}
	
	expectedDeadline := time.Now().Add(timeout)
	if deadline.Before(expectedDeadline.Add(-1*time.Second)) || deadline.After(expectedDeadline.Add(1*time.Second)) {
		t.Error("Context deadline should match requested timeout")
	}
}

func TestClientNewContextWithCancel(t *testing.T) {
	config := NewConfigWithAPIKey("test-api-key")
	client, err := NewClient(config)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()
	
	ctx, cancel := client.NewContextWithCancel(nil)
	defer cancel()
	
	if ctx == nil {
		t.Error("NewContextWithCancel should not return nil context")
	}
	if cancel == nil {
		t.Error("NewContextWithCancel should not return nil cancel function")
	}
}

func TestClientNewContextWithDeadline(t *testing.T) {
	config := NewConfigWithAPIKey("test-api-key")
	client, err := NewClient(config)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()
	
	deadline := time.Now().Add(10 * time.Second)
	ctx, cancel := client.NewContextWithDeadline(nil, deadline)
	defer cancel()
	
	if ctx == nil {
		t.Error("NewContextWithDeadline should not return nil context")
	}
	if cancel == nil {
		t.Error("NewContextWithDeadline should not return nil cancel function")
	}
	
	actualDeadline, ok := ctx.Deadline()
	if !ok {
		t.Error("Context should have a deadline")
	}
	
	if !actualDeadline.Equal(deadline) {
		t.Errorf("Expected deadline %v, got %v", deadline, actualDeadline)
	}
}

func TestClientCloseWithContext(t *testing.T) {
	config := NewConfigWithAPIKey("test-api-key")
	client, err := NewClient(config)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	
	// Test normal close
	ctx := context.Background()
	err = client.CloseWithContext(ctx)
	if err != nil {
		t.Errorf("CloseWithContext should not return error: %v", err)
	}
}

func TestClientHealthCheck(t *testing.T) {
	config := NewConfigWithAPIKey("test-api-key")
	client, err := NewClient(config)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()
	
	ctx := context.Background()
	err = client.HealthCheck(ctx)
	// Health check may fail due to connection issues in test environment,
	// but it shouldn't panic or return nil error in unexpected ways
	if err != nil {
		// This is expected if the server is not actually available
		t.Logf("Health check failed (expected in test environment): %v", err)
	}
}

func TestClientWithTimeout(t *testing.T) {
	config := NewConfigWithAPIKey("test-api-key")
	client, err := NewClient(config)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()
	
	newTimeout := 45 * time.Second
	newClient := client.WithTimeout(newTimeout)
	
	if newClient == nil {
		t.Error("WithTimeout should not return nil")
	}
	
	if newClient.Config().Timeout != newTimeout {
		t.Errorf("Expected timeout %v, got %v", newTimeout, newClient.Config().Timeout)
	}
}

func TestClientWithAPIKey(t *testing.T) {
	config := NewConfigWithAPIKey("original-api-key")
	client, err := NewClient(config)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()
	
	newAPIKey := "new-api-key"
	newClient := client.WithAPIKey(newAPIKey)
	
	if newClient == nil {
		t.Error("WithAPIKey should not return nil")
	}
	
	if newClient.Config().APIKey != newAPIKey {
		t.Errorf("Expected API key %s, got %s", newAPIKey, newClient.Config().APIKey)
	}
}

func TestClientString(t *testing.T) {
	config := NewConfigWithAPIKey("1234567890abcdef")
	config.Host = "test.host"
	config.Environment = "test"
	
	client, err := NewClient(config)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()
	
	str := client.String()
	
	if str == "" {
		t.Error("Client string should not be empty")
	}
	
	// Should contain host information
	if !containsString(str, "test.host") {
		t.Errorf("Client string should contain host: %s", str)
	}
	
	// Should contain environment information
	if !containsString(str, "test") {
		t.Errorf("Client string should contain environment: %s", str)
	}
}

func TestClientGetHealthStatus(t *testing.T) {
	config := NewConfigWithAPIKey("test-api-key")
	client, err := NewClient(config)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()
	
	status := client.GetHealthStatus()
	
	if status.Status == "" {
		t.Error("Health status should not have empty status")
	}
	
	if status.Timestamp.IsZero() {
		t.Error("Health status timestamp should not be zero")
	}
	
	if status.Status == "closed" && !client.IsClosed() {
		t.Error("Should not return closed status for non-closed client")
	}
}

func TestClientEnsureGRPCConnection(t *testing.T) {
	config := NewConfigWithAPIKey("test-api-key")
	client, err := NewClient(config)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()
	
	err = client.EnsureGRPCConnection()
	if err != nil {
		t.Errorf("EnsureGRPCConnection should not return error: %v", err)
	}
}

func BenchmarkNewClient(b *testing.B) {
	for i := 0; i < b.N; i++ {
		client, _ := NewClient(NewConfigWithAPIKey("benchmark-api-key"))
		if client != nil {
			client.Close()
		}
	}
}

func BenchmarkClientNewContext(b *testing.B) {
	client, _ := NewClient(NewConfigWithAPIKey("benchmark-api-key"))
	defer client.Close()
	
	ctx := context.Background()
	for i := 0; i < b.N; i++ {
		_ = client.NewContext(ctx)
	}
}

// Helper functions
func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || (len(s) > len(substr) && findString(s, substr)))
}

func findString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}