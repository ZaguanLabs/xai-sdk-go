# xAI SDK for Go

[![Go Version](https://img.shields.io/badge/Go-1.22+-blue.svg)](https://golang.org)
[![Release](https://img.shields.io/github/v/release/ZaguanLabs/xai-sdk-go)](https://github.com/ZaguanLabs/xai-sdk-go/releases/latest)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

**The unofficial Go SDK for xAI** provides a first-class, idiomatic Go interface to xAI's powerful AI capabilities. This SDK enables Go developers to integrate chat completions, streaming responses, embeddings, file operations, image generation, document search, and more.

> **Note**: This is an unofficial, community-maintained SDK and is not affiliated with or endorsed by xAI.

> **Status**: **v0.3.0 Released** - Production-ready with 100% API coverage (11/11 APIs), comprehensive examples, integration tests, performance optimizations, and security enhancements!

## ✨ Features

### Core APIs
- **🤖 Chat Completions** - Synchronous and streaming chat with message builders
- **🛠️ Function Calling** - Define and use tools in your chat completions
- **🧠 Reasoning & Search** - Control reasoning effort and perform searches
- **📝 Structured Outputs** - Get structured JSON and JSON schema outputs
- **🎯 Models** - List and retrieve available models

### REST APIs
- **🖼️ Image Generation** - Text-to-image and image-to-image generation
- **📄 Embeddings** - Generate embeddings for text and images
- **📁 Files** - Upload, download, list, and delete files
- **📚 Collections** - Manage document collections with 11 methods
- **🔍 Document Search** - Search across document collections
- **🔐 Auth** - API key validation and management
- **⏳ Deferred Completions** - Long-running completion support
- **🔤 Tokenizer** - Text tokenization utilities
- **📋 Sample** - Legacy text completion (Chat API recommended)

### Infrastructure & Performance
- **🔐 Secure Authentication** - API key and Bearer token support with TLS
- **⚙️ Flexible Configuration** - Environment variables and programmatic config
- **🔄 Connection Management** - Health checks, retries, and keepalive
- **🛡️ Error Handling** - Comprehensive error types with gRPC and REST integration
- **🧪 Well Tested** - Unit tests + integration tests for all components
- **⚡ High Performance** - Connection pooling, HTTP/2, buffer pooling (2-10x faster)
- **📊 Production Ready** - Optimized for low latency and high throughput
- **📚 Comprehensive Examples** - 8 detailed examples covering all APIs

## 🚀 Quick Start

### Installation

```bash
go get github.com/ZaguanLabs/xai-sdk-go@v0.3.0
```

### Basic Usage

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/ZaguanLabs/xai-sdk-go/xai"
	"github.com/ZaguanLabs/xai-sdk-go/xai/chat"
)

func main() {
	// Create client with API key from environment
	client, err := xai.NewClient(&xai.Config{APIKey: os.Getenv("XAI_API_KEY")})
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	// Create chat request
	req := chat.NewRequest("grok-1.5-flash",
		chat.WithTemperature(0.7),
		chat.WithMaxTokens(1000),
		chat.WithMessages(
			chat.System(chat.Text("You are a helpful assistant.")),
			chat.User(chat.Text("What is the capital of France?")),
		),
	)

	// Get response
	resp, err := req.Sample(context.Background(), client.Chat())
	if err != nil {
		log.Fatalf("Chat completion failed: %v", err)
	}

	fmt.Printf("Response: %s\n", resp.Content())
}
```

### Streaming Usage

```go
// Stream response in real-time
stream, err := req.Stream(context.Background(), client.Chat())
if err != nil {
    log.Fatalf("Stream failed: %v", err)
}
defer stream.Close()

for stream.Next() {
    chunk := stream.Current()
    if content := chunk.Content(); content != "" {
        fmt.Print(content)
    }
}
```

### Environment Setup

Set your xAI API key:

```bash
export XAI_API_KEY=your_api_key_here
```

Or configure programmatically:

```go
config := &xai.Config{
    APIKey:     "your-api-key",
    Timeout:    30 * time.Second,
}

client, err := xai.NewClient(config)
```

## 📁 Project Structure

```
.
├── docs/                  # Planning and reference documentation
├── examples/              # Usage examples and tutorials
├── proto/                 # Protocol buffer definitions and generated code
├── xai/                   # Main SDK source
│   ├── chat/              # Chat completion functionality
│   ├── internal/          # Shared utilities and interceptors
│   ├── client.go          # Main client implementation
│   └── config.go         # Configuration management
└── .github/workflows/     # CI/CD workflows
```

## 🔧 Configuration

The SDK supports flexible configuration through environment variables or programmatic setup:

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `XAI_API_KEY` | Your xAI API key | - |
| `XAI_HOST` | API host | `api.x.ai` |
| `XAI_TIMEOUT` | Request timeout | `30s` |
| `XAI_INSECURE` | Disable TLS (for testing) | `false` |
| `XAI_MAX_RETRIES` | Maximum retry attempts | `3` |

### Programmatic Configuration

```go
config := &xai.Config{
    APIKey:     "your-api-key",
    Host:       "api.x.ai",
    Timeout:    30 * time.Second,
    Insecure:   false,
    MaxRetries: 3,
}

client, err := xai.NewClient(config)
```

## 🧪 Development

### Prerequisites

- Go **1.22** or newer
- [Buf](https://buf.build) for protobuf code generation
- Make (for build automation)

### Building

```bash
# Format code
make fmt

# Run tests
make test

# Generate protobuf code
make proto

# Clean build artifacts
make clean
```

### Running Examples

```bash
# Basic chat example
go run examples/chat/basic/main.go

# Streaming example
go run examples/chat/streaming/main.go
```

## 🔒 Security Best Practices

### Error Logging

**⚠️ IMPORTANT**: When logging errors from API calls, be careful not to expose sensitive information:

```go
// ❌ DON'T: Log full error messages which may contain sensitive data
if err != nil {
    log.Printf("API error: %v", err) // May expose API keys or tokens in response body
}

// ✅ DO: Use SafeError() for HTTP errors
if err != nil {
    if httpErr, ok := err.(*rest.HTTPError); ok {
        log.Printf("API error: %s", httpErr.SafeError()) // Safe: only logs status code
    } else {
        log.Printf("API error: %v", err)
    }
}
```

### Configuration Security

- **Never hardcode API keys** - Use environment variables or secure configuration management
- **Use `Insecure` and `SkipVerify` only in local/test environments** - These flags disable TLS security
- **Rotate API keys regularly** - Follow security best practices for credential management

## 📚 Documentation

- **Development Plan**: [`docs/development-plan.md`](docs/development-plan.md) - Comprehensive phase-by-phase implementation plan
- **Examples**: [`examples/README.md`](examples/README.md) - Usage examples and tutorials
- **API Reference**: Available via [godoc.org](https://pkg.go.dev/github.com/ZaguanLabs/xai-sdk-go) (once published)

## � API Coverage

**100% Complete** - All 11 APIs from the xAI Python SDK are implemented!

| API | Transport | Status | Methods |
|-----|-----------|--------|---------|
| Chat | gRPC | ✅ Production Ready | All |
| Models | gRPC | ✅ Production Ready | All |
| Embed | REST | ✅ Complete | 1/1 |
| Files | REST | ✅ Complete | 6/6 |
| Auth | REST | ✅ Complete | 3/3 |
| Collections | REST | ✅ Complete | 11/11 |
| Image | REST | ✅ Complete | 1/1 |
| Deferred | REST | ✅ Complete | 2/2 |
| Documents | REST | ✅ Complete | 1/1 |
| Sample | REST | ✅ Complete | 1/1 |
| Tokenizer | REST | ✅ Complete | 1/1 |

**Total**: 28+ methods across 11 APIs

## 🗺️ Roadmap

### Released
- ✅ **v0.1.x**: Foundation, proto alignment, Chat and Models APIs
- ✅ **v0.2.0**: 100% proto alignment with xAI Python SDK v1.4.0
- ✅ **v0.2.1**: Hotfix for compilation errors

### Current (v0.3.0 - Released 2025-11-16)
- ✅ **REST Client Foundation**: Complete HTTP infrastructure with connection pooling
- ✅ **All 11 APIs Implemented**: 100% API coverage (28+ methods)
- ✅ **Production Ready**: Chat and Models tested in production
- ✅ **Feature Complete**: All REST APIs fully functional
- ✅ **Comprehensive Examples**: 8 detailed examples for all APIs
- ✅ **Integration Tests**: 15+ tests with build tag isolation
- ✅ **Performance Optimized**: 2-10x faster with connection pooling and HTTP/2
- ✅ **Security Hardened**: Audit completed with all P0-P1 issues resolved
- ✅ **Well Documented**: Performance guide, testing guide, audit reports, and API docs

### What's New in v0.3.0
- 🎉 **100% API Coverage**: All 11 xAI APIs fully implemented
- ⚡ **2-10x Performance**: Connection pooling, HTTP/2, buffer pooling
- 🔒 **Security Enhanced**: File size limits, safe error logging, TLS warnings
- 🧪 **Production Tested**: Comprehensive audit with all critical fixes
- 📚 **Complete Documentation**: Examples, guides, and best practices

### Next Steps
- 📝 Community feedback and improvements
- 🔧 Bug fixes and enhancements
- 🌟 Additional features based on user requests
- 🚀 v0.4.0 planning

See [`docs/SDK_STATUS.md`](docs/SDK_STATUS.md) for detailed status and [`CHANGELOG.md`](CHANGELOG.md) for release notes.

## 🤝 Contributing

Contributions are welcome! We encourage you to submit pull requests and bug reports to help improve the SDK. Please review our contributing guidelines:

- [`CONTRIBUTING.md`](CONTRIBUTING.md) - Contribution guidelines and process
- [`CODE_OF_CONDUCT.md`](CODE_OF_CONDUCT.md) - Community standards

### Development Workflow

1. Fork the repository
2. Create a feature branch
3. Implement your changes with tests
4. Run `make fmt test`
5. Submit a pull request

## 📄 License

This project is licensed under the Apache License 2.0 – see [`LICENSE`](LICENSE) for details.

## 🔗 Links

- [xAI API Documentation](https://docs.x.ai/)
- [Python SDK Reference](https://github.com/xai-org/xai-sdk-python)
- [Buf Protocol Buffers](https://buf.build/)
- [Go Documentation](https://golang.org/doc/)
