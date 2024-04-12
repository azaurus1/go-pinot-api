package model

type GetTableIdealStateResponse struct {
	Offline  map[string]map[string]string `json:"OFFLINE,omitempty"`
	Realtime map[string]map[string]string `json:"REALTIME,omitempty"`
}
