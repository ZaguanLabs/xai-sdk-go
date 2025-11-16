package sample

import "errors"

var ErrClientNotInitialized = errors.New("sample client not initialized: REST client is nil")
