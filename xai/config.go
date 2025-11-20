// Package xai provides the main client for interacting with xAI services.
package xai

import (
	"crypto/tls"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ZaguanLabs/xai-sdk-go/xai/internal/auth"
	"github.com/ZaguanLabs/xai-sdk-go/xai/internal/constants"
	"github.com/ZaguanLabs/xai-sdk-go/xai/internal/errors"
	"github.com/ZaguanLabs/xai-sdk-go/xai/internal/grpcutil"
	"github.com/ZaguanLabs/xai-sdk-go/xai/internal/metadata"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

// Config represents the client configuration.
type Config struct {
	// APIKey is the xAI API key for authentication.
	APIKey string `json:"api_key"`

	// Host is the xAI API host (default: api.x.ai).
	Host string `json:"host"`

	// GRPCPort is the gRPC port (default: 443).
	GRPCPort string `json:"grpc_port"`

	// HTTPHost is the HTTP API host (default: api.x.ai).
	HTTPHost string `json:"http_host"`

	// HTTPPort is the HTTP API port (default: 80).
	HTTPPort string `json:"http_port"`

	// Timeout is the default request timeout (default: 30s).
	Timeout time.Duration `json:"timeout"`

	// ConnectTimeout is the connection timeout (default: 10s).
	ConnectTimeout time.Duration `json:"connect_timeout"`

	// KeepAliveTimeout is the keep-alive timeout (default: 20s).
	KeepAliveTimeout time.Duration `json:"keep_alive_timeout"`

	// StreamTimeout is the streaming timeout (default: 300s).
	StreamTimeout time.Duration `json:"stream_timeout"`

	// Insecure controls whether to use TLS (default: false).
	// WARNING: Setting this to true disables TLS encryption entirely.
	// Only use in local development or testing environments. Never use in production.
	Insecure bool `json:"insecure"`

	// SkipVerify controls whether to skip TLS certificate verification (default: false).
	// WARNING: Setting this to true disables certificate validation, making connections
	// vulnerable to man-in-the-middle attacks. Only use in local development or testing
	// environments with self-signed certificates. Never use in production.
	SkipVerify bool `json:"skip_verify"`

	// MaxRetries is the maximum number of retries (default: 3).
	MaxRetries int `json:"max_retries"`

	// RetryBackoff is the retry backoff duration (default: 1s).
	RetryBackoff time.Duration `json:"retry_backoff"`

	// MaxBackoff is the maximum backoff duration (default: 60s).
	MaxBackoff time.Duration `json:"max_backoff"`

	// Environment is the deployment environment (default: production).
	Environment string `json:"environment"`

	// UserAgent is the user agent string (default: xai-sdk-go/version).
	UserAgent string `json:"user_agent"`

	// EnableTelemetry controls whether to enable telemetry (default: true).
	EnableTelemetry bool `json:"enable_telemetry"`

	// CustomTLSConfig allows providing a custom TLS configuration.
	CustomTLSConfig *tls.Config `json:"-"`

	// Logger is an optional logger (if nil, no logging is performed).
	// Logger logger.Logger `json:"-"`
}

// DefaultConfig returns a Config with default values.
func DefaultConfig() *Config {
	return &Config{
		Host:             constants.DefaultAPIV1Host,
		GRPCPort:         constants.DefaultGRPCPort,
		HTTPHost:         constants.DefaultHTTPHost,
		HTTPPort:         "80",
		Timeout:          constants.DefaultTimeout,
		ConnectTimeout:   constants.DefaultConnectTimeout,
		KeepAliveTimeout: constants.DefaultKeepAliveTimeout,
		StreamTimeout:    constants.DefaultStreamTimeout,
		Insecure:         false,
		SkipVerify:       false,
		MaxRetries:       constants.DefaultMaxRetries,
		RetryBackoff:     constants.DefaultRetryBackoff,
		MaxBackoff:       constants.DefaultMaxBackoff,
		Environment:      "production",
		UserAgent:        constants.DefaultUserAgent,
		EnableTelemetry:  true,
	}
}

// NewConfig creates a new Config with default values and applies environment variable overrides.
func NewConfig() *Config {
	config := DefaultConfig()
	config.LoadFromEnvironment()
	return config
}

// NewConfigWithAPIKey creates a new Config with the provided API key.
func NewConfigWithAPIKey(apiKey string) *Config {
	config := NewConfig()
	config.APIKey = apiKey
	return config
}

// LoadFromEnvironment loads configuration from environment variables.
func (c *Config) LoadFromEnvironment() {
	c.loadHostConfig()
	c.loadTimeoutConfig()
	c.loadSecurityConfig()
	c.loadRetryConfig()
	c.loadOtherConfig()
}

func (c *Config) loadHostConfig() {
	// API Key
	if apiKey := os.Getenv("XAI_API_KEY"); apiKey != "" {
		c.APIKey = apiKey
	}

	// Host configuration
	if host := os.Getenv("XAI_HOST"); host != "" {
		c.Host = host
	}

	if grpcPort := os.Getenv("XAI_GRPC_PORT"); grpcPort != "" {
		c.GRPCPort = grpcPort
	}

	if httpHost := os.Getenv("XAI_HTTP_HOST"); httpHost != "" {
		c.HTTPHost = httpHost
	}

	if httpPort := os.Getenv("XAI_HTTP_PORT"); httpPort != "" {
		c.HTTPPort = httpPort
	}
}

func (c *Config) loadTimeoutConfig() {
	// Timeouts
	if timeoutStr := os.Getenv("XAI_TIMEOUT"); timeoutStr != "" {
		if timeout, err := parseDuration(timeoutStr); err == nil {
			c.Timeout = timeout
		}
	}

	if connectTimeoutStr := os.Getenv("XAI_CONNECT_TIMEOUT"); connectTimeoutStr != "" {
		if timeout, err := parseDuration(connectTimeoutStr); err == nil {
			c.ConnectTimeout = timeout
		}
	}

	if keepAliveTimeoutStr := os.Getenv("XAI_KEEPALIVE_TIMEOUT"); keepAliveTimeoutStr != "" {
		if timeout, err := parseDuration(keepAliveTimeoutStr); err == nil {
			c.KeepAliveTimeout = timeout
		}
	}

	if streamTimeoutStr := os.Getenv("XAI_STREAM_TIMEOUT"); streamTimeoutStr != "" {
		if timeout, err := parseDuration(streamTimeoutStr); err == nil {
			c.StreamTimeout = timeout
		}
	}
}

func (c *Config) loadSecurityConfig() {
	// Security settings
	if insecureStr := os.Getenv("XAI_INSECURE"); insecureStr != "" {
		if insecure, err := strconv.ParseBool(insecureStr); err == nil {
			c.Insecure = insecure
		}
	}

	if skipVerifyStr := os.Getenv("XAI_SKIP_VERIFY"); skipVerifyStr != "" {
		if skipVerify, err := strconv.ParseBool(skipVerifyStr); err == nil {
			c.SkipVerify = skipVerify
		}
	}
}

func (c *Config) loadRetryConfig() {
	// Retry settings
	if maxRetriesStr := os.Getenv("XAI_MAX_RETRIES"); maxRetriesStr != "" {
		if maxRetries, err := strconv.Atoi(maxRetriesStr); err == nil && maxRetries >= 0 {
			c.MaxRetries = maxRetries
		}
	}

	if retryBackoffStr := os.Getenv("XAI_RETRY_BACKOFF"); retryBackoffStr != "" {
		if backoff, err := parseDuration(retryBackoffStr); err == nil {
			c.RetryBackoff = backoff
		}
	}

	if maxBackoffStr := os.Getenv("XAI_MAX_BACKOFF"); maxBackoffStr != "" {
		if backoff, err := parseDuration(maxBackoffStr); err == nil {
			c.MaxBackoff = backoff
		}
	}
}

func (c *Config) loadOtherConfig() {
	// Other settings
	if environment := os.Getenv("XAI_ENVIRONMENT"); environment != "" {
		c.Environment = environment
	}

	if userAgent := os.Getenv("XAI_USER_AGENT"); userAgent != "" {
		c.UserAgent = userAgent
	}

	if enableTelemetryStr := os.Getenv("XAI_ENABLE_TELEMETRY"); enableTelemetryStr != "" {
		if enableTelemetry, err := strconv.ParseBool(enableTelemetryStr); err == nil {
			c.EnableTelemetry = enableTelemetry
		}
	}
}

// Validate validates the configuration.
func (c *Config) Validate() error {
	if c.APIKey == "" {
		return errors.NewConfigError("API key is required. Set XAI_API_KEY environment variable or call WithAPIKey().")
	}

	if err := c.validateHost(); err != nil {
		return err
	}

	if err := c.validateTimeouts(); err != nil {
		return err
	}

	if err := c.validateRetries(); err != nil {
		return err
	}

	// Validate environment
	if c.Environment == "" {
		c.Environment = "production"
	}

	// Validate user agent
	if c.UserAgent == "" {
		c.UserAgent = constants.DefaultUserAgent
	}

	return nil
}

func (c *Config) validateHost() error {
	// Validate host configuration
	if c.Host == "" {
		c.Host = constants.DefaultAPIV1Host
	}

	if c.GRPCPort == "" {
		c.GRPCPort = constants.DefaultGRPCPort
	}

	if c.HTTPHost == "" {
		c.HTTPHost = constants.DefaultHTTPHost
	}

	if c.HTTPPort == "" {
		c.HTTPPort = "80"
	}
	return nil
}

func (c *Config) validateTimeouts() error {
	// Validate timeouts
	if c.Timeout <= 0 {
		return errors.NewConfigError("timeout must be positive")
	}

	if c.ConnectTimeout <= 0 {
		return errors.NewConfigError("connect_timeout must be positive")
	}

	if c.KeepAliveTimeout <= 0 {
		return errors.NewConfigError("keep_alive_timeout must be positive")
	}

	if c.StreamTimeout <= 0 {
		return errors.NewConfigError("stream_timeout must be positive")
	}
	return nil
}

func (c *Config) validateRetries() error {
	// Validate retry settings
	if c.MaxRetries < 0 {
		return errors.NewConfigError("max_retries must be non-negative")
	}

	if c.RetryBackoff <= 0 {
		return errors.NewConfigError("retry_backoff must be positive")
	}

	if c.MaxBackoff <= 0 {
		return errors.NewConfigError("max_backoff must be positive")
	}

	if c.MaxBackoff < c.RetryBackoff {
		return errors.NewConfigError("max_backoff must be greater than or equal to retry_backoff")
	}
	return nil
}

// GRPCAddress returns the gRPC address for this configuration.
func (c *Config) GRPCAddress() string {
	return net.JoinHostPort(c.Host, c.GRPCPort)
}

// HTTPAddress returns the HTTP address for this configuration.
func (c *Config) HTTPAddress() string {
	return net.JoinHostPort(c.HTTPHost, c.HTTPPort)
}

// ToSDKMetadata converts the config to SDK metadata.
func (c *Config) ToSDKMetadata() *metadata.SDKMetadata {
	return &metadata.SDKMetadata{
		APIKey:        c.APIKey,
		ClientVersion: c.UserAgent,
		Environment:   c.Environment,
	}
}

// CreateGRPCDialOptions creates gRPC dial options based on the configuration.
func (c *Config) CreateGRPCDialOptions() ([]grpc.DialOption, error) {
	var opts []grpc.DialOption

	// Configure interceptors
	var unaryInterceptors []grpc.UnaryClientInterceptor
	var streamInterceptors []grpc.StreamClientInterceptor

	// Add authentication interceptor
	if c.APIKey != "" {
		authInterceptor := auth.NewAPIKeyAuthInterceptor(c.APIKey, false)
		unaryInterceptors = append(unaryInterceptors, authInterceptor.UnaryInterceptor())
		streamInterceptors = append(streamInterceptors, authInterceptor.StreamInterceptor())
	}

	// Add timeout interceptor
	timeoutInterceptor := grpcutil.NewTimeoutInterceptor(c.Timeout, c.StreamTimeout)
	unaryInterceptors = append(unaryInterceptors, timeoutInterceptor.UnaryInterceptor())
	streamInterceptors = append(streamInterceptors, timeoutInterceptor.StreamInterceptor())

	// Note: Content-Type header is automatically handled by gRPC
	// Adding it manually can cause "malformed header" errors

	// Apply interceptors
	if len(unaryInterceptors) > 0 {
		opts = append(opts, grpc.WithChainUnaryInterceptor(unaryInterceptors...))
	}
	if len(streamInterceptors) > 0 {
		opts = append(opts, grpc.WithChainStreamInterceptor(streamInterceptors...))
	}

	// Configure transport credentials
	if c.Insecure {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	} else {
		tlsConfig := c.CustomTLSConfig
		if tlsConfig == nil {
			tlsConfig = &tls.Config{
				MinVersion:         tls.VersionTLS12,
				InsecureSkipVerify: c.SkipVerify, //nolint:gosec
			}
		}
		opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)))
	}

	// Configure keep-alive
	opts = append(opts, grpc.WithKeepaliveParams(keepalive.ClientParameters{
		Time:    c.KeepAliveTimeout,
		Timeout: c.KeepAliveTimeout,
	}))

	return opts, nil
}

