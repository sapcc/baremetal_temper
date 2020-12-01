package ironic

import (
	"fmt"

	"github.com/stmcginnis/gofish"
)

func (n Node) LoadRedfishInfo() (i InspectorCallbackData, err error) {
	fmt.Println(n.IP)
	cfg := gofish.ClientConfig{
		Endpoint:  fmt.Sprintf("https://%s", n.IP),
		Username:  n.IronicUser,
		Password:  n.IronicPassword,
		Insecure:  true,
		BasicAuth: false,
	}
	c, err := gofish.Connect(cfg)
	if err != nil {
		return
	}
	defer c.Logout()
	service := c.Service
	chassis, err := service.Chassis()
	if err != nil {
		return
	}
	for _, chass := range chassis {
		n, err := chass.NetworkAdapters()
		if err != nil {
			continue
		}
		if len(n) == 0 {
			continue
		}
		i.Interfaces = make([]Interface, len(n))
		f, err := n[0].NetworkDeviceFunctions()
		i.Interfaces[0].MacAddress = f[0].Ethernet.MACAddress
		fmt.Println(f[0].Ethernet.MACAddress)
		fmt.Printf("Chassis: %#v\n\n", chass.Manufacturer)
	}
	return
}
