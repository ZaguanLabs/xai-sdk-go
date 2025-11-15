// Package metadata provides utilities for working with gRPC metadata and headers.
package metadata

import (
	"context"
	"fmt"
	"strings"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

// Common metadata keys used by the xAI SDK.
const (
	// APIKeyKey is the metadata key for the API key.
	APIKeyKey = "x-api-key"

	// RequestIDKey is the metadata key for the request ID.
	RequestIDKey = "x-request-id"

	// ClientVersionKey is the metadata key for the client version.
	ClientVersionKey = "x-client-version"

	// EnvironmentKey is the metadata key for the environment.
	EnvironmentKey = "x-environment"

	// ConversationIDKey is the metadata key for the conversation ID.
	ConversationIDKey = "x-conversation-id"

	// UserAgentKey is the metadata key for the user agent.
	UserAgentKey = "user-agent"

	// AuthorizationKey is the metadata key for authorization.
	AuthorizationKey = "authorization"

	// ContentTypeKey is the metadata key for content type.
	ContentTypeKey = "content-type"
)

// SDKMetadata contains common metadata for SDK requests.
type SDKMetadata struct {
	APIKey         string
	RequestID      string
	ClientVersion  string
	Environment    string
	ConversationID string
	UserAgent      string
}

// NewSDKMetadata creates a new SDKMetadata with default values.
func NewSDKMetadata(apiKey string) *SDKMetadata {
	return &SDKMetadata{
		APIKey:        apiKey,
		ClientVersion: "xai-sdk-go/0.1.2",
		Environment:   "production",
	}
}

// ToMetadata converts SDKMetadata to gRPC metadata.
func (m *SDKMetadata) ToMetadata() metadata.MD {
	md := metadata.MD{}

	if m.APIKey != "" {
		md.Set(APIKeyKey, m.APIKey)
	}

	if m.RequestID != "" {
		md.Set(RequestIDKey, m.RequestID)
	}

	if m.ClientVersion != "" {
		md.Set(ClientVersionKey, m.ClientVersion)
	}

	if m.Environment != "" {
		md.Set(EnvironmentKey, m.Environment)
	}

	if m.ConversationID != "" {
		md.Set(ConversationIDKey, m.ConversationID)
	}

	if m.UserAgent != "" {
		md.Set(UserAgentKey, m.UserAgent)
	}

	return md
}

// AddToOutgoingContext adds SDK metadata to the outgoing context.
// Uses AppendToOutgoingContext to preserve any existing metadata (e.g., gRPC's content-type).
func (m *SDKMetadata) AddToOutgoingContext(ctx context.Context) context.Context {
	// Convert to key-value pairs for AppendToOutgoingContext
	pairs := make([]string, 0)

	if m.APIKey != "" {
		pairs = append(pairs, APIKeyKey, m.APIKey)
	}

	if m.RequestID != "" {
		pairs = append(pairs, RequestIDKey, m.RequestID)
	}

	if m.ClientVersion != "" {
		pairs = append(pairs, ClientVersionKey, m.ClientVersion)
	}

	if m.Environment != "" {
		pairs = append(pairs, EnvironmentKey, m.Environment)
	}

	if m.ConversationID != "" {
		pairs = append(pairs, ConversationIDKey, m.ConversationID)
	}

	if m.UserAgent != "" {
		pairs = append(pairs, UserAgentKey, m.UserAgent)
	}

	return metadata.AppendToOutgoingContext(ctx, pairs...)
}

// ExtractFromIncomingContext extracts SDK metadata from the incoming context.
func (m *SDKMetadata) ExtractFromIncomingContext(ctx context.Context) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return
	}

	if values := md.Get(APIKeyKey); len(values) > 0 {
		m.APIKey = values[0]
	}

	if values := md.Get(RequestIDKey); len(values) > 0 {
		m.RequestID = values[0]
	}

	if values := md.Get(ClientVersionKey); len(values) > 0 {
		m.ClientVersion = values[0]
	}

	if values := md.Get(EnvironmentKey); len(values) > 0 {
		m.Environment = values[0]
	}

	if values := md.Get(ConversationIDKey); len(values) > 0 {
		m.ConversationID = values[0]
	}

	if values := md.Get(UserAgentKey); len(values) > 0 {
		m.UserAgent = values[0]
	}
}

