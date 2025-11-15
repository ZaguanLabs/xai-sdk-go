// Package constants provides default constants used throughout the xAI SDK.
package constants

import (
	"time"
)

// Default service endpoints and hosts
const (
	// DefaultAPIV1Host is the default API v1 host.
	DefaultAPIV1Host = "api.x.ai"

	// DefaultGRPCPort is the default gRPC port.
	DefaultGRPCPort = "443"

	// DefaultHTTPHost is the default HTTP host.
	DefaultHTTPHost = "api.x.ai"

	// DefaultChatEndpoint is the default chat endpoint path.
	DefaultChatEndpoint = "/v1/chat/completions"

	// DefaultFilesEndpoint is the default files endpoint path.
	DefaultFilesEndpoint = "/v1/files"

	// DefaultImageEndpoint is the default image endpoint path.
	DefaultImageEndpoint = "/v1/images"

	// DefaultModelsEndpoint is the default models endpoint path.
	DefaultModelsEndpoint = "/v1/models"

	// DefaultTokenizerEndpoint is the default tokenizer endpoint path.
	DefaultTokenizerEndpoint = "/v1/tokenizer"

	// DefaultCollectionsEndpoint is the default collections endpoint path.
	DefaultCollectionsEndpoint = "/v1/collections"

	// DefaultAuthEndpoint is the default auth endpoint path.
	DefaultAuthEndpoint = "/v1/auth"
)

// Default timeout values
const (
	// DefaultTimeout is the default request timeout.
	DefaultTimeout = 30 * time.Second

	// DefaultConnectTimeout is the default connection timeout.
	DefaultConnectTimeout = 10 * time.Second

	// DefaultKeepAliveTimeout is the default keep-alive timeout.
	DefaultKeepAliveTimeout = 20 * time.Second

	// DefaultStreamTimeout is the default streaming timeout.
	DefaultStreamTimeout = 300 * time.Second
)

// Default size and limit values
const (
	// DefaultMaxRetries is the default maximum number of retries.
	DefaultMaxRetries = 3

	// DefaultRetryBackoff is the default retry backoff duration.
	DefaultRetryBackoff = 1 * time.Second

	// DefaultMaxBackoff is the default maximum backoff duration.
	DefaultMaxBackoff = 60 * time.Second

	// DefaultChunkSize is the default chunk size for file uploads.
	DefaultChunkSize = 3 * 1024 * 1024 // 3 MiB

	// DefaultMaxChunkSize is the maximum chunk size for file uploads.
	DefaultMaxChunkSize = 10 * 1024 * 1024 // 10 MiB

	// DefaultBufferSize is the default buffer size for streaming.
	DefaultBufferSize = 64 * 1024 // 64 KiB

	// DefaultMaxImageSize is the default maximum image size in bytes.
	DefaultMaxImageSize = 10 * 1024 * 1024 // 10 MiB

	// DefaultMaxFileSize is the default maximum file size in bytes.
	DefaultMaxFileSize = 100 * 1024 * 1024 // 100 MiB

	// DefaultMaxTokens is the default maximum number of tokens.
	DefaultMaxTokens = 4096

	// DefaultMaxPromptTokens is the default maximum number of prompt tokens.
	DefaultMaxPromptTokens = 8192

	// DefaultMaxToolResults is the default maximum number of tool results.
	DefaultMaxToolResults = 50

	// DefaultMaxToolDescriptions is the default maximum number of tool descriptions.
	DefaultMaxToolDescriptions = 200
)

// Header names
const (
	// HeaderAuthorization is the Authorization header name.
	HeaderAuthorization = "Authorization"

	// HeaderContentType is the Content-Type header name.
	HeaderContentType = "Content-Type"

	// HeaderUserAgent is the User-Agent header name.
	HeaderUserAgent = "User-Agent"

	// HeaderXAIClient is the X-AI-Client header name.
	HeaderXAIClient = "X-AI-Client"

	// HeaderXAIRequestID is the X-AI-Request-ID header name.
	HeaderXAIRequestID = "X-AI-Request-ID"

	// HeaderXAIAPIVersion is the X-AI-API-Version header name.
	HeaderXAIAPIVersion = "X-AI-API-Version"

	// HeaderXAIEnvironment is the X-AI-Environment header name.
	HeaderXAIEnvironment = "X-AI-Environment"

	// HeaderXAIConversationID is the X-AI-Conversation-ID header name.
	HeaderXAIConversationID = "X-AI-Conversation-ID"
)

// Content types
const (
	// ContentTypeJSON is the JSON content type.
	ContentTypeJSON = "application/json"

	// ContentTypeProtobuf is the Protocol Buffers content type.
	ContentTypeProtobuf = "application/x-protobuf"

	// ContentTypeOctetStream is the octet stream content type.
	ContentTypeOctetStream = "application/octet-stream"

	// ContentTypeTextPlain is the plain text content type.
	ContentTypeTextPlain = "text/plain"

	// ContentTypeImageJPEG is the JPEG image content type.
	ContentTypeImageJPEG = "image/jpeg"

	// ContentTypeImagePNG is the PNG image content type.
	ContentTypeImagePNG = "image/png"
)

// User agent
const (
	// DefaultUserAgent is the default User-Agent string.
	DefaultUserAgent = "xai-sdk-go/" + "0.1.4"
)
