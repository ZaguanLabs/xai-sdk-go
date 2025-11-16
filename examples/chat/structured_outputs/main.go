package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/ZaguanLabs/xai-sdk-go/xai"
	"github.com/ZaguanLabs/xai-sdk-go/xai/chat"
)

// WeatherResponse represents the structured weather data we want
type WeatherResponse struct {
	Location    string `json:"location"`
	Temperature string `json:"temperature"`
	Unit        string `json:"unit"`
	Description string `json:"description"`
	Conditions  string `json:"conditions"`
	Humidity    int    `json:"humidity"`
}

// TaskResponse represents a structured task
type TaskResponse struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Priority    string `json:"priority"`
	DueDate     string `json:"due_date"`
	Completed   bool   `json:"completed"`
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

	// Example 1: JSON object response format
	fmt.Println("=== Example 1: JSON Object Format ===")

	jsonRequest := chat.NewRequest("grok-1.5-flash",
		chat.WithMessages(
			chat.User(chat.Text("What's the weather like in Tokyo? Provide the response in JSON format with fields: location, temperature, unit, description, conditions, humidity.")),
		),
		chat.WithResponseFormat(chat.ResponseFormatJSONObject),
		chat.WithTemperature(0.1),
		chat.WithMaxTokens(500),
	)

	var weatherData WeatherResponse
	err = jsonRequest.Parse(context.Background(), client.Chat(), &weatherData)
	if err != nil {
		log.Fatalf("Failed to parse weather data: %v", err)
	}

	fmt.Printf("Weather Data: %+v\n", weatherData)
	fmt.Printf("Temperature: %s %s\n", weatherData.Temperature, weatherData.Unit)
	fmt.Printf("Location: %s\n", weatherData.Location)
	fmt.Printf("Description: %s\n", weatherData.Description)
	fmt.Println()

	// Example 2: JSON schema response format
	fmt.Println("=== Example 2: JSON Schema Format ===")

	schema := map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"id": map[string]interface{}{
				"type":        "string",
				"description": "Unique identifier for the task",
			},
			"title": map[string]interface{}{
				"type":        "string",
				"description": "Title of the task",
			},
			"description": map[string]interface{}{
				"type":        "string",
				"description": "Detailed description of what needs to be done",
			},
			"priority": map[string]interface{}{
				"type":        "string",
				"enum":        []string{"low", "medium", "high"},
				"description": "Task priority level",
			},
			"due_date": map[string]interface{}{
				"type":        "string",
				"description": "When the task should be completed (ISO format)",
			},
			"completed": map[string]interface{}{
				"type":        "boolean",
				"description": "Whether the task has been completed",
			},
		},
		"required": []string{"id", "title", "description"},
	}

	schemaRequest := chat.NewRequest("grok-1.5-flash",
		chat.WithMessages(
			chat.User(chat.Text("Create a task for 'Review quarterly report' with medium priority due next week. Return the response in the specified JSON schema format.")),
		),
		chat.WithResponseFormatOption(&chat.ResponseFormatOption{
			Type:   chat.ResponseFormatJSONSchema,
			Schema: schema,
		}),
		chat.WithTemperature(0.1),
		chat.WithMaxTokens(500),
	)

	var taskData TaskResponse
	err = schemaRequest.Parse(context.Background(), client.Chat(), &taskData)
	if err != nil {
		log.Fatalf("Failed to parse task data: %v", err)
	}

	fmt.Printf("Task Data: %+v\n", taskData)
	fmt.Printf("Task ID: %s\n", taskData.ID)
	fmt.Printf("Title: %s\n", taskData.Title)
	fmt.Printf("Priority: %s\n", taskData.Priority)
	fmt.Printf("Due Date: %s\n", taskData.DueDate)
	fmt.Printf("Completed: %t\n", taskData.Completed)
	fmt.Println()

	// Example 3: Parse to map interface
	fmt.Println("=== Example 3: Parse to Map Interface ===")

	mapRequest := chat.NewRequest("grok-1.5-flash",
		chat.WithMessages(
			chat.User(chat.Text("List 3 programming languages with their key features. Return as a simple map structure.")),
		),
		chat.WithResponseFormat(chat.ResponseFormatJSONObject),
		chat.WithTemperature(0.1),
		chat.WithMaxTokens(300),
	)

	var languages map[string]interface{}
	err = mapRequest.Parse(context.Background(), client.Chat(), &languages)
	if err != nil {
		log.Fatalf("Failed to parse languages data: %v", err)
	}

	fmt.Printf("Languages: %+v\n", languages)
	for name, features := range languages {
		if featuresMap, ok := features.(map[string]interface{}); ok {
			fmt.Printf("- %s: %v\n", name, featuresMap["features"])
		}
	}
	fmt.Println("=== Structured Outputs Examples Complete ===")
}
