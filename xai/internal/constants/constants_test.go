package constants

import (
	"strings"
	"testing"
	"time"

	"github.com/ZaguanLabs/xai-sdk-go/xai/internal/version"
)

func TestDefaultEndpoints(t *testing.T) {
	tests := []struct {
		name     string
		constant string
		expected string
	}{
		{"DefaultAPIV1Host", DefaultAPIV1Host, "api.x.ai"},
		{"DefaultGRPCPort", DefaultGRPCPort, "443"},
		{"DefaultHTTPHost", DefaultHTTPHost, "api.x.ai"},
		{"DefaultChatEndpoint", DefaultChatEndpoint, "/v1/chat/completions"},
		{"DefaultFilesEndpoint", DefaultFilesEndpoint, "/v1/files"},
		{"DefaultImageEndpoint", DefaultImageEndpoint, "/v1/images"},
		{"DefaultModelsEndpoint", DefaultModelsEndpoint, "/v1/models"},
		{"DefaultTokenizerEndpoint", DefaultTokenizerEndpoint, "/v1/tokenizer"},
		{"DefaultCollectionsEndpoint", DefaultCollectionsEndpoint, "/v1/collections"},
		{"DefaultAuthEndpoint", DefaultAuthEndpoint, "/v1/auth"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.constant != tt.expected {
				t.Errorf("Expected %s to be %s, got %s", tt.name, tt.expected, tt.constant)
			}
		})
	}
}

func TestDefaultTimeouts(t *testing.T) {
	timeoutTests := []struct {
		name     string
		duration time.Duration
		expected time.Duration
	}{
		{"DefaultTimeout", DefaultTimeout, 30 * time.Second},
		{"DefaultConnectTimeout", DefaultConnectTimeout, 10 * time.Second},
		{"DefaultKeepAliveTimeout", DefaultKeepAliveTimeout, 20 * time.Second},
		{"DefaultStreamTimeout", DefaultStreamTimeout, 300 * time.Second},
	}

	for _, tt := range timeoutTests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.duration != tt.expected {
				t.Errorf("Expected %s to be %v, got %v", tt.name, tt.expected, tt.duration)
			}
		})
	}
}

func TestDefaultSizeValues(t *testing.T) {
	sizeTests := []struct {
		name     string
		value    int
		expected int
	}{
		{"DefaultMaxRetries", DefaultMaxRetries, 3},
		{"DefaultChunkSize", DefaultChunkSize, 3 * 1024 * 1024},
		{"DefaultMaxChunkSize", DefaultMaxChunkSize, 10 * 1024 * 1024},
		{"DefaultBufferSize", DefaultBufferSize, 64 * 1024},
		{"DefaultMaxImageSize", DefaultMaxImageSize, 10 * 1024 * 1024},
		{"DefaultMaxFileSize", DefaultMaxFileSize, 100 * 1024 * 1024},
		{"DefaultMaxTokens", DefaultMaxTokens, 4096},
		{"DefaultMaxPromptTokens", DefaultMaxPromptTokens, 8192},
		{"DefaultMaxToolResults", DefaultMaxToolResults, 50},
		{"DefaultMaxToolDescriptions", DefaultMaxToolDescriptions, 200},
	}

	for _, tt := range sizeTests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.value != tt.expected {
				t.Errorf("Expected %s to be %d, got %d", tt.name, tt.expected, tt.value)
			}
		})
	}
}

func TestRetryBackoff(t *testing.T) {
	if DefaultRetryBackoff <= 0 {
		t.Error("DefaultRetryBackoff should be positive")
	}

	if DefaultMaxBackoff <= DefaultRetryBackoff {
		t.Error("DefaultMaxBackoff should be greater than DefaultRetryBackoff")
	}

	if DefaultMaxBackoff <= 0 {
		t.Error("DefaultMaxBackoff should be positive")
	}
}

func TestHeaderNames(t *testing.T) {
	headerTests := []struct {
		name     string
		constant string
		expected string
	}{
		{"HeaderAuthorization", HeaderAuthorization, "Authorization"},
		{"HeaderContentType", HeaderContentType, "Content-Type"},
		{"HeaderUserAgent", HeaderUserAgent, "User-Agent"},
		{"HeaderXAIClient", HeaderXAIClient, "X-AI-Client"},
		{"HeaderXAIRequestID", HeaderXAIRequestID, "X-AI-Request-ID"},
		{"HeaderXAIAPIVersion", HeaderXAIAPIVersion, "X-AI-API-Version"},
		{"HeaderXAIEnvironment", HeaderXAIEnvironment, "X-AI-Environment"},
		{"HeaderXAIConversationID", HeaderXAIConversationID, "X-AI-Conversation-ID"},
	}

	for _, tt := range headerTests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.constant != tt.expected {
				t.Errorf("Expected %s to be %s, got %s", tt.name, tt.expected, tt.constant)
			}
		})
	}
}

