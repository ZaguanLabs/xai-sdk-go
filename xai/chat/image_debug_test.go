package chat

import (
	"encoding/json"
	"testing"

	xaiv1 "github.com/ZaguanLabs/xai-sdk-go/proto/gen/go/xai/v1"
	"google.golang.org/protobuf/encoding/protojson"
)

// TestImageProtoSerialization verifies that images are correctly serialized to proto
func TestImageProtoSerialization(t *testing.T) {
	// Create a message with a base64 image (like OpenWebUI would send)
	base64Image := "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+M9QDwADhgGAWjR9awAAAABJRU5ErkJggg=="

	msg := User(
		Text("What's in this image?"),
		Image(base64Image, ImageDetailHigh),
	)

	// Verify the proto structure
	if msg.proto == nil {
		t.Fatal("Message proto is nil")
	}

	if len(msg.proto.Content) != 2 {
		t.Fatalf("Expected 2 content parts, got %d", len(msg.proto.Content))
	}

	// Verify text content
	textContent := msg.proto.Content[0]
	if textContent.Text != "What's in this image?" {
		t.Errorf("Text content mismatch: got %q", textContent.Text)
	}
	if textContent.ImageUrl != nil {
		t.Error("Text content should not have ImageUrl")
	}
	if textContent.File != nil {
		t.Error("Text content should not have File")
	}

	// Verify image content
	imageContent := msg.proto.Content[1]
	if imageContent.Text != "" {
		t.Errorf("Image content should not have text, got %q", imageContent.Text)
	}
	if imageContent.File != nil {
		t.Error("Image content should not have File")
	}
	if imageContent.ImageUrl == nil {
		t.Fatal("Image content should have ImageUrl")
	}
	if imageContent.ImageUrl.ImageUrl != base64Image {
		t.Errorf("ImageUrl mismatch:\nExpected: %q\nGot: %q", base64Image, imageContent.ImageUrl.ImageUrl)
	}
	if imageContent.ImageUrl.Detail != xaiv1.ImageDetail_DETAIL_HIGH {
		t.Errorf("Expected DETAIL_HIGH, got %v", imageContent.ImageUrl.Detail)
	}

	// Serialize to JSON to see what gets sent
	jsonBytes, err := protojson.Marshal(msg.proto)
	if err != nil {
		t.Fatalf("Failed to marshal to JSON: %v", err)
	}

	t.Logf("Proto JSON representation:\n%s", string(jsonBytes))

	// Verify JSON structure
	var jsonData map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &jsonData); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	content, ok := jsonData["content"].([]interface{})
	if !ok || len(content) != 2 {
		t.Fatalf("Expected content array with 2 elements, got %v", jsonData["content"])
	}

	// Check image content in JSON
	imageJSON, ok := content[1].(map[string]interface{})
	if !ok {
		t.Fatalf("Second content element should be a map, got %T", content[1])
	}

	imageUrlData, ok := imageJSON["imageUrl"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected imageUrl field, got %v", imageJSON)
	}

	if imageUrlData["imageUrl"] != base64Image {
		t.Errorf("JSON imageUrl mismatch:\nExpected: %q\nGot: %q", base64Image, imageUrlData["imageUrl"])
	}
}

// TestRequestWithImageProto verifies the full request proto structure
func TestRequestWithImageProto(t *testing.T) {
	base64Image := "data:image/jpeg;base64,/9j/4AAQSkZJRgABAQEAYABgAAD/2wBDAAgGBgcGBQgHBwcJCQgKDBQNDAsLDBkSEw8UHRofHh0aHBwgJC4nICIsIxwcKDcpLDAxNDQ0Hyc5PTgyPC4zNDL/2wBDAQkJCQwLDBgNDRgyIRwhMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjL/wAARCAABAAEDASIAAhEBAxEB/8QAFQABAQAAAAAAAAAAAAAAAAAAAAv/xAAUEAEAAAAAAAAAAAAAAAAAAAAA/8QAFQEBAQAAAAAAAAAAAAAAAAAAAAX/xAAUEQEAAAAAAAAAAAAAAAAAAAAA/9oADAMBAAIRAxEAPwCwAA8A/9k="

	req := NewRequest("grok-2-vision",
		WithMessages(
			User(
				Text("Describe this image"),
				Image(base64Image, ImageDetailHigh),
			),
		),
		WithMaxTokens(100),
	)

	proto := req.Proto()
	if proto == nil {
		t.Fatal("Request proto is nil")
	}

	if len(proto.Messages) != 1 {
		t.Fatalf("Expected 1 message, got %d", len(proto.Messages))
	}

	msg := proto.Messages[0]
	if len(msg.Content) != 2 {
		t.Fatalf("Expected 2 content parts, got %d", len(msg.Content))
	}

	// Verify the image content
	imageContent := msg.Content[1]
	if imageContent.ImageUrl == nil {
		t.Fatal("ImageUrl is nil in request proto")
	}
	if imageContent.ImageUrl.ImageUrl != base64Image {
		t.Error("Base64 image not preserved in request proto")
	}

	// Serialize the full request to JSON
	jsonBytes, err := protojson.Marshal(proto)
	if err != nil {
		t.Fatalf("Failed to marshal request to JSON: %v", err)
	}

	t.Logf("Full request proto JSON:\n%s", string(jsonBytes))

	// Verify the JSON contains the image
	var jsonData map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &jsonData); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	messages, ok := jsonData["messages"].([]interface{})
	if !ok || len(messages) != 1 {
		t.Fatalf("Expected messages array with 1 element")
	}

	msgJSON, ok := messages[0].(map[string]interface{})
	if !ok {
		t.Fatal("Message should be a map")
	}

	content, ok := msgJSON["content"].([]interface{})
	if !ok || len(content) != 2 {
		t.Fatalf("Expected content array with 2 elements")
	}

	imageJSON, ok := content[1].(map[string]interface{})
	if !ok {
		t.Fatal("Second content element should be a map")
	}

	if _, hasImageUrl := imageJSON["imageUrl"]; !hasImageUrl {
		t.Errorf("JSON missing imageUrl field. Got: %v", imageJSON)
	}
}
