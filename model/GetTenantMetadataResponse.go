package model

type GetTenantMetadataResponse struct {
	ServerInstances         []string `json:"ServerInstances"`
	BrokerInstances         []string `json:"BrokerInstances"`
	OfflineServerInstances  []string `json:"OfflineServerInstances"`
	RealtimeServerInstances []string `json:"RealtimeServerInstances"`
	TenantName              string   `json:"TenantName"`
}
