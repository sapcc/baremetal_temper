module github.com/sapcc/baremetal_temper

go 1.14

require (
	github.com/aristanetworks/goeapi v0.5.0
	github.com/certifi/gocertifi v0.0.0-20200922220541-2c3bb06c6054 // indirect
	github.com/ciscoecosystem/aci-go-client v1.10.1
	github.com/evalphobia/logrus_sentry v0.8.2
	github.com/getsentry/raven-go v0.2.0 // indirect
	github.com/go-openapi/runtime v0.19.21
	github.com/go-ping/ping v0.0.0-20201115131931-3300c582a663
	github.com/gophercloud/gophercloud v0.14.0
	github.com/gorilla/mux v1.8.0
	github.com/netbox-community/go-netbox v0.0.0-20200923200002-49832662a6fd
	github.com/sirupsen/logrus v1.7.0
	github.com/spf13/cobra v0.0.3
	github.com/spf13/viper v1.7.1
	github.com/stmcginnis/gofish v0.7.0
	github.com/stretchr/testify v1.6.1
	github.com/vaughan0/go-ini v0.0.0-20130923145212-a98ad7ee00ec // indirect
	gopkg.in/yaml.v2 v2.3.0
	k8s.io/apimachinery v0.20.0

)

replace github.com/netbox-community/go-netbox => github.com/stefanhipfel/go-netbox v0.0.0-20200928114340-fcd4119414a4

replace github.com/stmcginnis/gofish => github.com/stefanhipfel/gofish v0.9.1-0.20210423073907-81e338649907
