package container

import (
	"context"
	"fmt"
	"github.com/hamba/avro/v2"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/kafka"
	"github.com/twmb/franz-go/pkg/kadm"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/sr"
	"log"
)

type Kafka struct {
	Container *kafka.KafkaContainer
	Brokers   []string
	TearDown  func()
}

func StartKafkaContainer(ctx context.Context) (*Kafka, error) {

	kafkaContainer, err := kafka.RunContainer(ctx,
		kafka.WithClusterID("test-cluster"),
		testcontainers.WithImage("confluentinc/confluent-local:7.5.0"),
	)
	if err != nil {
		log.Fatalf("failed to start container: %s", err)
	}

	// Clean up the container after
	tearDown := func() {
		if err := kafkaContainer.Terminate(ctx); err != nil {
			log.Panicf("failed to terminate container: %s", err)
		}
	}

	brokers, err := kafkaContainer.Brokers(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get brokers: %s", err)
	}

	return &Kafka{
		Container: kafkaContainer,
		Brokers:   brokers,
		TearDown:  tearDown,
	}, nil

}

func NewSchemaRegistryClient(schemaRegistryUrl string) (*sr.Client, error) {

	client, err := sr.NewClient(sr.URLs(schemaRegistryUrl))
	if err != nil {
		return nil, fmt.Errorf("unable to create schema registry client %w", err)
	}

	return client, nil
}

func NewKafkaClient(brokers []string, consumerGroup string) (*kgo.Client, error) {

	client, err := kgo.NewClient(
		kgo.SeedBrokers(brokers...),
		kgo.ConsumerGroup(consumerGroup),
	)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func ProduceMessage(ctx context.Context, client *kgo.Client, records ...*kgo.Record) error {

	response := client.ProduceSync(ctx, records...)
	if response.FirstErr() != nil {
		return response.FirstErr()
	}

	return nil
}

func ConsumeMessage(ctx context.Context, client *kgo.Client, topic string, maxPollRecords int) ([]*kgo.Record, error) {

	client.AddConsumeTopics(topic)

	fetches := client.PollRecords(ctx, maxPollRecords)
	if fetches.Err() != nil {
		return nil, fetches.Err()
	}

	return fetches.Records(), nil
}

func NewSerde[T any](ctx context.Context, srClient *sr.Client, schemaBytes []byte, subject string) (*sr.Serde, error) {

	subjectSchema, err := srClient.CreateSchema(ctx, subject, sr.Schema{
		Schema: string(schemaBytes),
		Type:   sr.TypeAvro,
	})
	if err != nil {
		return nil, fmt.Errorf("unable to create schema %w", err)
	}

	avroSchema, err := avro.Parse(string(schemaBytes))
	if err != nil {
		return nil, fmt.Errorf("unable to parse schema %w", err)
	}

	var schemaType T
	var serde sr.Serde
	serde.Register(
		subjectSchema.ID,
		schemaType,
		sr.EncodeFn(func(a any) ([]byte, error) {
			return avro.Marshal(avroSchema, a)
		}),
		sr.DecodeFn(func(bytes []byte, a any) error {
			return avro.Unmarshal(avroSchema, bytes, a)
		}),
	)

	return &serde, nil
}

func CreateTopic(ctx context.Context, client *kgo.Client, topicToCreate string) error {

	kadmin := kadm.NewClient(client)

	topicDetails, err := kadmin.ListTopics(ctx)
	if err != nil {
		return err
	}

	if topicDetails.Has(topicToCreate) {
		fmt.Printf("Topic %v already exists\n", topicToCreate)
		return nil
	}

	log.Printf("🚧 Creating topic %v\n", topicToCreate)

	createTopicResponse, err := kadmin.CreateTopic(ctx, 1, 1, nil, topicToCreate)
	if err != nil {
		return fmt.Errorf("unable to create topic %v\n", err)
	}

	log.Printf("✅ Successfully created topic %v\n", createTopicResponse.Topic)
	return nil

}
