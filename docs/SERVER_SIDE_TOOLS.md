# Server-Side Tools Guide

This guide explains how to use server-side tools in the xAI SDK for Go. Server-side tools are executed by the xAI backend and provide powerful capabilities like web search, code execution, and document search.

## Overview

The xAI SDK supports two types of tools:

1. **Client-Side Tools (Functions)**: Custom functions you define and execute in your application
2. **Server-Side Tools**: Built-in tools executed by xAI's backend

## Available Server-Side Tools

### 1. Web Search (`WebSearchTool`)

Enables the model to search the web for current information.

```go
req := chat.NewRequest("grok-beta",
    chat.WithMessage(chat.User(chat.Text("What are the latest AI developments?"))),
    chat.WithServerTool(
        chat.WebSearchTool(
            chat.WithAllowedDomains("example.com", "news.com"),
            chat.WithExcludedDomains("spam.com"),
            chat.WithImageUnderstanding(true),
        ),
    ),
)
```

**Options:**
- `WithAllowedDomains(domains ...string)` - Restrict search to specific domains
- `WithExcludedDomains(domains ...string)` - Exclude specific domains
- `WithImageUnderstanding(bool)` - Enable image analysis in search results

### 2. X (Twitter) Search (`XSearchTool`)

Enables the model to search X/Twitter for posts and information.

```go
now := time.Now()
lastWeek := now.AddDate(0, 0, -7)

req := chat.NewRequest("grok-beta",
    chat.WithMessage(chat.User(chat.Text("What is @xai saying recently?"))),
    chat.WithServerTool(
        chat.XSearchTool(
            chat.WithAllowedXHandles("xai", "elonmusk"),
            chat.WithExcludedXHandles("spam_account"),
            chat.WithXDateRange(lastWeek, now),
            chat.WithXImageUnderstanding(true),
            chat.WithXVideoUnderstanding(true),
        ),
    ),
)
```

**Options:**
- `WithAllowedXHandles(handles ...string)` - Restrict to specific X handles
- `WithExcludedXHandles(handles ...string)` - Exclude specific X handles
- `WithXDateRange(from, to time.Time)` - Set date range for posts
- `WithXImageUnderstanding(bool)` - Enable image analysis in posts
- `WithXVideoUnderstanding(bool)` - Enable video analysis in posts

### 3. Code Execution (`CodeExecutionTool`)

Enables the model to execute code (Python, etc.) to perform calculations or data processing.

```go
req := chat.NewRequest("grok-beta",
    chat.WithMessage(chat.User(chat.Text("Calculate the first 10 Fibonacci numbers"))),
    chat.WithServerTool(chat.CodeExecutionTool()),
)
```

**No configuration options** - Simply enable it and the model will execute code as needed.

### 4. Collections Search (`CollectionsSearchTool`)

Enables the model to search within your document collections.

```go
req := chat.NewRequest("grok-beta",
    chat.WithMessage(chat.User(chat.Text("Search my technical docs for API information"))),
    chat.WithServerTool(
        chat.CollectionsSearchTool(
            []string{"collection-id-1", "collection-id-2"},
            chat.WithCollectionsLimit(10),
        ),
    ),
)
```

**Options:**
- `collectionIDs []string` - Required: IDs of collections to search
- `WithCollectionsLimit(int32)` - Maximum number of results

### 5. Document Search (`DocumentSearchTool`)

Enables the model to search within uploaded documents.

```go
req := chat.NewRequest("grok-beta",
    chat.WithMessage(chat.User(chat.Text("Find information about API keys in my documents"))),
    chat.WithServerTool(
        chat.DocumentSearchTool(
            chat.WithDocumentLimit(5),
        ),
    ),
)
```

**Options:**
- `WithDocumentLimit(int32)` - Maximum number of document results

### 6. MCP (Model Context Protocol) (`MCPTool`)

Enables the model to interact with MCP servers for extended functionality.

```go
req := chat.NewRequest("grok-beta",
    chat.WithMessage(chat.User(chat.Text("Query the database"))),
    chat.WithServerTool(
        chat.MCPTool(
            "my-mcp-server",
            "https://mcp.example.com",
            chat.WithMCPDescription("Database query server"),
            chat.WithMCPAllowedTools("query", "search"),
            chat.WithMCPAuthorization("Bearer token123"),
            chat.WithMCPExtraHeaders(map[string]string{
                "X-Custom-Header": "value",
            }),
        ),
    ),
)
```

**Options:**
- `serverLabel string` - Required: Label for the MCP server
- `serverURL string` - Required: URL of the MCP server
- `WithMCPDescription(string)` - Description of the server
- `WithMCPAllowedTools(...string)` - Restrict which MCP tools can be called
- `WithMCPAuthorization(string)` - Authorization header
- `WithMCPExtraHeaders(map[string]string)` - Additional HTTP headers

