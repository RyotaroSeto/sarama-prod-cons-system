package main

import (
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/IBM/sarama"
)

var (
	kafkaBrokers = []string{os.Getenv("SUBSCRIPTION_HOST")}
	KafkaTopic   = os.Getenv("SUBSCRIPTION_TOPIC")
	enqueued     int
)

func main() {
	producer, err := setupAsyncProducer()
	if err != nil {
		panic(err)
	} else {
		log.Println("Kafka AsyncProducer up and running!")
	}

	// Trap SIGINT to trigger a graceful shutdown.
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	produceMessages(producer, signals)

	log.Printf("Kafka AsyncProducer finished with %d messages produced.", enqueued)
}

// setupAsyncProducer will create a AsyncProducer and returns it
func setupAsyncProducer() (sarama.AsyncProducer, error) {
	config := sarama.NewConfig()
	return sarama.NewAsyncProducer(kafkaBrokers, config)
}

// setupSyncProducer will create a SyncProducer and returns it
func setupSyncProducer() (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	return sarama.NewSyncProducer(kafkaBrokers, config)
}

// produceMessages will send 'testing 123' to KafkaTopic each second, until receive a os signal to stop e.g. control + c
// by the user in terminal
func produceMessages(producer sarama.AsyncProducer, signals chan os.Signal) {
	for {
		time.Sleep(time.Second)
		message := &sarama.ProducerMessage{Topic: KafkaTopic, Value: sarama.StringEncoder("testing 123")}
		select {
		case producer.Input() <- message:
			enqueued++
			log.Println("New Message produced")
		case <-signals:
			producer.AsyncClose() // Trigger a shutdown of the producer.
			return
		}
	}
}
