package main

import (
	"context"
	"fmt"
	"github.com/hamba/avro/v2"
	"github.com/twmb/franz-go/pkg/kadm"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/sr"
	"golang.org/x/crypto/sha3"
	"log"
	"math/rand"
	"os"
	"time"
)

const (
	BootstrapServers  = "localhost:29092"
	SchemaRegistryUrl = "http://localhost:8081"
	topic             = "ethereum.mainnet.blocks"
	BatchSize         = 10000
)

type BlockHeader struct {
	Number     int64  `json:"number" avro:"number"`
	Hash       string `json:"hash" avro:"hash"`
	ParentHash string `json:"parent_hash" avro:"parent_hash"`
	GasUsed    int64  `json:"gas_used" avro:"gas_used"`
	Timestamp  int64  `json:"timestamp" avro:"timestamp"`
}

func main() {

	kafkaClient, err := CreateKafkaClient()
	if err != nil {
		panic(err)
	}

	err = CreateTopic(kafkaClient, topic)
	if err != nil {
		log.Fatalf("unable to create topic %v\n", err)
	}

	schemaRegistryClient, err := NewSchemaRegistryClient()
	if err != nil {
		panic(err)
	}

	// register the schema and Create Serializer/Deserializer
	serde, err := NewSerde[BlockHeader](schemaRegistryClient, "block_header.avsc", "ethereum.mainnet.blocks-value")
	if err != nil {
		panic(err)
	}

	currentBlockNum := int64(10000)
	currentTimestamp := time.Now()

	var batch []*kgo.Record
	for i := 0; i < 1000000; i++ {

		timestamp := currentTimestamp.Add(time.Duration(30) * time.Second)
		block := GenerateBlock(currentBlockNum, timestamp)

		// monotonic increasing block number and increasing timestamp
		currentBlockNum = block.Number
		currentTimestamp = timestamp

		if i%BatchSize == 0 && len(batch) > 0 {

			response := kafkaClient.ProduceSync(context.Background(), batch...)
			if response.FirstErr() != nil {
				panic(response.FirstErr())
			}

			fmt.Printf("Produced %d blocks\n", i)
			// prob better way to do this
			batch = []*kgo.Record{}
		}

		batch = append(batch, &kgo.Record{
			Topic: topic,
			Key:   []byte(fmt.Sprintf("%d", block.Number)),
			Value: serde.MustEncode(block),
		})

	}

}

func CreateKafkaClient() (*kgo.Client, error) {

	kafkaClient, err := kgo.NewClient(kgo.SeedBrokers(BootstrapServers))
	if err != nil {
		return nil, fmt.Errorf("unable to create kafka client %w", err)
	}

	return kafkaClient, nil
}

func NewSchemaRegistryClient() (*sr.Client, error) {

	schemaRegistryClient, err := sr.NewClient(sr.URLs(SchemaRegistryUrl))
	if err != nil {
		return nil, fmt.Errorf("unable to create schema registry client %w", err)
	}

	return schemaRegistryClient, nil
}

func NewSerde[T any](srClient *sr.Client, schemaFilePath string, subject string) (*sr.Serde, error) {

	schemaTextBytes, err := os.ReadFile(schemaFilePath)
	if err != nil {
		return nil, fmt.Errorf("unable to read schema file %w", err)
	}

	subjectSchema, err := srClient.CreateSchema(context.Background(), subject, sr.Schema{
		Schema: string(schemaTextBytes),
		Type:   sr.TypeAvro,
	})
	if err != nil {
		return nil, fmt.Errorf("unable to create schema %w", err)
	}

	avroSchema, err := avro.Parse(string(schemaTextBytes))
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

func GenerateBlock(blockNumber int64, timestamp time.Time) BlockHeader {

	blockHash := sha3.New256()
	blockHash.Write([]byte(fmt.Sprintf("%d", blockNumber)))

	parentHash := sha3.New256()
	parentHash.Write([]byte(fmt.Sprintf("%d", blockNumber-1)))

	return BlockHeader{
		Number:     blockNumber + 1,
		Hash:       fmt.Sprintf("%x", blockHash.Sum(nil)),
		ParentHash: fmt.Sprintf("%x", parentHash.Sum(nil)),
		GasUsed:    rand.Int63n(10000000),
		Timestamp:  timestamp.Unix(),
	}
}

func CreateTopic(kafkaClient *kgo.Client, topicToCreate string) error {

	ctx := context.Background()
	kadmin := kadm.NewClient(kafkaClient)

	topicDetails, err := kadmin.ListTopics(ctx)
	if err != nil {
		return err
	}

	if topicDetails.Has(topicToCreate) {
		fmt.Printf("Topic %v already exists\n", topicToCreate)
		return nil
	}

	fmt.Printf("Creating topic %v\n", topicToCreate)

	createTopicResponse, err := kadmin.CreateTopic(ctx, 1, 1, nil, topicToCreate)
	if err != nil {
		return fmt.Errorf("unable to create topic %v\n", err)
	}

	fmt.Printf("Successfully created topic %v\n", createTopicResponse.Topic)
	return nil

}
