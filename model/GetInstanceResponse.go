package model

type SystemResourceInfo struct {
	NumCores      string `json:"numCores"`
	TotalMemoryMB string `json:"totalMemoryMB"`
	MaxHeapSizeMB string `json:"maxHeapSizeMB"`
}

type GetInstanceResponse struct {
	InstanceName       string             `json:"instanceName"`
	Hostname           string             `json:"hostname"`
	Enabled            bool               `json:"enabled"`
	Port               string             `json:"port"`
	Tags               []string           `json:"tags"`
	Pools              []string           `json:"pools"`
	GRPCPort           int                `json:"grpcPort"`
	AdminPort          int                `json:"adminPort"`
	QueryServicePort   int                `json:"queryServicePort"`
	QueryMailboxPort   int                `json:"queryMailboxPort"`
	SystemResourceInfo SystemResourceInfo `json:"systemResourceInfo,omitempty"`
}
