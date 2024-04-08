package main

import (
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"os"

	pinot "github.com/azaurus1/go-pinot-api"
	"github.com/azaurus1/go-pinot-api/model"
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
	// demoSchemaFromBytesFunctionality(client)
	// demoTableFunctionality(client)
	demoGetTableFunctionality(client)
	// demoUserFunctionality(client)
	// demoSegmentFunctionality(client)
	// demoClusterFunctionality(client)
	// demoTenantFunctionality(client)
	// demoInstanceFunctionality(client)

}

func boolPtr(b bool) *bool {
	return &b
}

func demoTableFunctionality(client *pinot.PinotAPIClient) {

	// Create Offline schema
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
				NotNull:          boolPtr(false),
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

	createSchemaResp, err := client.CreateSchemaFromBytes(schemaBytes)
	if err != nil {
		log.Panic(err)
	}

	fmt.Println(createSchemaResp.Status)

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
		Metadata: &pinotModel.TableMetadata{
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
		IngestionConfig: &pinotModel.TableIngestionConfig{
			SegmentTimeValueCheck: true,
			TransformConfigs: []pinotModel.TransformConfig{
				{
					ColumnName:        "timestamp",
					TransformFunction: "toEpochHours(millis)",
				},
			},
			ContinueOnError:   true,
			RowTimeValueCheck: true,
		},
		TierConfigs: []*pinotModel.TierConfig{
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

	// Get Table
	for _, info := range tablesResp.Tables {
		if info == table.TableName {
			fmt.Println("Reading Table:", info)
		}
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
		Metadata: &pinotModel.TableMetadata{
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
		IngestionConfig: &pinotModel.TableIngestionConfig{
			SegmentTimeValueCheck: true,
			TransformConfigs: []pinotModel.TransformConfig{
				{
					ColumnName:        "timestamp",
					TransformFunction: "toEpochHours(millis)",
				},
			},
			ContinueOnError:   true,
			RowTimeValueCheck: true,
		},
		TierConfigs: []*pinotModel.TierConfig{
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
	deleteTableResp, err := client.DeleteTable(table.TableName)
	if err != nil {
		log.Panic(err)
	}

	fmt.Println(deleteTableResp.Status)

	// Delete Schema
	deleteSchemaResp, err := client.DeleteSchema(schema.SchemaName)
	if err != nil {
		log.Panic(err)
	}

	fmt.Println(deleteSchemaResp.Status)
}

func demoGetTableFunctionality(client *pinot.PinotAPIClient) {

	tableResp, err := client.GetTable("realtime_ethereum_mainnet_block_headers")
	if err != nil {
		log.Panic(err)
	}

	fmt.Println("Reading Table:")
	fmt.Println(tableResp.REALTIME.TableName)

}

func demoSchemaFromBytesFunctionality(client *pinot.PinotAPIClient) {

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
				NotNull:          boolPtr(false),
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

	// validate
	validateResp, err := client.ValidateSchema(schema)
	if err != nil {
		log.Panic(err)
	}

	fmt.Println(validateResp)

	createResp, err := client.CreateSchemaFromBytes(schemaBytes)
	if err != nil {
		log.Panic(err)
	}

	fmt.Println(createResp.Status)

	// Get Schema
	getSchemaResp, err := client.GetSchema(schema.SchemaName)
	if err != nil {
		log.Panic(err)
	}

	fmt.Println("Reading Schema:")
	fmt.Println(getSchemaResp)

	// Update Schema
	updateSchema := model.Schema{
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
				NotNull:          boolPtr(false),
				SingleValueField: boolPtr(false),
			},
			{
				Name:     "test",
				DataType: "STRING",
				NotNull:  boolPtr(false),
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

	fmt.Println("Updating Schema:")
	updateSchemaBytes, err := json.Marshal(updateSchema)

	updateResp, err := client.UpdateSchemaFromBytes(updateSchemaBytes)
	if err != nil {
		log.Panic(err)
	}

	fmt.Println(updateResp.Status)

	// delete schema

	deleteResp, err := client.DeleteSchema(schema.SchemaName)
	if err != nil {
		log.Panic(err)
	}

	fmt.Println(deleteResp.Status)

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

	// Delete Schema
	deleteResp, err := client.DeleteSchema(schema.SchemaName)
	if err != nil {
		log.Panic(err)
	}

	fmt.Println(deleteResp.Status)

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

func demoInstanceFunctionality(client *pinot.PinotAPIClient) {

	// Get Instances
	instancesResp, err := client.GetInstances()
	if err != nil {
		log.Panic(err)
	}

	fmt.Println("Reading Instances:")
	for _, instance := range instancesResp.Instances {
		fmt.Println(instance)
	}

	// Create Instance
	instance := pinotModel.Instance{
		Host: "localhost",
		Port: 9000,
		Type: "CONTROLLER",
	}

	instanceBytes, err := json.Marshal(instance)
	if err != nil {
		log.Panic(err)
	}

	createInstanceResp, err := client.CreateInstance(instanceBytes)
	if err != nil {
		log.Panic(err)
	}

	fmt.Println(createInstanceResp.Status)

	// Get Instance
	getInstanceResp, err := client.GetInstance("Controller_localhost_9000") // this is the name of the instance we created
	if err != nil {
		log.Panic(err)
	}

	fmt.Println("Reading Instance:")
	fmt.Println(getInstanceResp.InstanceName)

	// Update Instance
	updateInstance := pinotModel.Instance{
		Host: "localhost",
		Port: 9000,
		Type: "BROKER",
	}

	updateInstanceBytes, err := json.Marshal(updateInstance)
	if err != nil {
		log.Panic(err)
	}

	updateInstanceResp, err := client.UpdateInstance("Controller_localhost_9000", updateInstanceBytes)
	if err != nil {
		log.Panic(err)
	}

	fmt.Println(updateInstanceResp.Status)

	// Delete Instance
	deleteInstanceResp, err := client.DeleteInstance("Controller_localhost_9000")
	if err != nil {
		log.Panic(err)
	}

	fmt.Println(deleteInstanceResp.Status)
}

func getOrDefault(defaultOption string, envKeys ...string) string {

	for _, envKey := range envKeys {
		if envVal := os.Getenv(envKey); envVal != "" {
			return envVal
		}
	}

	return defaultOption

}
