package errors

import (
	"errors"
	"io"
	"strings"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestSentinelErrors(t *testing.T) {
	tests := []struct {
		name string
		err  error
	}{
		{"ErrInvalidConfig", ErrInvalidConfig},
		{"ErrInvalidAPIKey", ErrInvalidAPIKey},
		{"ErrUnauthorized", ErrUnauthorized},
		{"ErrPermissionDenied", ErrPermissionDenied},
		{"ErrNotFound", ErrNotFound},
		{"ErrRateLimit", ErrRateLimit},
		{"ErrQuotaExceeded", ErrQuotaExceeded},
		{"ErrInternal", ErrInternal},
		{"ErrServiceUnavailable", ErrServiceUnavailable},
		{"ErrDeadlineExceeded", ErrDeadlineExceeded},
		{"ErrCanceled", ErrCanceled},
		{"ErrInvalidRequest", ErrInvalidRequest},
		{"ErrTooManyRequests", ErrTooManyRequests},
		{"ErrConnection", ErrConnection},
		{"ErrTimeout", ErrTimeout},
		{"ErrParsing", ErrParsing},
		{"ErrValidation", ErrValidation},
		{"ErrFileTooLarge", ErrFileTooLarge},
		{"ErrUnsupportedContentType", ErrUnsupportedContentType},
		{"ErrInvalidModel", ErrInvalidModel},
		{"ErrInvalidParameters", ErrInvalidParameters},
		{"ErrStreamClosed", ErrStreamClosed},
		{"ErrMalformedData", ErrMalformedData},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err == nil {
				t.Error("Sentinel error should not be nil")
			}
			if tt.err.Error() == "" {
				t.Error("Sentinel error message should not be empty")
			}
		})
	}
}

func TestNewError(t *testing.T) {
	err := NewError(ErrorTypeAPI, codes.InvalidArgument, "test error")

	if err.Type() != ErrorTypeAPI {
		t.Errorf("Expected error type %s, got %s", ErrorTypeAPI, err.Type())
	}

	if err.Code() != codes.InvalidArgument {
		t.Errorf("Expected error code %v, got %v", codes.InvalidArgument, err.Code())
	}

	if err.Error() == "" {
		t.Error("Error message should not be empty")
	}

	if err.Unwrap() != nil {
		t.Error("Error cause should be nil")
	}

	if err.Context() == nil {
		t.Error("Error context should not be nil")
	}

	if err.StackTrace() == nil {
		t.Error("Error stack trace should not be nil")
	}
}

func TestNewErrorWithCause(t *testing.T) {
	cause := errors.New("original error")
	err := NewErrorWithCause(ErrorTypeAuth, codes.Unauthenticated, "auth failed", cause)

	if err.Type() != ErrorTypeAuth {
		t.Errorf("Expected error type %s, got %s", ErrorTypeAuth, err.Type())
	}

	if err.Code() != codes.Unauthenticated {
		t.Errorf("Expected error code %v, got %v", codes.Unauthenticated, err.Code())
	}

	if err.Error() == "" {
		t.Error("Error message should not be empty")
	}

	if err.Unwrap() == nil {
		t.Error("Error cause should not be nil")
	}

	if err.Unwrap() != cause {
		t.Error("Error cause should match the provided cause")
	}
}

func TestNewErrorWithContext(t *testing.T) {
	context := map[string]interface{}{
		"user_id":      "12345",
		"request_path": "/api/test",
	}

	err := NewErrorWithContext(ErrorTypeValidation, codes.InvalidArgument, "validation failed", context)

	if err.Type() != ErrorTypeValidation {
		t.Errorf("Expected error type %s, got %s", ErrorTypeValidation, err.Type())
	}

	if err.Code() != codes.InvalidArgument {
		t.Errorf("Expected error code %v, got %v", codes.InvalidArgument, err.Code())
	}

	if err.Error() == "" {
		t.Error("Error message should not be empty")
	}

	if err.Unwrap() != nil {
		t.Error("Error cause should be nil")
	}

	if len(err.Context()) != 2 {
		t.Error("Error context should have 2 entries")
	}

	if err.Context()["user_id"] != "12345" {
		t.Errorf("Expected user_id to be '12345', got %v", err.Context()["user_id"])
	}

	if err.Context()["request_path"] != "/api/test" {
		t.Errorf("Expected request_path to be '/api/test', got %v", err.Context()["request_path"])
	}
}

