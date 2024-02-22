package main

import (
	"fmt"
	"os"

	goPinotAPI "github.com/azaurus1/go-pinot-api"
)

func main() {
	PINOT_URL := os.Getenv("PINOT_URL")

	client := goPinotAPI.PinotAPIClient{Host: PINOT_URL}

	userBody := []byte(`{
		"liam"
	}`)

	tableBody := []byte(`{
		"testTable"
	}`)

	fmt.Println(client.GetUsers())
	fmt.Println(client.CreateUser(userBody))

	fmt.Println(client.GetTenants())

	fmt.Println(client.GetSchemas())

	fmt.Println(client.GetTables())
	fmt.Println(client.CreateTable(tableBody))
}
