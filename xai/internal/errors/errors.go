// Package errors provides error types and utilities for the xAI SDK.
package errors

import (
	"fmt"
	"io"
	"runtime"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Sentinel errors for common error conditions.
var (
	// ErrInvalidConfig indicates an invalid configuration.
	ErrInvalidConfig = fmt.Errorf("invalid configuration")

	// ErrInvalidAPIKey indicates an invalid or missing API key.
	ErrInvalidAPIKey = fmt.Errorf("invalid or missing API key")

	// ErrUnauthorized indicates the request was not authorized.
	ErrUnauthorized = fmt.Errorf("unauthorized")

	// ErrPermissionDenied indicates the request was denied due to insufficient permissions.
	ErrPermissionDenied = fmt.Errorf("permission denied")

	// ErrNotFound indicates the requested resource was not found.
	ErrNotFound = fmt.Errorf("not found")

	// ErrRateLimit indicates the request was rate limited.
	ErrRateLimit = fmt.Errorf("rate limited")

	// ErrQuotaExceeded indicates the request exceeded quota limits.
	ErrQuotaExceeded = fmt.Errorf("quota exceeded")

	// ErrInternal indicates an internal server error.
	ErrInternal = fmt.Errorf("internal server error")

	// ErrServiceUnavailable indicates the service is temporarily unavailable.
	ErrServiceUnavailable = fmt.Errorf("service unavailable")

	// ErrDeadlineExceeded indicates the request deadline was exceeded.
	ErrDeadlineExceeded = fmt.Errorf("deadline exceeded")

	// ErrCanceled indicates the request was canceled.
	ErrCanceled = fmt.Errorf("canceled")

	// ErrInvalidRequest indicates an invalid request.
	ErrInvalidRequest = fmt.Errorf("invalid request")

	// ErrTooManyRequests indicates too many requests.
	ErrTooManyRequests = fmt.Errorf("too many requests")

	// ErrConnection indicates a connection error.
	ErrConnection = fmt.Errorf("connection error")

	// ErrTimeout indicates a timeout error.
	ErrTimeout = fmt.Errorf("timeout")

	// ErrParsing indicates a parsing error.
	ErrParsing = fmt.Errorf("parsing error")

	// ErrValidation indicates a validation error.
	ErrValidation = fmt.Errorf("validation error")

	// ErrFileTooLarge indicates the file is too large.
	ErrFileTooLarge = fmt.Errorf("file too large")

	// ErrUnsupportedContentType indicates unsupported content type.
	ErrUnsupportedContentType = fmt.Errorf("unsupported content type")

	// ErrInvalidModel indicates an invalid model.
	ErrInvalidModel = fmt.Errorf("invalid model")

	// ErrInvalidParameters indicates invalid parameters.
	ErrInvalidParameters = fmt.Errorf("invalid parameters")

	// ErrStreamClosed indicates the stream was closed.
	ErrStreamClosed = fmt.Errorf("stream closed")

	// ErrMalformedData indicates malformed data.
	ErrMalformedData = fmt.Errorf("malformed data")
)

// ErrorType represents different types of errors.
type ErrorType string

const (
	// ErrorTypeConfig represents configuration errors.
	ErrorTypeConfig ErrorType = "config"
	
	// ErrorTypeAuth represents authentication errors.
	ErrorTypeAuth ErrorType = "auth"
	
	// ErrorTypeAPI represents API errors.
	ErrorTypeAPI ErrorType = "api"
	
	// ErrorTypeNetwork represents network errors.
	ErrorTypeNetwork ErrorType = "network"
	
	// ErrorTypeValidation represents validation errors.
	ErrorTypeValidation ErrorType = "validation"
	
	// ErrorTypeRateLimit represents rate limiting errors.
	ErrorTypeRateLimit ErrorType = "rate_limit"
	
	// ErrorTypeQuota represents quota errors.
	ErrorTypeQuota ErrorType = "quota"
	
	// ErrorTypeInternal represents internal errors.
	ErrorTypeInternal ErrorType = "internal"
	
	// ErrorTypeService represents service errors.
	ErrorTypeService ErrorType = "service"
	
	// ErrorTypeTimeout represents timeout errors.
	ErrorTypeTimeout ErrorType = "timeout"
	
	// ErrorTypeCanceled represents cancellation errors.
	ErrorTypeCanceled ErrorType = "canceled"
	
	// ErrorTypeStream represents stream errors.
	ErrorTypeStream ErrorType = "stream"
	
	// ErrorTypeParsing represents parsing errors.
	ErrorTypeParsing ErrorType = "parsing"
	
	// ErrorTypeFile represents file-related errors.
	ErrorTypeFile ErrorType = "file"
)

// Error represents a detailed error with type, code, and context.
type Error struct {
	errorType  ErrorType
	grpcCode   codes.Code
	message    string
	cause      error
	context    map[string]interface{}
	stackTrace []string
}

// NewError creates a new Error with the given type, code, and message.
func NewError(errorType ErrorType, code codes.Code, message string) *Error {
	return &Error{
		errorType:  errorType,
		grpcCode:    code,
		message:    message,
		context:     make(map[string]interface{}),
		stackTrace:  getStackTrace(),
	}
}

// NewErrorWithCause creates a new Error with the given type, code, message, and cause.
func NewErrorWithCause(errorType ErrorType, code codes.Code, message string, cause error) *Error {
	return &Error{
		errorType:  errorType,
		grpcCode:    code,
		message:    message,
		cause:      cause,
		context:     make(map[string]interface{}),
		stackTrace:  getStackTrace(),
	}
}

// NewErrorWithContext creates a new Error with the given type, code, message, and context.
func NewErrorWithContext(errorType ErrorType, code codes.Code, message string, context map[string]interface{}) *Error {
	return &Error{
		errorType:  errorType,
		grpcCode:    code,
		message:    message,
		context:     context,
		stackTrace:  getStackTrace(),
	}
}

// NewErrorFull creates a new Error with all parameters.
func NewErrorFull(errorType ErrorType, code codes.Code, message string, cause error, context map[string]interface{}) *Error {
	return &Error{
		errorType:  errorType,
		grpcCode:    code,
		message:     message,
		cause:       cause,
		context:     context,
		stackTrace:   getStackTrace(),
	}
}

// Error returns the error message.
func (e *Error) Error() string {
	var parts []string
	parts = append(parts, fmt.Sprintf("[%s] %s", e.errorType, e.message))

	if e.cause != nil {
		parts = append(parts, fmt.Sprintf("cause: %v", e.cause))
	}

	if len(e.context) > 0 {
		var ctxParts []string
		for k, v := range e.context {
			ctxParts = append(ctxParts, fmt.Sprintf("%s=%v", k, v))
		}
		parts = append(parts, fmt.Sprintf("context: {%s}", strings.Join(ctxParts, ", ")))
	}

	return strings.Join(parts, " ")
}

// Unwrap returns the cause error.
func (e *Error) Unwrap() error {
	return e.cause
}

// Type returns the error type.
func (e *Error) Type() ErrorType {
	return e.errorType
}

// Code returns the gRPC code.
func (e *Error) Code() codes.Code {
	return e.grpcCode
}

// Context returns the error context.
func (e *Error) Context() map[string]interface{} {
	return e.context
}

// StackTrace returns the stack trace.
func (e *Error) StackTrace() []string {
	return e.stackTrace
}

// GRPCStatus returns the gRPC status.
func (e *Error) GRPCStatus() *status.Status {
	return status.New(e.grpcCode, e.message)
}

// FromGRPC converts a gRPC status error to an xAI SDK error.
func FromGRPC(err error) error {
	if err == nil {
		return nil
	}

	// If it's already an Error type, return as is
	if sdkError, ok := err.(*Error); ok {
		return sdkError
	}

	// Convert gRPC status error
	if st, ok := status.FromError(err); ok {
		errorType := mapGRPCCodeToErrorType(st.Code())
		return NewError(errorType, st.Code(), st.Message())
	}

	// Handle io.EOF as a special case
	if err == io.EOF {
		return io.EOF
	}

	// Handle context errors
	if strings.Contains(err.Error(), "context") {
		if strings.Contains(err.Error(), "canceled") {
			return NewError(ErrorTypeCanceled, codes.Canceled, err.Error())
		}
		if strings.Contains(err.Error(), "deadline exceeded") {
			return NewError(ErrorTypeTimeout, codes.DeadlineExceeded, err.Error())
		}
	}

	// Generic error
	return NewError(ErrorTypeInternal, codes.Unknown, err.Error())
}

// mapGRPCCodeToErrorType maps gRPC codes to xAI SDK error types.
func mapGRPCCodeToErrorType(code codes.Code) ErrorType {
	switch code {
	case codes.InvalidArgument:
		return ErrorTypeValidation
	case codes.Unauthenticated:
		return ErrorTypeAuth
	case codes.PermissionDenied:
		return ErrorTypeAuth
	case codes.NotFound:
		return ErrorTypeAPI
	case codes.AlreadyExists:
		return ErrorTypeAPI
	case codes.ResourceExhausted:
		return ErrorTypeRateLimit
	case codes.FailedPrecondition:
		return ErrorTypeValidation
	case codes.Aborted:
		return ErrorTypeAPI
	case codes.OutOfRange:
		return ErrorTypeValidation
	case codes.Unimplemented:
		return ErrorTypeAPI
	case codes.Internal:
		return ErrorTypeInternal
	case codes.Unavailable:
		return ErrorTypeService
	case codes.DataLoss:
		return ErrorTypeInternal
	case codes.DeadlineExceeded:
		return ErrorTypeTimeout
	case codes.Canceled:
		return ErrorTypeCanceled
	default:
		return ErrorTypeAPI
	}
}

// getStackTrace returns the current stack trace.
func getStackTrace() []string {
	var stack []string
	pcs := make([]uintptr, 10)
	n := runtime.Callers(2, pcs)
	if n == 0 {
		return stack
	}

	frames := runtime.CallersFrames(pcs[:n])
	for {
		frame, more := frames.Next()
		if !more {
			break
		}
		stack = append(stack, fmt.Sprintf("%s:%d %s", frame.File, frame.Line, frame.Function))
	}

	return stack
}

// Helper functions for creating specific error types

// NewConfigError creates a configuration error.
func NewConfigError(message string) *Error {
	return NewError(ErrorTypeConfig, codes.InvalidArgument, message)
}

// NewAuthError creates an authentication error.
func NewAuthError(message string) *Error {
	return NewError(ErrorTypeAuth, codes.Unauthenticated, message)
}

// NewAPIError creates an API error.
func NewAPIError(code codes.Code, message string) *Error {
	return NewError(ErrorTypeAPI, code, message)
}

// NewNetworkError creates a network error.
func NewNetworkError(message string) *Error {
	return NewError(ErrorTypeNetwork, codes.Unavailable, message)
}

// NewValidationError creates a validation error.
func NewValidationError(message string) *Error {
	return NewError(ErrorTypeValidation, codes.InvalidArgument, message)
}