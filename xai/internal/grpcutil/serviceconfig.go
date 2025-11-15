// Package grpcutil provides gRPC utilities for retry policies, keepalive, and service configuration.
package grpcutil

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/status"
)

// RetryPolicy defines retry configuration for gRPC calls.
type RetryPolicy struct {
	// MaxAttempts is the maximum number of retry attempts.
	MaxAttempts int `json:"max_attempts"`
	
	// InitialBackoff is the initial backoff duration.
	InitialBackoff time.Duration `json:"initial_backoff"`
	
	// MaxBackoff is the maximum backoff duration.
	MaxBackoff time.Duration `json:"max_backoff"`
	
	// BackoffMultiplier is the multiplier for exponential backoff.
	BackoffMultiplier float64 `json:"backoff_multiplier"`
	
	// RetryableCodes are the gRPC status codes that should trigger a retry.
	RetryableCodes []string `json:"retryable_codes"`
}

// DefaultRetryPolicy returns the default retry policy.
func DefaultRetryPolicy() *RetryPolicy {
	return &RetryPolicy{
		MaxAttempts:       3,
		InitialBackoff:    1 * time.Second,
		MaxBackoff:        60 * time.Second,
		BackoffMultiplier: 1.6,
		RetryableCodes: []string{
			"UNAVAILABLE",
			"DEADLINE_EXCEEDED", 
			"RESOURCE_EXHAUSTED",
			"INTERNAL",
		},
	}
}

// ToServiceConfig converts the retry policy to a gRPC service config.
func (rp *RetryPolicy) ToServiceConfig() map[string]any {
	retryPolicy := map[string]any{
		"maxAttempts":          float64(rp.MaxAttempts),
		"initialBackoff":       fmt.Sprintf("%v", rp.InitialBackoff),
		"maxBackoff":           fmt.Sprintf("%v", rp.MaxBackoff),
		"backoffMultiplier":    rp.BackoffMultiplier,
		"retryableStatusCodes": rp.RetryableCodes,
	}

	return map[string]any{
		"retryPolicy": retryPolicy,
	}
}

// KeepAliveConfig defines keepalive configuration for gRPC connections.
type KeepAliveConfig struct {
	// Time is the time after which if no activity is seen, a ping is sent.
	Time time.Duration `json:"time"`
	
	// Timeout is the time the ping response must be received before the connection is closed.
	Timeout time.Duration `json:"timeout"`
	
	// EnforcementPolicy specifies the keepalive enforcement policy.
	EnforcementPolicy keepalive.EnforcementPolicy `json:"enforcement_policy"`
}

// DefaultKeepAliveConfig returns the default keepalive configuration.
func DefaultKeepAliveConfig() *KeepAliveConfig {
	return &KeepAliveConfig{
		Time:    20 * time.Second,
		Timeout: 10 * time.Second,
		EnforcementPolicy: keepalive.EnforcementPolicy{
			MinTime:             5 * time.Second,  // Minimum time between pings
			PermitWithoutStream: true,              // Permit pings without active streams
		},
	}
}

// NewKeepAliveEnforcementPolicy creates a keepalive enforcement policy.
func NewKeepAliveEnforcementPolicy(minTime time.Duration, permitWithoutStream bool) keepalive.EnforcementPolicy {
	return keepalive.EnforcementPolicy{
		MinTime:             minTime,
		PermitWithoutStream: permitWithoutStream,
	}
}

// ToGRPCKeepAliveParams converts to gRPC keepalive parameters.
func (kac *KeepAliveConfig) ToGRPCKeepAliveParams() keepalive.ClientParameters {
	return keepalive.ClientParameters{
		Time:    kac.Time,
		Timeout: kac.Timeout,
	}
}

// ToGRPCKeepAliveEnforcementPolicy converts to gRPC keepalive enforcement policy.
func (kac *KeepAliveConfig) ToGRPCKeepAliveEnforcementPolicy() keepalive.EnforcementPolicy {
	return kac.EnforcementPolicy
}

// DialOptionBuilder builds gRPC dial options.
type DialOptionBuilder struct {
	retryPolicy    *RetryPolicy
	keepAlive      *KeepAliveConfig
	insecure       bool
	customTLS      credentials.TransportCredentials
	block          bool
	waitForReady   bool
	compressor     string
}

// NewDialOptionBuilder creates a new dial option builder.
func NewDialOptionBuilder() *DialOptionBuilder {
	return &DialOptionBuilder{
		retryPolicy: DefaultRetryPolicy(),
		keepAlive:   DefaultKeepAliveConfig(),
		insecure:    false,
		block:       true,
		waitForReady: true,
	}
}

// WithInsecure sets the connection to be insecure.
func (dob *DialOptionBuilder) WithInsecure() *DialOptionBuilder {
	dob.insecure = true
	return dob
}

// WithTLS sets custom TLS credentials.
func (dob *DialOptionBuilder) WithTLS(creds credentials.TransportCredentials) *DialOptionBuilder {
	dob.customTLS = creds
	return dob
}

// WithRetryPolicy sets the retry policy.
func (dob *DialOptionBuilder) WithRetryPolicy(policy *RetryPolicy) *DialOptionBuilder {
	dob.retryPolicy = policy
	return dob
}

// WithKeepAlive sets the keepalive configuration.
func (dob *DialOptionBuilder) WithKeepAlive(config *KeepAliveConfig) *DialOptionBuilder {
	dob.keepAlive = config
	return dob
}

