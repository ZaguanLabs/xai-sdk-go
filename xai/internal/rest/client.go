// Package rest provides a REST client for xAI APIs.
package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"sync"
	"time"
)

// Buffer pool for JSON encoding to reduce allocations
var bufferPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

const (
	// MaxResponseSize is the maximum size of a response body (100MB)
	MaxResponseSize = 100 * 1024 * 1024
)

// Client is a REST client for xAI APIs.
type Client struct {
	httpClient *http.Client
	baseURL    string
	apiKey     string
	userAgent  string
}

// Config contains configuration for the REST client.
type Config struct {
	BaseURL   string
	APIKey    string
	UserAgent string
	Timeout   time.Duration
}

// NewClient creates a new REST client with optimized connection pooling.
func NewClient(cfg Config) *Client {
	if cfg.BaseURL == "" {
		cfg.BaseURL = "https://api.x.ai/v1"
	}
	if cfg.Timeout == 0 {
		cfg.Timeout = 30 * time.Second
	}
	if cfg.UserAgent == "" {
		cfg.UserAgent = "xai-sdk-go"
	}

	// Create optimized HTTP transport with connection pooling
	transport := &http.Transport{
		// Connection pooling settings
		MaxIdleConns:        100,              // Maximum idle connections across all hosts
		MaxIdleConnsPerHost: 10,               // Maximum idle connections per host
		MaxConnsPerHost:     100,              // Maximum total connections per host
		IdleConnTimeout:     90 * time.Second, // How long idle connections stay open

		// Timeout settings
		DialContext: (&net.Dialer{
			Timeout:   10 * time.Second, // Connection timeout
			KeepAlive: 30 * time.Second, // TCP keepalive
		}).DialContext,
		TLSHandshakeTimeout:   10 * time.Second,
		ResponseHeaderTimeout: 10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,

		// Performance optimizations
		DisableCompression: false, // Enable gzip compression
		ForceAttemptHTTP2:  true,  // Use HTTP/2 when available
	}

	return &Client{
		httpClient: &http.Client{
			Timeout:   cfg.Timeout,
			Transport: transport,
		},
		baseURL:   cfg.BaseURL,
		apiKey:    cfg.APIKey,
		userAgent: cfg.UserAgent,
	}
}

// Request represents an HTTP request.
type Request struct {
	Method  string
	Path    string
	Body    interface{}
	Headers map[string]string
}

// Response represents an HTTP response.
type Response struct {
	StatusCode int
	Body       []byte
	Headers    http.Header
}

// Do executes an HTTP request.
func (c *Client) Do(ctx context.Context, req Request) (*Response, error) {
	url := c.baseURL + req.Path

	var body io.Reader
	if req.Body != nil {
		// Use buffer pool to reduce allocations
		buf := bufferPool.Get().(*bytes.Buffer)
		buf.Reset()
		defer bufferPool.Put(buf)

		if err := json.NewEncoder(buf).Encode(req.Body); err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}

		// Create a copy since buf will be returned to pool
		bodyBytes := make([]byte, buf.Len())
		copy(bodyBytes, buf.Bytes())
		body = bytes.NewReader(bodyBytes)
	}

	httpReq, err := http.NewRequestWithContext(ctx, req.Method, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set default headers
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("User-Agent", c.userAgent)
	if c.apiKey != "" {
		httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)
	}

	// Set custom headers
	for key, value := range req.Headers {
		httpReq.Header.Set(key, value)
	}

	httpResp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer httpResp.Body.Close()

	// Limit response size to prevent memory exhaustion
	limitedReader := io.LimitReader(httpResp.Body, MaxResponseSize)
	respBody, err := io.ReadAll(limitedReader)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	resp := &Response{
		StatusCode: httpResp.StatusCode,
		Body:       respBody,
		Headers:    httpResp.Header,
	}

	// Check for HTTP errors
	if httpResp.StatusCode >= 400 {
		return resp, &HTTPError{
			StatusCode: httpResp.StatusCode,
			Body:       respBody,
		}
	}

	return resp, nil
}

// Get executes a GET request.
func (c *Client) Get(ctx context.Context, path string) (*Response, error) {
	return c.Do(ctx, Request{
		Method: http.MethodGet,
		Path:   path,
	})
}

// Post executes a POST request.
func (c *Client) Post(ctx context.Context, path string, body interface{}) (*Response, error) {
	return c.Do(ctx, Request{
		Method: http.MethodPost,
		Path:   path,
		Body:   body,
	})
}

// Put executes a PUT request.
func (c *Client) Put(ctx context.Context, path string, body interface{}) (*Response, error) {
	return c.Do(ctx, Request{
		Method: http.MethodPut,
		Path:   path,
		Body:   body,
	})
}

// Delete executes a DELETE request.
func (c *Client) Delete(ctx context.Context, path string) (*Response, error) {
	return c.Do(ctx, Request{
		Method: http.MethodDelete,
		Path:   path,
	})
}

// Close closes idle connections in the HTTP client's connection pool.
// This should be called when the client is no longer needed to free resources.
func (c *Client) Close() {
	if c.httpClient != nil && c.httpClient.Transport != nil {
		if transport, ok := c.httpClient.Transport.(*http.Transport); ok {
			transport.CloseIdleConnections()
		}
	}
}

// DecodeJSON decodes a JSON response into a target struct.
func (r *Response) DecodeJSON(target interface{}) error {
	if err := json.Unmarshal(r.Body, target); err != nil {
		return fmt.Errorf("failed to decode JSON response: %w", err)
	}
	return nil
}
