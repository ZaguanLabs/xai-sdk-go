# Protocol Buffer Definitions

This directory contains the protobuf definitions for the xAI SDK Go module.

## Structure

```
proto/
├── buf.yaml              # Buf configuration
├── buf.lock              # Buf dependency lock file
├── xai/v1/               # Proto definitions (v1 API)
│   ├── chat.proto        # Chat service definitions
│   ├── files.proto       # Files service definitions
│   └── models.proto      # Models service definitions
└── gen/go/               # Generated Go code
    └── xai/v1/           # Generated packages
```

## Regenerating Code

To regenerate Go bindings from proto definitions:

```bash
make proto
```

Or directly with buf:

```bash
~/go/bin/buf generate proto
```

## Important Notes

**Current Status**: The proto definitions in this directory are **stub implementations** created to bootstrap the project. They provide minimal service definitions for chat, files, and models APIs.

**TODO**: These stubs need to be replaced with the actual proto definitions from the xAI API once they become available. The Python SDK (`docs/xai-sdk-python`) only contains generated Python code, not the original `.proto` files.

## Adding New Services

1. Create a new `.proto` file in `proto/xai/v1/`
2. Define your service and messages
3. Run `make proto` to generate Go bindings
4. Import the generated package in your Go code

## Dependencies

- **Buf**: Install from https://buf.build
- **protoc-gen-go**: Installed automatically via buf plugins
- **protoc-gen-go-grpc**: Installed automatically via buf plugins

## Updating Proto Workflow

When the Python SDK proto sources are updated:

1. Copy updated `.proto` files to `proto/xai/v1/`
2. Run `~/go/bin/buf dep update proto` to update dependencies
3. Run `make proto` to regenerate Go bindings
4. Run `go mod tidy` to update Go module dependencies
5. Commit both the `.proto` files and generated code
