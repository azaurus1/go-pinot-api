package goPinotAPI_test

import (
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	goPinotAPI "github.com/azaurus1/go-pinot-api"
	model "github.com/azaurus1/go-pinot-api/model"
	"github.com/stretchr/testify/assert"
)

const (
	RouteUsers                                        = "/users"
	RouteUser                                         = "/users/test"
	RouteInstances                                    = "/instances"
	RouteInstance                                     = "/instances/Minion_172.19.0.2_9514"
	RouteClusterInfo                                  = "/cluster/info"
	RouteClusterConfigs                               = "/cluster/configs"
	RouteClusterConfigsDelete                         = "/cluster/configs/allowParticipantAutoJoin"
	RouteTenants                                      = "/tenants"
	RouteTenantsInstances                             = "/tenants/DefaultTenant"
	RouteTenantsTables                                = "/tenants/DefaultTenant/tables"
	RouteTenantsMetadata                              = "/tenants/DefaultTenant/metadata"
	RouteSegmentsTest                                 = "/segments/test"
	RouteSegmentsTestReload                           = "/segments/test/reload"
	RouteSegmentTestReload                            = "/segments/test/test_1/reload"
	RouteSegmentTestResetAll                          = "/segments/test_OFFLINE/reset"
	RouteSegmentTestReset                             = "/segments/test_OFFLINE/test_OFFLINE_16071_16071_0/reset"
	RouteV2Segments                                   = "/v2/segments"
	RouteSchemas                                      = "/schemas"
	RouteSchemasTest                                  = "/schemas/test"
	RouteSchemasFieldSpec                             = "/schemas/fieldSpec"
	RouteTables                                       = "/tables"
	RouteTablesTest                                   = "/tables/test"
	RouteTablesTestExternalView                       = "/tables/test/externalview"
	RouteTablesTestIdealState                         = "/tables/test/idealstate"
	RouteTablesTestIndexes                            = "/tables/test/indexes"
	RouteTablesTestInstances                          = "/tables/test/instances"
	RouteTablesLiveBrokers                            = "/tables/livebrokers"
	RouteTablesTestLiveBrokers                        = "/tables/test/livebrokers"
	RouteTablesTestMetadata                           = "/tables/test/metadata"
	RouteTablesTestRebuildBrokerResourceFromHelixTags = "/tables/test/rebuildBrokerResourceFromHelixTags"
	RouteTablesTestSchema                             = "/tables/test/schema"
	RouteTablesTestSize                               = "/tables/test/size"
	RouteTablesTestState                              = "/tables/test/state"
	RouteTablesTestStats                              = "/tables/test/stats"
	RoutePinotControllerAdmin                         = "/pinot-controller/admin"
	RouteHealth                                       = "/health"
)

func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Basic YWRtaW46YWRtaW4K" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next(w, r)
	}

}

func handleGetUsers(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, `{"users": {"test_BROKER": {"username": "test","password": "$2a$10$3KYPvIy4fBM3CWfdKSc54u/BOn1rPlgb7u4P66s4upqAct30C/q6a","component": "BROKER","role": "ADMIN","usernameWithComponent": "test_BROKER"}}}`)
}

func handleGetUser(w http.ResponseWriter, r *http.Request) {
	component := r.URL.Query().Get("component")

	if component == "BROKER" {
		fmt.Fprint(w, `{"test_BROKER": {"username": "test","password": "$2a$10$3KYPvIy4fBM3CWfdKSc54u/BOn1rPlgb7u4P66s4upqAct30C/q6a","component": "BROKER","role": "ADMIN","usernameWithComponent": "test_BROKER"}}`)
	}
}

func handlePutUser(w http.ResponseWriter, r *http.Request) {
	component := r.URL.Query().Get("component")

	if component == "BROKER" {
		fmt.Fprint(w, `{"status": "User config update for test_BROKER"}`)
	}

}

func handleDeleteUser(w http.ResponseWriter, r *http.Request) {

	component := r.URL.Query().Get("component")

	if component == "BROKER" {
		fmt.Fprint(w, `{"status": "User: test_BROKER has been successfully deleted"}`)
	} else if component == "" {
		http.Error(w, `{"code": 400,"error": "Name is null"}`, http.StatusBadRequest)
	}
}

func handlePostUsers(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, `{"status": "User testUser_BROKER has been successfully added!"}`)
}

func handleGetInstances(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, `{"instances": ["Minion_172.19.0.2_9514","Server_172.19.0.7_8098","Broker_cdba1ba98e74_8099","Controller_8684b6757488_9000"]}`)
}

func handlePostInstances(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, `{"status": "Added instance: Broker_localhost_1234"}`)
}

func handleGetInstance(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `{"instanceName": "Minion_172.19.0.2_9514","hostName": "172.19.0.2","enabled": true,"port": "9514","tags": ["minion_untagged"],"pools": null,"grpcPort": -1,"adminPort": -1,"queryServicePort": -1,"queryMailboxPort": -1,"systemResourceInfo": null}`)

}

func handleGetClusterInfo(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, `{"clusterName": "PinotCluster"}`)
}

func handleGetClusterConfigs(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `{"allowParticipantAutoJoin": "true","enable.case.insensitive": "true","pinot.broker.enable.query.limit.override": "false","default.hyperloglog.log2m": "8"}`)
}

func handlePostClusterConfigs(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, `{"status": "Updated cluster config."}`)
}

func handleDeleteClusterConfigs(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, `{"status": "Deleted cluster config: allowParticipantAutoJoin"}`)
}

func handleGetTenants(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, `{"SERVER_TENANTS": ["DefaultTenant"],"BROKER_TENANTS": ["DefaultTenant"]}`)
}

func handleCreateTenant(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, `{"status": "Successfully created tenant"}`)
}

func handleUpdateTenant(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, `{"status": "Updated tenant"}`)
}

func handleDeleteTenant(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, `{"status": "Successfully deleted tenant DefaultTenant"}`)
}

func handleGetTenantInstances(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, `{"ServerInstances": ["Server_172.19.0.7_8098"],"BrokerInstances": ["Broker_91e4732e326d_8099"],"tenantName": "DefaultTenant"}`)
}

func handleGetTenantTables(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, `{"tables": []}`)
}

func handleGetTenantMetadata(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, `{"ServerInstances": ["Server_172.19.0.7_8098"],"OfflineServerInstances": null,"RealtimeServerInstances": null,"BrokerInstances": ["Broker_91e4732e326d_8099"],"tenantName": "DefaultTenant"}`)
}

func handleGetSegments(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, `[{"OFFLINE": ["test_1"],"REALTIME": ["test_1"]}]`)
}

func handleReloadTableSegments(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, `{"status": "{\"test_OFFLINE\":{\"reloadJobId\":\"f9db13c7-3ad6-45a8-a08f-75cd03c42fb5\",\"reloadJobMetaZKStorageStatus\":\"SUCCESS\",\"numMessagesSent\":\"1\"}}"}`)
}

func handleReloadTableSegment(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, `{"status": "Submitted reload job id: ce4650a9-774a-4b22-919b-4cb22b5c8129, sent 1 reload messages. Job meta ZK storage status: SUCCESS"}`)
}

func handleGetSchemas(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, `["test"]`)
}

func handleGetSchema(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, `{"schemaName": "test","enableColumnBasedNullHandling": false,"dimensionFieldSpecs": [{"name": "id","dataType": "STRING","notNull": false},{"name": "type","dataType": "STRING","notNull": false},{"name": "actor","dataType": "JSON","notNull": false},{"name": "repo","dataType": "JSON","notNull": false},{"name": "payload","dataType": "JSON","notNull": false},{"name": "public","dataType": "BOOLEAN","notNull": false}],"dateTimeFieldSpecs": [{"name": "created_at","dataType": "STRING","notNull": false,"format": "1:SECONDS:SIMPLE_DATE_FORMAT:yyyy-MM-dd'T'HH:mm:ss'Z'","granularity": "1:SECONDS"},{"name": "created_at_timestamp","dataType": "TIMESTAMP","notNull": false,"format": "1:MILLISECONDS:TIMESTAMP","granularity": "1:SECONDS"}]}`)
}

func handleCreateSchema(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, `{"unrecognizedProperties": {},"status": "test successfully added"}`)
}

func handleUpdateSchema(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, `{"unrecognizedProperties": {},"status": "test successfully added"}`)
}

func handleDeleteSchema(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, `{"status": "Schema test deleted"}`)
}

func handleValidateSchema(w http.ResponseWriter, r *http.Request) {
	// return 200
	fmt.Fprint(w, `{"ok": "true"}`)
}

func handleGetTables(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, `{"tables": ["test"]}`)
}

func handleGetTableConfig(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `{
		"OFFLINE": {
		  "tableName": "test_OFFLINE",
		  "tableType": "OFFLINE",
		  "segmentsConfig": {
			"timeColumnName": "DaysSinceEpoch",
			"replication": "1",
			"timeType": "DAYS",
			"minimizeDataMovement": false,
			"segmentAssignmentStrategy": "BalanceNumSegmentAssignmentStrategy",
			"segmentPushType": "APPEND"
		  },
		  "tenants": {
			"broker": "DefaultTenant",
			"server": "DefaultTenant"
		  },
		  "tableIndexConfig": {
			"enableDefaultStarTree": false,
			"starTreeIndexConfigs": [
			  {
				"dimensionsSplitOrder": [
				  "AirlineID",
				  "Origin",
				  "Dest"
				],
				"functionColumnPairs": [
				  "COUNT__*",
				  "MAX__ArrDelay"
				],
				"maxLeafRecords": 10
			  }
			],
			"tierOverwrites": {
			  "hotTier": {
				"starTreeIndexConfigs": [
				  {
					"dimensionsSplitOrder": [
					  "Carrier",
					  "CancellationCode",
					  "Origin",
					  "Dest"
					],
					"skipStarNodeCreationForDimensions": [],
					"functionColumnPairs": [
					  "MAX__CarrierDelay",
					  "AVG__CarrierDelay"
					],
					"maxLeafRecords": 10
				  }
				]
			  },
			  "coldTier": {
				"starTreeIndexConfigs": []
			  }
			},
			"enableDynamicStarTreeCreation": true,
			"aggregateMetrics": false,
			"nullHandlingEnabled": false,
			"columnMajorSegmentBuilderEnabled": false,
			"optimizeDictionary": false,
			"optimizeDictionaryForMetrics": false,
			"noDictionarySizeRatioThreshold": 0.85,
			"rangeIndexVersion": 2,
			"autoGeneratedInvertedIndex": false,
			"createInvertedIndexDuringSegmentGeneration": false,
			"loadMode": "MMAP"
		  },
		  "metadata": {
			"customConfigs": {}
		  },
		  "fieldConfigList": [
			{
			  "name": "ts",
			  "encodingType": "DICTIONARY",
			  "indexType": "TIMESTAMP",
			  "indexTypes": [
				"TIMESTAMP"
			  ],
			  "timestampConfig": {
				"granularities": [
				  "DAY",
				  "WEEK",
				  "MONTH"
				]
			  },
			  "indexes": null,
			  "tierOverwrites": null
			},
			{
			  "name": "ArrTimeBlk",
			  "encodingType": "DICTIONARY",
			  "indexTypes": [],
			  "indexes": {
				"inverted": {
				  "enabled": "true"
				}
			  },
			  "tierOverwrites": {
				"hotTier": {
				  "encodingType": "DICTIONARY",
				  "indexes": {
					"bloom": {
					  "enabled": "true"
					}
				  }
				},
				"coldTier": {
				  "encodingType": "RAW",
				  "indexes": {
					"text": {
					  "enabled": "true"
					}
				  }
				}
			  }
			}
		  ],
		  "ingestionConfig": {
			"segmentTimeValueCheck": true,
			"transformConfigs": [
			  {
				"columnName": "ts",
				"transformFunction": "fromEpochDays(DaysSinceEpoch)"
			  },
			  {
				"columnName": "tsRaw",
				"transformFunction": "fromEpochDays(DaysSinceEpoch)"
			  }
			],
			"continueOnError": false,
			"rowTimeValueCheck": false
		  },
		  "tierConfigs": [
			{
			  "name": "hotTier",
			  "segmentSelectorType": "time",
			  "segmentAge": "3130d",
			  "storageType": "pinot_server",
			  "serverTag": "DefaultTenant_OFFLINE"
			},
			{
			  "name": "coldTier",
			  "segmentSelectorType": "time",
			  "segmentAge": "3140d",
			  "storageType": "pinot_server",
			  "serverTag": "DefaultTenant_OFFLINE"
			}
		  ],
		  "isDimTable": false
		},
		"REALTIME":{
			"tableName":"realtime_ethereum_mainnet_block_headers_REALTIME",
			"tableType":"REALTIME",
			"segmentsConfig":{
				"replication":"1",
				"retentionTimeUnit":"DAYS",
				"retentionTimeValue":"7",
				"timeType":"MILLISECONDS",
				"replicasPerPartition":"1",
				"timeColumnName":"timestamp",
				"deletedSegmentsRetentionPeriod":"1d",
				"minimizeDataMovement":false
			},
			"tenants":{
				"broker":"DefaultTenant",
				"server":"DefaultTenant"
			},
			"tableIndexConfig":{
				"segmentNameGeneratorType":"",
				"columnMajorSegmentBuilderEnabled":false,
				"optimizeDictionary":false,
				"optimizeDictionaryForMetrics":false,
				"noDictionarySizeRatioThreshold":0.85,
				"noDictionaryColumns":["hash"],
				"rangeIndexVersion":2,
				"autoGeneratedInvertedIndex":false,
				"createInvertedIndexDuringSegmentGeneration":false,
				"sortedColumn":["number"],
				"loadMode":"MMAP",
				"varLengthDictionaryColumns":["parent_hash"],
				"enableDefaultStarTree":false,
				"enableDynamicStarTreeCreation":false,
				"aggregateMetrics":false,
				"nullHandlingEnabled":false
			},
			"metadata":{
				"customConfigs":{
					"customKey":"customValue"
				}
			},
			"routing":{
				"segmentPrunerTypes":["partition"],
				"instanceSelectorType":"strictReplicaGroup"
			},
			"upsertConfig":{
				"mode":"FULL",
				"dropOutOfOrderRecord":false,
				"hashFunction":"NONE",
				"defaultPartialUpsertStrategy":"OVERWRITE",
				"metadataTTL":0.0,
				"deletedKeysTTL":0.0,
				"enablePreload":false,
				"enableSnapshot":true
			},
			"ingestionConfig":{
				"streamIngestionConfig":{
					"streamConfigMaps":[{
						"stream.kafka.broker.list":"kafka:9092",
						"stream.kafka.consumer.factory.class.name":"org.apache.pinot.plugin.stream.kafka20.KafkaConsumerFactory",
						"stream.kafka.consumer.prop.auto.offset.reset":"smallest",
						"stream.kafka.consumer.type":"high-level",
						"stream.kafka.decoder.class.name":"org.apache.pinot.plugin.stream.kafka.KafkaJSONMessageDecoder",
						"stream.kafka.decoder.prop.schema.registry.rest.url":"http://schema-registry:8081",
						"stream.kafka.decoder.prop.schema.registry.schema.name":"ethereum_mainnet_block_headers-value",
						"stream.kafka.decoder.prop.schema.registry.schema.version":"latest",
						"stream.kafka.topic.name":"ethereum_mainnet_block_headers",
						"stream.kafka.zk.broker.url":"kafka:2181",
						"streamType":"kafka"
					}],
					"columnMajorSegmentBuilderEnabled":false
				},
				"transformConfigs":[{
					"columnName":"hash_json",
					"transformFunction":"json_format(hash)"
				}],
				"continueOnError":true,
				"rowTimeValueCheck":true,
				"segmentTimeValueCheck":true
			},
			"isDimTable":false}}
	}`)
}

