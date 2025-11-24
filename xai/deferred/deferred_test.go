package deferred

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

func TestStart(t *testing.T) {
	tests := []struct {
		name       string
		request    interface{}
		statusCode int
		response   *xaiv1.StartDeferredResponse
		wantErr    bool
	}{
		{
			name:       "success",
			request:    map[string]string{"prompt": "test"},
			statusCode: http.StatusOK,
			response: &xaiv1.StartDeferredResponse{
				RequestId: "req-123",
			},
			wantErr: false,
		},
		{
			name:       "server error",
			request:    map[string]string{"prompt": "test"},
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
				if r.URL.Path != "/deferred/start" {
					t.Errorf("Expected /deferred/start, got %s", r.URL.Path)
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

			resp, err := client.Start(context.Background(), tt.request)

			if (err != nil) != tt.wantErr {
				t.Errorf("Start() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if resp == nil {
					t.Fatal("Start() returned nil response")
				}
				if resp.RequestID != tt.response.RequestId {
					t.Errorf("RequestID = %v, want %v", resp.RequestID, tt.response.RequestId)
				}
			}
		})
	}
}

func TestGet(t *testing.T) {
	tests := []struct {
		name       string
		requestID  string
		statusCode int
		response   map[string]interface{}
		wantErr    bool
	}{
		{
			name:       "success - pending",
			requestID:  "req-123",
			statusCode: http.StatusOK,
			response: map[string]interface{}{
				"request_id": "req-123",
				"status":     int(xaiv1.DeferredStatus_PENDING),
			},
			wantErr: false,
		},
		{
			name:       "success - done",
			requestID:  "req-456",
			statusCode: http.StatusOK,
			response: map[string]interface{}{
				"request_id": "req-456",
				"status":     int(xaiv1.DeferredStatus_DONE),
				"result":     map[string]string{"text": "completion result"},
			},
			wantErr: false,
		},
		{
			name:       "server error",
			requestID:  "req-999",
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
				if r.URL.Path != "/deferred/get" {
					t.Errorf("Expected /deferred/get, got %s", r.URL.Path)
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

			status, err := client.Get(context.Background(), tt.requestID)

			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if status == nil {
					t.Fatal("Get() returned nil status")
				}
				if status.RequestID != tt.requestID {
					t.Errorf("RequestID = %v, want %v", status.RequestID, tt.requestID)
				}
			}
		})
	}
}

func TestStartClientNotInitialized(t *testing.T) {
	client := &Client{restClient: nil}

	_, err := client.Start(context.Background(), map[string]string{})
	if err != ErrClientNotInitialized {
		t.Errorf("Start() error = %v, want %v", err, ErrClientNotInitialized)
	}
}

func TestGetClientNotInitialized(t *testing.T) {
	client := &Client{restClient: nil}

	_, err := client.Get(context.Background(), "req-123")
	if err != ErrClientNotInitialized {
		t.Errorf("Get() error = %v, want %v", err, ErrClientNotInitialized)
	}
}
