package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/Shopify/sarama"
)

var (
	brokers  = ""
	topics   = ""
	username = ""
	password = ""
	group    = ""
	version  = ""
)

func init() {
	flag.StringVar(&brokers, "brokers", LookupEnvOrString("KAFKA_BROKERS", brokers), "Kafka bootstrap brokers to connect to, as a comma separated list")
	flag.StringVar(&topics, "topics", LookupEnvOrString("KAFKA_TOPICS", topics), "Kafka topics to be consumed, as a comma seperated list")
	flag.StringVar(&username, "username", LookupEnvOrString("KAFKA_USER", username), "Kafka username to use")
	flag.StringVar(&password, "password", LookupEnvOrString("KAFKA_PASS", password), "Kafka password to use")
	flag.StringVar(&group, "group", LookupEnvOrString("KAFKA_GROUP", group), "Kafka consumer group definition")
	flag.StringVar(&version, "version", "2.1.1", "Kafka cluster version")
	flag.Parse()

	if len(brokers) == 0 {
		panic("no Kafka bootstrap brokers defined, please set the -brokers flag")
	}

	if len(topics) == 0 {
		panic("no topics given to be consumed, please set the -topics flag")
	}

	if len(group) == 0 {
		panic("no Kafka consumer group defined, please set the -group flag")
	}
}

func LookupEnvOrString(key string, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultVal
}

func main() {
	Init(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)

	Info.Println("starting consumer")

	config := sarama.NewConfig()

	config.Net.TLS.Enable = true
	config.Net.SASL.Enable = true
	config.Net.SASL.User = username
	config.Net.SASL.Password = password
	config.Version = sarama.V0_10_2_0
	config.ClientID = username

	//config.Consumer.Offsets.Initial = sarama.OffsetOldest

	consumer := Consumer{
		ready: make(chan bool),
	}

	ctx, cancel := context.WithCancel(context.Background())
	client, err := sarama.NewConsumerGroup(strings.Split(brokers, ","), group, config)
	if err != nil {
		log.Panicf("Error creating consumer group client: %v", err)
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			// `Consume` should be called inside an infinite loop, when a
			// server-side rebalance happens, the consumer session will need to be
			// recreated to get the new claims
			if err := client.Consume(ctx, strings.Split(topics, ","), &consumer); err != nil {
				log.Panicf("Error from consumer: %v", err)
			}
			// check if context was cancelled, signaling that the consumer should stop
			if ctx.Err() != nil {
				return
			}
			consumer.ready = make(chan bool)
		}
	}()

	<-consumer.ready // Await till the consumer has been set up
	Info.Println("Sarama consumer up and running!...")

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-ctx.Done():
		Info.Println("terminating: context cancelled")
	case <-sigterm:
		Info.Println("terminating: via signal")
	}
	cancel()
	wg.Wait()
	if err = client.Close(); err != nil {
		log.Panicf("Error closing client: %v", err)
	}

}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (consumer *Consumer) Setup(sarama.ConsumerGroupSession) error {
	// Mark the consumer as ready
	close(consumer.ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (consumer *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (consumer *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {

	for message := range claim.Messages() {
		msg, err := ParseMessage(message)
		if err != nil {
			Error.Println(err)
		} else {
			Info.Println(msg)
		}
		session.MarkMessage(message, "")
	}

	return nil
}

func ParseMessage(msg *sarama.ConsumerMessage) (Message, error) {
	var m Message
	if err := json.Unmarshal(msg.Value, &m); err != nil {
		return Message{}, errors.New("invalid format")
	}
	return m, nil
}
