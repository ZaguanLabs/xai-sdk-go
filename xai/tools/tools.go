package tools

import (
	"strings"
	"time"

	xaiv1 "github.com/ZaguanLabs/xai-sdk-go/proto/gen/go/xai/api/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type WebSearchOptions struct {
	ExcludedDomains          []string
	AllowedDomains           []string
	EnableImageUnderstanding bool
	UserLocationCountry      string
	UserLocationCity         string
	UserLocationRegion       string
	UserLocationTimezone     string
}

type XSearchOptions struct {
	FromDate                 time.Time
	ToDate                   time.Time
	AllowedXHandles          []string
	ExcludedXHandles         []string
	EnableImageUnderstanding bool
	EnableVideoUnderstanding bool
}

type CollectionsSearchOptions struct {
	CollectionIDs []string
	Limit         *int32
	Instructions  string
	RetrievalMode string
}

type MCPOptions struct {
	ServerURL         string
	ServerLabel       string
	ServerDescription string
	AllowedToolNames  []string
	Authorization     string
	ExtraHeaders      map[string]string
}

func WebSearch(opts WebSearchOptions) *xaiv1.Tool {
	ws := &xaiv1.WebSearch{
		ExcludedDomains:          opts.ExcludedDomains,
		AllowedDomains:           opts.AllowedDomains,
		EnableImageUnderstanding: &opts.EnableImageUnderstanding,
	}
	if opts.UserLocationCountry != "" || opts.UserLocationCity != "" || opts.UserLocationRegion != "" || opts.UserLocationTimezone != "" {
		ws.UserLocation = &xaiv1.WebSearchUserLocation{
			Country:  stringPtr(opts.UserLocationCountry),
			City:     stringPtr(opts.UserLocationCity),
			Region:   stringPtr(opts.UserLocationRegion),
			Timezone: stringPtr(opts.UserLocationTimezone),
		}
	}
	return &xaiv1.Tool{Tool: &xaiv1.Tool_WebSearch{WebSearch: ws}}
}

func XSearch(opts XSearchOptions) *xaiv1.Tool {
	xs := &xaiv1.XSearch{
		AllowedXHandles:          opts.AllowedXHandles,
		ExcludedXHandles:         opts.ExcludedXHandles,
		EnableImageUnderstanding: &opts.EnableImageUnderstanding,
		EnableVideoUnderstanding: &opts.EnableVideoUnderstanding,
	}
	if !opts.FromDate.IsZero() {
		xs.FromDate = timestamppb.New(opts.FromDate)
	}
	if !opts.ToDate.IsZero() {
		xs.ToDate = timestamppb.New(opts.ToDate)
	}
	return &xaiv1.Tool{Tool: &xaiv1.Tool_XSearch{XSearch: xs}}
}

func CodeExecution() *xaiv1.Tool {
	return &xaiv1.Tool{Tool: &xaiv1.Tool_CodeExecution{CodeExecution: &xaiv1.CodeExecution{}}}
}

func CollectionsSearch(opts CollectionsSearchOptions) *xaiv1.Tool {
	cs := &xaiv1.CollectionsSearch{
		CollectionIds: opts.CollectionIDs,
		Limit:         opts.Limit,
		Instructions:  stringPtr(opts.Instructions),
	}
	switch opts.RetrievalMode {
	case "", "hybrid":
		if opts.RetrievalMode != "" {
			cs.RetrievalMode = &xaiv1.CollectionsSearch_HybridRetrieval{HybridRetrieval: &xaiv1.HybridRetrieval{}}
		}
	case "semantic":
		cs.RetrievalMode = &xaiv1.CollectionsSearch_SemanticRetrieval{SemanticRetrieval: &xaiv1.SemanticRetrieval{}}
	case "keyword":
		cs.RetrievalMode = &xaiv1.CollectionsSearch_KeywordRetrieval{KeywordRetrieval: &xaiv1.KeywordRetrieval{}}
	}
	return &xaiv1.Tool{Tool: &xaiv1.Tool_CollectionsSearch{CollectionsSearch: cs}}
}

func MCP(opts MCPOptions) *xaiv1.Tool {
	return &xaiv1.Tool{Tool: &xaiv1.Tool_Mcp{Mcp: &xaiv1.MCP{
		ServerUrl:         opts.ServerURL,
		ServerLabel:       opts.ServerLabel,
		ServerDescription: opts.ServerDescription,
		AllowedToolNames:  opts.AllowedToolNames,
		Authorization:     stringPtr(opts.Authorization),
		ExtraHeaders:      opts.ExtraHeaders,
	}}}
}

func GetToolCallType(toolCall *xaiv1.ToolCall) string {
	if toolCall == nil {
		return ""
	}
	return strings.ToLower(strings.TrimPrefix(toolCall.Type.String(), "TOOL_CALL_TYPE_"))
}

func stringPtr(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}
