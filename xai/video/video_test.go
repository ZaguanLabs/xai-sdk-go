package video

import (
	"testing"

	xaiv1 "github.com/ZaguanLabs/xai-sdk-go/proto/gen/go/xai/api/v1"
)

func TestNewGenerateRequestWithOptions(t *testing.T) {
	duration := int32(8)
	aspectRatio := xaiv1.VideoAspectRatio_VIDEO_ASPECT_RATIO_16_9
	resolution := xaiv1.VideoResolution_VIDEO_RESOLUTION_720P

	req := NewGenerateRequestWithOptions("prompt", "model", &GenerateOptions{
		ImageURL:           "https://example.com/image.jpg",
		VideoURL:           "https://example.com/video.mp4",
		ReferenceImageURLs: []string{"https://example.com/ref1.jpg", "https://example.com/ref2.jpg"},
		Duration:           &duration,
		AspectRatio:        &aspectRatio,
		Resolution:         &resolution,
	})

	if req.Prompt != "prompt" || req.Model != "model" {
		t.Fatalf("unexpected request identity: %q %q", req.Prompt, req.Model)
	}
	if req.GetImage().GetImageUrl() != "https://example.com/image.jpg" {
		t.Errorf("image url = %q", req.GetImage().GetImageUrl())
	}
	if req.GetVideo().GetUrl() != "https://example.com/video.mp4" {
		t.Errorf("video url = %q", req.GetVideo().GetUrl())
	}
	if len(req.ReferenceImages) != 2 {
		t.Fatalf("reference image count = %d, want 2", len(req.ReferenceImages))
	}
	if req.GetDuration() != duration {
		t.Errorf("duration = %d, want %d", req.GetDuration(), duration)
	}
	if req.GetAspectRatio() != aspectRatio {
		t.Errorf("aspect ratio = %v, want %v", req.GetAspectRatio(), aspectRatio)
	}
	if req.GetResolution() != resolution {
		t.Errorf("resolution = %v, want %v", req.GetResolution(), resolution)
	}
}

func TestNewExtendRequest(t *testing.T) {
	duration := int32(6)
	req := NewExtendRequest("continue", "model", "https://example.com/video.mp4", &duration)

	if req.Prompt != "continue" || req.Model != "model" {
		t.Fatalf("unexpected request identity: %q %q", req.Prompt, req.Model)
	}
	if req.GetVideo().GetUrl() != "https://example.com/video.mp4" {
		t.Errorf("video url = %q", req.GetVideo().GetUrl())
	}
	if req.GetDuration() != duration {
		t.Errorf("duration = %d, want %d", req.GetDuration(), duration)
	}
}

func TestPrepare(t *testing.T) {
	batchReq := Prepare("prompt", "model", "request-1", nil)

	if batchReq.GetBatchRequestId() != "request-1" {
		t.Errorf("batch request id = %q, want request-1", batchReq.GetBatchRequestId())
	}
	if batchReq.GetVideoRequest() == nil {
		t.Fatal("video request is nil")
	}
	if batchReq.GetVideoRequest().Model != "model" {
		t.Errorf("model = %q, want model", batchReq.GetVideoRequest().Model)
	}
}

func TestResponse(t *testing.T) {
	ticks := int64(100)
	resp := NewResponse(&xaiv1.VideoResponse{
		Model: "video-model",
		Usage: &xaiv1.SamplingUsage{CostInUsdTicks: &ticks},
		Video: &xaiv1.GeneratedVideo{
			Url:               "https://example.com/video.mp4",
			Duration:          6,
			RespectModeration: true,
		},
	})

	if resp.Model() != "video-model" {
		t.Errorf("Model() = %q, want video-model", resp.Model())
	}
	url, err := resp.URL()
	if err != nil {
		t.Fatalf("URL() error = %v", err)
	}
	if url != "https://example.com/video.mp4" {
		t.Errorf("URL() = %q", url)
	}
	if resp.Duration() != 6 {
		t.Errorf("Duration() = %d, want 6", resp.Duration())
	}
	if cost, ok := resp.CostUSD(); !ok || cost != float64(ticks)*1e-10 {
		t.Errorf("CostUSD() = %v, %v", cost, ok)
	}
}
