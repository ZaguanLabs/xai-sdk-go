package embed

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	xaiv1 "github.com/ZaguanLabs/xai-sdk-go/proto/gen/go/xai/api/v1"
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

func TestText(t *testing.T) {
	input := Text("hello world")

	if input.proto == nil {
		t.Fatal("Text() returned nil proto")
	}
	if input.proto.GetString_() != "hello world" {
		t.Errorf("String_ = %v, want hello world", input.proto.GetString_())
	}
}

func TestImage(t *testing.T) {
	input := Image("http://example.com/image.jpg", xaiv1.ImageDetail_DETAIL_HIGH)

	if input.proto == nil {
		t.Fatal("Image() returned nil proto")
	}
	if input.proto.GetImageUrl() == nil {
		t.Fatal("ImageUrl is nil")
	}
	if input.proto.GetImageUrl().ImageUrl != "http://example.com/image.jpg" {
		t.Errorf("ImageUrl = %v, want http://example.com/image.jpg", input.proto.GetImageUrl().ImageUrl)
	}
}

func TestNewRequest(t *testing.T) {
	req := NewRequest("embed-model", Text("test1"), Text("test2"))

	if req == nil {
		t.Fatal("NewRequest returned nil")
	}
	if req.proto == nil {
		t.Fatal("proto is nil")
	}
	if req.proto.Model != "embed-model" {
		t.Errorf("Model = %v, want embed-model", req.proto.Model)
	}
	if len(req.proto.Input) != 2 {
		t.Errorf("len(Input) = %v, want 2", len(req.proto.Input))
	}
}

func TestRequestBuilders(t *testing.T) {
	req := NewRequest("embed-model", Text("test")).
		WithEncodingFormat(xaiv1.EmbedEncodingFormat_FORMAT_FLOAT).
		WithUser("test-user")

	if req.proto.EncodingFormat != xaiv1.EmbedEncodingFormat_FORMAT_FLOAT {
		t.Errorf("EncodingFormat = %v, want FLOAT", req.proto.EncodingFormat)
	}
	if req.proto.User != "test-user" {
		t.Errorf("User = %v, want test-user", req.proto.User)
	}
}

func TestGenerate(t *testing.T) {
	tests := []struct {
		name       string
		req        *Request
		statusCode int
		response   *xaiv1.EmbedResponse
		wantErr    bool
	}{
		{
			name:       "success",
			req:        NewRequest("embed-model", Text("hello")),
			statusCode: http.StatusOK,
			response: &xaiv1.EmbedResponse{
				Id:    "emb-123",
				Model: "embed-model",
				Embeddings: []*xaiv1.Embedding{
					{
						Index: 0,
						Embeddings: []*xaiv1.FeatureVector{
							{FloatArray: []float32{0.1, 0.2, 0.3}},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name:       "server error",
			req:        NewRequest("embed-model", Text("test")),
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
				if r.URL.Path != "/embeddings" {
					t.Errorf("Expected /embeddings, got %s", r.URL.Path)
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
				if resp.ID() != tt.response.Id {
					t.Errorf("ID() = %v, want %v", resp.ID(), tt.response.Id)
				}
				if resp.Model() != tt.response.Model {
					t.Errorf("Model() = %v, want %v", resp.Model(), tt.response.Model)
				}
			}
		})
	}
}

func TestGenerateClientNotInitialized(t *testing.T) {
	client := &Client{restClient: nil}
	req := NewRequest("embed-model", Text("test"))

	_, err := client.Generate(context.Background(), req)
	if err != ErrClientNotInitialized {
		t.Errorf("Generate() error = %v, want %v", err, ErrClientNotInitialized)
	}
}

func TestResponseMethods(t *testing.T) {
	resp := &Response{
		proto: &xaiv1.EmbedResponse{
			Id:                "emb-123",
			Model:             "embed-model",
			SystemFingerprint: "fp-456",
			Embeddings: []*xaiv1.Embedding{
				{Index: 0},
				{Index: 1},
			},
			Usage: &xaiv1.EmbeddingUsage{},
		},
	}

	if resp.ID() != "emb-123" {
		t.Errorf("ID() = %v, want emb-123", resp.ID())
	}
	if resp.Model() != "embed-model" {
		t.Errorf("Model() = %v, want embed-model", resp.Model())
	}
	if resp.SystemFingerprint() != "fp-456" {
		t.Errorf("SystemFingerprint() = %v, want fp-456", resp.SystemFingerprint())
	}
	if len(resp.Embeddings()) != 2 {
		t.Errorf("len(Embeddings()) = %v, want 2", len(resp.Embeddings()))
	}
	if resp.Usage() == nil {
		t.Error("Usage() returned nil")
	}
}

func TestResponseNilProto(t *testing.T) {
	resp := &Response{proto: nil}

	if resp.ID() != "" {
		t.Errorf("ID() = %v, want empty", resp.ID())
	}
	if resp.Model() != "" {
		t.Errorf("Model() = %v, want empty", resp.Model())
	}
	if resp.Embeddings() != nil {
		t.Error("Embeddings() should return nil")
	}
}

func TestEmbeddingMethods(t *testing.T) {
	emb := &Embedding{
		proto: &xaiv1.Embedding{
			Index: 5,
			Embeddings: []*xaiv1.FeatureVector{
				{FloatArray: []float32{0.1, 0.2}},
			},
		},
	}

	if emb.Index() != 5 {
		t.Errorf("Index() = %v, want 5", emb.Index())
	}
	if len(emb.Vectors()) != 1 {
		t.Errorf("len(Vectors()) = %v, want 1", len(emb.Vectors()))
	}
}

func TestFeatureVectorMethods(t *testing.T) {
	fv := &FeatureVector{
		proto: &xaiv1.FeatureVector{
			FloatArray:  []float32{0.1, 0.2, 0.3},
			Base64Array: "AQIDBA==",
		},
	}

	if len(fv.FloatArray()) != 3 {
		t.Errorf("len(FloatArray()) = %v, want 3", len(fv.FloatArray()))
	}
	if fv.Base64Array() != "AQIDBA==" {
		t.Errorf("Base64Array() = %v, want AQIDBA==", fv.Base64Array())
	}
}

func TestFeatureVectorNilProto(t *testing.T) {
	fv := &FeatureVector{proto: nil}

	if fv.FloatArray() != nil {
		t.Error("FloatArray() should return nil")
	}
	if fv.Base64Array() != "" {
		t.Error("Base64Array() should return empty")
	}
}
