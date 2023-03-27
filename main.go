package main

import (
	"context"
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

	ctx := context.Background()
	res, err := client.
	if err != nil {
		panic(err)
	}

	fmt.Println(res)
}
