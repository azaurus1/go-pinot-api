package model

type ColumnToIndex struct {
	VectorIndex     int `json:"vector_index,omitempty"`
	NullValueVector int `json:"nullvalue_vector,omitempty"`
	H3Index         int `json:"h3_index,omitempty"`
	Dictionary      int `json:"dictionary,omitempty"`
	JsonIndex       int `json:"json_index,omitempty"`
	RangeIndex      int `json:"range_index,omitempty"`
	ForwardIndex    int `json:"forward_index,omitempty"`
	BloomFilter     int `json:"bloom_filter,omitempty"`
	InvertedIndex   int `json:"inverted_index,omitempty"`
	TextIndex       int `json:"text_index,omitempty"`
	FSTIndex        int `json:"fst_index,omitempty"`
}

type GetTableIndexesResponse struct {
	TotalOnlineSegments  int                      `json:"totalOnlineSegments,omitempty"`
	ColumnToIndexesCount map[string]ColumnToIndex `json:"columnToIndexesCount,omitempty"`
}
