package clients

import (
	"fmt"
	"math"
	"net"

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

// reverseaddr returns the in-addr.arpa.
func reverseaddr(addr string) (arpa string, err error) {
	ip := net.ParseIP(addr)
	if ip == nil {
		return "", fmt.Errorf("unrecognized address: %s", addr)
	}

	if ip.To4() != nil {
		arpa = fmt.Sprint(ip[15]) + "." + fmt.Sprint(ip[14]) + "." + fmt.Sprint(ip[13]) + "." + fmt.Sprint(ip[12]) + ".in-addr.arpa."
		return
	}
	return "", fmt.Errorf("no ip4 address")

}

func reverseZone(addr string) (arpaZone string, err error) {
	ip := net.ParseIP(addr)
	if ip == nil {
		return "", fmt.Errorf("unrecognized address: %s", addr)
	}

	if ip.To4() != nil {
		arpaZone = fmt.Sprint(ip[14]) + "." + fmt.Sprint(ip[13]) + "." + fmt.Sprint(ip[12]) + ".in-addr.arpa."
		return
	}
	return "", fmt.Errorf("no ip4 address")
}
