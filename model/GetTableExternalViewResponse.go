package model

type GetTableExternalViewResponse struct {
	Offline  map[string]map[string]string `json:"offline,omitempty"`
	Realtime map[string]map[string]string `json:"realtime,omitempty"`
}