func handleCreateTable(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, `{"unrecognizedProperties": {"/fieldConfigList/0/indexTypes": null},"status": "Table test_OFFLINE successfully added"}`)
}

func handleUpdateTable(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, `{"unrecognizedProperties": {"/fieldConfigList/0/indexTypes": null},"status": "Table config updated for test_OFFLINE"}`)
}

func handleDeleteTable(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, `{"status": "Tables: [test_OFFLINE] deleted"}`)
}

func handlePinotControllerAdmin(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "GOOD")
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "OK")
}

func handleGetFieldSpecs(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `
	{
		"fieldTypes": {
		  "METRIC": {
			"allowedDataTypes": {
			  "LONG": {
				"nullDefault": 0
			  },
			  "BIG_DECIMAL": {
				"nullDefault": 0
			  },
			  "DOUBLE": {
				"nullDefault": 0
			  },
			  "FLOAT": {
				"nullDefault": 0
			  },
			  "BYTES": {
				"nullDefault": ""
			  },
			  "INT": {
				"nullDefault": 0
			  }
			}
		  },
		  "DATE_TIME": {
			"allowedDataTypes": {
			  "LONG": {
				"nullDefault": -9223372036854776000
			  },
			  "JSON": {
				"nullDefault": "null"
			  },
			  "BIG_DECIMAL": {
				"nullDefault": 0
			  },
			  "BOOLEAN": {
				"nullDefault": 0
			  },
			  "STRING": {
				"nullDefault": "null"
			  },
			  "DOUBLE": {
				"nullDefault": "-Infinity"
			  },
			  "TIMESTAMP": {
				"nullDefault": 0
			  },
			  "FLOAT": {
				"nullDefault": "-Infinity"
			  },
			  "BYTES": {
				"nullDefault": ""
			  },
			  "INT": {
				"nullDefault": -2147483648
			  }
			}
		  },
		  "DIMENSION": {
			"allowedDataTypes": {
			  "LONG": {
				"nullDefault": -9223372036854776000
			  },
			  "JSON": {
				"nullDefault": "null"
			  },
			  "BIG_DECIMAL": {
				"nullDefault": 0
			  },
			  "BOOLEAN": {
				"nullDefault": 0
			  },
			  "STRING": {
				"nullDefault": "null"
			  },
			  "DOUBLE": {
				"nullDefault": "-Infinity"
			  },
			  "TIMESTAMP": {
				"nullDefault": 0
			  },
			  "FLOAT": {
				"nullDefault": "-Infinity"
			  },
			  "BYTES": {
				"nullDefault": ""
			  },
			  "INT": {
				"nullDefault": -2147483648
			  }
			}
		  },
		  "COMPLEX": {
			"allowedDataTypes": {
			  "STRUCT": {
				"nullDefault": null
			  },
			  "MAP": {
				"nullDefault": null
			  },
			  "LIST": {
				"nullDefault": null
			  }
			}
		  },
		  "TIME": {
			"allowedDataTypes": {
			  "LONG": {
				"nullDefault": -9223372036854776000
			  },
			  "JSON": {
				"nullDefault": "null"
			  },
			  "BIG_DECIMAL": {
				"nullDefault": 0
			  },
			  "BOOLEAN": {
				"nullDefault": 0
			  },
			  "STRING": {
				"nullDefault": "null"
			  },
			  "DOUBLE": {
				"nullDefault": "-Infinity"
			  },
			  "TIMESTAMP": {
				"nullDefault": 0
			  },
			  "FLOAT": {
				"nullDefault": "-Infinity"
			  },
			  "BYTES": {
				"nullDefault": ""
			  },
			  "INT": {
				"nullDefault": -2147483648
			  }
			}
		  }
		},
		"dataTypes": {
		  "LONG": {
			"storedType": "LONG",
			"size": 8,
			"sortable": true,
			"numeric": true
		  },
		  "JSON": {
			"storedType": "STRING",
			"size": -1,
			"sortable": false,
			"numeric": false
		  },
		  "STRUCT": {
			"storedType": "STRUCT",
			"size": -1,
			"sortable": false,
			"numeric": false
		  },
		  "BOOLEAN": {
			"storedType": "INT",
			"size": 4,
			"sortable": true,
			"numeric": false
		  },
		  "STRING": {
			"storedType": "STRING",
			"size": -1,
			"sortable": true,
			"numeric": false
		  },
		  "TIMESTAMP": {
			"storedType": "LONG",
			"size": 8,
			"sortable": true,
			"numeric": false
		  },
		  "FLOAT": {
			"storedType": "FLOAT",
			"size": 4,
			"sortable": true,
			"numeric": true
		  },
		  "UNKNOWN": {
			"storedType": "UNKNOWN",
			"size": -1,
			"sortable": true,
			"numeric": false
		  },
		  "MAP": {
			"storedType": "MAP",
			"size": -1,
			"sortable": false,
			"numeric": false
		  },
		  "LIST": {
			"storedType": "LIST",
			"size": -1,
			"sortable": false,
			"numeric": false
		  },
		  "BIG_DECIMAL": {
			"storedType": "BIG_DECIMAL",
			"size": -1,
			"sortable": true,
			"numeric": true
		  },
		  "DOUBLE": {
			"storedType": "DOUBLE",
			"size": 8,
			"sortable": true,
			"numeric": true
		  },
		  "BYTES": {
			"storedType": "BYTES",
			"size": -1,
			"sortable": false,
			"numeric": false
		  },
		  "INT": {
			"storedType": "INT",
			"size": 4,
			"sortable": true,
			"numeric": true
		  }
		}
	  }`)
}

func handleTableExternalView(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `{
		"OFFLINE": {
		  "test_OFFLINE_0": {
			"Server_172.17.0.3_7050": "ONLINE"
		  }
		},
		"REALTIME": null
	  }`)
}

func handleTableIdealState(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `{
		"OFFLINE": {
		  "test_OFFLINE_0": {
			"Server_172.17.0.3_7050": "ONLINE"
		  }
		},
		"REALTIME": null
	  }`)
}

