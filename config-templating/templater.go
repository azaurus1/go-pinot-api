package config_templating

import (
	"fmt"
	"os"
	"text/template"
)

type TableConfigTemplateParameters struct {
	KafkaBrokers      string
	KafkaTopic        string
	SchemaRegistryUrl string
}

func TemplateTableConfig(inputTemplate, outputConfigFile string, params TableConfigTemplateParameters) error {

	tableConfigTpl, err := os.ReadFile(inputTemplate)
	if err != nil {
		return err
	}

	tableConfigTemplate, err := template.New(outputConfigFile).Parse(string(tableConfigTpl))
	if err != nil {
		return err
	}

	outputFile, err := os.Create(outputConfigFile)
	if err != nil {
		return fmt.Errorf("error creating output file %s: %s", outputConfigFile, err)
	}

	defer outputFile.Close()

	err = tableConfigTemplate.Execute(outputFile, params)
	if err != nil {
		return err
	}

	return nil

}
