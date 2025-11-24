package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	xaiv1 "github.com/ZaguanLabs/xai-sdk-go/proto/gen/go/xai/api/v1"
	"github.com/ZaguanLabs/xai-sdk-go/xai/internal/rest"
	"google.golang.org/protobuf/types/known/timestamppb"
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

func TestFromProto(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name  string
		input *xaiv1.ApiKey
		want  *ApiKey
	}{
		{
			name:  "nil input",
			input: nil,
			want:  nil,
		},
		{
			name: "full api key",
			input: &xaiv1.ApiKey{
				RedactedApiKey: "xai-***123",
				UserId:         "user-123",
				Name:           "Test Key",
				TeamId:         "team-456",
				Acls:           []string{"read", "write"},
				ApiKeyId:       "key-789",
				ApiKeyBlocked:  false,
				ModifiedBy:     "admin",
				Disabled:       false,
				TeamBlocked:    false,
				CreateTime:     timestamppb.New(now),
				ModifyTime:     timestamppb.New(now),
			},
			want: &ApiKey{
				RedactedApiKey: "xai-***123",
				UserID:         "user-123",
				Name:           "Test Key",
				TeamID:         "team-456",
				ACLs:           []string{"read", "write"},
				ApiKeyID:       "key-789",
				ApiKeyBlocked:  false,
				ModifiedBy:     "admin",
				Disabled:       false,
				TeamBlocked:    false,
				CreateTime:     now,
				ModifyTime:     now,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := fromProto(tt.input)

			if tt.want == nil {
				if got != nil {
					t.Errorf("fromProto() = %v, want nil", got)
				}
				return
			}

			if got == nil {
				t.Fatal("fromProto() returned nil")
			}

			if got.RedactedApiKey != tt.want.RedactedApiKey {
				t.Errorf("RedactedApiKey = %v, want %v", got.RedactedApiKey, tt.want.RedactedApiKey)
			}
			if got.UserID != tt.want.UserID {
				t.Errorf("UserID = %v, want %v", got.UserID, tt.want.UserID)
			}
			if got.ApiKeyID != tt.want.ApiKeyID {
				t.Errorf("ApiKeyID = %v, want %v", got.ApiKeyID, tt.want.ApiKeyID)
			}
		})
	}
}

func TestValidateKey(t *testing.T) {
	tests := []struct {
		name       string
		apiKey     string
		statusCode int
		response   *xaiv1.ApiKey
		wantErr    bool
	}{
		{
			name:       "valid key",
			apiKey:     "xai-test-key",
			statusCode: http.StatusOK,
			response: &xaiv1.ApiKey{
				RedactedApiKey: "xai-***key",
				UserId:         "user-123",
				ApiKeyId:       "key-456",
			},
			wantErr: false,
		},
		{
			name:       "invalid key",
			apiKey:     "invalid",
			statusCode: http.StatusUnauthorized,
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
				if r.URL.Path != "/auth/validate" {
					t.Errorf("Expected /auth/validate, got %s", r.URL.Path)
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

			key, err := client.ValidateKey(context.Background(), tt.apiKey)

			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && key == nil {
				t.Error("ValidateKey() returned nil key")
			}
		})
	}
}

func TestGetKey(t *testing.T) {
	tests := []struct {
		name       string
		keyID      string
		statusCode int
		response   *xaiv1.ApiKey
		wantErr    bool
	}{
		{
			name:       "existing key",
			keyID:      "key-123",
			statusCode: http.StatusOK,
			response: &xaiv1.ApiKey{
				ApiKeyId:       "key-123",
				RedactedApiKey: "xai-***123",
			},
			wantErr: false,
		},
		{
			name:       "non-existent key",
			keyID:      "key-999",
			statusCode: http.StatusNotFound,
			response:   nil,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Errorf("Expected GET, got %s", r.Method)
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

			key, err := client.GetKey(context.Background(), tt.keyID)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && key == nil {
				t.Error("GetKey() returned nil key")
			}
		})
	}
}

func TestListKeys(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		response   map[string]interface{}
		wantCount  int
		wantErr    bool
	}{
		{
			name:       "multiple keys",
			statusCode: http.StatusOK,
			response: map[string]interface{}{
				"keys": []map[string]interface{}{
					{"api_key_id": "key-1", "redacted_api_key": "xai-***1"},
					{"api_key_id": "key-2", "redacted_api_key": "xai-***2"},
				},
			},
			wantCount: 2,
			wantErr:   false,
		},
		{
			name:       "empty list",
			statusCode: http.StatusOK,
			response: map[string]interface{}{
				"keys": []map[string]interface{}{},
			},
			wantCount: 0,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Errorf("Expected GET, got %s", r.Method)
				}
				if r.URL.Path != "/auth/keys" {
					t.Errorf("Expected /auth/keys, got %s", r.URL.Path)
				}

				w.WriteHeader(tt.statusCode)
				json.NewEncoder(w).Encode(tt.response)
			}))
			defer server.Close()

			restClient := rest.NewClient(rest.Config{
				BaseURL: server.URL,
				APIKey:  "test",
			})
			client := NewClient(restClient)

			keys, err := client.ListKeys(context.Background())

			if (err != nil) != tt.wantErr {
				t.Errorf("ListKeys() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && len(keys) != tt.wantCount {
				t.Errorf("ListKeys() returned %d keys, want %d", len(keys), tt.wantCount)
			}
		})
	}
}

func TestClientNotInitialized(t *testing.T) {
	client := &Client{restClient: nil}

	_, err := client.ValidateKey(context.Background(), "test")
	if err != ErrClientNotInitialized {
		t.Errorf("ValidateKey() error = %v, want %v", err, ErrClientNotInitialized)
	}

	_, err = client.GetKey(context.Background(), "test")
	if err != ErrClientNotInitialized {
		t.Errorf("GetKey() error = %v, want %v", err, ErrClientNotInitialized)
	}

	_, err = client.ListKeys(context.Background())
	if err != ErrClientNotInitialized {
		t.Errorf("ListKeys() error = %v, want %v", err, ErrClientNotInitialized)
	}
}
