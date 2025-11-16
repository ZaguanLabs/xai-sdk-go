package chat

import (
	"testing"

	xaiv1 "github.com/ZaguanLabs/xai-sdk-go/proto/gen/go/xai/v1"
)

func TestTextPart(t *testing.T) {
	text := Text("Hello, world!")

	if text.Content() != "Hello, world!" {
		t.Errorf("Content() = %q, want %q", text.Content(), "Hello, world!")
	}

	if text.Type() != PartTypeText {
		t.Errorf("Type() = %v, want %v", text.Type(), PartTypeText)
	}

	t.Log("✅ TextPart works correctly")
}

func TestImagePart(t *testing.T) {
	tests := []struct {
		name           string
		url            string
		detail         []ImageDetail
		expectedDetail xaiv1.ImageDetail
	}{
		{
			name:           "default detail (auto)",
			url:            "https://example.com/image.jpg",
			detail:         nil,
			expectedDetail: xaiv1.ImageDetail_DETAIL_AUTO,
		},
		{
			name:           "low detail",
			url:            "https://example.com/image.jpg",
			detail:         []ImageDetail{ImageDetailLow},
			expectedDetail: xaiv1.ImageDetail_DETAIL_LOW,
		},
		{
			name:           "high detail",
			url:            "https://example.com/image.jpg",
			detail:         []ImageDetail{ImageDetailHigh},
			expectedDetail: xaiv1.ImageDetail_DETAIL_HIGH,
		},
		{
			name:           "base64 image",
			url:            "data:image/jpeg;base64,/9j/4AAQSkZJRg...",
			detail:         []ImageDetail{ImageDetailHigh},
			expectedDetail: xaiv1.ImageDetail_DETAIL_HIGH,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var img Part
			if tt.detail == nil {
				img = Image(tt.url)
			} else {
				img = Image(tt.url, tt.detail[0])
			}

			if img.Content() != tt.url {
				t.Errorf("Content() = %q, want %q", img.Content(), tt.url)
			}

			if img.Type() != PartTypeImage {
				t.Errorf("Type() = %v, want %v", img.Type(), PartTypeImage)
			}

			imgPart, ok := img.(*ImagePart)
			if !ok {
				t.Fatal("Image() did not return *ImagePart")
			}

			if imgPart.ImageURL() != tt.url {
				t.Errorf("ImageURL() = %q, want %q", imgPart.ImageURL(), tt.url)
			}

			if imgPart.Detail() != tt.expectedDetail {
				t.Errorf("Detail() = %v, want %v", imgPart.Detail(), tt.expectedDetail)
			}
		})
	}

	t.Log("✅ ImagePart works correctly with all detail levels")
}

func TestFilePart(t *testing.T) {
	fileID := "file-abc123"
	file := File(fileID)

	if file.Content() != fileID {
		t.Errorf("Content() = %q, want %q", file.Content(), fileID)
	}

	if file.Type() != PartTypeFile {
		t.Errorf("Type() = %v, want %v", file.Type(), PartTypeFile)
	}

	filePart, ok := file.(*FilePart)
	if !ok {
		t.Fatal("File() did not return *FilePart")
	}

	if filePart.FileID() != fileID {
		t.Errorf("FileID() = %q, want %q", filePart.FileID(), fileID)
	}

	t.Log("✅ FilePart works correctly")
}

func TestMessageWithImage(t *testing.T) {
	msg := User(
		Text("What's in this image?"),
		Image("https://example.com/photo.jpg", ImageDetailHigh),
	)

	// Verify message structure
	if msg.Role() != "user" {
		t.Errorf("Role() = %q, want %q", msg.Role(), "user")
	}

	proto := msg.Proto()
	if len(proto.Content) != 2 {
		t.Fatalf("Expected 2 content parts, got %d", len(proto.Content))
	}

	// First part should be text
	if proto.Content[0].Text != "What's in this image?" {
		t.Errorf("Content[0].Text = %q, want %q", proto.Content[0].Text, "What's in this image?")
	}

	// Second part should be image
	if proto.Content[1].ImageUrl == nil {
		t.Fatal("Content[1].ImageUrl is nil")
	}

	if proto.Content[1].ImageUrl.ImageUrl != "https://example.com/photo.jpg" {
		t.Errorf("ImageUrl = %q, want %q", proto.Content[1].ImageUrl.ImageUrl, "https://example.com/photo.jpg")
	}

	if proto.Content[1].ImageUrl.Detail != xaiv1.ImageDetail_DETAIL_HIGH {
		t.Errorf("Detail = %v, want %v", proto.Content[1].ImageUrl.Detail, xaiv1.ImageDetail_DETAIL_HIGH)
	}

	t.Log("✅ Message with image content works correctly")
}

