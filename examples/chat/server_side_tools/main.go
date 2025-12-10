// Package main demonstrates server-side tool usage including web search,
// X search, code execution, and attachment search.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ZaguanLabs/xai-sdk-go/xai"
	"github.com/ZaguanLabs/xai-sdk-go/xai/chat"
)

func main() {
	// Get API key from environment
	apiKey := os.Getenv("XAI_API_KEY")
	if apiKey == "" {
		log.Fatal("XAI_API_KEY environment variable is required")
	}

	// Create client
	client, err := xai.NewClient(&xai.Config{
		APIKey: apiKey,
	})
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	ctx := context.Background()

	// Example 1: Web Search Tool
	fmt.Println("=== Example 1: Web Search Tool ===")
	webSearchExample(ctx, client)

	// Example 2: X (Twitter) Search Tool
	fmt.Println("\n=== Example 2: X (Twitter) Search Tool ===")
	xSearchExample(ctx, client)

	// Example 3: Code Execution Tool
	fmt.Println("\n=== Example 3: Code Execution Tool ===")
	codeExecutionExample(ctx, client)

	// Example 4: Attachment Search Tool
	fmt.Println("\n=== Example 4: Attachment Search Tool ===")
	attachmentSearchExample(ctx, client)

	// Example 5: Multiple Server-Side Tools
	fmt.Println("\n=== Example 5: Multiple Server-Side Tools ===")
	multipleToolsExample(ctx, client)

	// Example 6: Mixed Client and Server-Side Tools
	fmt.Println("\n=== Example 6: Mixed Client and Server-Side Tools ===")
	mixedToolsExample(ctx, client)

	fmt.Println("\n✅ All server-side tool examples completed successfully!")
}

func webSearchExample(ctx context.Context, client *xai.Client) {
	// Create a chat with web search enabled
	req := chat.NewRequest("grok-beta",
		chat.WithMessage(chat.User(chat.Text("What are the latest developments in AI from xAI?"))),
		chat.WithServerTool(
			chat.WebSearchTool(
				chat.WithAllowedDomains("x.ai", "xai.com"),
				chat.WithImageUnderstanding(true),
			),
		),
	)
	req.SetMaxTokens(200)

	resp, err := req.Sample(ctx, client.Chat())
	if err != nil {
		log.Printf("Web search example failed: %v", err)
		return
	}

	fmt.Printf("Q: What are the latest developments in AI from xAI?\n")
	fmt.Printf("A (with web search): %s\n", resp.Content())
	fmt.Printf("Note: The model used web search to find current information\n")
}

func xSearchExample(ctx context.Context, client *xai.Client) {
	// Create a chat with X (Twitter) search enabled
	now := time.Now()
	lastWeek := now.AddDate(0, 0, -7)

	req := chat.NewRequest("grok-beta",
		chat.WithMessage(chat.User(chat.Text("What is @xai saying about Grok recently?"))),
		chat.WithServerTool(
			chat.XSearchTool(
				chat.WithAllowedXHandles("xai", "elonmusk"),
				chat.WithXDateRange(lastWeek, now),
				chat.WithXImageUnderstanding(true),
			),
		),
	)
	req.SetMaxTokens(200)

	resp, err := req.Sample(ctx, client.Chat())
	if err != nil {
		log.Printf("X search example failed: %v", err)
		return
	}

	fmt.Printf("Q: What is @xai saying about Grok recently?\n")
	fmt.Printf("A (with X search): %s\n", resp.Content())
	fmt.Printf("Note: The model searched X/Twitter for recent posts\n")
}

func codeExecutionExample(ctx context.Context, client *xai.Client) {
	// Create a chat with code execution enabled
	req := chat.NewRequest("grok-beta",
		chat.WithMessage(chat.User(chat.Text("Calculate the first 10 Fibonacci numbers using Python"))),
		chat.WithServerTool(chat.CodeExecutionTool()),
	)
	req.SetMaxTokens(300)

	resp, err := req.Sample(ctx, client.Chat())
	if err != nil {
		log.Printf("Code execution example failed: %v", err)
		return
	}

	fmt.Printf("Q: Calculate the first 10 Fibonacci numbers using Python\n")
	fmt.Printf("A (with code execution): %s\n", resp.Content())
	fmt.Printf("Note: The model executed Python code to calculate the result\n")
}

func attachmentSearchExample(ctx context.Context, client *xai.Client) {
	// Create a chat with attachment search enabled
	req := chat.NewRequest("grok-beta",
		chat.WithMessage(chat.User(chat.Text("Search my attachments for information about API keys"))),
		chat.WithServerTool(
			chat.AttachmentSearchTool(
				chat.WithAttachmentLimit(5),
			),
		),
	)
	req.SetMaxTokens(200)

	resp, err := req.Sample(ctx, client.Chat())
	if err != nil {
		log.Printf("Attachment search example failed: %v", err)
		return
	}

	fmt.Printf("Q: Search my attachments for information about API keys\n")
	fmt.Printf("A (with attachment search): %s\n", resp.Content())
	fmt.Printf("Note: The model searched uploaded file attachments\n")
}

func multipleToolsExample(ctx context.Context, client *xai.Client) {
	// Create a chat with multiple server-side tools
	req := chat.NewRequest("grok-beta",
		chat.WithMessage(chat.User(chat.Text("Find recent news about AI and analyze it with code"))),
		chat.WithServerTool(
			chat.WebSearchTool(chat.WithImageUnderstanding(true)),
			chat.CodeExecutionTool(),
		),
	)
	req.SetMaxTokens(300)

	resp, err := req.Sample(ctx, client.Chat())
	if err != nil {
		log.Printf("Multiple tools example failed: %v", err)
		return
	}

	fmt.Printf("Q: Find recent news about AI and analyze it with code\n")
	fmt.Printf("A (with web search + code execution): %s\n", resp.Content())
	fmt.Printf("Note: The model can use both web search and code execution\n")
}

func mixedToolsExample(ctx context.Context, client *xai.Client) {
	// Create a client-side function tool
	weatherTool := chat.NewTool("get_weather", "Get the current weather for a city")
	weatherTool.WithParameter("city", "string", "The city name", true)
	weatherTool.WithParameter("units", "string", "Temperature units (C or F)", true)

	// Create a chat with both client-side and server-side tools
	req := chat.NewRequest("grok-beta",
		chat.WithMessage(chat.System(chat.Text("You are a helpful assistant with access to tools."))),
		chat.WithMessage(chat.User(chat.Text("What's the weather in San Francisco and what are people saying about it on X?"))),
		chat.WithTool(weatherTool),              // Client-side function
		chat.WithServerTool(chat.XSearchTool()), // Server-side X search
	)
	req.SetMaxTokens(300)

	resp, err := req.Sample(ctx, client.Chat())
	if err != nil {
		log.Printf("Mixed tools example failed: %v", err)
		return
	}

	fmt.Printf("Q: What's the weather in San Francisco and what are people saying about it on X?\n")
	fmt.Printf("A (with function + X search): %s\n", resp.Content())
	fmt.Printf("Note: The model can use both client-side functions and server-side tools\n")
}
