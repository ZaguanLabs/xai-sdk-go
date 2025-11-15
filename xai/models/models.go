// Package models provides model information functionality for xAI SDK.
package models

import (
	"context"
	"fmt"
	"time"

	xaiv1 "github.com/ZaguanLabs/xai-sdk-go/proto/gen/go/xai/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// LanguageModel represents an xAI language model.
type LanguageModel struct {
	name                     string
	aliases                  []string
	version                  string
	inputModalities          []string
	outputModalities         []string
	promptTextTokenPrice     int64
	promptImageTokenPrice    int64
	cachedPromptTokenPrice   int64
	completionTextTokenPrice int64
	searchPrice              int64
	created                  time.Time
	maxPromptLength          int32
	systemFingerprint        string
}

// EmbeddingModel represents an xAI embedding model.
type EmbeddingModel struct {
	name                  string
	aliases               []string
	version               string
	inputModalities       []string
	outputModalities      []string
	promptTextTokenPrice  int64
	promptImageTokenPrice int64
	created               time.Time
	systemFingerprint     string
}

// ImageGenerationModel represents an xAI image generation model.
type ImageGenerationModel struct {
	name              string
	aliases           []string
	version           string
	inputModalities   []string
	outputModalities  []string
	imagePrice        int64
	created           time.Time
	maxPromptLength   int32
	systemFingerprint string
}

// Client provides model information functionality.
type Client struct {
	grpcClient xaiv1.ModelsClient
}

// NewClient creates a new models client.
func NewClient(grpcClient xaiv1.ModelsClient) *Client {
	return &Client{
		grpcClient: grpcClient,
	}
}

// ListLanguageModels lists all available language models.
func (c *Client) ListLanguageModels(ctx context.Context) ([]*LanguageModel, error) {
	resp, err := c.grpcClient.ListLanguageModels(ctx, &emptypb.Empty{})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.Unauthenticated:
				return nil, fmt.Errorf("authentication failed: %s", st.Message())
			case codes.PermissionDenied:
				return nil, fmt.Errorf("permission denied: %s", st.Message())
			case codes.Unavailable:
				return nil, fmt.Errorf("service unavailable: %s", st.Message())
			default:
				return nil, fmt.Errorf("list language models failed (%s): %s", st.Code().String(), st.Message())
			}
		}
		return nil, fmt.Errorf("list language models failed: %w", err)
	}

	if resp == nil {
		return nil, fmt.Errorf("received nil response")
	}

	models := make([]*LanguageModel, 0, len(resp.Models))
	for _, m := range resp.Models {
		models = append(models, convertLanguageModel(m))
	}

	return models, nil
}

// ListEmbeddingModels lists all available embedding models.
func (c *Client) ListEmbeddingModels(ctx context.Context) ([]*EmbeddingModel, error) {
	resp, err := c.grpcClient.ListEmbeddingModels(ctx, &emptypb.Empty{})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.Unauthenticated:
				return nil, fmt.Errorf("authentication failed: %s", st.Message())
			case codes.PermissionDenied:
				return nil, fmt.Errorf("permission denied: %s", st.Message())
			case codes.Unavailable:
				return nil, fmt.Errorf("service unavailable: %s", st.Message())
			default:
				return nil, fmt.Errorf("list embedding models failed (%s): %s", st.Code().String(), st.Message())
			}
		}
		return nil, fmt.Errorf("list embedding models failed: %w", err)
	}

	if resp == nil {
		return nil, fmt.Errorf("received nil response")
	}

	models := make([]*EmbeddingModel, 0, len(resp.Models))
	for _, m := range resp.Models {
		models = append(models, convertEmbeddingModel(m))
	}

	return models, nil
}

// ListImageGenerationModels lists all available image generation models.
func (c *Client) ListImageGenerationModels(ctx context.Context) ([]*ImageGenerationModel, error) {
	resp, err := c.grpcClient.ListImageGenerationModels(ctx, &emptypb.Empty{})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.Unauthenticated:
				return nil, fmt.Errorf("authentication failed: %s", st.Message())
			case codes.PermissionDenied:
				return nil, fmt.Errorf("permission denied: %s", st.Message())
			case codes.Unavailable:
				return nil, fmt.Errorf("service unavailable: %s", st.Message())
			default:
				return nil, fmt.Errorf("list image generation models failed (%s): %s", st.Code().String(), st.Message())
			}
		}
		return nil, fmt.Errorf("list image generation models failed: %w", err)
	}

	if resp == nil {
		return nil, fmt.Errorf("received nil response")
	}

	models := make([]*ImageGenerationModel, 0, len(resp.Models))
	for _, m := range resp.Models {
		models = append(models, convertImageGenerationModel(m))
	}

	return models, nil
}

