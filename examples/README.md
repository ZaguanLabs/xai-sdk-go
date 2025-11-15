# xAI SDK Go Examples

This directory contains example code demonstrating how to use the xAI SDK for Go.

## Prerequisites

Before running the examples, you need to set your xAI API key as an environment variable:

```bash
export XAI_API_KEY=your_api_key_here
```

## Chat Examples

### Basic Chat Completion

The `basic.go` example demonstrates how to perform a basic chat completion request:

```bash
go run examples/chat/basic.go
```

### Streaming Chat Completion

The `streaming.go` example demonstrates how to perform a streaming chat completion request:

```bash
go run examples/chat/streaming.go
```

## Running Examples

To run any example, use the `go run` command followed by the path to the example file:

```bash
go run examples/chat/basic.go
```

Make sure you're in the root directory of the repository when running the commands.