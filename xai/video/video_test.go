package video

import (
	"testing"
	"time"

	xaiv1 "github.com/ZaguanLabs/xai-sdk-go/proto/gen/go/xai/api/v1"
	xaifiles "github.com/ZaguanLabs/xai-sdk-go/xai/files"
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

func TestNewGenerateRequestWithV117Options(t *testing.T) {
	req := NewGenerateRequestWithOptions("prompt", "model", &GenerateOptions{
		ImageFileID:           "image-file",
		VideoFileID:           "video-file",
		ReferenceImageFileIDs: []string{"ref-file"},
		ReferenceImageURLs:    []string{"https://example.com/ref.jpg"},
		Storage: &xaifiles.StorageOptions{
			Filename:  "output.mp4",
			PublicURL: xaifiles.PublicURLWithDefaults(),
		},
	})

	if req.GetImage().GetFileId() != "image-file" {
		t.Fatalf("image file id = %q", req.GetImage().GetFileId())
	}
	if req.GetVideo().GetFileId() != "video-file" {
		t.Fatalf("video file id = %q", req.GetVideo().GetFileId())
	}
	if len(req.ReferenceImages) != 2 {
		t.Fatalf("reference image count = %d, want 2", len(req.ReferenceImages))
	}
	if req.ReferenceImages[0].GetFileId() != "ref-file" {
		t.Fatalf("first reference file id = %q", req.ReferenceImages[0].GetFileId())
	}
	if req.ReferenceImages[1].GetImageUrl() != "https://example.com/ref.jpg" {
		t.Fatalf("second reference image url = %q", req.ReferenceImages[1].GetImageUrl())
	}
	if req.GetStorageOptions().GetFilename() != "output.mp4" || req.GetStorageOptions().GetPublicUrl() == nil {
		t.Fatalf("storage options = %+v", req.GetStorageOptions())
	}
}

func TestNewExtendRequestWithV117Options(t *testing.T) {
	duration := int32(6)
	req := NewExtendRequestWithOptions("continue", "model", "", &duration, &GenerateOptions{
		VideoFileID: "video-file",
		Storage: &xaifiles.StorageOptions{
			Filename:     "extended.mp4",
			ExpiresAfter: time.Hour,
		},
	})
	if req.GetVideo().GetFileId() != "video-file" {
		t.Fatalf("video file id = %q", req.GetVideo().GetFileId())
	}
	if req.GetStorageOptions().GetFilename() != "extended.mp4" || req.GetStorageOptions().GetExpiresAfter() != int64(time.Hour.Seconds()) {
		t.Fatalf("storage options = %+v", req.GetStorageOptions())
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

func TestPrepareExtension(t *testing.T) {
	duration := int32(6)
	batchReq := PrepareExtension("continue", "model", "https://example.com/video.mp4", "extension-1", &duration)

	if batchReq.GetBatchRequestId() != "extension-1" {
		t.Errorf("batch request id = %q, want extension-1", batchReq.GetBatchRequestId())
	}
	if batchReq.GetVideoExtensionRequest() == nil {
		t.Fatal("video extension request is nil")
	}
	if batchReq.GetVideoExtensionRequest().GetVideo().GetUrl() != "https://example.com/video.mp4" {
		t.Errorf("video url = %q", batchReq.GetVideoExtensionRequest().GetVideo().GetUrl())
	}
	if batchReq.GetVideoExtensionRequest().GetDuration() != duration {
		t.Errorf("duration = %d, want %d", batchReq.GetVideoExtensionRequest().GetDuration(), duration)
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

func TestResponseStorageAccessors(t *testing.T) {
	publicURL := "https://public.example/video.mp4"
	publicURLError := "public URL failed"
	storageError := "storage failed"
	resp := NewResponse(&xaiv1.VideoResponse{
		Video: &xaiv1.GeneratedVideo{
			FileOutput: &xaiv1.FileOutput{
				FileId:         "file-1",
				Filename:       "video.mp4",
				PublicUrl:      &publicURL,
				PublicUrlError: &publicURLError,
			},
			StorageError: &storageError,
		},
	})
	if resp.FileOutput().GetFileId() != "file-1" || resp.PublicURL() != publicURL || resp.PublicURLError() != publicURLError || resp.StorageError() != storageError {
		t.Fatalf("unexpected storage accessors")
	}
}
