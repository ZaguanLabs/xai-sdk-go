package types

import "testing"

func TestModelConstants(t *testing.T) {
	if ChatModelGrok420 != "grok-4.20" {
		t.Fatalf("ChatModelGrok420 = %q", ChatModelGrok420)
	}
	if ImageGenerationModelGrokImagineImagePro != "grok-imagine-image-pro" {
		t.Fatalf("ImageGenerationModelGrokImagineImagePro = %q", ImageGenerationModelGrokImagineImagePro)
	}
	if VideoGenerationModelGrokImagineVideo != "grok-imagine-video" {
		t.Fatalf("VideoGenerationModelGrokImagineVideo = %q", VideoGenerationModelGrokImagineVideo)
	}
	if len(ChatModels) == 0 || len(ImageGenerationModels) == 0 || len(VideoGenerationModels) == 0 {
		t.Fatal("model slices should not be empty")
	}
}

func TestTypeConstants(t *testing.T) {
	if ImageAspectRatio16x9 != "16:9" || VideoResolution720P != "720p" || IncludeOptionVerboseStreaming != "verbose_streaming" {
		t.Fatal("unexpected type constants")
	}
}
