package collections

import "errors"

// ErrNotImplemented is returned when a method is not yet implemented.
var ErrNotImplemented = errors.New("collections API method not yet implemented")

// ErrClientNotInitialized is returned when the REST client is not initialized.
var ErrClientNotInitialized = errors.New("collections client not initialized: REST client is nil")
