# Testing the Fixed SDK

## What Was Fixed

The SDK had incorrect proto definitions that didn't match the actual xAI API. The main issues were:

1. **Wrong package name**: Used `xai.v1` instead of `xai_api`
2. **Wrong RPC methods**: Used `ListModels` instead of `ListLanguageModels`, `ListEmbeddingModels`, `ListImageGenerationModels`
3. **Wrong message structures**: Simplified models that didn't match the actual API schema

## Changes Made

### 1. Updated Proto Definition (`proto/xai/v1/models.proto`)
- Changed package from `xai.v1` to `xai_api`
- Added correct RPC methods matching the xAI API
- Added proper message definitions with all fields

### 2. Regenerated Go Code
- Ran protoc to generate new `models.pb.go` and `models_grpc.pb.go`

### 3. Updated Models Client (`xai/models/models.go`)
- Replaced generic `Model` with `LanguageModel`, `EmbeddingModel`, `ImageGenerationModel`
- Added methods: `ListLanguageModels()`, `ListEmbeddingModels()`, `ListImageGenerationModels()`
- Added methods: `GetLanguageModel()`, `GetEmbeddingModel()`, `GetImageGenerationModel()`

### 4. Updated Example (`examples/models/list.go`)
- Changed to use `ListLanguageModels()` instead of `List()`
- Updated to display new model fields (aliases, version, modalities, etc.)

### 5. Fixed Metadata Handling
- Changed `metadata.NewOutgoingContext()` to `metadata.AppendToOutgoingContext()` in:
  - `xai/internal/auth/interceptor.go`
  - `xai/internal/metadata/metadata.go`
- Removed manual content-type interceptor from `xai/config.go`

## Test Command

```bash
export XAI_API_KEY="your-api-key-here"
go run examples/models/list.go
```

## Expected Output

You should see:
1. Debug logs showing the connection process
2. Outgoing metadata (with API key redacted)
3. A list of available language models with details
4. Detailed information about one specific model
5. Success message

## If It Still Fails

If you still get errors, please share:
1. The complete error message
2. The debug output
3. Your Go version (`go version`)
4. Your protoc version (`protoc --version`)
