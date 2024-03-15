package model

type ClusterConfig struct {
	AllowParticipantAutoJoin            string `json:"allowParticipantAutoJoin"`
	EnableCaseInsensitive               string `json:"enable.case.insensitive"`
	DefaultHyperlogLogLog2m             string `json:"default.hyperloglog.log2m"`
	PinotBrokerEnableQueryLimitOverride string `json:"pinot.broker.enable.query.limit.override"`
}
