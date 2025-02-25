package config_templating

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/azaurus1/go-pinot-api/model"
	"github.com/stretchr/testify/assert"
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
		assert.Equal(t, params.KafkaBrokers, templatedConfig.IngestionConfig.StreamIngestionConfig.StreamConfigMaps[0].StreamKafkaBrokerList)
		assert.Equal(t, params.KafkaTopic, templatedConfig.IngestionConfig.StreamIngestionConfig.StreamConfigMaps[0].StreamKafkaTopicName)
		assert.Equal(t, params.SchemaRegistryUrl, templatedConfig.IngestionConfig.StreamIngestionConfig.StreamConfigMaps[0].StreamKafkaDecoderPropSchemaRegistryRestUrl)
		assert.Equal(t, params.SchemaRegistryUsername+":"+params.SchemaRegistryPassword, templatedConfig.IngestionConfig.StreamIngestionConfig.StreamConfigMaps[0].StreamKafkaDecoderPropSchemaRegistryBasicAuthUserInfo)
		assert.Equal(t, params.KafkaSecurityProtocol, templatedConfig.IngestionConfig.StreamIngestionConfig.StreamConfigMaps[0].SecurityProtocol)
		assert.Equal(t, params.KafkaSaslMechanism, templatedConfig.IngestionConfig.StreamIngestionConfig.StreamConfigMaps[0].SaslMechanism)
		assert.Equal(t, "org.apache.kafka.common.security.plain.PlainLoginModule required \n username=\""+params.KafkaSaslUsername+"\" \n password=\""+params.KafkaSaslPassword+"\";", templatedConfig.IngestionConfig.StreamIngestionConfig.StreamConfigMaps[0].SaslJaasConfig)

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
		assert.Equal(t, params.KafkaBrokers, templatedConfig.IngestionConfig.StreamIngestionConfig.StreamConfigMaps[0].StreamKafkaBrokerList)
		assert.Equal(t, params.KafkaTopic, templatedConfig.IngestionConfig.StreamIngestionConfig.StreamConfigMaps[0].StreamKafkaTopicName)
		assert.Equal(t, params.SchemaRegistryUrl, templatedConfig.IngestionConfig.StreamIngestionConfig.StreamConfigMaps[0].StreamKafkaDecoderPropSchemaRegistryRestUrl)
		assert.Equal(t, params.SchemaRegistryUsername+":"+params.SchemaRegistryPassword, templatedConfig.IngestionConfig.StreamIngestionConfig.StreamConfigMaps[0].StreamKafkaDecoderPropSchemaRegistryBasicAuthUserInfo)
		assert.Equal(t, params.KafkaSecurityProtocol, templatedConfig.IngestionConfig.StreamIngestionConfig.StreamConfigMaps[0].SecurityProtocol)
		assert.Equal(t, params.KafkaSaslMechanism, templatedConfig.IngestionConfig.StreamIngestionConfig.StreamConfigMaps[0].SaslMechanism)
		assert.Equal(t, "org.apache.kafka.common.security.plain.PlainLoginModule required \n username=\""+params.KafkaSaslUsername+"\" \n password=\""+params.KafkaSaslPassword+"\";", templatedConfig.IngestionConfig.StreamIngestionConfig.StreamConfigMaps[0].SaslJaasConfig)

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
		IngestionConfig: &model.TableIngestionConfig{
			StreamIngestionConfig: &model.StreamIngestionConfig{
				StreamConfigMaps: []model.StreamConfig{
					{
						StreamKafkaTopicName:                                           "{{ .KafkaTopic }}",
						StreamKafkaBrokerList:                                          "{{ .KafkaBrokers }}",
						StreamKafkaDecoderPropSchemaRegistryRestUrl:                    "{{ .SchemaRegistryUrl }}",
						StreamKafkaDecoderPropSchemaRegistryBasicAuthCredentialsSource: "USER_INFO",
						StreamKafkaDecoderPropSchemaRegistryBasicAuthUserInfo:          "{{ .SchemaRegistryUsername }}:{{ .SchemaRegistryPassword }}",
						StreamKafkaDecoderPropFormat:                                   "AVRO",
						AuthenticationType:                                             "SASL",
						SecurityProtocol:                                               "{{ .KafkaSecurityProtocol }}",
						SaslMechanism:                                                  "{{ .KafkaSaslMechanism }}",
						SaslJaasConfig:                                                 "org.apache.kafka.common.security.plain.PlainLoginModule required \n username=\"{{ .KafkaSaslUsername }}\" \n password=\"{{ .KafkaSaslPassword }}\";",
					},
				},
			},
		},
	}
}

