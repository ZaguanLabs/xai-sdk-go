# xAI SDK Go Port Plan

## 1. Goals & Constraints
- Deliver a first-class Go module that mirrors the Python SDK's functionality (sync + async ergonomics, feature coverage, performance), while embracing idiomatic Go APIs.
- Expose a stable module path hosted at `github.com/ZaguanLabs/xai-sdk-go`, ensuring compatibility with Go modules, semantic import versioning, and internal dependency tagging.
- Maintain feature parity with Python surfaces: authentication, chat/completions (including streaming, multi-turn state), files, collections, models, tokenizer, search, image, tools, telemetry, and management APIs (@docs/xai-sdk-python/src/xai_sdk/client.py#46-176).
- Support both secure and optional insecure channels, configurable retries/timeouts, and metadata propagation equivalent to Python defaults.
- Provide robust documentation, tests, CI, and release automation for Go developers.

## 2. Python SDK Inventory (source to port)
1. **Client bootstrap**: `BaseClient` composes shared channel options, metadata, timeout, and credential wiring for both sync (`sync.Client`) and async (`aio.Client`) variants (@docs/xai-sdk-python/src/xai_sdk/client.py#46-176, @docs/xai-sdk-python/src/xai_sdk/sync/client.py#1-108, @docs/xai-sdk-python/src/xai_sdk/aio/client.py#1-128).
2. **Chat module**: Rich request builder supporting conversation state, sampling, streaming, tool calls, search parameters, response formatting, telemetry annotations, **deferred chat** (long-running polling), **stored chat** (conversation branching via `store_messages`/`previous_response_id`), and **encrypted reasoning content** for ZDR users (@docs/xai-sdk-python/src/xai_sdk/chat.py#1-216).
3. **Files module**: Upload/download helper with chunking, progress callbacks, ordering and sorting helpers, plus async counterparts. Also supports **files in chat** via file references (@docs/xai-sdk-python/src/xai_sdk/files.py#1-345).
4. **Search & Tools helpers**: Dataclasses/constructors for search sources (web, news, X, RSS), tool call metadata, and conversions to protobuf messages. Includes **server-side tools** (web_search, x_search, code_execution, collections_search, MCP) (@docs/xai-sdk-python/src/xai_sdk/search.py#1-151, @docs/xai-sdk-python/src/xai_sdk/tools.py#1-217).
5. **Ancillary modules**: `image.py`, `models.py`, `tokenizer.py`, `collections.py` (dual-channel with management API), `auth.py`, telemetry exporters, interceptors (auth + timeout for sync/async), and proto bindings (`src/xai_sdk/proto`).
6. **Examples/tests/docs**: 16 examples per client type (sync/aio) covering: basic chat, streaming, function calling, structured outputs, image generation/understanding, reasoning, search, telemetry, deferred chat, stored chat, files, files_chat, server-side tools, models, tokenizer, auth (@docs/xai-sdk-python/examples/).
7. **Error handling & proto access**: README documents gRPC error codes (UNAUTHENTICATED, PERMISSION_DENIED, RESOURCE_EXHAUSTED, etc.) and `.proto` attribute for accessing raw protobuf objects (@docs/xai-sdk-python/README.md#415-456).

## 3. Target Go Module Architecture
1. **Module skeleton**
   - `go.mod` with minimum Go version (≥1.22 for generics/slices) and dependencies (`google.golang.org/grpc`, `google.golang.org/protobuf`, `github.com/grpc-ecosystem/go-grpc-middleware`, `golang.org/x/sync`, telemetry libs).
   - Folder layout:
     - `/xai` top-level package exposing `Client` (sync) & `AsyncClient` equivalent (using goroutines/channels for async patterns) plus config structs.
     - `/xai/chat`, `/xai/files`, `/xai/image`, `/xai/models`, `/xai/tokenizer`, `/xai/collections`, `/xai/tools`, `/xai/search`, `/xai/telemetry` packages mirroring Python modules.
     - `/xai/internal/grpcutil`, `/xai/internal/auth`, `/xai/internal/stream`, `/xai/internal/testutil` for shared helpers.
     - `/proto` containing generated Go code from `docs/xai-sdk-python/src/xai_sdk/proto/*.proto` (decide whether to vendor `.proto` or fetch upstream).
2. **Code generation**
   - Use `buf` or `protoc` with `protoc-gen-go` and `protoc-gen-go-grpc`. Provide `Makefile`/`mage` targets for regeneration.
   - Automate version tagging between Python and Go packages; embed SDK version string via `internal/version` referencing `git describe` or manual constant.
3. **Configuration API**
   - `type Config struct { APIKey string; ManagementKey string; APIHost string; ManagementHost string; Metadata []metadata.MD; DialOptions []grpc.DialOption; Timeout time.Duration; Insecure bool }` with sensible defaults.
   - Provide environment variable fallbacks for API keys as Python does (`XAI_API_KEY`, `XAI_MANAGEMENT_KEY`).
4. **Client lifecycle**
   - `Client` struct holds `conn *grpc.ClientConn`, `mgmtConn *grpc.ClientConn`, sub-clients, and `Closer` interface.
   - Provide `Close()` and context-aware `NewClient(ctx context.Context, cfg Config)`.
   - Support connection reuse and custom interceptors (auth, timeout, telemetry).

## 4. API Surface Mapping & Implementation Steps
1. **Authentication & Metadata**
   - Implement unary/stream interceptors injecting `authorization: Bearer <key>` and telemetry metadata mirroring Python `_APIAuthPlugin` and metadata additions (`xai-sdk-version`, `xai-sdk-language=go/<version>`).
   - Provide `LocalCredentials` vs `TLS` selection when host starts with `localhost:` similar to `create_channel_credentials`.
2. **Retry/timeout policy**
   - Reproduce `_DEFAULT_CHANNEL_OPTIONS` using gRPC service config JSON or `grpc.WithDefaultServiceConfig` in Go, including keepalive settings and `max_send/receive` message sizes.
   - Expose overrides in Config.
3. **Chat API**
   - Build request builder struct `chat.Request` owning `chatpb.CreateChatRequest` plus helper methods: `Append`, `Sample(ctx)`, `Stream(ctx)` returning channel of deltas, multi-turn state, tool-call decode/encode.
   - Provide typed wrappers for tool choice, response format (JSON schema via `jsonschema` from Go types), reasoning effort, search parameters (mapping to `chatpb.SearchParameters`).
   - Provide helper constructors `chat.System(content string)`, `chat.User(parts ...chat.Part)`, `chat.Image(urlOrData string, detail chat.ImageDetail)`, etc.
   - Provide `chat.Parse[T any](ctx)` to unmarshal JSON-object responses into Go structs akin to Python `BaseModel` support.
   - **Deferred chat**: Implement polling-based `SampleDeferred(ctx)` for long-running requests returning job ID + status checker.
   - **Stored chat**: Support `StoreMessages bool` and `PreviousResponseID string` for conversation branching and retrieval via `GetStoredCompletion`/`DeleteStoredCompletion`.
   - **Encrypted content**: Support `UseEncryptedContent bool` for ZDR-enabled users to maintain reasoning trace continuity without server-side storage.
4. **Search helpers**
   - Offer builder functions `search.WebSource(opts ...)`, `search.NewsSource`, `search.XSource`, etc., returning proto structs and validations replicating Python enforcement.
   - Implement `SearchParameters.ToProto()` and validation (mode must be `auto|on|off`).
5. **Tools helpers**
   - Mirror `tools.web_search`, `tools.x_search`, `code_execution`, `collections_search`, `mcp`, including datetime handling and domain/x-handle constraints.
6. **Files API**
   - Provide streaming upload/download using chunk size 3 MiB with progress callbacks. Support sync + async (async via goroutines returning `chan Progress` and `context.Context`).
   - File iteration should reuse `io.Reader` chunker; expose `UploadBytes`, `UploadFile(path)`, `UploadReader`. Provide `Download` returning streaming reader.
   - **Files in chat**: Support file references in chat messages via file IDs, enabling document-based conversations (example: `files_chat.py`).
7. **Collections**
   - Provide dual-channel support with API + optional management connection. Validate presence of management key before management RPCs (matching Python guard in `collections.Client`).
8. **Image / Models / Tokenizer**
   - Straightforward wrappers around proto stubs with typed helper structs for requests/responses.
9. **Telemetry**
   - Provide optional OpenTelemetry instrumentation: `telemetry.SetupConsoleExporter`, `telemetry.SetupOTLPExporter`, and trace spans per call with attributes (user prompts, response metadata) matching Python flavor described in README.
   - Provide env toggles `XAI_SDK_DISABLE_TRACING`, `XAI_SDK_DISABLE_SENSITIVE_TELEMETRY_ATTRIBUTES`.
10. **Streaming ergonomics**
    - Use Go channels to surface streaming deltas; return aggregated `Response` plus `<-chan chat.Chunk` similar to Python `for response, chunk in chat.stream()` semantics. Provide helper to iterate, plus context cancellation support.
11. **Async patterns**
    - Since Go lacks async/await, provide goroutine utilities: `chat.Stream(ctx)` returns `StreamIterator` with `Next() (Chunk, error)`; multi-turn state maintained inside request builder.
12. **Error handling**
    - Convert `status.Status` to typed errors; wrap with sentinel errors matching Python/README documentation: `ErrUnauthenticated`, `ErrPermissionDenied`, `ErrResourceExhausted`, `ErrDeadlineExceeded`, `ErrNotFound`, `ErrInvalidArgument`, `ErrInternal`, `ErrUnavailable`, `ErrDataLoss`, `ErrUnknown`.
    - Provide helper `errors.FromGRPC(err)` returning typed error with `.Code()`, `.Message()`, and `.Details()` accessors.
    - Support accessing raw proto objects via `.Proto()` method on response types for advanced use cases (document when/why to use).

## 5. Proto Source Strategy
1. Copy `.proto` files from `docs/xai-sdk-python/src/xai_sdk/proto` into `/proto` folder while keeping them authoritative; document update workflow.
2. Create `buf.yaml`/`buf.gen.yaml` for lint/build.
3. Generate Go code under `proto/xai/...` with go_package options for stable import paths.
4. Provide CI guard verifying generated code matches `.proto`.

## 6. Testing & QA
1. **Unit tests**: Table-driven tests per package (chat builders, search/tool conversions, config validation, metadata injection, file chunkers).
2. **Integration tests**: Use `grpc-go` test server or recorded fixtures to validate RPC wiring. Provide gating for live tests using env var `XAI_SDK_E2E=1` and real API keys (skipped in CI otherwise).
3. **Concurrency tests**: Ensure `Client` methods are goroutine-safe; add race-detector coverage in CI.
4. **Proto compatibility tests**: Validate generated Go structs match Python examples by marshaling/unmarshaling sample payloads from `docs/xai-sdk-python/examples`.
5. **Telemetry tests**: Use OTLP test collectors to assert exported spans/attributes when enabled/disabled.

## 7. Tooling, CI, and Release
1. **Build tooling**: Provide `make fmt`, `make lint`, `make test`, `make proto`, `make examples`. Use `golangci-lint` with modules (`govet`, `staticcheck`, `ineffassign`, `gofmt`).
2. **CI pipeline** (GitHub Actions): matrix over Go versions (1.22, 1.21) running unit tests + race + lint + `buf lint` + `buf breaking`.
3. **Release automation**: Tag-based GitHub workflow building modules, running `goreleaser` to publish docs/examples and attach checksums.
4. **Versioning**: Track parity with Python releases via changelog entry referencing upstream tag; maintain `CHANGELOG.md` and `docs/` updates.

## 8. Documentation & Examples
1. Port README content to Go-specific instructions (installation via `go get`, usage snippets covering chat, streaming, tools, telemetry) referencing equivalent Python scenarios. Include sections on error codes, versioning (SemVer), accessing proto objects, and determining installed version.
2. Re-create `examples/` tree mirroring Python structure with 16 examples: `auth`, `chat`, `deferred_chat`, `files`, `files_chat`, `function_calling`, `image_generation`, `image_understanding`, `models`, `reasoning`, `search`, `server_side_tools`, `stored_chat`, `structured_outputs`, `telemetry`, `tokenizer`. Provide both blocking and concurrent variants where applicable.
3. Document environment variables (`XAI_API_KEY`, `XAI_MANAGEMENT_KEY`, `XAI_SDK_DISABLE_TRACING`, `XAI_SDK_DISABLE_SENSITIVE_TELEMETRY_ATTRIBUTES`, OTEL vars), telemetry configuration, retries/timeouts, file upload best practices, and error handling patterns.
4. Provide migration guide for developers moving from Python to Go, highlighting API differences (async/await → goroutines/channels, context propagation, error handling), language idioms, and feature mapping table.

## 9. Migration & Delivery Phases
1. **Phase 1 – Foundations**: Setup repo, module scaffolding, proto generation, base client with auth interceptors, config, connection mgmt, simple chat completion + streaming.
2. **Phase 2 – Core Features**: Port chat (multi-turn, tools, search, response formats, Parse), files (upload/download/progress), models, tokenizer, image generation/understanding.
3. **Phase 3 – Advanced Chat**: Implement deferred chat (polling), stored chat (branching), encrypted content, files_chat integration.
4. **Phase 4 – Tools & Collections**: Port server-side tools (web_search, x_search, code_execution, collections_search, MCP), collections API with dual-channel management.
5. **Phase 5 – Telemetry & Observability**: Implement OTEL exporters (console/OTLP), span attributes, env toggles, context propagation, documentation.
6. **Phase 6 – QA & Hardening**: Add integration tests, race detector, load tests, error handling coverage, finalize all 16 examples per variant, run beta with internal consumers.
7. **Phase 7 – GA Release**: Tag v0.1.0 following SemVer, publish docs, announce, set up issue templates, establish changelog workflow.

## 10. Risks & Mitigations
- **Proto drift**: Upstream `.proto` updates could desync. Mitigate with `buf breaking` checks against upstream main.
- **Async semantics mismatch**: Python exposes `AsyncClient`; Go lacks direct analog. Provide idiomatic concurrency helpers and document differences.
- **Telemetry privacy expectations**: Ensure env toggles and safe defaults mimic Python behavior (no sensitive traces unless enabled).
- **Large file uploads**: Manage memory usage with streaming readers/writers and configurable chunk size.
- **API stability**: Use interfaces and internal modules to allow future addition of REST fallback if gRPC shape changes.

## 11. Next Actions
1. Bootstrap repo on GitHub (`github.com/ZaguanLabs/xai-sdk-go`) with Go module structure, buf config, Makefile, LICENSE (Apache-2.0 matching Python), and GitHub workflow scaffolding.
2. Import proto definitions from `docs/xai-sdk-python/src/xai_sdk/proto/*.proto` and generate Go bindings, committing both `.proto` sources and generated code with Buf linting in CI.
3. Implement base client + auth/timeout interceptors (sync + async variants) + config, validate defaults against Python behavior (channel options, keepalive, retry policy, metadata injection) using sample requests and embed SDK metadata (`xai-sdk-version: go/<version>`, `xai-sdk-language: go/<runtime_version>`).
4. Create a parity tracking doc (`docs/feature-parity.md`) with checklist covering all 7 inventory categories, 16 example scenarios, and link from README to communicate progress transparently.
5. Begin porting chat functionality (request builder, Append, Sample, Stream, Parse, message helpers) as first major surface with unit tests, then proceed module-by-module following 7-phase delivery plan.
6. Define CONTRIBUTING.md, CODE_OF_CONDUCT.md, issue/PR templates, and GitHub Projects board early so external collaborators can start filing/working on tasks before GA.

## 12. Repository Workflow & Governance
- **Primary remote**: https://github.com/ZaguanLabs/xai-sdk-go. Default branch `main`, protected with required checks (lint, test, buf) before merge.
- **Branching model**: feature branches prefixed by area (`feature/chat-stream`, `fix/files-upload`), PRs referencing tracking issues labeled by workstream (e.g., `chat`, `files`, `telemetry`).
- **Issue tracking**: Use GitHub Projects board with columns (`Backlog`, `In progress`, `Review`, `Done`). Each task links to sections of this plan for context.
- **Release process**: Tag releases (`v0.1.0`, `v0.2.0` …) on `main`. GitHub Action runs `goreleaser` to publish artifacts and updates `docs/CHANGELOG.md`. Mirror releases to Go proxy by ensuring tags follow semantic import versioning.
- **Documentation updates**: Every feature PR must update relevant docs/examples plus parity tracker. Add `docs/adr/` (architecture decision records) for foundational decisions (e.g., auth interceptor strategy, telemetry defaults).
- **Security & compliance**: Enforce Dependabot updates, use CodeQL workflow, and require signed commits for maintainers. Sensitive config (API keys) pulled from GitHub secrets for CI integration tests gated behind opt-in workflow.
