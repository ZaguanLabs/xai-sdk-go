package image

import (
	"context"
	"testing"

	xaiv1 "github.com/ZaguanLabs/xai-sdk-go/proto/gen/go/xai/v1"
	"google.golang.org/grpc"
)

// mockImageServiceClient implements ImageServiceClient for testing
type mockImageServiceClient struct {
	images []*xaiv1.Image
	err    error
}

func (m *mockImageServiceClient) GenerateImage(ctx context.Context, req *xaiv1.GenerateImageRequest, opts ...grpc.CallOption) (*xaiv1.GenerateImageResponse, error) {
	if m.err != nil {
		return nil, m.err
	}
	return &xaiv1.GenerateImageResponse{Images: m.images}, nil
}

func TestGenerate(t *testing.T) {
	mockClient := &mockImageServiceClient{
		images: []*xaiv1.Image{
			{Url: "https://example.com/image1.png", RevisedPrompt: "A beautiful sunset"},
			{Url: "https://example.com/image2.png", RevisedPrompt: "A beautiful sunset over mountains"},
		},
	}

	client := NewClient(mockClient)
	req := NewGenerateRequest("A beautiful sunset", "dall-e-3")
	
	images, err := client.Generate(context.Background(), req)

	if err != nil {
		t.Fatalf("Generate() returned error: %v", err)
	}

	if len(images) != 2 {
		t.Errorf("Expected 2 images, got %d", len(images))
	}

	if images[0].URL() != "https://example.com/image1.png" {
		t.Errorf("Expected first image URL to be 'https://example.com/image1.png', got '%s'", images[0].URL())
	}

	if images[0].RevisedPrompt() != "A beautiful sunset" {
		t.Errorf("Expected first image revised prompt to be 'A beautiful sunset', got '%s'", images[0].RevisedPrompt())
	}
}

func TestGenerateRequestValidation(t *testing.T) {
	tests := []struct {
		name    string
		req     *GenerateRequest
		wantErr bool
	}{
		{
			name:    "empty prompt",
			req:     NewGenerateRequest("", "dall-e-3"),
			wantErr: true,
		},
		{
			name:    "empty model",
			req:     NewGenerateRequest("A beautiful sunset", ""),
			wantErr: true,
		},
		{
			name:    "invalid n",
			req:     NewGenerateRequest("A beautiful sunset", "dall-e-3").WithN(15),
			wantErr: true,
		},
		{
			name:    "invalid size",
			req:     NewGenerateRequest("A beautiful sunset", "dall-e-3").WithSize("invalid"),
			wantErr: true,
		},
		{
			name:    "invalid quality",
			req:     NewGenerateRequest("A beautiful sunset", "dall-e-3").WithQuality("invalid"),
			wantErr: true,
		},
		{
			name:    "invalid style",
			req:     NewGenerateRequest("A beautiful sunset", "dall-e-3").WithStyle("invalid"),
			wantErr: true,
		},
		{
			name:    "valid request",
			req:     NewGenerateRequest("A beautiful sunset", "dall-e-3"),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestImageMethods(t *testing.T) {
	img := &Image{
		url:           "https://example.com/image.png",
		data:          []byte("image data"),
		revisedPrompt: "A beautiful sunset",
	}

	if img.URL() != "https://example.com/image.png" {
		t.Errorf("Expected URL to be 'https://example.com/image.png', got '%s'", img.URL())
	}

	if string(img.Data()) != "image data" {
		t.Errorf("Expected data to be 'image data', got '%s'", string(img.Data()))
	}

	if img.RevisedPrompt() != "A beautiful sunset" {
		t.Errorf("Expected revised prompt to be 'A beautiful sunset', got '%s'", img.RevisedPrompt())
	}
}

func TestNewGenerateRequest(t *testing.T) {
	req := NewGenerateRequest("A beautiful sunset", "dall-e-3")
	
	if req.prompt != "A beautiful sunset" {
		t.Errorf("Expected prompt to be 'A beautiful sunset', got '%s'", req.prompt)
	}

	if req.model != "dall-e-3" {
		t.Errorf("Expected model to be 'dall-e-3', got '%s'", req.model)
	}

	if req.size != "1024x1024" {
		t.Errorf("Expected default size to be '1024x1024', got '%s'", req.size)
	}

	if req.quality != "standard" {
		t.Errorf("Expected default quality to be 'standard', got '%s'", req.quality)
	}

	if req.style != "vivid" {
		t.Errorf("Expected default style to be 'vivid', got '%s'", req.style)
	}

	if req.n != 1 {
		t.Errorf("Expected default n to be 1, got %d", req.n)
	}
}