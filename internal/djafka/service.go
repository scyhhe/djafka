package djafka

import (
	"context"
	"fmt"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

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
	producer *kafka.Producer
}

func NewService() (*Service, error) {
	client, err := kafka.NewAdminClient(&kafka.ConfigMap{
		"bootstrap.servers": "localhost",
	})

	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost",
		"group.id":          "testis",
	})

	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost",
	})

	return &Service{client, consumer, producer}, err
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

func (s *Service) CreateTopic(name string) (string, error) {
	topicSpec := kafka.TopicSpecification{Topic: name, NumPartitions: 1, ReplicationFactor: 1, Config: map[string]string{}}
	_, err := s.client.CreateTopics(context.Background(), []kafka.TopicSpecification{topicSpec})

	if err != nil {
		return err.Error(), fmt.Errorf("Failed to create new topic: %w", err)
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
				ctp = append(ctp, ConsumerTopicPartition{*topicParts.Topic, int64(topicParts.Offset), topicParts.Partition})
			}

			consumer.ConsumerId = member.ConsumerID
			consumer.TopicPartitions = ctp
		}
		consumers = append(consumers, consumer)
	}

	return consumers, nil
}

func (s *Service) FetchMessages(topic string) ([]string, error) {

	s.consumer.SubscribeTopics([]string{topic, "^aRegex.*[Tt]opic"}, nil)

	for true {
		msg, err := s.consumer.ReadMessage(time.Second)
		if err == nil {
			fmt.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))
		} else if !err.(kafka.Error).IsTimeout() {
			// The client will automatically try to recover from all errors.
			// Timeout is not considered an error because it is raised by
			// ReadMessage in absence of messages.
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
		}
	}

	s.consumer.Close()
	return nil, nil

}
