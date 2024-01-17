package main

import (
	"log"
	"os"

	"github.com/IBM/sarama"
)

func main() {
	var config *sarama.Config

	client, err := sarama.NewClient([]string{"kafka:9092"}, config)
	if err != nil {
		panic(err)
	} else {
		os.Stderr.WriteString("> connected\n")
	}
	defer client.Close()

	consumer, err := sarama.NewConsumer([]string{"kafka:9092"}, config)
	if err != nil {
		panic(err)
	}
	log.Println("consumer created")
	defer func() {
		if err := consumer.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	log.Println("commence consuming")
	partitionConsumer, err := consumer.ConsumePartition("TEST", 0, sarama.OffsetOldest)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := partitionConsumer.Close(); err != nil {
			log.Fatalln(err)
		}
	}()
}