func dummyRawTable() string {
	return `{
  "tableName": "realtime_ethereum_mainnet_block_headers_REALTIME",
  "tableType": "REALTIME",
  "segmentsConfig": {
    "timeColumnName": "timestamp",
    "timeType": "MILLISECONDS",
    "replication": "1",
    "segmentAssignmentStrategy": "BalanceNumSegmentAssignmentStrategy",
    "segmentPushType": "APPEND",
    "minimizeDataMovement": true
  },
  "tenants": {
    "broker": "DefaultTenant",
    "server": "DefaultTenant"
  },
  "tableIndexConfig": {
    "rangeIndexVersion": 0,
    "autoGeneratedInvertedIndex": false,
    "createInvertedIndexDuringSegmentGeneration": false,
    "loadMode": "MMAP",
    "enableDefaultStarTree": false,
    "tierOverwrites": {
      "hotTier": {
        "starTreeIndexConfigs": null
      },
      "coldTier": {
        "starTreeIndexConfigs": null
      }
    },
    "enableDynamicStarTreeCreation": false,
    "aggregateMetrics": false,
    "nullHandlingEnabled": false,
    "columnMajorSegmentBuilderEnabled": false,
    "optimizeDictionary": false,
    "optimizeDictionaryForMetrics": false,
    "noDictionarySizeRatioThreshold": 0
  },
  "metadata": {
    "customConfigs": {
      "customKey": "customValue"
    }
  },
  "fieldConfigList": [
    {
      "name": "number",
      "encodingType": "RAW",
      "indexType": "SORTED",
      "indexTypes": [
        "SORTED"
      ],
      "timestampConfig": {},
      "indexes": {
        "inverted": {
          "enabled": ""
        }
      },
      "tierOverwrites": null
    }
  ],
  "ingestionConfig": {
    "segmentTimeValueCheck": true,
    "streamIngestionConfig": {
      "streamConfigMaps": [
        {
          "streamType": "kafka",
          "stream.kafka.consumer.type": "high-level",
          "stream.kafka.topic.name": "ethereum_mainnet_block_headers",
          "stream.kafka.decoder.class.name": "org.apache.pinot.plugin.stream.kafka.KafkaJSONMessageDecoder",
          "stream.kafka.consumer.factory.class.name": "org.apache.pinot.plugin.stream.kafka20.KafkaConsumerFactory",
          "stream.kafka.broker.list": "kafka:9092",
          "stream.kafka.zk.broker.url": "kafka:2181",
          "stream.kafka.consumer.prop.auto.offset.reset": "smallest",
          "stream.kafka.decoder.prop.schema.registry.rest.url": "http://schema-registry:8081",
          "stream.kafka.decoder.prop.schema.registry.schema.name": "ethereum_mainnet_block_headers-value",
          "stream.kafka.decoder.prop.schema.registry.schema.version": "latest"
        }
      ],
      "columnMajorSegmentBuilderEnabled": false
    },
    "continueOnError": true,
    "rowTimeValueCheck": true
  },
  "tierConfigs": [
    {
      "name": "hotTier",
      "segmentSelectorType": "time",
      "segmentAge": "3130d",
      "storageType": "pinot_server",
      "serverTag": "DefaultTenant_REALTIME"
    }
  ],
  "isDimTable": false
}`
}
