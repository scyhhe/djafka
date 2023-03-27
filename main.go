package main

import (
	"fmt"

	"github.com/scyhhe/djafka/internal/djafka"
)

func main() {
	service, err := djafka.NewService()
	if err != nil {
		panic(err)
	}
	defer service.Close()

	topics, err := service.ListTopics()
	if err != nil {
		panic(err)
	}

	for _, topic := range topics {
		fmt.Println(topic)
	}
}
