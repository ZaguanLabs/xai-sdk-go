package chat

import xaiv1 "github.com/ZaguanLabs/xai-sdk-go/proto/gen/go/xai/api/v1"

// RequestSettings represents the settings that were used for a request.
// This is returned in the response to show what settings were actually applied.
type RequestSettings struct {
	proto *xaiv1.RequestSettings
}

// MaxTokens returns the max tokens setting.
func (rs *RequestSettings) MaxTokens() int32 {
	if rs.proto == nil || rs.proto.MaxTokens == nil {
		return 0
	}
	return *rs.proto.MaxTokens
}

// ParallelToolCalls returns whether parallel tool calls were enabled.
func (rs *RequestSettings) ParallelToolCalls() bool {
	if rs.proto == nil {
		return false
	}
	return rs.proto.ParallelToolCalls
}

// PreviousResponseID returns the previous response ID if set.
func (rs *RequestSettings) PreviousResponseID() string {
	if rs.proto == nil || rs.proto.PreviousResponseId == nil {
		return ""
	}
	return *rs.proto.PreviousResponseId
}

// ReasoningEffort returns the reasoning effort level.
func (rs *RequestSettings) ReasoningEffort() string {
	if rs.proto == nil || rs.proto.ReasoningEffort == nil {
		return ""
	}
	return reasoningEffortFromProto(*rs.proto.ReasoningEffort)
}

// Temperature returns the temperature setting.
func (rs *RequestSettings) Temperature() float32 {
	if rs.proto == nil || rs.proto.Temperature == nil {
		return 0
	}
	return *rs.proto.Temperature
}

// TopP returns the top_p setting.
func (rs *RequestSettings) TopP() float32 {
	if rs.proto == nil || rs.proto.TopP == nil {
		return 0
	}
	return *rs.proto.TopP
}

// User returns the user identifier.
func (rs *RequestSettings) User() string {
	if rs.proto == nil {
		return ""
	}
	return rs.proto.User
}

// StoreMessages returns whether message storage was enabled.
func (rs *RequestSettings) StoreMessages() bool {
	if rs.proto == nil {
		return false
	}
	return rs.proto.StoreMessages
}

// UseEncryptedContent returns whether encrypted content was enabled.
func (rs *RequestSettings) UseEncryptedContent() bool {
	if rs.proto == nil {
		return false
	}
	return rs.proto.UseEncryptedContent
}

// Proto returns the underlying protobuf message.
func (rs *RequestSettings) Proto() *xaiv1.RequestSettings {
	return rs.proto
}

// DebugOutput contains debugging information from the API.
type DebugOutput struct {
	proto *xaiv1.DebugOutput
}

// Attempts returns the number of attempts made.
func (d *DebugOutput) Attempts() int32 {
	if d.proto == nil {
		return 0
	}
	return d.proto.Attempts
}

// Request returns the request string.
func (d *DebugOutput) Request() string {
	if d.proto == nil {
		return ""
	}
	return d.proto.Request
}

// Prompt returns the prompt string.
func (d *DebugOutput) Prompt() string {
	if d.proto == nil {
		return ""
	}
	return d.proto.Prompt
}

// Responses returns the response strings.
func (d *DebugOutput) Responses() []string {
	if d.proto == nil {
		return nil
	}
	return d.proto.Responses
}

// CacheReadCount returns the number of cache reads.
func (d *DebugOutput) CacheReadCount() uint32 {
	if d.proto == nil {
		return 0
	}
	return d.proto.CacheReadCount
}

// CacheReadInputBytes returns the bytes read from cache.
func (d *DebugOutput) CacheReadInputBytes() uint64 {
	if d.proto == nil {
		return 0
	}
	return d.proto.CacheReadInputBytes
}

// CacheWriteCount returns the number of cache writes.
func (d *DebugOutput) CacheWriteCount() uint32 {
	if d.proto == nil {
		return 0
	}
	return d.proto.CacheWriteCount
}

// CacheWriteInputBytes returns the bytes written to cache.
func (d *DebugOutput) CacheWriteInputBytes() uint64 {
	if d.proto == nil {
		return 0
	}
	return d.proto.CacheWriteInputBytes
}

// EngineRequest returns the engine request string.
func (d *DebugOutput) EngineRequest() string {
	if d.proto == nil {
		return ""
	}
	return d.proto.EngineRequest
}

// LBAddress returns the load balancer address.
func (d *DebugOutput) LBAddress() string {
	if d.proto == nil {
		return ""
	}
	return d.proto.LbAddress
}

// SamplerTag returns the sampler tag.
func (d *DebugOutput) SamplerTag() string {
	if d.proto == nil {
		return ""
	}
	return d.proto.SamplerTag
}

// Chunks returns the chunk strings.
func (d *DebugOutput) Chunks() []string {
	if d.proto == nil {
		return nil
	}
	return d.proto.Chunks
}

// Proto returns the underlying protobuf message.
func (d *DebugOutput) Proto() *xaiv1.DebugOutput {
	return d.proto
}

// LogProb represents log probability information for a token.
type LogProb struct {
	proto *xaiv1.LogProb
}

// Token returns the token string.
func (lp *LogProb) Token() string {
	if lp.proto == nil {
		return ""
	}
	return lp.proto.Token
}

// Logprob returns the log probability value.
func (lp *LogProb) Logprob() float32 {
	if lp.proto == nil {
		return 0
	}
	return lp.proto.Logprob
}

// Bytes returns the token bytes.
func (lp *LogProb) Bytes() []byte {
	if lp.proto == nil {
		return nil
	}
	return lp.proto.Bytes
}

// TopLogProbs returns the top log probabilities.
func (lp *LogProb) TopLogProbs() []*TopLogProb {
	if lp.proto == nil || len(lp.proto.TopLogprobs) == 0 {
		return nil
	}
	result := make([]*TopLogProb, len(lp.proto.TopLogprobs))
	for i, tlp := range lp.proto.TopLogprobs {
		result[i] = &TopLogProb{proto: tlp}
	}
	return result
}

// TopLogProb represents a top log probability entry.
type TopLogProb struct {
	proto *xaiv1.TopLogProb
}

// Token returns the token string.
func (tlp *TopLogProb) Token() string {
	if tlp.proto == nil {
		return ""
	}
	return tlp.proto.Token
}

// Logprob returns the log probability value.
func (tlp *TopLogProb) Logprob() float32 {
	if tlp.proto == nil {
		return 0
	}
	return tlp.proto.Logprob
}

// Bytes returns the token bytes.
func (tlp *TopLogProb) Bytes() []byte {
	if tlp.proto == nil {
		return nil
	}
	return tlp.proto.Bytes
}

// LogProbs represents log probabilities for content.
type LogProbs struct {
	proto *xaiv1.LogProbs
}

// Content returns the log probability content.
func (lps *LogProbs) Content() []*LogProb {
	if lps.proto == nil || len(lps.proto.Content) == 0 {
		return nil
	}
	result := make([]*LogProb, len(lps.proto.Content))
	for i, lp := range lps.proto.Content {
		result[i] = &LogProb{proto: lp}
	}
	return result
}
