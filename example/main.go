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
		Component: "Broker",
		Role:      "admin",
	}

	userBytes, err := json.Marshal(user)
	if err != nil {
		log.Panic(err)
	}

	_, err = client.CreateUser(userBytes)
	if err != nil {
		log.Panic(err)
	}

	userResp, err := client.GetUsers()
	if err != nil {
		log.Panic(err)
	}

	for userName, info := range userResp.Users {
		fmt.Println(userName, info)
	}

}
