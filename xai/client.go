// Package xai provides the main client for interacting with xAI services.
package xai

import (
	"context"
	"fmt"
	"sync"
	"time"

	xaiv1 "github.com/ZaguanLabs/xai-sdk-go/proto/gen/go/xai/v1"
	"github.com/ZaguanLabs/xai-sdk-go/xai/auth"
	"github.com/ZaguanLabs/xai-sdk-go/xai/chat"
	"github.com/ZaguanLabs/xai-sdk-go/xai/collections"
	"github.com/ZaguanLabs/xai-sdk-go/xai/deferred"
	"github.com/ZaguanLabs/xai-sdk-go/xai/documents"
	"github.com/ZaguanLabs/xai-sdk-go/xai/embed"
	"github.com/ZaguanLabs/xai-sdk-go/xai/files"
	"github.com/ZaguanLabs/xai-sdk-go/xai/image"
	"github.com/ZaguanLabs/xai-sdk-go/xai/internal/errors"
	"github.com/ZaguanLabs/xai-sdk-go/xai/internal/metadata"
	"github.com/ZaguanLabs/xai-sdk-go/xai/internal/rest"
	"github.com/ZaguanLabs/xai-sdk-go/xai/models"
	"github.com/ZaguanLabs/xai-sdk-go/xai/sample"
	"github.com/ZaguanLabs/xai-sdk-go/xai/tokenizer"
	"google.golang.org/grpc"
)

// Client represents the main xAI SDK client.
type Client struct {
	config       *Config
	grpcConn     *grpc.ClientConn
	grpcClient   *grpc.ClientConn // Alias for consistency
	restClient   *rest.Client
	chatClient   xaiv1.ChatClient
	modelsClient xaiv1.ModelsClient
	mu           sync.RWMutex
	isClosed     bool
	createdAt    time.Time
	metadata     *metadata.SDKMetadata
}

// NewClient creates a new xAI client with the given configuration.
func NewClient(config *Config) (*Client, error) {
	if config == nil {
		config = DefaultConfig()
	}

	// Validate configuration
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	// Create SDK metadata
	metadata := config.ToSDKMetadata()

	client := &Client{
		config:    config,
		metadata:  metadata,
		createdAt: time.Now(),
	}

	// Create REST client
	baseURL := fmt.Sprintf("https://%s/v1", config.HTTPHost)
	if config.Insecure {
		baseURL = fmt.Sprintf("http://%s/v1", config.HTTPHost)
	}
	client.restClient = rest.NewClient(rest.Config{
		BaseURL:   baseURL,
		APIKey:    config.APIKey,
		UserAgent: metadata.UserAgent,
		Timeout:   config.Timeout,
	})

	// Create gRPC connection
	if err := client.createGRPCConnection(); err != nil {
		return nil, fmt.Errorf("failed to create gRPC connection: %w", err)
	}

	// Create chat client
	client.chatClient = xaiv1.NewChatClient(client.grpcConn)

	return client, nil
}

// NewClientWithAPIKey creates a new client with the given API key and default configuration.
func NewClientWithAPIKey(apiKey string) (*Client, error) {
	config := NewConfigWithAPIKey(apiKey)
	return NewClient(config)
}

// NewClientFromEnvironment creates a new client loading configuration from environment variables.
func NewClientFromEnvironment() (*Client, error) {
	config := NewConfig()
	return NewClient(config)
}

// createGRPCConnection creates and configures the gRPC connection.
func (c *Client) createGRPCConnection() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.isClosed {
		return errors.NewConfigError("client is closed")
	}

	// Create gRPC dial options
	dialOptions, err := c.config.CreateGRPCDialOptions()
	if err != nil {
		return fmt.Errorf("failed to create gRPC dial options: %w", err)
	}

	// Create gRPC connection using NewClient (replaces deprecated DialContext)
	// Note: NewClient doesn't establish connection immediately - it's lazy
	// Connection happens on first RPC call
	grpcConn, err := grpc.NewClient(
		c.config.GRPCAddress(),
		dialOptions...,
	)

	if err != nil {
		return fmt.Errorf("failed to dial gRPC server at %s: %w", c.config.GRPCAddress(), err)
	}

	c.grpcConn = grpcConn
	c.grpcClient = grpcConn // Alias for compatibility
	c.chatClient = xaiv1.NewChatClient(grpcConn)
	c.modelsClient = xaiv1.NewModelsClient(grpcConn)

	return nil
}

