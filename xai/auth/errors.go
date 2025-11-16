package auth

import "errors"

// ErrNotImplemented is returned when a method is not yet implemented.
var ErrNotImplemented = errors.New("auth API method not yet implemented")

// ErrClientNotInitialized is returned when the REST client is not initialized.
var ErrClientNotInitialized = errors.New("auth client not initialized: REST client is nil")
