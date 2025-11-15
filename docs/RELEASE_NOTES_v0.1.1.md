# Release v0.1.1 - Bug Fix Release

This is a critical bug fix release that resolves issues with the models API.

## 🐛 Bug Fixes

### Models API Proto Definitions
- **Fixed proto package name**: Changed from `xai.v1` to `xai_api` to match the actual xAI API
- **Fixed RPC methods**: Replaced generic `ListModels` with type-specific methods:
  - `ListLanguageModels()`
  - `ListEmbeddingModels()`
  - `ListImageGenerationModels()`
- **Fixed proto field numbers**: Corrected field ordering in all model message types to match the server's wire format
  - `LanguageModel`: Moved `aliases` from field 2 → 11, `cached_prompt_token_price` from 8 → 12, `search_price` from 10 → 13
  - `EmbeddingModel`: Moved `aliases` from field 2 → 11
  - `ImageGenerationModel`: Moved `aliases` from field 2 → 11

### Metadata Handling
- **Fixed metadata operations**: Changed from `metadata.NewOutgoingContext()` to `metadata.AppendToOutgoingContext()` to preserve gRPC internal headers
- **Removed content-type interceptor**: Eliminated manual content-type header that was interfering with gRPC's automatic handling

## ✨ Improvements

- **Enhanced models example**: Added comprehensive debug logging for better troubleshooting
- **Updated models client API**: Now uses type-specific methods instead of generic `List()` and `Get()`
- **Better error messages**: Improved error handling and reporting

## 🎯 What's Working Now

The SDK now successfully:
- ✅ Lists all available language models (8 models including grok-2, grok-3, grok-4 variants)
- ✅ Retrieves detailed information about specific models
- ✅ Handles authentication and metadata correctly
- ✅ Properly unmarshals proto responses from the xAI API

## 📦 Installation

```bash
go get github.com/ZaguanLabs/xai-sdk-go@v0.1.1
```

## 🚀 Quick Example

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/ZaguanLabs/xai-sdk-go/xai"
)

func main() {
    client, err := xai.NewClientWithAPIKey("your-api-key")
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()
    
    modelsClient := client.Models()
    ctx := client.NewContext(context.Background())
    
    models, err := modelsClient.ListLanguageModels(ctx)
    if err != nil {
        log.Fatal(err)
    }
    
    for _, model := range models {
        fmt.Printf("Model: %s (v%s)\n", model.Name(), model.Version())
    }
}
```

## 📝 Full Changelog

See [CHANGELOG.md](https://github.com/ZaguanLabs/xai-sdk-go/blob/main/CHANGELOG.md) for complete details.

## 🙏 Acknowledgments

Thanks to everyone who reported issues and helped test the fixes!

---

**Full Changelog**: https://github.com/ZaguanLabs/xai-sdk-go/compare/v0.1.0...v0.1.1
