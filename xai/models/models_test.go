package models

import (
	"context"
	"testing"
	"time"

	xaiv1 "github.com/ZaguanLabs/xai-sdk-go/proto/gen/go/xai/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Mock gRPC client
type mockModelsClient struct {
	xaiv1.ModelsClient
	listLangModels     func(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*xaiv1.ListLanguageModelsResponse, error)
	listEmbedModels    func(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*xaiv1.ListEmbeddingModelsResponse, error)
	listImageGenModels func(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*xaiv1.ListImageGenerationModelsResponse, error)
	getLangModel       func(ctx context.Context, in *xaiv1.GetModelRequest, opts ...grpc.CallOption) (*xaiv1.LanguageModel, error)
	getEmbedModel      func(ctx context.Context, in *xaiv1.GetModelRequest, opts ...grpc.CallOption) (*xaiv1.EmbeddingModel, error)
	getImageGenModel   func(ctx context.Context, in *xaiv1.GetModelRequest, opts ...grpc.CallOption) (*xaiv1.ImageGenerationModel, error)
}

func (m *mockModelsClient) ListLanguageModels(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*xaiv1.ListLanguageModelsResponse, error) {
	if m.listLangModels != nil {
		return m.listLangModels(ctx, in, opts...)
	}
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (m *mockModelsClient) ListEmbeddingModels(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*xaiv1.ListEmbeddingModelsResponse, error) {
	if m.listEmbedModels != nil {
		return m.listEmbedModels(ctx, in, opts...)
	}
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (m *mockModelsClient) ListImageGenerationModels(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*xaiv1.ListImageGenerationModelsResponse, error) {
	if m.listImageGenModels != nil {
		return m.listImageGenModels(ctx, in, opts...)
	}
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (m *mockModelsClient) GetLanguageModel(ctx context.Context, in *xaiv1.GetModelRequest, opts ...grpc.CallOption) (*xaiv1.LanguageModel, error) {
	if m.getLangModel != nil {
		return m.getLangModel(ctx, in, opts...)
	}
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (m *mockModelsClient) GetEmbeddingModel(ctx context.Context, in *xaiv1.GetModelRequest, opts ...grpc.CallOption) (*xaiv1.EmbeddingModel, error) {
	if m.getEmbedModel != nil {
		return m.getEmbedModel(ctx, in, opts...)
	}
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (m *mockModelsClient) GetImageGenerationModel(ctx context.Context, in *xaiv1.GetModelRequest, opts ...grpc.CallOption) (*xaiv1.ImageGenerationModel, error) {
	if m.getImageGenModel != nil {
		return m.getImageGenModel(ctx, in, opts...)
	}
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func TestNewClient(t *testing.T) {
	mockClient := &mockModelsClient{}
	client := NewClient(mockClient)

	if client == nil {
		t.Fatal("NewClient returned nil")
	}
	if client.grpcClient != mockClient {
		t.Error("grpcClient not set correctly")
	}
}

func TestListLanguageModels(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name      string
		mockFunc  func(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*xaiv1.ListLanguageModelsResponse, error)
		wantCount int
		wantErr   bool
	}{
		{
			name: "success",
			mockFunc: func(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*xaiv1.ListLanguageModelsResponse, error) {
				return &xaiv1.ListLanguageModelsResponse{
					Models: []*xaiv1.LanguageModel{
						{Name: "grok-1", Version: "1.0", Created: timestamppb.New(now)},
						{Name: "grok-2", Version: "2.0", Created: timestamppb.New(now)},
					},
				}, nil
			},
			wantCount: 2,
			wantErr:   false,
		},
		{
			name: "unauthenticated",
			mockFunc: func(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*xaiv1.ListLanguageModelsResponse, error) {
				return nil, status.Error(codes.Unauthenticated, "invalid API key")
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockModelsClient{listLangModels: tt.mockFunc}
			client := NewClient(mock)

			models, err := client.ListLanguageModels(context.Background())

			if (err != nil) != tt.wantErr {
				t.Errorf("ListLanguageModels() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && len(models) != tt.wantCount {
				t.Errorf("ListLanguageModels() returned %d models, want %d", len(models), tt.wantCount)
			}
		})
	}
}

func TestGetLanguageModel(t *testing.T) {
	tests := []struct {
		name      string
		modelName string
		mockFunc  func(ctx context.Context, in *xaiv1.GetModelRequest, opts ...grpc.CallOption) (*xaiv1.LanguageModel, error)
		wantErr   bool
	}{
		{
			name:      "success",
			modelName: "grok-1",
			mockFunc: func(ctx context.Context, in *xaiv1.GetModelRequest, opts ...grpc.CallOption) (*xaiv1.LanguageModel, error) {
				return &xaiv1.LanguageModel{
					Name:    "grok-1",
					Version: "1.0",
				}, nil
			},
			wantErr: false,
		},
		{
			name:      "empty name",
			modelName: "",
			wantErr:   true,
		},
		{
			name:      "not found",
			modelName: "nonexistent",
			mockFunc: func(ctx context.Context, in *xaiv1.GetModelRequest, opts ...grpc.CallOption) (*xaiv1.LanguageModel, error) {
				return nil, status.Error(codes.NotFound, "model not found")
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockModelsClient{getLangModel: tt.mockFunc}
			client := NewClient(mock)

			model, err := client.GetLanguageModel(context.Background(), tt.modelName)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetLanguageModel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && model == nil {
				t.Error("GetLanguageModel() returned nil model")
			}
		})
	}
}

func TestLanguageModelMethods(t *testing.T) {
	now := time.Now()
	model := &LanguageModel{
		name:                     "grok-1",
		aliases:                  []string{"grok"},
		version:                  "1.0",
		inputModalities:          []string{"text"},
		outputModalities:         []string{"text"},
		promptTextTokenPrice:     100,
		promptImageTokenPrice:    200,
		cachedPromptTokenPrice:   50,
		completionTextTokenPrice: 150,
		searchPrice:              75,
		created:                  now,
		maxPromptLength:          4096,
		systemFingerprint:        "fp123",
	}

	if model.Name() != "grok-1" {
		t.Errorf("Name() = %v, want grok-1", model.Name())
	}
	if len(model.Aliases()) != 1 || model.Aliases()[0] != "grok" {
		t.Errorf("Aliases() = %v, want [grok]", model.Aliases())
	}
	if model.Version() != "1.0" {
		t.Errorf("Version() = %v, want 1.0", model.Version())
	}
	if model.MaxPromptLength() != 4096 {
		t.Errorf("MaxPromptLength() = %v, want 4096", model.MaxPromptLength())
	}
	if model.PromptTextTokenPrice() != 100 {
		t.Errorf("PromptTextTokenPrice() = %v, want 100", model.PromptTextTokenPrice())
	}

	str := model.String()
	if str == "" {
		t.Error("String() returned empty string")
	}
}

func TestConvertLanguageModel(t *testing.T) {
	now := time.Now()
	proto := &xaiv1.LanguageModel{
		Name:                     "grok-1",
		Aliases:                  []string{"grok"},
		Version:                  "1.0",
		InputModalities:          []xaiv1.Modality{xaiv1.Modality_TEXT},
		OutputModalities:         []xaiv1.Modality{xaiv1.Modality_TEXT},
		PromptTextTokenPrice:     100,
		PromptImageTokenPrice:    200,
		CachedPromptTokenPrice:   50,
		CompletionTextTokenPrice: 150,
		SearchPrice:              75,
		Created:                  timestamppb.New(now),
		MaxPromptLength:          4096,
		SystemFingerprint:        "fp123",
	}

	model := convertLanguageModel(proto)

	if model == nil {
		t.Fatal("convertLanguageModel returned nil")
	}
	if model.Name() != "grok-1" {
		t.Errorf("Name = %v, want grok-1", model.Name())
	}
	if model.Version() != "1.0" {
		t.Errorf("Version = %v, want 1.0", model.Version())
	}
}

func TestModalityToStrings(t *testing.T) {
	modalities := []xaiv1.Modality{
		xaiv1.Modality_TEXT,
		xaiv1.Modality_IMAGE,
	}

	result := modalityToStrings(modalities)

	if len(result) != 2 {
		t.Errorf("modalityToStrings() returned %d items, want 2", len(result))
	}
}
