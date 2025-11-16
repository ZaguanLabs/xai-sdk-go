package embed

import "errors"

// ErrClientNotInitialized is returned when the REST client is not initialized.
var ErrClientNotInitialized = errors.New("embed client not initialized: REST client is nil")
