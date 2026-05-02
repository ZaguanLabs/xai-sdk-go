package batch

import (
	"testing"

	xaiv1 "github.com/ZaguanLabs/xai-sdk-go/proto/gen/go/xai/api/v1"
	"github.com/ZaguanLabs/xai-sdk-go/xai/image"
	"github.com/ZaguanLabs/xai-sdk-go/xai/video"
	"google.golang.org/genproto/googleapis/rpc/status"
)

func TestRequestFromImageRequest(t *testing.T) {
	req := image.NewRequest("prompt", "image-model").WithCount(2)
	batchReq := RequestFromImageRequest(req, "image-1")

	if batchReq.GetBatchRequestId() != "image-1" {
		t.Errorf("batch request id = %q, want image-1", batchReq.GetBatchRequestId())
	}
	if batchReq.GetImageRequest() == nil {
		t.Fatal("image request is nil")
	}
	if batchReq.GetImageRequest().GetN() != 2 {
		t.Errorf("n = %d, want 2", batchReq.GetImageRequest().GetN())
	}
}

func TestRequestFromVideoRequest(t *testing.T) {
	req := video.NewGenerateRequest("prompt", "video-model")
	batchReq := RequestFromVideoRequest(req, "video-1")

	if batchReq.GetBatchRequestId() != "video-1" {
		t.Errorf("batch request id = %q, want video-1", batchReq.GetBatchRequestId())
	}
	if batchReq.GetVideoRequest() == nil {
		t.Fatal("video request is nil")
	}
	if batchReq.GetVideoRequest().Model != "video-model" {
		t.Errorf("model = %q, want video-model", batchReq.GetVideoRequest().Model)
	}
}

func TestPrepareVideoRequest(t *testing.T) {
	duration := int32(8)
	batchReq := PrepareVideoRequest("prompt", "video-model", "video-2", &video.GenerateOptions{Duration: &duration})

	if batchReq.GetBatchRequestId() != "video-2" {
		t.Errorf("batch request id = %q, want video-2", batchReq.GetBatchRequestId())
	}
	if batchReq.GetVideoRequest().GetDuration() != duration {
		t.Errorf("duration = %d, want %d", batchReq.GetVideoRequest().GetDuration(), duration)
	}
}

func TestBatchResultDataSupportsImageAndVideo(t *testing.T) {
	imageResult := &xaiv1.BatchResultData{Response: &xaiv1.BatchResultData_ImageResponse{ImageResponse: &xaiv1.ImageResponse{Model: "image-model"}}}
	videoResult := &xaiv1.BatchResultData{Response: &xaiv1.BatchResultData_VideoResponse{VideoResponse: &xaiv1.VideoResponse{Model: "video-model"}}}

	if imageResult.GetImageResponse().Model != "image-model" {
		t.Errorf("image response model = %q", imageResult.GetImageResponse().Model)
	}
	if videoResult.GetVideoResponse().Model != "video-model" {
		t.Errorf("video response model = %q", videoResult.GetVideoResponse().Model)
	}
}

func TestResultWrappers(t *testing.T) {
	success := NewResult(&xaiv1.BatchResult{
		BatchRequestId: "request-1",
		Result: &xaiv1.BatchResult_Response{Response: &xaiv1.BatchResultData{
			Response: &xaiv1.BatchResultData_ImageResponse{ImageResponse: &xaiv1.ImageResponse{Model: "image-model"}},
		}},
	})

	if !success.IsSuccess() || success.HasError() {
		t.Fatal("success result status mismatch")
	}
	if success.BatchRequestID() != "request-1" || success.ImageResponse().Model != "image-model" {
		t.Fatalf("success wrapper mismatch: %v", success)
	}

	failed := NewResult(&xaiv1.BatchResult{
		BatchRequestId: "request-2",
		Result:         &xaiv1.BatchResult_Error{Error: &status.Status{Message: "failed"}},
	})

	if !failed.HasError() || failed.IsSuccess() || failed.ErrorMessage() != "failed" {
		t.Fatal("failed result status mismatch")
	}

	list := NewListBatchResultsResponse(&xaiv1.ListBatchResultsResponse{
		Results: []*xaiv1.BatchResult{success.Proto(), failed.Proto()},
	})
	if len(list.Succeeded()) != 1 || len(list.Failed()) != 1 || len(list.Results()) != 2 {
		t.Fatal("list result wrapper mismatch")
	}
}
