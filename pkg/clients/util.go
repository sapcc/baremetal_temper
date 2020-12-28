package clients

import (
	"math"

	"github.com/stmcginnis/gofish/redfish"
)

func calcTotalMemory(mem []*redfish.Memory) (totalMem int) {
	totalMem = 0
	for _, m := range mem {
		totalMem = totalMem + m.CapacityMiB
	}
	return
}

func calcDelta(a, b int) (delta float64) {
	diff := math.Abs(float64(a - b))
	delta = (diff / float64(b))
	return
}
