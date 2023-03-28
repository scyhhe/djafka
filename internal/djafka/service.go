package djafka

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type Connection struct {
	Name            string `json:"name"`
	BootstrapServer string `json:"bootstrapServer"`
}

type Config struct {
	Connections []Connection `json:"connections"`
}

func (c *Config) FindConnection(name string) (string, error) {
	for _, conn := range c.Connections {
		if conn.Name == name {
			return conn.BootstrapServer, nil
		}
	}

	return "", fmt.Errorf("Failed to find connection for name '%s'.", name)
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

type Service struct {
	client *kafka.AdminClient
}

func NewService() (*Service, error) {
	client, err := kafka.NewAdminClient(&kafka.ConfigMap{
		"bootstrap.servers": "localhost",
	})

	return &Service{client}, err
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
