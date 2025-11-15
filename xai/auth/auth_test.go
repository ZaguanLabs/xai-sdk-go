package auth

import (
	"context"
	"testing"

	xaiv1 "github.com/ZaguanLabs/xai-sdk-go/proto/gen/go/xai/v1"
	"google.golang.org/grpc"
)

// mockAuthServiceClient implements AuthServiceClient for testing
type mockAuthServiceClient struct {
	resp *xaiv1.ValidateKeyResponse
	err  error
}

func (m *mockAuthServiceClient) ValidateKey(ctx context.Context, req *xaiv1.ValidateKeyRequest, opts ...grpc.CallOption) (*xaiv1.ValidateKeyResponse, error) {
	if m.err != nil {
		return nil, m.err
	}
	if m.resp != nil {
		return m.resp, nil
	}
	return &xaiv1.ValidateKeyResponse{
		Valid:        true,
		Message:      "API key is valid",
		Organization: "TestOrg",
		Project:      "TestProject",
	}, nil
}

func TestValidate(t *testing.T) {
	mockClient := &mockAuthServiceClient{}
	client := NewClient(mockClient)

	result, err := client.Validate(context.Background(), "test-api-key")

	if err != nil {
		t.Fatalf("Validate() returned error: %v", err)
	}

	if result == nil {
		t.Fatal("Expected non-nil result")
	}

	if !result.IsValid() {
		t.Errorf("Expected result to be valid")
	}

	if result.Organization() != "TestOrg" {
		t.Errorf("Expected organization to be 'TestOrg', got '%s'", result.Organization())
	}

	if result.Project() != "TestProject" {
		t.Errorf("Expected project to be 'TestProject', got '%s'", result.Project())
	}
}

func TestValidateInvalidKey(t *testing.T) {
	mockClient := &mockAuthServiceClient{
		resp: &xaiv1.ValidateKeyResponse{
			Valid:   false,
			Message: "Invalid API key",
		},
	}
	client := NewClient(mockClient)

	result, err := client.Validate(context.Background(), "invalid-api-key")

	if err != nil {
		t.Fatalf("Validate() returned error: %v", err)
	}

	if result.IsValid() {
		t.Errorf("Expected result to be invalid")
	}

	if result.Message() != "Invalid API key" {
		t.Errorf("Expected message to be 'Invalid API key', got '%s'", result.Message())
	}
}

func TestValidateEmptyAPIKey(t *testing.T) {
	mockClient := &mockAuthServiceClient{}
	client := NewClient(mockClient)

	_, err := client.Validate(context.Background(), "")

	if err == nil {
		t.Fatal("Expected error for empty API key, got nil")
	}
}

func TestValidationResultString(t *testing.T) {
	result := &ValidationResult{
		valid:        true,
		organization: "TestOrg",
		project:      "TestProject",
	}

	expected := "ValidationResult{Valid: true, Organization: TestOrg, Project: TestProject}"
	if result.String() != expected {
		t.Errorf("Expected String() to return '%s', got '%s'", expected, result.String())
	}

	// Test invalid result
	result.valid = false
	result.message = "Invalid key"
	expected = "ValidationResult{Valid: false, Message: Invalid key}"
	if result.String() != expected {
		t.Errorf("Expected String() to return '%s', got '%s'", expected, result.String())
	}
}