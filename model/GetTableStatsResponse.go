package model

type TableStats struct {
	CreationTime string `json:"creationTime"`
}

type GetTableStatsResponse map[string]TableStats
