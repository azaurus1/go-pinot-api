package goPinotAPI

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/azaurus1/go-pinot-api/model"
)

type pinotHttp struct {
	httpClient         *http.Client
	pinotControllerUrl *url.URL
	httpAuthWriter     httpAuthWriter
}

type httpAuthWriter func(*http.Request)

type PinotAPIClient struct {
	pinotControllerUrl *url.URL
	pinotHttp          *pinotHttp
	Host               string
}

func NewPinotAPIClient(pinotController string) *PinotAPIClient {

	pinotUrl, err := url.Parse(pinotController)
	if err != nil {
		log.Panic(err)
	}

	// handle authenticated requests
	pinotAuthToken := os.Getenv("PINOT_AUTH_TOKEN")
	httpAuthWriterFunc := func(req *http.Request) {
		if pinotAuthToken != "" {
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", pinotAuthToken))
		}
	}

	return &PinotAPIClient{
		pinotControllerUrl: pinotUrl,
		pinotHttp: &pinotHttp{
			httpClient:         &http.Client{},
			pinotControllerUrl: pinotUrl,
			httpAuthWriter:     httpAuthWriterFunc,
		},
		Host: pinotController,
	}
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
func (c *PinotAPIClient) FetchData(endpoint string, result any) error {

	fullURL := fullUrl(c.pinotControllerUrl, endpoint)

	resp, err := http.Get(fullURL)
	if err != nil {
		return fmt.Errorf("client: could not create request: %w", err)
	}

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(result)
	if err != nil {
		return fmt.Errorf("client: could not unmarshal JSON: %w", err)
	}

	return nil
}

func (c *PinotAPIClient) CreateObject(endpoint string, body []byte, result interface{}) error {

	fullURL := fullUrl(c.pinotControllerUrl, endpoint)

	req, err := http.NewRequest("POST", fullURL, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("client: could not create request: %w", err)
	}

	res, err := c.pinotHttp.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("client: could not send request: %w", err)
	}

	err = json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		return fmt.Errorf("client: could not unmarshal JSON: %w", err)
	}

	return nil
}

// users
func (c *PinotAPIClient) GetUsers() (*model.GetUsersResponse, error) {
	var result model.GetUsersResponse
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

func fullUrl(url *url.URL, path string) string {
	return fmt.Sprintf("http://%s:%s%s", url.Hostname(), url.Port(), path)
}
