# xAI SDK Go - Development Plan

This document breaks down the porting effort into discrete, committable phases. Each phase represents a logical unit of work that can be committed to git independently.

---

## Phase 0: Repository Bootstrap
**Goal**: Initialize repository structure, tooling, and foundational files.

### Tasks
- Initialize `go.mod` with module path `github.com/ZaguanLabs/xai-sdk-go` (Go 1.22+)
- Create directory structure: `/xai`, `/xai/internal`, `/proto`, `/examples`, `/docs`, `/.github/workflows`
- Add `LICENSE` (Apache-2.0), `README.md`, `.gitignore`, `CONTRIBUTING.md`, `CODE_OF_CONDUCT.md`, `CHANGELOG.md`
- Setup `Makefile` with targets: `help`, `fmt`, `lint`, `test`, `proto`, `clean`
- Configure Buf: `buf.yaml`, `buf.gen.yaml`, `buf.lock`
- Setup CI skeleton: `.github/workflows/ci.yml` with Go matrix (1.22, 1.23)

### Commit
```
chore: initialize repository structure and tooling
```

---

## Phase 1: Proto Definitions & Code Generation
**Goal**: Import proto files and generate Go bindings.

### Tasks
- Copy `.proto` files from `docs/xai-sdk-python/src/xai_sdk/proto/` to `/proto/`
- Configure `buf.gen.yaml` with Go plugin, set `go_package` options
- Generate Go code under `/proto/gen/go/`
- Update `Makefile` with `make proto` target
- Add proto dependencies to `go.mod`
- Document proto workflow in `proto/README.md`

### Commit
```
feat(proto): import proto definitions and generate Go bindings
```

---

## Phase 2: Internal Utilities
**Goal**: Create internal packages for shared utilities.

### Tasks
- Create `xai/internal/version/version.go` with SDK version and Go runtime detection
- Create `xai/internal/constants/constants.go` with defaults (hosts, timeout, sizes)
- Create `xai/internal/errors/errors.go` with sentinel errors and `FromGRPC()` converter
- Create `xai/internal/metadata/metadata.go` with SDK metadata helpers
- Add comprehensive unit tests

### Commit
```
feat(internal): add internal utilities for version, errors, and metadata
```

---

## Phase 3: Configuration & Client Skeleton
**Goal**: Implement client configuration and basic structure.

### Tasks
- Create `xai/config.go` with `Config` type, defaults, env var fallbacks, validation
- Create `xai/client.go` with `Client` struct, `NewClient()`, `Close()`
- Implement channel creation with TLS/local credentials
- Create `xai/internal/grpcutil/serviceconfig.go` with retry policy and keepalive
- Add tests for config validation and client initialization

### Commit
```
feat(client): implement configuration and client skeleton
```

---

## Phase 4: Authentication Interceptors
**Goal**: Implement gRPC interceptors for auth and timeout.

### Tasks
- Create `xai/internal/auth/interceptor.go` with unary/stream auth interceptors
- Create `xai/internal/grpcutil/timeout.go` with timeout interceptors
- Integrate interceptors into `Client.NewClient()`
- Support both secure and insecure channels
- Add tests with mock gRPC server and benchmarks

### Commit
```
feat(auth): implement authentication and timeout interceptors
```

---

## Phase 5: Chat Foundation
**Goal**: Implement basic chat with message builders.

### Tasks
- Create `xai/chat/chat.go` with `Request`, `Response`, `Chunk` types
- Implement message builders: `System()`, `User()`, `Assistant()` in `xai/chat/message.go`
- Create content helpers: `Text()`, `Image()` in `xai/chat/content.go`
- Implement `Request` builder with functional options
- Add chat client stub attached to main `Client`
- Add tests and `examples/chat/basic.go` skeleton

### Commit
```
feat(chat): implement chat foundation with message builders
```

---

## Phase 6: Chat Sample & Streaming
**Goal**: Implement synchronous sampling and streaming.

### Tasks
- Implement `(r *Request) Sample(ctx) (*Response, error)`
- Implement `(r *Request) Stream(ctx) (*StreamIterator, error)`
- Create `StreamIterator` with `Next()` method in `xai/chat/stream.go`
- Add response helpers: `Content()`, `ToolCalls()`, `Proto()`
- Add tests with mock server
- Complete `examples/chat/basic.go` and add `examples/chat/streaming.go`

