package model

type GetSegmentTiersResponse struct {
	SegmentTiers map[string]map[string]string `json:"segmentTiers"`
	TableName    string                       `json:"tableName"`
}
