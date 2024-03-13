package config_templating

import (
	"bytes"
	"fmt"
	"os"
	"text/template"
)

type TableConfigTemplateParameters struct {
	KafkaBrokers      string
	KafkaTopic        string
	SchemaRegistryUrl string
}

func TemplateTableConfigToFile(inputTemplate, outputConfigFile string, params TableConfigTemplateParameters) error {

	tableConfigTpl, err := os.ReadFile(inputTemplate)
	if err != nil {
		return err
	}

	outputConfigBytes, err := templateTableConfigToBytes(tableConfigTpl, params)
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

func templateTableConfigToBytes(inputTemplate []byte, params TableConfigTemplateParameters) ([]byte, error) {

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
