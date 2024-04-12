package goPinotAPI

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/azaurus1/go-pinot-api/model"
)

type PinotAPIClient struct {
	pinotControllerUrl *url.URL
	pinotHttp          *pinotHttp
	Host               string
	log                *slog.Logger
}

func NewPinotAPIClient(opts ...Opt) *PinotAPIClient {

	clientCfg, pinotControllerUrl, err := validateOpts(opts...)
	if err != nil {
		log.Panic(err)
	}

	return &PinotAPIClient{
		pinotControllerUrl: pinotControllerUrl,
		pinotHttp: &pinotHttp{
			httpClient:         &http.Client{},
			pinotControllerUrl: pinotControllerUrl,
			httpAuthWriter:     clientCfg.httpAuthWriter,
		},
		Host: pinotControllerUrl.Hostname(),
		log:  clientCfg.logger,
	}
}

func (c *PinotAPIClient) FetchData(endpoint string, result any) error {

	fullURL := prepareRequestURL(c, endpoint)

	request, err := http.NewRequest(http.MethodGet, fullURL.String(), nil)
	if err != nil {
		return fmt.Errorf("client: could not create request: %w", err)
	}

	c.log.Debug(fmt.Sprintf("attempting GET %s", fullURL))

	resp, err := c.pinotHttp.Do(request)
	if err != nil {
		c.logErrorResp(resp)
		return fmt.Errorf("client: could not send request: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// extract body contents to add to error message
		bodyContents, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("client: request failed with status code: %d", resp.StatusCode)
		}
		return fmt.Errorf("client: request failed with status code: %d, body: %s", resp.StatusCode, string(bodyContents))
	}

	bodyContents, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("client: could not read response body: %w", err)
	}

	err = json.NewDecoder(bytes.NewReader(bodyContents)).Decode(result)
	if err != nil {
		c.log.Debug(fmt.Sprintf("unable to decode response from successful request: %s", err))
		return fmt.Errorf("client: could not unmarshal response JSON: %w\n%s", err, string(bodyContents))
	}

	return nil
}

func (c *PinotAPIClient) FetchPlainText(endpoint string, result *model.PlainTextAPIResponse) error {

	fullURL := prepareRequestURL(c, endpoint)

	request, err := http.NewRequest(http.MethodGet, fullURL.String(), nil)
	if err != nil {
		return fmt.Errorf("client: could not create request: %w", err)
	}

	c.log.Debug(fmt.Sprintf("attempting GET %s", fullURL))

	resp, err := c.pinotHttp.Do(request)
	if err != nil {
		c.logErrorResp(resp)
		return fmt.Errorf("client: could not send request: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// extract body contents to add to error message
		bodyContents, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("client: request failed with status code: %d", resp.StatusCode)
		}
		return fmt.Errorf("client: request failed with status code: %d, body: %s", resp.StatusCode, string(bodyContents))
	}

	bodyContents, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("client: could not read response body: %w", err)
	}

	c.log.Debug(fmt.Sprintf("response from successful request: %s", string(bodyContents)))

	result.Response = string(bodyContents)

	return nil

}

