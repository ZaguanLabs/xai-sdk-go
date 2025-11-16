package deferred

import "errors"

// ErrClientNotInitialized is returned when the REST client is not initialized.
var ErrClientNotInitialized = errors.New("deferred client not initialized: REST client is nil")
