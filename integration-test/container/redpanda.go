package container

import (
	"context"
	"fmt"
	"github.com/testcontainers/testcontainers-go/modules/redpanda"
	"log"
	"time"
)

type RedPanda struct {
	Container         *redpanda.Container
	Brokers           []string
	SchemaRegistryUrl string
	TearDown          func()
}

func StartRedPandaContainer() (*RedPanda, error) {

	ctx := context.Background()
	start := time.Now()

	redpandaContainer, err := redpanda.RunContainer(ctx,
		redpanda.WithAutoCreateTopics(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to start container: %s", err)
	}

	log.Printf(" ✅ Redpanda took %v\n to start", time.Since(start))

	tearDown := func() {
		if err := redpandaContainer.Terminate(ctx); err != nil {
			log.Panicf("failed to terminate container: %s", err)
		}
	}

	brokers, err := redpandaContainer.ContainerIPs(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get brokers: %s", err)
	}

	schemaRegistryUrl, err := redpandaContainer.SchemaRegistryAddress(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get schema registry url: %s", err)
	}

	return &RedPanda{
		Container:         redpandaContainer,
		Brokers:           brokers,
		SchemaRegistryUrl: schemaRegistryUrl,
		TearDown:          tearDown,
	}, nil

}