func handleTableIndexes(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, `{
		"totalOnlineSegments": 31,
		"columnToIndexesCount": {
		  "Quarter": {
			"vector_index": 0,
			"nullvalue_vector": 0,
			"h3_index": 0,
			"dictionary": 31,
			"json_index": 0,
			"range_index": 0,
			"forward_index": 31,
			"bloom_filter": 0,
			"inverted_index": 31,
			"text_index": 0,
			"fst_index": 0
		  },
		  "Origin": {
			"vector_index": 0,
			"nullvalue_vector": 0,
			"h3_index": 0,
			"dictionary": 31,
			"json_index": 0,
			"range_index": 0,
			"forward_index": 31,
			"bloom_filter": 0,
			"inverted_index": 0,
			"text_index": 0,
			"fst_index": 0
		  },
		  "FlightNum": {
			"vector_index": 0,
			"nullvalue_vector": 0,
			"h3_index": 0,
			"dictionary": 31,
			"json_index": 0,
			"range_index": 0,
			"forward_index": 31,
			"bloom_filter": 0,
			"inverted_index": 0,
			"text_index": 0,
			"fst_index": 0
		  },
		  "LateAircraftDelay": {
			"vector_index": 0,
			"nullvalue_vector": 0,
			"h3_index": 0,
			"dictionary": 31,
			"json_index": 0,
			"range_index": 0,
			"forward_index": 31,
			"bloom_filter": 0,
			"inverted_index": 0,
			"text_index": 0,
			"fst_index": 0
		  },
		  "DivWheelsOns": {
			"vector_index": 0,
			"nullvalue_vector": 0,
			"h3_index": 0,
			"dictionary": 31,
			"json_index": 0,
			"range_index": 0,
			"forward_index": 31,
			"bloom_filter": 0,
			"inverted_index": 0,
			"text_index": 0,
			"fst_index": 0
		  },
		  "DivActualElapsedTime": {
			"vector_index": 0,
			"nullvalue_vector": 0,
			"h3_index": 0,
			"dictionary": 31,
			"json_index": 0,
			"range_index": 0,
			"forward_index": 31,
			"bloom_filter": 0,
			"inverted_index": 0,
			"text_index": 0,
			"fst_index": 0
		  },
		  "DivWheelsOffs": {
			"vector_index": 0,
			"nullvalue_vector": 0,
			"h3_index": 0,
			"dictionary": 31,
			"json_index": 0,
			"range_index": 0,
			"forward_index": 31,
			"bloom_filter": 0,
			"inverted_index": 0,
			"text_index": 0,
			"fst_index": 0
		  },
		  "ArrDel15": {
			"vector_index": 0,
			"nullvalue_vector": 0,
			"h3_index": 0,
			"dictionary": 31,
			"json_index": 0,
			"range_index": 0,
			"forward_index": 31,
			"bloom_filter": 0,
			"inverted_index": 0,
			"text_index": 0,
			"fst_index": 0
		  },
		  "AirTime": {
			"vector_index": 0,
			"nullvalue_vector": 0,
			"h3_index": 0,
			"dictionary": 31,
			"json_index": 0,
			"range_index": 0,
			"forward_index": 31,
			"bloom_filter": 0,
			"inverted_index": 0,
			"text_index": 0,
			"fst_index": 0
		  },
		  "DivTotalGTimes": {
			"vector_index": 0,
			"nullvalue_vector": 0,
			"h3_index": 0,
			"dictionary": 31,
			"json_index": 0,
			"range_index": 0,
			"forward_index": 31,
			"bloom_filter": 0,
			"inverted_index": 0,
			"text_index": 0,
			"fst_index": 0
		  },
		  "DepTimeBlk": {
			"vector_index": 0,
			"nullvalue_vector": 0,
			"h3_index": 0,
			"dictionary": 31,
			"json_index": 0,
			"range_index": 0,
			"forward_index": 31,
			"bloom_filter": 0,
			"inverted_index": 0,
			"text_index": 0,
			"fst_index": 0
		  },
		  "DestCityMarketID": {
			"vector_index": 0,
			"nullvalue_vector": 0,
			"h3_index": 0,
			"dictionary": 31,
			"json_index": 0,
			"range_index": 0,
			"forward_index": 31,
			"bloom_filter": 0,
			"inverted_index": 0,
			"text_index": 0,
			"fst_index": 0
		  },
		  "$ts$WEEK": {
			"vector_index": 0,
			"nullvalue_vector": 0,
			"h3_index": 0,
			"dictionary": 31,
			"json_index": 0,
			"range_index": 31,
			"forward_index": 31,
			"bloom_filter": 0,
			"inverted_index": 31,
			"text_index": 0,
			"fst_index": 0
		  },
		  "DivAirportSeqIDs": {
			"vector_index": 0,
			"nullvalue_vector": 0,
			"h3_index": 0,
			"dictionary": 31,
			"json_index": 0,
			"range_index": 0,
			"forward_index": 31,
			"bloom_filter": 0,
			"inverted_index": 0,
			"text_index": 0,
			"fst_index": 0
		  },
		  "DaysSinceEpoch": {
			"vector_index": 0,
			"nullvalue_vector": 0,
			"h3_index": 0,
			"dictionary": 31,
			"json_index": 0,
			"range_index": 0,
			"forward_index": 31,
			"bloom_filter": 0,
			"inverted_index": 31,
			"text_index": 0,
			"fst_index": 0
		  },
		  "DepTime": {
			"vector_index": 0,
			"nullvalue_vector": 0,
			"h3_index": 0,
			"dictionary": 31,
			"json_index": 0,
			"range_index": 0,
			"forward_index": 31,
			"bloom_filter": 0,
			"inverted_index": 0,
			"text_index": 0,
			"fst_index": 0
		  },
		  "Month": {
			"vector_index": 0,
			"nullvalue_vector": 0,
			"h3_index": 0,
			"dictionary": 31,
			"json_index": 0,
			"range_index": 0,
			"forward_index": 31,
			"bloom_filter": 0,
			"inverted_index": 31,
			"text_index": 0,
			"fst_index": 0
		  },
		  "$segmentName": {
			"vector_index": 0,
			"nullvalue_vector": 0,
			"h3_index": 0,
			"dictionary": 31,
			"json_index": 0,
			"range_index": 0,
			"forward_index": 31,
			"bloom_filter": 0,
			"inverted_index": 31,
			"text_index": 0,
			"fst_index": 0
		  },
		  "DestStateName": {
			"vector_index": 0,
			"nullvalue_vector": 0,
			"h3_index": 0,
			"dictionary": 31,
			"json_index": 0,
			"range_index": 0,
			"forward_index": 31,
			"bloom_filter": 0,
			"inverted_index": 0,
			"text_index": 0,
			"fst_index": 0
		  },
		  "CRSElapsedTime": {
			"vector_index": 0,
			"nullvalue_vector": 0,
			"h3_index": 0,
			"dictionary": 31,
			"json_index": 0,
			"range_index": 0,
			"forward_index": 31,
			"bloom_filter": 0,
			"inverted_index": 0,
			"text_index": 0,
			"fst_index": 0
		  },
		  "Carrier": {
			"vector_index": 0,
			"nullvalue_vector": 0,
			"h3_index": 0,
			"dictionary": 31,
			"json_index": 0,
			"range_index": 0,
			"forward_index": 31,
			"bloom_filter": 0,
			"inverted_index": 0,
			"text_index": 0,
			"fst_index": 0
		  },
		  "Distance": {
			"vector_index": 0,
			"nullvalue_vector": 0,
			"h3_index": 0,
			"dictionary": 31,
			"json_index": 0,
			"range_index": 0,
			"forward_index": 31,
			"bloom_filter": 0,
			"inverted_index": 0,
			"text_index": 0,
			"fst_index": 0
		  },
		  "DestAirportID": {
			"vector_index": 0,
			"nullvalue_vector": 0,
			"h3_index": 0,
			"dictionary": 31,
			"json_index": 0,
			"range_index": 0,
			"forward_index": 31,
			"bloom_filter": 0,
			"inverted_index": 0,
			"text_index": 0,
			"fst_index": 0
		  },
		  "ArrTimeBlk": {
			"vector_index": 0,
			"nullvalue_vector": 0,
			"h3_index": 0,
			"dictionary": 31,
			"json_index": 0,
			"range_index": 0,
			"forward_index": 31,
			"bloom_filter": 0,
			"inverted_index": 31,
			"text_index": 0,
			"fst_index": 0
		  },
		  "SecurityDelay": {
			"vector_index": 0,
			"nullvalue_vector": 0,
			"h3_index": 0,
			"dictionary": 31,
			"json_index": 0,
			"range_index": 0,
			"forward_index": 31,
			"bloom_filter": 0,
			"inverted_index": 0,
			"text_index": 0,
			"fst_index": 0
		  },
		  "DivArrDelay": {
			"vector_index": 0,
			"nullvalue_vector": 0,
			"h3_index": 0,
			"dictionary": 31,
			"json_index": 0,
			"range_index": 0,
			"forward_index": 31,
			"bloom_filter": 0,
			"inverted_index": 0,
			"text_index": 0,
			"fst_index": 0
		  },
		  "LongestAddGTime": {
			"vector_index": 0,
			"nullvalue_vector": 0,
			"h3_index": 0,
			"dictionary": 31,
			"json_index": 0,
			"range_index": 0,
			"forward_index": 31,
			"bloom_filter": 0,
			"inverted_index": 5,
			"text_index": 0,
			"fst_index": 0
		  },
		  "OriginWac": {
			"vector_index": 0,
			"nullvalue_vector": 0,
			"h3_index": 0,
			"dictionary": 31,
			"json_index": 0,
			"range_index": 0,
			"forward_index": 31,
			"bloom_filter": 0,
			"inverted_index": 0,
			"text_index": 0,
			"fst_index": 0
		  },
		  "WheelsOff": {
			"vector_index": 0,
			"nullvalue_vector": 0,
			"h3_index": 0,
			"dictionary": 31,
			"json_index": 0,
			"range_index": 0,
			"forward_index": 31,
			"bloom_filter": 0,
			"inverted_index": 0,
			"text_index": 0,
			"fst_index": 0
		  },
		  "UniqueCarrier": {
			"vector_index": 0,
			"nullvalue_vector": 0,
			"h3_index": 0,
			"dictionary": 31,
			"json_index": 0,
			"range_index": 0,
			"forward_index": 31,
			"bloom_filter": 0,
			"inverted_index": 0,
			"text_index": 0,
			"fst_index": 0
		  },
		  "DestAirportSeqID": {
			"vector_index": 0,
			"nullvalue_vector": 0,
			"h3_index": 0,
			"dictionary": 31,
			"json_index": 0,
			"range_index": 0,
			"forward_index": 31,
			"bloom_filter": 0,
			"inverted_index": 0,
			"text_index": 0,
			"fst_index": 0
		  },
		  "DivReachedDest": {
			"vector_index": 0,
			"nullvalue_vector": 0,
			"h3_index": 0,
			"dictionary": 31,
			"json_index": 0,
			"range_index": 0,
			"forward_index": 31,
			"bloom_filter": 0,
			"inverted_index": 0,
			"text_index": 0,
			"fst_index": 0
		  },
		  "Diverted": {
			"vector_index": 0,
			"nullvalue_vector": 0,
			"h3_index": 0,
			"dictionary": 31,
			"json_index": 0,
			"range_index": 0,
			"forward_index": 31,
			"bloom_filter": 0,
			"inverted_index": 0,
			"text_index": 0,
			"fst_index": 0
		  },
		  "ActualElapsedTime": {
			"vector_index": 0,
			"nullvalue_vector": 0,
			"h3_index": 0,
			"dictionary": 31,
			"json_index": 0,
			"range_index": 0,
			"forward_index": 31,
			"bloom_filter": 0,
			"inverted_index": 0,
			"text_index": 0,
			"fst_index": 0
		  },
		  "tsRaw": {
			"vector_index": 0,
			"nullvalue_vector": 0,
			"h3_index": 0,
			"dictionary": 31,
			"json_index": 0,
			"range_index": 0,
			"forward_index": 31,
			"bloom_filter": 0,
			"inverted_index": 31,
			"text_index": 0,
			"fst_index": 0
		  },
		  "OriginStateName": {
			"vector_index": 0,
			"nullvalue_vector": 0,
			"h3_index": 0,
			"dictionary": 31,
			"json_index": 0,
			"range_index": 0,
			"forward_index": 31,
			"bloom_filter": 0,
			"inverted_index": 0,
			"text_index": 0,
			"fst_index": 0
		  },
		  "AirlineID": {
			"vector_index": 0,
			"nullvalue_vector": 0,
			"h3_index": 0,
			"dictionary": 31,
			"json_index": 0,
			"range_index": 0,
			"forward_index": 31,
			"bloom_filter": 0,
			"inverted_index": 0,
			"text_index": 0,
			"fst_index": 0
		  },
		  "FlightDate": {
			"vector_index": 0,
			"nullvalue_vector": 0,
			"h3_index": 0,
			"dictionary": 31,
			"json_index": 0,
			"range_index": 0,
			"forward_index": 31,
			"bloom_filter": 0,
			"inverted_index": 31,
			"text_index": 0,
			"fst_index": 0
		  },
		  "DivAirportLandings": {
			"vector_index": 0,
			"nullvalue_vector": 0,
			"h3_index": 0,
			"dictionary": 31,
			"json_index": 0,
			"range_index": 0,
			"forward_index": 31,
			"bloom_filter": 0,
			"inverted_index": 0,
			"text_index": 0,
			"fst_index": 0
		  },
		  "DepartureDelayGroups": {
			"vector_index": 0,
			"nullvalue_vector": 0,
			"h3_index": 0,
			"dictionary": 31,
			"json_index": 0,
			"range_index": 0,
			"forward_index": 31,
			"bloom_filter": 0,
			"inverted_index": 0,
			"text_index": 0,
			"fst_index": 0
		  },
		  "OriginCityName": {
			"vector_index": 0,
			"nullvalue_vector": 0,
			"h3_index": 0,
			"dictionary": 31,
			"json_index": 0,
			"range_index": 0,
			"forward_index": 31,
			"bloom_filter": 0,
			"inverted_index": 0,
			"text_index": 0,
			"fst_index": 0
		  },
		  "OriginStateFips": {
			"vector_index": 0,
			"nullvalue_vector": 0,
			"h3_index": 0,
			"dictionary": 31,
			"json_index": 0,
			"range_index": 0,
			"forward_index": 31,
			"bloom_filter": 0,
			"inverted_index": 0,
			"text_index": 0,
			"fst_index": 0
		  },
		  "OriginState": {
			"vector_index": 0,
			"nullvalue_vector": 0,
			"h3_index": 0,
			"dictionary": 31,
			"json_index": 0,
			"range_index": 0,
			"forward_index": 31,
			"bloom_filter": 0,
			"inverted_index": 0,
			"text_index": 0,
			"fst_index": 0
		  },
		  "DistanceGroup": {
			"vector_index": 0,
			"nullvalue_vector": 0,
			"h3_index": 0,
			"dictionary": 31,
			"json_index": 0,
			"range_index": 0,
			"forward_index": 31,
			"bloom_filter": 0,
			"inverted_index": 0,
			"text_index": 0,
			"fst_index": 0
		  },
		  "WeatherDelay": {
			"vector_index": 0,
			"nullvalue_vector": 0,
			"h3_index": 0,
			"dictionary": 31,
			"json_index": 0,
			"range_index": 0,
			"forward_index": 31,
			"bloom_filter": 0,
			"inverted_index": 0,
			"text_index": 0,
			"fst_index": 0
		  },
		  "DestWac": {
			"vector_index": 0,
			"nullvalue_vector": 0,
			"h3_index": 0,
			"dictionary": 31,
			"json_index": 0,
			"range_index": 0,
			"forward_index": 31,
			"bloom_filter": 0,
			"inverted_index": 0,
			"text_index": 0,
			"fst_index": 0
		  },
		  "WheelsOn": {
			"vector_index": 0,
			"nullvalue_vector": 0,
			"h3_index": 0,
			"dictionary": 31,
			"json_index": 0,
			"range_index": 0,
			"forward_index": 31,
			"bloom_filter": 0,
			"inverted_index": 0,
			"text_index": 0,
			"fst_index": 0
		  },
		  "OriginAirportID": {
			"vector_index": 0,
			"nullvalue_vector": 0,
			"h3_index": 0,
			"dictionary": 31,
			"json_index": 0,
			"range_index": 0,
			"forward_index": 31,
			"bloom_filter": 0,
			"inverted_index": 0,
			"text_index": 0,
			"fst_index": 0
		  },
		  "OriginCityMarketID": {
			"vector_index": 0,
			"nullvalue_vector": 0,
			"h3_index": 0,
			"dictionary": 31,
			"json_index": 0,
			"range_index": 0,
			"forward_index": 31,
			"bloom_filter": 0,
			"inverted_index": 0,
			"text_index": 0,
			"fst_index": 0
		  },
		  "NASDelay": {
			"vector_index": 0,
			"nullvalue_vector": 0,
			"h3_index": 0,
			"dictionary": 31,
			"json_index": 0,
			"range_index": 0,
			"forward_index": 31,
			"bloom_filter": 0,
			"inverted_index": 0,
			"text_index": 0,
			"fst_index": 0
		  },
		  "$ts$DAY": {
			"vector_index": 0,
			"nullvalue_vector": 0,
			"h3_index": 0,
			"dictionary": 31,
			"json_index": 0,
			"range_index": 31,
			"forward_index": 31,
			"bloom_filter": 0,
			"inverted_index": 31,
			"text_index": 0,
			"fst_index": 0
		  },
		  "DestState": {
			"vector_index": 0,
			"nullvalue_vector": 0,
			"h3_index": 0,
			"dictionary": 31,
			"json_index": 0,
			"range_index": 0,
			"forward_index": 31,
			"bloom_filter": 0,
			"inverted_index": 0,
			"text_index": 0,
			"fst_index": 0
		  },
		  "ArrTime": {
			"vector_index": 0,
			"nullvalue_vector": 0,
			"h3_index": 0,
			"dictionary": 31,
			"json_index": 0,
			"range_index": 0,
			"forward_index": 31,
			"bloom_filter": 0,
			"inverted_index": 0,
			"text_index": 0,
			"fst_index": 0
		  },
		  "$hostName": {
			"vector_index": 0,
			"nullvalue_vector": 0,
			"h3_index": 0,
			"dictionary": 31,
			"json_index": 0,
			"range_index": 0,
			"forward_index": 31,
			"bloom_filter": 0,
			"inverted_index": 31,
			"text_index": 0,
			"fst_index": 0
		  },
		  "ArrivalDelayGroups": {
			"vector_index": 0,
			"nullvalue_vector": 0,
			"h3_index": 0,
			"dictionary": 31,
			"json_index": 0,
			"range_index": 0,
			"forward_index": 31,
			"bloom_filter": 0,
			"inverted_index": 0,
			"text_index": 0,
			"fst_index": 0
		  },
		  "Flights": {
			"vector_index": 0,
			"nullvalue_vector": 0,
			"h3_index": 0,
			"dictionary": 31,
			"json_index": 0,
			"range_index": 0,
			"forward_index": 31,
			"bloom_filter": 0,
			"inverted_index": 31,
			"text_index": 0,
			"fst_index": 0
		  },
		  "RandomAirports": {
			"vector_index": 0,
			"nullvalue_vector": 0,
			"h3_index": 0,
			"dictionary": 31,
			"json_index": 0,
			"range_index": 0,
			"forward_index": 31,
			"bloom_filter": 0,
			"inverted_index": 0,
			"text_index": 0,
			"fst_index": 0
		  },
		  "DayofMonth": {
			"vector_index": 0,
			"nullvalue_vector": 0,
			"h3_index": 0,
			"dictionary": 31,
			"json_index": 0,
			"range_index": 0,
			"forward_index": 31,
			"bloom_filter": 0,
			"inverted_index": 31,
			"text_index": 0,
			"fst_index": 0
		  },
		  "TotalAddGTime": {
			"vector_index": 0,
			"nullvalue_vector": 0,
			"h3_index": 0,
			"dictionary": 31,
			"json_index": 0,
			"range_index": 0,
			"forward_index": 31,
			"bloom_filter": 0,
			"inverted_index": 5,
			"text_index": 0,
			"fst_index": 0
		  },
		  "CRSDepTime": {
			"vector_index": 0,
			"nullvalue_vector": 0,
			"h3_index": 0,
			"dictionary": 31,
			"json_index": 0,
			"range_index": 0,
			"forward_index": 31,
			"bloom_filter": 0,
			"inverted_index": 0,
			"text_index": 0,
			"fst_index": 0
		  },
		  "DayOfWeek": {
			"vector_index": 0,
			"nullvalue_vector": 0,
			"h3_index": 0,
			"dictionary": 31,
			"json_index": 0,
			"range_index": 0,
			"forward_index": 31,
			"bloom_filter": 0,
			"inverted_index": 31,
			"text_index": 0,
			"fst_index": 0
		  },
		  "FirstDepTime": {
			"vector_index": 0,
			"nullvalue_vector": 0,
			"h3_index": 0,
			"dictionary": 31,
			"json_index": 0,
			"range_index": 0,
			"forward_index": 31,
			"bloom_filter": 0,
			"inverted_index": 5,
			"text_index": 0,
			"fst_index": 0
		  },
		  "Dest": {
			"vector_index": 0,
			"nullvalue_vector": 0,
			"h3_index": 0,
			"dictionary": 31,
			"json_index": 0,
			"range_index": 0,
			"forward_index": 31,
			"bloom_filter": 0,
			"inverted_index": 0,
			"text_index": 0,
			"fst_index": 0
		  },
		  "CancellationCode": {
			"vector_index": 0,
			"nullvalue_vector": 0,
			"h3_index": 0,
			"dictionary": 31,
			"json_index": 0,
			"range_index": 0,
			"forward_index": 31,
			"bloom_filter": 0,
			"inverted_index": 1,
			"text_index": 0,
			"fst_index": 0
		  },
		  "DivTailNums": {
			"vector_index": 0,
			"nullvalue_vector": 0,
			"h3_index": 0,
			"dictionary": 31,
			"json_index": 0,
			"range_index": 0,
			"forward_index": 31,
			"bloom_filter": 0,
			"inverted_index": 0,
			"text_index": 0,
			"fst_index": 0
		  },
		  "DepDelayMinutes": {
			"vector_index": 0,
			"nullvalue_vector": 0,
			"h3_index": 0,
			"dictionary": 31,
			"json_index": 0,
			"range_index": 0,
			"forward_index": 31,
			"bloom_filter": 0,
			"inverted_index": 0,
			"text_index": 0,
			"fst_index": 0
		  },
		  "TaxiIn": {
			"vector_index": 0,
			"nullvalue_vector": 0,
			"h3_index": 0,
			"dictionary": 31,
			"json_index": 0,
			"range_index": 0,
			"forward_index": 31,
			"bloom_filter": 0,
			"inverted_index": 0,
			"text_index": 0,
			"fst_index": 0
		  },
		  "DepDelay": {
			"vector_index": 0,
			"nullvalue_vector": 0,
			"h3_index": 0,
			"dictionary": 31,
			"json_index": 0,
			"range_index": 0,
			"forward_index": 31,
			"bloom_filter": 0,
			"inverted_index": 0,
			"text_index": 0,
			"fst_index": 0
		  },
		  "OriginAirportSeqID": {
			"vector_index": 0,
			"nullvalue_vector": 0,
			"h3_index": 0,
			"dictionary": 31,
			"json_index": 0,
			"range_index": 0,
			"forward_index": 31,
			"bloom_filter": 0,
			"inverted_index": 0,
			"text_index": 0,
			"fst_index": 0
		  },
		  "DestStateFips": {
			"vector_index": 0,
			"nullvalue_vector": 0,
			"h3_index": 0,
			"dictionary": 31,
			"json_index": 0,
			"range_index": 0,
			"forward_index": 31,
			"bloom_filter": 0,
			"inverted_index": 0,
			"text_index": 0,
			"fst_index": 0
		  },
		  "Cancelled": {
			"vector_index": 0,
			"nullvalue_vector": 0,
			"h3_index": 0,
			"dictionary": 31,
			"json_index": 0,
			"range_index": 0,
			"forward_index": 31,
			"bloom_filter": 0,
			"inverted_index": 1,
			"text_index": 0,
			"fst_index": 0
		  },
		  "ArrDelay": {
			"vector_index": 0,
			"nullvalue_vector": 0,
			"h3_index": 0,
			"dictionary": 31,
			"json_index": 0,
			"range_index": 0,
			"forward_index": 31,
			"bloom_filter": 0,
			"inverted_index": 0,
			"text_index": 0,
			"fst_index": 0
		  },
		  "DivAirportIDs": {
			"vector_index": 0,
			"nullvalue_vector": 0,
			"h3_index": 0,
			"dictionary": 31,
			"json_index": 0,
			"range_index": 0,
			"forward_index": 31,
			"bloom_filter": 0,
			"inverted_index": 0,
			"text_index": 0,
			"fst_index": 0
		  },
		  "$docId": {
			"vector_index": 0,
			"nullvalue_vector": 0,
			"h3_index": 0,
			"dictionary": 31,
			"json_index": 0,
			"range_index": 0,
			"forward_index": 31,
			"bloom_filter": 0,
			"inverted_index": 31,
			"text_index": 0,
			"fst_index": 0
		  },
		  "TaxiOut": {
			"vector_index": 0,
			"nullvalue_vector": 0,
			"h3_index": 0,
			"dictionary": 31,
			"json_index": 0,
			"range_index": 0,
			"forward_index": 31,
			"bloom_filter": 0,
			"inverted_index": 0,
			"text_index": 0,
			"fst_index": 0
		  },
		  "DivLongestGTimes": {
			"vector_index": 0,
			"nullvalue_vector": 0,
			"h3_index": 0,
			"dictionary": 31,
			"json_index": 0,
			"range_index": 0,
			"forward_index": 31,
			"bloom_filter": 0,
			"inverted_index": 0,
			"text_index": 0,
			"fst_index": 0
		  },
		  "DepDel15": {
			"vector_index": 0,
			"nullvalue_vector": 0,
			"h3_index": 0,
			"dictionary": 31,
			"json_index": 0,
			"range_index": 0,
			"forward_index": 31,
			"bloom_filter": 0,
			"inverted_index": 0,
			"text_index": 0,
			"fst_index": 0
		  },
		  "CarrierDelay": {
			"vector_index": 0,
			"nullvalue_vector": 0,
			"h3_index": 0,
			"dictionary": 31,
			"json_index": 0,
			"range_index": 0,
			"forward_index": 31,
			"bloom_filter": 0,
			"inverted_index": 0,
			"text_index": 0,
			"fst_index": 0
		  },
		  "DivAirports": {
			"vector_index": 0,
			"nullvalue_vector": 0,
			"h3_index": 0,
			"dictionary": 31,
			"json_index": 0,
			"range_index": 0,
			"forward_index": 31,
			"bloom_filter": 0,
			"inverted_index": 0,
			"text_index": 0,
			"fst_index": 0
		  },
		  "DivDistance": {
			"vector_index": 0,
			"nullvalue_vector": 0,
			"h3_index": 0,
			"dictionary": 31,
			"json_index": 0,
			"range_index": 0,
			"forward_index": 31,
			"bloom_filter": 0,
			"inverted_index": 0,
			"text_index": 0,
			"fst_index": 0
		  },
		  "Year": {
			"vector_index": 0,
			"nullvalue_vector": 0,
			"h3_index": 0,
			"dictionary": 31,
			"json_index": 0,
			"range_index": 0,
			"forward_index": 31,
			"bloom_filter": 0,
			"inverted_index": 31,
			"text_index": 0,
			"fst_index": 0
		  },
		  "$ts$MONTH": {
			"vector_index": 0,
			"nullvalue_vector": 0,
			"h3_index": 0,
			"dictionary": 31,
			"json_index": 0,
			"range_index": 31,
			"forward_index": 31,
			"bloom_filter": 0,
			"inverted_index": 31,
			"text_index": 0,
			"fst_index": 0
		  },
		  "CRSArrTime": {
			"vector_index": 0,
			"nullvalue_vector": 0,
			"h3_index": 0,
			"dictionary": 31,
			"json_index": 0,
			"range_index": 0,
			"forward_index": 31,
			"bloom_filter": 0,
			"inverted_index": 0,
			"text_index": 0,
			"fst_index": 0
		  },
		  "ArrDelayMinutes": {
			"vector_index": 0,
			"nullvalue_vector": 0,
			"h3_index": 0,
			"dictionary": 31,
			"json_index": 0,
			"range_index": 0,
			"forward_index": 31,
			"bloom_filter": 0,
			"inverted_index": 0,
			"text_index": 0,
			"fst_index": 0
		  },
		  "TailNum": {
			"vector_index": 0,
			"nullvalue_vector": 0,
			"h3_index": 0,
			"dictionary": 31,
			"json_index": 0,
			"range_index": 0,
			"forward_index": 31,
			"bloom_filter": 0,
			"inverted_index": 0,
			"text_index": 0,
			"fst_index": 0
		  },
		  "ts": {
			"vector_index": 0,
			"nullvalue_vector": 0,
			"h3_index": 0,
			"dictionary": 31,
			"json_index": 0,
			"range_index": 0,
			"forward_index": 31,
			"bloom_filter": 0,
			"inverted_index": 31,
			"text_index": 0,
			"fst_index": 0
		  },
		  "DestCityName": {
			"vector_index": 0,
			"nullvalue_vector": 0,
			"h3_index": 0,
			"dictionary": 31,
			"json_index": 0,
			"range_index": 0,
			"forward_index": 31,
			"bloom_filter": 0,
			"inverted_index": 0,
			"text_index": 0,
			"fst_index": 0
		  }
		}
	  }`)
}

