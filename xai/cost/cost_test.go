package cost

import (
	"testing"

	xaiv1 "github.com/ZaguanLabs/xai-sdk-go/proto/gen/go/xai/api/v1"
)

func TestUSDFromUsage(t *testing.T) {
	ticks := int64(123)
	got, ok := USDFromUsage(&xaiv1.SamplingUsage{CostInUsdTicks: &ticks})
	if !ok {
		t.Fatal("USDFromUsage() ok = false, want true")
	}
	want := float64(ticks) * USDPerTick
	if got != want {
		t.Errorf("USDFromUsage() = %v, want %v", got, want)
	}
}

func TestUSDFromUsageMissingCost(t *testing.T) {
	if _, ok := USDFromUsage(&xaiv1.SamplingUsage{}); ok {
		t.Fatal("USDFromUsage() ok = true, want false")
	}
	if _, ok := USDFromUsage(nil); ok {
		t.Fatal("USDFromUsage(nil) ok = true, want false")
	}
}
