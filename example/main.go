package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	pinot "github.com/azaurus1/go-pinot-api"
	pinotModel "github.com/azaurus1/go-pinot-api/model"
)

var PinotUrl = "http://localhost:9000"
var PinotAuth = "YWRtaW46dmVyeXNlY3JldA" // Default Admin password=verysecret  admin:verysecret (b64 encoded)

func main() {

	envPinotUrl := os.Getenv("PINOT_URL")
	if envPinotUrl != "" {
		PinotUrl = envPinotUrl
	}

	envPinotAuth := os.Getenv("PINOT_AUTH")
	if envPinotAuth != "" {
		PinotAuth = envPinotAuth
	}

	client := pinot.NewPinotAPIClient(PinotUrl, PinotAuth)

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

	// Schema
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
	deleteTableResp, err := client.DeleteTable(table.TableName)
	if err != nil {
		log.Panic(err)
	}

	fmt.Println(deleteTableResp.Status)

}

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
