package integrationtest

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"testing"
	"time"

	goPinotAPI "github.com/azaurus1/go-pinot-api"
	"github.com/azaurus1/go-pinot-api/model"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestCreateUser(t *testing.T) {
	// Context for the test
	ctx := context.Background()

	// Define the Pinot container request
	pinotReq := testcontainers.ContainerRequest{
		Image: "apachepinot/pinot:latest", // Specify the Pinot image
		ExposedPorts: []string{
			"2123/tcp", // Assuming these ports are necessary for your application
			"9000/tcp", // Controller
			"8000/tcp", // Broker
			"7050/tcp", // Server
			"6000/tcp", // Minion
		},
		Cmd:        []string{"QuickStart", "-type", "batch"},
		Files: []ContainerFile(
			{
				HostFilePath: absPath,
				ContainerFilePath: "/controller.conf",
				FileMode: 0o700,
			}
		)
		WaitingFor: wait.ForLog("You can always go to http://localhost:9000 to play around in the query console").WithStartupTimeout(4 * time.Minute),
	}

	// Start the Pinot container
	pinotContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: pinotReq,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("Failed to start Pinot container: %s", err)
	}

	// Clean up after the test
	defer pinotContainer.Terminate(ctx)

	// Use the first exposed port for communication, assuming it's the controller
	pinotPort, err := pinotContainer.MappedPort(ctx, "9000")
	if err != nil {
		t.Fatalf("Failed to get mapped port: %s", err)
	}

	pinotHost := fmt.Sprintf("http://localhost:%s", pinotPort.Port())

	client := goPinotAPI.NewPinotAPIClient(pinotHost)

	// Test logic to create a user
	user := model.User{
		Username:  "testUser",
		Password:  "password",
		Component: "Broker",
		Role:      "admin",
	}

	userBytes, err := json.Marshal(user)
	if err != nil {
		log.Fatal(err)
	}

	_, err = client.CreateUser(userBytes)
	if err != nil {
		log.Fatal(err)
	}

	// Fetch and log the users
	userResp, err := client.GetUsers()
	if err != nil {
		log.Fatal(err)
	}

	for userName, info := range userResp.Users {
		if userName == user.Username {
			t.Errorf("Expected matching username, got non-matching")
		}
		if info.Password == user.Password {
			t.Errorf("Expected matching password, got non-matching")
		}
		if info.Component == user.Component {
			t.Errorf("Expected matching component, got non-matching")
		}
		if info.Role == user.Role {
			t.Errorf("Expected matching role, got non-matching")
		}
	}

}
