package main

import (
	"log"
	"os"

	"github.com/IBM/sarama"
)

var (
	kafkaBrokers = []string{os.Getenv("SUBSCRIPTION_HOST")}
	KafkaTopic   = os.Getenv("SUBSCRIPTION_TOPIC")
)

func main() {
	config := sarama.NewConfig()
	// config.Consumer.Fetch.Min = 1
	// config.Consumer.Fetch.Default = 1024 * 1024
	// config.Consumer.Retry.Backoff = 2 * time.Second
	// config.Consumer.MaxWaitTime = 500 * time.Millisecond
	// config.Consumer.MaxProcessingTime = 100 * time.Millisecond
	// config.Consumer.Return.Errors = false
	// config.Consumer.Offsets.AutoCommit.Enable = true
	// config.Consumer.Offsets.AutoCommit.Interval = 1 * time.Second
	// config.Consumer.Offsets.Initial = sarama.OffsetNewest
	// config.Consumer.Offsets.Retry.Max = 3

	client, err := sarama.NewClient(kafkaBrokers, config)
	if err != nil {
		panic(err)
	} else {
		os.Stderr.WriteString("> connected\n")
	}
	defer client.Close()

	consumer, err := sarama.NewConsumer(kafkaBrokers, config)
	if err != nil {
		panic(err)
	}
	log.Println("consumer created")
	defer func() {
		if err := consumer.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	partitionConsumer, err := consumer.ConsumePartition(KafkaTopic, 0, sarama.OffsetOldest)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := partitionConsumer.Close(); err != nil {
			log.Fatalln(err)
		}
	}()
}