### Commit
```
feat(chat): implement Sample and Stream methods
```

---

## Phase 7: Chat Tool Support
**Goal**: Add function calling and tool choice.

### Tasks
- Create `xai/chat/tool.go` with `Tool`, `Function`, `ToolCall`, `ToolResult` types
- Implement `Function()` builder with JSON schema generation
- Add `ToolChoice` type and `RequiredTool()` helper
- Implement tool call parsing in responses
- Add tests and `examples/chat/function_calling.go`

### Commit
```
feat(chat): add function calling and tool support
```

---

## Phase 8: Chat Advanced Features
**Goal**: Implement Parse, response formats, search, reasoning.

### Tasks
- Implement `(r *Request) Parse(ctx, v any) error` with generics in `xai/chat/parse.go`
- Add `ResponseFormat` type (text, json_object, json_schema)
- Create `xai/search/search.go` with `SearchParameters`, source builders
- Add `ReasoningEffort` type (low, high)
- Add all sampling parameters (temperature, top_p, etc.)
- Add tests and examples: `structured_outputs.go`, `reasoning.go`, `search.go`

### Commit
```
feat(chat): add Parse, response formats, search, and reasoning
```

---

## Phase 9: Chat Deferred & Stored
**Goal**: Implement deferred and stored chat.

### Tasks
- Implement `SampleDeferred()` with polling in `xai/chat/deferred.go`
- Implement stored chat: `WithStoreMessages()`, `WithPreviousResponseID()`, `GetStoredCompletion()`, `DeleteStoredCompletion()`
- Add encrypted content support: `WithEncryptedContent()`
- Add tests and examples: `deferred_chat.go`, `stored_chat.go`

### Commit
```
feat(chat): add deferred and stored chat support
```

---

## Phase 10: Files Module
**Goal**: Implement file upload/download with chunking.

### Tasks
- Create `xai/files/files.go` with `Client`, `File`, `ProgressCallback` types
- Implement upload methods: `Upload()`, `UploadBytes()`, `UploadReader()` with 3 MiB chunking
- Implement `Download()` with streaming
- Implement `List()`, `Get()`, `Delete()` with sorting/ordering
- Add tests and examples: `upload.go`, `download.go`

### Commit
```
feat(files): implement file upload/download with chunking
```

---

## Phase 11: Files in Chat
**Goal**: Integrate file references into chat.

### Tasks
- Add `File(fileID string) Part` to `xai/chat/content.go`
- Update message builders to support file references
- Add tests and `examples/chat/files_chat.go`

### Commit
```
feat(chat): add file reference support in messages
```

---

## Phase 12: Server-Side Tools
**Goal**: Implement server-side tool builders.

### Tasks
- Create `xai/tools/tools.go` with tool builder functions
- Implement `WebSearch()`, `XSearch()`, `CodeExecution()`, `CollectionsSearch()`, `MCP()`
- Add validation (domain limits, collection limits)
- Add `GetToolCallType()` helper
- Add tests and `examples/chat/server_side_tools.go`

### Commit
```
feat(tools): implement server-side tool builders
```

---

## Phase 13: Image Module
**Goal**: Implement image generation API.

### Tasks
- Create `xai/image/image.go` with `Client`, `Image`, `GenerateRequest` types
- Implement `Generate()` method with prompt, model, size, quality, style
- Add image helpers: `URL()`, `Data()`, `Save()`
- Add tests and `examples/image/generation.go`

### Commit
```
feat(image): implement image generation API
```

---

## Phase 14: Models, Tokenizer, Auth
**Goal**: Implement remaining simple modules.

### Tasks
- Create `xai/models/models.go` with `List()`, `Get()`
- Create `xai/tokenizer/tokenizer.go` with `Encode()`, `Decode()`, `Count()`
- Create `xai/auth/auth.go` with `Validate()`
- Integrate with main `Client`
- Add tests and examples: `models/list.go`, `tokenizer/encode.go`, `auth/validate.go`

### Commit
```
feat: implement models, tokenizer, and auth modules
```

---

## Phase 15: Collections Module
**Goal**: Implement collections API with dual-channel.

