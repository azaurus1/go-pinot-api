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

	// Create User
	createResp, err := client.CreateUser(userBytes)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(createResp.Status)

	// Read Users
	userResp, err := client.GetUsers()
	if err != nil {
		log.Panic(err)
	}

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
