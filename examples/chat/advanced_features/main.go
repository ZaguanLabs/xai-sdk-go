// Package main demonstrates advanced chat features including conversation continuity,
// parallel tool calls, message storage, and encrypted content.
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

	// Example 1: Conversation Continuity with PreviousResponseID
	fmt.Println("=== Example 1: Conversation Continuity ===")

	// First message in conversation
	req1 := chat.NewRequest("grok-beta",
		chat.WithMessage(chat.User(chat.Text("What is the capital of France?"))),
		chat.WithStoreMessages(true), // Store this conversation
	)
	req1.SetMaxTokens(50)

	resp1, err := req1.Sample(ctx, client.Chat())
	if err != nil {
		log.Fatalf("First message failed: %v", err)
	}
	fmt.Printf("Q: What is the capital of France?\n")
	fmt.Printf("A: %s\n", resp1.Content())
	fmt.Printf("Response ID: %s\n\n", resp1.ID())

	// Follow-up message referencing previous response
	req2 := chat.NewRequest("grok-beta",
		chat.WithMessage(chat.User(chat.Text("What is its population?"))),
		chat.WithPreviousResponseID(resp1.ID()), // Reference previous conversation
		chat.WithStoreMessages(true),
	)
	req2.SetMaxTokens(50)

	resp2, err := req2.Sample(ctx, client.Chat())
	if err != nil {
		log.Fatalf("Follow-up message failed: %v", err)
	}
	fmt.Printf("Q: What is its population? (referencing previous response)\n")
	fmt.Printf("A: %s\n\n", resp2.Content())

	// Example 2: Parallel Tool Calls
	fmt.Println("=== Example 2: Parallel Tool Calls ===")
	req3 := chat.NewRequest("grok-beta",
		chat.WithMessage(chat.System(chat.Text("You are a helpful assistant with access to tools."))),
		chat.WithMessage(chat.User(chat.Text("Get the weather in Paris and New York."))),
		chat.WithParallelToolCalls(true), // Execute tool calls in parallel
	)
	req3.SetMaxTokens(100)

	resp3, err := req3.Sample(ctx, client.Chat())
	if err != nil {
		log.Fatalf("Parallel tool calls example failed: %v", err)
	}
	fmt.Printf("Response: %s\n", resp3.Content())
	fmt.Printf("Note: With parallel_tool_calls=true, multiple tools can execute simultaneously\n\n")

	// Example 3: Sequential Tool Calls
	fmt.Println("=== Example 3: Sequential Tool Calls ===")
	req4 := chat.NewRequest("grok-beta",
		chat.WithMessage(chat.System(chat.Text("You are a helpful assistant with access to tools."))),
		chat.WithMessage(chat.User(chat.Text("Get the weather in Paris and New York."))),
		chat.WithParallelToolCalls(false), // Execute tool calls sequentially
	)
	req4.SetMaxTokens(100)

	resp4, err := req4.Sample(ctx, client.Chat())
	if err != nil {
		log.Fatalf("Sequential tool calls example failed: %v", err)
	}
	fmt.Printf("Response: %s\n", resp4.Content())
	fmt.Printf("Note: With parallel_tool_calls=false, tools execute one at a time\n\n")

	// Example 4: Message Storage Control
	fmt.Println("=== Example 4: Message Storage Control ===")

	// Store messages for future reference
	req5 := chat.NewRequest("grok-beta",
		chat.WithMessage(chat.User(chat.Text("Remember this: My favorite color is blue."))),
		chat.WithStoreMessages(true), // Explicitly store messages
	)
	req5.SetMaxTokens(50)

	resp5, err := req5.Sample(ctx, client.Chat())
	if err != nil {
		log.Fatalf("Store messages example failed: %v", err)
	}
	fmt.Printf("Stored message - Response: %s\n", resp5.Content())

	// Don't store messages (ephemeral conversation)
	req6 := chat.NewRequest("grok-beta",
		chat.WithMessage(chat.User(chat.Text("This is a temporary message."))),
		chat.WithStoreMessages(false), // Don't store this conversation
	)
	req6.SetMaxTokens(50)

	resp6, err := req6.Sample(ctx, client.Chat())
	if err != nil {
		log.Fatalf("Ephemeral message example failed: %v", err)
	}
	fmt.Printf("Ephemeral message - Response: %s\n\n", resp6.Content())

	// Example 5: Encrypted Content
	fmt.Println("=== Example 5: Encrypted Content ===")
	req7 := chat.NewRequest("grok-beta",
		chat.WithMessage(chat.User(chat.Text("Sensitive information: Account number 12345"))),
		chat.WithUseEncryptedContent(true), // Enable encryption
	)
	req7.SetMaxTokens(50)

	resp7, err := req7.Sample(ctx, client.Chat())
	if err != nil {
		log.Fatalf("Encrypted content example failed: %v", err)
	}
	fmt.Printf("Response (encrypted): %s\n", resp7.Content())
	fmt.Printf("Note: Content is encrypted for enhanced security\n\n")

	// Example 6: All Advanced Features Combined
	fmt.Println("=== Example 6: All Advanced Features Combined ===")
	req8 := chat.NewRequest("grok-beta")
	req8.SetMessages(
		*chat.System(chat.Text("You are a secure assistant.")),
		*chat.User(chat.Text("Process this sensitive request.")),
	)

	// Advanced features
	req8.SetParallelToolCalls(true)
	req8.SetStoreMessages(true)
	req8.SetUseEncryptedContent(true)

	// Standard parameters
	req8.SetTemperature(0.7)
	req8.SetMaxTokens(100)
	req8.SetUser("secure-user-123")

	resp8, err := req8.Sample(ctx, client.Chat())
	if err != nil {
		log.Fatalf("Combined features example failed: %v", err)
	}
	fmt.Printf("Response: %s\n", resp8.Content())
	fmt.Printf("Features enabled:\n")
	fmt.Printf("  - Parallel tool calls: true\n")
	fmt.Printf("  - Store messages: true\n")
	fmt.Printf("  - Encrypted content: true\n")
	fmt.Printf("  - User tracking: secure-user-123\n")

	fmt.Println("\n✅ All advanced feature examples completed successfully!")
	fmt.Println("\n📊 Advanced features demonstrated:")
	fmt.Println("  - Conversation continuity (PreviousResponseID)")
	fmt.Println("  - Parallel tool execution")
	fmt.Println("  - Sequential tool execution")
	fmt.Println("  - Message storage control")
	fmt.Println("  - Encrypted content")
	fmt.Println("  - Combined advanced features")
}
