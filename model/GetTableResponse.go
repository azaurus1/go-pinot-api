package model

type GetTableResponse struct {
	OFFLINE  Table `json:"OFFLINE"`
	REALTIME Table `json:"REALTIME"`
}
