package model

type TableSegmentsConfig struct {
	TimeType                  string `json:"timeType"`
	Replication               string `json:"replication"`
	ReplicasPerPartition      string `json:"replicasPerPartition,omitempty"`
	TimeColumnName            string `json:"timeColumnName"`
	SegmentAssignmentStrategy string `json:"segmentAssignmentStrategy,omitempty"`
	SegmentPushType           string `json:"segmentPushType,omitempty"`
	MinimizeDataMovement      bool   `json:"minimizeDataMovement"`
	RetentionTimeUnit         string `json:"retentionTimeUnit,omitempty"`
	RetentionTimeValue        string `json:"retentionTimeValue,omitempty"`
	CrypterClassName          string `json:"crypterClassName,omitempty"`
	PeerSegmentDownloadScheme string `json:"peerSegmentDownloadScheme,omitempty"`
}

type TableTenant struct {
	Broker string `json:"broker"`
	Server string `json:"server"`
}

type TableIndexConfig struct {
	CreateInvertedIndexDuringSegmentGeneration bool                   `json:"createInvertedIndexDuringSegmentGeneration"`
	SortedColumn                               []string               `json:"sortedColumn"`
	StarTreeIndexConfigs                       []StarTreeIndexConfig  `json:"starTreeIndexConfigs"`
	EnableDefaultStarTree                      bool                   `json:"enableDefaultStarTree"`
	EnableDynamicStarTreeCreation              bool                   `json:"enableDynamicStarTreeCreation"`
	SegmentPartitionConfig                     SegmentPartitionConfig `json:"segmentPartitionConfig"`
	LoadMode                                   string                 `json:"loadMode"`
	ColumnMinMaxValueGeneratorMode             string                 `json:"columnMinMaxValueGeneratorMode"`
	NullHandlingEnabled                        bool                   `json:"nullHandlingEnabled"`
	AggregateMetrics                           bool                   `json:"aggregateMetrics"`
	OptimizeDictionary                         bool                   `json:"optimizeDictionary"`
	OptimizeDictionaryForMetrics               bool                   `json:"optimizeDictionaryForMetrics"`
	NoDictionarySizeRatioThreshold             float64                `json:"noDictionarySizeRatioThreshold"`
	SegmentNameGeneratorType                   string                 `json:"segmentNameGeneratorType"`
}

type ColumnPartitionMapConfig struct {
	FunctionName  string `json:"functionName"`
	NumPartitions int    `json:"numPartitions"`
}

type SegmentPartitionConfig struct {
	ColumnPartitionMap map[string]ColumnPartitionMapConfig `json:"columnPartitionMap"`
}

type StarTreeIndexConfig struct {
	DimensionsSplitOrder              []string `json:"dimensionsSplitOrder"`
	SkipStarNodeCreationForDimensions []string `json:"skipStarNodeCreationForDimensions"`
	FunctionColumnPairs               []string `json:"functionColumnPairs"`
	MaxLeafRecords                    int      `json:"maxLeafRecords"`
}

type TierOverwrite struct {
	StarTreeIndexConfigs []StarTreeIndexConfig `json:"starTreeIndexConfigs"`
}

type TierOverwrites struct {
	HotTier  TierOverwrite `json:"hotTier"`
	ColdTier TierOverwrite `json:"coldTier"`
}

type TableMetadata struct {
	CustomConfigs map[string]string `json:"customConfigs"`
}

type TimestampConfig struct {
	Granulatities []string `json:"granularities"`
}

type FiendIndexInverted struct {
	Enabled string `json:"enabled"`
}

type FieldIndexes struct {
	Inverted FiendIndexInverted `json:"inverted"`
}

type FieldConfig struct {
	Name            string          `json:"name"`
	EncodingType    string          `json:"encodingType"`
	IndexType       string          `json:"indexType"`
	IndexTypes      []string        `json:"indexTypes"`
	TimestampConfig TimestampConfig `json:"timestampConfig"`
	Indexes         FieldIndexes    `json:"indexes"`
}

type TransformConfig struct {
	ColumnName        string `json:"columnName"`
	TransformFunction string `json:"transformFunction"`
}

type TableIngestionConfig struct {
	SegmentTimeValueCheck bool                  `json:"segmentTimeValueCheck,omitempty"`
	TransformConfigs      []TransformConfig     `json:"transformConfigs,omitempty"`
	ContinueOnError       bool                  `json:"continueOnError,omitempty"`
	RowTimeValueCheck     bool                  `json:"rowTimeValueCheck,omitempty"`
	StreamIngestionConfig StreamIngestionConfig `json:"streamIngestionConfig,omitempty"`
}

type StreamIngestionConfig struct {
	StreamConfigMaps []map[string]string `json:"streamConfigMaps"`
}

type TierConfig struct {
	Name                string `json:"name"`
	SegmentSelectorType string `json:"segmentSelectorType"`
	SegmentAge          string `json:"segmentAge"`
	StorageType         string `json:"storageType"`
	ServerTag           string `json:"serverTag"`
}

type Table struct {
	TableName        string               `json:"tableName"`
	TableType        string               `json:"tableType"`
	SegmentsConfig   TableSegmentsConfig  `json:"segmentsConfig"`
	Tenants          TableTenant          `json:"tenants"`
	TableIndexConfig TableIndexConfig     `json:"tableIndexConfig"`
	Metadata         TableMetadata        `json:"metadata"`
	FieldConfigList  []FieldConfig        `json:"fieldConfigList,omitempty"`
	IngestionConfig  TableIngestionConfig `json:"ingestionConfig,omitempty"`
	TierConfigs      []TierConfig         `json:"tierConfigs,omitempty"`
	IsDimTable       bool                 `json:"isDimTable"`
}
