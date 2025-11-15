// Package grpcutil provides gRPC utility functions and interceptors.
package grpcutil

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// ContentTypeInterceptor creates a unary interceptor that adds a content-type header.
type ContentTypeInterceptor struct {
	contentType string
}

// NewContentTypeInterceptor creates a new content-type interceptor.
func NewContentTypeInterceptor(contentType string) *ContentTypeInterceptor {
	return &ContentTypeInterceptor{
		contentType: contentType,
	}
}

// UnaryInterceptor returns a unary interceptor function for adding content-type.
func (c *ContentTypeInterceptor) UnaryInterceptor() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		// Add content-type metadata to context
		ctx = metadata.AppendToOutgoingContext(ctx, "Content-Type", c.contentType)

		// Call the invoker with the modified context
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

// StreamInterceptor returns a stream interceptor function for adding content-type.
func (c *ContentTypeInterceptor) StreamInterceptor() grpc.StreamClientInterceptor {
	return func(
		ctx context.Context,
		desc *grpc.StreamDesc,
		cc *grpc.ClientConn,
		method string,
		streamer grpc.Streamer,
		opts ...grpc.CallOption,
	) (grpc.ClientStream, error) {
		// Add content-type metadata to context
		ctx = metadata.AppendToOutgoingContext(ctx, "Content-Type", c.contentType)

		// Call the streamer with the modified context
		return streamer(ctx, desc, cc, method, opts...)
	}
}
