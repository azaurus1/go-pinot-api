package model

import "encoding/json"

type FieldSpec struct {
	Name        string `json:"name"`
	DataType    string `json:"dataType"`
	Format      string `json:"format,omitempty"`
	Granularity string `json:"granularity,omitempty"`
}

type Schema struct {
	SchemaName          string      `json:"schemaName"`
	DimensionFieldSpecs []FieldSpec `json:"dimensionFieldSpecs"`
	MetricFieldSpecs    []FieldSpec `json:"metricFieldSpecs"`
	DateTimeFieldSpecs  []FieldSpec `json:"dateTimeFieldSpecs"`
	PrimaryKeyColumns   []string    `json:"primaryKeyColumns"`
}

func (schema *Schema) AsBytes() ([]byte, error) {

	schemaBytes, err := json.Marshal(schema)
	if err != nil {
		return nil, err
	}

	return schemaBytes, nil

}