func handleGetTableInstances(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `{
		"tableName": "test",
		"brokers": [
		  {
			"tableType": "offline",
			"instances": [
			  "Broker_172.17.0.3_8000"
			]
		  }
		],
		"server": []
	  }`)
}

func handleGetTableLiveBrokers(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `{
		"test_OFFLINE": [
		  {
			"instanceName": "Broker_172.17.0.3_8000",
			"port": 8000,
			"host": "172.17.0.3"
		  }
		],
		"baseballStats_OFFLINE": [
		  {
			"instanceName": "Broker_172.17.0.3_8000",
			"port": 8000,
			"host": "172.17.0.3"
		  }
		],
		"dimBaseballTeams_OFFLINE": [
		  {
			"instanceName": "Broker_172.17.0.3_8000",
			"port": 8000,
			"host": "172.17.0.3"
		  }
		],
		"fineFoodReviews_OFFLINE": [
		  {
			"instanceName": "Broker_172.17.0.3_8000",
			"port": 8000,
			"host": "172.17.0.3"
		  }
		],
		"starbucksStores_OFFLINE": [
		  {
			"instanceName": "Broker_172.17.0.3_8000",
			"port": 8000,
			"host": "172.17.0.3"
		  }
		],
		"githubEvents_OFFLINE": [
		  {
			"instanceName": "Broker_172.17.0.3_8000",
			"port": 8000,
			"host": "172.17.0.3"
		  }
		],
		"billing_OFFLINE": [
		  {
			"instanceName": "Broker_172.17.0.3_8000",
			"port": 8000,
			"host": "172.17.0.3"
		  }
		],
		"githubComplexTypeEvents_OFFLINE": [
		  {
			"instanceName": "Broker_172.17.0.3_8000",
			"port": 8000,
			"host": "172.17.0.3"
		  }
		]
	  }`)
}

func handleGetTableTestLiveBrokers(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `[
		"Broker_172.17.0.3_8000"
	  ]`)
}

func handleGetTableTestMetadata(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `{
		"tableName": "test_OFFLINE",
		"diskSizeInBytes": 72068,
		"numSegments": 1,
		"numRows": 694,
		"columnLengthMap": {},
		"columnCardinalityMap": {},
		"maxNumMultiValuesMap": {},
		"columnIndexSizeMap": {},
		"upsertPartitionToServerPrimaryKeyCountMap": {}
	  }`)
}

func handleTableRebuildBrokerResourceFromHelixTags(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `{
		"status": "Broker resource is not rebuilt because ideal state is the same for table: test_OFFLINE"
	  }`)
}

