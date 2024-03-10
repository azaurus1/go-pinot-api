package model

type GetSegmentMetadataResponse struct {
	SegmentStartTime    string `json:"segment.start.time"`
	SegmentTimeUnit     string `json:"segment.time.unit"`
	SegmentSizeInBytes  string `json:"segment.size.in.bytes"`
	SegmentEndTime      string `json:"segment.end.time"`
	SegmentTotalDocs    string `json:"segment.total.docs"`
	SegmentCreationTime string `json:"segment.creation.time"`
	SegmentPushTime     string `json:"segment.push.time"`
	SegmentEndTimeRaw   string `json:"segment.end.time.raw"`
	SegmentTier         string `json:"segment.tier"`
	SegmentStartTimeRaw string `json:"segment.start.time.raw"`
	SegmentIndexVersion string `json:"segment.index.version"`
	CustomMap           string `json:"custom.map"`
	SegmentCrc          string `json:"segment.crc"`
	SegmentDownloadUrl  string `json:"segment.download.url"`
}