## Combining Multiple Tools

You can provide multiple server-side tools in a single request:

```go
req := chat.NewRequest("grok-beta",
    chat.WithMessage(chat.User(chat.Text("Research AI trends and analyze with code"))),
    chat.WithServerTool(
        chat.WebSearchTool(),
        chat.CodeExecutionTool(),
    ),
)
```

## Mixing Client-Side and Server-Side Tools

You can combine both types of tools:

```go
// Define a client-side function
weatherTool := chat.NewTool("get_weather", "Get weather for a city")
weatherTool.WithParameter("city", "string", "City name", true)

req := chat.NewRequest("grok-beta",
    chat.WithMessage(chat.User(chat.Text("What's the weather and what are people saying about it?"))),
    chat.WithTool(weatherTool),              // Client-side function
    chat.WithServerTool(chat.XSearchTool()), // Server-side X search
)
```

## Best Practices

### 1. Choose the Right Tool

- **Web Search**: For current events, news, or general web information
- **X Search**: For social media trends, public opinion, or X-specific content
- **Code Execution**: For calculations, data processing, or algorithmic tasks
- **Document/Collections Search**: For searching your own uploaded content
- **MCP**: For custom integrations with external systems

### 2. Configure Appropriately

- Use domain/handle restrictions to improve relevance and reduce noise
- Set reasonable limits to control response size and cost
- Enable image/video understanding only when needed

### 3. Handle Tool Calls

Server-side tools are executed automatically by the backend. The model will:
1. Decide when to use a tool based on the conversation
2. Execute the tool on the server
3. Incorporate the results into its response

You don't need to handle tool execution yourself - just provide the tools and the model handles the rest.

### 4. Error Handling

```go
resp, err := req.Sample(ctx, client.Chat())
if err != nil {
    log.Printf("Request failed: %v", err)
    return
}

// Check if tools were used (implementation depends on response structure)
fmt.Printf("Response: %s\n", resp.Content())
```

## Feature Parity with Python SDK

The Go SDK now has **100% feature parity** with the Python SDK for server-side tools:

| Feature | Python SDK | Go SDK |
|---------|-----------|--------|
| Web Search | ✅ | ✅ |
| X Search | ✅ | ✅ |
| Code Execution | ✅ | ✅ |
| Collections Search | ✅ | ✅ |
| Document Search | ✅ | ✅ |
| MCP | ✅ | ✅ |
| Client-Side Functions | ✅ | ✅ |
| Mixed Tools | ✅ | ✅ |

## Examples

See the complete example at `examples/chat/server_side_tools/main.go` for demonstrations of:
- Each individual server-side tool
- Multiple tools in one request
- Mixing client-side and server-side tools

## API Reference

### Core Functions

- `WebSearchTool(opts ...WebSearchOption) *ServerTool`
- `XSearchTool(opts ...XSearchOption) *ServerTool`
- `CodeExecutionTool() *ServerTool`
- `CollectionsSearchTool(collectionIDs []string, opts ...CollectionsSearchOption) *ServerTool`
- `DocumentSearchTool(opts ...DocumentSearchOption) *ServerTool`
- `MCPTool(serverLabel, serverURL string, opts ...MCPOption) *ServerTool`

### Request Options

- `WithServerTool(tools ...*ServerTool) RequestOption` - Add server-side tools
- `WithTool(tools ...*Tool) RequestOption` - Add client-side function tools

## Troubleshooting

### Tool Not Being Used

If the model isn't using a tool when you expect it to:
1. Make sure your prompt clearly indicates the need for that tool
2. Check that the tool is properly configured
3. Verify the model supports the tool (use `grok-beta` or later)

### Tool Execution Errors

If you encounter errors:
1. Check your tool configuration (domains, handles, limits, etc.)
2. Verify API permissions for the tool type
3. Ensure collection/document IDs are valid
4. Check MCP server availability and authentication

### Performance Considerations

- Server-side tools add latency (web searches, code execution take time)
- Use `SetMaxTokens()` to control response length
- Consider using `SetParallelToolCalls(true)` for multiple tools

## Migration from Client-Side Only

If you were previously using only client-side tools:

**Before:**
```go
req := chat.NewRequest("grok-beta",
    chat.WithTool(myFunction),
)
```

**After (with server-side tools):**
```go
req := chat.NewRequest("grok-beta",
    chat.WithTool(myFunction),                // Client-side
    chat.WithServerTool(chat.WebSearchTool()), // Server-side
)
```

No breaking changes - existing code continues to work!
