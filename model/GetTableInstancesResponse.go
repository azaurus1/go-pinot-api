package model

type BrokerInstances struct {
	TableType string   `json:"tableType,omitempty"`
	Instances []string `json:"instances,omitempty"`
}

type ServerInstances struct {
	TableType string   `json:"tableType,omitempty"`
	Instances []string `json:"instances,omitempty"`
}

type GetTableInstancesResponse struct {
	TableName string            `json:"tableName,omitempty"`
	Brokers   []BrokerInstances `json:"brokers,omitempty"`
	Servers   []ServerInstances `json:"servers,omitempty"`
}
