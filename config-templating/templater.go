package config_templating

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/azaurus1/go-pinot-api/model"
	"os"
	"text/template"
)

type TableConfigTemplateParameters struct {
	PinotSegmentsReplication string
	PinotTenantBroker        string
	PinotTenantServer        string
	KafkaBrokers             string
	KafkaTopic               string
	KafkaSaslUsername        string
	KafkaSaslPassword        string
	KafkaSaslMechanism       string
	KafkaSecurityProtocol    string
	SchemaRegistryUrl        string
	SchemaRegistryUsername   string
	SchemaRegistryPassword   string
}

func TemplateTableConfigToFile(inputTemplate, outputConfigFile string, params TableConfigTemplateParameters) error {

	tableConfigTpl, err := os.ReadFile(inputTemplate)
	if err != nil {
		return err
	}

	outputConfigBytes, err := templateTableConfig(tableConfigTpl, params)
	if err != nil {
		return err
	}

	outputFile, err := os.Create(outputConfigFile)
	if err != nil {
		return fmt.Errorf("error creating output file %s: %s", outputConfigFile, err)
	}

	defer outputFile.Close()

	_, err = outputFile.Write(outputConfigBytes)
	if err != nil {
		return fmt.Errorf("error writing to output file %s: %s", outputConfigFile, err)
	}

	return nil

}

func TemplateTableConfig(inputTemplate []byte, params TableConfigTemplateParameters) (*model.Table, error) {

	renderedBytes, err := templateTableConfig(inputTemplate, params)
	if err != nil {
		return nil, err
	}

	var table model.Table
	err = json.Unmarshal(renderedBytes, &table)
	if err != nil {
		return nil, err
	}

	return &table, nil
}

func templateTableConfig(inputTemplate []byte, params TableConfigTemplateParameters) ([]byte, error) {

	tableConfigTemplate, err := template.New("tableConfig").Parse(string(inputTemplate))
	if err != nil {
		return nil, err
	}

	var outputConfig bytes.Buffer

	err = tableConfigTemplate.Execute(&outputConfig, params)
	if err != nil {
		return nil, err
	}

	return outputConfig.Bytes(), nil

}
