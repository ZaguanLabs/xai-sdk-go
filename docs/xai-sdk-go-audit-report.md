## 1. Security Audit

### 1.1 Authentication & Authorization

- **API key handling (REST & gRPC)**  
  - **Findings**  
    - REST: [internal/rest.Client](cci:2://file:///home/stig/dev/ai/zaguan/labs/xai-sdk-go/xai/client.go:28:0-39:1) only uses API key in `Authorization: Bearer <key>` header; no logging in library code.  
    - gRPC: [internal/auth](cci:7://file:///home/stig/dev/ai/zaguan/labs/xai-sdk-go/xai/internal/auth:0:0-0:0) interceptors attach API key via metadata (`authorization: Bearer`, `x-api-key`) using `AppendToOutgoingContext`; no logging.  
    - [internal/metadata.SanitizeMetadata](cci:1://file:///home/stig/dev/ai/zaguan/labs/xai-sdk-go/xai/internal/metadata/metadata.go:260:0-275:1) masks API key before logging; [Config.String()](cci:1://file:///home/stig/dev/ai/zaguan/labs/xai-sdk-go/xai/config.go:432:0-445:1) also masks the key.  
    - Examples log errors but not the API key directly.  
  - **Status**: Largely **PASS**, with caveats  
  - **Gaps / Risks**  
    - If an upstream service ever echoes secrets in HTTP error bodies, [rest.HTTPError.Error()](cci:2://file:///home/stig/dev/ai/zaguan/labs/xai-sdk-go/xai/internal/errors/errors.go:133:0-140:1) will include those in the error string; user code might log them.  
  - **Recommendations (P1)**  
    - Consider adding an option or wrapper to *redact response bodies* in [HTTPError.Error()](cci:2://file:///home/stig/dev/ai/zaguan/labs/xai-sdk-go/xai/internal/errors/errors.go:133:0-140:1) or expose a `SafeError()` variant for logging.  
    - Add a short “do not log raw errors from sensitive calls” warning in README/examples.

- **Hardcoded credentials**  
  - No hardcoded API keys or secrets found (only env vars and tests).  
  - **Status**: **PASS**

- **Bearer token format & transmission**  
  - `Authorization: Bearer <key>` used consistently in REST and gRPC (via metadata).  
  - **Status**: **PASS**

- **Environment variable handling**  
  - [Config.LoadFromEnvironment](cci:1://file:///home/stig/dev/ai/zaguan/labs/xai-sdk-go/xai/config.go:119:0-214:1) uses `XAI_API_KEY` and other `XAI_*` vars; parsing is robust (duration parsing, bools, ints).  
  - **Status**: **PASS**

### 1.1 TLS/SSL Configuration

- **HTTPS enforcement (REST)**  
  - Default base URL is `https://api.x.ai/v1`; `Insecure` flag switches to `http://`.  
  - **Status**: **PASS**, with explicit escape hatch  
- **Certificate validation**  
  - gRPC: TLS config uses `MinVersion: TLS1.2`. `SkipVerify` is configurable and defaults to false.  
  - REST: uses default `http.Transport` TLS handling (no explicit `InsecureSkipVerify`).  
  - **Status**: **PASS**, but `SkipVerify` is dangerous if enabled.  
- **Recommendations (P2)**  
  - Document clearly that `Insecure` and `SkipVerify` must only be used in local/test environments.  
  - Optionally emit a warning when `SkipVerify` is true (e.g., via a debug log hook if you add logging later).

### 1.1 Input Validation

- **User inputs**  
  - Chat & Models: validate required fields and ranges (e.g., [chat.Request.validate()](cci:1://file:///home/stig/dev/ai/zaguan/labs/xai-sdk-go/xai/chat/chat.go:492:0-548:1) checks model, messages, temperature, max_tokens).  
  - REST APIs: rely mostly on proto-level type safety; some methods validate required IDs (`Get*` methods).  
  - **Status**: **PARTIAL PASS**  
- **File upload size limits**  
  - [files.Client.Upload](cci:1://file:///home/stig/dev/ai/zaguan/labs/xai-sdk-go/xai/files/files.go:81:0-118:1) uses `io.ReadAll(reader)` with no size limit before marshalling into JSON. This ignores `constants.DefaultMaxFileSize`.  
  - **Status**: **FAIL (P1)**  
  - **Recommendation (P0/P1)**:  
    - Introduce an upload size limit (e.g., configurable with default 100MB) that:  
      - Wraps `reader` with `io.LimitedReader`;  
      - Returns `ErrFileTooLarge` (or similar) if exceeded.  
- **Path traversal / filesystem**  
  - SDK doesn’t touch local filesystem paths (only uses IDs in URLs and `io.Reader`/`io.ReadCloser`); no obvious traversal risk.  
  - **Status**: **PASS**  
- **SSRF / URL inputs**  
  - Image & embed APIs accept URLs but only forward them to xAI API; SDK doesn’t perform outbound fetches. SSRF is server-side concern.  
  - **Status**: **PASS (client-side)**  
- **Injection (SQL, command, etc.)**  
  - No shell or SQL usage; only JSON/Protobuf over HTTP/gRPC.  
  - **Status**: **PASS**

### 1.2 Data Protection

- **Logging of PII / contents**  
  - SDK library code does not log. Only examples log coarse-grained errors.  
  - [internal/metadata.SanitizeMetadata](cci:1://file:///home/stig/dev/ai/zaguan/labs/xai-sdk-go/xai/internal/metadata/metadata.go:260:0-275:1) masks keys if used.  
  - **Status**: **PASS**, with same caveat about logging HTTP error bodies.  
- **File contents not logged**  
  - [files.Upload](cci:1://file:///home/stig/dev/ai/zaguan/labs/xai-sdk-go/xai/files/files.go:81:0-118:1) reads content but never logs; only sends to server.  
  - **Status**: **PASS**
- **Response size limits**  
  - REST: `MaxResponseSize = 100MB` enforced via `io.LimitReader`.  
  - **Status**: **PASS**  
- **Malformed JSON handling**  
  - REST wrappers consistently check `protojson.Unmarshal` errors and return them; no panic.  
  - **Status**: **PASS**

### 1.3 Dependencies

- **Dependencies** ([go.mod](cci:7://file:///home/stig/dev/ai/zaguan/labs/xai-sdk-go/go.mod:0:0-0:0))  
  - Only gRPC, Protobuf, and standard `golang.org/x/*` packages. No obvious risky deps.  
  - **Status**: **PASS (structure)**, **Needs tool verification**  
- **Recommended commands**  
  - `go list -m all`  
  - `govulncheck ./...`  
  - For license scanning: use `go-licenses` or similar across modules.

### 1.4 Code Security

- **Static analysis / race detector**  
  - Not wired into CI.  
  - **Recommendations**  
    - `gosec ./...`  
    - `go test -race ./...` (likely fine; code looks race-safe).  
    - Add optional job in CI for the above.  

---

## 2. Performance Audit

### 2.1 HTTP Client Performance

- **Connection pooling & HTTP/2**  
  - [internal/rest.NewClient](cci:1://file:///home/stig/dev/ai/zaguan/labs/xai-sdk-go/xai/client.go:41:0-82:1) configures `http.Transport` with pooling, HTTP/2, gzip, and appropriate timeouts; matches [docs/PERFORMANCE.md](cci:7://file:///home/stig/dev/ai/zaguan/labs/xai-sdk-go/docs/PERFORMANCE.md:0:0-0:0).  
  - **Status**: **PASS**
- **Connection leaks**  
  - [rest.Client.Close()](cci:1://file:///home/stig/dev/ai/zaguan/labs/xai-sdk-go/xai/client.go:226:0-247:1) calls `CloseIdleConnections()`.  
  - **Gap**: [xai.Client.Close()](cci:1://file:///home/stig/dev/ai/zaguan/labs/xai-sdk-go/xai/client.go:226:0-247:1) never calls [restClient.Close()](cci:1://file:///home/stig/dev/ai/zaguan/labs/xai-sdk-go/xai/client.go:226:0-247:1), so HTTP idle conns remain.  
  - **Status**: **PARTIAL**  
  - **Recommendation (P1)**: call [c.restClient.Close()](cci:1://file:///home/stig/dev/ai/zaguan/labs/xai-sdk-go/xai/client.go:226:0-247:1) inside [xai.Client.Close()](cci:1://file:///home/stig/dev/ai/zaguan/labs/xai-sdk-go/xai/client.go:226:0-247:1) if non-nil.  
- **Idle timeouts, MaxConnsPerHost**  
  - Configured with sensible defaults (100 conns, 90s idle).  
  - **Status**: **PASS**

### 2.1 Memory Management

- **Buffer pooling & response limiting**  
  - Buffer pool for request JSON, `MaxResponseSize` cap; documented in PERFORMANCE.md.  
  - **Status**: **PASS**  
- **Upload path**  
  - [files.Upload](cci:1://file:///home/stig/dev/ai/zaguan/labs/xai-sdk-go/xai/files/files.go:81:0-118:1) reads entire file into memory before send.  
  - **Status**: **FAIL (same as 1.1)**; risk of high memory usage.  

### 2.1 Concurrency

- **Thread safety of REST & main client**  
  - [http.Client](cci:2://file:///home/stig/dev/ai/zaguan/labs/xai-sdk-go/xai/client.go:28:0-39:1) is safe for concurrent use; [rest.Client](cci:2://file:///home/stig/dev/ai/zaguan/labs/xai-sdk-go/xai/internal/rest/client.go:28:0-33:1) holds no mutable state.  
  - [xai.Client](cci:2://file:///home/stig/dev/ai/zaguan/labs/xai-sdk-go/xai/client.go:28:0-39:1) uses `sync.RWMutex` to guard config/metadata/conn; service accessors use `RLock()`.  
  - No obvious goroutine leaks; streaming uses per-call gRPC streams.  
  - **Status**: **PASS**, pending `go test -race`.

### 2.2 Benchmarking & 2.3 Resource Usage

- Only a couple of micro-benchmarks (config/client) exist; no HTTP/REST or chat-level benchmarks or pprof scripts.  
- **Status**: **NOT IMPLEMENTED** for plan items 2.2/2.3.  
- **Recommendation (P2)**  
  - Add benchmarks for at least:  
    - Chat [Sample](cci:1://file:///home/stig/dev/ai/zaguan/labs/xai-sdk-go/xai/client.go:491:0-496:1)/[Stream](cci:2://file:///home/stig/dev/ai/zaguan/labs/xai-sdk-go/xai/chat/chat.go:150:0-154:1) via [client.Chat()](cci:1://file:///home/stig/dev/ai/zaguan/labs/xai-sdk-go/xai/client.go:170:0-175:1);  
    - REST [Embed().Generate](cci:1://file:///home/stig/dev/ai/zaguan/labs/xai-sdk-go/xai/client.go:442:0-447:1) and [Files().Upload/List](cci:1://file:///home/stig/dev/ai/zaguan/labs/xai-sdk-go/xai/client.go:449:0-454:1).  
  - Add `go test -bench=. -benchmem ./xai/...` guidance in docs.

---

## 3. Code Quality Audit

### 3.1 Go Best Practices

- **Style & naming**  
  - Code is gofmt’d; CI enforces `gofmt -l .`.  
  - Naming is idiomatic; packages focused and small.  
  - **Status**: **PASS**, but no goimports/golangci-lint job.  
- **Error handling**  
  - No use of `panic` in library code.  
  - Clear error messages, with mapping from gRPC `codes.*` to human-readable strings (chat, models, errors pkg).  
  - Internal error type [internal/errors.Error](cci:2://file:///home/stig/dev/ai/zaguan/labs/xai-sdk-go/xai/internal/errors/errors.go:133:0-140:1) has type, code, stack, context; not yet used broadly in wrappers (they mostly use `fmt.Errorf`).  
  - **Status**: **PASS**, with potential enhancement to use structured errors more consistently.  
- **Doc comments**  
  - Most exported types/functions in `xai/*` have doc comments.  
  - Could be gaps in some smaller helpers; overall OK.  
  - **Status**: **PARTIAL PASS**

### 3.2 Complexity & Duplication

- Functions are short and focused; the largest file is [chat/chat.go](cci:7://file:///home/stig/dev/ai/zaguan/labs/xai-sdk-go/xai/chat/chat.go:0:0-0:0) but complexity per function appears modest. No obvious deeply nested control-flow.  
- No automated gocyclo/deadcode runs.  
- **Recommendation**: run `gocyclo -over 15 .` and `golangci-lint run` once and keep results.

### 3.3 Testing

- **Unit tests**  
  - [xai/client_test.go](cci:7://file:///home/stig/dev/ai/zaguan/labs/xai-sdk-go/xai/client_test.go:0:0-0:0), [xai/config_test.go](cci:7://file:///home/stig/dev/ai/zaguan/labs/xai-sdk-go/xai/config_test.go:0:0-0:0), `internal/*_test.go`, [chat/chat_test.go](cci:7://file:///home/stig/dev/ai/zaguan/labs/xai-sdk-go/xai/chat/chat_test.go:0:0-0:0) exist and cover config, metadata, client, chat builders.  
  - **Status**: **GOOD but not quantified**  
  - Recommendation: add coverage job: `go test -cover ./...` in CI.  
- **Integration tests**  
  - Present for chat, embed, files, image, auth (with build tags / env gating).  
  - [xai/INTEGRATION_TESTS.md](cci:7://file:///home/stig/dev/ai/zaguan/labs/xai-sdk-go/xai/INTEGRATION_TESTS.md:0:0-0:0) documents how to run them and cleanup patterns.  
  - **Status**: **PASS vs plan**  
- **Edge cases**  
  - Some validation paths are tested (config), but API wrappers often lack tests for invalid input (e.g., empty IDs).  
  - **Recommendation**: add more tests for nil/empty inputs and timeout/ctx cancellation.

---

## 4. API Design Audit

### 4.1 Interface Design

- **Consistency & ergonomics**  
  - All REST API clients: [NewClient(rest *rest.Client)](cci:1://file:///home/stig/dev/ai/zaguan/labs/xai-sdk-go/xai/client.go:41:0-82:1) + methods mirroring proto operations.  
  - Chat: fluent builder with `RequestOption` pattern; consistent with README.  
  - **Status**: **PASS**, but with a few **bugs/gaps**:
    - **Bug (P1)**: [xai.Client.WithTimeout](cci:1://file:///home/stig/dev/ai/zaguan/labs/xai-sdk-go/xai/config.go:372:0-376:1) and [WithAPIKey](cci:1://file:///home/stig/dev/ai/zaguan/labs/xai-sdk-go/xai/client.go:328:0-347:1) create a new [Client](cci:2://file:///home/stig/dev/ai/zaguan/labs/xai-sdk-go/xai/client.go:28:0-39:1) but **do not copy** `restClient`, `chatClient`, `modelsClient`, or mutex. Resulting client will have nil service clients; any use of [NewClient.WithTimeout().Embed()](cci:1://file:///home/stig/dev/ai/zaguan/labs/xai-sdk-go/xai/config.go:372:0-376:1) or [.Chat()](cci:1://file:///home/stig/dev/ai/zaguan/labs/xai-sdk-go/xai/client.go:170:0-175:1) will fail at runtime.  
    - **Bug (P2)**: [collections.Client.ListDocuments](cci:1://file:///home/stig/dev/ai/zaguan/labs/xai-sdk-go/xai/collections/collections.go:273:0-310:1) accepts `opts *ListDocumentsOptions` but uses `opts.CollectionID` outside nil-check; [ListDocuments(ctx, nil)](cci:1://file:///home/stig/dev/ai/zaguan/labs/xai-sdk-go/xai/collections/collections.go:273:0-310:1) would panic.  

### 4.2 Backwards Compatibility

- CHANGELOG documents breaking changes between 0.1.x, 0.2.x, 0.3.0.  
- README claims “v0.3.0 Ready” and `go get ...@v0.2.1`, while internal constants use `0.2.1`; docs are not fully synchronized.  
- **Status**: **DOCS INCONSISTENT**  
- **Recommendation (P1)**: Align README, CHANGELOG, internal version constants (`constants.DefaultUserAgent`, metadata default) and actual tags.

### 4.3 Context Usage

- All public API calls accept [context.Context](cci:1://file:///home/stig/dev/ai/zaguan/labs/xai-sdk-go/xai/internal/errors/errors.go:223:0-226:1).  
- `xai.Client.NewContext*` helper methods add metadata and timeouts.  
- `context.Background()` appears in **examples** and internal setup (creating the gRPC conn), but not in user-facing operations.  
- **Status**: **PASS** per practical intent; if you want strict “no Background in library code”, treat the gRPC dial usage as a deliberate exception.

---

## 5. Concurrency & Thread Safety

- **Shared state**  
  - All shared state inside [xai.Client](cci:2://file:///home/stig/dev/ai/zaguan/labs/xai-sdk-go/xai/client.go:28:0-39:1) guarded by `sync.RWMutex`.  
  - [rest.Client](cci:2://file:///home/stig/dev/ai/zaguan/labs/xai-sdk-go/xai/internal/rest/client.go:28:0-33:1) and service clients are stateless wrappers.  
- **Goroutine lifecycle**  
  - No long-lived goroutines spawned by SDK except those internal to gRPC/http.  
- **Race detector**  
  - Not yet run; no race-prone patterns spotted by inspection.  
- **Status**: **PASS (design)**, **Needs `go test -race`**.

---

## 6. Error Handling & Resilience

- **Error types & messages**  
  - Clear mapping from gRPC status codes to friendly errors in chat/models; HTTP errors via `HTTPError`.  
  - [internal/errors](cci:7://file:///home/stig/dev/ai/zaguan/labs/xai-sdk-go/xai/internal/errors:0:0-0:0) provides richer structured errors but is underused by high-level wrappers.  
  - **Status**: **PASS**, with improvement potential.  
- **Retries & timeouts**  
  - Timeouts: enforced via config and gRPC timeout interceptor; REST uses `http.Client.Timeout`.  
  - Retries: [grpcutil.RetryWithBackoff](cci:1://file:///home/stig/dev/ai/zaguan/labs/xai-sdk-go/xai/internal/grpcutil/serviceconfig.go:276:0-315:1) exists but is **not used**; Config has `MaxRetries`, `RetryBackoff`, `MaxBackoff` but no integration in client code.  
  - **Status**: **NO BUILT-IN RETRIES**  
  - **Recommendation (P2)**: either wire [RetryWithBackoff](cci:1://file:///home/stig/dev/ai/zaguan/labs/xai-sdk-go/xai/internal/grpcutil/serviceconfig.go:276:0-315:1) into critical operations or clearly document that retries are left to the caller.

---

## 7. Documentation Audit

- **Code documentation (godoc)**  
  - Generally good for the main packages; minor gaps in some internal helpers.  
  - `go doc` should work; not wired into CI.  
- **User docs (README, CHANGELOG, guides)**  
  - Very comprehensive, but **versions and status docs are out of sync**:  
    - README says v0.3.0 ready and 100% API coverage, while [SDK_STATUS.md](cci:7://file:///home/stig/dev/ai/zaguan/labs/xai-sdk-go/docs/SDK_STATUS.md:0:0-0:0) is dated v0.1.6 and says most REST APIs are not implemented; CHANGELOG [Unreleased] describes 0.3.0 with all APIs implemented, matching current code.  
    - README shows `go get ...@v0.2.1`; but code claims v0.3.0-ready.  
  - **Status**: **PARTIAL / OUTDATED**  
  - **Recommendation (P1)**:  
    - Update [SDK_STATUS.md](cci:7://file:///home/stig/dev/ai/zaguan/labs/xai-sdk-go/docs/SDK_STATUS.md:0:0-0:0) to reflect the current wrapper implementations and tests.  
    - Make README’s version, coverage table, and `go get` instruction consistent with the next tag.

- **Examples**  
  - They compile (by inspection) and use best practices (defer Close, env-based API key, basic error handling).  
  - Some examples use `context.Background()` directly rather than [client.NewContext](cci:1://file:///home/stig/dev/ai/zaguan/labs/xai-sdk-go/xai/client.go:177:0-197:1); not unsafe but slightly inconsistent with best practices.

---

## 8. Build & Deployment

- **Build matrix & Go version**  
  - CI uses Go 1.22 and 1.23; [go.mod](cci:7://file:///home/stig/dev/ai/zaguan/labs/xai-sdk-go/go.mod:0:0-0:0) says `go 1.24.0` / toolchain 1.24.6.  
  - That mismatch will matter once 1.24 semantics diverge.  
  - **Status**: **INCONSISTENT**  
  - **Recommendation (P1)**:  
    - Align `go` directive with actual supported version (e.g., `go 1.22` or `1.23`) or bump CI matrix to include 1.24 once available.  
- **Dependencies**  
  - Simple and pinned; `go mod tidy` is not in CI but Makefile likely covers it.  
- **CI**  
  - Currently: gofmt check + `go test ./...` + placeholder Buf lint.  
  - Missing: race, coverage, static analysis, govulncheck.  

---

## 9. Proto & gRPC

- [buf.yaml](cci:7://file:///home/stig/dev/ai/zaguan/labs/xai-sdk-go/buf.yaml:0:0-0:0) / [buf.gen.yaml](cci:7://file:///home/stig/dev/ai/zaguan/labs/xai-sdk-go/buf.gen.yaml:0:0-0:0) present and minimal; CHANGELOG + SDK_STATUS document full proto alignment.  
- gRPC dial options ([Config.CreateGRPCDialOptions](cci:1://file:///home/stig/dev/ai/zaguan/labs/xai-sdk-go/xai/config.go:305:0-357:1)) set min TLS 1.2, keepalive, interceptors.  
- No signs of manual edits to generated code under `proto/gen/go`.  
- **Status**: **PASS**

---

## 10. Compliance & Legal

- **LICENSE**: Apache 2.0 present.  
- **License headers**: Some files (e.g., [chat_integration_test.go](cci:7://file:///home/stig/dev/ai/zaguan/labs/xai-sdk-go/xai/chat/chat_integration_test.go:0:0-0:0)) have Apache headers; many others don’t.  
- **Unofficial disclaimer**: README has clear “unofficial, community-maintained” disclaimer and attribution to xAI.  
- **Dependencies**: No GPL dependencies.  
- **Status**:  
  - Licensing: **PASS**  
  - Headers: **PARTIAL** (low priority unless required by your policy).

---

## 11. Specific API Audits (Highlights)

- **REST client ([xai/internal/rest](cci:7://file:///home/stig/dev/ai/zaguan/labs/xai-sdk-go/xai/internal/rest:0:0-0:0))**  
  - Pooling, HTTP/2, gzip, MaxResponseSize, [Close()](cci:1://file:///home/stig/dev/ai/zaguan/labs/xai-sdk-go/xai/client.go:226:0-247:1) implemented.  
  - Gap: not hooked into [xai.Client.Close()](cci:1://file:///home/stig/dev/ai/zaguan/labs/xai-sdk-go/xai/client.go:226:0-247:1).  
- **Chat API (`xai/chat`)**  
  - Streaming & non-streaming paths robust, with detailed gRPC error mapping and request validation.  
  - **Function calling/tooling**: tool support is partially stubbed (comments and placeholder implementations for [ToolCalls](cci:1://file:///home/stig/dev/ai/zaguan/labs/xai-sdk-go/xai/chat/chat.go:748:0-753:1), [Choice](cci:2://file:///home/stig/dev/ai/zaguan/labs/xai-sdk-go/xai/chat/chat.go:85:0-87:1), [SetTools](cci:1://file:///home/stig/dev/ai/zaguan/labs/xai-sdk-go/xai/chat/chat.go:366:0-377:1)).  
  - **Status**: **Core chat solid**, **tooling incomplete vs README claims**.  
- **Files API (`xai/files`)**  
  - CRUD operations implemented; integration tests exist.  
  - **Issues**: no upload size limit; full file read into memory.  
- **Embed, Collections, Documents, Image, Deferred, Tokenizer, Sample**  
  - All wrappers implemented as thin, straightforward Protobuf/JSON marshaling layers.  
  - Deferred `Status.Result` is `interface{}`; not strongly typed.  
  - Collections [ListDocuments](cci:1://file:///home/stig/dev/ai/zaguan/labs/xai-sdk-go/xai/collections/collections.go:273:0-310:1) nil-opts bug (see above).  

---

## 12–16. Phases, Tools, Success Criteria, Issue Tracking

- **Automated tools**: Not wired into CI, but Makefile + docs already outline test and proto workflows. Security/perf tools from §13 are not yet integrated.  
- **Success criteria**:  
  - No panics in library code: **PASS** (panic search returned none).  
  - All tests pass: seems true for unit tests; integration tests gated by env.  
  - Thread safety: design looks sound; needs `-race` to fully confirm.  
  - Coverage, lints, performance benchmarks: not yet enforced.

---

## Recommended Next Actions (Prioritized)

- **P0–P1 (Before v0.3.0 release)**  
  - **Fix file upload size handling** ([files.Upload](cci:1://file:///home/stig/dev/ai/zaguan/labs/xai-sdk-go/xai/files/files.go:81:0-118:1)): enforce max size using a limit and return `ErrFileTooLarge`.  
  - **Teach [xai.Client.Close()](cci:1://file:///home/stig/dev/ai/zaguan/labs/xai-sdk-go/xai/client.go:226:0-247:1) to close REST client** ([restClient.Close()](cci:1://file:///home/stig/dev/ai/zaguan/labs/xai-sdk-go/xai/client.go:226:0-247:1)), so HTTP pools don’t leak.  
  - **Fix [Client.WithTimeout](cci:1://file:///home/stig/dev/ai/zaguan/labs/xai-sdk-go/xai/client.go:310:0-326:1) / [WithAPIKey](cci:1://file:///home/stig/dev/ai/zaguan/labs/xai-sdk-go/xai/client.go:328:0-347:1)** to carry over `restClient`, `chatClient`, `modelsClient`, mutex, and other fields.  
  - **Fix [collections.ListDocuments](cci:1://file:///home/stig/dev/ai/zaguan/labs/xai-sdk-go/xai/collections/collections.go:273:0-310:1) nil-opts bug** (don’t deref `opts` if nil).  
  - **Align versions/docs**: README, CHANGELOG, SDK_STATUS, default user-agent string, and `go get` snippet.

- **P1–P2 (Short term)**  
  - Add CI jobs for: `go test -race ./...`, `go test -cover ./...`, and `govulncheck ./...`.  
  - Add a minimal `golangci-lint` config and run it in CI.  
  - Decide and document stance on retries (no built-in vs limited automatic retries on clearly retryable codes).  
  - Clarify in docs/examples that logging full `HTTPError` values can expose response bodies; suggest redaction patterns.

- **P2 (Nice-to-have before or after 0.3.0)**  
  - Add benchmarks for chat & REST APIs.  
  - Strengthen integration tests around edge cases (timeouts, invalid inputs).  
  - Complete function-calling tooling in `chat` or update README if it’s best-effort/partial.

---

### Status Summary

- **Security**: Solid foundation; main concerns are unbounded upload reads and potential logging of sensitive response bodies; TLS and auth are correctly configured.  
- **Performance**: REST layer is well-optimized; main performance risk is large in-memory file uploads.  
- **API coverage**: All 11 APIs implemented as claimed, but function-calling support and some wrappers (deferred/documents result typing) are still incomplete vs “perfect” spec.  
- **Docs & build**: Very good overall, but internal status docs and versioning need a pass to sync with current implementation and planned v0.3.0 release.