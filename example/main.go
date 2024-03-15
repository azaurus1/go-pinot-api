package main

import (
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"os"

	pinot "github.com/azaurus1/go-pinot-api"
	pinotModel "github.com/azaurus1/go-pinot-api/model"
)

const LocalPinotUrl = "http://localhost:9000"
const LocalPinotAuthToken = "YWRtaW46dmVyeXNlY3JldA" // Default Admin password=verysecret  admin:verysecret (b64 encoded)

func getSchema() pinotModel.Schema {

	schemaFilePath := "./example/data-gen/block_header_schema.json"

	f, err := os.Open(schemaFilePath)
	if err != nil {
		log.Panic(err)
	}

	defer f.Close()

	var schema pinotModel.Schema
	err = json.NewDecoder(f).Decode(&schema)
	if err != nil {
		log.Panic(err)
	}

	return schema
}

func main() {

	pinotUrl := getOrDefault(LocalPinotUrl, "PINOT_URL", "PINOT_CONTROLLER_URL")
	authToken := getOrDefault(LocalPinotAuthToken, "PINOT_AUTH", "PINOT_AUTH_TOKEN")

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	client := pinot.NewPinotAPIClient(
		pinot.ControllerUrl(pinotUrl),
		pinot.AuthToken(authToken),
		pinot.Logger(logger))

	// demoSchemaFunctionality(client)
	// demoTableFunctionality(client)
	// demoUserFunctionality(client)
	// demoSegmentFunctionality(client)
	// demoClusterFunctionality(client)
	demoTenantFunctionality(client)
}

