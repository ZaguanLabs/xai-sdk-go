# xAI SDK for Go

> **Status**: Pre-alpha (Phase 0 bootstrap). The SDK is under active development and not yet suitable for production use.

The goal of this repository is to deliver a first-class Go module that mirrors the xAI Python SDK (`docs/xai-sdk-python`) while embracing idiomatic Go APIs. The Go SDK will expose clients for chat, files, collections, models, tokenization, search tools, telemetry, and more over gRPC.

## Project Structure

```
.
├── docs/                  # Planning documents and reference material
├── examples/              # Go usage examples (coming soon)
├── proto/                 # Protobuf definitions and generated code (coming soon)
├── xai/                   # Go module source
│   └── internal/          # Shared internal utilities
└── .github/workflows/     # Continuous integration workflows
```

## Getting Started

1. Install Go **1.22** or newer.
2. Clone the repository and install Buf (https://buf.build) for protobuf tooling.
3. Follow the development phases outlined in [`docs/development-plan.md`](docs/development-plan.md).

## Documentation

- High-level porting plan: [`docs/xai-sdk-go-plan.md`](docs/xai-sdk-go-plan.md)
- Phase-by-phase development plan: [`docs/development-plan.md`](docs/development-plan.md)

Additional documentation (examples, API reference, telemetry, etc.) will be added as the project progresses through the planned phases.

## Contributing

Contributions are welcome once the initial milestones land. Please review [`CONTRIBUTING.md`](CONTRIBUTING.md) and [`CODE_OF_CONDUCT.md`](CODE_OF_CONDUCT.md) before opening an issue or pull request.

## License

This project is licensed under the Apache License 2.0 – see [`LICENSE`](LICENSE) for details.
