package models

import (
	"context"
	"testing"

	xaiv1 "github.com/ZaguanLabs/xai-sdk-go/proto/gen/go/xai/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// mockModelServiceClient implements ModelServiceClient for testing
type mockModelServiceClient struct {
	models []*xaiv1.Model
	err    error
}

func (m *mockModelServiceClient) ListModels(ctx context.Context, req *xaiv1.ListModelsRequest, opts ...grpc.CallOption) (*xaiv1.ListModelsResponse, error) {
	if m.err != nil {
		return nil, m.err
	}
	return &xaiv1.ListModelsResponse{Models: m.models}, nil
}

func (m *mockModelServiceClient) GetModel(ctx context.Context, req *xaiv1.GetModelRequest, opts ...grpc.CallOption) (*xaiv1.Model, error) {
	if m.err != nil {
		return nil, m.err
	}
	for _, model := range m.models {
		if model.Id == req.ModelId {
			return model, nil
		}
	}
	return nil, status.Error(codes.NotFound, "model not found")
}

func TestList(t *testing.T) {
	mockClient := &mockModelServiceClient{
		models: []*xaiv1.Model{
			{Id: "grok-beta", Name: "Grok Beta", Description: "Grok Beta model", MaxTokens: 8192},
			{Id: "grok-vision", Name: "Grok Vision", Description: "Grok Vision model", MaxTokens: 4096},
		},
	}

	client := NewClient(mockClient)
	models, err := client.List(context.Background())

	if err != nil {
		t.Fatalf("List() returned error: %v", err)
	}

	if len(models) != 2 {
		t.Errorf("Expected 2 models, got %d", len(models))
	}

	if models[0].ID() != "grok-beta" {
		t.Errorf("Expected first model ID to be 'grok-beta', got '%s'", models[0].ID())
	}

	if models[0].Name() != "Grok Beta" {
		t.Errorf("Expected first model name to be 'Grok Beta', got '%s'", models[0].Name())
	}
}

func TestGet(t *testing.T) {
	mockClient := &mockModelServiceClient{
		models: []*xaiv1.Model{
			{Id: "grok-beta", Name: "Grok Beta", Description: "Grok Beta model", MaxTokens: 8192},
		},
	}

	client := NewClient(mockClient)
	model, err := client.Get(context.Background(), "grok-beta")

	if err != nil {
		t.Fatalf("Get() returned error: %v", err)
	}

	if model.ID() != "grok-beta" {
		t.Errorf("Expected model ID to be 'grok-beta', got '%s'", model.ID())
	}

	if model.MaxTokens() != 8192 {
		t.Errorf("Expected max tokens to be 8192, got %d", model.MaxTokens())
	}
}

func TestGetNotFound(t *testing.T) {
	mockClient := &mockModelServiceClient{
		models: []*xaiv1.Model{},
	}

	client := NewClient(mockClient)
	_, err := client.Get(context.Background(), "nonexistent")

	if err == nil {
		t.Fatal("Expected error for nonexistent model, got nil")
	}
}

func TestGetEmptyID(t *testing.T) {
	mockClient := &mockModelServiceClient{}
	client := NewClient(mockClient)

	_, err := client.Get(context.Background(), "")

	if err == nil {
		t.Fatal("Expected error for empty model ID, got nil")
	}
}

func TestModelString(t *testing.T) {
	model := &Model{
		id:        "test-model",
		name:      "Test Model",
		maxTokens: 4096,
	}

	expected := "Model{ID: test-model, Name: Test Model, MaxTokens: 4096}"
	if model.String() != expected {
		t.Errorf("Expected String() to return '%s', got '%s'", expected, model.String())
	}
}