func demoTableFunctionality(client *pinot.PinotAPIClient) {
	// Create Table
	fmt.Println("Creating Table:")

	table := pinotModel.Table{
		TableName: "ethereum_mainnet_block_headers",
		TableType: "OFFLINE",
		SegmentsConfig: pinotModel.TableSegmentsConfig{
			TimeColumnName:            "timestamp",
			TimeType:                  "MILLISECONDS",
			Replication:               "1",
			SegmentAssignmentStrategy: "BalanceNumSegmentAssignmentStrategy",
			SegmentPushType:           "APPEND",
			MinimizeDataMovement:      true,
		},
		Tenants: pinotModel.TableTenant{
			Broker: "DefaultTenant",
			Server: "DefaultTenant",
		},
		TableIndexConfig: pinotModel.TableIndexConfig{
			LoadMode: "MMAP",
		},
		Metadata: pinotModel.TableMetadata{
			CustomConfigs: map[string]string{
				"customKey": "customValue",
			},
		},
		FieldConfigList: []pinotModel.FieldConfig{
			{
				Name:         "number",
				EncodingType: "RAW",
				IndexType:    "SORTED",
			},
		},
		IngestionConfig: pinotModel.TableIngestionConfig{
			SegmentTimeValueCheckType: "EPOCH",
			TransformConfigs: []pinotModel.TransformConfig{
				{
					ColumnName:        "timestamp",
					TransformFunction: "fromEpochMilliseconds(timestamp)",
				},
			},
			ContinueOnError:   true,
			RowTimeValueCheck: true,
		},
		TierConfigs: []pinotModel.TierConfig{
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
		log.Panic(err)
	}

	fmt.Println(string(tableBytes))

	createTableResp, err := client.CreateTable(tableBytes)
	if err != nil {
		log.Panic(err)
	}

	fmt.Println(createTableResp.Status)

	// Get Tables
	tablesResp, err := client.GetTables()
	if err != nil {
		log.Panic(err)
	}

	fmt.Println("Reading Tables:")

	// Get Table
	_, err = client.GetTable(tablesResp.Tables[0])
	if err != nil {
		log.Panic(err)
	}

	// Update Table
	updateTable := pinotModel.Table{
		TableName: "ethereum_mainnet_block_headers",
		TableType: "OFFLINE",
		SegmentsConfig: pinotModel.TableSegmentsConfig{
			TimeColumnName:            "timestamp",
			TimeType:                  "MILLISECONDS",
			Replication:               "1",
			SegmentAssignmentStrategy: "BalanceNumSegmentAssignmentStrategy",
			SegmentPushType:           "APPEND",
			MinimizeDataMovement:      true,
		},
		Tenants: pinotModel.TableTenant{
			Broker: "DefaultTenant",
			Server: "DefaultTenant",
		},
		TableIndexConfig: pinotModel.TableIndexConfig{
			LoadMode: "MMAP",
		},
		Metadata: pinotModel.TableMetadata{
			CustomConfigs: map[string]string{
				"customKey": "customValue",
			},
		},
		FieldConfigList: []pinotModel.FieldConfig{
			{
				Name:         "number",
				EncodingType: "RAW",
				IndexType:    "SORTED",
			},
		},
		IngestionConfig: pinotModel.TableIngestionConfig{
			SegmentTimeValueCheckType: "EPOCH",
			TransformConfigs: []pinotModel.TransformConfig{
				{
					ColumnName:        "timestamp",
					TransformFunction: "fromEpochMilliseconds(timestamp)",
				},
			},
			ContinueOnError:   true,
			RowTimeValueCheck: true,
		},
		TierConfigs: []pinotModel.TierConfig{
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

	fmt.Println(updateTableResp.Status)

	// Delete Table
	//deleteTableResp, err := client.DeleteTable(table.TableName)
	//if err != nil {
	//	log.Panic(err)
	//}
	//
	//fmt.Println(deleteTableResp.Status)
}

func demoSchemaFunctionality(client *pinot.PinotAPIClient) {

	schema := getSchema()

	// Create Schema will validate the schema first anyway
	validateResp, err := client.ValidateSchema(schema)
	if err != nil {
		log.Panic(err)
	}

	if !validateResp.Ok {
		log.Panic(validateResp.Error)
	}

	_, err = client.CreateSchema(schema)
	if err != nil {
		log.Panic(err)
	}

	currentSchemas, err := client.GetSchemas()
	if err != nil {
		log.Panic(err)
	}

	currentSchemas.ForEachSchema(func(schemaName string) {

		schemaResp, err := client.GetSchema(schemaName)
		if err != nil {
			log.Panic(err)
		}

		fmt.Println("Reading Schema:")
		fmt.Println(schemaResp)

	})

}

func demoUserFunctionality(client *pinot.PinotAPIClient) {

	user := pinotModel.User{
		Username:  "liam1",
		Password:  "password",
		Component: "BROKER",
		Role:      "admin",
	}

	userBytes, err := json.Marshal(user)
	if err != nil {
		log.Panic(err)
	}

	updateUser := pinotModel.User{
		Username:  "liam1",
		Password:  "password",
		Component: "BROKER",
		Role:      "user",
	}

	updateUserBytes, err := json.Marshal(user)
	if err != nil {
		log.Panic(err)
	}

	fmt.Println("Creating User:")

	// Create User
	createResp, err := client.CreateUser(userBytes)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(createResp.Status)

	// Read User
	getUserResp, err := client.GetUser(user.Username, user.Component)
	if err != nil {
		log.Panic(err)
	}

	fmt.Println("Reading User:")
	fmt.Println(getUserResp.UsernameWithComponent)

	// Read Users
	userResp, err := client.GetUsers()
	if err != nil {
		log.Panic(err)
	}

	fmt.Println("Reading Users:")
	for userName, info := range userResp.Users {
		fmt.Println(userName, info)
	}

	// Update User
	updateResp, err := client.UpdateUser(updateUser.Username, updateUser.Component, false, updateUserBytes)
	if err != nil {
		log.Panic(err)
	}
	fmt.Println(updateResp.Status)

	// Delete User
	delResp, err := client.DeleteUser(user.Username, user.Component)
	if err != nil {
		log.Panic(err)
	}

	fmt.Println(delResp.Status)
}

func demoSegmentFunctionality(client *pinot.PinotAPIClient) {

	// Create Segment
	// fmt.Println("Creating Segment:")
	// segment := pinotModel.Segment{
	// 	ContentDisposition: pinotModel.ContentDisposition{
	// 		Type:       "type",
	// 		Parameters: map[string]string{"param": "value"},
	// 		FileName:   "fileName",
	// 		Size:       100,
	// 	},
	// 	Entity: map[string]interface{}{
	// 		"key": "value",
	// 	},
	// 	Headers: map[string][]string{
	// 		"header": {"value"},
	// 	},
	// 	MediaType: pinotModel.MediaType{
	// 		Type:            "type",
	// 		Subtype:         "subtype",
	// 		Parameters:      map[string]string{"param": "value"},
	// 		WildcardType:    true,
	// 		WildcardSubtype: true,
	// 	},
	// 	MessageBodyWorkers: map[string]interface{}{
	// 		"key": "value",
	// 	},
	// 	Parent: pinotModel.Parent{
	// 		ContentDisposition: pinotModel.ContentDisposition{
	// 			Type:       "type",
	// 			Parameters: map[string]string{"param": "value"},
	// 			FileName:   "fileName",
	// 			Size:       100,
	// 		},
	// 		Entity: map[string]interface{}{
	// 			"key": "value",
	// 		},
	// 		Headers: map[string][]string{
	// 			"header": {"value"},
	// 		},
	// 		MediaType: pinotModel.MediaType{
	// 			Type:            "type",
	// 			Subtype:         "subtype",
	// 			Parameters:      map[string]string{"param": "value"},
	// 			WildcardType:    true,
	// 			WildcardSubtype: true,
	// 		},
	// 		MessageBodyWorkers: map[string]interface{}{
	// 			"key": "value",
	// 		},
	// 		Parent: "parent",
	// 		Providers: map[string]interface{}{
	// 			"key": "value",
	// 		},
	// 		BodyParts: []pinotModel.BodyPart{
	// 			{
	// 				ContentDisposition: pinotModel.ContentDisposition{
	// 					Type:       "type",
	// 					Parameters: map[string]string{"param": "value"},
	// 					FileName:   "fileName",
	// 					Size:       100,
	// 				},
	// 				Entity: map[string]interface{}{
	// 					"key": "value",
	// 				},
	// 				Headers: map[string][]string{
	// 					"header": {"value"},
	// 				},
	// 				MediaType: pinotModel.MediaType{
	// 					Type:            "type",
	// 					Subtype:         "subtype",
	// 					Parameters:      map[string]string{"param": "value"},
	// 					WildcardType:    true,
	// 					WildcardSubtype: true,
	// 				},
	// 				MessageBodyWorkers: map[string]interface{}{
	// 					"key": "value",
	// 				},
	// 				Parent: "parent",
	// 				Providers: map[string]interface{}{
	// 					"key": "value",
	// 				},
	// 				ParameterizedHeaders: map[string][]pinotModel.ParameterizedHeader{
	// 					"key": {
	// 						{
	// 							Value: "value",
	// 						},
	// 					},
	// 				},
	// 			},
	// 		},
	// 		ParameterizedHeaders: map[string][]pinotModel.ParameterizedHeader{
	// 			"key": {
	// 				{
	// 					Value: "value",
	// 				},
	// 			},
	// 		},
	// 	},
	// 	Providers: map[string]interface{}{
	// 		"key": "value",
	// 	},
	// 	ParameterizedHeaders: map[string][]pinotModel.ParameterizedHeader{
	// 		"key": {
	// 			{
	// 				Value: "value",
	// 			},
	// 		},
	// 	},
	// }

	// segmentBytes, err := json.Marshal(segment)
	// if err != nil {
	// 	log.Panic(err)
	// }

	// createSegmentResp, err := client.CreateSegment(segmentBytes)
	// if err != nil {
	// 	log.Panic(err)
	// }

	// fmt.Println(createSegmentResp.Status)

	// Get Segments
	// segmentsResp, err := client.GetSegments("githubComplexTypeEvents")
	// if err != nil {
	// 	log.Panic(err)
	// }
	// fmt.Println("Reading Segments:")
	// if len(segmentsResp) > 0 {
	// 	if len(segmentsResp[0].Offline) > 0 {
	// 		fmt.Println("OFFLINE Data:", segmentsResp[0].Offline)
	// 	} else {
	// 		fmt.Println("No OFFLINE data available.")
	// 	}
	// 	if len(segmentsResp[0].Realtime) > 0 {
	// 		fmt.Println("REALTIME Data:", segmentsResp[0].Realtime)
	// 	} else {
	// 		fmt.Println("No REALTIME data available.")
	// 	}
	// }

	// // Get Segment Metadata
	// segmentMetadataResp, err := client.GetSegmentMetadata("githubComplexTypeEvents", segmentsResp[0].Offline[0])
	// if err != nil {
	// 	log.Panic(err)
	// }

	// fmt.Println("Reading Segment Metadata:")
	// fmt.Println(segmentMetadataResp.SegmentStartTime)

	// Delete Segment
	//deleteSegmentResp, err := client.DeleteSegment("airlineStats", segmentsResp.OFFLINE[0])
	//if err != nil {
	//	log.Panic(err)
	//}

	// Reload All Segments
	// reloadTableSegmentsResp, err := client.ReloadTableSegments("githubComplexTypeEvents")
	// if err != nil {
	// 	log.Panic(err)
	// }

	// fmt.Println(reloadTableSegmentsResp.Status)

	// // Reload Segment
	// reloadSegmentResp, err := client.ReloadSegment("githubComplexTypeEvents", segmentsResp[0].Offline[0])
	// if err != nil {
	// 	log.Panic(err)
	// }

	// fmt.Println(reloadSegmentResp.Status)

}

func demoClusterFunctionality(client *pinot.PinotAPIClient) {

	// Get Cluster
	clusterInfoResp, err := client.GetClusterInfo()
	if err != nil {
		log.Panic(err)
	}

	fmt.Println("Reading Cluster Info:")
	fmt.Println(clusterInfoResp.ClusterName)

	// Get Cluster Config
	clusterConfigResp, err := client.GetClusterConfigs()
	if err != nil {
		log.Panic(err)
	}

	fmt.Println("Reading Cluster Config:")
	fmt.Println(clusterConfigResp.AllowParticipantAutoJoin)

	// Update Cluster Config
	updateClusterConfig := pinotModel.ClusterConfig{
		AllowParticipantAutoJoin:            "true",
		EnableCaseInsensitive:               "true",
		DefaultHyperlogLogLog2m:             "14",
		PinotBrokerEnableQueryLimitOverride: "true",
	}

	updateClusterConfigBytes, err := json.Marshal(updateClusterConfig)
	if err != nil {
		log.Panic(err)
	}

	updateClusterConfigResp, err := client.UpdateClusterConfigs(updateClusterConfigBytes)
	if err != nil {
		log.Panic(err)
	}

	fmt.Println(updateClusterConfigResp.Status)

	// Delete Cluster Config
	// This deletes the cluster config - Theres no Create cluster config operation .....
	// deleteClusterConfigResp, err := client.DeleteClusterConfig(clusterInfoResp.ClusterName)
	// if err != nil {
	// 	log.Panic(err)
	// }

	// fmt.Println(deleteClusterConfigResp.Status)

}

func demoTenantFunctionality(client *pinot.PinotAPIClient) {

	// Create Tenant
	tenant := pinotModel.Tenant{
		TenantName: "test",
		TenantRole: "BROKER",
	}

	tenantBytes, err := json.Marshal(tenant)
	if err != nil {
		log.Panic(err)
	}

	createTenantResp, err := client.CreateTenant(tenantBytes)
	if err != nil {
		log.Panic(err)
	}

	fmt.Println(createTenantResp.Status)

	// Get Tenants
	tenantsResp, err := client.GetTenants()
	if err != nil {
		log.Panic(err)
	}

	fmt.Println("Reading Broker Tenants:")
	for _, tenant := range tenantsResp.BrokerTenants {
		fmt.Println(tenant)
	}

	fmt.Println("Reading Server Tenants:")
	for _, tenant := range tenantsResp.ServerTenants {
		fmt.Println(tenant)
	}

	// Get Tenant
	getTenantResp, err := client.GetTenantMetadata("test")
	if err != nil {
		log.Panic(err)
	}

	fmt.Println("Reading Tenant:")
	fmt.Println(getTenantResp.TenantName)

	// Update Tenant
	updateTenant := pinotModel.Tenant{
		TenantName: "test",
		TenantRole: "SERVER",
	}

	updateTenantBytes, err := json.Marshal(updateTenant)
	if err != nil {
		log.Panic(err)
	}

	updateTenantResp, err := client.UpdateTenant(updateTenantBytes)
	if err != nil {
		log.Panic(err)
	}

	fmt.Println(updateTenantResp.Status)

	// Delete Tenant
	deleteTenantResp, err := client.DeleteTenant("test", "SERVER")
	if err != nil {
		log.Panic(err)
	}

	fmt.Println(deleteTenantResp.Status)

}

func getOrDefault(defaultOption string, envKeys ...string) string {

	for _, envKey := range envKeys {
		if envVal := os.Getenv(envKey); envVal != "" {
			return envVal
		}
	}

	return defaultOption

}
