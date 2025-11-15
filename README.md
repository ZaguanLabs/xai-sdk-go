# xAI SDK for Go

[![Go Version](https://img.shields.io/badge/Go-1.22+-blue.svg)](https://golang.org)
[![Release](https://img.shields.io/github/v/release/ZaguanLabs/xai-sdk-go)](https://github.com/ZaguanLabs/xai-sdk-go/releases/latest)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

The official Go SDK for xAI provides a first-class, idiomatic Go interface to xAI's powerful AI capabilities. This SDK enables Go developers to integrate chat completions, streaming responses, and upcoming features like file operations, image generation, and more.

> **Status**: **v0.1.2** - Bug fix release with corrected proto definitions for models and chat APIs.

## ✨ Features

- **🤖 Chat Completions** - Synchronous and streaming chat with message builders
- **🛠️ Function Calling** - Define and use tools in your chat completions
- **🧠 Reasoning & Search** - Control reasoning effort and perform searches
- **📝 Structured Outputs** - Get structured JSON and JSON schema outputs
- **🔐 Secure Authentication** - API key and Bearer token support with TLS
- **⚙️ Flexible Configuration** - Environment variables and programmatic config
- **🔄 Connection Management** - Health checks, retries, and keepalive
- **🛡️ Error Handling** - Comprehensive error types with gRPC integration
- **📊 Telemetry Ready** - Foundation for observability (coming soon)
- **🧪 Well Tested** - Comprehensive test coverage for all components

## 🚀 Quick Start

### Installation

```bash
go get github.com/ZaguanLabs/xai-sdk-go@v0.1.2
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

## 📚 Documentation

- **Development Plan**: [`docs/development-plan.md`](docs/development-plan.md) - Comprehensive phase-by-phase implementation plan
- **Examples**: [`examples/README.md`](examples/README.md) - Usage examples and tutorials
- **API Reference**: Available via [godoc.org](https://pkg.go.dev/github.com/ZaguanLabs/xai-sdk-go) (once published)

## 🗺️ Roadmap

The SDK is being developed in phases. Current status:

- ✅ **v0.1.0**: Foundation, proto, configuration, client, auth, basic and advanced chat
- ✅ **v0.1.1**: Fixed models API proto definitions and metadata handling
- ✅ **v0.1.2**: Fixed chat API proto definitions (package name, method names, message structures)
- 🚧 **Upcoming**: Deferred chat, files, image generation, and more.

See [`docs/development-plan.md`](docs/development-plan.md) for detailed progress.

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
