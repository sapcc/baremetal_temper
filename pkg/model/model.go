package model

type IronicNode struct {
	Name   string
	IP     string
	UUID   string `json:"uuid"`
	Region string
	Host   string
}
