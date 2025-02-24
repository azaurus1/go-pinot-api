package model

type TableSegmentsConfig struct {
	TimeType                       string            `json:"timeType"`
	Replication                    string            `json:"replication"`
	ReplicasPerPartition           string            `json:"replicasPerPartition,omitempty"`
	TimeColumnName                 string            `json:"timeColumnName"`
	SegmentAssignmentStrategy      string            `json:"segmentAssignmentStrategy,omitempty"`
	SegmentPushType                string            `json:"segmentPushType,omitempty"`
	MinimizeDataMovement           bool              `json:"minimizeDataMovement"`
	RetentionTimeUnit              string            `json:"retentionTimeUnit,omitempty"`
	RetentionTimeValue             string            `json:"retentionTimeValue,omitempty"`
	CrypterClassName               string            `json:"crypterClassName,omitempty"`
	PeerSegmentDownloadScheme      string            `json:"peerSegmentDownloadScheme,omitempty"`
	CompletionMode                 string            `json:"completionMode,omitempty"`
	DeletedSegmentsRetentionPeriod string            `json:"deletedSegmentsRetentionPeriod,omitempty"`
	SchemaName                     string            `json:"schemaName,omitempty"`
	CompletionConfig               *CompletionConfig `json:"completionConfig,omitempty"`
}

type CompletionConfig struct {
	CompletionMode string `json:"completionMode"`
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
	MetadataTTL                  float64           `json:"metadataTTL,omitempty"`
	DeletedKeysTTL               float64           `json:"deletedKeysTTL,omitempty"`
	HashFunction                 string            `json:"hashFunction,omitempty"`
	EnableSnapshot               *bool             `json:"enableSnapshot,omitempty"`
	EnablePreload                *bool             `json:"enablePreload,omitempty"`
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
	SegmentTimeValueCheck *bool                  `json:"segmentTimeValueCheck,omitempty"`
	TransformConfigs      []TransformConfig      `json:"transformConfigs,omitempty"`
	FilterConfig          *FilterConfig          `json:"filterConfig,omitempty"`
	ContinueOnError       bool                   `json:"continueOnError,omitempty"`
	RowTimeValueCheck     bool                   `json:"rowTimeValueCheck,omitempty"`
	StreamIngestionConfig *StreamIngestionConfig `json:"streamIngestionConfig,omitempty"`
}

type StreamIngestionConfig struct {
	StreamConfigMaps                 []StreamConfig `json:"streamConfigMaps"`
	ColumnMajorSegmentBuilderEnabled *bool          `json:"columnMajorSegmentBuilderEnabled,omitempty"`
	TrackFilteredMessageOffsets      *bool          `json:"trackFilteredMessageOffsets,omitempty"`
}

