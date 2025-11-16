package image

import "errors"

// ErrClientNotInitialized is returned when the REST client is not initialized.
var ErrClientNotInitialized = errors.New("image client not initialized: REST client is nil")
