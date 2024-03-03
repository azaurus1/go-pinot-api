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