func (c *PinotAPIClient) CreateObject(endpoint string, body []byte, result any) error {

	fullURL := prepareRequestURL(c, endpoint)

	req, err := http.NewRequest(http.MethodPost, fullURL.String(), bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("client: could not create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	c.log.Debug(fmt.Sprintf("attempting POST %s", fullURL))

	res, err := c.pinotHttp.Do(req)
	if err != nil {
		c.logErrorResp(res)
		return fmt.Errorf("client: could not send request: %w", err)
	}

	// Check the status code
	if res.StatusCode != http.StatusOK {

		errRespMessage, err := c.extractErrorMessage(res)
		if err != nil {
			return fmt.Errorf("client: could not extract error message: %w", err)
		}

		var errMsg string

		// From client perspective, 409 isnt a failed request
		if res.StatusCode == 409 {
			errMsg = "client: conflict, object exists - "
		} else if res.StatusCode == 403 {
			errMsg = "client: forbidden - "
		} else {
			errMsg = "client: "
		}

		return fmt.Errorf("%srequest failed: status %d\n%s", errMsg, res.StatusCode, errRespMessage)
	}

	err = json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		return fmt.Errorf("client: could not unmarshal JSON: %w", err)
	}

	return nil
}

func (c *PinotAPIClient) CreateFormDataObject(endpoint string, body []byte, result any) error {

	fullURL := c.pinotControllerUrl.JoinPath(endpoint).String()

	req, err := http.NewRequest(http.MethodPost, fullURL, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("client: could not create request: %w", err)
	}

	req.Header.Set("Content-Type", "multipart/form-data")

	res, err := c.pinotHttp.Do(req)
	if err != nil {
		c.logErrorResp(res)
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

	fullURL := prepareRequestURL(c, endpoint)

	c.encodeParams(fullURL, queryParams)

	request, err := http.NewRequest(http.MethodDelete, fullURL.String(), nil)
	if err != nil {
		return fmt.Errorf("client: could not create request: %w", err)
	}

	c.log.Debug(fmt.Sprintf("attempting DELETE %s", fullURL.String()))

	res, err := c.pinotHttp.Do(request)
	if err != nil {
		c.logErrorResp(res)
		return fmt.Errorf("client: could not send request: %w", err)
	}

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusMultipleChoices {

		errRespMessage, err := c.extractErrorMessage(res)
		if err != nil {
			return fmt.Errorf("client: could not extract error message: %w", err)
		}

		var errMsg string
		// From client perspective, 409 isnt a failed request
		switch res.StatusCode {
		case http.StatusNotFound:
			errMsg = "client: object can not be found - "
		default:
			errMsg = "client: "
		}

		return fmt.Errorf("%srequest failed: status %d\n%s", errMsg, res.StatusCode, errRespMessage)

	}

	err = json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		return fmt.Errorf("client: could not unmarshal JSON: %w", err)
	}

	return nil
}

func (c *PinotAPIClient) UpdateObject(endpoint string, queryParams map[string]string, body []byte, result any) error {

	fullURL := c.pinotControllerUrl.JoinPath(endpoint)

	c.encodeParams(fullURL, queryParams)

	req, err := http.NewRequest(http.MethodPut, fullURL.String(), bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("client: could not create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	c.log.Debug(fmt.Sprintf("attempting PUT %s", fullURL.String()))

	res, err := c.pinotHttp.Do(req)
	if err != nil {
		c.logErrorResp(res)
		return fmt.Errorf("client: could not send request: %w", err)
	}

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusMultipleChoices {

		errRespMessage, err := c.extractErrorMessage(res)
		if err != nil {
			return fmt.Errorf("client: could not extract error message: %w", err)
		}

		var errMsg string
		// From client perspective, 409 isnt a failed request
		if res.StatusCode == 404 {
			errMsg = "client: object can not be found - "
		} else {
			errMsg = "client: "
		}

		return fmt.Errorf("%srequest failed: status %d\n%s", errMsg, res.StatusCode, errRespMessage)
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

func (c *PinotAPIClient) CreateTable(body []byte) (*model.CreateTablesResponse, error) {
	result := &model.CreateTablesResponse{}
	err := c.CreateObject("/tables", body, result)
	return result, err
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

func (c *PinotAPIClient) CreateTableFromFile(tableConfigFile string) (*model.CreateTablesResponse, error) {

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

func (c *PinotAPIClient) GetTableExternalView(tableName string) (*model.GetTableExternalViewResponse, error) {
	var result model.GetTableExternalViewResponse
	endpoint := fmt.Sprintf("/tables/%s/externalview", tableName)
	err := c.FetchData(endpoint, &result)
	return &result, err
}

func (c *PinotAPIClient) GetTableIdealState(tableName string) (*model.GetTableIdealStateResponse, error) {
	var result model.GetTableIdealStateResponse
	endpoint := fmt.Sprintf("/tables/%s/idealstate", tableName)
	err := c.FetchData(endpoint, &result)
	return &result, err
}

func (c *PinotAPIClient) GetTableIndexes(tableName string) (*model.GetTableIndexesResponse, error) {
	var result model.GetTableIndexesResponse
	endpoint := fmt.Sprintf("/tables/%s/indexes", tableName)
	err := c.FetchData(endpoint, &result)
	return &result, err
}

func (c *PinotAPIClient) GetTableInstances(tableName string) (*model.GetTableInstancesResponse, error) {
	var result model.GetTableInstancesResponse
	endpoint := fmt.Sprintf("/tables/%s/instances", tableName)
	err := c.FetchData(endpoint, &result)
	return &result, err
}

func (c *PinotAPIClient) GetAllTableLiveBrokers() (*model.GetLiveBrokersResponse, error) {
	var result model.GetLiveBrokersResponse
	err := c.FetchData("/tables/livebrokers", &result)
	return &result, err
}

func (c *PinotAPIClient) GetTableLiveBrokers(tableName string) (*[]string, error) {
	var result []string
	endpoint := fmt.Sprintf("/tables/%s/livebrokers", tableName)
	err := c.FetchData(endpoint, &result)
	return &result, err
}

func (c *PinotAPIClient) GetTableMetadata(tableName string) (*model.GetTableMetadataResponse, error) {
	var result model.GetTableMetadataResponse
	endpoint := fmt.Sprintf("/tables/%s/metadata", tableName)
	err := c.FetchData(endpoint, &result)
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

func (c *PinotAPIClient) CreateSchemaFromBytes(schemaBytes []byte) (*model.CreateSchemaResponse, error) {

	var schema model.Schema

	// Validate first
	json.Unmarshal(schemaBytes, &schema)

	schemaResp, err := c.ValidateSchema(schema)
	if err != nil {
		return nil, fmt.Errorf("unable to validate schema: %w", err)
	}

	if !schemaResp.Ok {
		return nil, fmt.Errorf("schema is invalid: %s", schemaResp.Error)
	}

	var result model.CreateSchemaResponse
	err = c.CreateObject("/schemas?override=false&force=false", schemaBytes, &result)

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
func (c *PinotAPIClient) ValidateSchema(schema model.Schema) (*model.ValidateSchemaResponse, error) {

	schemaBytes, err := schema.AsBytes()
	if err != nil {
		return nil, fmt.Errorf("unable to marshal schema: %w", err)
	}

	fullUrl := c.pinotControllerUrl.JoinPath("schemas", "validate").String()

	req, err := http.NewRequest("POST", fullUrl, bytes.NewBuffer(schemaBytes))
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

		return &model.ValidateSchemaResponse{
			Ok:    false,
			Error: result["error"],
		}, nil
	}

	return &model.ValidateSchemaResponse{Ok: true}, nil
}

func (c *PinotAPIClient) UpdateSchemaFromBytes(schemaBytes []byte) (*model.UserActionResponse, error) {

	var schema model.Schema

	// Validate first
	json.Unmarshal(schemaBytes, &schema)

	schemaResp, err := c.ValidateSchema(schema)
	if err != nil {
		return nil, fmt.Errorf("unable to validate schema: %w", err)
	}

	if !schemaResp.Ok {
		return nil, fmt.Errorf("schema is invalid: %s", schemaResp.Error)
	}

	var result model.UserActionResponse
	err = c.CreateObject("/schemas", schemaBytes, &result)
	return &result, err
}

func (c *PinotAPIClient) UpdateSchema(schema model.Schema) (*model.UserActionResponse, error) {

	var result model.UserActionResponse

	schemaBytes, err := schema.AsBytes()
	if err != nil {
		return nil, fmt.Errorf("unable to marshal schema: %w", err)
	}

	err = c.CreateObject("/schemas", schemaBytes, result) // Should be PUT?
	return &result, err

}

func (c *PinotAPIClient) DeleteSchema(schemaName string) (*model.UserActionResponse, error) {

	getTablesRes, err := c.GetTables()
	if err != nil {
		return nil, fmt.Errorf("unable to get tables names to check: %w", err)
	}

	for _, tableName := range getTablesRes.Tables {
		if tableName == schemaName {
			return nil, fmt.Errorf("can not delete schema %s, it is used by table %s", schemaName, tableName)
		}
	}

	// proceed with deletion
	var result model.UserActionResponse
	err = c.DeleteObject(fmt.Sprintf("/schemas/%s", schemaName), nil, &result)

	return &result, err
}

func (c *PinotAPIClient) GetSchemaFieldSpecs() (*model.GetSchemaFieldSpecsResponse, error) {
	var result model.GetSchemaFieldSpecsResponse
	err := c.FetchData(fmt.Sprintf("/schemas/fieldSpec"), &result)
	return &result, err
}

// Segments
// TODO: Implement Create, Get, GetMetadata, Delete
// func (c *PinotAPIClient) CreateSegment(body []byte) (*model.UserActionResponse, error) {
// 	var result model.UserActionResponse
// 	err := c.CreateFormDataObject("/v2/segments", body, &result)
// 	return &result, err
// }

func (c *PinotAPIClient) GetSegments(tableName string) (model.GetSegmentsResponse, error) {
	var result model.GetSegmentsResponse
	err := c.FetchData(fmt.Sprintf("/segments/%s", tableName), &result)
	return result, err
}

// func (c *PinotAPIClient) GetSegmentMetadata(tableName string, segmentName string) (*model.GetSegmentMetadataResponse, error) {
// 	var result model.GetSegmentMetadataResponse
// 	err := c.FetchData(fmt.Sprintf("/segments/%s/%s/metadata", tableName, segmentName), &result)
// 	return &result, err
// }

// func (c *PinotAPIClient) DeleteSegment(tableName string, segmentName string) (*model.UserActionResponse, error) {
// 	var result model.UserActionResponse
// 	err := c.DeleteObject(fmt.Sprintf("/segments/%s/%s", tableName, segmentName), nil, &result)
// 	return &result, err
// }

func (c *PinotAPIClient) ReloadTableSegments(tableName string) (*model.UserActionResponse, error) {
	var result model.UserActionResponse
	err := c.CreateObject(fmt.Sprintf("/segments/%s/reload", tableName), nil, &result)
	return &result, err
}

func (c *PinotAPIClient) ReloadSegment(tableName string, segmentName string) (*model.UserActionResponse, error) {
	var result model.UserActionResponse
	err := c.CreateObject(fmt.Sprintf("/segments/%s/%s/reload", tableName, segmentName), nil, &result)
	return &result, err
}

// Cluster

func (c *PinotAPIClient) GetClusterInfo() (*model.GetClusterResponse, error) {
	var result model.GetClusterResponse
	err := c.FetchData("/cluster/info", &result)

	return &result, err
}

func (c *PinotAPIClient) GetClusterConfigs() (*model.GetClusterConfigResponse, error) {
	var result model.GetClusterConfigResponse
	err := c.FetchData("/cluster/configs", &result)

	return &result, err
}

func (c *PinotAPIClient) UpdateClusterConfigs(body []byte) (*model.UserActionResponse, error) {
	var result model.UserActionResponse
	err := c.CreateObject("/cluster/configs", body, &result)
	return &result, err
}

func (c *PinotAPIClient) DeleteClusterConfig(configName string) (*model.UserActionResponse, error) {
	var result model.UserActionResponse
	err := c.DeleteObject(fmt.Sprintf("/cluster/configs/%s", configName), nil, &result)
	return &result, err
}

// Tenants

func (c *PinotAPIClient) GetTenants() (*model.GetTenantsResponse, error) {
	var result model.GetTenantsResponse
	err := c.FetchData("/tenants", &result)
	return &result, err
}

func (c *PinotAPIClient) GetTenantInstances(tenantName string) (*model.GetTenantResponse, error) {
	var result model.GetTenantResponse
	err := c.FetchData(fmt.Sprintf("/tenants/%s", tenantName), &result)
	return &result, err

}

func (c *PinotAPIClient) GetTenantTables(tenantName string) (*model.GetTablesResponse, error) {
	var result model.GetTablesResponse
	err := c.FetchData(fmt.Sprintf("/tenants/%s/tables", tenantName), &result)
	return &result, err
}

func (c *PinotAPIClient) GetTenantMetadata(tenantName string) (*model.GetTenantMetadataResponse, error) {
	var result model.GetTenantMetadataResponse
	err := c.FetchData(fmt.Sprintf("/tenants/%s/metadata", tenantName), &result)
	return &result, err
}

func (c *PinotAPIClient) CreateTenant(body []byte) (*model.UserActionResponse, error) {
	var result model.UserActionResponse
	err := c.CreateObject("/tenants", body, &result)
	return &result, err
}

func (c *PinotAPIClient) UpdateTenant(body []byte) (*model.UserActionResponse, error) {
	var result model.UserActionResponse
	err := c.UpdateObject("/tenants", nil, body, &result)
	return &result, err
}

func (c *PinotAPIClient) DeleteTenant(tenantName string, tenantType string) (*model.UserActionResponse, error) {
	var result model.UserActionResponse
	err := c.DeleteObject(fmt.Sprintf("/tenants/%s?type=%s", tenantName, tenantType), nil, &result)
	return &result, err
}

func (c *PinotAPIClient) RebalanceTenant(tenantName string) (*model.UserActionResponse, error) {
	var result model.UserActionResponse
	err := c.CreateObject(fmt.Sprintf("/tenants/%s/rebalance", tenantName), nil, &result)
	return &result, err
}

// Instances
func (c *PinotAPIClient) GetInstances() (*model.GetInstancesResponse, error) {
	var result model.GetInstancesResponse
	err := c.FetchData("/instances", &result)
	return &result, err
}

func (c *PinotAPIClient) GetInstance(instanceName string) (*model.GetInstanceResponse, error) {
	var result model.GetInstanceResponse
	err := c.FetchData(fmt.Sprintf("/instances/%s", instanceName), &result)
	return &result, err
}

func (c *PinotAPIClient) CreateInstance(body []byte) (*model.UserActionResponse, error) {
	var result model.UserActionResponse
	err := c.CreateObject("/instances", body, &result)
	return &result, err
}

func (c *PinotAPIClient) UpdateInstance(instanceName string, body []byte) (*model.UserActionResponse, error) {
	var result model.UserActionResponse
	err := c.UpdateObject(fmt.Sprintf("/instances/%s", instanceName), nil, body, &result)
	return &result, err
}

func (c *PinotAPIClient) DeleteInstance(instanceName string) (*model.UserActionResponse, error) {
	var result model.UserActionResponse
	err := c.DeleteObject(fmt.Sprintf("/instances/%s", instanceName), nil, &result)
	return &result, err
}

func (c *PinotAPIClient) CheckPinotControllerAdminHealth() (*model.PlainTextAPIResponse, error) {
	// Returns text/plain
	var result model.PlainTextAPIResponse
	err := c.FetchPlainText("/pinot-controller/admin", &result)
	return &result, err
}

func (c *PinotAPIClient) CheckPinotControllerHealth() (*model.PlainTextAPIResponse, error) {
	// Returns text/plain
	var result model.PlainTextAPIResponse
	err := c.FetchPlainText("/health", &result)
	return &result, err
}

func (c *PinotAPIClient) extractErrorMessage(resp *http.Response) (string, error) {
	resultBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("unable to decode response from failed request: %s", err)
	}
	return string(resultBytes), nil
}

func (c *PinotAPIClient) logErrorResp(r *http.Response) {

	var responseContent map[string]any

	err := json.NewDecoder(r.Body).Decode(&responseContent)
	if err != nil {
		c.log.Debug(fmt.Sprintf("unable to decode response from failed request: %s", err))
		return
	}

	c.log.Debug(fmt.Sprintf("response from failed request: %s", responseContent))

}

func (c *PinotAPIClient) encodeParams(fullUrl *url.URL, params map[string]string) {
	query := fullUrl.Query()
	for key, value := range params {
		query.Set(key, value)
	}
	fullUrl.RawQuery = query.Encode()
}

func (c *PinotAPIClient) generateQueryParams(queryString string) map[string]string {
	parts := strings.SplitN(queryString, "=", 2)

	// Create a map and add the key-value pair
	m := make(map[string]string)
	if len(parts) == 2 {
		m[parts[0]] = parts[1]
	}
	return m
}

func prepareRequestURL(c *PinotAPIClient, endpoint string) *url.URL {
	pathAndQuery := strings.SplitN(endpoint, "?", 2)
	var path string
	if len(pathAndQuery) > 0 {
		path = pathAndQuery[0]
	}
	var queryString string
	if len(pathAndQuery) > 1 {
		queryString = pathAndQuery[1]
	}

	queryMap := c.generateQueryParams(queryString)
	fullURL := c.pinotControllerUrl.JoinPath(path)
	c.encodeParams(fullURL, queryMap)

	return fullURL
}
