package chat

import (
	"testing"

	xaiv1 "github.com/ZaguanLabs/xai-sdk-go/proto/gen/go/xai/api/v1"
	"google.golang.org/protobuf/proto"
)

func TestGuidelines1210(t *testing.T) {
	r := NewRequest("grok", WithUser("legacy"), WithSafetyIdentifier("opaque"))
	for _, category := range []GuidelineCategory{GuidelineCategoryNSFW, GuidelineCategoryViolence, GuidelineCategoryProfanity} {
		if err := r.SetGuidelines(Guideline{Type: "policy", Category: category, Behavior: GuidelineBehaviorDecline}); err != nil {
			t.Fatal(err)
		}
		data, err := proto.Marshal(r.Proto())
		if err != nil {
			t.Fatal(err)
		}
		var wire xaiv1.GetCompletionsRequest
		if err := proto.Unmarshal(data, &wire); err != nil {
			t.Fatal(err)
		}
		if wire.GetSafetyIdentifier() != "opaque" || wire.User != "legacy" || len(wire.Guidelines) != 1 || wire.Guidelines[0].GetPolicy().Behavior != 1 {
			t.Fatalf("wire: %v", &wire)
		}
	}
	if r.Proto().Guidelines[0].GetPolicy().Category != 5 {
		t.Fatal("profanity must use wire value 5")
	}
	for _, invalid := range []Guideline{{Type: "unknown", Category: GuidelineCategoryNSFW, Behavior: GuidelineBehaviorDecline}, {Type: "policy", Category: "unknown", Behavior: GuidelineBehaviorDecline}, {Type: "policy", Category: GuidelineCategoryNSFW, Behavior: "allow"}} {
		before := proto.Clone(r.Proto())
		if err := r.SetGuidelines(invalid); err == nil {
			t.Fatal("expected validation error")
		}
		if !proto.Equal(before, r.Proto()) {
			t.Fatal("failed validation changed request")
		}
	}
	descriptor := r.Proto().ProtoReflect().Descriptor().Fields()
	if descriptor.ByName("safety_identifier").Number() != 34 || descriptor.ByName("guidelines").Number() != 39 {
		t.Fatal("incorrect upstream field numbers")
	}
	g := r.Proto().Guidelines[0]
	copy := NewRequest("grok", WithGuidelines(g), WithSafetyIdentifier(""))
	g.GetPolicy().Category = 1
	if copy.Proto().Guidelines[0].GetPolicy().Category != 5 || copy.Proto().SafetyIdentifier == nil {
		t.Fatal("proto copying or optional presence lost")
	}
}
