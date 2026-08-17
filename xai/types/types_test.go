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

func TestPython1180ModelCatalogParity(t *testing.T) {
	wantChat := []string{
		"grok-4", "grok-4-0709", "grok-4-latest", "grok-4-1-fast",
		"grok-4-1-fast-reasoning", "grok-4-1-fast-reasoning-latest",
		"grok-4-1-fast-non-reasoning", "grok-4-1-fast-non-reasoning-latest",
		"grok-4-fast", "grok-4-fast-reasoning", "grok-4-fast-reasoning-latest",
		"grok-4-fast-non-reasoning", "grok-4-fast-non-reasoning-latest",
		"grok-4.20-0309-reasoning", "grok-4.20", "grok-4.20-0309",
		"grok-4.20-reasoning-latest", "grok-4.20-0309-non-reasoning",
		"grok-4.20-non-reasoning", "grok-4.20-non-reasoning-latest",
		"grok-4.20-multi-agent", "grok-4.20-multi-agent-0309", "grok-4.20-multi-agent-latest",
		"grok-4.3", "grok-4.3-latest", "grok-4.5", "grok-4.5-latest", "grok-4.6",
		"grok-code-fast-1", "grok-build-0.1", "grok-3", "grok-3-latest", "grok-3-mini",
		"grok-3-fast", "grok-3-fast-latest", "grok-3-mini-fast", "grok-3-mini-fast-latest",
	}
	if len(ChatModels) != len(wantChat) {
		t.Fatalf("ChatModels length = %d, want %d", len(ChatModels), len(wantChat))
	}
	for i := range wantChat {
		if ChatModels[i] != wantChat[i] {
			t.Fatalf("ChatModels[%d] = %q, want %q", i, ChatModels[i], wantChat[i])
		}
	}

	wantImages := []string{"grok-imagine-image", "grok-imagine-image-pro", "grok-imagine-image-quality"}
	for i := range wantImages {
		if ImageGenerationModels[i] != wantImages[i] {
			t.Fatalf("ImageGenerationModels[%d] = %q, want %q", i, ImageGenerationModels[i], wantImages[i])
		}
	}
	wantVideos := []string{"grok-imagine-video", "grok-imagine-video-1.5-preview"}
	for i := range wantVideos {
		if VideoGenerationModels[i] != wantVideos[i] {
			t.Fatalf("VideoGenerationModels[%d] = %q, want %q", i, VideoGenerationModels[i], wantVideos[i])
		}
	}
}

func TestTypeConstants(t *testing.T) {
	if ImageAspectRatio16x9 != "16:9" || VideoResolution720P != "720p" || IncludeOptionVerboseStreaming != "verbose_streaming" {
		t.Fatal("unexpected type constants")
	}
}