// WithAPIKey sets the API key.
func (c *Config) WithAPIKey(apiKey string) *Config {
	c.APIKey = apiKey
	return c
}

// WithHost sets the API host.
func (c *Config) WithHost(host string) *Config {
	c.Host = host
	c.HTTPHost = host // Keep HTTP host in sync
	return c
}

// WithTimeout sets the request timeout.
func (c *Config) WithTimeout(timeout time.Duration) *Config {
	c.Timeout = timeout
	return c
}

// WithInsecure sets the insecure flag.
func (c *Config) WithInsecure(insecure bool) *Config {
	c.Insecure = insecure
	return c
}

// WithSkipVerify sets the skip verify flag.
func (c *Config) WithSkipVerify(skipVerify bool) *Config {
	c.SkipVerify = skipVerify
	return c
}

// WithTLSConfig sets a custom TLS configuration.
func (c *Config) WithTLSConfig(tlsConfig *tls.Config) *Config {
	c.CustomTLSConfig = tlsConfig
	return c
}

// WithEnvironment sets the environment.
func (c *Config) WithEnvironment(env string) *Config {
	c.Environment = env
	return c
}

// WithUserAgent sets the user agent.
func (c *Config) WithUserAgent(userAgent string) *Config {
	c.UserAgent = userAgent
	return c
}

