package xai

import (
	"crypto/tls"
	"os"
	"testing"
	"time"

	"github.com/ZaguanLabs/xai-sdk-go/xai/internal/constants"
	"github.com/ZaguanLabs/xai-sdk-go/xai/internal/errors"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	// Test that defaults are set correctly
	if config.Host != constants.DefaultAPIV1Host {
		t.Errorf("Expected host %s, got %s", constants.DefaultAPIV1Host, config.Host)
	}

	if config.GRPCPort != constants.DefaultGRPCPort {
		t.Errorf("Expected gRPC port %s, got %s", constants.DefaultGRPCPort, config.GRPCPort)
	}

	if config.HTTPHost != constants.DefaultHTTPHost {
		t.Errorf("Expected HTTP host %s, got %s", constants.DefaultHTTPHost, config.HTTPHost)
	}

	if config.Timeout != constants.DefaultTimeout {
		t.Errorf("Expected timeout %v, got %v", constants.DefaultTimeout, config.Timeout)
	}

	if config.MaxRetries != constants.DefaultMaxRetries {
		t.Errorf("Expected max retries %d, got %d", constants.DefaultMaxRetries, config.MaxRetries)
	}

	if config.Environment != "production" {
		t.Errorf("Expected environment 'production', got %s", config.Environment)
	}

	if config.EnableTelemetry != true {
		t.Errorf("Expected telemetry enabled, got %t", config.EnableTelemetry)
	}
}

func TestNewConfig(t *testing.T) {
	config := NewConfig()

	// Should have default values
	if config == nil {
		t.Error("NewConfig should not return nil")
	}

	if config.Host == "" {
		t.Error("Config host should not be empty")
	}

	// Test that environment variables are loaded
	os.Setenv("XAI_API_KEY", "test-key")
	defer os.Unsetenv("XAI_API_KEY")

	config = NewConfig()
	if config.APIKey != "test-key" {
		t.Errorf("Expected API key 'test-key', got %s", config.APIKey)
	}
}

func TestNewConfigWithAPIKey(t *testing.T) {
	apiKey := "test-api-key"
	config := NewConfigWithAPIKey(apiKey)

	if config == nil {
		t.Error("NewConfigWithAPIKey should not return nil")
	}

	if config.APIKey != apiKey {
		t.Errorf("Expected API key %s, got %s", apiKey, config.APIKey)
	}

	// Should still have other defaults
	if config.Host != constants.DefaultAPIV1Host {
		t.Errorf("Expected host %s, got %s", constants.DefaultAPIV1Host, config.Host)
	}
}

func TestConfigLoadFromEnvironment(t *testing.T) {
	// Set environment variables
	os.Setenv("XAI_API_KEY", "env-api-key")
	os.Setenv("XAI_HOST", "custom.host")
	os.Setenv("XAI_TIMEOUT", "60")
	os.Setenv("XAI_INSECURE", "true")
	os.Setenv("XAI_MAX_RETRIES", "5")
	defer func() {
		os.Unsetenv("XAI_API_KEY")
		os.Unsetenv("XAI_HOST")
		os.Unsetenv("XAI_TIMEOUT")
		os.Unsetenv("XAI_INSECURE")
		os.Unsetenv("XAI_MAX_RETRIES")
	}()

	config := DefaultConfig()
	config.LoadFromEnvironment()

	if config.APIKey != "env-api-key" {
		t.Errorf("Expected API key 'env-api-key', got %s", config.APIKey)
	}

	if config.Host != "custom.host" {
		t.Errorf("Expected host 'custom.host', got %s", config.Host)
	}

	if config.Timeout != 60*time.Second {
		t.Errorf("Expected timeout 60s, got %v", config.Timeout)
	}

	if !config.Insecure {
		t.Error("Expected insecure to be true")
	}

	if config.MaxRetries != 5 {
		t.Errorf("Expected max retries 5, got %d", config.MaxRetries)
	}
}

func TestConfigValidate(t *testing.T) {
	t.Run("ValidConfig", func(t *testing.T) {
		config := NewConfigWithAPIKey("test-key")
		err := config.Validate()
		if err != nil {
			t.Errorf("Valid config should not return error: %v", err)
		}
	})

	t.Run("MissingAPIKey", func(t *testing.T) {
		config := DefaultConfig()
		config.APIKey = ""
		err := config.Validate()
		if err == nil {
			t.Error("Missing API key should return error")
		}
		if _, ok := err.(*errors.Error); !ok {
			t.Error("Error should be of type *errors.Error")
		}
	})

	t.Run("InvalidTimeouts", func(t *testing.T) {
		config := NewConfigWithAPIKey("test-key")
		config.Timeout = 0
		err := config.Validate()
		if err == nil {
			t.Error("Zero timeout should return error")
		}

		config.KeepAliveTimeout = -1
		err = config.Validate()
		if err == nil {
			t.Error("Negative keep-alive timeout should return error")
		}
	})

	t.Run("InvalidRetrySettings", func(t *testing.T) {
		config := NewConfigWithAPIKey("test-key")
		config.MaxRetries = -1
		err := config.Validate()
		if err == nil {
			t.Error("Negative max retries should return error")
		}

		config.RetryBackoff = 0
		err = config.Validate()
		if err == nil {
			t.Error("Zero retry backoff should return error")
		}

		config.RetryBackoff = 5 * time.Second
		config.MaxBackoff = 3 * time.Second // Less than retry backoff
		err = config.Validate()
		if err == nil {
			t.Error("Max backoff less than retry backoff should return error")
		}
	})
}