func TestContentTypes(t *testing.T) {
	contentTypeTests := []struct {
		name     string
		constant string
		expected string
	}{
		{"ContentTypeJSON", ContentTypeJSON, "application/json"},
		{"ContentTypeProtobuf", ContentTypeProtobuf, "application/x-protobuf"},
		{"ContentTypeOctetStream", ContentTypeOctetStream, "application/octet-stream"},
		{"ContentTypeTextPlain", ContentTypeTextPlain, "text/plain"},
		{"ContentTypeImageJPEG", ContentTypeImageJPEG, "image/jpeg"},
		{"ContentTypeImagePNG", ContentTypeImagePNG, "image/png"},
	}

	for _, tt := range contentTypeTests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.constant != tt.expected {
				t.Errorf("Expected %s to be %s, got %s", tt.name, tt.expected, tt.constant)
			}
		})
	}
}

func TestDefaultUserAgent(t *testing.T) {
	if DefaultUserAgent == "" {
		t.Error("DefaultUserAgent should not be empty")
	}

	if !strings.Contains(DefaultUserAgent, "xai-sdk-go") {
		t.Errorf("DefaultUserAgent should contain 'xai-sdk-go', got: %s", DefaultUserAgent)
	}

	// Verify it contains the version from the version package
	if !strings.Contains(DefaultUserAgent, version.SDKVersion) {
		t.Errorf("DefaultUserAgent should contain version '%s', got: %s", version.SDKVersion, DefaultUserAgent)
	}
}

func TestConstantsConsistency(t *testing.T) {
	// Test that related constants are consistent
	if DefaultGRPCPort != "443" {
		t.Errorf("DefaultGRPCPort should be 443 for TLS, got %s", DefaultGRPCPort)
	}

	if DefaultHTTPHost != DefaultAPIV1Host {
		t.Errorf("DefaultHTTPHost should match DefaultAPIV1Host for consistency")
	}

	// Test timeout relationships
	if DefaultConnectTimeout >= DefaultTimeout {
		t.Error("DefaultConnectTimeout should be less than DefaultTimeout")
	}

	if DefaultKeepAliveTimeout > DefaultTimeout {
		t.Error("DefaultKeepAliveTimeout should be less than or equal to DefaultTimeout")
	}

	// Test size relationships
	if DefaultChunkSize >= DefaultMaxChunkSize {
		t.Error("DefaultChunkSize should be less than DefaultMaxChunkSize")
	}

	if DefaultMaxImageSize >= DefaultMaxFileSize {
		t.Error("DefaultMaxImageSize should be less than DefaultMaxFileSize")
	}
}

func TestPortConsistency(t *testing.T) {
	// Test that we have the standard ports
	grpcPort := DefaultGRPCPort
	if grpcPort != "443" && grpcPort != "80" {
		t.Errorf("DefaultGRPCPort should be a standard port (80 or 443), got: %s", grpcPort)
	}
}

func TestEndpointPaths(t *testing.T) {
	endpoints := []string{
		DefaultChatEndpoint,
		DefaultFilesEndpoint,
		DefaultImageEndpoint,
		DefaultModelsEndpoint,
		DefaultTokenizerEndpoint,
		DefaultCollectionsEndpoint,
		DefaultAuthEndpoint,
	}

	for _, endpoint := range endpoints {
		if endpoint == "" {
			t.Error("Endpoint should not be empty")
		}

		if !strings.HasPrefix(endpoint, "/") {
			t.Errorf("Endpoint should start with '/', got: %s", endpoint)
		}

		if !strings.Contains(endpoint, "/v1/") {
			t.Errorf("Endpoint should contain '/v1/', got: %s", endpoint)
		}
	}
}

func TestBufferSizeReasonable(t *testing.T) {
	// Test that buffer sizes are reasonable
	if DefaultBufferSize < 1024 {
		t.Error("DefaultBufferSize should be at least 1KB")
	}

	if DefaultBufferSize > 1024*1024 {
		t.Error("DefaultBufferSize should not exceed 1MB for reasonable memory usage")
	}
}

func BenchmarkConstants(b *testing.B) {
	// Benchmark accessing constants
	for i := 0; i < b.N; i++ {
		_ = DefaultTimeout
		_ = DefaultChunkSize
		_ = HeaderAuthorization
		_ = ContentTypeJSON
		_ = DefaultUserAgent
	}
}
