package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/scyhhe/djafka/internal/djafka"
)

func main() {
	service, err := djafka.NewService(djafka.Connection{Name: "test", BootstrapServer: "localhost"})
	if err != nil {
		panic(err)
	}
	defer service.Close()

	topic, err := service.CreateTopic("tschusch", 1)
	if err != nil {
		panic(err)
	}

	// Publish test messages.
	go func() {
		i := 0
		for {
			err = service.PublishMessage(topic.Name, "1", "Vasko-"+strconv.Itoa(i), nil)
			if err != nil {
				panic(err)
			}
			i++
			time.Sleep(time.Second)
		}
	}()

	fmt.Println("consumer groups:", topic)

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

	control, output, err := service.FetchMessages(topic.Name)
	if err != nil {
		panic(err)
	}

	for i := 10; i > 0; i-- {
		msg := <-output
		fmt.Printf("Received message: %s\n", string(msg.Value))
	}
	fmt.Println("end")
	close(control)

	// djafka.Run()
}
