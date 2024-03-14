package container

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	goPinotAPI "github.com/azaurus1/go-pinot-api"
	"github.com/azaurus1/go-pinot-api/model"
	"github.com/google/uuid"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/network"
	"github.com/testcontainers/testcontainers-go/wait"
)

type Pinot struct {
	Container testcontainers.Container
	TearDown  func()
	URI       string
}

func randomShortUUID() string {
	return uuid.New().String()[:8]
}

func randomisePinotURI(uri string) string {
	return fmt.Sprintf("%s-%s", uri, randomShortUUID())
}

func RunPinotContainer(ctx context.Context) (*Pinot, error) {

	zkURI := randomisePinotURI("pinot-zk")
	pinotURI := randomisePinotURI("pinot-controller")

	absPath, err := filepath.Abs(filepath.Join(".", "testdata", "pinot-controller.conf"))
	if err != nil {
		return nil, fmt.Errorf("failed to add data: %s", err)
	}
	newNetwork, err := network.New(ctx, network.WithCheckDuplicate())
	if err != nil {
		return nil, fmt.Errorf("failed to create network: %s", err)
	}

	networkName := newNetwork.Name

	_, err = testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Networks: []string{networkName},
			NetworkAliases: map[string][]string{
				networkName: {"pinot-zk"},
			},
			Name:         zkURI,
			Image:        "apachepinot/pinot:latest",
			ExposedPorts: []string{"2181/tcp"},
			Cmd:          []string{"StartZookeeper"},
			WaitingFor:   wait.ForLog("Start zookeeper at localhost:2181 in thread main").WithStartupTimeout(4 * time.Minute),
		},
		Started: true,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to start container: %s", err)
	}

	fmt.Println("ZK URI: ", zkURI)

	pinotContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Networks: []string{networkName},
			NetworkAliases: map[string][]string{
				networkName: {"pinot-controller"}, // this does work btw
			},
			Name:         pinotURI,
			Image:        "apachepinot/pinot:latest",
			ExposedPorts: []string{"2123/tcp", "9000/tcp", "8000/tcp", "7050/tcp", "6000/tcp"},
			Files: []testcontainers.ContainerFile{
				{
					HostFilePath:      absPath,
					ContainerFilePath: "/config/pinot-controller.conf",
					FileMode:          0o700,
				},
			},
			Cmd:        []string{"StartController", "-configFileName", "/config/pinot-controller.conf"}, //"StartController"  "-configFileName", "/config/pinot-controller.conf"
			WaitingFor: wait.ForLog("INFO [StartServiceManagerCommand] [main] Started Pinot [CONTROLLER] instance").WithStartupTimeout(4 * time.Minute),
		},
		Started: true,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to start container: %s", err)
	}

	tearDown := func() {
		if err := pinotContainer.Terminate(ctx); err != nil {
			log.Panicf("failed to terminate container: %s", err)
		}
	}

	pinotControllerMappedPort, err := pinotContainer.MappedPort(ctx, "9000")
	if err != nil {
		return nil, fmt.Errorf("failed to get mapped port: %s", err)
	}

	pinotContainerHost, err := pinotContainer.Host(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get container host: %s", err)
	}

	pinotControllerURI := fmt.Sprintf("%s:%v", pinotContainerHost, pinotControllerMappedPort.Port())

	return &Pinot{
		Container: pinotContainer,
		TearDown:  tearDown,
		URI:       pinotControllerURI,
	}, nil
}

func (p *Pinot) CreateUser(_ context.Context, userBytes []byte) (*model.UserActionResponse, error) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	client := goPinotAPI.NewPinotAPIClient(
		goPinotAPI.ControllerUrl(p.URI),
		goPinotAPI.AuthToken("YWRtaW46YWRtaW4="),
		goPinotAPI.Logger(logger),
	)

	userCreationResponse, err := client.CreateUser(userBytes)
	if err != nil {
		log.Fatal(err)
	}

	return userCreationResponse, nil

}

func (p *Pinot) GetUsers(_ context.Context) (*model.GetUsersResponse, error) {
	client := goPinotAPI.NewPinotAPIClient(
		goPinotAPI.ControllerUrl(p.URI),
		goPinotAPI.AuthToken("YWRtaW46YWRtaW4="),
	)

	userResp, err := client.GetUsers()
	if err != nil {
		log.Fatal(err)
	}

	return userResp, nil
}

func (p *Pinot) GetSegment(tableName string) (*model.GetSegmentsResponse, error) {
	client := goPinotAPI.NewPinotAPIClient(
		goPinotAPI.ControllerUrl(p.URI),
		goPinotAPI.AuthToken("YWRtaW46YWRtaW4="),
	)

	segments, err := client.GetSegments(tableName)
	if err != nil {
		log.Fatal(err)
	}

	return &segments, nil
}

func (p *Pinot) ReloadSegment(tableName string, segmentName string) (*model.UserActionResponse, error) {
	client := goPinotAPI.NewPinotAPIClient(
		goPinotAPI.ControllerUrl(p.URI),
		goPinotAPI.AuthToken("YWRtaW46YWRtaW4="),
	)

	reloadSegmentResponse, err := client.ReloadSegment(tableName, segmentName)
	if err != nil {
		log.Fatal(err)
	}

	return reloadSegmentResponse, nil

}

func (p *Pinot) ReloadTableSegments(tableName string) (*model.UserActionResponse, error) {
	client := goPinotAPI.NewPinotAPIClient(
		goPinotAPI.ControllerUrl(p.URI),
		goPinotAPI.AuthToken("YWRtaW46YWRtaW4="),
	)

	reloadTableSegmentsResponse, err := client.ReloadTableSegments(tableName)
	if err != nil {
		log.Fatal(err)
	}

	return reloadTableSegmentsResponse, nil
}
