package chat

import (
	"testing"

	xaiv1 "github.com/ZaguanLabs/xai-sdk-go/proto/gen/go/xai/v1"
)

func TestImageBase64Encoding(t *testing.T) {
	// Test with a base64 data URI (simulating OpenWebUI)
	base64Image := "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+M9QDwADhgGAWjR9awAAAABJRU5ErkJggg=="

	msg := User(
		Text("What's in this image?"),
		Image(base64Image),
	)

	if msg.proto == nil {
		t.Fatal("Message proto is nil")
	}

	if len(msg.proto.Content) != 2 {
		t.Fatalf("Expected 2 content parts, got %d", len(msg.proto.Content))
	}

	// First part should be text
	if msg.proto.Content[0].Text != "What's in this image?" {
		t.Errorf("Expected text content, got %q", msg.proto.Content[0].Text)
	}

	// Second part should be image
	if msg.proto.Content[1].ImageUrl == nil {
		t.Fatal("ImageUrl is nil")
	}

	if msg.proto.Content[1].ImageUrl.ImageUrl != base64Image {
		t.Errorf("Expected base64 image URL %q, got %q", base64Image, msg.proto.Content[1].ImageUrl.ImageUrl)
	}

	if msg.proto.Content[1].ImageUrl.Detail != xaiv1.ImageDetail_DETAIL_AUTO {
		t.Errorf("Expected DETAIL_AUTO, got %v", msg.proto.Content[1].ImageUrl.Detail)
	}
}

func TestImageWithHighDetail(t *testing.T) {
	imageURL := "https://example.com/image.png"

	msg := User(
		Image(imageURL, ImageDetailHigh),
	)

	if msg.proto == nil {
		t.Fatal("Message proto is nil")
	}

	if len(msg.proto.Content) != 1 {
		t.Fatalf("Expected 1 content part, got %d", len(msg.proto.Content))
	}

	if msg.proto.Content[0].ImageUrl == nil {
		t.Fatal("ImageUrl is nil")
	}

	if msg.proto.Content[0].ImageUrl.ImageUrl != imageURL {
		t.Errorf("Expected image URL %q, got %q", imageURL, msg.proto.Content[0].ImageUrl.ImageUrl)
	}

	if msg.proto.Content[0].ImageUrl.Detail != xaiv1.ImageDetail_DETAIL_HIGH {
		t.Errorf("Expected DETAIL_HIGH, got %v", msg.proto.Content[0].ImageUrl.Detail)
	}
}

func TestMultipleImages(t *testing.T) {
	img1 := "https://example.com/image1.png"
	img2 := "data:image/jpeg;base64,/9j/4AAQSkZJRgABAQEAYABgAAD/2wBDAAgGBgcGBQgHBwcJCQgKDBQNDAsLDBkSEw8UHRofHh0aHBwgJC4nICIsIxwcKDcpLDAxNDQ0Hyc5PTgyPC4zNDL/2wBDAQkJCQwLDBgNDRgyIRwhMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjL/wAARCAABAAEDASIAAhEBAxEB/8QAFQABAQAAAAAAAAAAAAAAAAAAAAv/xAAUEAEAAAAAAAAAAAAAAAAAAAAA/8QAFQEBAQAAAAAAAAAAAAAAAAAAAAX/xAAUEQEAAAAAAAAAAAAAAAAAAAAA/9oADAMBAAIRAxEAPwCwAA8A/9k="

	msg := User(
		Text("Compare these images"),
		Image(img1, ImageDetailLow),
		Image(img2, ImageDetailHigh),
	)

	if msg.proto == nil {
		t.Fatal("Message proto is nil")
	}

	if len(msg.proto.Content) != 3 {
		t.Fatalf("Expected 3 content parts, got %d", len(msg.proto.Content))
	}

	// Check text
	if msg.proto.Content[0].Text != "Compare these images" {
		t.Errorf("Expected text content, got %q", msg.proto.Content[0].Text)
	}

	// Check first image (URL with low detail)
	if msg.proto.Content[1].ImageUrl == nil {
		t.Fatal("First ImageUrl is nil")
	}
	if msg.proto.Content[1].ImageUrl.ImageUrl != img1 {
		t.Errorf("Expected first image URL %q, got %q", img1, msg.proto.Content[1].ImageUrl.ImageUrl)
	}
	if msg.proto.Content[1].ImageUrl.Detail != xaiv1.ImageDetail_DETAIL_LOW {
		t.Errorf("Expected DETAIL_LOW for first image, got %v", msg.proto.Content[1].ImageUrl.Detail)
	}

	// Check second image (base64 with high detail)
	if msg.proto.Content[2].ImageUrl == nil {
		t.Fatal("Second ImageUrl is nil")
	}
	if msg.proto.Content[2].ImageUrl.ImageUrl != img2 {
		t.Errorf("Expected second image URL %q, got %q", img2, msg.proto.Content[2].ImageUrl.ImageUrl)
	}
	if msg.proto.Content[2].ImageUrl.Detail != xaiv1.ImageDetail_DETAIL_HIGH {
		t.Errorf("Expected DETAIL_HIGH for second image, got %v", msg.proto.Content[2].ImageUrl.Detail)
	}
}
