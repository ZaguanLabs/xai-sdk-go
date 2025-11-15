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

package chat_test

import (
	"context"
	"os"
	"testing"

	"github.com/ZaguanLabs/xai-sdk-go/xai"
	"github.com/ZaguanLabs/xai-sdk-go/xai/chat"
)

func TestChat_Sample_Integration(t *testing.T) {
	if os.Getenv("XAI_SDK_E2E") == "" {
		t.Skip("skipping integration test; set XAI_SDK_E2E to run")
	}
	apiKey := os.Getenv("XAI_API_KEY")
	if apiKey == "" {
		t.Fatal("XAI_API_KEY is not set")
	}

	client, err := xai.NewClient(&xai.Config{APIKey: apiKey})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	defer client.Close()

	req := chat.NewRequest("grok-1.5-flash",
		chat.WithMessages(
			chat.System(chat.Text("You are a helpful assistant.")),
			chat.User(chat.Text("What is the capital of France?")),
		),
	)

	resp, err := req.Sample(context.Background(), client.Chat())
	if err != nil {
		t.Fatalf("failed to sample chat: %v", err)
	}

	if resp.Content() == "" {
		t.Error("expected content to be non-empty")
	}
}

func TestChat_Stream_Integration(t *testing.T) {
	if os.Getenv("XAI_SDK_E2E") == "" {
		t.Skip("skipping integration test; set XAI_SDK_E2E to run")
	}
	apiKey := os.Getenv("XAI_API_KEY")
	if apiKey == "" {
		t.Fatal("XAI_API_KEY is not set")
	}

	client, err := xai.NewClient(&xai.Config{APIKey: apiKey})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	defer client.Close()

	req := chat.NewRequest("grok-1.5-flash",
		chat.WithMessages(
			chat.System(chat.Text("You are a helpful assistant.")),
			chat.User(chat.Text("Tell me a story about a brave knight in 100 words.")),
		),
	)

	stream, err := req.Stream(context.Background(), client.Chat())
	if err != nil {
		t.Fatalf("failed to stream chat: %v", err)
	}
	defer stream.Close()

	var content string
	for stream.Next() {
		chunk := stream.Current()
		content += chunk.Content()
	}

	if stream.Err() != nil {
		t.Fatalf("streaming error: %v", stream.Err())
	}

	if content == "" {
		t.Error("expected content to be non-empty")
	}
}