func TestNewErrorFull(t *testing.T) {
	cause := errors.New("original error")
	context := map[string]interface{}{
		"retry_count": 3,
	}

	err := NewErrorFull(ErrorTypeNetwork, codes.Unavailable, "connection failed", cause, context)

	if err.Type() != ErrorTypeNetwork {
		t.Errorf("Expected error type %s, got %s", ErrorTypeNetwork, err.Type())
	}

	if err.Code() != codes.Unavailable {
		t.Errorf("Expected error code %v, got %v", codes.Unavailable, err.Code())
	}

	if err.Error() == "" {
		t.Error("Error message should not be empty")
	}

	if err.Unwrap() == nil {
		t.Error("Error cause should not be nil")
	}

	if err.Unwrap() != cause {
		t.Error("Error cause should match the provided cause")
	}

	if err.Context()["retry_count"] != 3 {
		t.Errorf("Expected retry_count to be 3, got %v", err.Context()["retry_count"])
	}
}

func TestErrorString(t *testing.T) {
	cause := errors.New("original error")
	context := map[string]interface{}{
		"user_id":   "12345",
		"operation": "upload",
	}

	err := NewErrorFull(ErrorTypeAPI, codes.NotFound, "resource not found", cause, context)
	msg := err.Error()

	// Check that the error message contains expected parts
	if !strings.Contains(msg, "[api]") {
		t.Errorf("Error message should contain '[api]', got: %s", msg)
	}

	if !strings.Contains(msg, "resource not found") {
		t.Errorf("Error message should contain 'resource not found', got: %s", msg)
	}

	if !strings.Contains(msg, "cause:") {
		t.Errorf("Error message should contain 'cause:', got: %s", msg)
	}

	if !strings.Contains(msg, "original error") {
		t.Errorf("Error message should contain 'original error', got: %s", msg)
	}

	if !strings.Contains(msg, "context:") {
		t.Errorf("Error message should contain 'context:', got: %s", msg)
	}
}

func TestErrorUnwrap(t *testing.T) {
	cause := errors.New("original error")
	err := NewErrorWithCause(ErrorTypeAPI, codes.InvalidArgument, "test error", cause)

	unwrapped := err.Unwrap()
	if unwrapped == nil {
		t.Error("Unwrap should return the cause error")
	}

	if unwrapped != cause {
		t.Error("Unwrap should return the exact cause error")
	}
}

func TestErrorMethods(t *testing.T) {
	err := NewError(ErrorTypeAuth, codes.Unauthenticated, "unauthorized")

	if err.Type() != ErrorTypeAuth {
		t.Errorf("Type() should return %s, got %s", ErrorTypeAuth, err.Type())
	}

	if err.Code() != codes.Unauthenticated {
		t.Errorf("Code() should return %v, got %v", codes.Unauthenticated, err.Code())
	}

	// Context should be initialized (even if empty)
	if err.Context() == nil {
		t.Error("Context() should not be nil")
	}

	if err.StackTrace() == nil {
		t.Error("StackTrace() should not be nil")
	}
}

func TestErrorGRPCStatus(t *testing.T) {
	err := NewError(ErrorTypeAPI, codes.NotFound, "not found")

	status := err.GRPCStatus()
	if status == nil {
		t.Error("GRPCStatus() should not return nil")
	}

	if status.Code() != codes.NotFound {
		t.Errorf("GRPC status code should be %v, got %v", codes.NotFound, status.Code())
	}

	if status.Message() != "not found" {
		t.Errorf("GRPC status message should be 'not found', got '%s'", status.Message())
	}
}

func TestFromGRPC(t *testing.T) {
	// Test with nil error
	result := FromGRPC(nil)
	if result != nil {
		t.Error("FromGRPC(nil) should return nil")
	}

	// Test with gRPC status error
	st := status.New(codes.NotFound, "resource not found")
	grpcErr := st.Err()
	result = FromGRPC(grpcErr)

	if result == nil {
		t.Error("FromGRPC should not return nil for gRPC error")
	}

	// Test with io.EOF
	eofErr := io.EOF
	result = FromGRPC(eofErr)
	if result != io.EOF {
		t.Error("FromGRPC should return io.EOF as-is")
	}

	// Test with context canceled error
	cancelErr := errors.New("context canceled")
	result = FromGRPC(cancelErr)
	if result == nil {
		t.Error("FromGRPC should not return nil for context error")
	}

	// Test with generic error
	genericErr := errors.New("generic error")
	result = FromGRPC(genericErr)
	if result == nil {
		t.Error("FromGRPC should not return nil for generic error")
	}
}

func TestFromGRPCAlreadyErrorType(t *testing.T) {
	// Test that if we pass an Error type, it returns as-is
	originalErr := NewError(ErrorTypeAPI, codes.InvalidArgument, "test error")
	result := FromGRPC(originalErr)

	if result != originalErr {
		t.Error("FromGRPC should return the same Error instance when passed an Error type")
	}
}

