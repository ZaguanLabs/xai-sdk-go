package sample

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

func TestNewRequest(t *testing.T) {
	req := NewRequest("grok-1", "Hello", "World")

	if req == nil {
		t.Fatal("NewRequest returned nil")
	}
	if req.Model != "grok-1" {
		t.Errorf("Model = %v, want grok-1", req.Model)
	}
	if len(req.Prompts) != 2 {
		t.Errorf("len(Prompts) = %v, want 2", len(req.Prompts))
	}
	if req.MaxTokens != 100 {
		t.Errorf("MaxTokens = %v, want 100", req.MaxTokens)
	}
	if req.Temperature != 1.0 {
		t.Errorf("Temperature = %v, want 1.0", req.Temperature)
	}
}

func TestRequestBuilders(t *testing.T) {
	req := NewRequest("grok-1", "test").
		WithMaxTokens(200).
		WithTemperature(0.5)

	if req.MaxTokens != 200 {
		t.Errorf("MaxTokens = %v, want 200", req.MaxTokens)
	}
	if req.Temperature != 0.5 {
		t.Errorf("Temperature = %v, want 0.5", req.Temperature)
	}
}

func TestSample(t *testing.T) {
	tests := []struct {
		name       string
		req        *Request
		statusCode int
		response   *xaiv1.SampleTextResponse
		wantErr    bool
	}{
		{
			name:       "success",
			req:        NewRequest("grok-1", "Hello"),
			statusCode: http.StatusOK,
			response: &xaiv1.SampleTextResponse{
				Choices: []*xaiv1.SampleChoice{
					{
						FinishReason: xaiv1.FinishReason_REASON_STOP,
						Index:        0,
						Text:         "World",
					},
				},
				Model: "grok-1",
			},
			wantErr: false,
		},
		{
			name:       "multiple choices",
			req:        NewRequest("grok-1", "Test").WithMaxTokens(50),
			statusCode: http.StatusOK,
			response: &xaiv1.SampleTextResponse{
				Choices: []*xaiv1.SampleChoice{
					{FinishReason: xaiv1.FinishReason_REASON_STOP, Index: 0, Text: "Response 1"},
					{FinishReason: xaiv1.FinishReason_REASON_STOP, Index: 1, Text: "Response 2"},
				},
				Model: "grok-1",
			},
			wantErr: false,
		},
		{
			name:       "server error",
			req:        NewRequest("grok-1", "test"),
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
				if r.URL.Path != "/completions" {
					t.Errorf("Expected /completions, got %s", r.URL.Path)
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

			resp, err := client.Sample(context.Background(), tt.req)

			if (err != nil) != tt.wantErr {
				t.Errorf("Sample() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if resp == nil {
					t.Fatal("Sample() returned nil response")
				}
				if len(resp.Choices) != len(tt.response.Choices) {
					t.Errorf("Sample() returned %d choices, want %d", len(resp.Choices), len(tt.response.Choices))
				}
				if resp.Model != tt.response.Model {
					t.Errorf("Model = %v, want %v", resp.Model, tt.response.Model)
				}
			}
		})
	}
}

func TestSampleClientNotInitialized(t *testing.T) {
	client := &Client{restClient: nil}
	req := NewRequest("grok-1", "test")

	_, err := client.Sample(context.Background(), req)
	if err != ErrClientNotInitialized {
		t.Errorf("Sample() error = %v, want %v", err, ErrClientNotInitialized)
	}
}

func TestChoice(t *testing.T) {
	choice := &Choice{
		FinishReason: "stop",
		Index:        0,
		Text:         "Hello world",
	}

	if choice.FinishReason != "stop" {
		t.Errorf("FinishReason = %v, want stop", choice.FinishReason)
	}
	if choice.Index != 0 {
		t.Errorf("Index = %v, want 0", choice.Index)
	}
	if choice.Text != "Hello world" {
		t.Errorf("Text = %v, want Hello world", choice.Text)
	}
}
