package tools

import "testing"

func TestToolHelpers(t *testing.T) {
	web := WebSearch(WebSearchOptions{ExcludedDomains: []string{"example.com"}, EnableImageUnderstanding: true, EnableImageSearch: true, UserLocationCountry: "US"})
	if web.GetWebSearch() == nil || !web.GetWebSearch().GetEnableImageUnderstanding() || !web.GetWebSearch().GetEnableImageSearch() || web.GetWebSearch().GetUserLocation().GetCountry() != "US" {
		t.Fatalf("WebSearch() = %v", web)
	}

	x := XSearch(XSearchOptions{AllowedXHandles: []string{"xai"}, EnableVideoUnderstanding: true})
	if x.GetXSearch() == nil || !x.GetXSearch().GetEnableVideoUnderstanding() || len(x.GetXSearch().GetAllowedXHandles()) != 1 {
		t.Fatalf("XSearch() = %v", x)
	}

	code := CodeExecution()
	if code.GetCodeExecution() == nil {
		t.Fatalf("CodeExecution() = %v", code)
	}

	limit := int32(3)
	collections := CollectionsSearch(CollectionsSearchOptions{CollectionIDs: []string{"col-1"}, Limit: &limit, RetrievalMode: "semantic"})
	if collections.GetCollectionsSearch() == nil || collections.GetCollectionsSearch().GetLimit() != limit || collections.GetCollectionsSearch().GetSemanticRetrieval() == nil {
		t.Fatalf("CollectionsSearch() = %v", collections)
	}

	mcp := MCP(MCPOptions{ServerURL: "https://example.com", ServerLabel: "server", Authorization: "token"})
	if mcp.GetMcp() == nil || mcp.GetMcp().GetServerUrl() != "https://example.com" || mcp.GetMcp().GetAuthorization() != "token" {
		t.Fatalf("MCP() = %v", mcp)
	}
}
