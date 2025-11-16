package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/ZaguanLabs/xai-sdk-go/xai"
	"github.com/ZaguanLabs/xai-sdk-go/xai/chat"
	"google.golang.org/protobuf/encoding/protojson"
)

func main() {
	// Get API key from environment
	apiKey := os.Getenv("XAI_API_KEY")
	if apiKey == "" {
		log.Fatal("XAI_API_KEY environment variable is required")
	}

	// Create client
	client, err := xai.NewClient(&xai.Config{APIKey: apiKey})
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	// Small 1x1 red pixel PNG in base64 (like OpenWebUI would send)
	base64Image := "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mP8z8DwHwAFBQIAX8jx0gAAAABJRU5ErkJggg=="

	fmt.Println("=== Image Base64 Diagnostic Test ===\n")
	fmt.Printf("Base64 image length: %d characters\n", len(base64Image))
	fmt.Printf("Image prefix: %s...\n\n", base64Image[:50])

	// Create a request with base64 image
	req := chat.NewRequest("grok-2-vision",
		chat.WithMessages(
			chat.User(
				chat.Text("What color is this 1x1 pixel image? It should be red."),
				chat.Image(base64Image, chat.ImageDetailHigh),
			),
		),
		chat.WithMaxTokens(100),
		chat.WithTemperature(0.0),
	)

	// Show the proto structure
	proto := req.Proto()
	fmt.Println("=== Proto Structure ===")
	fmt.Printf("Model: %s\n", proto.Model)
	fmt.Printf("Messages count: %d\n", len(proto.Messages))

	if len(proto.Messages) > 0 {
		msg := proto.Messages[0]
		fmt.Printf("Message role: %v\n", msg.Role)
		fmt.Printf("Content parts: %d\n", len(msg.Content))

		for i, content := range msg.Content {
			fmt.Printf("\nContent[%d]:\n", i)
			if content.Text != "" {
				fmt.Printf("  Type: Text\n")
				fmt.Printf("  Text: %q\n", content.Text)
			}
			if content.ImageUrl != nil {
				fmt.Printf("  Type: Image\n")
				fmt.Printf("  ImageUrl length: %d\n", len(content.ImageUrl.ImageUrl))
				fmt.Printf("  ImageUrl prefix: %s...\n", content.ImageUrl.ImageUrl[:50])
				fmt.Printf("  Detail: %v\n", content.ImageUrl.Detail)
			}
			if content.File != nil {
				fmt.Printf("  Type: File\n")
				fmt.Printf("  FileId: %s\n", content.File.FileId)
			}
		}
	}

	// Serialize to JSON to see exactly what's being sent
	jsonBytes, err := protojson.Marshal(proto)
	if err != nil {
		log.Fatalf("Failed to marshal proto to JSON: %v", err)
	}

	fmt.Println("\n=== Proto JSON (formatted) ===")
	var jsonData map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &jsonData); err != nil {
		log.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	prettyJSON, err := json.MarshalIndent(jsonData, "", "  ")
	if err != nil {
		log.Fatalf("Failed to format JSON: %v", err)
	}
	fmt.Println(string(prettyJSON))

	// Actually send the request
	fmt.Println("\n=== Sending Request to xAI API ===")
	ctx := context.Background()

	response, err := req.Sample(ctx, client.Chat())
	if err != nil {
		log.Fatalf("Request failed: %v", err)
	}

	// Show the response
	fmt.Println("\n=== Response ===")
	fmt.Printf("ID: %s\n", response.ID())
	fmt.Printf("Model: %s\n", response.Model())
	fmt.Printf("Finish Reason: %s\n", response.FinishReason())

	if len(response.Choices()) > 0 {
		choice := response.Choices()[0]
		fmt.Printf("\nContent:\n%s\n", choice.Message().Content())
	}

	// Show usage
	if usage := response.Usage(); usage != nil {
		fmt.Printf("\nUsage:\n")
		fmt.Printf("  Prompt tokens: %d\n", usage.PromptTokens())
		fmt.Printf("  Completion tokens: %d\n", usage.CompletionTokens())
		fmt.Printf("  Total tokens: %d\n", usage.TotalTokens())
	}

	fmt.Println("\n=== Test Complete ===")
	fmt.Println("If the model correctly identified the red pixel, image handling is working!")
}
