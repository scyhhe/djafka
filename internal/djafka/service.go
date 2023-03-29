package djafka

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type Connection struct {
	Name            string `json:"name"`
	BootstrapServer string `json:"bootstrapServer"`
}

type Config struct {
	Connections []Connection `json:"connections"`
}

func (c *Config) FindConnection(name string) (Connection, error) {
	for _, conn := range c.Connections {
		if conn.Name == name {
			return conn, nil
		}
	}

	return Connection{}, fmt.Errorf("Failed to find connection for name '%s'.", name)
}

func ReadConfig() (*Config, error) {
	file, err := os.Open("config.json")
	if err != nil {
		return nil, fmt.Errorf("Failed to read config file: %w", err)
	}
	defer file.Close()

	config := Config{}
	if err := json.NewDecoder(file).Decode(&config); err != nil {
		return nil, fmt.Errorf("Failed to decode config file: %w", err)
	}

	return &config, nil
}

type Consumer struct {
	GroupId         string
	ConsumerId      string
	State           string
	TopicPartitions []ConsumerTopicPartition
}

type ConsumerTopicPartition struct {
	TopicName string
	Offset    int64
	Partition int32
}

type Service struct {
	client   *kafka.AdminClient
	consumer *kafka.Consumer
	// producer *kafka.Producer
}

func NewService(conn Connection) (*Service, error) {
	client, err := kafka.NewAdminClient(&kafka.ConfigMap{
		"bootstrap.servers": "localhost",
	})

	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost",
		"group.id":          "testis",
	})

	// producer, err := kafka.NewProducer(&kafka.ConfigMap{
	// 	"bootstrap.servers": "localhost",
	// })

	return &Service{client, consumer}, err
}

func (s *Service) Close() {
	s.client.Close()
}

func (s *Service) ListTopics() ([]string, error) {
	metaData, err := s.client.GetMetadata(nil, true, 5000)
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch meta data: %w", err)
	}

	topics := []string{}
	for key := range metaData.Topics {
		topics = append(topics, key)
	}

	return topics, nil
}

func (s *Service) CreateTopic(name string, partitions int, replicationFactor int) (string, error) {
	topicSpec := kafka.TopicSpecification{Topic: name, NumPartitions: partitions, ReplicationFactor: replicationFactor}
	fmt.Println("spec", topicSpec)
	res, err := s.client.CreateTopics(context.Background(), []kafka.TopicSpecification{topicSpec})

	if err != nil {
		return err.Error(), fmt.Errorf("Failed to create new topic: %w", err)
	}
	for _, r := range res {
		if r.Error.Code() != kafka.ErrNoError {
			return r.Error.String(), fmt.Errorf("Failed to create new topic: %w", r.Error)
		}
	}
	return topicSpec.Topic, nil

}

func (s *Service) ListConsumerGroups() ([]string, error) {
	consumerGroups, err := s.client.ListConsumerGroups(context.Background())
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch meta data: %w", err)
	}

	groupIds := []string{}
	for _, group := range consumerGroups.Valid {
		groupIds = append(groupIds, group.GroupID)
	}

	return groupIds, nil
}

func (s *Service) ListConsumers(groupIds []string) ([]Consumer, error) {
	consumerGroups, err := s.client.DescribeConsumerGroups(context.Background(), groupIds)
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch meta data: %w", err)
	}

	consumers := []Consumer{}

	for _, consumerDescription := range consumerGroups.ConsumerGroupDescriptions {
		consumer := Consumer{
			consumerDescription.GroupID,
			"",
			consumerDescription.State.String(),
			nil,
		}

		for _, member := range consumerDescription.Members {
			ctp := []ConsumerTopicPartition{}

			for _, topicParts := range member.Assignment.TopicPartitions {
				fmt.Println("Topic Metadata", topicParts.Metadata)
				ctp = append(ctp, ConsumerTopicPartition{*topicParts.Topic, int64(topicParts.Offset), topicParts.Partition})
			}

			consumer.ConsumerId = member.ConsumerID
			consumer.TopicPartitions = ctp
		}
		consumers = append(consumers, consumer)
	}

	return consumers, nil
}

// func (s *Service) PublishMessage(topic string, key string, message string, channel chan kafka.Event) error {
// 	kafkaMsg := kafka.Message{
// 		TopicPartition: kafka.TopicPartition{
// 			Topic:     &topic,
// 			Partition: kafka.PartitionAny,
// 		},
// 		Key:   []byte(key),
// 		Value: []byte(message),
// 	}
// 	fmt.Println("Publishing message to %s:", topic)
// 	s.producer.Produce(&kafkaMsg, channel)
// 	return nil
// }

func (s *Service) FetchMessages(topic string, channel chan string) error {

	s.consumer.SubscribeTopics([]string{topic}, nil)
	defer s.consumer.Close()

	running := true

	for running {
		select {
		case receivedMsg := <-channel:
			fmt.Println(receivedMsg)
			running = false
		default:
		}
		msg, err := s.consumer.ReadMessage(time.Second)

		if err == nil {
			channel <- msg.String()
			fmt.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))
		} else if !err.(kafka.Error).IsTimeout() {
			// The client will automatically try to recover from all errors.
			// Timeout is not considered an error because it is raised by
			// ReadMessage in absence of messages.
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
		}
	}

	return nil
}

func (s *Service) GetTopicMetadata(topic string) (kafka.TopicMetadata, error) {
	result, err := s.client.GetMetadata(&topic, false, 5000)
	if err != nil {
		panic(err)
	}
	fmt.Println("GetTopicMetadata.getTopicMetadataResult", result)
	return result.Topics[topic], nil
}

func (s *Service) ResetConsumerOffsets(group string, topic string, offset int64) error {
	topicMetadata, err := s.GetTopicMetadata(topic)
	if err != nil {
		panic(err)
	}

	partitionArg := []kafka.TopicPartition{}
	for _, partition := range topicMetadata.Partitions {
		partitionArg = append(partitionArg, kafka.TopicPartition{
			Topic:     &topic,
			Partition: partition.ID,
			Offset:    kafka.Offset(offset),
		})
	}
	fmt.Println("ResetConsumerOffsets.partitionArg", partitionArg)
	result, err := s.client.AlterConsumerGroupOffsets(context.Background(), []kafka.ConsumerGroupTopicPartitions{
		{
			Group:      group,
			Partitions: partitionArg,
		},
	})
	if err != nil {
		panic(err)
	}

	for _, res := range result.ConsumerGroupsTopicPartitions {
		for _, resPartion := range res.Partitions {
			if resPartion.Error != nil {
				panic(resPartion.Error)
			}
		}
	}
	fmt.Println("ResetConsumerOffsets.AlterConsumerGroupOffsets result", result)
	return nil
}
