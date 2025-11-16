package metadata

import (
	"context"
	"testing"

	"google.golang.org/grpc/metadata"
)

func TestNewSDKMetadata(t *testing.T) {
	apiKey := "test-api-key"
	md := NewSDKMetadata(apiKey)

	if md == nil {
		t.Error("NewSDKMetadata should not return nil")
	}

	if md.APIKey != apiKey {
		t.Errorf("Expected API key %s, got %s", apiKey, md.APIKey)
	}

	if md.ClientVersion == "" {
		t.Error("Client version should not be empty")
	}

	if md.Environment == "" {
		t.Error("Environment should not be empty")
	}

	// Check default values
	expectedClientVersion := "xai-sdk-go/0.2.1"
	if md.ClientVersion != expectedClientVersion {
		t.Errorf("Expected client version %s, got %s", expectedClientVersion, md.ClientVersion)
	}

	expectedEnvironment := "production"
	if md.Environment != expectedEnvironment {
		t.Errorf("Expected environment %s, got %s", expectedEnvironment, md.Environment)
	}
}

func TestSDKMetadataToMetadata(t *testing.T) {
	md := &SDKMetadata{
		APIKey:         "test-api-key",
		RequestID:      "req-123",
		ClientVersion:  "xai-sdk-go/1.0.0",
		Environment:    "test",
		ConversationID: "conv-456",
		UserAgent:      "test-agent",
	}

	grpcMD := md.ToMetadata()

	// Test that all fields are converted correctly
	if values := grpcMD.Get(APIKeyKey); len(values) != 1 || values[0] != "test-api-key" {
		t.Error("API key should be converted correctly")
	}

	if values := grpcMD.Get(RequestIDKey); len(values) != 1 || values[0] != "req-123" {
		t.Error("Request ID should be converted correctly")
	}

	if values := grpcMD.Get(ClientVersionKey); len(values) != 1 || values[0] != "xai-sdk-go/1.0.0" {
		t.Error("Client version should be converted correctly")
	}

	if values := grpcMD.Get(EnvironmentKey); len(values) != 1 || values[0] != "test" {
		t.Error("Environment should be converted correctly")
	}

	if values := grpcMD.Get(ConversationIDKey); len(values) != 1 || values[0] != "conv-456" {
		t.Error("Conversation ID should be converted correctly")
	}

	if values := grpcMD.Get(UserAgentKey); len(values) != 1 || values[0] != "test-agent" {
		t.Error("User agent should be converted correctly")
	}
}

func TestSDKMetadataToMetadataWithEmptyFields(t *testing.T) {
	md := &SDKMetadata{
		APIKey: "test-api-key",
		// All other fields are empty
	}

	grpcMD := md.ToMetadata()

	// Should only contain API key
	if values := grpcMD.Get(APIKeyKey); len(values) != 1 || values[0] != "test-api-key" {
		t.Error("API key should be present")
	}

	if values := grpcMD.Get(RequestIDKey); len(values) != 0 {
		t.Error("Empty request ID should not be included")
	}

	if values := grpcMD.Get(ClientVersionKey); len(values) != 0 {
		t.Error("Empty client version should not be included")
	}

	if values := grpcMD.Get(EnvironmentKey); len(values) != 0 {
		t.Error("Empty environment should not be included")
	}
}

func TestAddToOutgoingContext(t *testing.T) {
	md := &SDKMetadata{
		APIKey: "test-api-key",
		ClientVersion: "test-version",
	}

	ctx := context.Background()
	newCtx := md.AddToOutgoingContext(ctx)

	// Extract metadata from the new context
	outgoingMD, ok := metadata.FromOutgoingContext(newCtx)
	if !ok {
		t.Error("Outgoing metadata should be present in context")
	}

	if values := outgoingMD.Get(APIKeyKey); len(values) != 1 || values[0] != "test-api-key" {
		t.Error("API key should be in outgoing context")
	}

	if values := outgoingMD.Get(ClientVersionKey); len(values) != 1 || values[0] != "test-version" {
		t.Error("Client version should be in outgoing context")
	}
}