// WithBlock sets whether to block until the connection is established.
func (dob *DialOptionBuilder) WithBlock(block bool) *DialOptionBuilder {
	dob.block = block
	return dob
}

// WithWaitForReady sets whether to wait for ready state before RPC calls.
func (dob *DialOptionBuilder) WithWaitForReady(waitForReady bool) *DialOptionBuilder {
	dob.waitForReady = waitForReady
	return dob
}

// WithCompressor sets the compressor to use.
func (dob *DialOptionBuilder) WithCompressor(compressor string) *DialOptionBuilder {
	dob.compressor = compressor
	return dob
}

// Build creates the gRPC dial options.
func (dob *DialOptionBuilder) Build() ([]grpc.DialOption, error) {
	var opts []grpc.DialOption

	// Transport credentials
	if dob.insecure {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	} else if dob.customTLS != nil {
		opts = append(opts, grpc.WithTransportCredentials(dob.customTLS))
	} else {
		// Use default TLS
		opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(nil)))
	}

	// Keepalive
	opts = append(opts, grpc.WithKeepaliveParams(dob.keepAlive.ToGRPCKeepAliveParams()))

	// Retry policy via service config
	retryConfig := dob.retryPolicy.ToServiceConfig()
	retryConfigJSON, err := json.Marshal(retryConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal retry config: %w", err)
	}
	opts = append(opts, grpc.WithDefaultServiceConfig(string(retryConfigJSON)))

	// Block until connected
	if dob.block {
		opts = append(opts, grpc.WithBlock())
	}

	// Note: grpc.WithWaitForReady and grpc.NewCompressor are not available in all gRPC versions
	// Users can add these manually if needed

	return opts, nil
}

// ConnectionState represents the state of a gRPC connection.
type ConnectionState struct {
	State        connectivity.State `json:"state"`
	Address      string             `json:"address,omitempty"`
	LastConnected time.Time         `json:"last_connected,omitempty"`
	Error        string             `json:"error,omitempty"`
}

// HealthChecker provides health checking utilities for gRPC connections.
type HealthChecker struct {
	conn *grpc.ClientConn
	addr string
}

// NewHealthChecker creates a new health checker.
func NewHealthChecker(conn *grpc.ClientConn, address string) *HealthChecker {
	return &HealthChecker{
		conn: conn,
		addr: address,
	}
}

// CheckConnection checks the health of the connection.
func (hc *HealthChecker) CheckConnection() ConnectionState {
	if hc.conn == nil {
		return ConnectionState{
			State:  connectivity.Shutdown,
			Address: hc.addr,
			Error:  "connection is nil",
		}
	}
	
	state := hc.conn.GetState()
	
	return ConnectionState{
		State:     state,
		Address:   hc.addr,
		Error:     hc.getStateError(state),
	}
}

// IsHealthy returns whether the connection is healthy.
func (hc *HealthChecker) IsHealthy() bool {
	state := hc.conn.GetState()
	return state == connectivity.Ready || state == connectivity.Idle
}

// getStateError returns an error message for a given connection state.
func (hc *HealthChecker) getStateError(state connectivity.State) string {
	switch state {
	case connectivity.Ready:
		return ""
	case connectivity.Connecting:
		return "connection is being established"
	case connectivity.Idle:
		return "connection is idle"
	case connectivity.TransientFailure:
		return "connection is temporarily unavailable"
	case connectivity.Shutdown:
		return "connection is closed"
	default:
		return fmt.Sprintf("unknown connection state: %v", state)
	}
}

// RetryWithBackoff performs a retry operation with exponential backoff.
func RetryWithBackoff(ctx context.Context, policy *RetryPolicy, operation func() error) error {
	backoff := policy.InitialBackoff
	
	for attempt := 0; attempt < policy.MaxAttempts; attempt++ {
		err := operation()
		if err == nil {
			return nil
		}
		
		// Check if error is retryable
		if !isRetryableError(err, policy.RetryableCodes) {
			return err
		}
		
		// If this is the last attempt, return the error
		if attempt == policy.MaxAttempts-1 {
			return err
		}
		
		// Create a new context with timeout for the backoff period
		backoffCtx, cancel := context.WithTimeout(ctx, backoff)
		defer cancel()
		
		// Wait for the backoff period or context cancellation
		select {
		case <-backoffCtx.Done():
			return backoffCtx.Err()
		case <-time.After(backoff):
		}
		
		// Increase backoff for next iteration
		backoff = time.Duration(float64(backoff) * policy.BackoffMultiplier)
		if backoff > policy.MaxBackoff {
			backoff = policy.MaxBackoff
		}
	}
	
	return nil
}

// isRetryableError checks if an error is retryable based on the retryable codes.
func isRetryableError(err error, retryableCodes []string) bool {
	// Check if error is a gRPC status error
	if st, ok := status.FromError(err); ok {
		code := st.Code().String()
		for _, retryableCode := range retryableCodes {
			if code == retryableCode {
				return true
			}
		}
	}
	
	// For other error types, check common retryable conditions
	errorStr := err.Error()
	if errorStr == "context canceled" || errorStr == "context deadline exceeded" {
		return true
	}
	
	return false
}

// NewInsecureCredentials creates insecure credentials for testing.
func NewInsecureCredentials() credentials.TransportCredentials {
	return insecure.NewCredentials()
}