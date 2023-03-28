package main

import (
	"github.com/scyhhe/djafka/internal/djafka"
)

func main() {
	service, err := djafka.NewService()
	if err != nil {
		panic(err)
	}
	defer service.Close()

	topicName, err := service.CreateTopic("tschusch")
	if err != nil {
		panic(err)
	}

	err = service.PublishMessage(topicName, "1", "Vasko", nil)
	if err != nil {
		panic(err)
	}

	fmt.Println("consumer groups:", topicName)

	topics, err := service.ListTopics()
	if err != nil {
		panic(err)
	}

	fmt.Println("consumer groups:", topics)

	consumerGroups, err := service.ListConsumerGroups()
	if err != nil {
		panic(err)
	}

	fmt.Println("consumer groups:", consumerGroups)

	consumers, err := service.ListConsumers(consumerGroups)
	if err != nil {
		panic(err)
	}

	fmt.Println("consumers:", consumers)

	messageChan := make(chan string)

	err = service.FetchMessages(topicName, messageChan)

	allMessages := []string{}
	for msg := range messageChan {
		allMessages = append(allMessages, msg)
		messageChan <- "STOP"
		break
	}

	fmt.Println("Messages:", allMessages)

	// if err != nil {
	// 	messageChan <- "STOP"
	// }


	djafka.Run()
}
