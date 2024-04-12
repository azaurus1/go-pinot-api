package model

type LiveBrokerInstance struct {
	InstanceName string `json:"instanceName"`
	Port         int    `json:"port"`
	Host         string `json:"host"`
}

type GetLiveBrokersResponse map[string][]LiveBrokerInstance
