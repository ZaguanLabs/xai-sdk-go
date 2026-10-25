package chat

import (
	"fmt"

	xaiv1 "github.com/ZaguanLabs/xai-sdk-go/proto/gen/go/xai/api/v1"
	"google.golang.org/protobuf/proto"
)

type GuidelineCategory string
type GuidelineBehavior string

const (
	GuidelineCategoryNSFW      GuidelineCategory = "nsfw"
	GuidelineCategoryViolence  GuidelineCategory = "violence"
	GuidelineCategoryProfanity GuidelineCategory = "profanity"
	GuidelineBehaviorDecline   GuidelineBehavior = "decline"
)

// Guideline controls moderation of one category of model output.
type Guideline struct {
	Type     string            `json:"type"`
	Category GuidelineCategory `json:"category"`
	Behavior GuidelineBehavior `json:"behavior"`
}

// Proto validates the same dictionary values accepted by the Python SDK.
func (g Guideline) Proto() (*xaiv1.Guideline, error) {
	if g.Type != "policy" {
		return nil, fmt.Errorf("invalid guideline type: %q", g.Type)
	}
	var category xaiv1.GuidelineCategory
	switch g.Category {
	case GuidelineCategoryNSFW:
		category = xaiv1.GuidelineCategory_GUIDELINE_CATEGORY_NSFW
	case GuidelineCategoryViolence:
		category = xaiv1.GuidelineCategory_GUIDELINE_CATEGORY_VIOLENCE
	case GuidelineCategoryProfanity:
		category = xaiv1.GuidelineCategory_GUIDELINE_CATEGORY_PROFANITY
	default:
		return nil, fmt.Errorf("invalid guideline category: %q", g.Category)
	}
	if g.Behavior != GuidelineBehaviorDecline {
		return nil, fmt.Errorf("invalid guideline behavior: %q", g.Behavior)
	}
	return &xaiv1.Guideline{Kind: &xaiv1.Guideline_Policy{Policy: &xaiv1.PolicyGuideline{Category: category, Behavior: xaiv1.GuidelineBehavior_GUIDELINE_BEHAVIOR_DECLINE}}}, nil
}

// SetGuidelines validates and replaces the guidelines atomically. Repeated
// categories and model compatibility are validated by the API, as in Python.
func (r *Request) SetGuidelines(guidelines ...Guideline) error {
	values := make([]*xaiv1.Guideline, len(guidelines))
	for i, g := range guidelines {
		value, err := g.Proto()
		if err != nil {
			return err
		}
		values[i] = value
	}
	r.proto.Guidelines = values
	return nil
}

// WithGuidelines accepts native protobuf guidelines, matching Python's proto path.
func WithGuidelines(guidelines ...*xaiv1.Guideline) RequestOption {
	return func(r *Request) { r.SetProtoGuidelines(guidelines...) }
}

func (r *Request) SetProtoGuidelines(guidelines ...*xaiv1.Guideline) *Request {
	r.proto.Guidelines = make([]*xaiv1.Guideline, len(guidelines))
	for i, g := range guidelines {
		if g != nil {
			r.proto.Guidelines[i] = proto.Clone(g).(*xaiv1.Guideline)
		}
	}
	return r
}

// WithSafetyIdentifier sets an opaque stable end-user identifier.
func WithSafetyIdentifier(identifier string) RequestOption {
	return func(r *Request) { r.SetSafetyIdentifier(identifier) }
}

func (r *Request) SetSafetyIdentifier(identifier string) *Request {
	r.proto.SafetyIdentifier = &identifier
	return r
}
