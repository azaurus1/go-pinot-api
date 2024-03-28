package model

import "encoding/json"

type FieldSpec struct {
	Name             string `json:"name"`
	DataType         string `json:"dataType"`
	Format           string `json:"format,omitempty"`
	Granularity      string `json:"granularity,omitempty"`
	NotNull          *bool  `json:"notNull,omitempty"`
	SingleValueField *bool  `json:"singleValueField,omitempty"`
}

type Schema struct {
	SchemaName          string      `json:"schemaName"`
	DimensionFieldSpecs []FieldSpec `json:"dimensionFieldSpecs"`
	MetricFieldSpecs    []FieldSpec `json:"metricFieldSpecs,omitempty"`
	DateTimeFieldSpecs  []FieldSpec `json:"dateTimeFieldSpecs,omitempty"`
	PrimaryKeyColumns   []string    `json:"primaryKeyColumns,omitempty"`
}

func (schema *Schema) AsBytes() ([]byte, error) {

	schemaBytes, err := json.Marshal(schema)
	if err != nil {
		return nil, err
	}

	return schemaBytes, nil

}

func (schema *Schema) String() string {
	jsonString, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		return err.Error()
	}
	return string(jsonString)
}