// GetLanguageModel retrieves information about a specific language model.
func (c *Client) GetLanguageModel(ctx context.Context, name string) (*LanguageModel, error) {
	if name == "" {
		return nil, fmt.Errorf("model name is required")
	}

	resp, err := c.grpcClient.GetLanguageModel(ctx, &xaiv1.GetModelRequest{Name: name})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.Unauthenticated:
				return nil, fmt.Errorf("authentication failed: %s", st.Message())
			case codes.PermissionDenied:
				return nil, fmt.Errorf("permission denied: %s", st.Message())
			case codes.NotFound:
				return nil, fmt.Errorf("model not found: %s", st.Message())
			case codes.Unavailable:
				return nil, fmt.Errorf("service unavailable: %s", st.Message())
			default:
				return nil, fmt.Errorf("get language model failed (%s): %s", st.Code().String(), st.Message())
			}
		}
		return nil, fmt.Errorf("get language model failed: %w", err)
	}

	return convertLanguageModel(resp), nil
}

// GetEmbeddingModel retrieves information about a specific embedding model.
func (c *Client) GetEmbeddingModel(ctx context.Context, name string) (*EmbeddingModel, error) {
	if name == "" {
		return nil, fmt.Errorf("model name is required")
	}

	resp, err := c.grpcClient.GetEmbeddingModel(ctx, &xaiv1.GetModelRequest{Name: name})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.Unauthenticated:
				return nil, fmt.Errorf("authentication failed: %s", st.Message())
			case codes.PermissionDenied:
				return nil, fmt.Errorf("permission denied: %s", st.Message())
			case codes.NotFound:
				return nil, fmt.Errorf("model not found: %s", st.Message())
			case codes.Unavailable:
				return nil, fmt.Errorf("service unavailable: %s", st.Message())
			default:
				return nil, fmt.Errorf("get embedding model failed (%s): %s", st.Code().String(), st.Message())
			}
		}
		return nil, fmt.Errorf("get embedding model failed: %w", err)
	}

	return convertEmbeddingModel(resp), nil
}

// GetImageGenerationModel retrieves information about a specific image generation model.
func (c *Client) GetImageGenerationModel(ctx context.Context, name string) (*ImageGenerationModel, error) {
	if name == "" {
		return nil, fmt.Errorf("model name is required")
	}

	resp, err := c.grpcClient.GetImageGenerationModel(ctx, &xaiv1.GetModelRequest{Name: name})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.Unauthenticated:
				return nil, fmt.Errorf("authentication failed: %s", st.Message())
			case codes.PermissionDenied:
				return nil, fmt.Errorf("permission denied: %s", st.Message())
			case codes.NotFound:
				return nil, fmt.Errorf("model not found: %s", st.Message())
			case codes.Unavailable:
				return nil, fmt.Errorf("service unavailable: %s", st.Message())
			default:
				return nil, fmt.Errorf("get image generation model failed (%s): %s", st.Code().String(), st.Message())
			}
		}
		return nil, fmt.Errorf("get image generation model failed: %w", err)
	}

	return convertImageGenerationModel(resp), nil
}

// Helper functions to convert proto models to SDK models

func convertLanguageModel(m *xaiv1.LanguageModel) *LanguageModel {
	var created time.Time
	if m.Created != nil {
		created = m.Created.AsTime()
	}

	return &LanguageModel{
		name:                     m.Name,
		aliases:                  m.Aliases,
		version:                  m.Version,
		inputModalities:          modalityToStrings(m.InputModalities),
		outputModalities:         modalityToStrings(m.OutputModalities),
		promptTextTokenPrice:     m.PromptTextTokenPrice,
		promptImageTokenPrice:    m.PromptImageTokenPrice,
		cachedPromptTokenPrice:   m.CachedPromptTokenPrice,
		completionTextTokenPrice: m.CompletionTextTokenPrice,
		searchPrice:              m.SearchPrice,
		created:                  created,
		maxPromptLength:          m.MaxPromptLength,
		systemFingerprint:        m.SystemFingerprint,
	}
}

