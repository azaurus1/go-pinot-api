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
	RouteUsers                = "/users"
	RouteUser                 = "/users/test"
	RouteInstances            = "/instances"
	RouteInstance             = "/instances/Minion_172.19.0.2_9514"
	RouteClusterInfo          = "/cluster/info"
	RouteClusterConfigs       = "/cluster/configs"
	RouteClusterConfigsDelete = "/cluster/configs/allowParticipantAutoJoin"
	RouteTenants              = "/tenants"
	RouteTenantsInstances     = "/tenants/DefaultTenant"
	RouteTenantsTables        = "/tenants/DefaultTenant/tables"
	RouteTenantsMetadata      = "/tenants/DefaultTenant/metadata"
	RouteSegmentsTest         = "/segments/test"
	RouteSegmentsTestReload   = "/segments/test/reload"
	RouteSegmentTestReload    = "/segments/test/test_1/reload"
	RouteV2Segments           = "/v2/segments"
	RouteSchemas              = "/schemas"
	RouteSchemasTest          = "/schemas/test"
	RouteTables               = "/tables"
	RouteTablesTest           = "/tables/test"
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
		"REALTIME":  {
			"tableName": "test_REALTIME",
			"tableType": "REALTIME",
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
				"serverTag": "DefaultTenant_REALTIME"
			  },
			  {
				"name": "coldTier",
				"segmentSelectorType": "time",
				"segmentAge": "3140d",
				"storageType": "pinot_server",
				"serverTag": "DefaultTenant_REALTIME"
			  }
			],
			"isDimTable": false
		  }
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

	assert.Equal(t, res.REALTIME.TableName, "test_REALTIME", "Expected table name to be test")

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
				NotNull:  false,
			},
			{
				Name:     "hash",
				DataType: "STRING",
				NotNull:  false,
			},
			{
				Name:     "parent_hash",
				DataType: "STRING",
				NotNull:  false,
			},
		},
		MetricFieldSpecs: []model.FieldSpec{
			{
				Name:     "gas_used",
				DataType: "LONG",
				NotNull:  false,
			},
		},
		DateTimeFieldSpecs: []model.FieldSpec{
			{
				Name:        "timestamp",
				DataType:    "LONG",
				NotNull:     false,
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
