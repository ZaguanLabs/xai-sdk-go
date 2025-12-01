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

// ============================================================================
// Include Options
// ============================================================================

// IncludeOption specifies additional output to include in responses.
type IncludeOption string

const (
	// IncludeWebSearchCallOutput includes web search call output in the response.
	IncludeWebSearchCallOutput IncludeOption = "web_search_call_output"
	// IncludeXSearchCallOutput includes X search call output in the response.
	IncludeXSearchCallOutput IncludeOption = "x_search_call_output"
	// IncludeCodeExecutionCallOutput includes code execution call output in the response.
	IncludeCodeExecutionCallOutput IncludeOption = "code_execution_call_output"
	// IncludeCollectionsSearchCallOutput includes collections search call output in the response.
	IncludeCollectionsSearchCallOutput IncludeOption = "collections_search_call_output"
	// IncludeDocumentSearchCallOutput includes document search call output in the response.
	IncludeDocumentSearchCallOutput IncludeOption = "document_search_call_output"
	// IncludeMCPCallOutput includes MCP call output in the response.
	IncludeMCPCallOutput IncludeOption = "mcp_call_output"
	// IncludeInlineCitations includes inline citations in the response.
	IncludeInlineCitations IncludeOption = "inline_citations"
)

// includeOptionToProto converts an IncludeOption to its proto enum value.
func includeOptionToProto(opt IncludeOption) xaiv1.IncludeOption {
	switch opt {
	case IncludeWebSearchCallOutput:
		return xaiv1.IncludeOption_INCLUDE_OPTION_WEB_SEARCH_CALL_OUTPUT
	case IncludeXSearchCallOutput:
		return xaiv1.IncludeOption_INCLUDE_OPTION_X_SEARCH_CALL_OUTPUT
	case IncludeCodeExecutionCallOutput:
		return xaiv1.IncludeOption_INCLUDE_OPTION_CODE_EXECUTION_CALL_OUTPUT
	case IncludeCollectionsSearchCallOutput:
		return xaiv1.IncludeOption_INCLUDE_OPTION_COLLECTIONS_SEARCH_CALL_OUTPUT
	case IncludeDocumentSearchCallOutput:
		return xaiv1.IncludeOption_INCLUDE_OPTION_DOCUMENT_SEARCH_CALL_OUTPUT
	case IncludeMCPCallOutput:
		return xaiv1.IncludeOption_INCLUDE_OPTION_MCP_CALL_OUTPUT
	case IncludeInlineCitations:
		return xaiv1.IncludeOption_INCLUDE_OPTION_INLINE_CITATIONS
	default:
		return xaiv1.IncludeOption_INCLUDE_OPTION_INVALID
	}
}

// includeOptionFromProto converts a proto enum value to an IncludeOption.
func includeOptionFromProto(opt xaiv1.IncludeOption) IncludeOption {
	switch opt {
	case xaiv1.IncludeOption_INCLUDE_OPTION_WEB_SEARCH_CALL_OUTPUT:
		return IncludeWebSearchCallOutput
	case xaiv1.IncludeOption_INCLUDE_OPTION_X_SEARCH_CALL_OUTPUT:
		return IncludeXSearchCallOutput
	case xaiv1.IncludeOption_INCLUDE_OPTION_CODE_EXECUTION_CALL_OUTPUT:
		return IncludeCodeExecutionCallOutput
	case xaiv1.IncludeOption_INCLUDE_OPTION_COLLECTIONS_SEARCH_CALL_OUTPUT:
		return IncludeCollectionsSearchCallOutput
	case xaiv1.IncludeOption_INCLUDE_OPTION_DOCUMENT_SEARCH_CALL_OUTPUT:
		return IncludeDocumentSearchCallOutput
	case xaiv1.IncludeOption_INCLUDE_OPTION_MCP_CALL_OUTPUT:
		return IncludeMCPCallOutput
	case xaiv1.IncludeOption_INCLUDE_OPTION_INLINE_CITATIONS:
		return IncludeInlineCitations
	default:
		return ""
	}
}

// ============================================================================
// Inline Citations
// ============================================================================

// InlineCitation represents an inline citation in the response content.
type InlineCitation struct {
	proto *xaiv1.InlineCitation
}

// ID returns the unique identifier for this citation.
func (ic *InlineCitation) ID() string {
	if ic.proto == nil {
		return ""
	}
	return ic.proto.Id
}

// StartIndex returns the character index in the content where this citation starts.
func (ic *InlineCitation) StartIndex() int32 {
	if ic.proto == nil {
		return 0
	}
	return ic.proto.StartIndex
}

// WebCitation returns the web citation details if this is a web citation.
func (ic *InlineCitation) WebCitation() *WebCitationInfo {
	if ic.proto == nil || ic.proto.WebCitation == nil {
		return nil
	}
	return &WebCitationInfo{proto: ic.proto.WebCitation}
}

