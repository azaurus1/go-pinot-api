package model

type TableSegmentsConfig struct {
	TimeType                       string `json:"timeType"`
	Replication                    string `json:"replication"`
	ReplicasPerPartition           string `json:"replicasPerPartition,omitempty"`
	TimeColumnName                 string `json:"timeColumnName"`
	SegmentAssignmentStrategy      string `json:"segmentAssignmentStrategy,omitempty"`
	SegmentPushType                string `json:"segmentPushType,omitempty"`
	MinimizeDataMovement           bool   `json:"minimizeDataMovement"`
	RetentionTimeUnit              string `json:"retentionTimeUnit,omitempty"`
	RetentionTimeValue             string `json:"retentionTimeValue,omitempty"`
	CrypterClassName               string `json:"crypterClassName,omitempty"`
	PeerSegmentDownloadScheme      string `json:"peerSegmentDownloadScheme,omitempty"`
	CompletionMode                 string `json:"completionMode,omitempty"`
	DeletedSegmentsRetentionPeriod string `json:"deletedSegmentsRetentionPeriod,omitempty"`
}

type TableTenant struct {
	Broker string `json:"broker"`
	Server string `json:"server"`
}

type TableIndexConfig struct {
	CreateInvertedIndexDuringSegmentGeneration bool                    `json:"createInvertedIndexDuringSegmentGeneration"`
	SortedColumn                               []string                `json:"sortedColumn,omitempty"`
	StarTreeIndexConfigs                       []*StarTreeIndexConfig  `json:"starTreeIndexConfigs,omitempty"`
	EnableDefaultStarTree                      bool                    `json:"enableDefaultStarTree,omitempty"`
	EnableDynamicStarTreeCreation              bool                    `json:"enableDynamicStarTreeCreation,omitempty"`
	SegmentPartitionConfig                     *SegmentPartitionConfig `json:"segmentPartitionConfig,omitempty"`
	LoadMode                                   string                  `json:"loadMode"`
	ColumnMinMaxValueGeneratorMode             string                  `json:"columnMinMaxValueGeneratorMode,omitempty"`
	NullHandlingEnabled                        bool                    `json:"nullHandlingEnabled"`
	AggregateMetrics                           bool                    `json:"aggregateMetrics,omitempty"`
	OptimizeDictionary                         bool                    `json:"optimizeDictionary"`
	OptimizeDictionaryForMetrics               bool                    `json:"optimizeDictionaryForMetrics"`
	NoDictionarySizeRatioThreshold             float64                 `json:"noDictionarySizeRatioThreshold"`
	SegmentNameGeneratorType                   string                  `json:"segmentNameGeneratorType"`
	InvertedIndexColumns                       []string                `json:"invertedIndexColumns,omitempty"`
	NoDictionaryColumns                        []string                `json:"noDictionaryColumns,omitempty"`
	RangeIndexColumns                          []string                `json:"rangeIndexColumns,omitempty"`
	OnHeapDictionaryColumns                    []string                `json:"onHeapDictionaryColumns,omitempty"`
	VarLengthDictionaryColumns                 []string                `json:"varLengthDictionaryColumns,omitempty"`
	BloomFilterColumns                         []string                `json:"bloomFilterColumns,omitempty"`
	RangeIndexVersion                          int                     `json:"rangeIndexVersion,omitempty"`
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
	StarTreeIndexConfigs []*StarTreeIndexConfig `json:"starTreeIndexConfigs"`
}

type TierOverwrites struct {
	HotTier  *TierOverwrite `json:"hotTier,omitempty"`
	ColdTier *TierOverwrite `json:"coldTier,omitempty"`
}

type TableMetadata struct {
	CustomConfigs map[string]string `json:"customConfigs"`
}

type TimestampConfig struct {
	Granularities []string `json:"granularities"`
}

type FiendIndexInverted struct {
	Enabled string `json:"enabled"`
}

type FieldIndexes struct {
	Inverted *FiendIndexInverted `json:"inverted,omitempty"`
}

type FieldConfig struct {
	Name            string           `json:"name"`
	EncodingType    string           `json:"encodingType"`
	IndexType       string           `json:"indexType"`
	IndexTypes      []string         `json:"indexTypes"`
	TimestampConfig *TimestampConfig `json:"timestampConfig,omitempty"`
	Indexes         *FieldIndexes    `json:"indexes,omitempty"`
}

type QueryConfig struct {
	TimeoutMs              string            `json:"timeoutMs"`
	DisableGroovy          bool              `json:"disableGroovy"`
	UseApproximateFunction bool              `json:"useApproximateFunction"`
	ExpressionOverrideMap  map[string]string `json:"expressionOverrideMap"`
	MaxQueryResponseBytes  string            `json:"maxQueryResponseBytes"`
	MaxServerResponseBytes string            `json:"maxServerResponseBytes"`
}

type RoutingConfig struct {
	SegmentPrunerTypes   []string `json:"segmentPrunerTypes"`
	InstanceSelectorType string   `json:"instanceSelectorType"`
}

type QuotaConfig struct {
	Storage             string `json:"storage"`
	MaxQueriesPerSecond int    `json:"maxQueriesPerSecond"`
}

