package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/ZaguanLabs/xai-sdk-go/xai"
	"github.com/ZaguanLabs/xai-sdk-go/xai/chat"
)

// getCurrentWeather is a sample function that our tool will call
func getCurrentWeather(location string, unit string) (string, error) {
	// In a real implementation, you would call a weather API here
	// For this example, we'll return mock data
	if location == "San Francisco" {
		return fmt.Sprintf(`{"temperature": "72", "unit": "%s", "description": "Sunny"}`, unit), nil
	}
	if location == "New York" {
		return fmt.Sprintf(`{"temperature": "45", "unit": "%s", "description": "Cloudy"}`, unit), nil
	}
	return fmt.Sprintf(`{"temperature": "65", "unit": "%s", "description": "Unknown location"}`, unit), nil
}

func main() {
	// Get API key from environment variable
	apiKey := os.Getenv("XAI_API_KEY")
	if apiKey == "" {
		log.Fatal("XAI_API_KEY environment variable is required")
	}

	// Create a new xAI client
	client, err := xai.NewClient(&xai.Config{APIKey: apiKey})
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	// Create a weather tool
	weatherTool := chat.NewTool(
		"get_current_weather",
		"Get current weather in a given location",
	).
		WithParameter("location", "string", "The city and state, e.g. San Francisco, CA", true).
		WithParameter("unit", "string", "The unit of temperature, e.g. celsius or fahrenheit", false)

	// Create a chat request with tool calling
	request := chat.NewRequest("grok-1.5-flash",
		chat.WithMessages(
			chat.User(chat.Text("What's the weather like in San Francisco?")),
		),
		chat.WithTool(weatherTool),
		chat.WithToolChoice(chat.ToolChoiceAuto),
	)

	// Perform the request
	response, err := request.Sample(context.Background(), client.Chat())
	if err != nil {
		log.Fatalf("Chat completion failed: %v", err)
	}

	// Check if there are tool calls in the response
	toolCalls := response.ToolCalls()
	if len(toolCalls) > 0 {
		// Handle the tool calls
		toolResults := make([]chat.ToolResult, 0)

		for _, toolCall := range toolCalls {
			if toolCall.Function().Name() == "get_current_weather" {
				// Parse the arguments
				args := toolCall.Function().Arguments()
				location, _ := args["location"].(string)
				unit, _ := args["unit"].(string)
				if unit == "" {
					unit = "fahrenheit" // default
				}

				// Call the actual function
				result, err := getCurrentWeather(location, unit)
				if err != nil {
					log.Printf("Function call failed: %v", err)
					toolResults = append(toolResults, *chat.NewToolResult(toolCall.ID(), fmt.Sprintf("Error: %v", err)))
				} else {
					toolResults = append(toolResults, *chat.NewToolResult(toolCall.ID(), result))
				}
			}
		}

		// Create a follow-up request with tool results
		followupRequest := chat.NewRequest("grok-1.5-flash",
			chat.WithMessages(
				chat.User(chat.Text("What's the weather like in San Francisco?")),
				chat.Assistant(chat.Text(response.Content())),
			),
			chat.WithToolResults(toolResults...),
			chat.WithToolChoice(chat.ToolChoiceAuto),
		)

		// Get the final response
		finalResponse, err := followupRequest.Sample(context.Background(), client.Chat())
		if err != nil {
			log.Fatalf("Follow-up completion failed: %v", err)
		}

		fmt.Printf("Assistant: %s\n", finalResponse.Content())
	} else {
		// No tool calls, just a regular response
		fmt.Printf("Assistant: %s\n", response.Content())
	}
}