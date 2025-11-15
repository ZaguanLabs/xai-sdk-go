# xAI SDK Go Examples

This directory contains examples for the xAI SDK for Go.

## Running the Examples

To run the examples, you need to have Go installed and the `XAI_API_KEY` environment variable set to your xAI API key.

You can run a specific example by running the following command:

```sh
go run examples/chat/basic.go
```

You can also build and run all examples using the `make examples` command from the root of the repository.

## Examples

- `auth/validate.go`: Validates the API key.
- `chat/basic.go`: Basic chat completion.
- `chat/function_calling.go`: Chat completion with function calling.
- `chat/reasoning.go`: Chat completion with reasoning.
- `chat/search.go`: Chat completion with search.
- `chat/streaming.go`: Streaming chat completion.
- `chat/structured_outputs.go`: Chat completion with structured outputs.
- `collections/manage.go`: Manage collections.
- `image/generation.go`: Image generation.
- `models/list.go`: List available models.
- `tokenizer/encode.go`: Tokenize text.
