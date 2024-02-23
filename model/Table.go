package model

type Table struct {
	TableName        string            `json:"tableName"`
	TableType        string            `json:"tableType"`
	SegmentsConfig   SegmentConfig     `json:""`
	Tenants          map[string]string `json:""`
	TableIndexConfig TableIndexConfig  `json:""`
	Metadata         Metadata          `json:""`
	FieldConfigList  FieldConfigList   `json:""`
	IngestionConfig  IngestionConfig   `json:""`
	TierConfigs      TierConfig        `json:"tierConfigs"`
	IsDimTable       bool              `json:"isDimTable"`
}