func TestConfigGRPCAddress(t *testing.T) {
	config := &Config{
		Host:     "localhost",
		GRPCPort: "8080",
	}

	addr := config.GRPCAddress()
	expected := "localhost:8080"
	if addr != expected {
		t.Errorf("Expected gRPC address %s, got %s", expected, addr)
	}
}

func TestConfigHTTPAddress(t *testing.T) {
	config := &Config{
		HTTPHost: "localhost",
		HTTPPort: "8080",
	}

	addr := config.HTTPAddress()
	expected := "localhost:8080"
	if addr != expected {
		t.Errorf("Expected HTTP address %s, got %s", expected, addr)
	}
}


func TestConfigCreateGRPCDialOptions(t *testing.T) {
	t.Run("InsecureConnection", func(t *testing.T) {
		config := &Config{
			Insecure: true,
		}

		opts, err := config.CreateGRPCDialOptions()
		if err != nil {
			t.Errorf("Should not return error for insecure connection: %v", err)
		}

		if len(opts) == 0 {
			t.Error("Should have at least one dial option")
		}
	})

	t.Run("SecureConnection", func(t *testing.T) {
		config := &Config{
			Insecure:  false,
			SkipVerify: true,
		}

		opts, err := config.CreateGRPCDialOptions()
		if err != nil {
			t.Errorf("Should not return error for secure connection: %v", err)
		}

		if len(opts) == 0 {
			t.Error("Should have at least one dial option")
		}
	})

	t.Run("CustomTLS", func(t *testing.T) {
		customTLS := &tls.Config{
			MinVersion: tls.VersionTLS13,
		}
		config := &Config{
			CustomTLSConfig: customTLS,
		}

		opts, err := config.CreateGRPCDialOptions()
		if err != nil {
			t.Errorf("Should not return error for custom TLS: %v", err)
		}

		if len(opts) == 0 {
			t.Error("Should have at least one dial option")
		}
	})
}

func TestConfigWithMethods(t *testing.T) {
	config := DefaultConfig()

	// Test WithAPIKey
	config = config.WithAPIKey("new-api-key")
	if config.APIKey != "new-api-key" {
		t.Errorf("Expected API key 'new-api-key', got %s", config.APIKey)
	}

	// Test WithHost
	config = config.WithHost("new.host")
	if config.Host != "new.host" {
		t.Errorf("Expected host 'new.host', got %s", config.Host)
	}

	// Test WithTimeout
	newTimeout := 45 * time.Second
	config = config.WithTimeout(newTimeout)
	if config.Timeout != newTimeout {
		t.Errorf("Expected timeout %v, got %v", newTimeout, config.Timeout)
	}

	// Test WithInsecure
	config = config.WithInsecure(true)
	if !config.Insecure {
		t.Error("Expected insecure to be true")
	}

	// Test WithEnvironment
	config = config.WithEnvironment("test")
	if config.Environment != "test" {
		t.Errorf("Expected environment 'test', got %s", config.Environment)
	}

	// Test WithUserAgent
	config = config.WithUserAgent("test-agent")
	if config.UserAgent != "test-agent" {
		t.Errorf("Expected user agent 'test-agent', got %s", config.UserAgent)
	}

	// Test WithMaxRetries
	config = config.WithMaxRetries(10)
	if config.MaxRetries != 10 {
		t.Errorf("Expected max retries 10, got %d", config.MaxRetries)
	}

	// Test WithEnableTelemetry
	config = config.WithEnableTelemetry(false)
	if config.EnableTelemetry {
		t.Error("Expected telemetry to be false")
	}
}

func TestConfigString(t *testing.T) {
	config := &Config{
		APIKey:         "1234567890abcdef",
		Host:           "api.example.com",
		GRPCPort:       "443",
		Insecure:       false,
		Environment:    "production",
		MaxRetries:     3,
		EnableTelemetry: true,
	}

	str := config.String()

	// Should contain non-sensitive information
	if !configContainsString(str, "api.example.com") {
		t.Errorf("Config string should contain host: %s", str)
	}

	if !configContainsString(str, "production") {
		t.Errorf("Config string should contain environment: %s", str)
	}

	// Should mask API key
	if !configContainsString(str, "***") {
		t.Errorf("Config string should mask API key: %s", str)
	}
}

func TestParseDuration(t *testing.T) {
	tests := []struct {
		input    string
		expected time.Duration
		hasError bool
	}{
		{"30", 30 * time.Second, false},
		{"60s", 60 * time.Second, false},
		{"2m", 2 * time.Minute, false},
		{"1h", 1 * time.Hour, false},
		{"invalid", 0, true},
		{"", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result, err := parseDuration(tt.input)

			if tt.hasError {
				if err == nil {
					t.Error("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
				}
				if result != tt.expected {
					t.Errorf("Expected duration %v, got %v", tt.expected, result)
				}
			}
		})
	}
}

func BenchmarkConfigValidate(b *testing.B) {
	config := NewConfigWithAPIKey("benchmark-api-key")
	for i := 0; i < b.N; i++ {
		_ = config.Validate()
	}
}

func BenchmarkConfigCreateGRPCDialOptions(b *testing.B) {
	config := NewConfigWithAPIKey("test-api-key")
	for i := 0; i < b.N; i++ {
		_, _ = config.CreateGRPCDialOptions()
	}
}

// Helper function to check if a string contains a substring
func configContainsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || (len(s) > len(substr) && configFindString(s, substr)))
}

func configFindString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}