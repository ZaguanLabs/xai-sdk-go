package image

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	xaiv1 "github.com/ZaguanLabs/xai-sdk-go/proto/gen/go/xai/v1"
	"github.com/ZaguanLabs/xai-sdk-go/xai/internal/rest"
)

func TestNewClient(t *testing.T) {
	restClient := &rest.Client{}
	client := NewClient(restClient)

	if client == nil {
		t.Fatal("NewClient returned nil")
	}
	if client.restClient != restClient {
		t.Error("restClient not set correctly")
	}
}

func TestNewRequest(t *testing.T) {
	req := NewRequest("A beautiful sunset", "image-model")

	if req == nil {
		t.Fatal("NewRequest returned nil")
	}
	if req.Prompt != "A beautiful sunset" {
		t.Errorf("Prompt = %v, want A beautiful sunset", req.Prompt)
	}
	if req.Model != "image-model" {
		t.Errorf("Model = %v, want image-model", req.Model)
	}
	if req.N != 1 {
		t.Errorf("N = %v, want 1", req.N)
	}
	if req.Format != xaiv1.ImageFormat_IMG_FORMAT_URL {
		t.Errorf("Format = %v, want IMG_FORMAT_URL", req.Format)
	}
}

func TestRequestBuilders(t *testing.T) {
	req := NewRequest("test", "model").
		WithCount(3).
		WithUser("test-user").
		WithFormat(xaiv1.ImageFormat_IMG_FORMAT_BASE64).
		WithImage("http://example.com/img.jpg", xaiv1.ImageDetail_DETAIL_HIGH)

	if req.N != 3 {
		t.Errorf("N = %v, want 3", req.N)
	}
	if req.User != "test-user" {
		t.Errorf("User = %v, want test-user", req.User)
	}
	if req.Format != xaiv1.ImageFormat_IMG_FORMAT_BASE64 {
		t.Errorf("Format = %v, want IMG_FORMAT_BASE64", req.Format)
	}
	if req.Image == nil {
		t.Fatal("Image is nil")
	}
	if req.Image.ImageURL != "http://example.com/img.jpg" {
		t.Errorf("ImageURL = %v, want http://example.com/img.jpg", req.Image.ImageURL)
	}
}

func TestGenerate(t *testing.T) {
	tests := []struct {
		name       string
		req        *GenerateRequest
		statusCode int
		response   *xaiv1.ImageResponse
		wantErr    bool
	}{
		{
			name:       "success",
			req:        NewRequest("A cat", "image-model"),
			statusCode: http.StatusOK,
			response: &xaiv1.ImageResponse{
				Images: []*xaiv1.GeneratedImage{
					{
						Url:               "http://example.com/cat.jpg",
						Base64:            "",
						UpSampledPrompt:   "A detailed cat",
						RespectModeration: true,
					},
				},
				Model: "image-model",
			},
			wantErr: false,
		},
		{
			name:       "multiple images",
			req:        NewRequest("A dog", "image-model").WithCount(2),
			statusCode: http.StatusOK,
			response: &xaiv1.ImageResponse{
				Images: []*xaiv1.GeneratedImage{
					{Url: "http://example.com/dog1.jpg"},
					{Url: "http://example.com/dog2.jpg"},
				},
				Model: "image-model",
			},
			wantErr: false,
		},
		{
			name:       "server error",
			req:        NewRequest("test", "model"),
			statusCode: http.StatusInternalServerError,
			response:   nil,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Errorf("Expected POST, got %s", r.Method)
				}
				if r.URL.Path != "/images/generations" {
					t.Errorf("Expected /images/generations, got %s", r.URL.Path)
				}

				w.WriteHeader(tt.statusCode)
				if tt.response != nil {
					json.NewEncoder(w).Encode(tt.response)
				}
			}))
			defer server.Close()

			restClient := rest.NewClient(rest.Config{
				BaseURL: server.URL,
				APIKey:  "test",
			})
			client := NewClient(restClient)

			resp, err := client.Generate(context.Background(), tt.req)

			if (err != nil) != tt.wantErr {
				t.Errorf("Generate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if resp == nil {
					t.Fatal("Generate() returned nil response")
				}
				if len(resp.Images) != len(tt.response.Images) {
					t.Errorf("Generate() returned %d images, want %d", len(resp.Images), len(tt.response.Images))
				}
				if resp.Model != tt.response.Model {
					t.Errorf("Model = %v, want %v", resp.Model, tt.response.Model)
				}
			}
		})
	}
}

func TestGenerateClientNotInitialized(t *testing.T) {
	client := &Client{restClient: nil}
	req := NewRequest("test", "model")

	_, err := client.Generate(context.Background(), req)
	if err != ErrClientNotInitialized {
		t.Errorf("Generate() error = %v, want %v", err, ErrClientNotInitialized)
	}
}

func TestGeneratedImage(t *testing.T) {
	img := &GeneratedImage{
		Base64:            "base64data",
		URL:               "http://example.com/img.jpg",
		UpsampledPrompt:   "detailed prompt",
		RespectModeration: true,
	}

	if img.Base64 != "base64data" {
		t.Errorf("Base64 = %v, want base64data", img.Base64)
	}
	if img.URL != "http://example.com/img.jpg" {
		t.Errorf("URL = %v, want http://example.com/img.jpg", img.URL)
	}
	if img.UpsampledPrompt != "detailed prompt" {
		t.Errorf("UpsampledPrompt = %v, want detailed prompt", img.UpsampledPrompt)
	}
	if !img.RespectModeration {
		t.Error("RespectModeration should be true")
	}
}

func TestResponse(t *testing.T) {
	resp := &Response{
		Images: []*GeneratedImage{
			{URL: "http://example.com/1.jpg"},
			{URL: "http://example.com/2.jpg"},
		},
		Model: "image-model",
	}

	if len(resp.Images) != 2 {
		t.Errorf("len(Images) = %v, want 2", len(resp.Images))
	}
	if resp.Model != "image-model" {
		t.Errorf("Model = %v, want image-model", resp.Model)
	}
}
