package search

import "testing"

func TestSourceHelpers(t *testing.T) {
	web := WebSource("US", []string{"example.com"}, nil, true)
	if web.GetWeb().GetCountry() != "US" || !web.GetWeb().GetSafeSearch() {
		t.Fatalf("WebSource() = %v", web.GetWeb())
	}

	news := NewsSource("UK", []string{"example.org"}, false)
	if news.GetNews().GetCountry() != "UK" || news.GetNews().GetSafeSearch() {
		t.Fatalf("NewsSource() = %v", news.GetNews())
	}

	favorites := int32(10)
	x := XSource([]string{"xai"}, []string{"spam"}, &favorites, nil)
	if len(x.GetX().GetIncludedXHandles()) != 1 || x.GetX().GetPostFavoriteCount() != favorites {
		t.Fatalf("XSource() = %v", x.GetX())
	}

	rss := RSSSource([]string{"https://example.com/rss"})
	if len(rss.GetRss().GetLinks()) != 1 {
		t.Fatalf("RSSSource() = %v", rss.GetRss())
	}
}

func TestParametersProto(t *testing.T) {
	max := int32(5)
	params, err := (Parameters{Mode: ModeOn, ReturnCitations: true, MaxSearchResults: &max}).Proto()
	if err != nil {
		t.Fatalf("Proto() error = %v", err)
	}
	if params.GetMode().String() != "ON_SEARCH_MODE" || !params.GetReturnCitations() || params.GetMaxSearchResults() != max {
		t.Fatalf("Proto() = %v", params)
	}
}
