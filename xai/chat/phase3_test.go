package chat

import (
	"testing"
	"time"
)

// Phase 3 Tests: SearchParameters date filtering and custom sources

func TestSearchParametersWithFromDate(t *testing.T) {
	params := NewSearchParameters()
	fromDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	params.WithFromDate(fromDate)

	proto := params.Proto()
	if proto.FromDate == nil {
		t.Fatal("FromDate is nil")
	}

	if proto.FromDate.AsTime().Unix() != fromDate.Unix() {
		t.Errorf("FromDate = %v, want %v", proto.FromDate.AsTime(), fromDate)
	}

	t.Log("✅ SearchParameters.WithFromDate() works correctly")
}

func TestSearchParametersWithToDate(t *testing.T) {
	params := NewSearchParameters()
	toDate := time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)

	params.WithToDate(toDate)

	proto := params.Proto()
	if proto.ToDate == nil {
		t.Fatal("ToDate is nil")
	}

	if proto.ToDate.AsTime().Unix() != toDate.Unix() {
		t.Errorf("ToDate = %v, want %v", proto.ToDate.AsTime(), toDate)
	}

	t.Log("✅ SearchParameters.WithToDate() works correctly")
}

func TestSearchParametersDateRange(t *testing.T) {
	params := NewSearchParameters()
	fromDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	toDate := time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)

	params.WithFromDate(fromDate).WithToDate(toDate)

	proto := params.Proto()
	if proto.FromDate == nil || proto.ToDate == nil {
		t.Fatal("Date range not set")
	}

	t.Log("✅ SearchParameters date range chaining works correctly")
}

func TestWebSource(t *testing.T) {
	web := NewWebSource().
		WithExcludedWebsites("spam.com", "ads.com").
		WithAllowedWebsites("trusted.com").
		WithCountry("US").
		WithSafeSearch(true)

	proto := web.Proto()

	if len(proto.ExcludedWebsites) != 2 {
		t.Errorf("ExcludedWebsites length = %d, want 2", len(proto.ExcludedWebsites))
	}

	if proto.ExcludedWebsites[0] != "spam.com" {
		t.Errorf("ExcludedWebsites[0] = %q, want %q", proto.ExcludedWebsites[0], "spam.com")
	}

	if len(proto.AllowedWebsites) != 1 {
		t.Errorf("AllowedWebsites length = %d, want 1", len(proto.AllowedWebsites))
	}

	if proto.AllowedWebsites[0] != "trusted.com" {
		t.Errorf("AllowedWebsites[0] = %q, want %q", proto.AllowedWebsites[0], "trusted.com")
	}

	if proto.Country != "US" {
		t.Errorf("Country = %q, want %q", proto.Country, "US")
	}

	if !proto.SafeSearch {
		t.Error("SafeSearch = false, want true")
	}

	t.Log("✅ WebSource works correctly")
}

func TestNewsSource(t *testing.T) {
	news := NewNewsSource().
		WithExcludedWebsites("fakenews.com").
		WithCountry("UK").
		WithSafeSearch(true)

	proto := news.Proto()

	if len(proto.ExcludedWebsites) != 1 {
		t.Errorf("ExcludedWebsites length = %d, want 1", len(proto.ExcludedWebsites))
	}

	if proto.Country != "UK" {
		t.Errorf("Country = %q, want %q", proto.Country, "UK")
	}

	if !proto.SafeSearch {
		t.Error("SafeSearch = false, want true")
	}

	t.Log("✅ NewsSource works correctly")
}

func TestXSource(t *testing.T) {
	x := NewXSource().
		WithIncludedHandles("@elonmusk", "@xai").
		WithExcludedHandles("@spammer").
		WithPostFavoriteCount(100).
		WithPostViewCount(1000)

	proto := x.Proto()

	if len(proto.IncludedXHandles) != 2 {
		t.Errorf("IncludedXHandles length = %d, want 2", len(proto.IncludedXHandles))
	}

	if proto.IncludedXHandles[0] != "@elonmusk" {
		t.Errorf("IncludedXHandles[0] = %q, want %q", proto.IncludedXHandles[0], "@elonmusk")
	}

	if len(proto.ExcludedXHandles) != 1 {
		t.Errorf("ExcludedXHandles length = %d, want 1", len(proto.ExcludedXHandles))
	}

	if proto.PostFavoriteCount != 100 {
		t.Errorf("PostFavoriteCount = %d, want 100", proto.PostFavoriteCount)
	}

	if proto.PostViewCount != 1000 {
		t.Errorf("PostViewCount = %d, want 1000", proto.PostViewCount)
	}

	t.Log("✅ XSource works correctly")
}

