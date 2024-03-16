package model

type Instance struct {
	Host string `json:"host"`
	Port int    `json:"port"`
	Type string `json:"type"` // Can be CONTROLLER, BROKER, SERVER, MINION
}
