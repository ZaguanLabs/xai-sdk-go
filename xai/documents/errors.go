package documents

import "errors"

// ErrClientNotInitialized is returned when the REST client is not initialized.
var ErrClientNotInitialized = errors.New("documents client not initialized: REST client is nil")
