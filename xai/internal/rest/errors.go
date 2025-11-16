package rest

import "fmt"

// HTTPError represents an HTTP error response.
//
// WARNING: The Error() method includes the full response body, which may contain
// sensitive information. When logging errors, consider using SafeError() instead
// to avoid exposing API keys, tokens, or other sensitive data that might be
// echoed in error responses.
type HTTPError struct {
	StatusCode int
	Body       []byte
}

// Error implements the error interface.
// WARNING: This includes the full response body which may contain sensitive data.
// Use SafeError() for logging purposes.
func (e *HTTPError) Error() string {
	return fmt.Sprintf("HTTP %d: %s", e.StatusCode, string(e.Body))
}

// SafeError returns a sanitized error message without the response body.
// Use this method when logging errors to avoid exposing sensitive information.
func (e *HTTPError) SafeError() string {
	return fmt.Sprintf("HTTP %d error", e.StatusCode)
}

// IsNotFound returns true if the error is a 404 Not Found.
func (e *HTTPError) IsNotFound() bool {
	return e.StatusCode == 404
}

// IsUnauthorized returns true if the error is a 401 Unauthorized.
func (e *HTTPError) IsUnauthorized() bool {
	return e.StatusCode == 401
}

// IsForbidden returns true if the error is a 403 Forbidden.
func (e *HTTPError) IsForbidden() bool {
	return e.StatusCode == 403
}

// IsRateLimited returns true if the error is a 429 Too Many Requests.
func (e *HTTPError) IsRateLimited() bool {
	return e.StatusCode == 429
}

// IsServerError returns true if the error is a 5xx server error.
func (e *HTTPError) IsServerError() bool {
	return e.StatusCode >= 500 && e.StatusCode < 600
}
