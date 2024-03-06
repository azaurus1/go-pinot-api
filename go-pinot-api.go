package goPinotAPI

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"

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

func NewPinotAPIClient(pinotController string, pinotAuthToken string) *PinotAPIClient {

	pinotUrl, err := url.Parse(pinotController)
	if err != nil {
		log.Panic(err)
	}

	// handle authenticated requests
	// pinotAuthToken := os.Getenv("PINOT_AUTH_TOKEN")
	httpAuthWriterFunc := func(req *http.Request) {
		if pinotAuthToken != "" {
			req.Header.Set("Authorization", fmt.Sprintf("Basic %s", pinotAuthToken))
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

type CreateTablesResponse struct {
	UnrecognizedProperties map[string][]string `json:"unrecognizedProperties"`
	Status                 string              `json:"status"`
}

type GetTenantsResponse struct {
	ServerTenants []string `json:"SERVER_TENANTS"`
	BrokerTenants []string `json:"BROKER_TENANTS"`
}

type ValidateSchemaResponse struct {
	Ok    bool
	Error string
}

func (c *PinotAPIClient) FetchData(endpoint string, result any) error {

	fullURL := fullUrl(c.pinotControllerUrl, endpoint)

	request, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return fmt.Errorf("client: could not create request: %w", err)
	}

	c.pinotHttp.httpAuthWriter(request)

	resp, err := c.pinotHttp.httpClient.Do(request)
	if err != nil {
		return fmt.Errorf("client: could not send request: %w", err)
	}

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(result)
	if err != nil {
		return fmt.Errorf("client: could not unmarshal JSON: %w", err)
	}

	return nil
}

func (c *PinotAPIClient) CreateObject(endpoint string, body []byte, result any) error {

	fullURL := fullUrl(c.pinotControllerUrl, endpoint)

	req, err := http.NewRequest("POST", fullURL, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("client: could not create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	c.pinotHttp.httpAuthWriter(req)

	res, err := c.pinotHttp.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("client: could not send request: %w", err)
	}

	// Check the status code
	if res.StatusCode != http.StatusOK {
		var errMsg string
		// From client perspective, 409 isnt a failed request
		if res.StatusCode == 409 {
			errMsg = "client: conflict, object exists - "
		} else if res.StatusCode == 403 {
			errMsg = "client: forbidden - "
		} else {
			errMsg = "client: "
		}
		return fmt.Errorf("%srequest failed with status code: %d", errMsg, res.StatusCode)
	}

	err = json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		return fmt.Errorf("client: could not unmarshal JSON: %w", err)
	}

	return nil
}

func (c *PinotAPIClient) DeleteObject(endpoint string, queryParams map[string]string, result any) error {
	fullURL := fullUrl(c.pinotControllerUrl, endpoint)

	parsedURL, err := url.Parse(fullURL)
	if err != nil {
		return fmt.Errorf("client: could not parse URL: %w", err)
	}

	if len(queryParams) > 0 {
		query := parsedURL.Query()
		for key, value := range queryParams {
			query.Set(key, value)
		}
		parsedURL.RawQuery = query.Encode()
	}

	req, err := http.NewRequest("DELETE", parsedURL.String(), nil)
	if err != nil {
		return fmt.Errorf("client: could not create request: %w", err)
	}

	c.pinotHttp.httpAuthWriter(req)

	res, err := c.pinotHttp.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("client: could not send request: %w", err)
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		var errMsg string
		// From client perspective, 409 isnt a failed request
		if res.StatusCode == 404 {
			errMsg = "client: object can not be found - "
		} else {
			errMsg = "client: "
		}
		return fmt.Errorf("%srequest failed with status code: %d", errMsg, res.StatusCode)
	}

	err = json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		return fmt.Errorf("client: could not unmarshal JSON: %w", err)
	}

	return nil
}

