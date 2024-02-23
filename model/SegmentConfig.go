package model

type SegmentConfig struct {
	SegmentAssignmentStrategy string `json:""`
	TimeColumnName            string `json:""`
	TimeType                  string `json:""`
	Replication               int    `json:""`
	SegmentPushType           string `json:""`
	MinimizeDataMovement      bool   `json:""`
}