func convertEmbeddingModel(m *xaiv1.EmbeddingModel) *EmbeddingModel {
	var created time.Time
	if m.Created != nil {
		created = m.Created.AsTime()
	}

	return &EmbeddingModel{
		name:                  m.Name,
		aliases:               m.Aliases,
		version:               m.Version,
		inputModalities:       modalityToStrings(m.InputModalities),
		outputModalities:      modalityToStrings(m.OutputModalities),
		promptTextTokenPrice:  m.PromptTextTokenPrice,
		promptImageTokenPrice: m.PromptImageTokenPrice,
		created:               created,
		systemFingerprint:     m.SystemFingerprint,
	}
}

func convertImageGenerationModel(m *xaiv1.ImageGenerationModel) *ImageGenerationModel {
	var created time.Time
	if m.Created != nil {
		created = m.Created.AsTime()
	}

	return &ImageGenerationModel{
		name:              m.Name,
		aliases:           m.Aliases,
		version:           m.Version,
		inputModalities:   modalityToStrings(m.InputModalities),
		outputModalities:  modalityToStrings(m.OutputModalities),
		imagePrice:        m.ImagePrice,
		created:           created,
		maxPromptLength:   m.MaxPromptLength,
		systemFingerprint: m.SystemFingerprint,
	}
}

func modalityToStrings(modalities []xaiv1.Modality) []string {
	result := make([]string, len(modalities))
	for i, m := range modalities {
		result[i] = m.String()
	}
	return result
}

// LanguageModel methods

func (m *LanguageModel) Name() string                    { return m.name }
func (m *LanguageModel) Aliases() []string               { return m.aliases }
func (m *LanguageModel) Version() string                 { return m.version }
func (m *LanguageModel) InputModalities() []string       { return m.inputModalities }
func (m *LanguageModel) OutputModalities() []string      { return m.outputModalities }
func (m *LanguageModel) PromptTextTokenPrice() int64     { return m.promptTextTokenPrice }
func (m *LanguageModel) PromptImageTokenPrice() int64    { return m.promptImageTokenPrice }
func (m *LanguageModel) CachedPromptTokenPrice() int64   { return m.cachedPromptTokenPrice }
func (m *LanguageModel) CompletionTextTokenPrice() int64 { return m.completionTextTokenPrice }
func (m *LanguageModel) SearchPrice() int64              { return m.searchPrice }
func (m *LanguageModel) Created() time.Time              { return m.created }
func (m *LanguageModel) MaxPromptLength() int32          { return m.maxPromptLength }
func (m *LanguageModel) SystemFingerprint() string       { return m.systemFingerprint }

func (m *LanguageModel) String() string {
	return fmt.Sprintf("LanguageModel{Name: %s, Version: %s, MaxPromptLength: %d}", m.name, m.version, m.maxPromptLength)
}

// EmbeddingModel methods

func (m *EmbeddingModel) Name() string                 { return m.name }
func (m *EmbeddingModel) Aliases() []string            { return m.aliases }
func (m *EmbeddingModel) Version() string              { return m.version }
func (m *EmbeddingModel) InputModalities() []string    { return m.inputModalities }
func (m *EmbeddingModel) OutputModalities() []string   { return m.outputModalities }
func (m *EmbeddingModel) PromptTextTokenPrice() int64  { return m.promptTextTokenPrice }
func (m *EmbeddingModel) PromptImageTokenPrice() int64 { return m.promptImageTokenPrice }
func (m *EmbeddingModel) Created() time.Time           { return m.created }
func (m *EmbeddingModel) SystemFingerprint() string    { return m.systemFingerprint }

func (m *EmbeddingModel) String() string {
	return fmt.Sprintf("EmbeddingModel{Name: %s, Version: %s}", m.name, m.version)
}

// ImageGenerationModel methods

func (m *ImageGenerationModel) Name() string               { return m.name }
func (m *ImageGenerationModel) Aliases() []string          { return m.aliases }
func (m *ImageGenerationModel) Version() string            { return m.version }
func (m *ImageGenerationModel) InputModalities() []string  { return m.inputModalities }
func (m *ImageGenerationModel) OutputModalities() []string { return m.outputModalities }
func (m *ImageGenerationModel) ImagePrice() int64          { return m.imagePrice }
func (m *ImageGenerationModel) Created() time.Time         { return m.created }
func (m *ImageGenerationModel) MaxPromptLength() int32     { return m.maxPromptLength }
func (m *ImageGenerationModel) SystemFingerprint() string  { return m.systemFingerprint }

func (m *ImageGenerationModel) String() string {
	return fmt.Sprintf("ImageGenerationModel{Name: %s, Version: %s, MaxPromptLength: %d}", m.name, m.version, m.maxPromptLength)
}
