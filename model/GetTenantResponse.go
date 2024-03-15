package model

// GetTenantResponse Get tenant response
type GetTenantResponse struct {
	ServerInstances []string `json:"ServerInstances"`
	BrokerInstances []string `json:"BrokerInstances"`
	TenantName      string   `json:"TenantName"`
}
