package cost

import xaiv1 "github.com/ZaguanLabs/xai-sdk-go/proto/gen/go/xai/api/v1"

const USDPerTick = 1e-10

func USDFromUsage(usage *xaiv1.SamplingUsage) (float64, bool) {
	if usage == nil || usage.CostInUsdTicks == nil {
		return 0, false
	}
	return float64(usage.GetCostInUsdTicks()) * USDPerTick, true
}
