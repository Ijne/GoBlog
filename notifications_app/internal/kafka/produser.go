package kafka

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/IBM/sarama"
	"github.com/joho/godotenv"
)

type Producer struct {
	producer sarama.AsyncProducer
}

func NewProducer(brokers []string) (*Producer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true

	producer, err := sarama.NewAsyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	go func() {
		for {
			select {
			case success, ok := <-producer.Successes():
				if !ok {
					return
				}
				log.Printf("Отправлено: topic=%s, partition=%d", success.Topic, success.Partition)
			case err := <-producer.Errors():
				log.Printf("Ошибка: %v", err)
			}
		}
	}()

	return &Producer{producer: producer}, nil
}

func (p *Producer) Send(topic string, key, value []byte) {
	p.producer.Input() <- &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.ByteEncoder(key),
		Value: sarama.ByteEncoder(value),
	}
	log.Println("Message sent", string(value))
}

func (p *Producer) Close() error {
	return p.producer.Close()
}

func SendMessage(body map[string]interface{}) {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	KAFKA_HOST := os.Getenv("KAFKA_HOST")

	log.Println("SendMessage")
	producer, err := NewProducer([]string{fmt.Sprintf("%s:9092", KAFKA_HOST)})
	if err != nil {
		panic(err)
	}
	defer producer.Close()

	eventJSON, _ := json.Marshal(body)
	producer.Send("notifications-out", nil, eventJSON)
}