// Config returns the client's configuration.
func (c *Client) Config() *Config {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.config
}

// Metadata returns the client's SDK metadata.
func (c *Client) Metadata() *metadata.SDKMetadata {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.metadata
}

// CreatedAt returns the time when the client was created.
func (c *Client) CreatedAt() time.Time {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.createdAt
}

// IsClosed returns whether the client has been closed.
func (c *Client) IsClosed() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.isClosed
}

// GRPCConnection returns the underlying gRPC connection.
// This should be used sparingly and callers should not close this connection.
func (c *Client) GRPCConnection() *grpc.ClientConn {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.grpcConn
}

// Chat returns the chat service client.
func (c *Client) Chat() xaiv1.ChatClient {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.chatClient
}

// NewContext creates a new context with the client's configuration and metadata.
func (c *Client) NewContext(ctx context.Context) context.Context {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if ctx == nil {
		ctx = context.Background()
	}

	// Add timeout if not already set
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, c.config.Timeout)
		// Note: We intentionally don't call cancel here as the context is returned to the caller
		// The caller is responsible for managing the context lifecycle
		_ = cancel
	}

	// Add metadata to context
	return c.metadata.AddToOutgoingContext(ctx)
}

// NewContextWithTimeout creates a new context with a specific timeout.
func (c *Client) NewContextWithTimeout(ctx context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	if ctx == nil {
		ctx = context.Background()
	}

	return context.WithTimeout(c.metadata.AddToOutgoingContext(ctx), timeout)
}

// NewContextWithCancel creates a new context with a cancel function.
func (c *Client) NewContextWithCancel(ctx context.Context) (context.Context, context.CancelFunc) {
	if ctx == nil {
		ctx = context.Background()
	}

	return context.WithCancel(c.metadata.AddToOutgoingContext(ctx))
}

// NewContextWithDeadline creates a new context with a deadline.
func (c *Client) NewContextWithDeadline(ctx context.Context, deadline time.Time) (context.Context, context.CancelFunc) {
	if ctx == nil {
		_ = context.Background() //nolint:ineffassign // Intentional reassignment for nil context
	}

	return context.WithDeadline(c.metadata.AddToOutgoingContext(ctx), deadline)
}

// Close closes the client and its underlying connections.
func (c *Client) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.isClosed {
		return nil // Already closed
	}

	c.isClosed = true

	// Close REST client
	if c.restClient != nil {
		c.restClient.Close()
	}

	// Close gRPC connection
	if c.grpcConn != nil {
		if err := c.grpcConn.Close(); err != nil {
			return fmt.Errorf("failed to close gRPC connection: %w", err)
		}
		c.grpcConn = nil
		c.grpcClient = nil
	}

	return nil
}

// CloseWithContext closes the client with a context for cancellation.
func (c *Client) CloseWithContext(ctx context.Context) error {
	// Try to close with context first
	done := make(chan struct{})
	var closeErr error

	go func() {
		defer close(done)
		closeErr = c.Close()
	}()

	select {
	case <-done:
		return closeErr
	case <-ctx.Done():
		return ctx.Err()
	}
}

// EnsureGRPCConnection ensures that the gRPC connection is healthy and reconnect if needed.
func (c *Client) EnsureGRPCConnection() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.isClosed {
		return errors.NewConfigError("client is closed")
	}

	// Check if connection is healthy
	if c.grpcConn != nil && c.grpcConn.GetState().String() != "TRANSIENT_FAILURE" {
		return nil
	}

	// Recreate connection if needed
	if err := c.createGRPCConnection(); err != nil {
		return fmt.Errorf("failed to recreate gRPC connection: %w", err)
	}

	return nil
}

// HealthCheck performs a health check on the client connection.
func (c *Client) HealthCheck(_ context.Context) error {
	// Ensure connection is healthy
	if err := c.EnsureGRPCConnection(); err != nil {
		return fmt.Errorf("connection health check failed: %w", err)
	}

	// Basic connectivity check
	state := c.grpcConn.GetState()
	if state.String() == "TRANSIENT_FAILURE" {
		return errors.NewNetworkError("gRPC connection is in transient failure state")
	}

	return nil
}

