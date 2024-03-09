package model

type CreateTablesResponse struct {
	UnrecognizedProperties map[string][]string `json:"unrecognizedProperties"`
	Status                 string              `json:"status"`
}