type StreamConfig struct {
	AccessKey                                                        string `json:"accessKey,omitempty"`
	AuthenticationType                                               string `json:"authentication.type,omitempty"`
	KeySerializer                                                    string `json:"key.serializer,omitempty"`
	MaxRecordsToFetch                                                string `json:"maxRecordsToFetch,omitempty"`
	RealtimeSegmentCommitTimeoutSeconds                              string `json:"realtime.segment.commit.timeoutSeconds,omitempty"`
	RealtimeSegmentFlushAutotuneInitialRows                          string `json:"realtime.segment.flush.autotune.initialRows,omitempty"`
	RealtimeSegmentFlushDesiredSize                                  string `json:"realtime.segment.flush.desired.size,omitempty"`
	RealtimeSegmentFlushThresholdRows                                string `json:"realtime.segment.flush.threshold.rows,omitempty"`
	RealtimeSegmentFlushThresholdTime                                string `json:"realtime.segment.flush.threshold.time,omitempty"`
	RealtimeSegmentFlushThresholdSegmentRows                         string `json:"realtime.segment.flush.threshold.segment.rows,omitempty"`
	RealtimeSegmentFlushThresholdSegmentSize                         string `json:"realtime.segment.flush.threshold.segment.size,omitempty"`
	RealtimeSegmentFlushThresholdSegmentTime                         string `json:"realtime.segment.flush.threshold.segment.time,omitempty"`
	RealtimeSegmentServerUploadToDeepStore                           string `json:"realtime.segment.serverUploadToDeepStore,omitempty"`
	Region                                                           string `json:"region,omitempty"`
	SaslJaasConfig                                                   string `json:"sasl.jaas.config,omitempty"`
	SaslMechanism                                                    string `json:"sasl.mechanism,omitempty"`
	SecretKey                                                        string `json:"secretKey,omitempty"`
	SecurityProtocol                                                 string `json:"security.protocol,omitempty"`
	ShardIteratorType                                                string `json:"shardIteratorType,omitempty"`
	StreamType                                                       string `json:"streamType,omitempty"`
	SslKeyPassword                                                   string `json:"ssl.key.password,omitempty"`
	SslKeystoreLocation                                              string `json:"ssl.keystore.location,omitempty"`
	SslKeystorePassword                                              string `json:"ssl.keystore.password,omitempty"`
	SslKeystoreType                                                  string `json:"ssl.keystore.type,omitempty"`
	SslTruststoreLocation                                            string `json:"ssl.truststore.location,omitempty"`
	SslTruststorePassword                                            string `json:"ssl.truststore.password,omitempty"`
	SslTruststoreType                                                string `json:"ssl.truststore.type,omitempty"`
	StreamKafkaBrokerList                                            string `json:"stream.kafka.broker.list,omitempty"`
	StreamKafkaBufferSize                                            string `json:"stream.kafka.buffer.size,omitempty"`
	StreamKafkaConsumerFactoryClassName                              string `json:"stream.kafka.consumer.factory.class.name,omitempty"`
	StreamKafkaConsumerPropAutoOffsetReset                           string `json:"stream.kafka.consumer.prop.auto.offset.reset,omitempty"`
	StreamKafkaConsumerType                                          string `json:"stream.kafka.consumer.type,omitempty"`
	StreamKafkaFetchTimeoutMillis                                    string `json:"stream.kafka.fetch.timeout.millis,omitempty"`
	StreamKafkaConnectionTimeoutMillis                               string `json:"stream.kafka.connection.timeout.millis,omitempty"`
	StreamKafkaDecoderClassName                                      string `json:"stream.kafka.decoder.class.name,omitempty"`
	StreamKafkaDecoderPropBasicAuthCredentialsSource                 string `json:"stream.kafka.decoder.prop.basic.auth.credentials.source,omitempty"`
	StreamKafkaDecoderPropDescriptorFile                             string `json:"stream.kafka.decoder.prop.descriptorFile,omitempty"`
	StreamKafkaDecoderPropProtoClassName                             string `json:"stream.kafka.decoder.prop.protoClassName,omitempty"`
	StreamKafkaDecoderPropFormat                                     string `json:"stream.kafka.decoder.prop.format,omitempty"`
	StreamKafkaDecoderPropSchemaRegistryBasicAuthUserInfo            string `json:"stream.kafka.decoder.prop.schema.registry.basic.auth.user.info,omitempty"`
	StreamKafkaDecoderPropSchemaRegistryBasicAuthCredentialsSource   string `json:"stream.kafka.decoder.prop.schema.registry.basic.auth.credentials.source,omitempty"`
	StreamKafkaDecoderPropSchemaRegistryRestUrl                      string `json:"stream.kafka.decoder.prop.schema.registry.rest.url,omitempty"`
	StreamKafkaDecoderPropSchemaRegistrySchemaName                   string `json:"stream.kafka.decoder.prop.schema.registry.schema.name,omitempty"`
	StreamKafkaDecoderPropSchemaRegistrySslKeystoreLocation          string `json:"stream.kafka.decoder.prop.schema.registry.ssl.keystore.location,omitempty"`
	StreamKafkaDecoderPropSchemaRegistrySslKeystorePassword          string `json:"stream.kafka.decoder.prop.schema.registry.ssl.keystore.password,omitempty"`
	StreamKafkaDecoderPropSchemaRegistrySslKeystoreType              string `json:"stream.kafka.decoder.prop.schema.registry.ssl.keystore.type,omitempty"`
	StreamKafkaDecoderPropSchemaRegistrySslTruststoreLocation        string `json:"stream.kafka.decoder.prop.schema.registry.ssl.truststore.location,omitempty"`
	StreamKafkaDecoderPropSchemaRegistrySslTruststorePassword        string `json:"stream.kafka.decoder.prop.schema.registry.ssl.truststore.password,omitempty"`
	StreamKafkaDecoderPropSchemaRegistrySslTruststoreType            string `json:"stream.kafka.decoder.prop.schema.registry.ssl.truststore.type,omitempty"`
	StreamKafkaDecoderPropSchemaRegistrySslProtocol                  string `json:"stream.kafka.decoder.prop.schema.registry.ssl.protocol,omitempty"`
	StreamKafkaFetcherMinBytes                                       string `json:"stream.kafka.fetcher.minBytes,omitempty"`
	StreamKafkaFetcherSize                                           string `json:"stream.kafka.fetcher.size,omitempty"`
	StreamKafkaHlcGroupId                                            string `json:"stream.kafka.hlc.group.id,omitempty"`
	StreamKafkaIdleTimeoutMillis                                     string `json:"stream.kafka.idle.timeout.millis,omitempty"`
	StreamKafkaIsolationLevel                                        string `json:"stream.kafka.isolation.level,omitempty"`
	StreamKafkaMetadataPopulate                                      string `json:"stream.kafka.metadata.populate,omitempty"`
	StreamKafkaSchemaRegistryUrl                                     string `json:"stream.kafka.schema.registry.url,omitempty"`
	StreamKafkaSocketTimeout                                         string `json:"stream.kafka.socket.timeout,omitempty"`
	StreamKafkaSslCertificateType                                    string `json:"stream.kafka.ssl.certificate.type,omitempty"`
	StreamKafkaSslClientCertificate                                  string `json:"stream.kafka.ssl.client.certificate,omitempty"`
	StreamKafkaSslClientKey                                          string `json:"stream.kafka.ssl.client.key,omitempty"`
	StreamKafkaSslClientKeyAlgorithm                                 string `json:"stream.kafka.ssl.client.key.algorithm,omitempty"`
	StreamKafkaSslServerCertificate                                  string `json:"stream.kafka.ssl.server.certificate,omitempty"`
	StreamKafkaTopicName                                             string `json:"stream.kafka.topic.name,omitempty"`
	StreamKafkaZkBrokerUrl                                           string `json:"stream.kafka.zk.broker.url,omitempty"`
	StreamKinesisConsumerFactoryClassName                            string `json:"stream.kinesis.consumer.factory.class.name,omitempty"`
	StreamKinesisConsumerType                                        string `json:"stream.kinesis.consumer.type,omitempty"`
	StreamKinesisDecoderClassName                                    string `json:"stream.kinesis.decoder.class.name,omitempty"`
	StreamKinesisFetchTimeoutMillis                                  string `json:"stream.kinesis.fetch.timeout.millis,omitempty"`
	StreamKinesisDecoderPropSchemaRegistryBasicAuthUserInfo          string `json:"stream.kinesis.decoder.prop.schema.registry.basic.auth.user.info,omitempty"`
	StreamKinesisDecoderPropSchemaRegistryBasicAuthCredentialsSource string `json:"stream.kinesis.decoder.prop.schema.registry.basic.auth.credentials.source,omitempty"`
	StreamKinesisDecoderPropSchemaRegistryRestUrl                    string `json:"stream.kinesis.decoder.prop.schema.registry.rest.url,omitempty"`
	StreamKinesisDecoderPropSchemaRegistrySchemaName                 string `json:"stream.kinesis.decoder.prop.schema.registry.schema.name,omitempty"`
	StreamKinesisDecoderPropSchemaRegistrySslKeystoreLocation        string `json:"stream.kinesis.decoder.prop.schema.registry.ssl.keystore.location,omitempty"`
	StreamKinesisDecoderPropSchemaRegistrySslKeystorePassword        string `json:"stream.kinesis.decoder.prop.schema.registry.ssl.keystore.password,omitempty"`
	StreamKinesisDecoderPropSchemaRegistrySslKeystoreType            string `json:"stream.kinesis.decoder.prop.schema.registry.ssl.keystore.type,omitempty"`
	StreamKinesisDecoderPropSchemaRegistrySslTruststoreLocation      string `json:"stream.kinesis.decoder.prop.schema.registry.ssl.truststore.location,omitempty"`
	StreamKinesisDecoderPropSchemaRegistrySslTruststorePassword      string `json:"stream.kinesis.decoder.prop.schema.registry.ssl.truststore.password,omitempty"`
	StreamKinesisDecoderPropSchemaRegistrySslTruststoreType          string `json:"stream.kinesis.decoder.prop.schema.registry.ssl.truststore.type,omitempty"`
	StreamKinesisDecoderPropSchemaRegistrySslProtocol                string `json:"stream.kinesis.decoder.prop.schema.registry.ssl.protocol,omitempty"`
	StreamKinesisTopicName                                           string `json:"stream.kinesis.topic.name,omitempty"`
	StreamPulsarAudience                                             string `json:"stream.pulsar.audience,omitempty"`
	StreamPulsarAuthenticationToken                                  string `json:"stream.pulsar.authenticationToken,omitempty"`
	StreamPulsarBootstrapServers                                     string `json:"stream.pulsar.bootstrap.servers,omitempty"`
	StreamPulsarCredsFilePath                                        string `json:"stream.pulsar.credsFilePath,omitempty"`
	StreamPulsarConsumerFactoryClassName                             string `json:"stream.pulsar.consumer.factory.class.name,omitempty"`
	StreamPulsarConsumerPropAutoOffsetReset                          string `json:"stream.pulsar.consumer.prop.auto.offset.reset,omitempty"`
	StreamPulsarConsumerType                                         string `json:"stream.pulsar.consumer.type,omitempty"`
	StreamPulsarDecoderClassName                                     string `json:"stream.pulsar.decoder.class.name,omitempty"`
	StreamPulsarFetchTimeoutMillis                                   string `json:"stream.pulsar.fetch.timeout.millis,omitempty"`
	StreamPulsarIssuerUrl                                            string `json:"stream.pulsar.issuerUrl,omitempty"`
	StreamPulsarMetadataPopulate                                     string `json:"stream.pulsar.metadata.populate,omitempty"`
	StreamPulsarMetadataFields                                       string `json:"stream.pulsar.metadata.fields,omitempty"`
	StreamPulsarDecoderPropSchemaRegistryBasicAuthUserInfo           string `json:"stream.pulsar.decoder.prop.schema.registry.basic.auth.user.info,omitempty"`
	StreamPulsarDecoderPropSchemaRegistryBasicAuthCredentialsSource  string `json:"stream.pulsar.decoder.prop.schema.registry.basic.auth.credentials.source,omitempty"`
	StreamPulsarDecoderPropSchemaRegistryRestUrl                     string `json:"stream.pulsar.decoder.prop.schema.registry.rest.url,omitempty"`
	StreamPulsarDecoderPropSchemaRegistrySchemaName                  string `json:"stream.pulsar.decoder.prop.schema.registry.schema.name,omitempty"`
	StreamPulsarDecoderPropSchemaRegistrySslKeystoreLocation         string `json:"stream.pulsar.decoder.prop.schema.registry.ssl.keystore.location,omitempty"`
	StreamPulsarDecoderPropSchemaRegistrySslKeystorePassword         string `json:"stream.pulsar.decoder.prop.schema.registry.ssl.keystore.password,omitempty"`
	StreamPulsarDecoderPropSchemaRegistrySslKeystoreType             string `json:"stream.pulsar.decoder.prop.schema.registry.ssl.keystore.type,omitempty"`
	StreamPulsarDecoderPropSchemaRegistrySslTruststoreLocation       string `json:"stream.pulsar.decoder.prop.schema.registry.ssl.truststore.location,omitempty"`
	StreamPulsarDecoderPropSchemaRegistrySslTruststorePassword       string `json:"stream.pulsar.decoder.prop.schema.registry.ssl.truststore.password,omitempty"`
	StreamPulsarDecoderPropSchemaRegistrySslTruststoreType           string `json:"stream.pulsar.decoder.prop.schema.registry.ssl.truststore.type,omitempty"`
	StreamPulsarDecoderPropSchemaRegistrySslProtocol                 string `json:"stream.pulsar.decoder.prop.schema.registry.ssl.protocol,omitempty"`
	StreamPulsarTlsTrustCertsFilePath                                string `json:"stream.pulsar.tlsTrustCertsFilePath,omitempty"`
	StreamPulsarTopicName                                            string `json:"stream.pulsar.topic.name,omitempty"`
	TopicConsumptionRateLimit                                        string `json:"topic.consumption.rate.limit,omitempty"`
	ValueSerializer                                                  string `json:"value.serializer,omitempty"`
}

