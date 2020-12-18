module github.com/sapcc/ironic_temper

go 1.14

require (
	github.com/certifi/gocertifi v0.0.0-20200922220541-2c3bb06c6054 // indirect
	github.com/evalphobia/logrus_sentry v0.8.2
	github.com/getsentry/raven-go v0.2.0 // indirect
	github.com/go-openapi/runtime v0.19.21
	github.com/gophercloud/gophercloud v0.14.0
	github.com/netbox-community/go-netbox v0.0.0-20200923200002-49832662a6fd
	github.com/sirupsen/logrus v1.7.0
	github.com/stmcginnis/gofish v0.7.0
	gopkg.in/yaml.v2 v2.3.0
	k8s.io/apimachinery v0.20.0

)

replace github.com/netbox-community/go-netbox => github.com/stefanhipfel/go-netbox v0.0.0-20200928114340-fcd4119414a4
