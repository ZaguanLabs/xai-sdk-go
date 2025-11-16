package files

import "errors"

// ErrNotImplemented is returned when a method is not yet implemented.
var ErrNotImplemented = errors.New("files API method not yet implemented")

// ErrClientNotInitialized is returned when the REST client is not initialized.
var ErrClientNotInitialized = errors.New("files client not initialized: REST client is nil")

// ErrFileTooLarge is returned when a file exceeds the maximum allowed size.
var ErrFileTooLarge = errors.New("file size exceeds maximum allowed size")
