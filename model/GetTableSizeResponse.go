package model

type ServerInfo struct {
	SegmentName     string `json:"segmentName"`
	DiskSizeInBytes int64  `json:"diskSizeInBytes"`
}

type SegmentSizeDetail struct {
	ReportedSizeInBytes              int64                 `json:"reportedSizeInBytes"`
	EstimatedSizeInBytes             int64                 `json:"estimatedSizeInBytes"`
	MaxReportedSizePerReplicainBytes int64                 `json:"maxReportedSizePerReplicaInBytes"`
	ServerInfo                       map[string]ServerInfo `json:"serverInfo"`
}

type TableSegments struct {
	ReportedSizeInBytes           int64                        `json:"reportedSizeInBytes"`
	EstimatedSizeInBytes          int64                        `json:"estimatedSizeInBytes"`
	MissingSegments               int                          `json:"missingSegments"`
	ReportedSizePerReplicaInBytes int64                        `json:"reportedSizePerReplicaInBytes"`
	Segments                      map[string]SegmentSizeDetail `json:"segments"`
}

type GetTableSizeResponse struct {
	TableName                     string        `json:"tableName"`
	ReportedSizeInBytes           int64         `json:"reportedSizeInBytes"`
	EstimatedSizeInBytes          int64         `json:"estimatedSizeInBytes"`
	ReportedSizePerReplicaInBytes int64         `json:"reportedSizePerReplicaInBytes"`
	OfflineSegments               TableSegments `json:"offlineSegments"`
	RealtimeSegments              TableSegments `json:"realtimeSegments"`
}