type UpsertConfig struct {
	Mode                         string            `json:"mode"`
	PartialUpsertStrategies      map[string]string `json:"partialUpsertStrategies,omitempty"`
	DefaultPartialUpsertStrategy string            `json:"defaultPartialUpsertStrategy,omitempty"`
	ComparisonColumns            string            `json:"comparisonColumn,omitempty"`
	DeleteRecordColumn           string            `json:"deleteRecordColumn,omitempty"`
	MetadataTTL                  int64             `json:"metadataTTL,omitempty"`
	DeletedKeysTTL               float64           `json:"deletedKeysTTL,omitempty"`
	HashFunction                 string            `json:"hashFunction,omitempty"`
	EnableSnapshot               *bool             `json:"enableSnapshot,omitempty"`
	EnablePreLoad                *bool             `json:"enablePreLoad,omitempty"`
	DropOutOfOrderRecord         *bool             `json:"dropOutOfOrderRecord,omitempty"`
	OutOfOrderRecordColumn       string            `json:"outOfOrderRecordColumn,omitempty"`
	MetadataManagerClass         string            `json:"metadataManagerClass,omitempty"`
	MetadataManagerConfigs       map[string]string `json:"metadataManagerConfigs,omitempty"`
}

type DedupConfig struct {
	DedupEnabled bool   `json:"dedupEnabled"`
	HashFunction string `json:"hashFunction"`
}

type TransformConfig struct {
	ColumnName        string `json:"columnName"`
	TransformFunction string `json:"transformFunction"`
}

type FilterConfig struct {
	FilterFunction string `json:"filterFunction"`
}

type TableIngestionConfig struct {
	SegmentTimeValueCheck bool                   `json:"segmentTimeValueCheck,omitempty"`
	TransformConfigs      []TransformConfig      `json:"transformConfigs,omitempty"`
	FilterConfig          *FilterConfig          `json:"filterConfig,omitempty"`
	ContinueOnError       bool                   `json:"continueOnError,omitempty"`
	RowTimeValueCheck     bool                   `json:"rowTimeValueCheck,omitempty"`
	StreamIngestionConfig *StreamIngestionConfig `json:"streamIngestionConfig,omitempty"`
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

type Task struct {
	TaskConfigMap map[string]string `json:"taskConfigMap"`
}

type Table struct {
	TableName        string                `json:"tableName"`
	TableType        string                `json:"tableType"`
	SegmentsConfig   TableSegmentsConfig   `json:"segmentsConfig"`
	Tenants          TableTenant           `json:"tenants"`
	TableIndexConfig TableIndexConfig      `json:"tableIndexConfig"`
	Metadata         *TableMetadata        `json:"metadata"`
	FieldConfigList  []FieldConfig         `json:"fieldConfigList,omitempty"`
	IngestionConfig  *TableIngestionConfig `json:"ingestionConfig,omitempty"`
	TierConfigs      []*TierConfig         `json:"tierConfigs,omitempty"`
	IsDimTable       bool                  `json:"isDimTable"`
	Query            *QueryConfig          `json:"query,omitempty"`
	Routing          *RoutingConfig        `json:"routing,omitempty"`
	Quota            *QuotaConfig          `json:"quota,omitempty"`
	UpsertConfig     *UpsertConfig         `json:"upsertConfig,omitempty"`
	DedupConfig      *DedupConfig          `json:"dedupConfig,omitempty"`
	Task             *Task                 `json:"task,omitempty"`
}

func (t *Table) IsEmpty() bool {

	return t.TableName == "" && t.TableType == "" && t.SegmentsConfig.TimeType == "" && t.Tenants.Broker == "" && t.Tenants.Server == "" && t.TableIndexConfig.LoadMode == "" && t.TableIndexConfig.ColumnMinMaxValueGeneratorMode == "" && t.TableIndexConfig.NoDictionarySizeRatioThreshold == 0 && t.TableIndexConfig.SegmentNameGeneratorType == "" && t.TableIndexConfig.OptimizeDictionary == false && t.TableIndexConfig.OptimizeDictionaryForMetrics == false && t.TableIndexConfig.NoDictionaryColumns == nil && t.TableIndexConfig.OnHeapDictionaryColumns == nil && t.TableIndexConfig.VarLengthDictionaryColumns == nil && t.TableIndexConfig.BloomFilterColumns == nil && t.TableIndexConfig.RangeIndexVersion == 0 && t.TableIndexConfig.CreateInvertedIndexDuringSegmentGeneration == false && t.TableIndexConfig.SortedColumn == nil && t.TableIndexConfig.StarTreeIndexConfigs == nil && t.TableIndexConfig.EnableDefaultStarTree == false && t.TableIndexConfig.EnableDynamicStarTreeCreation == false && t.TableIndexConfig.AggregateMetrics == false && t.TableIndexConfig.RangeIndexColumns == nil && t.TableIndexConfig.InvertedIndexColumns == nil && t.TableIndexConfig.SegmentPartitionConfig == nil && t.Metadata == nil && t.FieldConfigList == nil && t.IngestionConfig == nil && t.TierConfigs == nil && t.IsDimTable == false && t.Query == nil && t.Routing == nil && t.Quota == nil && t.UpsertConfig == nil && t.DedupConfig == nil
}
