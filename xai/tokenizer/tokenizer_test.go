package tokenizer

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

func TestTokenize(t *testing.T) {
	tests := []struct {
		name       string
		text       string
		model      string
		user       string
		statusCode int
		response   *xaiv1.TokenizeTextResponse
		wantErr    bool
	}{
		{
			name:       "success",
			text:       "Hello world",
			model:      "grok-1",
			user:       "test-user",
			statusCode: http.StatusOK,
			response: &xaiv1.TokenizeTextResponse{
				Tokens: []*xaiv1.Token{
					{TokenId: 1, StringToken: "Hello", TokenBytes: []byte("Hello")},
					{TokenId: 2, StringToken: "world", TokenBytes: []byte("world")},
				},
				Model: "grok-1",
			},
			wantErr: false,
		},
		{
			name:       "empty text",
			text:       "",
			model:      "grok-1",
			user:       "",
			statusCode: http.StatusOK,
			response: &xaiv1.TokenizeTextResponse{
				Tokens: []*xaiv1.Token{},
				Model:  "grok-1",
			},
			wantErr: false,
		},
		{
			name:       "server error",
			text:       "test",
			model:      "grok-1",
			user:       "",
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
				if r.URL.Path != "/tokenize" {
					t.Errorf("Expected /tokenize, got %s", r.URL.Path)
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

			resp, err := client.Tokenize(context.Background(), tt.text, tt.model, tt.user)

			if (err != nil) != tt.wantErr {
				t.Errorf("Tokenize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if resp == nil {
					t.Fatal("Tokenize() returned nil response")
				}
				if len(resp.Tokens) != len(tt.response.Tokens) {
					t.Errorf("Tokenize() returned %d tokens, want %d", len(resp.Tokens), len(tt.response.Tokens))
				}
				if resp.Model != tt.response.Model {
					t.Errorf("Model = %v, want %v", resp.Model, tt.response.Model)
				}
			}
		})
	}
}

func TestTokenizeClientNotInitialized(t *testing.T) {
	client := &Client{restClient: nil}

	_, err := client.Tokenize(context.Background(), "test", "grok-1", "")
	if err != ErrClientNotInitialized {
		t.Errorf("Tokenize() error = %v, want %v", err, ErrClientNotInitialized)
	}
}

func TestToken(t *testing.T) {
	token := &Token{
		TokenID:     123,
		StringToken: "test",
		TokenBytes:  []byte("test"),
	}

	if token.TokenID != 123 {
		t.Errorf("TokenID = %v, want 123", token.TokenID)
	}
	if token.StringToken != "test" {
		t.Errorf("StringToken = %v, want test", token.StringToken)
	}
	if string(token.TokenBytes) != "test" {
		t.Errorf("TokenBytes = %v, want test", string(token.TokenBytes))
	}
}

func TestResponse(t *testing.T) {
	resp := &Response{
		Tokens: []*Token{
			{TokenID: 1, StringToken: "hello"},
			{TokenID: 2, StringToken: "world"},
		},
		Model: "grok-1",
	}

	if len(resp.Tokens) != 2 {
		t.Errorf("len(Tokens) = %v, want 2", len(resp.Tokens))
	}
	if resp.Model != "grok-1" {
		t.Errorf("Model = %v, want grok-1", resp.Model)
	}
}