func handleGetTableSchema(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `{
		"schemaName": "test",
		"enableColumnBasedNullHandling": false,
		"dimensionFieldSpecs": [
		  {
			"name": "ActualElapsedTime",
			"dataType": "INT",
			"notNull": false
		  },
		  {
			"name": "AirTime",
			"dataType": "INT",
			"notNull": false
		  },
		  {
			"name": "AirlineID",
			"dataType": "INT",
			"notNull": false
		  },
		  {
			"name": "ArrDel15",
			"dataType": "INT",
			"notNull": false
		  },
		  {
			"name": "ArrDelay",
			"dataType": "INT",
			"notNull": false
		  },
		  {
			"name": "ArrDelayMinutes",
			"dataType": "INT",
			"notNull": false
		  },
		  {
			"name": "ArrTime",
			"dataType": "INT",
			"notNull": false
		  },
		  {
			"name": "ArrTimeBlk",
			"dataType": "STRING",
			"notNull": false
		  },
		  {
			"name": "ArrivalDelayGroups",
			"dataType": "INT",
			"notNull": false
		  },
		  {
			"name": "CRSArrTime",
			"dataType": "INT",
			"notNull": false
		  },
		  {
			"name": "CRSDepTime",
			"dataType": "INT",
			"notNull": false
		  },
		  {
			"name": "CRSElapsedTime",
			"dataType": "INT",
			"notNull": false
		  },
		  {
			"name": "CancellationCode",
			"dataType": "STRING",
			"notNull": false
		  },
		  {
			"name": "Cancelled",
			"dataType": "INT",
			"notNull": false
		  },
		  {
			"name": "Carrier",
			"dataType": "STRING",
			"notNull": false
		  },
		  {
			"name": "CarrierDelay",
			"dataType": "INT",
			"notNull": false
		  },
		  {
			"name": "DayOfWeek",
			"dataType": "INT",
			"notNull": false
		  },
		  {
			"name": "DayofMonth",
			"dataType": "INT",
			"notNull": false
		  },
		  {
			"name": "DepDel15",
			"dataType": "INT",
			"notNull": false
		  },
		  {
			"name": "DepDelay",
			"dataType": "INT",
			"notNull": false
		  },
		  {
			"name": "DepDelayMinutes",
			"dataType": "INT",
			"notNull": false
		  },
		  {
			"name": "DepTime",
			"dataType": "INT",
			"notNull": false
		  },
		  {
			"name": "DepTimeBlk",
			"dataType": "STRING",
			"notNull": false
		  },
		  {
			"name": "DepartureDelayGroups",
			"dataType": "INT",
			"notNull": false
		  },
		  {
			"name": "Dest",
			"dataType": "STRING",
			"notNull": false
		  },
		  {
			"name": "DestAirportID",
			"dataType": "INT",
			"notNull": false
		  },
		  {
			"name": "DestAirportSeqID",
			"dataType": "INT",
			"notNull": false
		  },
		  {
			"name": "DestCityMarketID",
			"dataType": "INT",
			"notNull": false
		  },
		  {
			"name": "DestCityName",
			"dataType": "STRING",
			"notNull": false
		  },
		  {
			"name": "DestState",
			"dataType": "STRING",
			"notNull": false
		  },
		  {
			"name": "DestStateFips",
			"dataType": "INT",
			"notNull": false
		  },
		  {
			"name": "DestStateName",
			"dataType": "STRING",
			"notNull": false
		  },
		  {
			"name": "DestWac",
			"dataType": "INT",
			"notNull": false
		  },
		  {
			"name": "Distance",
			"dataType": "INT",
			"notNull": false
		  },
		  {
			"name": "DistanceGroup",
			"dataType": "INT",
			"notNull": false
		  },
		  {
			"name": "DivActualElapsedTime",
			"dataType": "INT",
			"notNull": false
		  },
		  {
			"name": "DivAirportIDs",
			"dataType": "INT",
			"singleValueField": false,
			"notNull": false
		  },
		  {
			"name": "DivAirportLandings",
			"dataType": "INT",
			"notNull": false
		  },
		  {
			"name": "DivAirportSeqIDs",
			"dataType": "INT",
			"singleValueField": false,
			"notNull": false
		  },
		  {
			"name": "DivAirports",
			"dataType": "STRING",
			"singleValueField": false,
			"notNull": false
		  },
		  {
			"name": "DivArrDelay",
			"dataType": "INT",
			"notNull": false
		  },
		  {
			"name": "DivDistance",
			"dataType": "INT",
			"notNull": false
		  },
		  {
			"name": "DivLongestGTimes",
			"dataType": "INT",
			"singleValueField": false,
			"notNull": false
		  },
		  {
			"name": "DivReachedDest",
			"dataType": "INT",
			"notNull": false
		  },
		  {
			"name": "DivTailNums",
			"dataType": "STRING",
			"singleValueField": false,
			"notNull": false
		  },
		  {
			"name": "DivTotalGTimes",
			"dataType": "INT",
			"singleValueField": false,
			"notNull": false
		  },
		  {
			"name": "DivWheelsOffs",
			"dataType": "INT",
			"singleValueField": false,
			"notNull": false
		  },
		  {
			"name": "DivWheelsOns",
			"dataType": "INT",
			"singleValueField": false,
			"notNull": false
		  },
		  {
			"name": "Diverted",
			"dataType": "INT",
			"notNull": false
		  },
		  {
			"name": "FirstDepTime",
			"dataType": "INT",
			"notNull": false
		  },
		  {
			"name": "FlightDate",
			"dataType": "STRING",
			"notNull": false
		  },
		  {
			"name": "FlightNum",
			"dataType": "INT",
			"notNull": false
		  },
		  {
			"name": "Flights",
			"dataType": "INT",
			"notNull": false
		  },
		  {
			"name": "LateAircraftDelay",
			"dataType": "INT",
			"notNull": false
		  },
		  {
			"name": "LongestAddGTime",
			"dataType": "INT",
			"notNull": false
		  },
		  {
			"name": "Month",
			"dataType": "INT",
			"notNull": false
		  },
		  {
			"name": "NASDelay",
			"dataType": "INT",
			"notNull": false
		  },
		  {
			"name": "Origin",
			"dataType": "STRING",
			"notNull": false
		  },
		  {
			"name": "OriginAirportID",
			"dataType": "INT",
			"notNull": false
		  },
		  {
			"name": "OriginAirportSeqID",
			"dataType": "INT",
			"notNull": false
		  },
		  {
			"name": "OriginCityMarketID",
			"dataType": "INT",
			"notNull": false
		  },
		  {
			"name": "OriginCityName",
			"dataType": "STRING",
			"notNull": false
		  },
		  {
			"name": "OriginState",
			"dataType": "STRING",
			"notNull": false
		  },
		  {
			"name": "OriginStateFips",
			"dataType": "INT",
			"notNull": false
		  },
		  {
			"name": "OriginStateName",
			"dataType": "STRING",
			"notNull": false
		  },
		  {
			"name": "OriginWac",
			"dataType": "INT",
			"notNull": false
		  },
		  {
			"name": "Quarter",
			"dataType": "INT",
			"notNull": false
		  },
		  {
			"name": "RandomAirports",
			"dataType": "STRING",
			"singleValueField": false,
			"notNull": false
		  },
		  {
			"name": "SecurityDelay",
			"dataType": "INT",
			"notNull": false
		  },
		  {
			"name": "TailNum",
			"dataType": "STRING",
			"notNull": false
		  },
		  {
			"name": "TaxiIn",
			"dataType": "INT",
			"notNull": false
		  },
		  {
			"name": "TaxiOut",
			"dataType": "INT",
			"notNull": false
		  },
		  {
			"name": "Year",
			"dataType": "INT",
			"notNull": false
		  },
		  {
			"name": "WheelsOn",
			"dataType": "INT",
			"notNull": false
		  },
		  {
			"name": "WheelsOff",
			"dataType": "INT",
			"notNull": false
		  },
		  {
			"name": "WeatherDelay",
			"dataType": "INT",
			"notNull": false
		  },
		  {
			"name": "UniqueCarrier",
			"dataType": "STRING",
			"notNull": false
		  },
		  {
			"name": "TotalAddGTime",
			"dataType": "INT",
			"notNull": false
		  }
		],
		"dateTimeFieldSpecs": [
		  {
			"name": "DaysSinceEpoch",
			"dataType": "INT",
			"notNull": false,
			"format": "1:DAYS:EPOCH",
			"granularity": "1:DAYS"
		  },
		  {
			"name": "ts",
			"dataType": "TIMESTAMP",
			"notNull": false,
			"format": "1:MILLISECONDS:TIMESTAMP",
			"granularity": "1:SECONDS"
		  },
		  {
			"name": "tsRaw",
			"dataType": "TIMESTAMP",
			"notNull": false,
			"format": "1:MILLISECONDS:TIMESTAMP",
			"granularity": "1:SECONDS"
		  }
		]
	  }`)
}

func handleGetTableSize(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `{
		"tableName": "test",
		"reportedSizeInBytes": 4723495,
		"estimatedSizeInBytes": 4723495,
		"reportedSizePerReplicaInBytes": 4723495,
		"offlineSegments": {
		  "reportedSizeInBytes": 4723495,
		  "estimatedSizeInBytes": 4723495,
		  "missingSegments": 0,
		  "reportedSizePerReplicaInBytes": 4723495,
		  "segments": {
			"test_OFFLINE_16071_16071_0": {
			  "reportedSizeInBytes": 148141,
			  "estimatedSizeInBytes": 148141,
			  "maxReportedSizePerReplicaInBytes": 148141,
			  "serverInfo": {
				"Server_172.17.0.3_7050": {
				  "segmentName": "test_OFFLINE_16071_16071_0",
				  "diskSizeInBytes": 148141
				}
			  }
			},
			"test_OFFLINE_16072_16072_0": {
			  "reportedSizeInBytes": 170479,
			  "estimatedSizeInBytes": 170479,
			  "maxReportedSizePerReplicaInBytes": 170479,
			  "serverInfo": {
				"Server_172.17.0.3_7050": {
				  "segmentName": "test_OFFLINE_16072_16072_0",
				  "diskSizeInBytes": 170479
				}
			  }
			},
			"test_OFFLINE_16074_16074_0": {
			  "reportedSizeInBytes": 152313,
			  "estimatedSizeInBytes": 152313,
			  "maxReportedSizePerReplicaInBytes": 152313,
			  "serverInfo": {
				"Server_172.17.0.3_7050": {
				  "segmentName": "test_OFFLINE_16074_16074_0",
				  "diskSizeInBytes": 152313
				}
			  }
			},
			"test_OFFLINE_16081_16081_0": {
			  "reportedSizeInBytes": 143424,
			  "estimatedSizeInBytes": 143424,
			  "maxReportedSizePerReplicaInBytes": 143424,
			  "serverInfo": {
				"Server_172.17.0.3_7050": {
				  "segmentName": "test_OFFLINE_16081_16081_0",
				  "diskSizeInBytes": 143424
				}
			  }
			},
			"test_OFFLINE_16073_16073_0": {
			  "reportedSizeInBytes": 172148,
			  "estimatedSizeInBytes": 172148,
			  "maxReportedSizePerReplicaInBytes": 172148,
			  "serverInfo": {
				"Server_172.17.0.3_7050": {
				  "segmentName": "test_OFFLINE_16073_16073_0",
				  "diskSizeInBytes": 172148
				}
			  }
			},
			"test_OFFLINE_16083_16083_0": {
			  "reportedSizeInBytes": 146487,
			  "estimatedSizeInBytes": 146487,
			  "maxReportedSizePerReplicaInBytes": 146487,
			  "serverInfo": {
				"Server_172.17.0.3_7050": {
				  "segmentName": "test_OFFLINE_16083_16083_0",
				  "diskSizeInBytes": 146487
				}
			  }
			},
			"test_OFFLINE_16082_16082_0": {
			  "reportedSizeInBytes": 144331,
			  "estimatedSizeInBytes": 144331,
			  "maxReportedSizePerReplicaInBytes": 144331,
			  "serverInfo": {
				"Server_172.17.0.3_7050": {
				  "segmentName": "test_OFFLINE_16082_16082_0",
				  "diskSizeInBytes": 144331
				}
			  }
			},
			"test_OFFLINE_16077_16077_0": {
			  "reportedSizeInBytes": 147533,
			  "estimatedSizeInBytes": 147533,
			  "maxReportedSizePerReplicaInBytes": 147533,
			  "serverInfo": {
				"Server_172.17.0.3_7050": {
				  "segmentName": "test_OFFLINE_16077_16077_0",
				  "diskSizeInBytes": 147533
				}
			  }
			},
			"test_OFFLINE_16076_16076_0": {
			  "reportedSizeInBytes": 160686,
			  "estimatedSizeInBytes": 160686,
			  "maxReportedSizePerReplicaInBytes": 160686,
			  "serverInfo": {
				"Server_172.17.0.3_7050": {
				  "segmentName": "test_OFFLINE_16076_16076_0",
				  "diskSizeInBytes": 160686
				}
			  }
			},
			"test_OFFLINE_16085_16085_0": {
			  "reportedSizeInBytes": 155212,
			  "estimatedSizeInBytes": 155212,
			  "maxReportedSizePerReplicaInBytes": 155212,
			  "serverInfo": {
				"Server_172.17.0.3_7050": {
				  "segmentName": "test_OFFLINE_16085_16085_0",
				  "diskSizeInBytes": 155212
				}
			  }
			},
			"test_OFFLINE_16092_16092_0": {
			  "reportedSizeInBytes": 146841,
			  "estimatedSizeInBytes": 146841,
			  "maxReportedSizePerReplicaInBytes": 146841,
			  "serverInfo": {
				"Server_172.17.0.3_7050": {
				  "segmentName": "test_OFFLINE_16092_16092_0",
				  "diskSizeInBytes": 146841
				}
			  }
			},
			"test_OFFLINE_16093_16093_0": {
			  "reportedSizeInBytes": 152480,
			  "estimatedSizeInBytes": 152480,
			  "maxReportedSizePerReplicaInBytes": 152480,
			  "serverInfo": {
				"Server_172.17.0.3_7050": {
				  "segmentName": "test_OFFLINE_16093_16093_0",
				  "diskSizeInBytes": 152480
				}
			  }
			},
			"test_OFFLINE_16075_16075_0": {
			  "reportedSizeInBytes": 173926,
			  "estimatedSizeInBytes": 173926,
			  "maxReportedSizePerReplicaInBytes": 173926,
			  "serverInfo": {
				"Server_172.17.0.3_7050": {
				  "segmentName": "test_OFFLINE_16075_16075_0",
				  "diskSizeInBytes": 173926
				}
			  }
			},
			"test_OFFLINE_16084_16084_0": {
			  "reportedSizeInBytes": 151416,
			  "estimatedSizeInBytes": 151416,
			  "maxReportedSizePerReplicaInBytes": 151416,
			  "serverInfo": {
				"Server_172.17.0.3_7050": {
				  "segmentName": "test_OFFLINE_16084_16084_0",
				  "diskSizeInBytes": 151416
				}
			  }
			},
			"test_OFFLINE_16089_16089_0": {
			  "reportedSizeInBytes": 140200,
			  "estimatedSizeInBytes": 140200,
			  "maxReportedSizePerReplicaInBytes": 140200,
			  "serverInfo": {
				"Server_172.17.0.3_7050": {
				  "segmentName": "test_OFFLINE_16089_16089_0",
				  "diskSizeInBytes": 140200
				}
			  }
			},
			"test_OFFLINE_16101_16101_0": {
			  "reportedSizeInBytes": 161946,
			  "estimatedSizeInBytes": 161946,
			  "maxReportedSizePerReplicaInBytes": 161946,
			  "serverInfo": {
				"Server_172.17.0.3_7050": {
				  "segmentName": "test_OFFLINE_16101_16101_0",
				  "diskSizeInBytes": 161946
				}
			  }
			},
			"test_OFFLINE_16096_16096_0": {
			  "reportedSizeInBytes": 157100,
			  "estimatedSizeInBytes": 157100,
			  "maxReportedSizePerReplicaInBytes": 157100,
			  "serverInfo": {
				"Server_172.17.0.3_7050": {
				  "segmentName": "test_OFFLINE_16096_16096_0",
				  "diskSizeInBytes": 157100
				}
			  }
			},
			"test_OFFLINE_16088_16088_0": {
			  "reportedSizeInBytes": 135803,
			  "estimatedSizeInBytes": 135803,
			  "maxReportedSizePerReplicaInBytes": 135803,
			  "serverInfo": {
				"Server_172.17.0.3_7050": {
				  "segmentName": "test_OFFLINE_16088_16088_0",
				  "diskSizeInBytes": 135803
				}
			  }
			},
			"test_OFFLINE_16094_16094_0": {
			  "reportedSizeInBytes": 159093,
			  "estimatedSizeInBytes": 159093,
			  "maxReportedSizePerReplicaInBytes": 159093,
			  "serverInfo": {
				"Server_172.17.0.3_7050": {
				  "segmentName": "test_OFFLINE_16094_16094_0",
				  "diskSizeInBytes": 159093
				}
			  }
			},
			"test_OFFLINE_16091_16091_0": {
			  "reportedSizeInBytes": 148804,
			  "estimatedSizeInBytes": 148804,
			  "maxReportedSizePerReplicaInBytes": 148804,
			  "serverInfo": {
				"Server_172.17.0.3_7050": {
				  "segmentName": "test_OFFLINE_16091_16091_0",
				  "diskSizeInBytes": 148804
				}
			  }
			},
			"test_OFFLINE_16086_16086_0": {
			  "reportedSizeInBytes": 152318,
			  "estimatedSizeInBytes": 152318,
			  "maxReportedSizePerReplicaInBytes": 152318,
			  "serverInfo": {
				"Server_172.17.0.3_7050": {
				  "segmentName": "test_OFFLINE_16086_16086_0",
				  "diskSizeInBytes": 152318
				}
			  }
			},
			"test_OFFLINE_16087_16087_0": {
			  "reportedSizeInBytes": 146948,
			  "estimatedSizeInBytes": 146948,
			  "maxReportedSizePerReplicaInBytes": 146948,
			  "serverInfo": {
				"Server_172.17.0.3_7050": {
				  "segmentName": "test_OFFLINE_16087_16087_0",
				  "diskSizeInBytes": 146948
				}
			  }
			},
			"test_OFFLINE_16090_16090_0": {
			  "reportedSizeInBytes": 147078,
			  "estimatedSizeInBytes": 147078,
			  "maxReportedSizePerReplicaInBytes": 147078,
			  "serverInfo": {
				"Server_172.17.0.3_7050": {
				  "segmentName": "test_OFFLINE_16090_16090_0",
				  "diskSizeInBytes": 147078
				}
			  }
			},
			"test_OFFLINE_16095_16095_0": {
			  "reportedSizeInBytes": 133078,
			  "estimatedSizeInBytes": 133078,
			  "maxReportedSizePerReplicaInBytes": 133078,
			  "serverInfo": {
				"Server_172.17.0.3_7050": {
				  "segmentName": "test_OFFLINE_16095_16095_0",
				  "diskSizeInBytes": 133078
				}
			  }
			},
			"test_OFFLINE_16078_16078_0": {
			  "reportedSizeInBytes": 154535,
			  "estimatedSizeInBytes": 154535,
			  "maxReportedSizePerReplicaInBytes": 154535,
			  "serverInfo": {
				"Server_172.17.0.3_7050": {
				  "segmentName": "test_OFFLINE_16078_16078_0",
				  "diskSizeInBytes": 154535
				}
			  }
			},
			"test_OFFLINE_16079_16079_0": {
			  "reportedSizeInBytes": 160246,
			  "estimatedSizeInBytes": 160246,
			  "maxReportedSizePerReplicaInBytes": 160246,
			  "serverInfo": {
				"Server_172.17.0.3_7050": {
				  "segmentName": "test_OFFLINE_16079_16079_0",
				  "diskSizeInBytes": 160246
				}
			  }
			},
			"test_OFFLINE_16080_16080_0": {
			  "reportedSizeInBytes": 159886,
			  "estimatedSizeInBytes": 159886,
			  "maxReportedSizePerReplicaInBytes": 159886,
			  "serverInfo": {
				"Server_172.17.0.3_7050": {
				  "segmentName": "test_OFFLINE_16080_16080_0",
				  "diskSizeInBytes": 159886
				}
			  }
			},
			"test_OFFLINE_16099_16099_0": {
			  "reportedSizeInBytes": 151587,
			  "estimatedSizeInBytes": 151587,
			  "maxReportedSizePerReplicaInBytes": 151587,
			  "serverInfo": {
				"Server_172.17.0.3_7050": {
				  "segmentName": "test_OFFLINE_16099_16099_0",
				  "diskSizeInBytes": 151587
				}
			  }
			},
			"test_OFFLINE_16098_16098_0": {
			  "reportedSizeInBytes": 147585,
			  "estimatedSizeInBytes": 147585,
			  "maxReportedSizePerReplicaInBytes": 147585,
			  "serverInfo": {
				"Server_172.17.0.3_7050": {
				  "segmentName": "test_OFFLINE_16098_16098_0",
				  "diskSizeInBytes": 147585
				}
			  }
			},
			"test_OFFLINE_16097_16097_0": {
			  "reportedSizeInBytes": 143544,
			  "estimatedSizeInBytes": 143544,
			  "maxReportedSizePerReplicaInBytes": 143544,
			  "serverInfo": {
				"Server_172.17.0.3_7050": {
				  "segmentName": "test_OFFLINE_16097_16097_0",
				  "diskSizeInBytes": 143544
				}
			  }
			},
			"test_OFFLINE_16100_16100_0": {
			  "reportedSizeInBytes": 158327,
			  "estimatedSizeInBytes": 158327,
			  "maxReportedSizePerReplicaInBytes": 158327,
			  "serverInfo": {
				"Server_172.17.0.3_7050": {
				  "segmentName": "test_OFFLINE_16100_16100_0",
				  "diskSizeInBytes": 158327
				}
			  }
			}
		  }
		},
		"realtimeSegments": null
	  }`)
}