type TierConfig struct {
	Name                string `json:"name"`
	SegmentSelectorType string `json:"segmentSelectorType"`
	SegmentAge          string `json:"segmentAge"`
	StorageType         string `json:"storageType"`
	ServerTag           string `json:"serverTag"`
}

// map[string]map[string]string is like
// {"field": {"inner_nested_fied": "value"}}
type Task struct {
	TaskTypeConfigsMap map[string]map[string]string `json:"taskTypeConfigsMap"`
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
	return t.TableName == "" && t.TableType == "" && t.SegmentsConfig.TimeType == "" && t.Tenants.Broker == "" && t.Tenants.Server == "" && t.TableIndexConfig.LoadMode == "" && t.TableIndexConfig.ColumnMinMaxValueGeneratorMode == "" && t.TableIndexConfig.NoDictionarySizeRatioThreshold == 0 && t.TableIndexConfig.SegmentNameGeneratorType == "" && !t.TableIndexConfig.OptimizeDictionary && !t.TableIndexConfig.OptimizeDictionaryForMetrics && t.TableIndexConfig.NoDictionaryColumns == nil && t.TableIndexConfig.OnHeapDictionaryColumns == nil && t.TableIndexConfig.VarLengthDictionaryColumns == nil && t.TableIndexConfig.BloomFilterColumns == nil && t.TableIndexConfig.RangeIndexVersion == 0 && !t.TableIndexConfig.CreateInvertedIndexDuringSegmentGeneration && t.TableIndexConfig.SortedColumn == nil && t.TableIndexConfig.StarTreeIndexConfigs == nil && !t.TableIndexConfig.EnableDefaultStarTree && !t.TableIndexConfig.EnableDynamicStarTreeCreation && !t.TableIndexConfig.AggregateMetrics && t.TableIndexConfig.RangeIndexColumns == nil && t.TableIndexConfig.InvertedIndexColumns == nil && t.TableIndexConfig.SegmentPartitionConfig == nil && t.Metadata == nil && t.FieldConfigList == nil && t.IngestionConfig == nil && t.TierConfigs == nil && !t.IsDimTable && t.Query == nil && t.Routing == nil && t.Quota == nil && t.UpsertConfig == nil && t.DedupConfig == nil
}
