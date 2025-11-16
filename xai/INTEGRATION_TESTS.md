# Integration Tests

This directory contains integration tests for the xAI SDK REST APIs. These tests make actual API calls to the xAI service and require valid credentials.

## Running Integration Tests

### Prerequisites

1. **API Key**: You need a valid xAI API key
2. **Environment Variable**: Set `XAI_API_KEY` environment variable

```bash
export XAI_API_KEY=your-api-key-here
```

### Run All Integration Tests

```bash
make test-integration
```

Or manually:

```bash
go test -tags=integration -v ./xai/embed ./xai/files ./xai/image ./xai/auth
```

### Run Specific API Tests

```bash
# Embeddings API
go test -tags=integration -v ./xai/embed

# Files API
go test -tags=integration -v ./xai/files

# Image Generation API
go test -tags=integration -v ./xai/image

# Auth API
go test -tags=integration -v ./xai/auth
```

## Test Coverage

### Embed API (`xai/embed/embed_integration_test.go`)
- ✅ Generate text embeddings
- ✅ Generate batch embeddings
- ✅ Generate image embeddings (if supported)

### Files API (`xai/files/files_integration_test.go`)
- ✅ Upload files
- ✅ List files
- ✅ Get file metadata
- ✅ Get download URL
- ✅ Download file content
- ✅ Delete files

### Image API (`xai/image/image_integration_test.go`)
- ✅ Generate single image
- ✅ Generate multiple images
- ✅ Generate Base64 format

### Auth API (`xai/auth/auth_integration_test.go`)
- ✅ Validate API key
- ✅ List API keys
- ✅ Get key by ID

## Build Tags

Integration tests use the `integration` build tag to prevent them from running during normal `go test` or CI builds. This ensures:

1. **No accidental API calls** during regular testing
2. **No API key required** for unit tests
3. **Explicit opt-in** for integration testing

## CI/CD

Integration tests are **not** run in CI by default. To run them in CI:

1. Set `XAI_API_KEY` as a secret
2. Add a separate CI job with the integration tag:

```yaml
- name: Integration Tests
  env:
    XAI_API_KEY: ${{ secrets.XAI_API_KEY }}
  run: make test-integration
```

## Notes

- **Cleanup**: Tests clean up resources they create (e.g., uploaded files)
- **Costs**: Some tests may incur API usage costs
- **Rate Limits**: Be aware of API rate limits when running tests
- **Skipping**: Tests automatically skip if `XAI_API_KEY` is not set
- **Failures**: Some endpoints may not be fully implemented yet and will skip gracefully

## Adding New Integration Tests

1. Create a new test file with `// +build integration` at the top
2. Use `_integration_test.go` suffix
3. Check for `XAI_API_KEY` and skip if not present
4. Clean up any resources created during the test
5. Add the package to `Makefile` test-integration target

Example:

```go
// +build integration

package myapi_test

import (
    "context"
    "os"
    "testing"
    
    "github.com/ZaguanLabs/xai-sdk-go/xai"
)

func TestMyAPIIntegration(t *testing.T) {
    apiKey := os.Getenv("XAI_API_KEY")
    if apiKey == "" {
        t.Skip("XAI_API_KEY not set, skipping integration test")
    }
    
    client, err := xai.NewClient(&xai.Config{APIKey: apiKey})
    if err != nil {
        t.Fatalf("Failed to create client: %v", err)
    }
    defer client.Close()
    
    // Your test code here
}
```
