package clients

import (
	"fmt"
	"net"

	"github.com/stmcginnis/gofish"
)

type RedfishClient struct {
	Host     string
	User     string
	Password string
	client   *gofish.APIClient
	service  *gofish.Service
	data     *InspectorCallbackData
}

func (r RedfishClient) LoadRedfishInfo(nodeIP string) (i *InspectorCallbackData, err error) {
	cfg := gofish.ClientConfig{
		Endpoint:  fmt.Sprintf("https://%s", nodeIP),
		Username:  r.User,
		Password:  r.Password,
		Insecure:  true,
		BasicAuth: true,
	}
	client, err := gofish.Connect(cfg)
	if err != nil {
		return
	}
	defer client.Logout()
	r.client = client
	r.data = &InspectorCallbackData{}
	r.service = client.Service
	if err = r.setBMCAddress(); err != nil {
		return
	}
	r.setInventory()
	i = r.data
	return
}

func (r RedfishClient) setBMCAddress() (err error) {
	m, err := r.service.Managers()
	if err != nil && len(m) == 0 {
		return fmt.Errorf("cannot set bmc address")
	}
	in, err := m[0].EthernetInterfaces()
	if err != nil || len(in) == 0 {
		return
	}
	addr, err := net.LookupAddr(in[0].IPv4Addresses[0].Address)
	if err != nil || len(addr) == 0 {
		return
	}
	r.data.Inventory.BmcAddress = addr[0]
	return
}

func (r RedfishClient) setInventory() (err error) {
	ch, err := r.service.Chassis()
	if err != nil || len(ch) == 0 {
		return
	}

	r.data.Inventory.Manufacturer = ch[0].Manufacturer
	r.data.Inventory.Serial = ch[0].SerialNumber
	r.data.Inventory.Model = ch[0].Model

	if err = r.setMemory(); err != nil {
		return
	}
	if err = r.setDisks(); err != nil {
		return
	}
	if err = r.setCPUs(); err != nil {
		return
	}
	if err = r.setNetworkDevicesData(); err != nil {
		return
	}
	return
}

func (r RedfishClient) setMemory() (err error) {
	s, err := r.service.Systems()
	if err != nil || len(s) == 0 {
		return
	}
	r.data.Inventory.Memory.PhysicalMb = int(s[0].MemorySummary.TotalSystemPersistentMemoryGiB)
	return
}

func (r RedfishClient) setDisks() (err error) {
	s, err := r.service.Systems()
	if err != nil || len(s) == 0 {
		return
	}
	st, err := s[0].Storage()
	rootDisk := RootDisk{
		Rotational: true,
	}
	r.data.Inventory.Disks = make([]Disk, 0)
	for _, s := range st {
		ds, err := s.Drives()
		if err != nil {
			continue
		}
		for _, s := range ds {
			rotational := true
			if s.RotationSpeedRPM == 0 {
				rotational = false
			}
			disk := Disk{
				Name:       s.Name,
				Model:      s.Model,
				Vendor:     s.Manufacturer,
				Size:       s.CapacityBytes,
				Rotational: rotational,
			}

			if s.CapacityBytes > rootDisk.Size {
				rootDisk.Size = s.CapacityBytes
				rootDisk.Name = s.Name
				rootDisk.Model = s.Model
				rootDisk.Vendor = s.Manufacturer
				if s.RotationSpeedRPM == 0 {
					rootDisk.Rotational = rotational
				}
			}
			r.data.Inventory.Disks = append(r.data.Inventory.Disks, disk)
		}
	}

	r.data.RootDisk = rootDisk
	return
}

func (r RedfishClient) setCPUs() (err error) {
	sys, err := r.service.Systems()
	cpu, err := sys[0].Processors()
	r.data.Inventory.CPU.Count = sys[0].ProcessorSummary.LogicalProcessorCount
	r.data.Inventory.CPU.Architecture = string(cpu[0].InstructionSet)
	return
}

func (r RedfishClient) setNetworkDevicesData() (err error) {
	sys, err := r.service.Systems()
	if err != nil {
		return
	}
	if len(sys) == 0 {
		return
	}
	ethInt, err := sys[0].EthernetInterfaces()
	if err != nil {
		return
	}
	if len(ethInt) == 0 {
		return
	}
	r.data.Inventory.Interfaces = make([]Interface, len(ethInt))
	for i, e := range ethInt {
		//r.data.BootInterface = e.MACAddress
		r.data.Inventory.Interfaces[i].MacAddress = e.MACAddress
		r.data.Inventory.Interfaces[i].Name = e.Name
		r.data.Inventory.Interfaces[i].ClientID = &e.ID
		r.data.Inventory.Interfaces[i].HasCarrier = true
	}

	return
}
