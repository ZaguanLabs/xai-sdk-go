// Package auth provides authentication interceptors for gRPC requests.
package auth

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	// AuthorizationHeader is the authorization header key.
	AuthorizationHeader = "authorization"
	// XAPIKeyHeader is the x-api-key header key.
	XAPIKeyHeader = "x-api-key"
)

// APIKeyAuthInterceptor creates a unary interceptor that authenticates requests with an API key.
type APIKeyAuthInterceptor struct {
	apiKey     string
	useXAPIKey bool
}

// NewAPIKeyAuthInterceptor creates a new API key authentication interceptor.
func NewAPIKeyAuthInterceptor(apiKey string, useXAPIKey bool) *APIKeyAuthInterceptor {
	return &APIKeyAuthInterceptor{
		apiKey:     apiKey,
		useXAPIKey: useXAPIKey,
	}
}

// UnaryInterceptor returns a unary interceptor function for authentication.
func (a *APIKeyAuthInterceptor) UnaryInterceptor() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		// Add authentication metadata to context
		authCtx, err := a.addAuthMetadata(ctx)
		if err != nil {
			return status.Error(codes.Unauthenticated, fmt.Sprintf("failed to add auth metadata: %v", err))
		}

		// Call the invoker with authenticated context
		return invoker(authCtx, method, req, reply, cc, opts...)
	}
}

// StreamInterceptor returns a stream interceptor function for authentication.
func (a *APIKeyAuthInterceptor) StreamInterceptor() grpc.StreamClientInterceptor {
	return func(
		ctx context.Context,
		desc *grpc.StreamDesc,
		cc *grpc.ClientConn,
		method string,
		streamer grpc.Streamer,
		opts ...grpc.CallOption,
	) (grpc.ClientStream, error) {
		// Add authentication metadata to context
		authCtx, err := a.addAuthMetadata(ctx)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, fmt.Sprintf("failed to add auth metadata: %v", err))
		}

		// Call the streamer with authenticated context
		return streamer(authCtx, desc, cc, method, opts...)
	}
}

// addAuthMetadata adds authentication metadata to the context.
// Uses AppendToOutgoingContext to preserve gRPC's internal metadata.
func (a *APIKeyAuthInterceptor) addAuthMetadata(ctx context.Context) (context.Context, error) {
	if a.apiKey == "" {
		return nil, fmt.Errorf("API key is empty")
	}

	// Add authentication metadata using AppendToOutgoingContext
	// This preserves any existing metadata including gRPC's content-type
	if a.useXAPIKey {
		// Use x-api-key header
		return metadata.AppendToOutgoingContext(ctx, XAPIKeyHeader, a.apiKey), nil
	}
	// Use authorization bearer token
	return metadata.AppendToOutgoingContext(ctx, AuthorizationHeader, fmt.Sprintf("Bearer %s", a.apiKey)), nil
}

// APIKeyAuthUnaryInterceptor creates a simple unary interceptor with the given API key.
func APIKeyAuthUnaryInterceptor(apiKey string) grpc.UnaryClientInterceptor {
	interceptor := NewAPIKeyAuthInterceptor(apiKey, false)
	return interceptor.UnaryInterceptor()
}

// APIKeyAuthStreamInterceptor creates a simple stream interceptor with the given API key.
func APIKeyAuthStreamInterceptor(apiKey string) grpc.StreamClientInterceptor {
	interceptor := NewAPIKeyAuthInterceptor(apiKey, false)
	return interceptor.StreamInterceptor()
}

// XAPIKeyAuthUnaryInterceptor creates a unary interceptor using x-api-key header.
func XAPIKeyAuthUnaryInterceptor(apiKey string) grpc.UnaryClientInterceptor {
	interceptor := NewAPIKeyAuthInterceptor(apiKey, true)
	return interceptor.UnaryInterceptor()
}

// XAPIKeyAuthStreamInterceptor creates a stream interceptor using x-api-key header.
func XAPIKeyAuthStreamInterceptor(apiKey string) grpc.StreamClientInterceptor {
	interceptor := NewAPIKeyAuthInterceptor(apiKey, true)
	return interceptor.StreamInterceptor()
}

// CombinedAuthInterceptor creates interceptors that support both authorization and x-api-key headers.
type CombinedAuthInterceptor struct {
	apiKey         string
	useXAPIKey     bool
	useBearerToken bool
}

// NewCombinedAuthInterceptor creates a new combined auth interceptor.
func NewCombinedAuthInterceptor(apiKey string, useXAPIKey, useBearerToken bool) *CombinedAuthInterceptor {
	return &CombinedAuthInterceptor{
		apiKey:         apiKey,
		useXAPIKey:     useXAPIKey,
		useBearerToken: useBearerToken,
	}
}

// UnaryInterceptor returns a unary interceptor function.
func (a *CombinedAuthInterceptor) UnaryInterceptor() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		authCtx, err := a.addAuthMetadata(ctx)
		if err != nil {
			return status.Error(codes.Unauthenticated, fmt.Sprintf("failed to add auth metadata: %v", err))
		}
		return invoker(authCtx, method, req, reply, cc, opts...)
	}
}

// StreamInterceptor returns a stream interceptor function.
func (a *CombinedAuthInterceptor) StreamInterceptor() grpc.StreamClientInterceptor {
	return func(
		ctx context.Context,
		desc *grpc.StreamDesc,
		cc *grpc.ClientConn,
		method string,
		streamer grpc.Streamer,
		opts ...grpc.CallOption,
	) (grpc.ClientStream, error) {
		authCtx, err := a.addAuthMetadata(ctx)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, fmt.Sprintf("failed to add auth metadata: %v", err))
		}
		return streamer(authCtx, desc, cc, method, opts...)
	}
}

// addAuthMetadata for combined interceptor adds both auth types if configured.
// Uses AppendToOutgoingContext to preserve gRPC's internal metadata.
func (a *CombinedAuthInterceptor) addAuthMetadata(ctx context.Context) (context.Context, error) {
	if a.apiKey == "" {
		return nil, fmt.Errorf("API key is empty")
	}

	// Build key-value pairs for AppendToOutgoingContext
	pairs := make([]string, 0, 4)

	if a.useXAPIKey {
		pairs = append(pairs, XAPIKeyHeader, a.apiKey)
	}

	if a.useBearerToken {
		pairs = append(pairs, AuthorizationHeader, fmt.Sprintf("Bearer %s", a.apiKey))
	}

	return metadata.AppendToOutgoingContext(ctx, pairs...), nil
}

// EmptyAuthInterceptor returns interceptors that don't add any auth (for testing).
func EmptyAuthInterceptor() (grpc.UnaryClientInterceptor, grpc.StreamClientInterceptor) {
	unary := func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		return invoker(ctx, method, req, reply, cc, opts...)
	}

	stream := func(
		ctx context.Context,
		desc *grpc.StreamDesc,
		cc *grpc.ClientConn,
		method string,
		streamer grpc.Streamer,
		opts ...grpc.CallOption,
	) (grpc.ClientStream, error) {
		return streamer(ctx, desc, cc, method, opts...)
	}

	return unary, stream
}
