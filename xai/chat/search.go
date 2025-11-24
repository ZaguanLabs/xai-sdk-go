package chat

import (
	"time"

	xaiv1 "github.com/ZaguanLabs/xai-sdk-go/proto/gen/go/xai/api/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// WebSource configures web search source.
type WebSource struct {
	proto *xaiv1.WebSource
}

// NewWebSource creates a new web search source.
func NewWebSource() *WebSource {
	return &WebSource{
		proto: &xaiv1.WebSource{},
	}
}

// WithExcludedWebsites sets websites to exclude from search.
func (ws *WebSource) WithExcludedWebsites(websites ...string) *WebSource {
	ws.proto.ExcludedWebsites = websites
	return ws
}

// WithAllowedWebsites sets websites to allow in search.
func (ws *WebSource) WithAllowedWebsites(websites ...string) *WebSource {
	ws.proto.AllowedWebsites = websites
	return ws
}

// WithCountry sets the country for search results.
func (ws *WebSource) WithCountry(country string) *WebSource {
	ws.proto.Country = &country
	return ws
}

// WithSafeSearch enables or disables safe search.
func (ws *WebSource) WithSafeSearch(enabled bool) *WebSource {
	ws.proto.SafeSearch = enabled
	return ws
}

// Proto returns the underlying protobuf message.
func (ws *WebSource) Proto() *xaiv1.WebSource {
	return ws.proto
}

// NewsSource configures news search source.
type NewsSource struct {
	proto *xaiv1.NewsSource
}

// NewNewsSource creates a new news search source.
func NewNewsSource() *NewsSource {
	return &NewsSource{
		proto: &xaiv1.NewsSource{},
	}
}

// WithExcludedWebsites sets websites to exclude from news search.
func (ns *NewsSource) WithExcludedWebsites(websites ...string) *NewsSource {
	ns.proto.ExcludedWebsites = websites
	return ns
}

// WithCountry sets the country for news results.
func (ns *NewsSource) WithCountry(country string) *NewsSource {
	ns.proto.Country = &country
	return ns
}

// WithSafeSearch enables or disables safe search for news.
func (ns *NewsSource) WithSafeSearch(enabled bool) *NewsSource {
	ns.proto.SafeSearch = enabled
	return ns
}

// Proto returns the underlying protobuf message.
func (ns *NewsSource) Proto() *xaiv1.NewsSource {
	return ns.proto
}

// XSource configures X (Twitter) search source.
type XSource struct {
	proto *xaiv1.XSource
}

// NewXSource creates a new X search source.
func NewXSource() *XSource {
	return &XSource{
		proto: &xaiv1.XSource{},
	}
}

// WithIncludedHandles sets X handles to include in search.
func (xs *XSource) WithIncludedHandles(handles ...string) *XSource {
	xs.proto.IncludedXHandles = handles
	return xs
}

// WithExcludedHandles sets X handles to exclude from search.
func (xs *XSource) WithExcludedHandles(handles ...string) *XSource {
	xs.proto.ExcludedXHandles = handles
	return xs
}

// WithPostFavoriteCount sets minimum favorite count for posts.
func (xs *XSource) WithPostFavoriteCount(count int32) *XSource {
	xs.proto.PostFavoriteCount = &count
	return xs
}

// WithPostViewCount sets minimum view count for posts.
func (xs *XSource) WithPostViewCount(count int32) *XSource {
	xs.proto.PostViewCount = &count
	return xs
}

// Proto returns the underlying protobuf message.
func (xs *XSource) Proto() *xaiv1.XSource {
	return xs.proto
}

// RssSource configures RSS feed search source.
type RssSource struct {
	proto *xaiv1.RssSource
}

// NewRssSource creates a new RSS search source.
func NewRssSource() *RssSource {
	return &RssSource{
		proto: &xaiv1.RssSource{},
	}
}

// WithLinks sets RSS feed links to search.
func (rs *RssSource) WithLinks(links ...string) *RssSource {
	rs.proto.Links = links
	return rs
}

// Proto returns the underlying protobuf message.
func (rs *RssSource) Proto() *xaiv1.RssSource {
	return rs.proto
}

// Source represents a search source configuration.
type Source struct {
	proto *xaiv1.Source
}

// WithWeb sets the web source configuration.
func (s *Source) WithWeb(web *WebSource) *Source {
	s.proto.Source = &xaiv1.Source_Web{
		Web: web.Proto(),
	}
	return s
}

// WithNews sets the news source configuration.
func (s *Source) WithNews(news *NewsSource) *Source {
	s.proto.Source = &xaiv1.Source_News{
		News: news.Proto(),
	}
	return s
}

// WithX sets the X source configuration.
func (s *Source) WithX(x *XSource) *Source {
	s.proto.Source = &xaiv1.Source_X{
		X: x.Proto(),
	}
	return s
}

// WithRss sets the RSS source configuration.
func (s *Source) WithRss(rss *RssSource) *Source {
	s.proto.Source = &xaiv1.Source_Rss{
		Rss: rss.Proto(),
	}
	return s
}

// NewWebSearchSource creates a source with web search configuration.
func NewWebSearchSource(web *WebSource) *Source {
	return &Source{
		proto: &xaiv1.Source{
			Source: &xaiv1.Source_Web{
				Web: web.Proto(),
			},
		},
	}
}

// NewNewsSearchSource creates a source with news search configuration.
func NewNewsSearchSource(news *NewsSource) *Source {
	return &Source{
		proto: &xaiv1.Source{
			Source: &xaiv1.Source_News{
				News: news.Proto(),
			},
		},
	}
}

// NewXSearchSource creates a source with X search configuration.
func NewXSearchSource(x *XSource) *Source {
	return &Source{
		proto: &xaiv1.Source{
			Source: &xaiv1.Source_X{
				X: x.Proto(),
			},
		},
	}
}

// NewRssSearchSource creates a source with RSS search configuration.
func NewRssSearchSource(rss *RssSource) *Source {
	return &Source{
		proto: &xaiv1.Source{
			Source: &xaiv1.Source_Rss{
				Rss: rss.Proto(),
			},
		},
	}
}

// Proto returns the underlying protobuf message.
func (s *Source) Proto() *xaiv1.Source {
	return s.proto
}

// WithFromDate sets the start date for search results.
func (p *SearchParameters) WithFromDate(fromDate time.Time) *SearchParameters {
	p.pb.FromDate = timestamppb.New(fromDate)
	return p
}

// WithToDate sets the end date for search results.
func (p *SearchParameters) WithToDate(toDate time.Time) *SearchParameters {
	p.pb.ToDate = timestamppb.New(toDate)
	return p
}

// WithSources sets custom search sources.
func (p *SearchParameters) WithSources(sources ...*Source) *SearchParameters {
	p.pb.Sources = make([]*xaiv1.Source, len(sources))
	for i, src := range sources {
		p.pb.Sources[i] = src.Proto()
	}
	return p
}
