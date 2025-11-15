# Feature Parity: Go vs. Python SDK

This document tracks the feature parity between the Go and Python xAI SDKs.

## Go SDK Features

- **Chat Completions**: Synchronous and streaming chat with message builders.
- **Function Calling**: Define and use tools in your chat completions.
- **Reasoning & Search**: Control reasoning effort and perform searches.
- **Structured Outputs**: Get structured JSON and JSON schema outputs.
- **Authentication**: API key and Bearer token support with TLS.
- **Configuration**: Environment variables and programmatic config.
- **Connection Management**: Health checks, retries, and keepalive.
- **Error Handling**: Comprehensive error types with gRPC integration.
- **Files**: Upload and download files.
- **Images**: Image generation.
- **Models**: List and get models.
- **Tokenizer**: Tokenize text.
- **Collections**: Manage collections and documents.
- **Telemetry**: Foundational support for OpenTelemetry.

## Python SDK Features

The Python SDK exposes the following modules:

- **Auth**: Authentication with API key.
- **Chat**: Synchronous and streaming chat completion.
- **Collections**: Manage collections and documents.
- **Files**: Upload and download files.
- **Image**: Image generation.
- **Models**: List and get models.
- **Tokenizer**: Tokenize text.

## Parity Status

**Conclusion**: The Go SDK has reached feature parity with the Python SDK in terms of the core modules provided. Both SDKs offer clients for `auth`, `chat`, `collections`, `files`, `image`, `models`, and `tokenizer`.