// GetFromContext extracts specific metadata from the context.
func GetFromContext(ctx context.Context, key string) []string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil
	}
	return md.Get(key)
}

// GetSingleFromContext extracts a single metadata value from the context.
func GetSingleFromContext(ctx context.Context, key string) (string, bool) {
	values := GetFromContext(ctx, key)
	if len(values) == 0 {
		return "", false
	}
	return values[0], true
}

// GetRequestID retrieves the request ID from the context.
func GetRequestID(ctx context.Context) (string, bool) {
	return GetSingleFromContext(ctx, RequestIDKey)
}

// GetConversationID retrieves the conversation ID from the context.
func GetConversationID(ctx context.Context) (string, bool) {
	return GetSingleFromContext(ctx, ConversationIDKey)
}

// GetClientVersion retrieves the client version from the context.
func GetClientVersion(ctx context.Context) (string, bool) {
	return GetSingleFromContext(ctx, ClientVersionKey)
}

// GetEnvironment retrieves the environment from the context.
func GetEnvironment(ctx context.Context) (string, bool) {
	return GetSingleFromContext(ctx, EnvironmentKey)
}

// GetAPIKey retrieves the API key from the context.
func GetAPIKey(ctx context.Context) (string, bool) {
	return GetSingleFromContext(ctx, APIKeyKey)
}

// GetUserAgent retrieves the user agent from the context.
func GetUserAgent(ctx context.Context) (string, bool) {
	return GetSingleFromContext(ctx, UserAgentKey)
}

// GetPeerInfo retrieves peer information from the context.
func GetPeerInfo(ctx context.Context) (*peer.Peer, bool) {
	pr, ok := peer.FromContext(ctx)
	return pr, ok
}

// ExtractClientInfo extracts client information from metadata.
func ExtractClientInfo(ctx context.Context) (map[string]string, error) {
	info := make(map[string]string)

	// Extract peer information
	if pr, ok := GetPeerInfo(ctx); ok {
		if pr.Addr != nil {
			info["peer_address"] = pr.Addr.String()
		}
		info["peer_auth_type"] = pr.AuthInfo.AuthType()
	}

	// Extract client version
	if version, ok := GetClientVersion(ctx); ok {
		info["client_version"] = version
	}

	// Extract user agent
	if userAgent, ok := GetUserAgent(ctx); ok {
		info["user_agent"] = userAgent
	}

	// Extract request ID
	if requestID, ok := GetRequestID(ctx); ok {
		info["request_id"] = requestID
	}

	// Extract environment
	if env, ok := GetEnvironment(ctx); ok {
		info["environment"] = env
	}

	return info, nil
}

// ValidateMetadata validates required metadata fields.
func ValidateMetadata(md *SDKMetadata) error {
	if md == nil {
		return fmt.Errorf("metadata is nil")
	}

	if md.APIKey == "" {
		return fmt.Errorf("API key is required")
	}

	if md.ClientVersion == "" {
		return fmt.Errorf("client version is required")
	}

	return nil
}

// SanitizeMetadata removes sensitive information from metadata for logging.
func SanitizeMetadata(md *SDKMetadata) map[string]string {
	sanitized := make(map[string]string)

	if md.APIKey != "" {
		sanitized["api_key"] = strings.Repeat("*", len(md.APIKey)-4) + md.APIKey[len(md.APIKey)-4:]
	}

	sanitized["request_id"] = md.RequestID
	sanitized["client_version"] = md.ClientVersion
	sanitized["environment"] = md.Environment
	sanitized["conversation_id"] = md.ConversationID
	sanitized["user_agent"] = md.UserAgent

	return sanitized
}

// BuildCommonHeaders builds common HTTP headers from SDK metadata.
func (m *SDKMetadata) BuildCommonHeaders() map[string]string {
	headers := make(map[string]string)

	if m.APIKey != "" {
		headers[AuthorizationKey] = fmt.Sprintf("Bearer %s", m.APIKey)
	}

	if m.RequestID != "" {
		headers[RequestIDKey] = m.RequestID
	}

	if m.UserAgent != "" {
		headers[UserAgentKey] = m.UserAgent
	}

	headers[ClientVersionKey] = m.ClientVersion
	headers[EnvironmentKey] = m.Environment

	if m.ConversationID != "" {
		headers[ConversationIDKey] = m.ConversationID
	}

	return headers
}