// WithTimeout creates a new client with a different timeout setting.
func (c *Client) WithTimeout(timeout time.Duration) *Client {
	c.mu.Lock()
	defer c.mu.Unlock()

	newConfig := *c.config
	newConfig.Timeout = timeout

	return &Client{
		config:       &newConfig,
		metadata:     c.metadata,
		grpcConn:     c.grpcConn,
		grpcClient:   c.grpcConn,
		restClient:   c.restClient,
		chatClient:   c.chatClient,
		modelsClient: c.modelsClient,
		createdAt:    c.createdAt,
		isClosed:     c.isClosed,
	}
}

// WithAPIKey creates a new client with a different API key.
func (c *Client) WithAPIKey(apiKey string) *Client {
	c.mu.Lock()
	defer c.mu.Unlock()

	newConfig := *c.config
	newConfig.APIKey = apiKey

	newMetadata := c.metadata
	newMetadata.APIKey = apiKey

	return &Client{
		config:       &newConfig,
		metadata:     newMetadata,
		grpcConn:     c.grpcConn,
		grpcClient:   c.grpcConn,
		restClient:   c.restClient,
		chatClient:   c.chatClient,
		modelsClient: c.modelsClient,
		createdAt:    c.createdAt,
		isClosed:     c.isClosed,
	}
}

// String returns a string representation of the client.
func (c *Client) String() string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return fmt.Sprintf("Client{Config:%s, CreatedAt:%s, Closed:%t}",
		c.config.String(),
		c.createdAt.Format(time.RFC3339),
		c.isClosed)
}

// HealthStatus represents the health status of the client.
type HealthStatus struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Message   string    `json:"message,omitempty"`
}

// GetHealthStatus returns the current health status of the client.
func (c *Client) GetHealthStatus() HealthStatus {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.isClosed {
		return HealthStatus{
			Status:    "closed",
			Timestamp: time.Now(),
			Message:   "Client is closed",
		}
	}

	if c.grpcConn == nil {
		return HealthStatus{
			Status:    "disconnected",
			Timestamp: time.Now(),
			Message:   "No gRPC connection available",
		}
	}

	state := c.grpcConn.GetState()

	switch state.String() {
	case "READY":
		return HealthStatus{
			Status:    "healthy",
			Timestamp: time.Now(),
			Message:   "Connection is ready",
		}
	case "CONNECTING":
		return HealthStatus{
			Status:    "connecting",
			Timestamp: time.Now(),
			Message:   "Connection is being established",
		}
	case "TRANSIENT_FAILURE":
		return HealthStatus{
			Status:    "unhealthy",
			Timestamp: time.Now(),
			Message:   "Connection is in transient failure state",
		}
	case "IDLE":
		return HealthStatus{
			Status:    "idle",
			Timestamp: time.Now(),
			Message:   "Connection is idle",
		}
	case "SHUTDOWN":
		return HealthStatus{
			Status:    "shutdown",
			Timestamp: time.Now(),
			Message:   "Connection is shutdown",
		}
	default:
		return HealthStatus{
			Status:    "unknown",
			Timestamp: time.Now(),
			Message:   fmt.Sprintf("Unknown connection state: %s", state.String()),
		}
	}
}

// NewChatRequest creates a new chat request with the specified model.
func (c *Client) NewChatRequest(model string, opts ...chat.RequestOption) *chat.Request {
	return chat.NewRequest(model, opts...)
}

// Models returns the models service client.
func (c *Client) Models() *models.Client {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return models.NewClient(c.modelsClient)
}

// Embed returns the embeddings service client.
func (c *Client) Embed() *embed.Client {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return embed.NewClient(c.restClient)
}

// Files returns the files service client.
func (c *Client) Files() *files.Client {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return files.NewClient(c.restClient)
}

// Collections returns the collections service client.
func (c *Client) Collections() *collections.Client {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return collections.NewClient(c.restClient)
}

// Auth returns the auth service client.
func (c *Client) Auth() *auth.Client {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return auth.NewClient(c.restClient)
}

// Images returns the image generation service client.
func (c *Client) Images() *image.Client {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return image.NewClient(c.restClient)
}

// Deferred returns the deferred completions service client.
func (c *Client) Deferred() *deferred.Client {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return deferred.NewClient(c.restClient)
}

// Documents returns the document search service client.
func (c *Client) Documents() *documents.Client {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return documents.NewClient(c.restClient)
}

// Sample returns the sample/completion service client (legacy).
func (c *Client) Sample() *sample.Client {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return sample.NewClient(c.restClient)
}

// Tokenizer returns the tokenization service client.
func (c *Client) Tokenizer() *tokenizer.Client {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return tokenizer.NewClient(c.restClient)
}
