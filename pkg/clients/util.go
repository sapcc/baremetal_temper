package clients

import "github.com/stmcginnis/gofish/redfish"

func calcTotalMemory(mem []*redfish.Memory) (totalMem int) {
	totalMem = 0
	for _, m := range mem {
		totalMem = totalMem + m.CapacityMiB
	}
	return
}
