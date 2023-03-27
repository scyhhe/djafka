package main

import (
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

func main() {
	client, err := kafka.NewAdminClient(&kafka.ConfigMap{
		"bootstrap.servers": "localhost",
	})
	if err != nil {
		panic(err)
	}
	defer client.Close()

	res, err := client.GetMetadata(nil, true, 5000)
	if err != nil {
		panic(err)
	}

	fmt.Println(res)
}
