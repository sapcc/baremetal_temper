module github.com/sapcc/ironic_temper

go 1.14

require (
	github.com/gophercloud/gophercloud v0.14.0
	github.com/netbox-community/go-netbox v0.0.0-20200923200002-49832662a6fd
	github.com/prometheus/common v0.15.0 // indirect
	github.com/prometheus/prometheus v2.5.0+incompatible // indirect
	github.com/sirupsen/logrus v1.7.0
	github.com/stmcginnis/gofish v0.7.0
	gopkg.in/yaml.v2 v2.3.0
	k8s.io/apimachinery v0.20.0

)

replace github.com/netbox-community/go-netbox => github.com/stefanhipfel/go-netbox v0.0.0-20200928114340-fcd4119414a4
