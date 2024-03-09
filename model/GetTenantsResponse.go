package model

type GetTenantsResponse struct {
	ServerTenants []string `json:"SERVER_TENANTS"`
	BrokerTenants []string `json:"BROKER_TENANTS"`
}
