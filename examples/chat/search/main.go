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

	// Example: Using search with a query.
	fmt.Println("--- Example: Using web search with a query ---")
	req := chat.NewRequest("grok-1.5-flash",
		chat.WithSearch(
			chat.NewSearchParameters().WithCount(10),
		),
	)
	req.AppendMessage(*chat.User(chat.Text("What are the latest developments in AI?")))
	resp, err := req.Sample(ctx, client.Chat())
	if err != nil {
		return fmt.Errorf("sampling with search failed: %w", err)
	}

	fmt.Printf("Response with search results:\n%s\n", resp.Content())

	return nil
}