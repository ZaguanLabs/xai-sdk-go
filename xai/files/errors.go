package files

import "errors"

// ErrNotImplemented is returned when a method is not yet implemented.
// The Files API in the official xAI Python SDK is REST-based and does not
// have a gRPC service definition. These methods are placeholders for when
// gRPC support is added or REST client is implemented.
var ErrNotImplemented = errors.New("files API not yet implemented: waiting for gRPC service definition or REST client")
