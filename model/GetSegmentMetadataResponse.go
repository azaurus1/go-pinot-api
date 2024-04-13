package model

type StarTreeIndex struct {
	DimensionColumns        []string `json:"dimension-columns"`
	MetricAggregations      []string `json:"metric-aggregations"`
	MaxLeafRecords          int64    `json:"max-leaf-records"`
	DimensionColumnsSkipped []string `json:"dimension-columns-skipped"`
}

type SegmentMetadata struct {
	SegmentName          string            `json:"segmentName"`
	SchemaName           string            `json:"schemaName"`
	CRC                  int64             `json:"crc"`
	CreationTimeMillis   int64             `json:"creationTimeMillis"`
	CreationTimeReadable string            `json:"creationTimeReadable"`
	TimeColumn           string            `json:"timeColumn"`
	TimeUnit             string            `json:"timeUnit"`
	TimeGranularitySec   int64             `json:"timeGranularitySec"`
	StartTimeMillis      int64             `json:"startTimeMillis"`
	StartTimeReadable    string            `json:"startTimeReadable"`
	EndTimeMillis        int64             `json:"endTimeMillis"`
	EndTimeReadable      string            `json:"endTimeReadable"`
	SegmentVersion       string            `json:"segmentVersion"`
	CreatorName          string            `json:"creatorName"`
	TotalDocs            int64             `json:"totalDocs"`
	Custom               map[string]string `json:"custom"`
	StartOffset          int64             `json:"startOffset"`
	EndOffset            int64             `json:"endOffset"`
	Columns              []string          `json:"columns"`
	Indexes              any               `json:"indexes"`
	StarTreeIndex        []StarTreeIndex   `json:"star-tree-index"`
}

type GetSegmentMetadataResponse map[string]SegmentMetadata