func TestMessageWithMultipleImages(t *testing.T) {
	msg := User(
		Text("Compare these images:"),
		Image("https://example.com/photo1.jpg"),
		Image("https://example.com/photo2.jpg", ImageDetailLow),
	)

	proto := msg.Proto()
	if len(proto.Content) != 3 {
		t.Fatalf("Expected 3 content parts, got %d", len(proto.Content))
	}

	// First part: text
	if proto.Content[0].Text != "Compare these images:" {
		t.Errorf("Content[0].Text = %q, want %q", proto.Content[0].Text, "Compare these images:")
	}

	// Second part: first image (auto detail)
	if proto.Content[1].ImageUrl == nil {
		t.Fatal("Content[1].ImageUrl is nil")
	}
	if proto.Content[1].ImageUrl.ImageUrl != "https://example.com/photo1.jpg" {
		t.Errorf("ImageUrl[1] = %q, want %q", proto.Content[1].ImageUrl.ImageUrl, "https://example.com/photo1.jpg")
	}
	if proto.Content[1].ImageUrl.Detail != xaiv1.ImageDetail_DETAIL_AUTO {
		t.Errorf("Detail[1] = %v, want %v", proto.Content[1].ImageUrl.Detail, xaiv1.ImageDetail_DETAIL_AUTO)
	}

	// Third part: second image (low detail)
	if proto.Content[2].ImageUrl == nil {
		t.Fatal("Content[2].ImageUrl is nil")
	}
	if proto.Content[2].ImageUrl.ImageUrl != "https://example.com/photo2.jpg" {
		t.Errorf("ImageUrl[2] = %q, want %q", proto.Content[2].ImageUrl.ImageUrl, "https://example.com/photo2.jpg")
	}
	if proto.Content[2].ImageUrl.Detail != xaiv1.ImageDetail_DETAIL_LOW {
		t.Errorf("Detail[2] = %v, want %v", proto.Content[2].ImageUrl.Detail, xaiv1.ImageDetail_DETAIL_LOW)
	}

	t.Log("✅ Message with multiple images works correctly")
}

func TestMessageWithFile(t *testing.T) {
	msg := User(
		Text("Analyze this document:"),
		File("file-xyz789"),
	)

	proto := msg.Proto()
	if len(proto.Content) != 2 {
		t.Fatalf("Expected 2 content parts, got %d", len(proto.Content))
	}

	// First part: text
	if proto.Content[0].Text != "Analyze this document:" {
		t.Errorf("Content[0].Text = %q, want %q", proto.Content[0].Text, "Analyze this document:")
	}

	// Second part: file
	if proto.Content[1].File == nil {
		t.Fatal("Content[1].File is nil")
	}

	if proto.Content[1].File.FileId != "file-xyz789" {
		t.Errorf("FileId = %q, want %q", proto.Content[1].File.FileId, "file-xyz789")
	}

	t.Log("✅ Message with file content works correctly")
}

func TestMessageWithMixedContent(t *testing.T) {
	msg := User(
		Text("Here's my question:"),
		Image("https://example.com/chart.png", ImageDetailHigh),
		Text("What does this chart show?"),
		File("file-data123"),
		Text("And compare it with this file."),
	)

	proto := msg.Proto()
	if len(proto.Content) != 5 {
		t.Fatalf("Expected 5 content parts, got %d", len(proto.Content))
	}

	// Verify each part type
	if proto.Content[0].Text != "Here's my question:" {
		t.Error("Content[0] should be text")
	}

	if proto.Content[1].ImageUrl == nil {
		t.Error("Content[1] should be image")
	}

	if proto.Content[2].Text != "What does this chart show?" {
		t.Error("Content[2] should be text")
	}

	if proto.Content[3].File == nil {
		t.Error("Content[3] should be file")
	}

	if proto.Content[4].Text != "And compare it with this file." {
		t.Error("Content[4] should be text")
	}

	t.Log("✅ Message with mixed content types works correctly")
}

func TestBase64Image(t *testing.T) {
	base64Data := "data:image/jpeg;base64,/9j/4AAQSkZJRgABAQAAAQABAAD..."

	msg := User(
		Text("Analyze this base64 image:"),
		Image(base64Data, ImageDetailHigh),
	)

	proto := msg.Proto()
	if len(proto.Content) != 2 {
		t.Fatalf("Expected 2 content parts, got %d", len(proto.Content))
	}

	if proto.Content[1].ImageUrl == nil {
		t.Fatal("Content[1].ImageUrl is nil")
	}

	if proto.Content[1].ImageUrl.ImageUrl != base64Data {
		t.Errorf("Base64 image URL not preserved correctly")
	}

	t.Log("✅ Base64 images work correctly")
}
