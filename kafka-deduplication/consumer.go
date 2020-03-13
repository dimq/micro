package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/Shopify/sarama"
	"github.com/dimq/micro/models"
	"github.com/hashicorp/go-retryablehttp"
)

const baseURL = "http://localhost:8080"

func Consume() {
	Info.Println("starting consumer")

	config := sarama.NewConfig()

	config.Net.TLS.Enable = true
	config.Net.SASL.Enable = true
	config.Net.SASL.User = username
	config.Net.SASL.Password = password
	config.Version = sarama.V0_10_2_0
	config.ClientID = username

	config.Consumer.Offsets.Initial = sarama.OffsetOldest

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
		msg, err := ParseMessage(message.Value)
		if err != nil {
			Error.Println(err)
			continue
		}

		Info.Println(fmt.Sprintf("%+v\n", msg))
		md5Msg := GenerateHash(string(message.Value))

		firstInsert, errUpsert := models.HashUpsert(db, md5Msg)
		if errUpsert != nil {
			Error.Println(errUpsert)
		}
		if firstInsert {
			pushMessage(string(message.Value))
		}
		fmt.Println(firstInsert)

		session.MarkMessage(message, "")
	}

	return nil
}

func ParseMessage(msg []byte) (Message, error) {
	var (
		m                  Message
		invalidFormatError = errors.New("invalid format")
	)

	if err := json.Unmarshal([]byte(msg), &m); err != nil {
		return m, invalidFormatError
	}

	if m.ID == 0 {
		return Message{}, invalidFormatError
	}

	if m.Code == "" {
		return Message{}, invalidFormatError
	}

	if m.Message == "" {
		return Message{}, invalidFormatError
	}

	return m, nil
}

func GenerateHash(text string) string {
	hasher := sha256.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func pushMessage(msg string) {
	req, err := http.NewRequest("POST", baseURL+"/message", strings.NewReader(msg))
	if err != nil {
		fmt.Println(err)
	}

	retryReq, err := retryablehttp.FromRequest(req)
	if err != nil {
		fmt.Println(err)
	}

	client := retryablehttp.NewClient()
	client.Backoff = retryablehttp.LinearJitterBackoff
	client.RetryWaitMin = 800 * time.Millisecond
	client.RetryWaitMax = 1200 * time.Millisecond
	client.RetryMax = 4
	client.ErrorHandler = retryablehttp.PassthroughErrorHandler

	_, err = client.Do(retryReq)
	if err != nil {
		fmt.Println(err)
	}
}
