package chat

import (
	"testing"

	xaiv1 "github.com/ZaguanLabs/xai-sdk-go/proto/gen/go/xai/api/v1"
)

func TestReasoningEffortPython1180Parity(t *testing.T) {
	tests := []struct {
		name string
		want xaiv1.ReasoningEffort
	}{
		{"none", xaiv1.ReasoningEffort_EFFORT_NONE},
		{"low", xaiv1.ReasoningEffort_EFFORT_LOW},
		{"medium", xaiv1.ReasoningEffort_EFFORT_MEDIUM},
		{"high", xaiv1.ReasoningEffort_EFFORT_HIGH},
		{"xhigh", xaiv1.ReasoningEffort_EFFORT_XHIGH},
	}
	for _, test := range tests {
		if got := reasoningEffortToProto(test.name); got != test.want {
			t.Fatalf("reasoningEffortToProto(%q) = %v, want %v", test.name, got, test.want)
		}
		if got := reasoningEffortFromProto(test.want); got != test.name {
			t.Fatalf("reasoningEffortFromProto(%v) = %q, want %q", test.want, got, test.name)
		}
	}
}