func (c *PinotAPIClient) UpdateObject(endpoint string, queryParams map[string]string, body []byte, result any) error {
	fullURL := fullUrl(c.pinotControllerUrl, endpoint)

	parsedURL, err := url.Parse(fullURL)
	if err != nil {
		return fmt.Errorf("client: could not parse URL: %w", err)
	}

	if len(queryParams) > 0 {
		query := parsedURL.Query()
		for key, value := range queryParams {
			query.Set(key, value)
		}
		parsedURL.RawQuery = query.Encode()
	}

	fmt.Println(parsedURL.String())

	req, err := http.NewRequest("PUT", parsedURL.String(), bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("client: could not create request: %w", err)
	}

	c.pinotHttp.httpAuthWriter(req)
	req.Header.Set("Content-Type", "application/json")

	res, err := c.pinotHttp.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("client: could not send request: %w", err)
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		var errMsg string
		// From client perspective, 409 isnt a failed request
		if res.StatusCode == 404 {
			errMsg = "client: object can not be found - "
		} else {
			errMsg = "client: "
		}
		return fmt.Errorf("%srequest failed with status code: %d", errMsg, res.StatusCode)
	}

	err = json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		return fmt.Errorf("client: could not unmarshal JSON: %w", err)
	}

	return nil

}

func (c *PinotAPIClient) GetUsers() (*model.GetUsersResponse, error) {
	var result model.GetUsersResponse
	err := c.FetchData("/users", &result)
	return &result, err
}

func (c *PinotAPIClient) GetUser(username string, component string) (*model.User, error) {
	var result map[string]model.User
	var resultUser model.User

	endpoint := fmt.Sprintf("/users/%s?component=%s", username, component)
	err := c.FetchData(endpoint, &result)

	usernameWithComponent := fmt.Sprintf("%s_%s", username, component)

	resultUser = result[usernameWithComponent]
	return &resultUser, err
}

func (c *PinotAPIClient) CreateUser(body []byte) (*model.UserActionResponse, error) {
	var result model.UserActionResponse
	err := c.CreateObject("/users", body, &result)
	return &result, err
}

func (c *PinotAPIClient) DeleteUser(username string, component string) (*model.UserActionResponse, error) {
	deletionQueryParams := make(map[string]string)
	deletionQueryParams["component"] = component

	endpoint := fmt.Sprintf("/users/%s", username)

	var result model.UserActionResponse
	err := c.DeleteObject(endpoint, deletionQueryParams, &result)
	return &result, err
}

func (c *PinotAPIClient) UpdateUser(username string, component string, passwordChanged bool, body []byte) (*model.UserActionResponse, error) {
	updateQueryParams := make(map[string]string)
	updateQueryParams["component"] = component
	updateQueryParams["passwordChanged"] = strconv.FormatBool(passwordChanged)

	var result model.UserActionResponse
	endpoint := fmt.Sprintf("/users/%s", username)

	err := c.UpdateObject(endpoint, updateQueryParams, body, &result)
	return &result, err
}

func (c *PinotAPIClient) GetTables() (*model.GetTablesResponse, error) {
	var result model.GetTablesResponse

	err := c.FetchData("/tables", &result)
	return &result, err
}

func (c *PinotAPIClient) GetTable(tableName string) (*model.GetTableResponse, error) {
	var result model.GetTableResponse
	endpoint := fmt.Sprintf("/tables/%s", tableName)
	err := c.FetchData(endpoint, &result)
	return &result, err
}

// TODO: ValidateTable (?)

// func (c *PinotAPIClient) ValidateTable(body []byte) (*model.UserActionResponse, error) {

// }

func (c *PinotAPIClient) CreateTable(body []byte) (*model.UserActionResponse, error) {
	var result model.UserActionResponse
	err := c.CreateObject("/tables", body, result)
	return &result, err
}

func (c *PinotAPIClient) UpdateTable(tableName string, body []byte) (*model.UserActionResponse, error) {
	var result model.UserActionResponse
	endpoint := fmt.Sprintf("/tables/%s", tableName)
	err := c.UpdateObject(endpoint, nil, body, &result)
	return &result, err
}

func (c *PinotAPIClient) DeleteTable(tableName string) (*model.UserActionResponse, error) {
	var result model.UserActionResponse
	endpoint := fmt.Sprintf("/tables/%s", tableName)
	err := c.DeleteObject(endpoint, nil, &result)
	return &result, err
}

func (c *PinotAPIClient) CreateTableFromFile(tableConfigFile string) (*model.UserActionResponse, error) {

	f, err := os.Open(tableConfigFile)
	if err != nil {
		return nil, fmt.Errorf("unable to open table config file: %w", err)
	}

	defer f.Close()

	var tableConfig model.Table
	err = json.NewDecoder(f).Decode(&tableConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal table config: %w", err)
	}

	tableConfigBytes, err := json.Marshal(tableConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal table config: %w", err)
	}

	return c.CreateTable(tableConfigBytes)
}

