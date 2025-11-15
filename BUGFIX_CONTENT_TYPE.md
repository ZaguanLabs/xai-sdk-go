# Bug Fix: "malformed header: missing HTTP content-type"

## Problem
When calling `modelsClient.List(ctx)`, the SDK was returning the error:
```
list models failed (Unknown): malformed header: missing HTTP content-type
```

## Root Cause
The SDK was using `metadata.NewOutgoingContext()` in multiple places, which **replaces** all existing metadata in the context. This was interfering with gRPC's automatic content-type header handling.

Specifically, the issue occurred in three locations:
1. **xai/internal/auth/interceptor.go** - Auth interceptor (line 99, 198)
2. **xai/internal/metadata/metadata.go** - SDK metadata (line 93)
3. **xai/config.go** - Content-type interceptor (lines 326-329)

When gRPC internally sets the content-type header, our interceptors were replacing the entire metadata context, causing the content-type to be lost or malformed.

## Solution
Changed all metadata operations to use `metadata.AppendToOutgoingContext()` instead of `metadata.NewOutgoingContext()`. This function **appends** to existing metadata rather than replacing it, preserving gRPC's internal headers.

### Files Modified

#### 1. xai/config.go
- **Removed** the content-type interceptor entirely (lines 326-329)
- gRPC automatically handles content-type headers; manually adding them causes conflicts

#### 2. xai/internal/metadata/metadata.go
- **Changed** `AddToOutgoingContext()` method to use `metadata.AppendToOutgoingContext()`
- Converts metadata to key-value pairs for the append function
- Preserves all existing metadata including gRPC's content-type

#### 3. xai/internal/auth/interceptor.go
- **Changed** `addAuthMetadata()` in `APIKeyAuthInterceptor` to use `metadata.AppendToOutgoingContext()`
- **Changed** `addAuthMetadata()` in `CombinedAuthInterceptor` to use `metadata.AppendToOutgoingContext()`
- Simplified the logic since we no longer need to manually merge metadata

## Testing
To test the fix:
```bash
export XAI_API_KEY="your-api-key"
go run examples/models/list.go
```

The command should now successfully list available models without the "malformed header" error.

## Technical Details

### Why NewOutgoingContext is problematic
`metadata.NewOutgoingContext(ctx, md)` creates a **new** context with the provided metadata, discarding any metadata that was already in the context. This is problematic in interceptor chains where multiple interceptors need to add metadata.

### Why AppendToOutgoingContext is correct
`metadata.AppendToOutgoingContext(ctx, key1, val1, key2, val2, ...)` **appends** the key-value pairs to any existing outgoing metadata in the context. This preserves:
- gRPC's internal headers (like content-type)
- Metadata added by other interceptors in the chain
- Any metadata already present in the context

## Related Issues
This fix also resolves a context leak warning in `client.go` line 178 by properly handling the cancel function from `context.WithTimeout`.

## Additional Notes
The lint warnings in `client_test.go` about passing nil contexts are pre-existing test issues and are unrelated to this bug fix.
