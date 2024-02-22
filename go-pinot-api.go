package goPinotAPI

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type PinotAPIClient struct {
	Host string
}

type User struct {
	Username              string `json:"username"`
	Password              string `json:"password"`
	Component             string `json:"component"`
	Role                  string `json:"role"`
	UsernameWithComponent string `json:"usernameWithComponent"`
}

type GetUsersResponse struct {
	Users map[string][]User `json:"users"`
}

type CreateUsersResponse struct {
	Status string `json:"status"`
}

type GetTablesResponse struct {
	Tables []string `json:"tables"`
}

type CreateTablesResponse struct {
	UnrecognizedProperties map[string][]string `json:"unrecognizedProperties"`
	Status                 string              `json:"status"`
}

type GetTenantsResponse struct {
	ServerTenants []string `json:"SERVER_TENANTS"`
	BrokerTenants []string `json:"BROKER_TENANTS"`
}

type GetSchemaResponse []string

// generic function
func (c *PinotAPIClient) FetchData(endpoint string, result interface{}) error {
	fullURL := fmt.Sprintf("%s%s", c.Host, endpoint)
	resp, err := http.Get(fullURL)
	if err != nil {
		return fmt.Errorf("client: could not create request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("client: could not read response body: %w", err)
	}

	err = json.Unmarshal(body, result)
	if err != nil {
		return fmt.Errorf("client: could not unmarshal JSON: %w", err)
	}

	return nil
}

func (c *PinotAPIClient) CreateObject(endpoint string, body []byte, result interface{}) error {
	fullURL := fmt.Sprintf("%s%s", c.Host, endpoint)

	req, err := http.NewRequest("POST", fullURL, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("client: could not create request: %w", err)
	}
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("client: could not send request: %w", err)
	}

	json.NewDecoder(res.Body).Decode(result)

	return nil
}

// users
func (c *PinotAPIClient) GetUsers() (*GetUsersResponse, error) {
	var result GetUsersResponse
	err := c.FetchData("/users", &result)
	return &result, err
}

func (c *PinotAPIClient) CreateUser(body []byte) (*CreateUsersResponse, error) {
	var result CreateUsersResponse
	err := c.CreateObject("/users", body, result)
	return &result, err
}

// tables
func (c *PinotAPIClient) GetTables() (*GetTablesResponse, error) {
	var result GetTablesResponse
	err := c.FetchData("/tables", &result)
	return &result, err
}

func (c *PinotAPIClient) CreateTable(body []byte) (*CreateUsersResponse, error) {
	var result CreateUsersResponse
	err := c.CreateObject("/tables", body, result)
	return &result, err
}

// tenants
func (c *PinotAPIClient) GetTenants() (*GetTenantsResponse, error) {
	var result GetTenantsResponse
	err := c.FetchData("/tenants", &result)
	return &result, err
}

// schemas
func (c *PinotAPIClient) GetSchemas() (*GetSchemaResponse, error) {
	var result GetSchemaResponse
	err := c.FetchData("/schemas", &result)
	return &result, err
}
