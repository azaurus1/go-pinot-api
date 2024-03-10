package model

// Root struct for the JSON array
type GetSegmentsResponse []SegmentDetail

// Event represents the structure within each array element,
// accommodating both offline and realtime events.
type SegmentDetail struct {
	Offline  []string `json:"OFFLINE,omitempty"`
	Realtime []string `json:"REALTIME,omitempty"`
}
