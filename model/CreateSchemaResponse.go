package model

type CreateSchemaResponse struct {
	UnrecognizedProperties map[string]string `json:"unrecognizedProperties,omitempty"`
	Status                 string            `json:"status"`
}
