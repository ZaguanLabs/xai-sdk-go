// Copyright 2024 Zaguan Labs, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	apiKey := os.Getenv("XAI_API_KEY")
	if apiKey == "" {
		return fmt.Errorf("XAI_API_KEY environment variable not set")
	}

	client, err := xai.NewClient(&xai.Config{APIKey: apiKey})
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}
	defer client.Close()

	ctx := context.Background()

	// Example 1: Using default (low) reasoning effort.
	fmt.Println("--- Example 1: Default (low) reasoning effort ---")
	reqLow := chat.NewRequest("grok-1.5-flash")
	reqLow.AppendMessage(*chat.User(chat.Text("What is the speed of light?")))
	respLow, err := reqLow.Sample(ctx, client.Chat())
	if err != nil {
		return fmt.Errorf("sampling with low reasoning effort failed: %w", err)
	}
	fmt.Printf("Response (low effort): %s\n", respLow.Content())

	// Example 2: Using high reasoning effort for a more complex query.
	fmt.Println("\n--- Example 2: High reasoning effort for a complex query ---")
	reqHigh := chat.NewRequest("grok-1.5-flash", chat.WithReasoningEffort(chat.ReasoningEffortHigh))
	reqHigh.AppendMessage(*chat.User(chat.Text("Explain the theory of relativity in simple terms.")))
	respHigh, err := reqHigh.Sample(ctx, client.Chat())
	if err != nil {
		return fmt.Errorf("sampling with high reasoning effort failed: %w", err)
	}
	fmt.Printf("Response (high effort): %s\n", respHigh.Content())

	return nil
}