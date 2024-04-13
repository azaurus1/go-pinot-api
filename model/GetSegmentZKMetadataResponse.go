package model

type ZkSegmentMetadata struct {
	CustomMap           string `json:"custom.map"`
	SegmentCRC          string `json:"segment.crc"`
	SegmentCreationTime string `json:"segment.creation.time"`
	SegmentDownloadURL  string `json:"segment.download.url"`
	SegmentEndTime      string `json:"segment.end.time"`
	SegmentEndTimeRaw   string `json:"segment.end.time.raw"`
	SegmentIndexVersion string `json:"segment.index.version"`
	SegmentPushTime     string `json:"segment.push.time"`
	SegmentSizeInBytes  string `json:"segment.size.bytes"`
	SegmentStartTime    string `json:"segment.start.time"`
	SegmentStartTimeRaw string `json:"segment.start.time.raw"`
	SegmentTier         string `json:"segment.tier"`
	SegmentTimeUnit     string `json:"segment.time.unit"`
	SegmentTotalDocs    string `json:"segment.total.docs"`
}

type GetSegmentZKMetadataResponse map[string]ZkSegmentMetadata
