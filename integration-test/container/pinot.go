package container

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type Pinot struct {
	Container testcontainers.Container
	TearDown  func()
}

func RunPinotContainer(ctx context.Context) (*Pinot, error) {

	absPath, err := filepath.Abs(filepath.Join(".", "testdata", "pinot-controller.conf"))
	if err != nil {
		return nil, fmt.Errorf("failed to add data: %s", err)
	}

	pinotZkContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "apachepinot/pinot:latest",
			ExposedPorts: []string{"2181/tcp"},
			Cmd:          []string{"StartZookeeper"}, //TODO: Change to run on all interfaces
			WaitingFor:   wait.ForLog("Start zookeeper at localhost:2181 in thread main").WithStartupTimeout(4 * time.Minute),
		},
		Started: true,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to start container: %s", err)
	}

	defer pinotZkContainer.Terminate(ctx)

	ip, _ := pinotZkContainer.Host(ctx)
	port, _ := pinotZkContainer.MappedPort(ctx, "2181")

	zkHostUrl := fmt.Sprintf("%s:%s", ip, port.Port())

	pinotContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "apachepinot/pinot:latest",
			ExposedPorts: []string{"2123/tcp", "9000/tcp", "8000/tcp", "7050/tcp", "6000/tcp"},
			Files: []testcontainers.ContainerFile{
				{
					HostFilePath:      absPath,
					ContainerFilePath: "/config/pinot-controller.conf",
					FileMode:          0o700,
				},
			},
			Cmd:        []string{"StartController", "-zkAddress", zkHostUrl}, //"StartController"  "-configFileName", "/config/pinot-controller.conf"
			WaitingFor: wait.ForLog("You can always go to http://localhost:9000 to play around in the query console").WithStartupTimeout(4 * time.Minute),
		},
		Started: true,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to start container: %s", err)
	}

	defer pinotContainer.Terminate(ctx)

	tearDown := func() {
		if err := pinotContainer.Terminate(ctx); err != nil {
			log.Panicf("failed to terminate container: %s", err)
		}
	}

	return &Pinot{
		Container: pinotContainer,
		TearDown:  tearDown,
	}, nil
}
