package kafkastream

import (
	"context"
	"io/ioutil"
	"log"
	"time"

	"os"

	"github.com/segmentio/kafka-go"
	"gopkg.in/yaml.v2"
)

//struct for the yaml file
type Config struct {
	Kafka_topic  string `yaml:"kafka_topic"`
	Kafka_broker string `yaml:"kafka_broker"`
}

// Method to read the struct and return the struct data with topic and broker.
func (c *Config) GetConf(yamlfile string) *Config {
	//yamlFile, err := ioutil.ReadFile("../config/data.yaml")
	yamlFile, err := ioutil.ReadFile(yamlfile)
	if err != nil {
		log.Printf("Error reading yaml file   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	return c
}

//Function that passes in the kafka broker address, topic and then the key/value pair from the main function
//This is what actually puts data into kafka.
//func StreamToKafka(brokerAddress string, topic string, Kafka_Key string, Kafka_value string) {
func StreamToKafka(brokerAddress, topic, Kafka_Key, Kafka_value string) {
	l := log.New(os.Stdout, "kafka writer: ", 0)
	ctx := context.Background()
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{brokerAddress},
		Topic:   topic,
		Logger:  l,
	})

	err := w.WriteMessages(ctx, kafka.Message{
		Key: []byte(Kafka_Key),
		// create an arbitrary message payload for the value
		Value: []byte(Kafka_value),
	})
	if err != nil {
		panic("could not write message " + err.Error())
	}
	time.Sleep(time.Second)
}