func TestExtractFromIncomingContext(t *testing.T) {
	md := &SDKMetadata{
		APIKey: "initial-api-key",
	}

	// Create incoming metadata
	incomingMD := metadata.MD{}
	incomingMD.Set(APIKeyKey, "new-api-key")
	incomingMD.Set(RequestIDKey, "req-123")
	incomingMD.Set(ClientVersionKey, "new-version")
	incomingMD.Set(EnvironmentKey, "staging")

	ctx := metadata.NewIncomingContext(context.Background(), incomingMD)
	md.ExtractFromIncomingContext(ctx)

	// Check that fields were updated
	if md.APIKey != "new-api-key" {
		t.Errorf("API key should be updated to 'new-api-key', got %s", md.APIKey)
	}

	if md.RequestID != "req-123" {
		t.Errorf("Request ID should be 'req-123', got %s", md.RequestID)
	}

	if md.ClientVersion != "new-version" {
		t.Errorf("Client version should be 'new-version', got %s", md.ClientVersion)
	}

	if md.Environment != "staging" {
		t.Errorf("Environment should be 'staging', got %s", md.Environment)
	}
}

func TestGetFromContext(t *testing.T) {
	md := metadata.MD{}
	md.Set(APIKeyKey, "test-key")
	md.Set(RequestIDKey, "req-123")
	md.Set(RequestIDKey, "req-456") // Set multiple values

	ctx := metadata.NewIncomingContext(context.Background(), md)

	// Test getting API key
	apiKeys := GetFromContext(ctx, APIKeyKey)
	if len(apiKeys) != 1 || apiKeys[0] != "test-key" {
		t.Error("Should get single API key value")
	}

	// Test getting request ID (multiple values)
	requestIDs := GetFromContext(ctx, RequestIDKey)
	// Note: gRPC metadata may handle multiple values differently
	if len(requestIDs) < 1 || len(requestIDs) > 2 {
		t.Errorf("Should get 1-2 request ID values, got %d", len(requestIDs))
	}

	// Test getting non-existent key
	nonExistent := GetFromContext(ctx, "non-existent")
	if len(nonExistent) != 0 {
		t.Error("Should get empty slice for non-existent key")
	}
}

func TestGetSingleFromContext(t *testing.T) {
	md := metadata.MD{}
	md.Set(APIKeyKey, "test-key")
	md.Set(RequestIDKey, "req-123")
	md.Set(RequestIDKey, "req-456")

	ctx := metadata.NewIncomingContext(context.Background(), md)

	// Test getting single value that exists
	apiKey, ok := GetSingleFromContext(ctx, APIKeyKey)
	if !ok {
		t.Error("Should find API key")
	}
	if apiKey != "test-key" {
		t.Errorf("Expected API key 'test-key', got %s", apiKey)
	}

	// Test getting single value with multiple entries (behavior may vary)
	requestID, ok := GetSingleFromContext(ctx, RequestIDKey)
	if !ok {
		t.Error("Should find request ID")
	}
	// Either the first or second value is acceptable
	if requestID != "req-123" && requestID != "req-456" {
		t.Errorf("Expected request ID to be 'req-123' or 'req-456', got %s", requestID)
	}

	// Test getting non-existent value
	nonExistent, ok := GetSingleFromContext(ctx, "non-existent")
	if ok {
		t.Error("Should not find non-existent value")
	}
	if nonExistent != "" {
		t.Error("Non-existent value should be empty string")
	}
}

