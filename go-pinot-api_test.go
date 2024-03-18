package goPinotAPI_test

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	goPinotAPI "github.com/azaurus1/go-pinot-api"
	model "github.com/azaurus1/go-pinot-api/model"
	"github.com/stretchr/testify/assert"
)

func createMockControllerServer() *httptest.Server {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Basic YWRtaW46YWRtaW4K" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		if r.URL.Path == "/users" && r.Method == "GET" {
			fmt.Fprint(w, `{"users": {"test_BROKER": {"username": "test","password": "$2a$10$3KYPvIy4fBM3CWfdKSc54u/BOn1rPlgb7u4P66s4upqAct30C/q6a","component": "BROKER","role": "ADMIN","usernameWithComponent": "test_BROKER"}}}`)
		}
		if r.URL.Path == "/users/test" && r.URL.Query().Get("component") == "BROKER" && r.Method == "GET" {
			fmt.Fprint(w, `{"test_BROKER": {"username": "test","password": "$2a$10$3KYPvIy4fBM3CWfdKSc54u/BOn1rPlgb7u4P66s4upqAct30C/q6a","component": "BROKER","role": "ADMIN","usernameWithComponent": "test_BROKER"}}`)
		}
		if r.URL.Path == "/users" && r.Method == "POST" {
			fmt.Fprint(w, `{"status": "User testUser_BROKER has been successfully added!"}`)
		}
		if r.URL.Path == "/users/test" && r.URL.Query().Get("component") == "BROKER" && r.Method == "PUT" {
			fmt.Fprint(w, `{"status": "User config update for test_BROKER"}`)
		}
		if r.URL.Path == "/users/test" && r.Method == "DELETE" && r.URL.Query().Get("component") == "" {
			http.Error(w, `{"code": 400,"error": "Name is null"}`, http.StatusBadRequest)
		}
		if r.URL.Path == "/users/test" && r.URL.Query().Get("component") == "BROKER" && r.Method == "DELETE" {
			fmt.Fprint(w, `{"status": "User: test_BROKER has been successfully deleted"}`)
		}
		if r.URL.Path == "/instances" && r.Method == "GET" {
			fmt.Fprint(w, `{"instances": ["Minion_172.19.0.2_9514","Server_172.19.0.7_8098","Broker_cdba1ba98e74_8099","Controller_8684b6757488_9000"]}`)
		}
		if r.URL.Path == "/instances/Minion_172.19.0.2_9514" && r.Method == "GET" {
			fmt.Fprintf(w, `{"instanceName": "Minion_172.19.0.2_9514","hostName": "172.19.0.2","enabled": true,"port": "9514","tags": ["minion_untagged"],"pools": null,"grpcPort": -1,"adminPort": -1,"queryServicePort": -1,"queryMailboxPort": -1,"systemResourceInfo": null}`)
		}
		if r.URL.Path == "/instances" && r.Method == "POST" {
			fmt.Fprint(w, `{"status": "Added instance: Broker_localhost_1234"}`)
		}
	}))

	return server

}

func createPinotClient(server *httptest.Server) *goPinotAPI.PinotAPIClient {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	return goPinotAPI.NewPinotAPIClient(
		goPinotAPI.ControllerUrl(server.URL),
		goPinotAPI.AuthToken("YWRtaW46YWRtaW4K"),
		goPinotAPI.Logger(logger),
	)
}

// TestFetchData
// func TestFetchData(t *testing.T) {

// }

// test FetchData unauthorized

// Test GetUsers
func TestGetUsers(t *testing.T) {

	server := createMockControllerServer()
	client := createPinotClient(server)

	res, err := client.GetUsers()

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// expect 1 user in the response username test
	assert.Equal(t, len(res.Users), 1, "Expected 1 user in the response")
	assert.Equal(t, res.Users["test_BROKER"].Username, "test", "Expected username to be test")

}

// Test GetUser
func TestGetUser(t *testing.T) {

	server := createMockControllerServer()
	client := createPinotClient(server)

	res, err := client.GetUser("test", "BROKER")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	fmt.Println(res.Username)

	// expect username test
	// expect component BROKER
	assert.Equal(t, res.Username, "test", "Expected username to be test")
	assert.Equal(t, res.Component, "BROKER", "Expected component to be BROKER")

}

// Test CreateUser
func TestCreateUser(t *testing.T) {

	server := createMockControllerServer()
	client := createPinotClient(server)

	user := model.User{
		Username:  "testUser",
		Password:  "test",
		Component: "BROKER",
		Role:      "ADMIN",
	}

	userBytes, err := json.Marshal(user)
	if err != nil {
		t.Errorf("Couldn't marshal user: %v", err)
	}

	res, err := client.CreateUser(userBytes)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// expect User testUser_BROKER has been successfully added!
	assert.Equal(t, res.Status, "User testUser_BROKER has been successfully added!", "Expected response to be User testUser_BROKER has been successfully added!")
}

// Test UpdateUser
func TestUpdateUser(t *testing.T) {
	server := createMockControllerServer()
	client := createPinotClient(server)

	user := model.User{
		Username:  "test",
		Password:  "test",
		Component: "BROKER",
		Role:      "USER",
	}

	userBytes, err := json.Marshal(user)
	if err != nil {
		t.Errorf("Couldn't marshal user: %v", err)
	}

	res, err := client.UpdateUser(user.Username, user.Component, true, userBytes)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// expect User test_BROKER has been successfully updated!
	assert.Equal(t, res.Status, "User config update for test_BROKER", "Expected response to be User config update for test_BROKER")
}

// Test DeleteUser
func TestDeleteUser(t *testing.T) {
	server := createMockControllerServer()
	client := createPinotClient(server)

	res, err := client.DeleteUser("test", "BROKER")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// expect User test_BROKER has been successfully deleted!
	assert.Equal(t, res.Status, "User: test_BROKER has been successfully deleted", "Expected response to be User: test_BROKER has been successfully deleted!")

}

func TestDeleteUserNoComponent(t *testing.T) {
	server := createMockControllerServer()
	client := createPinotClient(server)

	_, err := client.DeleteUser("test", "")
	if err != nil {
		assert.Equal(t, err.Error(), "client: request failed: status 400\n{\"code\": 400,\"error\": \"Name is null\"}\n", "Expected error to be Name is null")
	}

}

func TestGetInstances(t *testing.T) {
	server := createMockControllerServer()
	client := createPinotClient(server)

	res, err := client.GetInstances()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	assert.Equal(t, len(res.Instances), 4, "Expected 4 instances in the response")
}

func TestGetInstance(t *testing.T) {
	server := createMockControllerServer()
	client := createPinotClient(server)

	res, err := client.GetInstance("Minion_172.19.0.2_9514")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	assert.Equal(t, res.InstanceName, "Minion_172.19.0.2_9514", "Expected instance name to be Minion_172.19.0.2_9514")
}

func TestCreateInstance(t *testing.T) {
	server := createMockControllerServer()
	client := createPinotClient(server)

	instance := model.Instance{
		Host: "localhost",
		Port: 1234,
		Type: "BROKER",
	}

	instanceBytes, err := json.Marshal(instance)
	if err != nil {
		t.Errorf("Couldn't marshal instance: %v", err)
	}

	res, err := client.CreateInstance(instanceBytes)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	assert.Equal(t, res.Status, "Added instance: Broker_localhost_1234", "Expected response to be Added instance: Broker_localhost_1234")

}
