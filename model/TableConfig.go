package model

import (
	"encoding/json"
	"os"
)

type Tenants struct {
	Broker            string         `json:"broker"`
	Server            string         `json:"server"`
	TagOverrideConfig map[string]any `json:"tagOverrideConfig"`
}

type SegmentsConfig struct {
	SchemaName                string `json:"schemaName"`
	TimeColumnName            any    `json:"timeColumnName"`
	Replication               string `json:"replication"`
	ReplicasPerPartition      string `json:"replicasPerPartition"`
	RetentionTimeUnit         any    `json:"retentionTimeUnit"`
	RetentionTimeValue        any    `json:"retentionTimeValue"`
	CompletionConfig          any    `json:"completionConfig"`
	CrypterClassName          any    `json:"crypterClassName"`
	PeerSegmentDownloadScheme any    `json:"peerSegmentDownloadScheme"`
}

type TableIndexConfig struct {
	LoadMode                                   string            `json:"loadMode"`
	InvertedIndexColumns                       []any             `json:"invertedIndexColumns"`
	CreateInvertedIndexDuringSegmentGeneration bool              `json:"createInvertedIndexDuringSegmentGeneration"`
	RangeIndexColumns                          []any             `json:"rangeIndexColumns"`
	SortedColumn                               []any             `json:"sortedColumn"`
	BloomFilterColumns                         []any             `json:"bloomFilterColumns"`
	BloomFilterConfigs                         any               `json:"bloomFilterConfigs"`
	NoDictionaryColumns                        []any             `json:"noDictionaryColumns"`
	OnHeapDictionaryColumns                    []any             `json:"onHeapDictionaryColumns"`
	VarLengthDictionaryColumns                 []any             `json:"varLengthDictionaryColumns"`
	EnableDefaultStarTree                      bool              `json:"enableDefaultStarTree"`
	StarTreeIndexConfigs                       any               `json:"starTreeIndexConfigs"`
	EnableDynamicStarTreeCreation              bool              `json:"enableDynamicStarTreeCreation"`
	SegmentPartitionConfig                     any               `json:"segmentPartitionConfig"`
	ColumnMinMaxValueGeneratorMode             any               `json:"columnMinMaxValueGeneratorMode"`
	AggregateMetrics                           bool              `json:"aggregateMetrics"`
	NullHandlingEnabled                        bool              `json:"nullHandlingEnabled"`
	StreamConfigs                              map[string]string `json:"streamConfigs"`
}

type IngestionConfig struct {
	FilterConfig     any `json:"filterConfig"`
	TransformConfigs any `json:"transformConfigs"`
}

type Quota struct {
	Storage             any `json:"storage"`
	MaxQueriesPerSecond any `json:"maxQueriesPerSecond"`
}

type Routing struct {
	SegmentPrunerTypes   any `json:"segmentPrunerTypes"`
	InstanceSelectorType any `json:"instanceSelectorType"`
}

type Query struct {
	TimeoutMs any `json:"timeoutMs"`
}

type Metadata map[string]any

type UpsertConfig any

type TierConfigs any

type TableConfig struct {
	TableName        string           `json:"tableName"`
	TableType        string           `json:"tableType"`
	Tenants          Tenants          `json:"tenants"`
	SegmentsConfig   SegmentsConfig   `json:"segmentsConfig"`
	TableIndexConfig TableIndexConfig `json:"tableIndexConfig"`
	Metadata         Metadata         `json:"metadata"`
	IngestionConfig  IngestionConfig  `json:"ingestionConfig"`
	Quota            Quota            `json:"quota"`
	Task             any              `json:"task"`
	Routing          Routing          `json:"routing"`
	Query            Query            `json:"query"`
	FieldConfigList  any              `json:"fieldConfigList"`
	UpsertConfig     UpsertConfig     `json:"upsertConfig"`
	TierConfigs      TierConfigs      `json:"tierConfigs"`
}

func NewTableConfigFromFile(fileName string) (TableConfig, error) {
	var tableConfig TableConfig
	f, err := os.Open(fileName)
	if err != nil {
		return tableConfig, err
	}

	defer f.Close()

	err = json.NewDecoder(f).Decode(&tableConfig)
	if err != nil {
		return tableConfig, err
	}

	return tableConfig, nil
}

func (tableConfig *TableConfig) AsBytes() ([]byte, error) {
	tableBytes, err := json.Marshal(tableConfig)
	if err != nil {
		return nil, err
	}

	return tableBytes, nil
}
