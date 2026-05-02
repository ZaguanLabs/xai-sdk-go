package types

const (
	ChatModelGrok4                           = "grok-4"
	ChatModelGrok40709                       = "grok-4-0709"
	ChatModelGrok4Latest                     = "grok-4-latest"
	ChatModelGrok41Fast                      = "grok-4-1-fast"
	ChatModelGrok41FastReasoning             = "grok-4-1-fast-reasoning"
	ChatModelGrok41FastReasoningLatest       = "grok-4-1-fast-reasoning-latest"
	ChatModelGrok41FastNonReasoning          = "grok-4-1-fast-non-reasoning"
	ChatModelGrok41FastNonReasoningLatest    = "grok-4-1-fast-non-reasoning-latest"
	ChatModelGrok4Fast                       = "grok-4-fast"
	ChatModelGrok4FastReasoning              = "grok-4-fast-reasoning"
	ChatModelGrok4FastReasoningLatest        = "grok-4-fast-reasoning-latest"
	ChatModelGrok4FastNonReasoning           = "grok-4-fast-non-reasoning"
	ChatModelGrok4FastNonReasoningLatest     = "grok-4-fast-non-reasoning-latest"
	ChatModelGrok4200309Reasoning            = "grok-4.20-0309-reasoning"
	ChatModelGrok420                         = "grok-4.20"
	ChatModelGrok4200309                     = "grok-4.20-0309"
	ChatModelGrok420ReasoningLatest          = "grok-4.20-reasoning-latest"
	ChatModelGrok4200309NonReasoning         = "grok-4.20-0309-non-reasoning"
	ChatModelGrok420NonReasoning             = "grok-4.20-non-reasoning"
	ChatModelGrok420NonReasoningLatest       = "grok-4.20-non-reasoning-latest"
	ChatModelGrok420MultiAgent               = "grok-4.20-multi-agent"
	ChatModelGrok420MultiAgent0309           = "grok-4.20-multi-agent-0309"
	ChatModelGrok420MultiAgentLatest         = "grok-4.20-multi-agent-latest"
	ChatModelGrok43                          = "grok-4.3"
	ChatModelGrok43Latest                    = "grok-4.3-latest"
	ChatModelGrokCodeFast1                   = "grok-code-fast-1"
	ChatModelGrok3                           = "grok-3"
	ChatModelGrok3Latest                     = "grok-3-latest"
	ChatModelGrok3Mini                       = "grok-3-mini"
	ChatModelGrok3Fast                       = "grok-3-fast"
	ChatModelGrok3FastLatest                 = "grok-3-fast-latest"
	ChatModelGrok3MiniFast                   = "grok-3-mini-fast"
	ChatModelGrok3MiniFastLatest             = "grok-3-mini-fast-latest"
	ImageGenerationModelGrokImagineImage     = "grok-imagine-image"
	ImageGenerationModelGrokImagineImagePro  = "grok-imagine-image-pro"
	VideoGenerationModelGrokImagineVideo     = "grok-imagine-video"
	ImageFormatBase64                        = "base64"
	ImageFormatURL                           = "url"
	ImageAspectRatio1x1                      = "1:1"
	ImageAspectRatio3x4                      = "3:4"
	ImageAspectRatio4x3                      = "4:3"
	ImageAspectRatio9x16                     = "9:16"
	ImageAspectRatio16x9                     = "16:9"
	ImageAspectRatio2x3                      = "2:3"
	ImageAspectRatio3x2                      = "3:2"
	ImageAspectRatio9x19_5                   = "9:19.5"
	ImageAspectRatio19_5x9                   = "19.5:9"
	ImageAspectRatio9x20                     = "9:20"
	ImageAspectRatio20x9                     = "20:9"
	ImageAspectRatio1x2                      = "1:2"
	ImageAspectRatio2x1                      = "2:1"
	ImageResolution1K                        = "1k"
	ImageResolution2K                        = "2k"
	VideoAspectRatio1x1                      = "1:1"
	VideoAspectRatio16x9                     = "16:9"
	VideoAspectRatio9x16                     = "9:16"
	VideoAspectRatio4x3                      = "4:3"
	VideoAspectRatio3x4                      = "3:4"
	VideoAspectRatio3x2                      = "3:2"
	VideoAspectRatio2x3                      = "2:3"
	VideoResolution480P                      = "480p"
	VideoResolution720P                      = "720p"
	ReasoningEffortLow                       = "low"
	ReasoningEffortHigh                      = "high"
	ImageDetailAuto                          = "auto"
	ImageDetailLow                           = "low"
	ImageDetailHigh                          = "high"
	ToolModeAuto                             = "auto"
	ToolModeNone                             = "none"
	ToolModeRequired                         = "required"
	ResponseFormatText                       = "text"
	ResponseFormatJSONObject                 = "json_object"
	IncludeOptionWebSearchCallOutput         = "web_search_call_output"
	IncludeOptionXSearchCallOutput           = "x_search_call_output"
	IncludeOptionCodeExecutionCallOutput     = "code_execution_call_output"
	IncludeOptionCollectionsSearchCallOutput = "collections_search_call_output"
	IncludeOptionAttachmentSearchCallOutput  = "attachment_search_call_output"
	IncludeOptionMCPCallOutput               = "mcp_call_output"
	IncludeOptionInlineCitations             = "inline_citations"
	IncludeOptionVerboseStreaming            = "verbose_streaming"
)

var ChatModels = []string{
	ChatModelGrok4,
	ChatModelGrok40709,
	ChatModelGrok4Latest,
	ChatModelGrok41Fast,
	ChatModelGrok41FastReasoning,
	ChatModelGrok41FastReasoningLatest,
	ChatModelGrok41FastNonReasoning,
	ChatModelGrok41FastNonReasoningLatest,
	ChatModelGrok4Fast,
	ChatModelGrok4FastReasoning,
	ChatModelGrok4FastReasoningLatest,
	ChatModelGrok4FastNonReasoning,
	ChatModelGrok4FastNonReasoningLatest,
	ChatModelGrok4200309Reasoning,
	ChatModelGrok420,
	ChatModelGrok4200309,
	ChatModelGrok420ReasoningLatest,
	ChatModelGrok4200309NonReasoning,
	ChatModelGrok420NonReasoning,
	ChatModelGrok420NonReasoningLatest,
	ChatModelGrok420MultiAgent,
	ChatModelGrok420MultiAgent0309,
	ChatModelGrok420MultiAgentLatest,
	ChatModelGrok43,
	ChatModelGrok43Latest,
	ChatModelGrokCodeFast1,
	ChatModelGrok3,
	ChatModelGrok3Latest,
	ChatModelGrok3Mini,
	ChatModelGrok3Fast,
	ChatModelGrok3FastLatest,
	ChatModelGrok3MiniFast,
	ChatModelGrok3MiniFastLatest,
}

var ImageGenerationModels = []string{ImageGenerationModelGrokImagineImage, ImageGenerationModelGrokImagineImagePro}
var VideoGenerationModels = []string{VideoGenerationModelGrokImagineVideo}