func handleGetTableState(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `{
		"state": "enabled"
	  }`)
}

func handleChangeTableState(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `{
		"status": "Request to enable table 'test_OFFLINE' is successful"
	  }`)
}

func handleGetTableStats(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `{
		"OFFLINE": {
		  "creationTime": "20240411T224913Z"
		}
	  }`)
}

func handleResetTableSegment(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `{
		"status": "Successfully reset segment: test_OFFLINE_16071_16071_0 of table: test_OFFLINE"
	  }`)
}

func handleResetTableSegments(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `{
		"status": "Successfully reset segments of table: test_OFFLINE"
	  }`)
}

func createMockControllerServer() *httptest.Server {

	mux := http.NewServeMux()

	mux.HandleFunc(RouteUsers, authMiddleware(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			handleGetUsers(w, r)
		case "POST":
			handlePostUsers(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	mux.HandleFunc(RouteUser, authMiddleware(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			handleGetUser(w, r)
		case "PUT":
			handlePutUser(w, r)
		case "DELETE":
			handleDeleteUser(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	mux.HandleFunc(RouteInstances, authMiddleware(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			handleGetInstances(w, r)
		case "POST":
			handlePostInstances(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	mux.HandleFunc(RouteInstance, authMiddleware(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			handleGetInstance(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	mux.HandleFunc(RouteClusterInfo, authMiddleware(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			handleGetClusterInfo(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	mux.HandleFunc(RouteClusterConfigs, authMiddleware(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			handleGetClusterConfigs(w, r)
		case "POST":
			handlePostClusterConfigs(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	mux.HandleFunc(RouteClusterConfigsDelete, authMiddleware(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "DELETE":
			handleDeleteClusterConfigs(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	mux.HandleFunc(RouteTenants, authMiddleware(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			handleGetTenants(w, r)
		case "POST":
			handleCreateTenant(w, r)
		case "PUT":
			handleUpdateTenant(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	mux.HandleFunc(RouteTenantsInstances, authMiddleware(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			handleGetTenantInstances(w, r)
		case "DELETE":
			handleDeleteTenant(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	mux.HandleFunc(RouteTenantsTables, authMiddleware(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			handleGetTenantTables(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	mux.HandleFunc(RouteTenantsMetadata, authMiddleware(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			handleGetTenantMetadata(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	mux.HandleFunc(RouteSegmentsTest, authMiddleware(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			handleGetSegments(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	mux.HandleFunc(RouteSegmentsTestReload, authMiddleware(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			handleReloadTableSegments(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	mux.HandleFunc(RouteSegmentTestReload, authMiddleware(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			handleReloadTableSegment(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	mux.HandleFunc(RouteSchemas, authMiddleware(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			handleGetSchemas(w, r)
		case "POST":
			handleCreateSchema(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	mux.HandleFunc(RouteSchemasTest, authMiddleware(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			handleGetSchema(w, r)
		case "PUT":
			handleUpdateSchema(w, r)
		case "DELETE":
			handleDeleteSchema(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	mux.HandleFunc(RouteTables, authMiddleware(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			handleGetTables(w, r)
		case "POST":
			handleCreateTable(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	mux.HandleFunc(RouteTablesTest, authMiddleware(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			handleGetTableConfig(w, r)
		case "PUT":
			handleUpdateTable(w, r)
		case "DELETE":
			handleDeleteTable(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	mux.HandleFunc(RouteTablesTestExternalView, authMiddleware(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			handleTableExternalView(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	mux.HandleFunc(RoutePinotControllerAdmin, authMiddleware(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			handlePinotControllerAdmin(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	mux.HandleFunc(RouteHealth, authMiddleware(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			handleHealth(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	mux.HandleFunc(RouteSchemasFieldSpec, authMiddleware(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			handleGetFieldSpecs(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	mux.HandleFunc(RouteTablesTestIdealState, authMiddleware(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			handleTableIdealState(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	mux.HandleFunc(RouteTablesTestIndexes, authMiddleware(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			handleTableIndexes(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	mux.HandleFunc(RouteTablesTestInstances, authMiddleware(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			handleGetTableInstances(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	mux.HandleFunc(RouteTablesLiveBrokers, authMiddleware(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			handleGetTableLiveBrokers(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	mux.HandleFunc(RouteTablesTestLiveBrokers, authMiddleware(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			handleGetTableTestLiveBrokers(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	mux.HandleFunc(RouteTablesTestMetadata, authMiddleware(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			handleGetTableTestMetadata(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	mux.HandleFunc(RouteTablesTestRebuildBrokerResourceFromHelixTags, authMiddleware(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			handleTableRebuildBrokerResourceFromHelixTags(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	mux.HandleFunc(RouteTablesTestSchema, authMiddleware(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			handleGetTableSchema(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	mux.HandleFunc(RouteTablesTestSize, authMiddleware(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			handleGetTableSize(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	mux.HandleFunc(RouteTablesTestState, authMiddleware(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			handleGetTableState(w, r)
		case "PUT":
			handleChangeTableState(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	mux.HandleFunc(RouteTablesTestStats, authMiddleware(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			handleGetTableStats(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	mux.HandleFunc(RouteSegmentTestReset, authMiddleware(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			handleResetTableSegment(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	mux.HandleFunc(RouteSegmentTestResetAll, authMiddleware(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			handleResetTableSegments(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	return httptest.NewServer(mux)

}

func createPinotClient(server *httptest.Server) *goPinotAPI.PinotAPIClient {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	return goPinotAPI.NewPinotAPIClient(
		goPinotAPI.ControllerUrl(server.URL),
		goPinotAPI.AuthToken("YWRtaW46YWRtaW4K"),
		goPinotAPI.Logger(logger),
	)
}

func getSchema() model.Schema {

	schemaFilePath := "./example/data-gen/block_header_schema.json"

	f, err := os.Open(schemaFilePath)
	if err != nil {
		log.Panic(err)
	}

	defer f.Close()

	var schema model.Schema
	err = json.NewDecoder(f).Decode(&schema)
	if err != nil {
		log.Panic(err)
	}

	return schema
}

func boolPtr(b bool) *bool {
	return &b
}

// TestFetchData
// func TestFetchData(t *testing.T) {

// }

// test FetchData unauthorized

// Test GetUsers
func TestGetUsers(t *testing.T) {

	server := createMockControllerServer()
	client := createPinotClient(server)

	res, err := client.GetUsers()

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// expect 1 user in the response username test
	assert.Equal(t, len(res.Users), 1, "Expected 1 user in the response")
	assert.Equal(t, res.Users["test_BROKER"].Username, "test", "Expected username to be test")

}

// Test GetUser
func TestGetUser(t *testing.T) {

	server := createMockControllerServer()
	client := createPinotClient(server)

	res, err := client.GetUser("test", "BROKER")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// expect username test
	// expect component BROKER
	assert.Equal(t, res.Username, "test", "Expected username to be test")
	assert.Equal(t, res.Component, "BROKER", "Expected component to be BROKER")

}

// Test CreateUser
func TestCreateUser(t *testing.T) {

	server := createMockControllerServer()
	client := createPinotClient(server)

	user := model.User{
		Username:  "testUser",
		Password:  "test",
		Component: "BROKER",
		Role:      "ADMIN",
	}

	userBytes, err := json.Marshal(user)
	if err != nil {
		t.Errorf("Couldn't marshal user: %v", err)
	}

	res, err := client.CreateUser(userBytes)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// expect User testUser_BROKER has been successfully added!
	assert.Equal(t, res.Status, "User testUser_BROKER has been successfully added!", "Expected response to be User testUser_BROKER has been successfully added!")
}

// Test UpdateUser
func TestUpdateUser(t *testing.T) {
	server := createMockControllerServer()
	client := createPinotClient(server)

	user := model.User{
		Username:  "test",
		Password:  "test",
		Component: "BROKER",
		Role:      "USER",
	}

	userBytes, err := json.Marshal(user)
	if err != nil {
		t.Errorf("Couldn't marshal user: %v", err)
	}

	res, err := client.UpdateUser(user.Username, user.Component, true, userBytes)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// expect User test_BROKER has been successfully updated!
	assert.Equal(t, res.Status, "User config update for test_BROKER", "Expected response to be User config update for test_BROKER")
}

// Test DeleteUser
func TestDeleteUser(t *testing.T) {
	server := createMockControllerServer()
	client := createPinotClient(server)

	res, err := client.DeleteUser("test", "BROKER")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// expect User test_BROKER has been successfully deleted!
	assert.Equal(t, res.Status, "User: test_BROKER has been successfully deleted", "Expected response to be User: test_BROKER has been successfully deleted!")

}

func TestDeleteUserNoComponent(t *testing.T) {
	server := createMockControllerServer()
	client := createPinotClient(server)

	_, err := client.DeleteUser("test", "")
	if err != nil {
		assert.Equal(t, err.Error(), "client: request failed: status 400\n{\"code\": 400,\"error\": \"Name is null\"}\n", "Expected error to be Name is null")
	}

}

func TestGetInstances(t *testing.T) {
	server := createMockControllerServer()
	client := createPinotClient(server)

	res, err := client.GetInstances()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	assert.Equal(t, len(res.Instances), 4, "Expected 4 instances in the response")
}

func TestGetInstance(t *testing.T) {
	server := createMockControllerServer()
	client := createPinotClient(server)

	res, err := client.GetInstance("Minion_172.19.0.2_9514")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	assert.Equal(t, res.InstanceName, "Minion_172.19.0.2_9514", "Expected instance name to be Minion_172.19.0.2_9514")
}

func TestCreateInstance(t *testing.T) {
	server := createMockControllerServer()
	client := createPinotClient(server)

	instance := model.Instance{
		Host: "localhost",
		Port: 1234,
		Type: "BROKER",
	}

	instanceBytes, err := json.Marshal(instance)
	if err != nil {
		t.Errorf("Couldn't marshal instance: %v", err)
	}

	res, err := client.CreateInstance(instanceBytes)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	assert.Equal(t, res.Status, "Added instance: Broker_localhost_1234", "Expected response to be Added instance: Broker_localhost_1234")

}

// DeleteInstance

func TestGetClusterInfo(t *testing.T) {
	server := createMockControllerServer()
	client := createPinotClient(server)

	res, err := client.GetClusterInfo()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	assert.Equal(t, res.ClusterName, "PinotCluster", "Expected cluster name to be PinotCluster")
}

func TestGetClusterConfigs(t *testing.T) {
	server := createMockControllerServer()
	client := createPinotClient(server)

	res, err := client.GetClusterConfigs()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	assert.Equal(t, res.AllowParticipantAutoJoin, "true", "Expected allowParticipantAutoJoin to be true")
}

func TestUpdateClusterConfig(t *testing.T) {
	server := createMockControllerServer()
	client := createPinotClient(server)

	config := model.ClusterConfig{
		AllowParticipantAutoJoin: "false",
	}

	configBytes, err := json.Marshal(config)
	if err != nil {
		t.Errorf("Couldn't marshal config: %v", err)
	}

	res, err := client.UpdateClusterConfigs(configBytes)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	assert.Equal(t, res.Status, "Updated cluster config.", "Expected response to be Updated cluster config.")

}

func TestDeleteClusterConfig(t *testing.T) {
	server := createMockControllerServer()
	client := createPinotClient(server)

	res, err := client.DeleteClusterConfig("allowParticipantAutoJoin")
	if err != nil {
		assert.Equal(t, err.Error(), "client: request failed: status 404\n404 page not found\n", "Expected error to be 404 page not found")
	}

	assert.Equal(t, res.Status, "Deleted cluster config: allowParticipantAutoJoin", "Expected response to be Deleted cluster config: allowParticipantAutoJoin")
}

func TestGetTenants(t *testing.T) {
	server := createMockControllerServer()
	client := createPinotClient(server)

	res, err := client.GetTenants()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	assert.Equal(t, len(res.BrokerTenants), 1, "Expected 1 Broker tenants in the response")
	assert.Equal(t, len(res.ServerTenants), 1, "Expected 1 Server tenants in the response")

}

// TestGetTenantInstances
func TestGetTenantInstances(t *testing.T) {
	server := createMockControllerServer()
	client := createPinotClient(server)

	res, err := client.GetTenantInstances("DefaultTenant")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	assert.Equal(t, len(res.BrokerInstances), 1, "Expected 1 Broker instance in the response")
	assert.Equal(t, len(res.ServerInstances), 1, "Expected 1 Server instance in the response")
	assert.Equal(t, res.TenantName, "DefaultTenant", "Expected tenant name to be DefaultTenant")
}

// TestGetTenantTables
func TestGetTenantTables(t *testing.T) {
	server := createMockControllerServer()
	client := createPinotClient(server)

	res, err := client.GetTenantTables("DefaultTenant")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	assert.Equal(t, len(res.Tables), 0, "Expected 0 tables in the response")
}

// TestGetTenantMetadata
func TestGetTenantMetadata(t *testing.T) {
	server := createMockControllerServer()
	client := createPinotClient(server)

	res, err := client.GetTenantMetadata("DefaultTenant")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	assert.Equal(t, len(res.BrokerInstances), 1, "Expected 1 Broker instance in the response")
	assert.Equal(t, len(res.ServerInstances), 1, "Expected 1 Server instance in the response")
	assert.Equal(t, res.OfflineServerInstances, []string([]string(nil)), "Expected OfflineServerInstances to be nil")
	assert.Equal(t, res.RealtimeServerInstances, []string([]string(nil)), "Expected RealtimeServerInstances to be nil")
	assert.Equal(t, res.TenantName, "DefaultTenant", "Expected tenant name to be DefaultTenant")
}

// TestCreateTenant
func TestCreateTenant(t *testing.T) {
	server := createMockControllerServer()
	client := createPinotClient(server)

	tenant := model.Tenant{
		TenantName: "test",
		TenantRole: "BROKER",
	}

	tenantBytes, err := json.Marshal(tenant)
	if err != nil {
		t.Errorf("Couldn't marshal tenant: %v", err)
	}

	res, err := client.CreateTenant(tenantBytes)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	assert.Equal(t, res.Status, "Successfully created tenant", "Expected response to be Successfully created tenant")
}

// TestUpdateTenant
func TestUpdateTenant(t *testing.T) {
	server := createMockControllerServer()
	client := createPinotClient(server)

	updateTenant := model.Tenant{
		TenantName: "test",
		TenantRole: "SERVER",
	}

	tenantBytes, err := json.Marshal(updateTenant)
	if err != nil {
		t.Errorf("Couldn't marshal tenant: %v", err)
	}

	res, err := client.UpdateTenant(tenantBytes)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	assert.Equal(t, res.Status, "Updated tenant", "Expected response to be Successfully updated tenant")
}

// TestDeleteTenant
func TestDeleteTenant(t *testing.T) {
	server := createMockControllerServer()
	client := createPinotClient(server)

	res, err := client.DeleteTenant("DefaultTenant", "SERVER")
	if err != nil {
		assert.Equal(t, err.Error(), "client: request failed: status 404\n404 page not found\n", "Expected error to be 404 page not found")
	}

	assert.Equal(t, res.Status, "Successfully deleted tenant DefaultTenant", "Expected response to be Successfully deleted tenant DefaultTenant")
}

// TestRebalanceTenant

// TestSegments
func TestSegments(t *testing.T) {
	server := createMockControllerServer()
	client := createPinotClient(server)

	res, err := client.GetSegments("test")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	assert.Equal(t, len(res[0].Offline), 1, "Expected 1 offline segment in the response")
	assert.Equal(t, len(res[0].Realtime), 1, "Expected 1 realtime segment in the response")
}

// TestReloadTableSegments
func TestReloadTableSegments(t *testing.T) {
	server := createMockControllerServer()
	client := createPinotClient(server)

	res, err := client.ReloadTableSegments("test")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	assert.Equal(t, res.Status, "{\"test_OFFLINE\":{\"reloadJobId\":\"f9db13c7-3ad6-45a8-a08f-75cd03c42fb5\",\"reloadJobMetaZKStorageStatus\":\"SUCCESS\",\"numMessagesSent\":\"1\"}}", "Expected response to be {\"test_OFFLINE\":{\"reloadJobId\":\"f9db13c7-3ad6-45a8-a08f-75cd03c42fb5\",\"reloadJobMetaZKStorageStatus\":\"SUCCESS\",\"numMessagesSent\":\"1\"}}")

}

// TestReloadTableSegment
func TestReloadTableSegment(t *testing.T) {
	server := createMockControllerServer()
	client := createPinotClient(server)

	res, err := client.ReloadSegment("test", "test_1")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	assert.Equal(t, res.Status, "Submitted reload job id: ce4650a9-774a-4b22-919b-4cb22b5c8129, sent 1 reload messages. Job meta ZK storage status: SUCCESS", "Expected response to be Submitted reload job id: ce4650a9-774a-4b22-919b-4cb22b5c8129, sent 1 reload messages. Job meta ZK storage status: SUCCESS")
}

// TestGetSchemas
func TestGetSchemas(t *testing.T) {
	server := createMockControllerServer()
	client := createPinotClient(server)

	schemas := []model.Schema{}

	res, err := client.GetSchemas()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	res.ForEachSchema(func(schemaName string) {

		schemaResp, _ := client.GetSchema(schemaName)

		// fmt.Println("Reading Schema:")
		schemas = append(schemas, *schemaResp)

	})

	assert.Equal(t, len(schemas), 1, "Expected 1 schema in the response")
}

// TestGetSchema
func TestGetSchema(t *testing.T) {
	server := createMockControllerServer()
	client := createPinotClient(server)

	res, err := client.GetSchema("test")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	assert.Equal(t, res.SchemaName, "test", "Expected schema name to be test")
}

// TestCreateSchema
// it appears that this is not returning the status...
func TestCreateSchema(t *testing.T) {
	server := createMockControllerServer()
	client := createPinotClient(server)

	schema := model.Schema{
		SchemaName: "test",
		DimensionFieldSpecs: []model.FieldSpec{
			{
				Name:     "test",
				DataType: "STRING",
			},
			{
				Name:             "multiValuedField",
				DataType:         "STRING",
				SingleValueField: boolPtr(false),
			},
		},
		MetricFieldSpecs: []model.FieldSpec{
			{
				Name:     "test",
				DataType: "INT",
			},
		},
	}

	schemaBytes, err := json.Marshal(schema)
	if err != nil {
		t.Errorf("Couldn't marshal schema: %v", err)
	}

	res, err := client.CreateSchemaFromBytes(schemaBytes)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// fmt.Println("Response: ", res.Status)

	assert.Equal(t, res.Status, "test successfully added", "Expected response to be test successfully added")

}

// // TestCreateSchemaFromFile
// func TestCreateSchemaFromFile(t *testing.T) {
// 	server := createMockControllerServer()
// 	client := createPinotClient(server)

// 	res, err := client.CreateSchemaFromFile("test")
// 	if err != nil {
// 		t.Errorf("Expected no error, got %v", err)
// 	}

// 	assert.Equal(t, res.Status, "Schema test has been successfully added!", "Expected response to be Schema test has been successfully added!")
// }

// TestUpdateSchema
func TestUpdateSchema(t *testing.T) {
	server := createMockControllerServer()
	client := createPinotClient(server)

	schema := model.Schema{
		SchemaName: "test",
		DimensionFieldSpecs: []model.FieldSpec{
			{
				Name:     "test",
				DataType: "STRING",
			},
		},
		MetricFieldSpecs: []model.FieldSpec{
			{
				Name:     "test",
				DataType: "INT",
			},
		},
	}

	schemaBytes, err := json.Marshal(schema)
	if err != nil {
		t.Errorf("Couldn't marshal schema: %v", err)
	}

	res, err := client.UpdateSchemaFromBytes(schemaBytes)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// fmt.Println("Response: ", res.Status)

	assert.Equal(t, res.Status, "test successfully added", "Expected response to be test successfully added")
}

// TestDeleteSchema
// Requires /tables to be implemented first
// func TestDeleteSchema(t *testing.T) {
// 	server := createMockControllerServer()
// 	client := createPinotClient(server)

// 	res, err := client.DeleteSchema("test")
// 	if err != nil {
// 		t.Errorf("Expected no error, got %v", err)
// 	}

// 	assert.Equal(t, res.Status, "Schema test deleted", "Expected response to be Schema test deleted")
// }

// TestValidateSchema
func TestValidateSchema(t *testing.T) {
	server := createMockControllerServer()
	client := createPinotClient(server)

	schema := model.Schema{
		SchemaName: "test",
		DimensionFieldSpecs: []model.FieldSpec{
			{
				Name:     "test",
				DataType: "STRING",
			},
		},
		MetricFieldSpecs: []model.FieldSpec{
			{
				Name:     "test",
				DataType: "INT",
			},
		},
	}

	// schemaBytes, err := json.Marshal(schema)
	// if err != nil {
	// 	t.Errorf("Couldn't marshal schema: %v", err)
	// }

	res, err := client.ValidateSchema(schema)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	assert.Equal(t, res.Ok, true, "Expected response to be Schema test is valid!")
}

func TestGetTables(t *testing.T) {
	server := createMockControllerServer()
	client := createPinotClient(server)

	res, err := client.GetTables()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	assert.Equal(t, len(res.Tables), 1, "Expected 1 table in the response")
}

func TestGetOfflineTable(t *testing.T) {
	server := createMockControllerServer()
	client := createPinotClient(server)

	res, err := client.GetTable("test")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	assert.Equal(t, res.OFFLINE.TableName, "test_OFFLINE", "Expected table name to be test")

}

func TestGetRealtimeTable(t *testing.T) {
	server := createMockControllerServer()
	client := createPinotClient(server)

	res, err := client.GetTable("test")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	assert.Equal(t, res.REALTIME.TableName, "realtime_ethereum_mainnet_block_headers_REALTIME", "Expected table name to be test")

}

func TestCreateTable(t *testing.T) {
	server := createMockControllerServer()
	client := createPinotClient(server)

	// Create Schema First
	schema := model.Schema{
		SchemaName: "ethereum_mainnet_block_headers",
		DimensionFieldSpecs: []model.FieldSpec{
			{
				Name:     "number",
				DataType: "LONG",
				NotNull:  boolPtr(false),
			},
			{
				Name:     "hash",
				DataType: "STRING",
				NotNull:  boolPtr(false),
			},
			{
				Name:     "parent_hash",
				DataType: "STRING",
				NotNull:  boolPtr(false),
			},
			{
				Name:             "tags",
				DataType:         "STRING",
				SingleValueField: boolPtr(false),
			},
		},
		MetricFieldSpecs: []model.FieldSpec{
			{
				Name:     "gas_used",
				DataType: "LONG",
				NotNull:  boolPtr(false),
			},
		},
		DateTimeFieldSpecs: []model.FieldSpec{
			{
				Name:        "timestamp",
				DataType:    "LONG",
				NotNull:     boolPtr(false),
				Format:      "1:MILLISECONDS:EPOCH",
				Granularity: "1:MILLISECONDS",
			},
		},
	}

	schemaBytes, err := json.Marshal(schema)
	if err != nil {
		log.Panic(err)
	}

	_, err = client.CreateSchemaFromBytes(schemaBytes)
	if err != nil {
		log.Panic(err)
	}

	table := model.Table{
		TableName: "ethereum_mainnet_block_headers",
		TableType: "OFFLINE",
		SegmentsConfig: model.TableSegmentsConfig{
			TimeColumnName:            "timestamp",
			TimeType:                  "MILLISECONDS",
			Replication:               "1",
			SegmentAssignmentStrategy: "BalanceNumSegmentAssignmentStrategy",
			SegmentPushType:           "APPEND",
			MinimizeDataMovement:      true,
		},
		Tenants: model.TableTenant{
			Broker: "DefaultTenant",
			Server: "DefaultTenant",
		},
		TableIndexConfig: model.TableIndexConfig{
			LoadMode: "MMAP",
		},
		Metadata: &model.TableMetadata{
			CustomConfigs: map[string]string{
				"customKey": "customValue",
			},
		},
		FieldConfigList: []model.FieldConfig{
			{
				Name:         "number",
				EncodingType: "RAW",
				IndexType:    "SORTED",
			},
		},
		IngestionConfig: &model.TableIngestionConfig{
			SegmentTimeValueCheck: true,
			TransformConfigs: []model.TransformConfig{
				{
					ColumnName:        "timestamp",
					TransformFunction: "toEpochHours(millis)",
				},
			},
			ContinueOnError:   true,
			RowTimeValueCheck: true,
		},
		TierConfigs: []*model.TierConfig{
			{
				Name:                "hotTier",
				SegmentSelectorType: "time",
				SegmentAge:          "3130d",
				StorageType:         "pinot_server",
				ServerTag:           "DefaultTenant_OFFLINE",
			},
		},
		IsDimTable: false,
	}

	tableBytes, err := json.Marshal(table)
	if err != nil {
		t.Errorf("Couldn't marshal table: %v", err)
	}

	res, err := client.CreateTable(tableBytes)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	assert.Equal(t, res.Status, "Table test_OFFLINE successfully added", "Expected response to be Table test_OFFLINE successfully added")
}

func TestUpdateTable(t *testing.T) {
	server := createMockControllerServer()
	client := createPinotClient(server)

	updateTable := model.Table{
		TableName: "test",
		TableType: "OFFLINE",
		SegmentsConfig: model.TableSegmentsConfig{
			TimeColumnName:            "timestamp",
			TimeType:                  "MILLISECONDS",
			Replication:               "1",
			SegmentAssignmentStrategy: "BalanceNumSegmentAssignmentStrategy",
			SegmentPushType:           "APPEND",
			MinimizeDataMovement:      true,
		},
		Tenants: model.TableTenant{
			Broker: "DefaultTenant",
			Server: "DefaultTenant",
		},
		TableIndexConfig: model.TableIndexConfig{
			LoadMode: "MMAP",
		},
		Metadata: &model.TableMetadata{
			CustomConfigs: map[string]string{
				"customKey": "customValue",
			},
		},
		FieldConfigList: []model.FieldConfig{
			{
				Name:         "number",
				EncodingType: "RAW",
				IndexType:    "SORTED",
			},
		},
		IngestionConfig: &model.TableIngestionConfig{
			SegmentTimeValueCheck: true,
			TransformConfigs: []model.TransformConfig{
				{
					ColumnName:        "timestamp",
					TransformFunction: "toEpochHours(millis)",
				},
			},
			ContinueOnError:   true,
			RowTimeValueCheck: true,
		},
		TierConfigs: []*model.TierConfig{
			{
				Name:                "hotTier",
				SegmentSelectorType: "time",
				SegmentAge:          "3130d",
				StorageType:         "pinot_server",
				ServerTag:           "DefaultTenant_OFFLINE",
			},
		},
		IsDimTable: false,
	}

	updateTableBytes, err := json.Marshal(updateTable)
	if err != nil {
		log.Panic(err)
	}

	updateTableResp, err := client.UpdateTable(updateTable.TableName, updateTableBytes)
	if err != nil {
		log.Panic(err)
	}

	assert.Equal(t, updateTableResp.Status, "Table config updated for test_OFFLINE", "Expected response to be Table config updated for test_OFFLINE")
}

func TestDeleteTable(t *testing.T) {
	server := createMockControllerServer()
	client := createPinotClient(server)

	res, err := client.DeleteTable("test")
	if err != nil {
		assert.Equal(t, err.Error(), "client: request failed: status 404\n404 page not found\n", "Expected error to be 404 page not found")
	}

	assert.Equal(t, res.Status, "Tables: [test_OFFLINE] deleted", "Expected response to be Tables: [test_OFFLINE] deleted")
}

func TestGetEmptyTable(t *testing.T) {
	// server := createMockControllerServer()
	// client := createPinotClient(server)

	emptyTable := model.Table{}

	assert.Equal(t, true, emptyTable.IsEmpty(), "Expected table to be empty")

}

func TestGetNonEmptyTable(t *testing.T) {
	// server := createMockControllerServer()
	// client := createPinotClient(server)

	table := model.Table{
		TableName: "test",
	}

	assert.Equal(t, false, table.IsEmpty(), "Expected table to not be empty")

}

func TestEmptyAuthType(t *testing.T) {
	mux := http.NewServeMux()

	mux.HandleFunc(RouteClusterInfo, func(w http.ResponseWriter, r *http.Request) {
		// Check the Authorization header
		authHeader := r.Header.Get("Authorization")
		assert.Equal(t, authHeader, "Basic your_token", "Expected Authorization header to be 'Basic your_token'")

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"clusterName": "PinotCluster"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	// Create a new client
	client := goPinotAPI.NewPinotAPIClient(
		goPinotAPI.AuthType(""),
		goPinotAPI.ControllerUrl(server.URL),
		goPinotAPI.AuthToken("your_token"),
	)

	// Make a request (replace this with an actual API call)
	_, err := client.GetClusterInfo()
	assert.NoError(t, err, "Expected no error from client.Get")
}

func TestBasicAuthType(t *testing.T) {
	mux := http.NewServeMux()

	mux.HandleFunc(RouteClusterInfo, func(w http.ResponseWriter, r *http.Request) {
		// Check the Authorization header
		authHeader := r.Header.Get("Authorization")
		assert.Equal(t, authHeader, "Basic your_token", "Expected Authorization header to be 'Basic your_token'")

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"clusterName": "PinotCluster"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	// Create a new client
	client := goPinotAPI.NewPinotAPIClient(
		goPinotAPI.AuthType("Basic"),
		goPinotAPI.ControllerUrl(server.URL),
		goPinotAPI.AuthToken("your_token"),
	)

	// Make a request (replace this with an actual API call)
	_, err := client.GetClusterInfo()
	assert.NoError(t, err, "Expected no error from client.Get")
}

func TestBearerAuthType(t *testing.T) {
	mux := http.NewServeMux()

	mux.HandleFunc(RouteClusterInfo, func(w http.ResponseWriter, r *http.Request) {
		// Check the Authorization header
		authHeader := r.Header.Get("Authorization")
		assert.Equal(t, authHeader, "Bearer your_token", "Expected Authorization header to be 'Bearer your_token'")

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"clusterName": "PinotCluster"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	// Create a new client
	client := goPinotAPI.NewPinotAPIClient(
		goPinotAPI.AuthType("Bearer"),
		goPinotAPI.ControllerUrl(server.URL),
		goPinotAPI.AuthToken("your_token"),
	)

	// Make a request (replace this with an actual API call)
	_, err := client.GetClusterInfo()
	assert.NoError(t, err, "Expected no error from client.Get")
}

func TestNoAuthType(t *testing.T) {
	mux := http.NewServeMux()

	mux.HandleFunc(RouteClusterInfo, func(w http.ResponseWriter, r *http.Request) {
		// Check the Authorization header
		authHeader := r.Header.Get("Authorization")
		assert.Equal(t, authHeader, "Basic your_token", "Expected Authorization header to be 'Basic your_token'")

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"clusterName": "PinotCluster"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	// Create a new client
	client := goPinotAPI.NewPinotAPIClient(
		goPinotAPI.ControllerUrl(server.URL),
		goPinotAPI.AuthToken("your_token"),
	)

	// Make a request (replace this with an actual API call)
	_, err := client.GetClusterInfo()
	assert.NoError(t, err, "Expected no error from client.Get")
}

func TestMultipleAuthTypes(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Logf("The code panicked with %v", r)
		}
	}()
	mux := http.NewServeMux()

	mux.HandleFunc(RouteClusterInfo, func(w http.ResponseWriter, r *http.Request) {
		// Check the Authorization header
		authHeader := r.Header.Get("Authorization")
		assert.Equal(t, authHeader, "Basic your_token", "Expected Authorization header to be 'Basic your_token'")

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"clusterName": "PinotCluster"}`)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	// Create a new client
	client := goPinotAPI.NewPinotAPIClient(
		goPinotAPI.AuthType("Bearer"),
		goPinotAPI.AuthType("Basic"),
		goPinotAPI.ControllerUrl(server.URL),
		goPinotAPI.AuthToken("your_token"),
	)

	// Make a request (replace this with an actual API call)
	_, err := client.GetClusterInfo()
	assert.Error(t, err, "Expected error from client.Get")
}

func TestPinotControllerAdmin(t *testing.T) {
	server := createMockControllerServer()
	client := createPinotClient(server)

	res, err := client.CheckPinotControllerAdminHealth()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	assert.Equal(t, res.Response, "GOOD", "Expected response to be GOOD")
}

func TestPinotHealth(t *testing.T) {
	server := createMockControllerServer()
	client := createPinotClient(server)

	res, err := client.CheckPinotControllerHealth()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	assert.Equal(t, res.Response, "OK", "Expected response to be OK")
}

func TestGetFieldSpecs(t *testing.T) {
	server := createMockControllerServer()
	client := createPinotClient(server)

	res, err := client.GetSchemaFieldSpecs()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	assert.Equal(t, (res.FieldTypes["METRIC"].AllowedDataTypes["LONG"].NullDefault), float64(0), "Expected METRIC field type to equal 0")
	assert.Equal(t, (res.DataTypes["LONG"].StoredType), "LONG", "Expected LONG data type to equal LONG")
}

func TestGetTableExternalView(t *testing.T) {
	server := createMockControllerServer()
	client := createPinotClient(server)

	res, err := client.GetTableExternalView("test")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	assert.Equal(t, res.Offline["test_OFFLINE_0"]["Server_172.17.0.3_7050"], "ONLINE", "Expected Server_172.17.0.3_7050 to be ONLINE")
}

func TestGetTableIdealState(t *testing.T) {
	server := createMockControllerServer()
	client := createPinotClient(server)

	res, err := client.GetTableIdealState("test")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	assert.Equal(t, res.Offline["test_OFFLINE_0"]["Server_172.17.0.3_7050"], "ONLINE", "Expected Server_172.17.0.3_7050 to be ONLINE")
}

func TestGetTableIndexes(t *testing.T) {
	server := createMockControllerServer()
	client := createPinotClient(server)

	res, err := client.GetTableIndexes("test")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	assert.Equal(t, res.TotalOnlineSegments, 31, "Expected 31 indexes in the response")
	assert.Equal(t, res.ColumnToIndexesCount["Quarter"].ForwardIndex, 31, "Expected ForwardIndex to be set to 31 in the response")
}

func TestGetTableInstances(t *testing.T) {
	server := createMockControllerServer()
	client := createPinotClient(server)

	res, err := client.GetTableInstances("test")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	assert.Equal(t, res.Brokers[0].TableType, "offline", "Expected broker_0 tableType to be offline")
	assert.Equal(t, res.Brokers[0].Instances[0], "Broker_172.17.0.3_8000", "Expected broker_0 instance_0 to be Broker_172.17.0.3_8000")
}

func TestGetTableLiveBrokers(t *testing.T) {
	server := createMockControllerServer()
	client := createPinotClient(server)

	res, err := client.GetAllTableLiveBrokers()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	brokers, ok := (*res)["test_OFFLINE"]
	if !ok {
		fmt.Println("No brokers found for test_OFFLINE")
	} else {
		for _, broker := range brokers {
			fmt.Printf("InstanceName: %s\n", broker.InstanceName)
			assert.Equal(t, broker.InstanceName, "Broker_172.17.0.3_8000", "Expected test broker_0 to be Broker_172.17.0.3_8000")
		}
	}

}

func TestGetTableTestLiveBrokers(t *testing.T) {
	server := createMockControllerServer()
	client := createPinotClient(server)

	res, err := client.GetTableLiveBrokers("test")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	assert.Equal(t, (*res)[0], "Broker_172.17.0.3_8000", "Expected test broker_0 to be Broker_172.17.0.3_8000")

}

func TestGetTableMetadata(t *testing.T) {
	server := createMockControllerServer()
	client := createPinotClient(server)

	res, err := client.GetTableMetadata("test")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	assert.Equal(t, res.TableName, "test_OFFLINE", "Expected table name to be test_OFFLINE")
}

func TestRebuildBrokerResourceFromHelixTags(t *testing.T) {
	server := createMockControllerServer()
	client := createPinotClient(server)

	res, err := client.RebuildBrokerResourceFromHelixTags("test")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	assert.Equal(t, res.Status, "Broker resource is not rebuilt because ideal state is the same for table: test_OFFLINE", "Expected response to be Broker resource is not rebuilt because ideal state is the same for table: test_OFFLINE")
}

func TestGetTableSchema(t *testing.T) {
	server := createMockControllerServer()
	client := createPinotClient(server)

	res, err := client.GetTableSchema("test")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	assert.Equal(t, res.SchemaName, "test", "Expected schema name to be test")
}

func TestGetTableSize(t *testing.T) {
	server := createMockControllerServer()
	client := createPinotClient(server)

	res, err := client.GetTableSize("test")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	assert.Equal(t, res.TableName, "test", "Expected table name to be test")
}

func TestGetTableState(t *testing.T) {
	server := createMockControllerServer()
	client := createPinotClient(server)

	res, err := client.GetTableState("test", "OFFLINE")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	assert.Equal(t, res.State, "enabled", "Expected table state to be enabled")
}

func TestChangeTableState(t *testing.T) {
	server := createMockControllerServer()
	client := createPinotClient(server)

	res, err := client.ChangeTableState("test", "OFFLINE", "enable")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	assert.Equal(t, res.Status, "Request to enable table 'test_OFFLINE' is successful", "Expected Request to enable table 'test_OFFLINE' is successful")
}

func TestGetTableStats(t *testing.T) {
	server := createMockControllerServer()
	client := createPinotClient(server)

	res, err := client.GetTableStats("test")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	assert.Equal(t, (*res)["OFFLINE"].CreationTime, "20240411T224913Z", "Expected table creationtime to be 20240411T224913Z")
}

func TestResetTableSegment(t *testing.T) {
	server := createMockControllerServer()
	client := createPinotClient(server)

	res, err := client.ResetTableSegment("test_OFFLINE", "test_OFFLINE_16071_16071_0")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	assert.Equal(t, res.Status, "Successfully reset segment: test_OFFLINE_16071_16071_0 of table: test_OFFLINE", "Expected response to be Successfully reset segment: test_OFFLINE_16071_16071_0 of table: test_OFFLINE")
}

func TestResetAllTableSegments(t *testing.T) {
	server := createMockControllerServer()
	client := createPinotClient(server)

	res, err := client.ResetTableSegments("test_OFFLINE")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	assert.Equal(t, res.Status, "Successfully reset segments of table: test_OFFLINE", "Expected response to be Successfully reset segments of table: test_OFFLINE")
}