func TestSpecificGetters(t *testing.T) {
	md := metadata.MD{}
	md.Set(APIKeyKey, "test-key")
	md.Set(RequestIDKey, "req-123")
	md.Set(ClientVersionKey, "v1.0.0")
	md.Set(EnvironmentKey, "test")
	md.Set(ConversationIDKey, "conv-456")
	md.Set(UserAgentKey, "test-agent")

	ctx := metadata.NewIncomingContext(context.Background(), md)

	// Test GetRequestID
	requestID, ok := GetRequestID(ctx)
	if !ok {
		t.Error("Should find request ID")
	}
	if requestID != "req-123" {
		t.Errorf("Expected request ID 'req-123', got %s", requestID)
	}

	// Test GetConversationID
	convID, ok := GetConversationID(ctx)
	if !ok {
		t.Error("Should find conversation ID")
	}
	if convID != "conv-456" {
		t.Errorf("Expected conversation ID 'conv-456', got %s", convID)
	}

	// Test GetClientVersion
	version, ok := GetClientVersion(ctx)
	if !ok {
		t.Error("Should find client version")
	}
	if version != "v1.0.0" {
		t.Errorf("Expected client version 'v1.0.0', got %s", version)
	}

	// Test GetEnvironment
	env, ok := GetEnvironment(ctx)
	if !ok {
		t.Error("Should find environment")
	}
	if env != "test" {
		t.Errorf("Expected environment 'test', got %s", env)
	}

	// Test GetAPIKey
	apiKey, ok := GetAPIKey(ctx)
	if !ok {
		t.Error("Should find API key")
	}
	if apiKey != "test-key" {
		t.Errorf("Expected API key 'test-key', got %s", apiKey)
	}

	// Test GetUserAgent
	userAgent, ok := GetUserAgent(ctx)
	if !ok {
		t.Error("Should find user agent")
	}
	if userAgent != "test-agent" {
		t.Errorf("Expected user agent 'test-agent', got %s", userAgent)
	}
}

func TestExtractClientInfo(t *testing.T) {
	md := metadata.MD{}
	md.Set(ClientVersionKey, "xai-sdk-go/1.0.0")
	md.Set(UserAgentKey, "test-agent")
	md.Set(RequestIDKey, "req-123")
	md.Set(EnvironmentKey, "test")

	ctx := metadata.NewIncomingContext(context.Background(), md)

	info, err := ExtractClientInfo(ctx)
	if err != nil {
		t.Errorf("ExtractClientInfo should not return error, got: %v", err)
	}

	// Check expected fields
	if info["client_version"] != "xai-sdk-go/1.0.0" {
		t.Errorf("Expected client_version in info, got: %v", info)
	}

	if info["user_agent"] != "test-agent" {
		t.Errorf("Expected user_agent in info, got: %v", info)
	}

	if info["request_id"] != "req-123" {
		t.Errorf("Expected request_id in info, got: %v", info)
	}

	if info["environment"] != "test" {
		t.Errorf("Expected environment in info, got: %v", info)
	}
}

func TestValidateMetadata(t *testing.T) {
	// Test valid metadata
	validMD := &SDKMetadata{
		APIKey:        "test-key",
		ClientVersion: "test-version",
	}

	err := ValidateMetadata(validMD)
	if err != nil {
		t.Errorf("Valid metadata should not return error, got: %v", err)
	}

	// Test nil metadata
	err = ValidateMetadata(nil)
	if err == nil {
		t.Error("Nil metadata should return error")
	}

	// Test metadata with empty API key
	invalidMD := &SDKMetadata{
		ClientVersion: "test-version",
		// APIKey is empty
	}

	err = ValidateMetadata(invalidMD)
	if err == nil {
		t.Error("Metadata with empty API key should return error")
	}

	// Test metadata with empty client version
	invalidMD2 := &SDKMetadata{
		APIKey: "test-key",
		// ClientVersion is empty
	}

	err = ValidateMetadata(invalidMD2)
	if err == nil {
		t.Error("Metadata with empty client version should return error")
	}
}