// XCitation returns the X citation details if this is an X/Twitter citation.
func (ic *InlineCitation) XCitation() *XCitationInfo {
	if ic.proto == nil || ic.proto.XCitation == nil {
		return nil
	}
	return &XCitationInfo{proto: ic.proto.XCitation}
}

// CollectionsCitation returns the collections citation details if this is a collections citation.
func (ic *InlineCitation) CollectionsCitation() *CollectionsCitationInfo {
	if ic.proto == nil || ic.proto.CollectionsCitation == nil {
		return nil
	}
	return &CollectionsCitationInfo{proto: ic.proto.CollectionsCitation}
}

// Proto returns the underlying protobuf message.
func (ic *InlineCitation) Proto() *xaiv1.InlineCitation {
	return ic.proto
}

// WebCitationInfo contains details about a web citation.
type WebCitationInfo struct {
	proto *xaiv1.WebCitation
}

// URL returns the URL of the web source.
func (wc *WebCitationInfo) URL() string {
	if wc.proto == nil {
		return ""
	}
	return wc.proto.Url
}

// XCitationInfo contains details about an X/Twitter citation.
type XCitationInfo struct {
	proto *xaiv1.XCitation
}

// URL returns the URL of the X post.
func (xc *XCitationInfo) URL() string {
	if xc.proto == nil {
		return ""
	}
	return xc.proto.Url
}

// CollectionsCitationInfo contains details about a collections citation.
type CollectionsCitationInfo struct {
	proto *xaiv1.CollectionsCitation
}

// FileID returns the file ID of the cited document.
func (cc *CollectionsCitationInfo) FileID() string {
	if cc.proto == nil {
		return ""
	}
	return cc.proto.FileId
}

// ChunkID returns the chunk ID within the file.
func (cc *CollectionsCitationInfo) ChunkID() string {
	if cc.proto == nil {
		return ""
	}
	return cc.proto.ChunkId
}

// ChunkContent returns the content of the cited chunk.
func (cc *CollectionsCitationInfo) ChunkContent() string {
	if cc.proto == nil {
		return ""
	}
	return cc.proto.ChunkContent
}

// Score returns the relevance score of this citation.
func (cc *CollectionsCitationInfo) Score() float32 {
	if cc.proto == nil {
		return 0
	}
	return cc.proto.Score
}

// CollectionIDs returns the collection IDs this citation belongs to.
func (cc *CollectionsCitationInfo) CollectionIDs() []string {
	if cc.proto == nil {
		return nil
	}
	return cc.proto.CollectionIds
}

// ============================================================================
// Tool Call Type
// ============================================================================

// ToolCallType indicates whether a tool call is client-side or server-side.
type ToolCallType string

const (
	// ToolCallTypeClientSide indicates a client-side tool that should be executed by the client.
	ToolCallTypeClientSide ToolCallType = "client_side"
	// ToolCallTypeWebSearch indicates a server-side web search tool.
	ToolCallTypeWebSearch ToolCallType = "web_search"
	// ToolCallTypeXSearch indicates a server-side X search tool.
	ToolCallTypeXSearch ToolCallType = "x_search"
	// ToolCallTypeCodeExecution indicates a server-side code execution tool.
	ToolCallTypeCodeExecution ToolCallType = "code_execution"
	// ToolCallTypeCollectionsSearch indicates a server-side collections search tool.
	ToolCallTypeCollectionsSearch ToolCallType = "collections_search"
	// ToolCallTypeMCP indicates a server-side MCP tool.
	ToolCallTypeMCP ToolCallType = "mcp"
	// ToolCallTypeDocumentSearch indicates a server-side document search tool.
	ToolCallTypeDocumentSearch ToolCallType = "document_search"
)

// toolCallTypeFromProto converts a proto enum value to a ToolCallType.
func toolCallTypeFromProto(t xaiv1.ToolCallType) ToolCallType {
	switch t {
	case xaiv1.ToolCallType_TOOL_CALL_TYPE_CLIENT_SIDE_TOOL:
		return ToolCallTypeClientSide
	case xaiv1.ToolCallType_TOOL_CALL_TYPE_WEB_SEARCH_TOOL:
		return ToolCallTypeWebSearch
	case xaiv1.ToolCallType_TOOL_CALL_TYPE_X_SEARCH_TOOL:
		return ToolCallTypeXSearch
	case xaiv1.ToolCallType_TOOL_CALL_TYPE_CODE_EXECUTION_TOOL:
		return ToolCallTypeCodeExecution
	case xaiv1.ToolCallType_TOOL_CALL_TYPE_COLLECTIONS_SEARCH_TOOL:
		return ToolCallTypeCollectionsSearch
	case xaiv1.ToolCallType_TOOL_CALL_TYPE_MCP_TOOL:
		return ToolCallTypeMCP
	case xaiv1.ToolCallType_TOOL_CALL_TYPE_DOCUMENT_SEARCH_TOOL:
		return ToolCallTypeDocumentSearch
	default:
		return ""
	}
}
