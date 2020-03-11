package main

import (
	"flag"
	"os"
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