func TestSanitizeMetadata(t *testing.T) {
	md := &SDKMetadata{
		APIKey:         "1234567890abcdef",
		RequestID:      "req-123",
		ClientVersion:  "xai-sdk-go/1.0.0",
		Environment:    "test",
		ConversationID: "conv-456",
		UserAgent:      "test-agent",
	}

	sanitized := SanitizeMetadata(md)

	// API key should be masked except last 4 characters
	if sanitized["api_key"] != "************cdef" {
		t.Errorf("API key should be masked, got: %s", sanitized["api_key"])
	}

	// Other fields should be preserved
	if sanitized["request_id"] != "req-123" {
		t.Errorf("Request ID should be preserved, got: %s", sanitized["request_id"])
	}

	if sanitized["client_version"] != "xai-sdk-go/1.0.0" {
		t.Errorf("Client version should be preserved, got: %s", sanitized["client_version"])
	}

	if sanitized["environment"] != "test" {
		t.Errorf("Environment should be preserved, got: %s", sanitized["environment"])
	}

	if sanitized["conversation_id"] != "conv-456" {
		t.Errorf("Conversation ID should be preserved, got: %s", sanitized["conversation_id"])
	}

	if sanitized["user_agent"] != "test-agent" {
		t.Errorf("User agent should be preserved, got: %s", sanitized["user_agent"])
	}
}

func TestBuildCommonHeaders(t *testing.T) {
	md := &SDKMetadata{
		APIKey:         "test-api-key",
		RequestID:      "req-123",
		ClientVersion:  "xai-sdk-go/1.0.0",
		Environment:    "test",
		ConversationID: "conv-456",
		UserAgent:      "test-agent",
	}

	headers := md.BuildCommonHeaders()

	// Check Authorization header
	if headers[AuthorizationKey] != "Bearer test-api-key" {
		t.Errorf("Expected Authorization header 'Bearer test-api-key', got: %s", headers[AuthorizationKey])
	}

	// Check other headers
	if headers[RequestIDKey] != "req-123" {
		t.Errorf("Expected RequestID header 'req-123', got: %s", headers[RequestIDKey])
	}

	if headers[UserAgentKey] != "test-agent" {
		t.Errorf("Expected UserAgent header 'test-agent', got: %s", headers[UserAgentKey])
	}

	if headers[ClientVersionKey] != "xai-sdk-go/1.0.0" {
		t.Errorf("Expected ClientVersion header 'xai-sdk-go/1.0.0', got: %s", headers[ClientVersionKey])
	}

	if headers[EnvironmentKey] != "test" {
		t.Errorf("Expected Environment header 'test', got: %s", headers[EnvironmentKey])
	}

	if headers[ConversationIDKey] != "conv-456" {
		t.Errorf("Expected ConversationID header 'conv-456', got: %s", headers[ConversationIDKey])
	}
}

func TestBuildCommonHeadersWithEmptyFields(t *testing.T) {
	md := &SDKMetadata{
		APIKey: "test-api-key",
		// All other fields are empty
	}

	headers := md.BuildCommonHeaders()

	// Should only contain Authorization header (API key)
	if headers[AuthorizationKey] != "Bearer test-api-key" {
		t.Error("Authorization header should be present")
	}

	// Should also contain default values for ClientVersion and Environment
	if len(headers) < 3 {
		t.Error("Should contain at least Authorization, ClientVersion, and Environment headers")
	}
}

func TestGetPeerInfo(t *testing.T) {
	// This test is limited because we can't easily mock peer information
	// But we can test that the function doesn't panic
	ctx := context.Background()
	peerInfo, ok := GetPeerInfo(ctx)
	
	// In a normal context without peer info, should return false
	if ok {
		t.Log("Peer info is available in this context")
	} else {
		t.Log("No peer info available in this context")
	}

	_ = peerInfo // Use the variable to avoid unused variable error
}

func BenchmarkSDKMetadataToMetadata(b *testing.B) {
	md := &SDKMetadata{
		APIKey:         "test-api-key",
		RequestID:      "req-123",
		ClientVersion:  "xai-sdk-go/1.0.0",
		Environment:    "test",
		ConversationID: "conv-456",
		UserAgent:      "test-agent",
	}

	for i := 0; i < b.N; i++ {
		_ = md.ToMetadata()
	}
}

func BenchmarkAddToOutgoingContext(b *testing.B) {
	md := &SDKMetadata{
		APIKey:        "test-api-key",
		ClientVersion: "xai-sdk-go/1.0.0",
	}

	ctx := context.Background()
	for i := 0; i < b.N; i++ {
		_ = md.AddToOutgoingContext(ctx)
	}
}