func (c *PinotAPIClient) GetTenants() (*GetTenantsResponse, error) {
	var result GetTenantsResponse
	err := c.FetchData("/tenants", &result)
	return &result, err
}

// GetSchemas returns a list of schemas
func (c *PinotAPIClient) GetSchemas() (*model.GetSchemaResponse, error) {
	var result model.GetSchemaResponse
	err := c.FetchData("/schemas", &result)
	return &result, err
}

// GetSchema returns a schema
func (c *PinotAPIClient) GetSchema(schemaName string) (*model.Schema, error) {
	var result model.Schema
	err := c.FetchData(fmt.Sprintf("/schemas/%s", schemaName), &result)
	return &result, err

}

// CreateSchema creates a new schema. if it already exists, it will nothing will happen
func (c *PinotAPIClient) CreateSchema(schema model.Schema) (*model.UserActionResponse, error) {

	// validate schema first
	schemaResp, err := c.ValidateSchema(schema)
	if err != nil {
		return nil, fmt.Errorf("unable to validate schema: %w", err)
	}

	if !schemaResp.Ok {
		return nil, fmt.Errorf("schema is invalid: %s", schemaResp.Error)
	}

	var result model.UserActionResponse

	schemaBytes, err := schema.AsBytes()
	if err != nil {
		return nil, fmt.Errorf("unable to marshal schema: %w", err)
	}

	err = c.CreateObject("/schemas", schemaBytes, result)
	return &result, err
}

// CreateSchemaFromFile creates a new schema from a file and uses CreateSchema
func (c *PinotAPIClient) CreateSchemaFromFile(schemaFilePath string) (*model.UserActionResponse, error) {

	f, err := os.Open(schemaFilePath)
	if err != nil {
		return nil, fmt.Errorf("unable to open schema file: %w", err)
	}

	defer f.Close()

	var schema model.Schema
	err = json.NewDecoder(f).Decode(&schema)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal schema: %w", err)
	}

	return c.CreateSchema(schema)

}

// ValidateSchema validates a schema
func (c *PinotAPIClient) ValidateSchema(schema model.Schema) (*ValidateSchemaResponse, error) {

	schemaBytes, err := schema.AsBytes()
	if err != nil {
		return nil, fmt.Errorf("unable to marshal schema: %w", err)
	}

	req, err := http.NewRequest("POST", fullUrl(c.pinotControllerUrl, "/schemas/validate"), bytes.NewBuffer(schemaBytes))
	if err != nil {
		return nil, fmt.Errorf("client: could not create request: %w", err)
	}

	res, err := c.pinotHttp.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("client: could not send request: %w", err)
	}

	if res.StatusCode == http.StatusInternalServerError {
		return nil, fmt.Errorf("client: internal server error")
	}

	// Invalid schema in body
	if res.StatusCode == http.StatusBadRequest {

		var result map[string]string

		err = json.NewDecoder(res.Body).Decode(&result)
		if err != nil {
			return nil, fmt.Errorf("client: could not unmarshal JSON: %w", err)
		}

		return &ValidateSchemaResponse{
			Ok:    false,
			Error: result["error"],
		}, nil
	}

	return &ValidateSchemaResponse{Ok: true}, nil
}

func (c *PinotAPIClient) UpdateSchema(schema model.Schema) (*model.UserActionResponse, error) {

	var result model.UserActionResponse

	schemaBytes, err := schema.AsBytes()
	if err != nil {
		return nil, fmt.Errorf("unable to marshal schema: %w", err)
	}

	err = c.CreateObject("/schemas", schemaBytes, result)
	return &result, err

}

func (c *PinotAPIClient) DeleteSchema(schemaName string) (*model.UserActionResponse, error) {
	var result model.UserActionResponse
	endpoint := fmt.Sprintf("/schemas/%s", schemaName)
	err := c.DeleteObject(endpoint, nil, &result)
	return &result, err
}

func fullUrl(url *url.URL, path string) string {
	return fmt.Sprintf("http://%s:%s%s", url.Hostname(), url.Port(), path)
}