// WithMaxRetries sets the maximum number of retries.
func (c *Config) WithMaxRetries(maxRetries int) *Config {
	c.MaxRetries = maxRetries
	return c
}

// WithRetryBackoff sets the retry backoff duration.
func (c *Config) WithRetryBackoff(backoff time.Duration) *Config {
	c.RetryBackoff = backoff
	return c
}

// WithMaxBackoff sets the maximum backoff duration.
func (c *Config) WithMaxBackoff(backoff time.Duration) *Config {
	c.MaxBackoff = backoff
	return c
}

// WithEnableTelemetry sets the telemetry enabled flag.
func (c *Config) WithEnableTelemetry(enable bool) *Config {
	c.EnableTelemetry = enable
	return c
}

// String returns a string representation of the config (without sensitive data).
func (c *Config) String() string {
	apiKeyMasked := ""
	if c.APIKey != "" {
		if len(c.APIKey) > 8 {
			apiKeyMasked = strings.Repeat("*", len(c.APIKey)-8) + c.APIKey[len(c.APIKey)-8:]
		} else {
			apiKeyMasked = strings.Repeat("*", len(c.APIKey))
		}
	}

	return fmt.Sprintf("Config{APIKey:%s, Host:%s, Port:%s, Insecure:%t, Environment:%s, MaxRetries:%d, EnableTelemetry:%t}",
		apiKeyMasked, c.Host, c.GRPCPort, c.Insecure, c.Environment, c.MaxRetries, c.EnableTelemetry)
}

// parseDuration parses a duration string that may have a suffix.
// If no suffix is provided, it defaults to seconds.
func parseDuration(s string) (time.Duration, error) {
	// Try parsing with Go's duration parser first
	if d, err := time.ParseDuration(s); err == nil {
		return d, nil
	}

	// If parsing failed, try treating it as seconds
	if num, err := strconv.Atoi(s); err == nil {
		return time.Duration(num) * time.Second, nil
	}

	return 0, fmt.Errorf("invalid duration: %s", s)
}