### Tasks
- Create `xai/collections/collections.go` with dual-channel support
- Implement collection operations: `ListCollections()`, `GetCollection()`, `CreateCollection()`, `DeleteCollection()`
- Implement document operations: `ListDocuments()`, `GetDocument()`, `AddDocument()`, `DeleteDocument()`
- Add sorting/ordering helpers
- Require management key validation
- Add tests and `examples/collections/manage.go`

### Commit
```
feat(collections): implement collections API with dual-channel support
```

---

## Phase 16: Telemetry Foundation
**Goal**: Implement OpenTelemetry foundation.

### Tasks
- Create `xai/telemetry/telemetry.go` with `Telemetry` struct
- Implement console exporter: `SetupConsoleExporter()`
- Add tracer creation with lazy initialization
- Support environment variables: `XAI_SDK_DISABLE_TRACING`, `XAI_SDK_DISABLE_SENSITIVE_TELEMETRY_ATTRIBUTES`
- Add tests and documentation

### Commit
```
feat(telemetry): implement OpenTelemetry foundation and console exporter
```

---

## Phase 17: Telemetry Instrumentation
**Goal**: Add OTLP exporter and instrument all operations.

### Tasks
- Implement OTLP exporter: `SetupOTLPExporter()` with HTTP/gRPC support
- Create span helpers in `xai/internal/telemetry/span.go`
- Instrument all operations: chat, files, image, models, tokenizer, auth, collections
- Follow OpenTelemetry GenAI semantic conventions
- Add sensitive attribute filtering
- Add tests and examples: `telemetry/console.go`, `telemetry/otlp.go`

### Commit
```
feat(telemetry): add OTLP exporter and instrument all operations
```

---

## Phase 18: Error Handling Enhancement
**Goal**: Enhance error handling with detailed types.

### Tasks
- Expand `xai/errors/errors.go` with detailed error structs for each gRPC code
- Add error helpers: `IsUnauthenticated()`, `IsPermissionDenied()`, etc.
- Add retry helpers: `IsRetryable()`, `ShouldRetry()`
- Update all packages to use typed errors consistently
- Add tests and error handling guide

### Commit
```
feat(errors): enhance error handling with detailed types and helpers
```

---

## Phase 19: Examples Completion
**Goal**: Complete all 16 examples.

### Tasks
- Verify and enhance existing examples
- Add missing examples to reach 16 total
- Ensure all examples compile, run, and include proper error handling
- Add comprehensive comments and documentation
- Create `examples/README.md` with index

### Commit
```
docs(examples): complete all 16 examples with documentation
```

---

## Phase 20: Documentation
**Goal**: Complete comprehensive documentation.

### Tasks
- Write main `README.md` with installation, quickstart, features
- Document all environment variables
- Create migration guide from Python to Go
- Add godoc comments to all exported types/functions
- Create `docs/feature-parity.md` tracking checklist
- Add architecture decision records in `docs/adr/`
- Document error codes and handling patterns

### Commit
```
docs: add comprehensive documentation and migration guide
```

---

## Phase 21: Testing & CI Enhancement
**Goal**: Achieve comprehensive test coverage and robust CI.

### Tasks
- Add integration tests with gating via `XAI_SDK_E2E=1`
- Add concurrency tests with race detector
- Add proto compatibility tests
- Enhance CI with coverage reporting, race detection, multiple Go versions
- Add `buf breaking` checks
- Setup Dependabot and CodeQL

### Commit
```
test: add comprehensive test coverage and enhance CI
```

---

## Phase 22: Release Preparation
**Goal**: Prepare for initial release.

### Tasks
- Setup `goreleaser` configuration
- Create release workflow in `.github/workflows/release.yml`
- Write `CHANGELOG.md` for v0.1.0
- Add GitHub issue templates and PR template
- Setup GitHub Projects board
- Tag v0.1.0-rc.1 for release candidate testing

### Commit
```
chore: prepare for v0.1.0 release
```

---

## Summary

**Total Phases**: 22 (Phase 0-21 + Phase 22)

**Estimated Commits**: 22 major commits (one per phase)

**Key Milestones**:
- Phase 0-4: Foundation (repo, proto, client, auth)
- Phase 5-9: Core chat functionality
- Phase 10-15: Additional modules (files, tools, image, models, collections)
- Phase 16-17: Telemetry and observability
- Phase 18-21: Polish (errors, examples, docs, tests)
- Phase 22: Release preparation

Each phase is designed to be independently committable with working, tested code.
