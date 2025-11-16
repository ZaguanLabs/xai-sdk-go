// Package grpcutil provides gRPC utility functions and interceptors.
package grpcutil

import (
	"context"
	"time"

	"google.golang.org/grpc"
)

// TimeoutInterceptor creates a unary interceptor that enforces request timeouts.
type TimeoutInterceptor struct {
	defaultTimeout time.Duration
	maxTimeout     time.Duration
}

// NewTimeoutInterceptor creates a new timeout interceptor.
func NewTimeoutInterceptor(defaultTimeout, maxTimeout time.Duration) *TimeoutInterceptor {
	if defaultTimeout <= 0 {
		defaultTimeout = 30 * time.Second
	}
	if maxTimeout <= 0 {
		maxTimeout = 5 * time.Minute
	}

	return &TimeoutInterceptor{
		defaultTimeout: defaultTimeout,
		maxTimeout:     maxTimeout,
	}
}

// UnaryInterceptor returns a unary interceptor function for timeout enforcement.
func (t *TimeoutInterceptor) UnaryInterceptor() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		// Check if context already has a deadline
		if _, hasDeadline := ctx.Deadline(); hasDeadline {
			// Context already has a deadline, use it as is
			return invoker(ctx, method, req, reply, cc, opts...)
		}

		// Apply default timeout
		timeoutCtx, cancel := context.WithTimeout(ctx, t.defaultTimeout)
		defer cancel()

		// Call the invoker with timeout context
		return invoker(timeoutCtx, method, req, reply, cc, opts...)
	}
}

// StreamInterceptor returns a stream interceptor function for timeout enforcement.
func (t *TimeoutInterceptor) StreamInterceptor() grpc.StreamClientInterceptor {
	return func(
		ctx context.Context,
		desc *grpc.StreamDesc,
		cc *grpc.ClientConn,
		method string,
		streamer grpc.Streamer,
		opts ...grpc.CallOption,
	) (grpc.ClientStream, error) {
		// For streaming, we need a longer timeout or no timeout
		streamTimeout := t.defaultTimeout
		if desc.ServerStreams {
			// For server streaming, use a much longer timeout or the configured stream timeout
			streamTimeout = 5 * time.Minute // Default for server streaming
		}

		// Check if context already has a deadline
		if _, hasDeadline := ctx.Deadline(); hasDeadline {
			// Context already has a deadline, use it as is
			return streamer(ctx, desc, cc, method, opts...)
		}

		// Apply stream timeout
		timeoutCtx, cancel := context.WithTimeout(ctx, streamTimeout)
		_ = cancel // Stream cancel cleanup is handled by client

		// Call the streamer with timeout context
		return streamer(timeoutCtx, desc, cc, method, opts...)
	}
}

// TimeoutUnaryInterceptor creates a simple unary timeout interceptor with the given duration.
func TimeoutUnaryInterceptor(timeout time.Duration) grpc.UnaryClientInterceptor {
	interceptor := NewTimeoutInterceptor(timeout, 0)
	return interceptor.UnaryInterceptor()
}

// TimeoutStreamInterceptor creates a simple stream timeout interceptor with the given duration.
func TimeoutStreamInterceptor(timeout time.Duration) grpc.StreamClientInterceptor {
	interceptor := NewTimeoutInterceptor(timeout, 0)
	return interceptor.StreamInterceptor()
}

// StreamTimeoutSettings configures timeout settings for streaming operations.
type StreamTimeoutSettings struct {
	ServerStreamTimeout  time.Duration
	ClientStreamTimeout  time.Duration
	BidirectionalTimeout time.Duration
}

// NewStreamTimeoutInterceptor creates an interceptor with custom stream timeout settings.
func NewStreamTimeoutInterceptor(settings StreamTimeoutSettings) *TimeoutInterceptor {
	return &TimeoutInterceptor{
		defaultTimeout: settings.ServerStreamTimeout,
		maxTimeout:     settings.BidirectionalTimeout,
	}
}

// TimeoutOptions provides options for configuring timeout behavior.
type TimeoutOptions struct {
	DefaultTimeout      time.Duration
	MaxTimeout          time.Duration
	EnableDeadlineCheck bool
}

// NewTimeoutInterceptorWithOptions creates a timeout interceptor with custom options.
func NewTimeoutInterceptorWithOptions(options TimeoutOptions) *TimeoutInterceptor {
	if options.DefaultTimeout <= 0 {
		options.DefaultTimeout = 30 * time.Second
	}
	if options.MaxTimeout <= 0 {
		options.MaxTimeout = 5 * time.Minute
	}

	return &TimeoutInterceptor{
		defaultTimeout: options.DefaultTimeout,
		maxTimeout:     options.MaxTimeout,
	}
}

// DeadlineInterceptor creates an interceptor that enforces request deadlines.
type DeadlineInterceptor struct {
	maxTimeout time.Duration
}

// NewDeadlineInterceptor creates a new deadline interceptor.
func NewDeadlineInterceptor(maxTimeout time.Duration) *DeadlineInterceptor {
	if maxTimeout <= 0 {
		maxTimeout = 5 * time.Minute
	}

	return &DeadlineInterceptor{
		maxTimeout: maxTimeout,
	}
}

// UnaryInterceptor returns a unary interceptor that enforces deadlines.
func (d *DeadlineInterceptor) UnaryInterceptor() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		// Check for existing deadline
		if deadline, ok := ctx.Deadline(); ok {
			// Check if deadline is too far in the future
			maxDeadline := time.Now().Add(d.maxTimeout)
			if deadline.After(maxDeadline) {
				// Cap the deadline to maxTimeout
				cappedCtx, cancel := context.WithDeadline(ctx, maxDeadline)
				defer cancel()
				ctx = cappedCtx
			}
		}

		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

// StreamInterceptor returns a stream interceptor that enforces deadlines.
func (d *DeadlineInterceptor) StreamInterceptor() grpc.StreamClientInterceptor {
	return func(
		ctx context.Context,
		desc *grpc.StreamDesc,
		cc *grpc.ClientConn,
		method string,
		streamer grpc.Streamer,
		opts ...grpc.CallOption,
	) (grpc.ClientStream, error) {
		// Check for existing deadline
		if deadline, ok := ctx.Deadline(); ok {
			// Check if deadline is too far in the future
			maxDeadline := time.Now().Add(d.maxTimeout)
			if deadline.After(maxDeadline) {
				// Cap the deadline to maxTimeout
				cappedCtx, cancel := context.WithDeadline(ctx, maxDeadline)
				_ = cancel // Stream cancel cleanup is handled by client
				ctx = cappedCtx
			}
		}

		return streamer(ctx, desc, cc, method, opts...)
	}
}

// GetRemainingTimeout returns the remaining timeout in the context.
func GetRemainingTimeout(ctx context.Context) (time.Duration, bool) {
	deadline, ok := ctx.Deadline()
	if !ok {
		return 0, false
	}

	return time.Until(deadline), true
}
