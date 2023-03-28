package main

import (
	"fmt"

	"github.com/scyhhe/djafka/internal/djafka"
)

func testResetOffsets(service *djafka.Service) {
	// kafka-consumer-groups --bootstrap-server broker:29092 --group consumer-group-3 --describe
	// kafka-consumer-groups --bootstrap-server broker:29092 --group consumer-group-3 --reset-offsets --to-earliest --all-topics --execute
	consumerGroups, err := service.ListConsumerGroups()
	if err != nil {
		panic(err)
	}

	fmt.Println("consumer groups:", consumerGroups)

	consumers, err := service.ListConsumers(consumerGroups[0:1])
	if err != nil {
		panic(err)
	}
	fmt.Println("target consumers", consumers)

	partition := consumers[0].TopicPartitions

	fmt.Println("target partition", partition)

	// for _, p := range consumers {
	// 	partition = append(partition, p.TopicPartitions)
	// }
	service.ResetConsumerOffsets(consumerGroups[0], 0, partition)
}
func main() {
	config, err := djafka.ReadConfig()
	service, err := djafka.NewService(config.Connections[0])
	if err != nil {
		panic(err)
	}
	defer service.Close()

	// topicName, err := service.CreateTopic("tschusch")
	// if err != nil {
	// 	panic(err)
	// }

	// err = service.PublishMessage(topicName, "1", "Vasko", nil)
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println("consumer groups:", topicName)

	// topics, err := service.ListTopics()
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println("consumer groups:", topics)

	// consumerGroups, err := service.ListConsumerGroups()
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println("consumer groups:", consumerGroups)

	// consumers, err := service.ListConsumers(consumerGroups)
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println("consumers:", consumers)

	// messageChan := make(chan string)

	// err = service.FetchMessages(topicName, messageChan)

	// allMessages := []string{}
	// for msg := range messageChan {
	// 	allMessages = append(allMessages, msg)
	// 	messageChan <- "STOP"
	// 	break
	// }

	// fmt.Println("Messages:", allMessages)

	// if err != nil {
	// 	messageChan <- "STOP"
	// }

	// djafka.Run()

	testResetOffsets(service)
}