func TestFromGRPCWithWrappedError(t *testing.T) {
	// Test with a wrapped gRPC error
	cause := status.New(codes.Unauthenticated, "auth failed").Err()
	// Note: This would require errors.Wrap from github.com/pkg/errors package
	// For now, test the basic functionality
	result := FromGRPC(cause)

	// Should convert to SDK error
	if _, ok := result.(*Error); !ok {
		t.Error("FromGRPC should convert gRPC error to SDK error")
	}
}

func TestMapGRPCCodeToErrorType(t *testing.T) {
	tests := []struct {
		code codes.Code
		want ErrorType
	}{
		{codes.InvalidArgument, ErrorTypeValidation},
		{codes.Unauthenticated, ErrorTypeAuth},
		{codes.PermissionDenied, ErrorTypeAuth},
		{codes.NotFound, ErrorTypeAPI},
		{codes.AlreadyExists, ErrorTypeAPI},
		{codes.ResourceExhausted, ErrorTypeRateLimit},
		{codes.FailedPrecondition, ErrorTypeValidation},
		{codes.Aborted, ErrorTypeAPI},
		{codes.OutOfRange, ErrorTypeValidation},
		{codes.Unimplemented, ErrorTypeAPI},
		{codes.Internal, ErrorTypeInternal},
		{codes.Unavailable, ErrorTypeService},
		{codes.DataLoss, ErrorTypeInternal},
		{codes.DeadlineExceeded, ErrorTypeTimeout},
		{codes.Canceled, ErrorTypeCanceled},
		{codes.Unknown, ErrorTypeAPI},
	}

	for _, tt := range tests {
		t.Run(tt.code.String(), func(t *testing.T) {
			got := mapGRPCCodeToErrorType(tt.code)
			if got != tt.want {
				t.Errorf("Expected error type %s for code %s, got %s", tt.want, tt.code, got)
			}
		})
	}
}

func TestHelperFunctions(t *testing.T) {
	// Test NewConfigError
	err := NewConfigError("invalid config")
	if err.Type() != ErrorTypeConfig {
		t.Errorf("NewConfigError should create config error, got %s", err.Type())
	}

	// Test NewAuthError
	err = NewAuthError("invalid api key")
	if err.Type() != ErrorTypeAuth {
		t.Errorf("NewAuthError should create auth error, got %s", err.Type())
	}

	// Test NewAPIError
	err = NewAPIError(codes.NotFound, "not found")
	if err.Type() != ErrorTypeAPI {
		t.Errorf("NewAPIError should create API error, got %s", err.Type())
	}

	// Test NewNetworkError
	err = NewNetworkError("connection failed")
	if err.Type() != ErrorTypeNetwork {
		t.Errorf("NewNetworkError should create network error, got %s", err.Type())
	}

	// Test NewValidationError
	err = NewValidationError("invalid parameters")
	if err.Type() != ErrorTypeValidation {
		t.Errorf("NewValidationError should create validation error, got %s", err.Type())
	}
}

func TestGetStackTrace(t *testing.T) {
	stack := getStackTrace()

	if stack == nil {
		t.Error("Stack trace should not be nil")
	}

	if len(stack) == 0 {
		t.Error("Stack trace should contain at least one frame")
	}

	// Check that stack frames contain expected information
	for _, frame := range stack {
		if frame == "" {
			t.Error("Stack frame should not be empty")
		}

		// Should contain file:line:function format
		if !strings.Contains(frame, ":") {
			t.Errorf("Stack frame should contain ':' separator, got: %s", frame)
		}
	}
}

func TestErrorTypes(t *testing.T) {
	allErrorTypes := []ErrorType{
		ErrorTypeConfig,
		ErrorTypeAuth,
		ErrorTypeAPI,
		ErrorTypeNetwork,
		ErrorTypeValidation,
		ErrorTypeRateLimit,
		ErrorTypeQuota,
		ErrorTypeInternal,
		ErrorTypeService,
		ErrorTypeTimeout,
		ErrorTypeCanceled,
		ErrorTypeStream,
		ErrorTypeParsing,
		ErrorTypeFile,
	}

	for _, errorType := range allErrorTypes {
		if string(errorType) == "" {
			t.Error("Error type should not be empty string")
		}
	}
}

func BenchmarkNewError(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = NewError(ErrorTypeAPI, codes.InvalidArgument, "benchmark error")
	}
}

func BenchmarkFromGRPC(b *testing.B) {
	grpcErr := status.New(codes.InvalidArgument, "grpc error").Err()

	for i := 0; i < b.N; i++ {
		_ = FromGRPC(grpcErr)
	}
}

func BenchmarkErrorString(b *testing.B) {
	err := NewErrorFull(ErrorTypeAPI, codes.InvalidArgument, "test error with context",
		errors.New("cause"),
		map[string]interface{}{"key": "value"})

	for i := 0; i < b.N; i++ {
		_ = err.Error()
	}
}