func TestRssSource(t *testing.T) {
	rss := NewRssSource().
		WithLinks("https://blog.example.com/feed", "https://news.example.com/rss")

	proto := rss.Proto()

	if len(proto.Links) != 2 {
		t.Errorf("Links length = %d, want 2", len(proto.Links))
	}

	if proto.Links[0] != "https://blog.example.com/feed" {
		t.Errorf("Links[0] = %q, want %q", proto.Links[0], "https://blog.example.com/feed")
	}

	t.Log("✅ RssSource works correctly")
}

func TestSearchParametersWithWebSource(t *testing.T) {
	web := NewWebSource().
		WithAllowedWebsites("trusted.com").
		WithSafeSearch(true)

	source := NewWebSearchSource(web)

	params := NewSearchParameters().
		WithSources(source)

	proto := params.Proto()

	if len(proto.Sources) != 1 {
		t.Fatalf("Sources length = %d, want 1", len(proto.Sources))
	}

	if proto.Sources[0].Web == nil {
		t.Fatal("Web source is nil")
	}

	if len(proto.Sources[0].Web.AllowedWebsites) != 1 {
		t.Errorf("AllowedWebsites length = %d, want 1", len(proto.Sources[0].Web.AllowedWebsites))
	}

	t.Log("✅ SearchParameters with WebSource works correctly")
}

func TestSearchParametersWithMultipleSources(t *testing.T) {
	web := NewWebSource().WithAllowedWebsites("trusted.com")
	news := NewNewsSource().WithCountry("US")
	x := NewXSource().WithIncludedHandles("@xai")
	rss := NewRssSource().WithLinks("https://blog.example.com/feed")

	params := NewSearchParameters().
		WithSources(
			NewWebSearchSource(web),
			NewNewsSearchSource(news),
			NewXSearchSource(x),
			NewRssSearchSource(rss),
		)

	proto := params.Proto()

	if len(proto.Sources) != 4 {
		t.Fatalf("Sources length = %d, want 4", len(proto.Sources))
	}

	// Verify each source type
	if proto.Sources[0].Web == nil {
		t.Error("First source should be Web")
	}

	if proto.Sources[1].News == nil {
		t.Error("Second source should be News")
	}

	if proto.Sources[2].X == nil {
		t.Error("Third source should be X")
	}

	if proto.Sources[3].Rss == nil {
		t.Error("Fourth source should be RSS")
	}

	t.Log("✅ SearchParameters with multiple sources works correctly")
}

func TestSearchParametersComplete(t *testing.T) {
	// Test complete search parameters with all features
	fromDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	toDate := time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)

	web := NewWebSource().
		WithAllowedWebsites("trusted.com").
		WithSafeSearch(true)

	params := NewSearchParameters().
		WithMode("auto").
		WithReturnCitations(true).
		WithMaxSearchResults(10).
		WithFromDate(fromDate).
		WithToDate(toDate).
		WithSources(NewWebSearchSource(web))

	proto := params.Proto()

	// Verify all fields are set
	if proto.FromDate == nil {
		t.Error("FromDate is nil")
	}

	if proto.ToDate == nil {
		t.Error("ToDate is nil")
	}

	if !proto.ReturnCitations {
		t.Error("ReturnCitations = false, want true")
	}

	if proto.MaxSearchResults != 10 {
		t.Errorf("MaxSearchResults = %d, want 10", proto.MaxSearchResults)
	}

	if len(proto.Sources) != 1 {
		t.Errorf("Sources length = %d, want 1", len(proto.Sources))
	}

	t.Log("✅ Complete SearchParameters with all Phase 3 features works correctly")
}
