package node

import (
	"fmt"
	"math"
	"net"
	"strings"

	"github.com/gophercloud/gophercloud/openstack/baremetal/v1/nodes"
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

func findNode(nodes []nodes.Node, n *Node) (err error) {
	for _, nd := range nodes {
		if nd.ProvisionState == "enroll" && nd.Name == "" {
			//node005r-ap017.cc.na-us-1.cloud.sap
			ipmi := fmt.Sprintf("%v", nd.DriverInfo["ipmi_address"])
			s := strings.Split(ipmi, ".")
			if len(s) == 5 {
				if strings.Replace(s[0], "r", "", 1) == n.Name {
					n.UUID = nd.UUID
					n.ProvisionState = nd.ProvisionState
					break
				}
			}
		}
	}
	return
}
