package search

import (
	"fmt"
	"time"

	xaiv1 "github.com/ZaguanLabs/xai-sdk-go/proto/gen/go/xai/api/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	ModeAuto = "auto"
	ModeOn   = "on"
	ModeOff  = "off"
)

type Parameters struct {
	Sources          []*xaiv1.Source
	Mode             string
	FromDate         time.Time
	ToDate           time.Time
	ReturnCitations  bool
	MaxSearchResults *int32
}

func (p Parameters) Proto() (*xaiv1.SearchParameters, error) {
	mode, err := modeToProto(p.Mode)
	if err != nil {
		return nil, err
	}
	pb := &xaiv1.SearchParameters{
		Sources:         p.Sources,
		Mode:            mode,
		ReturnCitations: p.ReturnCitations,
	}
	if !p.FromDate.IsZero() {
		pb.FromDate = timestamppb.New(p.FromDate)
	}
	if !p.ToDate.IsZero() {
		pb.ToDate = timestamppb.New(p.ToDate)
	}
	if p.MaxSearchResults != nil {
		pb.MaxSearchResults = p.MaxSearchResults
	}
	return pb, nil
}

func WebSource(country string, excludedWebsites, allowedWebsites []string, safeSearch bool) *xaiv1.Source {
	return &xaiv1.Source{Source: &xaiv1.Source_Web{Web: &xaiv1.WebSource{
		Country:          stringPtr(country),
		ExcludedWebsites: excludedWebsites,
		AllowedWebsites:  allowedWebsites,
		SafeSearch:       safeSearch,
	}}}
}

func NewsSource(country string, excludedWebsites []string, safeSearch bool) *xaiv1.Source {
	return &xaiv1.Source{Source: &xaiv1.Source_News{News: &xaiv1.NewsSource{
		Country:          stringPtr(country),
		ExcludedWebsites: excludedWebsites,
		SafeSearch:       safeSearch,
	}}}
}

func XSource(includedXHandles, excludedXHandles []string, postFavoriteCount, postViewCount *int32) *xaiv1.Source {
	return &xaiv1.Source{Source: &xaiv1.Source_X{X: &xaiv1.XSource{
		IncludedXHandles:  includedXHandles,
		ExcludedXHandles:  excludedXHandles,
		PostFavoriteCount: postFavoriteCount,
		PostViewCount:     postViewCount,
	}}}
}

func RSSSource(links []string) *xaiv1.Source {
	return &xaiv1.Source{Source: &xaiv1.Source_Rss{Rss: &xaiv1.RssSource{Links: links}}}
}

func modeToProto(mode string) (xaiv1.SearchMode, error) {
	switch mode {
	case "", ModeAuto:
		return xaiv1.SearchMode_AUTO_SEARCH_MODE, nil
	case ModeOn:
		return xaiv1.SearchMode_ON_SEARCH_MODE, nil
	case ModeOff:
		return xaiv1.SearchMode_OFF_SEARCH_MODE, nil
	default:
		return xaiv1.SearchMode_INVALID_SEARCH_MODE, fmt.Errorf("invalid search mode: %s", mode)
	}
}

func stringPtr(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}
