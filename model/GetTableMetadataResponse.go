package model

type GetTableMetadataResponse struct {
	TableName                                 string           `json:"tableName"`
	DiskSizeInBytes                           int64            `json:"diskSizeInBytes"`
	NumSegments                               int              `json:"numSegments"`
	NumRows                                   int64            `json:"numRows"`
	ColumnLengthMap                           map[string]int   `json:"columnLengthMap"`
	ColumnCardinalityMap                      map[string]int   `json:"columnCardinalityMap"`
	MaxNumMultiValuesMap                      map[string]int   `json:"maxNumMultiValuesMap"`
	ColumnIndexSizeMap                        map[string]int64 `json:"columnIndexSizeMap"`
	UpsertPartitionToServerPrimaryKeyCountMap map[string]int64 `json:"upsertPartitionToServerPrimaryKeyCountMap"`
}
