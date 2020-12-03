package clients

import (
	"fmt"

	"github.com/stmcginnis/gofish"
)

type RedfishClient struct {
	Host     string
	User     string
	Password string
}

func (r RedfishClient) LoadRedfishInfo(nodeIP string) (i InspectorCallbackData, err error) {
	fmt.Println(nodeIP)
	cfg := gofish.ClientConfig{
		Endpoint:  fmt.Sprintf("https://%s", nodeIP),
		Username:  r.User,
		Password:  r.Password,
		Insecure:  true,
		BasicAuth: false,
	}
	client, err := gofish.Connect(cfg)
	if err != nil {
		return
	}
	defer client.Logout()
	service := client.Service
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
