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

func main() {

	envPinotUrl := os.Getenv("PINOT_URL")
	if envPinotUrl != "" {
		PinotUrl = envPinotUrl
	}

	client := pinot.NewPinotAPIClient(PinotUrl)

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

	// Get Tables
	tablesResp, err := client.GetTables()
	if err != nil {
		log.Panic(err)
	}

	fmt.Println("Reading Tables:")
	fmt.Println(tablesResp.Tables[0])

	// Get Table
	_, err = client.GetTable("airlineStats")
	if err != nil {
		log.Panic(err)
	}

	// Create Table
	table := pinotModel.Table{
		TableName: "blockHeader",
		TableType: "OFFLINE",
		SegmentsConfig: pinotModel.TableSegmentsConfig{
			TimeColumnName:            "time",
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
				Name:         "blockNumber",
				EncodingType: "RAW",
				IndexType:    "SORTED",
				IndexTypes:   []string{"SORTED"},
				TimestampConfig: pinotModel.TimestampConfig{
					Granulatities: []string{"DAY"},
				},
				Indexes: pinotModel.FieldIndexes{
					Inverted: pinotModel.FiendIndexInverted{
						Enabled: "false",
					},
				},
			},
		},
		IngestionConfig: pinotModel.TableIngestionConfig{
			SegmentTimeValueCheckType: "EPOCH",
			TransformConfigs: []pinotModel.TransformConfig{
				{
					ColumnName:        "blockNumber",
					TransformFunction: "fromEpochDays(DaysSinceEpoch)",
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

	// Update Table
	updateTable := pinotModel.Table{
		TableName: "blockHeader",
		TableType: "OFFLINE",
		SegmentsConfig: pinotModel.TableSegmentsConfig{
			TimeColumnName:            "time",
			TimeType:                  "SECONDS",
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
				Name:         "blockNumber",
				EncodingType: "RAW",
				IndexType:    "SORTED",
				IndexTypes:   []string{"SORTED"},
				TimestampConfig: pinotModel.TimestampConfig{
					Granulatities: []string{"DAY"},
				},
				Indexes: pinotModel.FieldIndexes{
					Inverted: pinotModel.FiendIndexInverted{
						Enabled: "false",
					},
				},
			},
		},
		IngestionConfig: pinotModel.TableIngestionConfig{
			SegmentTimeValueCheckType: "EPOCH",
			TransformConfigs: []pinotModel.TransformConfig{
				{
					ColumnName:        "blockNumber",
					TransformFunction: "fromEpochDays(DaysSinceEpoch)",
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

}

func getSchema() pinotModel.Schema {

	schemaFilePath := "./data-gen/block_header_schema.json"

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
