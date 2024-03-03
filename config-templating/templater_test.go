package config_templating

import (
	"github.com/azaurus1/go-pinot-api/model"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestTemplater(t *testing.T) {

	t.Run("Params correctly injected into template", func(t *testing.T) {

		// create temp file
		f, err := os.CreateTemp("", "test")
		if err != nil {
			t.Fatal(err)
		}

		defer os.Remove(f.Name())

		// write to it
		config := dummyTableConfig()
		configBytes, err := config.AsBytes()
		assert.NoError(t, err)

		_, err = f.Write(configBytes)
		assert.NoError(t, err)

		// call TemplateTableConfig
		params := TableConfigTemplateParameters{
			KafkaBrokers:      "localhost:9092",
			KafkaTopic:        "test",
			SchemaRegistryUrl: "http://localhost:8081",
		}

		err = TemplateTableConfig(f.Name(), f.Name(), params)
		assert.NoError(t, err)

		// read from it
		templatedConfig, err := model.NewTableConfigFromFile(f.Name())
		assert.NoError(t, err)

		// assert that the template was correctly injected
		assert.Equal(t, "localhost:9092", templatedConfig.TableIndexConfig.StreamConfigs["stream.kafka.broker.list"])
		assert.Equal(t, "test", templatedConfig.TableIndexConfig.StreamConfigs["stream.kafka.topic.name"])
		assert.Equal(t, "http://localhost:8081", templatedConfig.TableIndexConfig.StreamConfigs["stream.kafka.decoder.prop.schema.registry.rest.url"])

	})

}

func dummyTableConfig() model.TableConfig {
	return model.TableConfig{
		TableName: "dummy",
		TableType: "REALTIME",
		TableIndexConfig: model.TableIndexConfig{
			StreamConfigs: map[string]string{
				"stream.kafka.topic.name":                            "{{ .KafkaTopic }}",
				"stream.kafka.broker.list":                           "{{ .KafkaBrokers }}",
				"stream.kafka.decoder.prop.schema.registry.rest.url": "{{ .SchemaRegistryUrl }}",
			},
		},
	}
}
