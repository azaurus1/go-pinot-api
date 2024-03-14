package config_templating

import (
	"encoding/json"
	"github.com/azaurus1/go-pinot-api/model"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestTemplater(t *testing.T) {

	t.Run("TemplateTableConfigToFile Params correctly injected into template", func(t *testing.T) {

		// create temp file
		f, err := os.CreateTemp("", "test")
		if err != nil {
			t.Fatal(err)
		}

		defer os.Remove(f.Name())

		// write to it
		config := dummyTableConfig()
		configBytes, err := json.Marshal(config)
		assert.NoError(t, err)

		_, err = f.Write(configBytes)
		assert.NoError(t, err)

		// call TemplateTableConfig

		params := TableConfigTemplateParameters{
			PinotSegmentsReplication: "1",
			PinotTenantBroker:        "custom-tenant1",
			PinotTenantServer:        "custom-tenant1",
			KafkaBrokers:             "localhost:9092",
			KafkaTopic:               "test",
			KafkaSaslUsername:        "user",
			KafkaSaslPassword:        "password",
			KafkaSaslMechanism:       "PLAIN",
			KafkaSecurityProtocol:    "SASL_PLAINTEXT",
			SchemaRegistryUrl:        "http://localhost:8081",
			SchemaRegistryUsername:   "user_sr",
			SchemaRegistryPassword:   "password_sr",
		}

		err = TemplateTableConfigToFile(f.Name(), f.Name(), params)
		assert.NoError(t, err)

		f, err = os.Open(f.Name())
		if err != nil {
			t.Fatal(err)
		}

		defer f.Close()

		var templatedConfig model.Table
		err = json.NewDecoder(f).Decode(&templatedConfig)
		assert.NoError(t, err)

		// assert that the template was correctly injected
		assert.Equal(t, params.PinotSegmentsReplication, templatedConfig.SegmentsConfig.Replication)
		assert.Equal(t, params.PinotTenantBroker, templatedConfig.Tenants.Broker)
		assert.Equal(t, params.PinotTenantServer, templatedConfig.Tenants.Server)
		assert.Equal(t, params.KafkaBrokers, templatedConfig.IngestionConfig.StreamIngestionConfig.StreamConfigMaps[0]["stream.kafka.broker.list"])
		assert.Equal(t, params.KafkaTopic, templatedConfig.IngestionConfig.StreamIngestionConfig.StreamConfigMaps[0]["stream.kafka.topic.name"])
		assert.Equal(t, params.SchemaRegistryUrl, templatedConfig.IngestionConfig.StreamIngestionConfig.StreamConfigMaps[0]["stream.kafka.decoder.prop.schema.registry.rest.url"])
		assert.Equal(t, params.SchemaRegistryUsername+":"+params.SchemaRegistryPassword, templatedConfig.IngestionConfig.StreamIngestionConfig.StreamConfigMaps[0]["stream.kafka.decoder.prop.schema.registry.basic.auth.user.info"])
		assert.Equal(t, params.KafkaSecurityProtocol, templatedConfig.IngestionConfig.StreamIngestionConfig.StreamConfigMaps[0]["security.protocol"])
		assert.Equal(t, params.KafkaSaslMechanism, templatedConfig.IngestionConfig.StreamIngestionConfig.StreamConfigMaps[0]["sasl.mechanism"])
		assert.Equal(t, "org.apache.kafka.common.security.plain.PlainLoginModule required \n username=\""+params.KafkaSaslUsername+"\" \n password=\""+params.KafkaSaslPassword+"\";", templatedConfig.IngestionConfig.StreamIngestionConfig.StreamConfigMaps[0]["sasl.jaas.config"])

	})

	t.Run("TemplateTableConfig correctly renders template", func(t *testing.T) {
		config := dummyTableConfig()
		configBytes, err := json.Marshal(config)
		assert.NoError(t, err)

		params := TableConfigTemplateParameters{
			PinotSegmentsReplication: "1",
			PinotTenantBroker:        "custom-tenant1",
			PinotTenantServer:        "custom-tenant1",
			KafkaBrokers:             "localhost:9092",
			KafkaTopic:               "test",
			KafkaSaslUsername:        "user",
			KafkaSaslPassword:        "password",
			KafkaSaslMechanism:       "PLAIN",
			KafkaSecurityProtocol:    "SASL_PLAINTEXT",
			SchemaRegistryUrl:        "http://localhost:8081",
			SchemaRegistryUsername:   "user_sr",
			SchemaRegistryPassword:   "password_sr",
		}

		templatedConfig, err := TemplateTableConfig(configBytes, params)
		assert.NoError(t, err)

		// assert that the template was correctly injected
		assert.Equal(t, params.PinotSegmentsReplication, templatedConfig.SegmentsConfig.Replication)
		assert.Equal(t, params.PinotTenantBroker, templatedConfig.Tenants.Broker)
		assert.Equal(t, params.PinotTenantServer, templatedConfig.Tenants.Server)
		assert.Equal(t, params.KafkaBrokers, templatedConfig.IngestionConfig.StreamIngestionConfig.StreamConfigMaps[0]["stream.kafka.broker.list"])
		assert.Equal(t, params.KafkaTopic, templatedConfig.IngestionConfig.StreamIngestionConfig.StreamConfigMaps[0]["stream.kafka.topic.name"])
		assert.Equal(t, params.SchemaRegistryUrl, templatedConfig.IngestionConfig.StreamIngestionConfig.StreamConfigMaps[0]["stream.kafka.decoder.prop.schema.registry.rest.url"])
		assert.Equal(t, params.SchemaRegistryUsername+":"+params.SchemaRegistryPassword, templatedConfig.IngestionConfig.StreamIngestionConfig.StreamConfigMaps[0]["stream.kafka.decoder.prop.schema.registry.basic.auth.user.info"])
		assert.Equal(t, params.KafkaSecurityProtocol, templatedConfig.IngestionConfig.StreamIngestionConfig.StreamConfigMaps[0]["security.protocol"])
		assert.Equal(t, params.KafkaSaslMechanism, templatedConfig.IngestionConfig.StreamIngestionConfig.StreamConfigMaps[0]["sasl.mechanism"])
		assert.Equal(t, "org.apache.kafka.common.security.plain.PlainLoginModule required \n username=\""+params.KafkaSaslUsername+"\" \n password=\""+params.KafkaSaslPassword+"\";", templatedConfig.IngestionConfig.StreamIngestionConfig.StreamConfigMaps[0]["sasl.jaas.config"])

	})

}

func dummyTableConfig() model.Table {
	return model.Table{
		TableName: "dummy",
		TableType: "REALTIME",
		SegmentsConfig: model.TableSegmentsConfig{
			Replication: "{{ .PinotSegmentsReplication }}",
		},
		Tenants: model.TableTenant{
			Broker: "{{ .PinotTenantBroker }}",
			Server: "{{ .PinotTenantServer }}",
		},
		IngestionConfig: model.TableIngestionConfig{
			StreamIngestionConfig: model.StreamIngestionConfig{
				StreamConfigMaps: []map[string]string{{
					"stream.kafka.topic.name":                                                 "{{ .KafkaTopic }}",
					"stream.kafka.broker.list":                                                "{{ .KafkaBrokers }}",
					"stream.kafka.decoder.prop.schema.registry.rest.url":                      "{{ .SchemaRegistryUrl }}",
					"stream.kafka.decoder.prop.schema.registry.basic.auth.credentials.source": "USER_INFO",
					"stream.kafka.decoder.prop.schema.registry.basic.auth.user.info":          "{{ .SchemaRegistryUsername }}:{{ .SchemaRegistryPassword }}",
					"stream.kafka.decoder.prop.format":                                        "AVRO",
					"authentication.type":                                                     "SASL",
					"security.protocol":                                                       "{{ .KafkaSecurityProtocol }}",
					"sasl.mechanism":                                                          "{{ .KafkaSaslMechanism }}",
					"sasl.jaas.config":                                                        "org.apache.kafka.common.security.plain.PlainLoginModule required \n username=\"{{ .KafkaSaslUsername }}\" \n password=\"{{ .KafkaSaslPassword }}\";",
				},
				},
			},
		},
	}
}
