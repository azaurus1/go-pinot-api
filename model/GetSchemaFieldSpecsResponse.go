package model

type AllowedDataType struct {
	NullDefault any `json:"nullDefault,omitempty"`
}

type FieldType struct {
	AllowedDataTypes map[string]AllowedDataType `json:"allowedDataTypes,omitempty"`
}

type DataType struct {
	StoredType string `json:"storedType,omitempty"`
	Size       int    `json:"size,omitempty"`
	Sortable   bool   `json:"sortable,omitempty"`
	Numeric    bool   `json:"numeric,omitempty"`
}

type GetSchemaFieldSpecsResponse struct {
	FieldTypes map[string]FieldType `json:"fieldTypes,omitempty"`
	DataTypes  map[string]DataType  `json:"dataTypes,omitempty"`
}
