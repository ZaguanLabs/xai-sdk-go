# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.1.6] - 2025-11-15

### Changed
- **Internal Version Strings**: Updated all internal version strings to match release tag (0.1.6)
- No functional changes from v0.1.5

### Note
This release ensures internal version strings match the git tag. v0.1.5 was tagged with internal version strings still set to 0.1.4.

## [0.1.5] - 2025-11-15

### Fixed
- **Message Proto Field Order**: Corrected Message proto to match official xAI Python SDK v1.4.0
  - Field 1: `repeated Content content` (was incorrectly field 2)
  - Field 2: `MessageRole role` (was incorrectly field 1)
  - Field 3: `string name` (was incorrectly field 4)
  - Field 4: `repeated ToolCall tool_calls` (was incorrectly field 5)
  - Field 5: `string reasoning_content` (was incorrectly field 3)
  - Field 6: `string encrypted_content` (unchanged)
- **Wire Format Encoding**: Now correctly encodes messages with content as field 1 (LengthDelimited) and role as field 2 (Varint)
- **Content Type**: Reverted content from string back to `repeated Content` array to match official API

### Changed
- Extracted official proto definitions from xAI Python SDK v1.4.0 distribution
- Updated all message builders and validators to use correct field order
- CompletionMessage (string content) now properly converts to Message (Content array)

## [0.1.4] - 2025-11-15

### Fixed
- **Message Content Type**: Changed Message.content from `repeated Content` to `string` (later reverted in v0.1.5)
- Attempted to fix wire type mismatch by simplifying content field

### Note
This version was superseded by v0.1.5 which uses the correct proto structure from official xAI SDK.

## [0.1.3] - 2025-11-15

### Fixed
- **Message Proto Field Order**: Attempted to fix wire type error by reordering Message fields
  - Changed role from field 2 to field 1
  - Changed content from field 1 to field 2
- Wire format encoding updated to match new field order

### Note
This version still had incorrect field order. The correct order was discovered in v0.1.5 by extracting official proto definitions.

## [0.1.2] - 2025-11-15

### Fixed
- **Chat API Proto Definitions**: Corrected proto package name from `xai.v1` to `xai_api` to match actual xAI API
- **Chat API RPC Methods**: Changed from `CreateChatCompletion`/`StreamChatCompletion` to `GetCompletion`/`GetCompletionChunk`
- **Chat Message Structures**: Updated to use enums for MessageRole, ReasoningEffort, ToolMode, FormatType, SearchMode
- **Chat Response Structure**: Changed from `Choices` to `Outputs` with `CompletionOutput` and `CompletionMessage`
- **Chat Streaming**: Updated to use `GetChatCompletionChunk` with `CompletionOutputChunk` and `Delta`
- **Token Usage**: Changed from `Usage` to `SamplingUsage` to match proto
- **Message Content**: Updated from string to array of `Content` objects
- **SearchParameters**: Updated fields to match new proto (Mode, ReturnCitations, MaxSearchResults)
- **ResponseFormat**: Changed from oneof to FormatType enum with schema string
- **Tool Definitions**: Updated to match new proto structure with Function field

### Changed
- All chat client wrapper code updated to work with new proto structures
- Message builders now convert between user-friendly strings and proto enums
- Helper functions added for enum conversions (roleToProto, roleFromProto, etc.)

## [0.1.1] - 2025-11-15

### Fixed
- **Models API Proto Definitions**: Corrected proto package name from `xai.v1` to `xai_api` to match actual xAI API
- **Models API RPC Methods**: Changed from generic `ListModels` to specific methods: `ListLanguageModels`, `ListEmbeddingModels`, `ListImageGenerationModels`
- **Proto Field Numbers**: Fixed field number ordering in all model message types (LanguageModel, EmbeddingModel, ImageGenerationModel) to match server wire format
- **Metadata Handling**: Changed from `metadata.NewOutgoingContext()` to `metadata.AppendToOutgoingContext()` to preserve gRPC internal headers
- **Content-Type Header**: Removed manual content-type interceptor that was interfering with gRPC's automatic header handling

### Changed
- Updated models client API to use type-specific methods instead of generic `List()` and `Get()`
- Enhanced models example with comprehensive debug logging and detailed model information display

### Added
- Debug logging throughout models example for better troubleshooting
- Support for all xAI language models including grok-2, grok-3, grok-4 variants

## [0.1.0] - 2025-11-15

### Added

- Initial release of the xAI SDK for Go.
- Core chat functionality with synchronous and streaming responses.
- Function calling, reasoning, and search capabilities.
- Structured outputs with JSON and JSON schema.
- Secure authentication with API key and Bearer token support.
- Flexible configuration with environment variables and programmatic setup.
- Connection management with health checks, retries, and keepalive.
- Comprehensive error handling with gRPC integration.
- Foundational support for telemetry.
- Comprehensive test coverage.
- Examples for all major features.
- Documentation and development plan.
- CI/CD workflows for testing and releases.
