package model

type IronicNode struct {
	Name         string
	IP           string
	UUID         string `json:"uuid"`
	InstanceUUID string
	Region       string
	Host         string
}
