package container

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type Pinot struct {
	Container testcontainers.Container
	TearDown  func()
}

func StartPinotContainer() (*Pinot, error) {
	ctx := context.Background()
	start := time.Now()

	pinotContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "apachepinot/pinot:latest",
			ExposedPorts: []string{"2123/tcp", "9000/tcp", "8000/tcp", "7050/tcp", "6000/tcp"},
			Cmd:          []string{"QuickStart", "-type", "batch"},
			WaitingFor:   wait.ForLog("You can always go to http://localhost:9000 to play around in the query console").WithStartupTimeout(4 * time.Minute),
		},
		Started: true,
	})
	if err != nil {
		fmt.Errorf("failed to start container: %s", err)
	}

	log.Printf(" ✅ Pinot took %v\n to start", time.Since(start))

	defer pinotContainer.Terminate(ctx)

	tearDown := func() {
		if err := pinotContainer.Terminate(ctx); err != nil {
			fmt.Errorf("failed to terminate container: %s", err)
		}
	}

	return &Pinot{
		Container: pinotContainer,
		TearDown:  tearDown,
	}, nil
}
