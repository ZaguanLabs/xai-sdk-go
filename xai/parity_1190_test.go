package xai

import (
	"testing"

	xaiv1 "github.com/ZaguanLabs/xai-sdk-go/proto/gen/go/xai/api/v1"
	"github.com/ZaguanLabs/xai-sdk-go/xai/image"
	"github.com/ZaguanLabs/xai-sdk-go/xai/tools"
	"github.com/ZaguanLabs/xai-sdk-go/xai/video"
)

func TestPython1190GenerationParity(t *testing.T) {
	quality := xaiv1.ImageQuality_IMG_QUALITY_LOW
	if got := image.NewRequest("p", "grok-imagine-image-2.0").WithQuality(quality).Proto().GetQuality(); got != quality {
		t.Fatalf("quality = %v", got)
	}
	generateAudio := false
	resolution := xaiv1.VideoResolution_VIDEO_RESOLUTION_1080P
	req := video.NewGenerateRequestWithOptions("p", "grok-imagine-video-1.5", &video.GenerateOptions{
		ReferenceAudioVoiceIDs: []string{"ara"}, GenerateAudio: &generateAudio, Resolution: &resolution,
	})
	if len(req.ReferenceAudios) != 1 || req.ReferenceAudios[0].GetVoiceId() != "ara" || req.GetGenerateAudio() || req.GetResolution() != resolution {
		t.Fatalf("unexpected video request: %#v", req)
	}
}

func TestPython1190ImageGenerationTool(t *testing.T) {
	tool := tools.ImageGeneration("edit")
	if tool.GetImageGeneration() == nil || tool.GetImageGeneration().GetAction() != "edit" {
		t.Fatalf("unexpected tool: %#v", tool)
	}
	if xaiv1.ToolCallType_TOOL_CALL_TYPE_IMAGE_GENERATION_TOOL.String() != "TOOL_CALL_TYPE_IMAGE_GENERATION_TOOL" {
		t.Fatal("missing image generation tool-call enum")
	}
}
