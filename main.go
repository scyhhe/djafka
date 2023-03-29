package main

import (
	"fmt"

	"github.com/scyhhe/djafka/internal/djafka"
)

func main() {
	service, err := djafka.NewService(djafka.Connection{Name: "test", BootstrapServer: "localhost"})
	if err != nil {
		panic(err)
	}
	defer service.Close()

	topicName, err := service.CreateTopic("tschusch", 1)
	if err != nil {
		panic(err)
	}

	err = service.PublishMessage(topicName.Name, "1", "Vasko", nil)
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

	_, err = service.FetchMessages(topicName.Name)
	if err != nil {
		panic(err)
	}

	allMessages := []string{}
	// for msg := range messageChan {
	// 	allMessages = append(allMessages, msg)
	// 	messageChan <- "STOP"
	// 	break
	// }

	fmt.Println("Messages:", allMessages)

	// if err != nil {
	// 	messageChan <- "STOP"
	// }

	djafka.Run()
}
