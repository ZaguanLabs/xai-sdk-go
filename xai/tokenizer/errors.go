package tokenizer

import "errors"

var ErrClientNotInitialized = errors.New("tokenizer client not initialized: REST client is nil")
