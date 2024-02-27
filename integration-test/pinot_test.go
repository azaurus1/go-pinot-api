package integration_test

import (
	"context"
	"encoding/json"
	"fmt"
	"integration-test/container"
	"log"
	"testing"
	"time"

	"github.com/azaurus1/go-pinot-api/model"
	"github.com/stretchr/testify/assert"
	"github.com/twmb/franz-go/pkg/kgo"
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

func TestGetUser(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Minute)
	defer cancel()

	t.Run("Get User", func(t *testing.T) {
		pinot, err := container.RunPinotContainer(ctx)
		assert.NoError(t, err)
		defer pinot.TearDown()

		userResp, err := pinot.GetUsers(ctx)
		if err != nil {
			log.Fatal(err)
		}

		for userName, info := range userResp.Users {
			fmt.Println(userName, info)
		}
		pinot.TearDown()
	})

}

func TestCreateUser(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Minute)
	defer cancel()

	t.Run("Create User", func(t *testing.T) {
		pinot, err := container.RunPinotContainer(ctx)
		assert.NoError(t, err)
		defer pinot.TearDown()

		user := model.User{
			Username:  "testUser",
			Password:  "password",
			Component: "BROKER",
			Role:      "admin",
		}

		userBytes, err := json.Marshal(user)
		if err != nil {
			log.Fatal(err)
		}

		createResp, err := pinot.CreateUser(ctx, userBytes)

		fmt.Println(createResp) // &{User testUser_BROKER has been successfully added!}

		//TODO: Add assertion

		pinot.TearDown()

	})
}
