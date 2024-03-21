package model

type CreateTablesResponse struct {
	UnrecognizedProperties map[string]any `json:"unrecognizedProperties"`
	Status                 string         `json:"status"`
}
