package integration_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/twmb/franz-go/pkg/kgo"
	"integration-test/container"
	"testing"
	"time"
)

func TestUser(t *testing.T) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	t.Run("Send Message to Kafka", func(t *testing.T) {

		redpandaInfo, err := container.StartRedPandaContainer()
		assert.NoError(t, err)

		kafkaClient, err := container.NewKafkaClient(redpandaInfo.Brokers, t.Name())
		assert.NoError(t, err)

		defer kafkaClient.Close()

		topicName := "test.topic"

		err = container.CreateTopic(ctx, kafkaClient, topicName)
		assert.NoError(t, err)

		message := &kgo.Record{
			Context: ctx,
			Topic:   topicName,
			Key:     []byte("hello-key"),
			Value:   []byte("hello world-value"),
		}

		err = container.ProduceMessage(ctx, kafkaClient, message)
		assert.NoError(t, err)

		consumedRecords, err := container.ConsumeMessage(ctx, kafkaClient, topicName, 1)
		assert.NoError(t, err)

		for _, record := range consumedRecords {
			assert.Equal(t, "hello-key", string(record.Key))
			assert.Equal(t, "hello world-value", string(record.Value))
		}

	})

}
